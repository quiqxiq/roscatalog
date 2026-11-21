package mikhmon

import (
	"context"
	"errors"
	"fmt"
	"strings"

	ros "example.com/roscatalog"
)

// ErrNotFound dipakai bila baris yang dicari tidak ada di router.
var ErrNotFound = errors.New("mikhmon: data tidak ditemukan")

// Executor adalah bagian rosrun.Runner yang dipakai Client. Dengan interface
// ini paket mikhmon tetap tidak bergantung pada goros dan mudah diuji.
type Executor interface {
	Once(ctx context.Context, req ros.Request) ([]map[string]string, error)
	Count(ctx context.Context, req ros.Request) (int, error)
	Exec(ctx context.Context, sentence []string) (string, error)
}

// Client membungkus workflow Mikhmon di atas Executor.
type Client struct {
	run Executor
}

// New membuat Client. Biasanya run adalah *rosrun.Runner.
func New(run Executor) *Client {
	return &Client{run: run}
}

// User adalah field hotspot user yang dipakai Mikhmon.
type User struct {
	Name            string
	Password        string
	Profile         string
	Server          string
	Comment         string
	Address         string
	MACAddress      string
	LimitUptime     string
	LimitBytesTotal string
	Email           string
}

func (u User) args() map[string]string {
	m := map[string]string{"name": u.Name}
	if u.Password != "" {
		m["password"] = u.Password
	}
	if u.Profile != "" {
		m["profile"] = u.Profile
	}
	if u.Server != "" {
		m["server"] = u.Server
	}
	if u.Comment != "" {
		m["comment"] = u.Comment
	}
	if u.Address != "" {
		m["address"] = u.Address
	}
	if u.MACAddress != "" {
		m["mac-address"] = u.MACAddress
	}
	if u.LimitUptime != "" {
		m["limit-uptime"] = u.LimitUptime
	}
	if u.LimitBytesTotal != "" {
		m["limit-bytes-total"] = u.LimitBytesTotal
	}
	if u.Email != "" {
		m["email"] = u.Email
	}
	return m
}

// MikhmonComment menambahkan prefix vc-/up- seperti post_add_user.php.
func MikhmonComment(name, password, comment string) string {
	if name == password {
		return "vc-" + comment
	}
	return "up-" + comment
}

// AddUser menambah satu hotspot user/voucher.
func (c *Client) AddUser(ctx context.Context, u User) (string, error) {
	if strings.TrimSpace(u.Name) == "" {
		return "", errors.New("mikhmon: user.Name wajib diisi")
	}
	return c.action(ctx, "hotspot-user", "add", "", u.args())
}

// AddUsers menambah beberapa user berurutan. Bila gagal di tengah, ID yang
// sudah dibuat tetap dikembalikan bersama error.
func (c *Client) AddUsers(ctx context.Context, users []User) ([]string, error) {
	out := make([]string, 0, len(users))
	for _, u := range users {
		id, err := c.AddUser(ctx, u)
		if err != nil {
			return out, fmt.Errorf("mikhmon: user %q: %w", u.Name, err)
		}
		out = append(out, id)
	}
	return out, nil
}

// Profile adalah field hotspot user profile yang dipakai Mikhmon.
type Profile struct {
	Name              string
	AddressPool       string
	RateLimit         string
	SharedUsers       string
	StatusAutorefresh string
	ParentQueue       string
	OnLogin           string
}

// AddProfile menambah profile. StatusAutorefresh default "1m" seperti Mikhmon.
func (c *Client) AddProfile(ctx context.Context, p Profile) (string, error) {
	if strings.TrimSpace(p.Name) == "" {
		return "", errors.New("mikhmon: profile.Name wajib diisi")
	}
	if err := checkOnLoginScript(p.OnLogin); err != nil {
		return "", err
	}
	args := map[string]string{
		"name":               p.Name,
		"status-autorefresh": p.StatusAutorefresh,
	}
	if args["status-autorefresh"] == "" {
		args["status-autorefresh"] = "1m"
	}
	if p.AddressPool != "" {
		args["address-pool"] = p.AddressPool
	}
	if p.RateLimit != "" {
		args["rate-limit"] = p.RateLimit
	}
	if p.SharedUsers != "" {
		args["shared-users"] = p.SharedUsers
	}
	if p.ParentQueue != "" {
		args["parent-queue"] = p.ParentQueue
	}
	if p.OnLogin != "" {
		args["on-login"] = p.OnLogin
	}
	return c.action(ctx, "hotspot-user-profile", "add", "", args)
}

