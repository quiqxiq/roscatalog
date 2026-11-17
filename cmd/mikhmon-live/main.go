// Command mikhmon-live menguji workflow Mikhmon v4 langsung ke router:
// profile + on-login, voucher, record penjualan, dan Expire Monitor.
//
//	go run ./cmd/mikhmon-live -addr 192.168.88.1:8728 -user admin -pass secret
//	go run ./cmd/mikhmon-live -addr 192.168.88.1:8728 -user admin -pass secret -json
//
// Command ini membuat data uji berprefix zzmik lalu menghapusnya lagi. Data
// Mikhmon yang sudah ada tidak diubah, kecuali scheduler Expire Monitor yang
// belum ada (dibuat lalu dihapus) dan script laporan uji.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	routeros "github.com/quiqxiq/goros/v4"

	ros "example.com/roscatalog"
	"example.com/roscatalog/mikhmon"
	"example.com/roscatalog/rosrun"
)

type check struct {
	Name   string `json:"name"`
	OK     bool   `json:"ok"`
	Detail string `json:"detail,omitempty"`
	Error  string `json:"error,omitempty"`
}

type report struct {
	Addr    string  `json:"addr"`
	Version string  `json:"version,omitempty"`
	Date    string  `json:"date,omitempty"`
	Time    string  `json:"time,omitempty"`
	Checks  []check `json:"checks"`
	Error   string  `json:"error,omitempty"`
	Kept    bool    `json:"kept,omitempty"`
}

func (r *report) add(name string, ok bool, detail string, err error) {
	c := check{Name: name, OK: ok, Detail: detail}
	if err != nil {
		c.Error = err.Error()
	}
	r.Checks = append(r.Checks, c)
}

func (r *report) ok() bool {
	if r.Error != "" {
		return false
	}
	for _, c := range r.Checks {
		if !c.OK {
			return false
		}
	}
	return true
}

func main() {
	addr := flag.String("addr", "", "alamat API, mis. 192.168.88.1:8728")
	user := flag.String("user", "admin", "user API")
	pass := flag.String("pass", "", "password API")
	timeout := flag.Duration("timeout", 2*time.Minute, "batas waktu seluruh skenario")
	asJSON := flag.Bool("json", false, "cetak laporan JSON")
	keep := flag.Bool("keep", false, "jangan hapus data uji")
	flag.Parse()

	if *addr == "" {
		fmt.Fprintln(os.Stderr, "-addr wajib diisi")
		os.Exit(2)
	}

	rep := &report{Addr: *addr, Kept: *keep}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	err := scenario(ctx, *addr, *user, *pass, *keep, rep)
	if err != nil {
		rep.Error = err.Error()
	}

	if *asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(rep)
	} else {
		printReport(rep)
	}
	if !rep.ok() {
		os.Exit(1)
	}
}

func printReport(r *report) {
	fmt.Printf("Router  : %s\n", r.Addr)
	if r.Version != "" {
		fmt.Printf("Versi   : %s\n", r.Version)
	}
	if r.Date != "" {
		fmt.Printf("Tanggal : %s %s\n", r.Date, r.Time)
	}
	for _, c := range r.Checks {
		status := "ok  "
		if !c.OK {
			status = "FAIL"
		}
		fmt.Printf("[%s] %s", status, c.Name)
		if c.Detail != "" {
			fmt.Printf(": %s", c.Detail)
		}
		if c.Error != "" {
			fmt.Printf(": %s", c.Error)
		}
		fmt.Println()
	}
	if r.Error != "" {
		fmt.Printf("ERROR: %s\n", r.Error)
	}
	if r.Kept {
		fmt.Println("Data uji dipertahankan (-keep).")
	}
}

type cleanup struct {
	name string
	fn   func() error
}

type tester struct {
	ctx    context.Context
	run    *rosrun.Runner
	client *mikhmon.Client
	rep    *report
	cleans []cleanup
}

