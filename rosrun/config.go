package rosrun

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	ros "example.com/roscatalog"
)

const (
	// Port bawaan API RouterOS.
	defaultAPIPort    = 8728
	defaultAPISSLPort = 8729

	// DefaultDialTimeout membatasi dial+login bila Config.Timeout kosong.
	DefaultDialTimeout = 10 * time.Second
	// DefaultPollEvery adalah jeda poll bawaan bila Config.PollEvery kosong.
	DefaultPollEvery = time.Second
	// DefaultQueueDepth adalah nilai client.Queue bawaan.
	DefaultQueueDepth = 1000
	// DefaultMaxBuffered membatasi antrean per stream.
	DefaultMaxBuffered = 8192
)

// Gate mengapit satu command yang benar-benar dikirim ke router. release
// dipanggil setelah command selesai (termasuk saat stream berakhir). Dipakai
// pemakai untuk limiter per router, audit, atau metrik.
//
// path adalah command API, mis. "/ip/hotspot/active/print".
type Gate func(ctx context.Context, path string) (release func(), err error)

// Config adalah parameter satu router: alamat, kredensial, dan perilaku
// koneksi. Satu Config selalu membuka dua koneksi (command dan feed stream).
type Config struct {
	// ID hanya dipakai untuk log. Boleh kosong.
	ID string

	// Address boleh "host" atau "host:port". Bila Port kosong, dipakai
	// 8728 (atau 8729 saat TLS).
	Address  string
	Port     int
	Username string
	Password string

	TLS      bool
	Insecure bool // hanya berlaku saat TLS; CHR memakai sertifikat sendiri

	// Timeout membatasi dial+login. Kosong = DefaultDialTimeout.
	Timeout time.Duration

	// Stream adalah preferensi delivery bawaan: true memakai stream, false
	// memakai poll. Plan eksplisit pada Request selalu menang. Lihat PlanFor.
	Stream bool

	// PollEvery adalah jeda poll bawaan. Kosong = DefaultPollEvery.
	PollEvery time.Duration

	// Queue adalah client.Queue (buffer sentence yang dialokasikan Listen).
	// Kosong = DefaultQueueDepth.
	Queue int

	// MaxBuffered membatasi antrean per stream. Kosong = DefaultMaxBuffered.
	MaxBuffered int

	// EmptySnapshots meneruskan Runner.EmptySnapshots.
	EmptySnapshots bool

	// Gate opsional, dipanggil sekali per command.
	Gate Gate

	// Logf opsional.
	Logf func(format string, args ...any)

	// backoff adalah jeda reconnect pertama. Tidak diekspos; test memakai
	// nilai kecil supaya cepat.
	backoff time.Duration
}

// Validate memeriksa parameter wajib. Nilai opsional yang kosong diisi
// denganDefaults saat Open.
func (c Config) Validate() error {
	if strings.TrimSpace(c.Address) == "" {
		return fmt.Errorf("rosrun: address is required")
	}
	if strings.TrimSpace(c.Username) == "" {
		return fmt.Errorf("rosrun: username is required")
	}
	if c.Port < 0 || c.Port > 65535 {
		return fmt.Errorf("rosrun: port %d out of range", c.Port)
	}
	if c.Timeout < 0 {
		return fmt.Errorf("rosrun: timeout must not be negative")
	}
	if c.PollEvery < 0 {
		return fmt.Errorf("rosrun: poll interval must not be negative")
	}
	if c.PollEvery > 0 && c.PollEvery < ros.MinPollEvery {
		return fmt.Errorf("rosrun: poll interval %s below minimum %s", c.PollEvery, ros.MinPollEvery)
	}
	if c.Queue < 0 || c.MaxBuffered < 0 {
		return fmt.Errorf("rosrun: queue and max buffered must not be negative")
	}
	return nil
}

// withDefaults mengembalikan salinan Config dengan nilai opsional terisi.
func (c Config) withDefaults() Config {
	if c.Timeout == 0 {
		c.Timeout = DefaultDialTimeout
	}
	if c.PollEvery == 0 {
		c.PollEvery = DefaultPollEvery
	}
	if c.Queue == 0 {
		c.Queue = DefaultQueueDepth
	}
	if c.MaxBuffered == 0 {
		c.MaxBuffered = DefaultMaxBuffered
	}
	return c
}

// addr mengembalikan alamat "host:port" siap dial.
func (c Config) addr() string {
	host := strings.TrimSpace(c.Address)
	if h, p, err := net.SplitHostPort(host); err == nil && h != "" && p != "" {
		return host
	}

	port := c.Port
	if port == 0 {
		if c.TLS {
			port = defaultAPISSLPort
		} else {
			port = defaultAPIPort
		}
	}
	return net.JoinHostPort(host, strconv.Itoa(port))
}
