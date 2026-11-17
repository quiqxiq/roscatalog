package rosrun

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/quiqxiq/goros/v4"
	"github.com/quiqxiq/goros/v4/proto"
)

// cmd adalah satu kalimat dari klien.
type cmd struct {
	path  string
	attrs map[string]string
	query []string
	tag   string
}

func (c cmd) has(k string) bool { _, ok := c.attrs[k]; return ok }

// mock adalah RouterOS API palsu: protokol kawat asli, perilaku buatan.
type mock struct {
	t  *testing.T
	ln net.Listener

	mu    sync.Mutex
	conns []net.Conn
	stops map[string]chan struct{} // tag -> sinyal stop
	calls []cmd
	wg    sync.WaitGroup

	rejectProplist map[string]bool // path -> tolak bila ada .proplist
	rejectPaths    map[string]bool // path -> selalu ditolak (mis. listen tidak didukung)
	printCalls     int32           // jumlah print one-shot hotspot-active
	streamsActive  int32
}

func newMock(t *testing.T) *mock {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	m := &mock{t: t, ln: ln, stops: map[string]chan struct{}{}, rejectProplist: map[string]bool{}, rejectPaths: map[string]bool{}}
	go m.accept()
	t.Cleanup(m.close)
	return m
}

func (m *mock) close() {
	m.ln.Close()
	m.mu.Lock()
	for _, c := range m.conns {
		c.Close()
	}
	m.mu.Unlock()
}

// client membuka koneksi goros/v4 sungguhan ke mock dan membungkusnya.
func (m *mock) client() (*routeros.Client, *Runner) {
	m.t.Helper()
	c, err := routeros.Dial(m.ln.Addr().String(), "admin", "")
	if err != nil {
		m.t.Fatal(err)
	}
	m.t.Cleanup(func() { c.Close() })
	r := New(c)
	r.Logf = func(f string, a ...any) { m.t.Logf(f, a...) }
	return c, r
}

func (m *mock) accept() {
	for {
		conn, err := m.ln.Accept()
		if err != nil {
			return
		}
		m.mu.Lock()
		m.conns = append(m.conns, conn)
		m.mu.Unlock()
		go m.serve(conn)
	}
}

// --- pembaca kalimat klien (proto.Reader menolak kata query "?") -----------

func readLen(r *bufio.Reader) (int, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	switch {
	case b&0x80 == 0:
		return int(b), nil
	case b&0xC0 == 0x80:
		b2, err := r.ReadByte()
		return int(b&^0xC0)<<8 | int(b2), err
	case b&0xE0 == 0xC0:
		var x [2]byte
		_, err := io.ReadFull(r, x[:])
		return int(b&^0xE0)<<16 | int(x[0])<<8 | int(x[1]), err
	}
	return 0, fmt.Errorf("panjang kata tidak didukung mock: %#x", b)
}

func readSentence(r *bufio.Reader) (cmd, error) {
	c := cmd{attrs: map[string]string{}}
	first := true
	for {
		n, err := readLen(r)
		if err != nil {
			return c, err
		}
		if n == 0 {
			return c, nil
		}
		b := make([]byte, n)
		if _, err := io.ReadFull(r, b); err != nil {
			return c, err
		}
		w := string(b)
		switch {
		case first:
			c.path, first = w, false
		case strings.HasPrefix(w, ".tag="):
			c.tag = w[5:]
		case strings.HasPrefix(w, "?"):
			c.query = append(c.query, w)
		case strings.HasPrefix(w, "="):
			kv := strings.SplitN(w[1:], "=", 2)
			if len(kv) == 1 {
				kv = append(kv, "")
			}
			c.attrs[kv[0]] = kv[1]
		}
	}
}

// --- penulisan balasan ------------------------------------------------------