func scenario(ctx context.Context, addr, user, pass string, keep bool, rep *report) error {
	c, err := routeros.Dial(addr, user, pass)
	if err != nil {
		return fmt.Errorf("dial %s: %w", addr, err)
	}
	defer c.Close()

	run := rosrun.New(c)
	t := &tester{ctx: ctx, run: run, client: mikhmon.New(run), rep: rep}
	defer func() {
		if !keep {
			t.cleanup()
		}
	}()

	if err := t.must("koneksi/resource", func() (string, error) {
		// Ambil versi saja. Proplist system-resource yang lengkap membuat
		// RouterOS 6.49.17 menutup koneksi API (field cpu pada proplist).
		reply, err := c.RunArgs([]string{"/system/resource/print", "=.proplist=version"})
		if err != nil {
			return "", err
		}
		if reply == nil || len(reply.Re) == 0 {
			return "", errors.New("system-resource kosong")
		}
		rep.Version = reply.Re[0].Map["version"]
		return rep.Version, nil
	}); err != nil {
		return err
	}

	var clockDate, clockTime string
	if err := t.must("koneksi/clock", func() (string, error) {
		rows, err := t.once("system-clock", "list")
		if err != nil {
			return "", err
		}
		if len(rows) == 0 {
			return "", errors.New("system-clock kosong")
		}
		clockDate, clockTime = rows[0]["date"], rows[0]["time"]
		rep.Date, rep.Time = clockDate, clockTime
		return clockDate + " " + clockTime, nil
	}); err != nil {
		return err
	}

	legacy, err := legacyDate(clockDate)
	if err != nil {
		return fmt.Errorf("tanggal router %q: %w", clockDate, err)
	}
	owner := (mikhmon.Sale{Date: legacy}).Owner()
	if owner == "" {
		return fmt.Errorf("owner tidak bisa diturunkan dari %q", legacy)
	}
	today := mmddyy(legacy)
	year := yearOf(legacy)
	if year == "" {
		year = fmt.Sprintf("%d", time.Now().Year())
	}

	tag := fmt.Sprintf("zzmik%d", time.Now().Unix()%100000000)
	profileName := tag + "-p"
	userName := tag + "-u"
	userX := tag + "-x"
	userN := tag + "-n"
	onLoginName := tag + "-onlogin"
	expMonName := tag + "-expmon"
	comment := mikhmon.MikhmonComment(userName, userName, tag+"-"+today)

	var onLogin string
	if err := t.must("on-login/build", func() (string, error) {
		s, err := mikhmon.OnLoginScript(mikhmon.OnLoginOptions{
			Mode:  mikhmon.ExpiryRemoveRecord,
			Price: "5000", Validity: "1d", SellingPrice: "7000",
			Profile: profileName,
		})
		if err != nil {
			return "", err
		}
		if strings.ContainsAny(s, "\r\n") {
			return "", errors.New("on-login mengandung newline")
		}
		onLogin = s
		return fmt.Sprintf("%d byte, satu baris", len(s)), nil
	}); err != nil {
		return err
	}

	if err := t.must("profile/add", func() (string, error) {
		id, err := t.client.AddProfile(ctx, mikhmon.Profile{
			Name: profileName, RateLimit: "1M/1M", SharedUsers: "1",
			OnLogin: onLogin,
		})
		if err != nil {
			return "", err
		}
		t.onCleanup("profile "+profileName, func() error {
			return t.removeByName("hotspot-user-profile", "detail", "name", profileName)
		})
		return "id=" + id, nil
	}); err != nil {
		return err
	}

	if err := t.must("profile/on-login tersimpan", func() (string, error) {
		rows, err := t.once("hotspot-user-profile", "detail",
			ros.Query{Prop: "name", Value: profileName})
		if err != nil {
			return "", err
		}
		if len(rows) == 0 {
			return "", errors.New("profile tidak ditemukan setelah add")
		}
		got := rows[0]["on-login"]
		if got == "" {
			return "", errors.New("on-login kosong")
		}
		if strings.ContainsAny(got, "\r\n") {
			return "", errors.New("on-login tersimpan multi-baris")
		}
		if !strings.Contains(got, `:local mode "X"`) || !strings.Contains(got, "comment=mikhmon") {
			return "", fmt.Errorf("on-login tersimpan tidak memuat mode/laporan: %.120s", got)
		}
		exact := "ya"
		if got != onLogin {
			exact = "tidak (router menormalkan)"
		}
		return fmt.Sprintf("%d byte, sama persis=%s", len(got), exact), nil
	}); err != nil {
		return err
	}

	if err := t.must("voucher/add", func() (string, error) {
		id, err := t.client.AddUser(ctx, mikhmon.User{
			Name: userName, Password: userName, Profile: profileName,
			Server: "all", Comment: comment, LimitUptime: "1d",
		})
		if err != nil {
			return "", err
		}
		t.onCleanup("user "+userName, func() error {
			return t.removeByName("hotspot-user", "list", "name", userName)
		})
		return "id=" + id, nil
	}); err != nil {
		return err
	}

	if err := t.must("voucher/query comment", func() (string, error) {
		rows, err := t.client.Vouchers(ctx, comment)
		if err != nil {
			return "", err
		}
		if len(rows) != 1 {
			return "", fmt.Errorf("voucher comment=%q: %d baris, want 1", comment, len(rows))
		}
		if rows[0]["name"] != userName || rows[0]["password"] != userName {
			return "", fmt.Errorf("voucher salah: %v", rows[0])
		}
		return fmt.Sprintf("%s / %s", rows[0]["name"], rows[0]["password"]), nil
	}); err != nil {
		return err
	}

	if err := t.must("voucher/count profile", func() (string, error) {
		n, err := t.client.CountUsers(ctx, profileName)
		if err != nil {
			return "", err
		}
		if n != 1 {
			return "", fmt.Errorf("count profile=%d, want 1", n)
		}
		return "1 user", nil
	}); err != nil {
		return err
	}

	// Jaring pengaman bila on-login gagal di tengah jalan.
	t.onCleanup("scheduler sementara "+userName, func() error {
		return t.removeByName("system-scheduler", "list", "name", userName)
	})
	t.onCleanup("sales scan "+userName, func() error {
		rows, err := t.salesForUser(userName)
		if err != nil {
			return err
		}
		for _, row := range rows {
			if err := t.removeByID("mikhmon-sales", row[".id"]); err != nil {
				return err
			}
		}
		return nil
	})

	if err := t.must("on-login/simulasi login", func() (string, error) {
		sim := fmt.Sprintf(`:local user "%s"; :local address "198.51.100.10"; :local "mac-address" "AA:BB:CC:DD:EE:FF"; %s`,
			userName, onLogin)
		if _, err := t.action("system-script", "add", "", map[string]string{
			"name": onLoginName, "source": sim, "comment": tag,
		}); err != nil {
			return "", err
		}
		t.onCleanup("script "+onLoginName, func() error {
			return t.removeByName("system-script", "list", "name", onLoginName)
		})
		if _, err := t.action("system-script", "run", onLoginName, nil); err != nil {
			return "", err
		}
		time.Sleep(3 * time.Second)
		return "script dijalankan", nil
	}); err != nil {
		return err
	}

	if err := t.must("on-login/comment user", func() (string, error) {
		rows, err := t.client.Users(ctx, profileName)
		if err != nil {
			return "", err
		}
		var got string
		for _, row := range rows {
			if row["name"] == userName {
				got = row["comment"]
				break
			}
		}
		if got == "" {
			return "", errors.New("comment user kosong/tidak ditemukan")
		}
		if got == comment {
			return "", errors.New("on-login tidak mengubah comment")
		}
		if !strings.Contains(got, " X") {
			return "", fmt.Errorf("comment tidak berisi mode X: %q", got)
		}
		return got, nil
	}); err != nil {
		return err
	}

	if err := t.must("on-login/scheduler sementara bersih", func() (string, error) {
		rows, err := t.once("system-scheduler", "list",
			ros.Query{Prop: "name", Value: userName})
		if err != nil {
			return "", err
		}
		if len(rows) != 0 {
			return "", fmt.Errorf("scheduler sementara %q masih ada", userName)
		}
		return "tidak ada sisa scheduler", nil
	}); err != nil {
		return err
	}

	var sale map[string]string
	if err := t.must("sales/record", func() (string, error) {
		rows, err := t.salesForUser(userName)
		if err != nil {
			return "", err
		}
		if len(rows) == 0 {
			return "", errors.New("tidak ada script laporan untuk user ini")
		}
		sale = rows[0]
		t.onCleanup("sales "+sale["name"], func() error {
			return t.removeByID("mikhmon-sales", sale[".id"])
		})
		parsed, err := mikhmon.ParseSaleName(sale["name"])
		if err != nil {
			return "", err
		}
		if parsed.User != userName {
			return "", fmt.Errorf("user di name=%q, want %q", parsed.User, userName)
		}
		if parsed.Price != "5000" || parsed.Validity != "1d" || parsed.Profile != profileName {
			return "", fmt.Errorf("field laporan salah: %+v", parsed)
		}
		if sale["source"] != legacy {
			return "", fmt.Errorf("source=%q, want %q", sale["source"], legacy)
		}
		return fmt.Sprintf("name=%s owner=%s source=%s", sale["name"], sale["owner"], sale["source"]), nil
	}); err != nil {
		return err
	}

	if err := t.must("sales/query source (get_report)", func() (string, error) {
		rows, err := t.once("mikhmon-sales", "report",
			ros.Query{Prop: "source", Value: legacy})
		if err != nil {
			return "", err
		}
		if !containsName(rows, userName) {
			return "", fmt.Errorf("query source=%q tidak menemukan laporan", legacy)
		}
		return fmt.Sprintf("source=%s ditemukan", legacy), nil
	}); err != nil {
		return err
	}

	if err := t.must("sales/query owner (dashboard)", func() (string, error) {
		rows, err := t.client.Sales(ctx, owner)
		if err != nil {
			return "", err
		}
		if containsName(rows, userName) {
			return "owner=" + owner + " ditemukan", nil
		}
		if sale["owner"] == "*sys" {
			return "owner=*sys (belum tampil di dashboard)", nil
		}
		return "", fmt.Errorf("query owner=%q tidak menemukan laporan (owner tersimpan %q)", owner, sale["owner"])
	}); err != nil {
		return err
	}

	// Uji persis skenario RouterOS 7.24+: script laporan dibuat dengan
	// owner=*sys, lalu dinormalisasi seperti get_dashboard.php.
	if err := t.must("sales/owner *sys + normalisasi", func() (string, error) {
		sysSale := mikhmon.Sale{
			Date: legacy, Time: "00:00:00", User: tag + "-sys",
			Price: "1", Validity: "1d", Profile: profileName,
			Comment: "vc-sys",
		}
		id, err := t.action("mikhmon-sales", "add", "", map[string]string{
			"name":    sysSale.Name(),
			"source":  sysSale.Source(),
			"owner":   "*sys",
			"comment": mikhmon.SalesComment,
		})
		if err != nil {
			return "", err
		}
		t.onCleanup("sales *sys "+sysSale.Name(), func() error {
			return t.removeByID("mikhmon-sales", id)
		})
		sysRows, err := t.once("mikhmon-sales", "report",
			ros.Query{Prop: "comment", Value: mikhmon.SalesComment},
			ros.Query{Prop: "owner", Value: "*sys"})
		if err != nil {
			return "", err
		}
		var found bool
		for _, row := range sysRows {
			if row["name"] != sysSale.Name() {
				continue
			}
			found = true
			if _, err := t.action("system-script", "set", row[".id"], map[string]string{
				"owner": owner,
			}); err != nil {
				return "", err
			}
		}
		if !found {
			return "", errors.New("query comment=mikhmon&owner=*sys tidak menemukan laporan uji")
		}
		rows, err := t.client.Sales(ctx, owner)
		if err != nil {
			return "", err
		}
		if !containsName(rows, sysSale.User) {
			return "", errors.New("set owner tidak membuat laporan tampil di query owner")
		}
		return fmt.Sprintf("owner=*sys ditemukan lalu dinormalisasi ke %s", owner), nil
	}); err != nil {
		return err
	}

	if err := t.must("expire-monitor/install", func() (string, error) {
		_, err := t.client.ExpireMonitor(ctx)
		if err == nil {
			return "sudah ada; tidak diubah", nil
		}
		if !errors.Is(err, mikhmon.ErrNotFound) {
			return "", err
		}
		if err := t.client.InstallExpireMonitor(ctx); err != nil {
			return "", err
		}
		t.onCleanup("scheduler "+mikhmon.ExpireMonitorName, func() error {
			return t.removeByName("system-scheduler", "list", "name", mikhmon.ExpireMonitorName)
		})
		return "dipasang", nil
	}); err != nil {
		return err
	}

	var expMon map[string]string
	if err := t.must("expire-monitor/verifikasi", func() (string, error) {
		row, err := t.client.ExpireMonitor(ctx)
		if err != nil {
			return "", err
		}
		expMon = row
		if row["disabled"] == "true" {
			return "", errors.New("scheduler disabled")
		}
		if row["interval"] != "00:01:00" && row["interval"] != "1m" {
			return "", fmt.Errorf("interval=%q, want 00:01:00/1m", row["interval"])
		}
		if len(row["on-event"]) < 1500 {
			return "", fmt.Errorf("on-event terlalu pendek: %d", len(row["on-event"]))
		}
		return fmt.Sprintf("interval=%s on-event=%d byte", row["interval"], len(row["on-event"])), nil
	}); err != nil {
		return err
	}

	pastX := fmt.Sprintf("jan/01/%s 00:00:00 X", year)
	pastN := fmt.Sprintf("jan/01/%s 00:00:00 N", year)
	for _, u := range []struct {
		name    string
		comment string
	}{{userX, pastX}, {userN, pastN}} {
		u := u
		if err := t.must("expire-monitor/user "+u.name, func() (string, error) {
			id, err := t.client.AddUser(ctx, mikhmon.User{
				Name: u.name, Password: u.name, Profile: profileName,
				Server: "all", Comment: u.comment,
			})
			if err != nil {
				return "", err
			}
			t.onCleanup("user "+u.name, func() error {
				return t.removeByName("hotspot-user", "list", "name", u.name)
			})
			return "id=" + id + " comment=" + u.comment, nil
		}); err != nil {
			return err
		}
	}

	if err := t.must("expire-monitor/jalankan", func() (string, error) {
		script := expMon["on-event"]
		if script == "" {
			script = mikhmon.ExpireMonitorScript()
		}
		if _, err := t.action("system-script", "add", "", map[string]string{
			"name": expMonName, "source": script, "comment": tag,
		}); err != nil {
			return "", err
		}
		t.onCleanup("script "+expMonName, func() error {
			return t.removeByName("system-script", "list", "name", expMonName)
		})
		if _, err := t.action("system-script", "run", expMonName, nil); err != nil {
			return "", err
		}
		time.Sleep(2 * time.Second)
		return "script on-event dijalankan manual", nil
	}); err != nil {
		return err
	}

	if err := t.must("expire-monitor/mode X menghapus user", func() (string, error) {
		rows, err := t.client.Users(ctx, profileName)
		if err != nil {
			return "", err
		}
		if findName(rows, userX) {
			return "", fmt.Errorf("user X %q masih ada", userX)
		}
		return userX + " terhapus", nil
	}); err != nil {
		return err
	}

	if err := t.must("expire-monitor/mode N limit-uptime", func() (string, error) {
		rows, err := t.client.Users(ctx, profileName)
		if err != nil {
			return "", err
		}
		for _, row := range rows {
			if row["name"] != userN {
				continue
			}
			lim := row["limit-uptime"]
			if lim != "1s" && lim != "00:00:01" {
				return "", fmt.Errorf("limit-uptime=%q, want 1s/00:00:01", lim)
			}
			return "limit-uptime=" + lim, nil
		}
		return "", fmt.Errorf("user N %q hilang, seharusnya tetap ada", userN)
	}); err != nil {
		return err
	}

	return nil
}