// SetOnLogin mengganti properti on-login sebuah profile. Script harus satu
// baris; gunakan OnLoginScript agar aman.
func (c *Client) SetOnLogin(ctx context.Context, profile, script string) error {
	if err := checkOnLoginScript(script); err != nil {
		return err
	}
	rows, err := c.once(ctx, "hotspot-user-profile", "detail",
		[]ros.Query{{Prop: "name", Value: profile}})
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return fmt.Errorf("%w: profile %q", ErrNotFound, profile)
	}
	id := rows[0][".id"]
	if id == "" {
		return fmt.Errorf("mikhmon: profile %q tidak mengembalikan .id", profile)
	}
	_, err = c.action(ctx, "hotspot-user-profile", "set", id,
		map[string]string{"on-login": script})
	return err
}

// InstallExpireMonitor membuat atau memperbarui scheduler
// Mikhmon-Expire-Monitor lalu memastikan scheduler aktif.
func (c *Client) InstallExpireMonitor(ctx context.Context) error {
	script := ExpireMonitorScript()
	rows, err := c.once(ctx, "system-scheduler", "list",
		[]ros.Query{{Prop: "name", Value: ExpireMonitorName}})
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		_, err = c.action(ctx, "system-scheduler", "add", "", map[string]string{
			"name":       ExpireMonitorName,
			"start-time": "00:00:00",
			"interval":   "00:01:00",
			"on-event":   script,
			"comment":    ExpireMonitorComment,
		})
		return err
	}
	id := rows[0][".id"]
	if id == "" {
		return errors.New("mikhmon: scheduler Expire Monitor tidak mengembalikan .id")
	}
	if _, err := c.action(ctx, "system-scheduler", "set", id, map[string]string{
		"interval": "00:01:00",
		"on-event": script,
		"comment":  ExpireMonitorComment,
	}); err != nil {
		return err
	}
	_, err = c.action(ctx, "system-scheduler", "enable", id, nil)
	return err
}

// InstallProfileMonitor membuat atau memperbarui scheduler Monitor Profile
// untuk satu user profile di /system/scheduler persis seperti adduserprofile.php di Mikhmon v4.
// Bila mode adalah ExpiryNone, scheduler monitor untuk profile ini dihapus.
func (c *Client) InstallProfileMonitor(ctx context.Context, profile string, mode ExpiryMode) error {
	if mode == ExpiryNone {
		return c.RemoveProfileMonitor(ctx, profile)
	}
	script, err := ProfileMonitorScript(profile, mode)
	if err != nil {
		return err
	}
	comment := ProfileMonitorCommentPrefix + profile
	rows, err := c.once(ctx, "system-scheduler", "list",
		[]ros.Query{{Prop: "name", Value: profile}})
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		_, err = c.action(ctx, "system-scheduler", "add", "", map[string]string{
			"name":       profile,
			"start-time": "01:00:00",
			"interval":   "00:02:30",
			"on-event":   script,
			"comment":    comment,
			"disabled":   "no",
		})
		return err
	}
	id := rows[0][".id"]
	if id == "" {
		return fmt.Errorf("mikhmon: scheduler Monitor Profile %q tidak mengembalikan .id", profile)
	}
	if _, err := c.action(ctx, "system-scheduler", "set", id, map[string]string{
		"interval": "00:02:30",
		"on-event": script,
		"comment":  comment,
	}); err != nil {
		return err
	}
	_, err = c.action(ctx, "system-scheduler", "enable", id, nil)
	return err
}

// RemoveProfileMonitor menghapus scheduler Monitor Profile jika ada.
func (c *Client) RemoveProfileMonitor(ctx context.Context, profile string) error {
	rows, err := c.once(ctx, "system-scheduler", "list",
		[]ros.Query{{Prop: "name", Value: profile}})
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}
	id := rows[0][".id"]
	if id == "" {
		return nil
	}
	_, err = c.action(ctx, "system-scheduler", "remove", id, nil)
	return err
}