func (m *mock) send(w proto.Writer, tag, word string, kv map[string]string) {
	w.BeginSentence()
	w.WriteWord(word)
	keys := make([]string, 0, len(kv))
	for k := range kv {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		w.WriteWord("=" + k + "=" + kv[k])
	}
	if tag != "" {
		w.WriteWord(".tag=" + tag)
	}
	w.EndSentence()
}

// project meniru =.proplist=: sisakan hanya properti yang diminta.
func project(c cmd, row map[string]string) map[string]string {
	pl, ok := c.attrs[".proplist"]
	if !ok {
		return row
	}
	want := map[string]bool{".dead": true}
	for _, p := range strings.Split(pl, ",") {
		want[p] = true
	}
	out := map[string]string{}
	for k, v := range row {
		if want[k] {
			out[k] = v
		}
	}
	return out
}

func (m *mock) serve(conn net.Conn) {
	rd := bufio.NewReader(conn)
	w := proto.NewWriter(conn)
	for {
		c, err := readSentence(rd)
		if err != nil {
			return
		}
		m.mu.Lock()
		m.calls = append(m.calls, c)
		m.mu.Unlock()

		switch c.path {
		case "/login":
			m.send(w, c.tag, "!done", nil)
		case "/cancel":
			m.mu.Lock()
			stop := m.stops[c.attrs["tag"]]
			delete(m.stops, c.attrs["tag"])
			m.mu.Unlock()
			if stop != nil {
				close(stop)
			}
			m.send(w, c.tag, "!done", nil)
		default:
			m.wg.Add(1)
			go func() {
				defer m.wg.Done()
				m.handle(w, c)
			}()
		}
	}
}

// stream menjalankan fn tiap interval sampai /cancel, lalu membalas seperti router asli.
func (m *mock) stream(w proto.Writer, c cmd, every time.Duration, fn func(i int)) {
	stop := make(chan struct{})
	m.mu.Lock()
	m.stops[c.tag] = stop
	m.mu.Unlock()
	atomic.AddInt32(&m.streamsActive, 1)
	defer atomic.AddInt32(&m.streamsActive, -1)

	t := time.NewTicker(every)
	defer t.Stop()
	for i := 0; ; i++ {
		fn(i)
		select {
		case <-stop:
			m.send(w, c.tag, "!trap", map[string]string{"category": "2", "message": "interrupted"})
			m.send(w, c.tag, "!done", nil)
			return
		case <-t.C:
		}
	}
}

func (m *mock) waitStreamsStopped(timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if atomic.LoadInt32(&m.streamsActive) == 0 {
			return true
		}
		time.Sleep(5 * time.Millisecond)
	}
	return false
}

func (m *mock) callsTo(path string) []cmd {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []cmd
	for _, c := range m.calls {
		if c.path == path {
			out = append(out, c)
		}
	}
	return out
}

func hotspotRow(id, user string, n int) map[string]string {
	return map[string]string{
		".id": id, "user": user, "address": "10.0.0." + id[1:], "mac-address": "AA:BB:CC:00:00:0" + id[1:],
		"server": "hs1", "login-by": "http-chap", "uptime": fmt.Sprintf("%ds", n),
		"session-time-left": "1h", "idle-time": "0s", "bytes-in": fmt.Sprint(n * 100), "bytes-out": "5",
		"comment": "tidak diminta",
	}
}