func (t *tester) must(name string, fn func() (string, error)) error {
	detail, err := fn()
	t.rep.add(name, err == nil, detail, err)
	return err
}

func (t *tester) onCleanup(name string, fn func() error) {
	t.cleans = append(t.cleans, cleanup{name: name, fn: fn})
}

func (t *tester) cleanup() {
	for i := len(t.cleans) - 1; i >= 0; i-- {
		c := t.cleans[i]
		err := c.fn()
		t.rep.add("cleanup/"+c.name, err == nil, "", err)
	}
}

func (t *tester) once(id, view string, q ...ros.Query) ([]map[string]string, error) {
	r, ok := ros.Lookup(id)
	if !ok {
		return nil, fmt.Errorf("resource %s tidak ada", id)
	}
	req, err := r.Build(view, ros.AnyVersion, ros.Once(), ros.Params{Query: q})
	if err != nil {
		return nil, err
	}
	return t.run.Once(t.ctx, req)
}

func (t *tester) action(id, name, selector string, args map[string]string) (string, error) {
	r, ok := ros.Lookup(id)
	if !ok {
		return "", fmt.Errorf("resource %s tidak ada", id)
	}
	a, ok := r.Action(name)
	if !ok {
		return "", fmt.Errorf("aksi %s/%s tidak ada", id, name)
	}
	sent, err := a.Sentence(selector, args, "")
	if err != nil {
		return "", err
	}
	return t.run.Exec(t.ctx, sent)
}