// ProfileMonitor membaca scheduler Monitor Profile untuk profile tertentu.
func (c *Client) ProfileMonitor(ctx context.Context, profile string) (map[string]string, error) {
	rows, err := c.once(ctx, "system-scheduler", "list",
		[]ros.Query{{Prop: "name", Value: profile}})
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("%w: scheduler monitor profile %q", ErrNotFound, profile)
	}
	return rows[0], nil
}

// RemoveProfile menghapus user profile sekaligus menghapus scheduler Monitor Profile-nya,
// sama seperti removeuserprofile.php di Mikhmon v4.
func (c *Client) RemoveProfile(ctx context.Context, profile string) error {
	rows, err := c.once(ctx, "hotspot-user-profile", "detail",
		[]ros.Query{{Prop: "name", Value: profile}})
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return fmt.Errorf("%w: profile %q", ErrNotFound, profile)
	}
	id := rows[0][".id"]
	if id == "" {
		return fmt.Errorf("mikhmon: profile %q tidak mengembalikan .id", profile)
	}
	if _, err := c.action(ctx, "hotspot-user-profile", "remove", id, nil); err != nil {
		return err
	}
	return c.RemoveProfileMonitor(ctx, profile)
}

// ExpireMonitor membaca scheduler Expire Monitor.
func (c *Client) ExpireMonitor(ctx context.Context) (map[string]string, error) {
	rows, err := c.once(ctx, "system-scheduler", "list",
		[]ros.Query{{Prop: "name", Value: ExpireMonitorName}})
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrNotFound
	}
	return rows[0], nil
}

// AddSale menulis satu laporan penjualan ke /system/script.
func (c *Client) AddSale(ctx context.Context, s Sale) (string, error) {
	if err := s.Validate(); err != nil {
		return "", err
	}
	return c.action(ctx, "mikhmon-sales", "add", "", map[string]string{
		"name":    s.Name(),
		"source":  s.Source(),
		"owner":   s.Owner(),
		"comment": SalesComment,
	})
}

// Sales membaca laporan berdasarkan owner (mis. nov2026).
func (c *Client) Sales(ctx context.Context, owner string) ([]map[string]string, error) {
	q := []ros.Query(nil)
	if owner != "" {
		q = append(q, ros.Query{Prop: "owner", Value: owner})
	}
	return c.once(ctx, "mikhmon-sales", "report", q)
}