func (m *mock) handle(w proto.Writer, c cmd) {
	if m.rejectPaths[c.path] {
		m.send(w, c.tag, "!trap", map[string]string{"message": "not supported on this menu"})
		m.send(w, c.tag, "!done", nil)
		return
	}
	switch c.path {
	case "/ip/hotspot/user/print":
		if c.has("count-only") {
			m.send(w, c.tag, "!done", map[string]string{"ret": "7"})
			return
		}
		m.send(w, c.tag, "!done", nil)

	case "/ip/hotspot/user/add":
		m.send(w, c.tag, "!done", map[string]string{"ret": "*A1"})

	case "/ip/hotspot/active/print":
		if m.rejectProplist[c.path] && c.has(".proplist") {
			m.send(w, c.tag, "!trap", map[string]string{"message": "unknown parameter"})
			m.send(w, c.tag, "!done", nil)
			return
		}
		if iv := c.attrs["interval"]; iv != "" {
			d, err := time.ParseDuration(iv)
			if err != nil {
				m.send(w, c.tag, "!trap", map[string]string{"message": "bad interval " + iv})
				return
			}
			m.stream(w, c, d, func(i int) {
				m.send(w, c.tag, "!re", project(c, hotspotRow("*1", "alice", i)))
				m.send(w, c.tag, "!re", project(c, hotspotRow("*2", "bob", i)))
			})
			return
		}
		n := int(atomic.AddInt32(&m.printCalls, 1))
		m.send(w, c.tag, "!re", project(c, hotspotRow("*1", "alice", n)))
		if n < 3 {
			m.send(w, c.tag, "!re", project(c, hotspotRow("*2", "bob", 1)))
		}
		if n >= 4 {
			m.send(w, c.tag, "!re", project(c, hotspotRow("*3", "carol", 1)))
		}
		m.send(w, c.tag, "!done", nil)

	case "/ip/hotspot/active/listen":
		stop := make(chan struct{})
		m.mu.Lock()
		m.stops[c.tag] = stop
		m.mu.Unlock()
		atomic.AddInt32(&m.streamsActive, 1)
		defer atomic.AddInt32(&m.streamsActive, -1)

		script := []map[string]string{
			hotspotRow("*3", "carol", 1),   // item baru
			{".id": "*1", "uptime": "99s"}, // perubahan parsial
			{".id": "*2", ".dead": "true"}, // item hilang
		}
		for _, row := range script {
			select {
			case <-time.After(60 * time.Millisecond):
				m.send(w, c.tag, "!re", project(c, row))
			case <-stop:
				goto cancelled
			}
		}
		<-stop
	cancelled:
		m.send(w, c.tag, "!trap", map[string]string{"category": "2", "message": "interrupted"})
		m.send(w, c.tag, "!done", nil)

	case "/interface/monitor-traffic":
		row := map[string]string{"name": c.attrs["interface"], "rx-bits-per-second": "1000", "tx-bits-per-second": "2000", "rx-drops-per-second": "0"}
		if c.has("once") {
			m.send(w, c.tag, "!re", project(c, row))
			m.send(w, c.tag, "!done", nil)
			return
		}
		m.stream(w, c, 15*time.Millisecond, func(i int) {
			m.send(w, c.tag, "!re", project(c, row))
		})

	case "/flood":
		// banjir cepat untuk menguji konsumen lambat
		m.stream(w, c, time.Millisecond, func(i int) {
			m.send(w, c.tag, "!re", map[string]string{"seq": fmt.Sprint(i)})
		})

	case "/log/print":
		if m.rejectProplist[c.path] && c.has(".proplist") {
			m.send(w, c.tag, "!trap", map[string]string{"message": "unknown parameter"})
			m.send(w, c.tag, "!done", nil)
			return
		}
		stop := make(chan struct{})
		m.mu.Lock()
		m.stops[c.tag] = stop
		m.mu.Unlock()
		atomic.AddInt32(&m.streamsActive, 1)
		defer atomic.AddInt32(&m.streamsActive, -1)
		for i := 1; i <= 3; i++ {
			m.send(w, c.tag, "!re", project(c, map[string]string{
				".id": fmt.Sprintf("*%d", i), "time": "10:00:0" + fmt.Sprint(i), "topics": "system,info", "message": fmt.Sprintf("baris %d", i),
			}))
			time.Sleep(10 * time.Millisecond)
		}
		<-stop
		m.send(w, c.tag, "!trap", map[string]string{"category": "2", "message": "interrupted"})
		m.send(w, c.tag, "!done", nil)

	default:
		m.send(w, c.tag, "!trap", map[string]string{"message": "no such command " + c.path})
		m.send(w, c.tag, "!done", nil)
	}
}