func (t *tester) removeByName(id, view, prop, value string) error {
	rows, err := t.once(id, view, ros.Query{Prop: prop, Value: value})
	if err != nil {
		return err
	}
	for _, row := range rows {
		if err := t.removeByID(id, row[".id"]); err != nil {
			return err
		}
	}
	return nil
}

func (t *tester) removeByID(id, itemID string) error {
	if itemID == "" {
		return nil
	}
	_, err := t.action(id, "remove", itemID, nil)
	return err
}

func (t *tester) salesForUser(user string) ([]map[string]string, error) {
	rows, err := t.once("mikhmon-sales", "report",
		ros.Query{Prop: "comment", Value: "mikhmon"})
	if err != nil {
		return nil, err
	}
	out := make([]map[string]string, 0)
	for _, row := range rows {
		if strings.Contains(row["name"], user) {
			out = append(out, row)
		}
	}
	return out, nil
}

func containsName(rows []map[string]string, user string) bool {
	for _, row := range rows {
		if strings.Contains(row["name"], user) {
			return true
		}
	}
	return false
}

func findName(rows []map[string]string, name string) bool {
	for _, row := range rows {
		if row["name"] == name {
			return true
		}
	}
	return false
}

func legacyDate(date string) (string, error) {
	date = strings.TrimSpace(date)
	if len(date) >= 10 && date[4] == '-' && date[7] == '-' {
		t, err := time.Parse("2006-01-02", date[:10])
		if err != nil {
			return "", err
		}
		return strings.ToLower(t.Format("Jan/02/2006")), nil
	}
	parts := strings.Split(date, "/")
	if len(parts) == 3 && len(parts[2]) == 4 {
		return strings.ToLower(date), nil
	}
	return "", fmt.Errorf("format tanggal tidak dikenal")
}

func mmddyy(legacy string) string {
	parts := strings.Split(legacy, "/")
	if len(parts) != 3 {
		return time.Now().Format("01.02.06")
	}
	months := map[string]string{
		"jan": "01", "feb": "02", "mar": "03", "apr": "04",
		"may": "05", "jun": "06", "jul": "07", "aug": "08",
		"sep": "09", "oct": "10", "nov": "11", "dec": "12",
	}
	month := months[strings.ToLower(parts[0])]
	if month == "" {
		month = time.Now().Format("01")
	}
	year := parts[2]
	if len(year) == 4 {
		year = year[2:]
	}
	return fmt.Sprintf("%s.%s.%s", month, parts[1], year)
}

func yearOf(legacy string) string {
	parts := strings.Split(legacy, "/")
	if len(parts) == 3 && len(parts[2]) == 4 {
		return parts[2]
	}
	return ""
}