// NormalizeSalesOwner memperbaiki laporan Mikhmon yang ditulis dengan
// owner=*sys (perilaku RouterOS 7.24+) menjadi bulan+tahun. Bila owner kosong,
// owner diturunkan dari source setiap laporan. Ini jalur yang dipakai
// get_dashboard.php sebelum menampilkan laporan bulanan.
func (c *Client) NormalizeSalesOwner(ctx context.Context, owner string) (int, error) {
	rows, err := c.once(ctx, "mikhmon-sales", "report",
		[]ros.Query{
			{Prop: "comment", Value: SalesComment},
			{Prop: "owner", Value: "*sys"},
		})
	if err != nil {
		return 0, err
	}
	n := 0
	for _, row := range rows {
		id := row[".id"]
		if id == "" {
			continue
		}
		target := owner
		if target == "" {
			target = (Sale{Date: row["source"]}).Owner()
		}
		if target == "" {
			continue
		}
		if _, err := c.action(ctx, "system-script", "set", id,
			map[string]string{"owner": target}); err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

// SaveQuickPrint membuat atau memperbarui preset QuickPrint berdasarkan nama.
func (c *Client) SaveQuickPrint(ctx context.Context, p QuickPrintPreset) (string, error) {
	if strings.TrimSpace(p.Package) == "" {
		return "", errors.New("mikhmon: QuickPrint.Package wajib diisi")
	}
	name := p.Name()
	args := map[string]string{
		"name":    name,
		"source":  p.Source(),
		"comment": QuickPrintComment,
	}
	rows, err := c.once(ctx, "system-script", "list",
		[]ros.Query{{Prop: "name", Value: name}})
	if err != nil {
		return "", err
	}
	if len(rows) > 0 {
		id := rows[0][".id"]
		if id == "" {
			return "", fmt.Errorf("mikhmon: script %q tidak mengembalikan .id", name)
		}
		return c.action(ctx, "system-script", "set", id, args)
	}
	return c.action(ctx, "system-script", "add", "", args)
}

// QuickPrints membaca semua preset QuickPrint beserta source-nya.
func (c *Client) QuickPrints(ctx context.Context) ([]map[string]string, error) {
	return c.once(ctx, "system-script", "detail",
		[]ros.Query{{Prop: "comment", Value: QuickPrintComment}})
}

// RemoveExpiredUsers menghapus user dengan limit-uptime=1s, penanda yang
// dipakai Expire Monitor Mikhmon.
func (c *Client) RemoveExpiredUsers(ctx context.Context) (int, error) {
	rows, err := c.once(ctx, "hotspot-user", "list",
		[]ros.Query{{Prop: "limit-uptime", Value: ExpiredLimitUptime}})
	if err != nil {
		return 0, err
	}
	n := 0
	for _, row := range rows {
		id := row[".id"]
		if id == "" {
			continue
		}
		if _, err := c.action(ctx, "hotspot-user", "remove", id, nil); err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

// KickUser memutus semua sesi aktif milik user.
func (c *Client) KickUser(ctx context.Context, user string) (int, error) {
	rows, err := c.once(ctx, "hotspot-active", "list",
		[]ros.Query{{Prop: "user", Value: user}})
	if err != nil {
		return 0, err
	}
	n := 0
	for _, row := range rows {
		id := row[".id"]
		if id == "" {
			continue
		}
		if _, err := c.action(ctx, "hotspot-active", "remove", id, nil); err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

// Vouchers membaca user berdasarkan comment batch. Comment harus cocok persis.
func (c *Client) Vouchers(ctx context.Context, comment string) ([]map[string]string, error) {
	return c.once(ctx, "hotspot-user", "voucher",
		[]ros.Query{{Prop: "comment", Value: comment}})
}

// Users membaca user, opsional difilter profile.
func (c *Client) Users(ctx context.Context, profile string) ([]map[string]string, error) {
	q := []ros.Query(nil)
	if profile != "" {
		q = append(q, ros.Query{Prop: "profile", Value: profile})
	}
	return c.once(ctx, "hotspot-user", "list", q)
}

// CountUsers menghitung user, opsional difilter profile.
func (c *Client) CountUsers(ctx context.Context, profile string) (int, error) {
	r, err := resource("hotspot-user")
	if err != nil {
		return 0, err
	}
	q := []ros.Query(nil)
	if profile != "" {
		q = append(q, ros.Query{Prop: "profile", Value: profile})
	}
	req, err := r.Build("count", ros.Once(), ros.Params{Query: q})
	if err != nil {
		return 0, err
	}
	return c.run.Count(ctx, req)
}

// CountActive menghitung sesi hotspot aktif.
func (c *Client) CountActive(ctx context.Context) (int, error) {
	r, err := resource("hotspot-active")
	if err != nil {
		return 0, err
	}
	req, err := r.Build("count", ros.Once(), ros.Params{})
	if err != nil {
		return 0, err
	}
	return c.run.Count(ctx, req)
}

func (c *Client) once(ctx context.Context, id, view string, q []ros.Query) ([]map[string]string, error) {
	r, err := resource(id)
	if err != nil {
		return nil, err
	}
	req, err := r.Build(view, ros.Once(), ros.Params{Query: q})
	if err != nil {
		return nil, err
	}
	return c.run.Once(ctx, req)
}

func (c *Client) action(ctx context.Context, id, name, selector string, args map[string]string) (string, error) {
	r, err := resource(id)
	if err != nil {
		return "", err
	}
	a, ok := r.Action(name)
	if !ok {
		return "", fmt.Errorf("mikhmon: aksi %s/%s tidak ada", id, name)
	}
	sent, err := a.Sentence(selector, args, "")
	if err != nil {
		return "", err
	}
	return c.run.Exec(ctx, sent)
}

func checkOnLoginScript(script string) error {
	if strings.ContainsAny(script, "\r\n") {
		return errors.New("mikhmon: on-login harus satu baris; newline akan memecah parameter API")
	}
	return nil
}
