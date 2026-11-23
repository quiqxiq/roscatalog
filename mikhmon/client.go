package mikhmon

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	ros "example.com/roscatalog"
	"example.com/roscatalog/fields"
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

// ToFields membungkus User ke typed struct fields.MikhmonUserAdd.
func (u User) ToFields() fields.MikhmonUserAdd {
	return fields.MikhmonUserAdd{
		Name:            u.Name,
		Password:        u.Password,
		Profile:         u.Profile,
		Server:          u.Server,
		Comment:         u.Comment,
		Address:         u.Address,
		MACAddress:      u.MACAddress,
		LimitUptime:     u.LimitUptime,
		LimitBytesTotal: u.LimitBytesTotal,
		Email:           u.Email,
	}
}

// UserFromVoucher mengonversi fields.MikhmonVoucherRow menjadi domain User.
func UserFromVoucher(v fields.MikhmonVoucherRow) User {
	return User{
		Name:            v.Name,
		Password:        v.Password,
		Profile:         v.Profile,
		Server:          v.Server,
		Comment:         v.Comment,
		Address:         v.Address,
		MACAddress:      v.MACAddress,
		LimitUptime:     v.LimitUptime,
		LimitBytesTotal: v.LimitBytesTotal,
	}
}

func (u User) args() map[string]string {
	m, err := ros.ArgsOf(u.ToFields())
	if err != nil || m == nil {
		m = map[string]string{"name": u.Name}
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

// ToFields membungkus Profile ke typed struct fields.MikhmonProfileAdd.
func (p Profile) ToFields() fields.MikhmonProfileAdd {
	refresh := p.StatusAutorefresh
	if refresh == "" {
		refresh = "1m"
	}
	return fields.MikhmonProfileAdd{
		Name:              p.Name,
		AddressPool:       p.AddressPool,
		RateLimit:         p.RateLimit,
		SharedUsers:       p.SharedUsers,
		StatusAutorefresh: refresh,
		ParentQueue:       p.ParentQueue,
		OnLogin:           p.OnLogin,
	}
}

// AddProfile menambah profile. StatusAutorefresh default "1m" seperti Mikhmon.
func (c *Client) AddProfile(ctx context.Context, p Profile) (string, error) {
	if strings.TrimSpace(p.Name) == "" {
		return "", errors.New("mikhmon: profile.Name wajib diisi")
	}
	if err := checkOnLoginScript(p.OnLogin); err != nil {
		return "", err
	}
	args, err := ros.ArgsOf(p.ToFields())
	if err != nil {
		return "", err
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

// ProfileMonitorRow membaca scheduler Monitor Profile dalam bentuk typed fields.MikhmonSchedulerRow.
func (c *Client) ProfileMonitorRow(ctx context.Context, profile string) (fields.MikhmonSchedulerRow, error) {
	row, err := c.ProfileMonitor(ctx, profile)
	if err != nil {
		return fields.MikhmonSchedulerRow{}, err
	}
	return ros.Decode[fields.MikhmonSchedulerRow](row)
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

// ExpireMonitorRow membaca scheduler Expire Monitor dalam bentuk typed fields.MikhmonSchedulerRow.
func (c *Client) ExpireMonitorRow(ctx context.Context) (fields.MikhmonSchedulerRow, error) {
	row, err := c.ExpireMonitor(ctx)
	if err != nil {
		return fields.MikhmonSchedulerRow{}, err
	}
	return ros.Decode[fields.MikhmonSchedulerRow](row)
}

// AddSale menulis satu laporan penjualan ke /system/script.
func (c *Client) AddSale(ctx context.Context, s Sale) (string, error) {
	if err := s.Validate(); err != nil {
		return "", err
	}
	args, err := ros.ArgsOf(s.ToFields())
	if err != nil {
		return "", err
	}
	return c.action(ctx, "mikhmon-sales", "add", "", args)
}

// Sales membaca laporan penjualan Mikhmon (/system/script dengan comment=mikhmon).
// Parameter monthOrOwner opsional:
// - Bila kosong, seluruh laporan penjualan Mikhmon dikembalikan.
// - Bila diisi (mis. "nov2026", "032024", atau "2024-03"), fungsi secara cerdas
//   mencoba query server-side terlebih dahulu (format teks & numerik), lalu fallback
//   ke penyaringan sisi klien agar kompatibel dengan seluruh versi RouterOS
//   (RouterOS 6, RouterOS 7 < 7.24, dan RouterOS 7.24+ owner=*sys).
func (c *Client) Sales(ctx context.Context, monthOrOwner string) ([]map[string]string, error) {
	monthOrOwner = strings.TrimSpace(monthOrOwner)
	if monthOrOwner == "" {
		return c.once(ctx, "mikhmon-sales", "report",
			[]ros.Query{{Prop: "comment", Value: SalesComment}})
	}

	textBucket, numBucket, isBucket := MonthBuckets(monthOrOwner)
	if isBucket {
		primary := textBucket
		alt := numBucket
		if monthOrOwner == numBucket {
			primary = numBucket
			alt = textBucket
		}

		// Percobaan 1: Server query dengan format utama
		rows, err := c.once(ctx, "mikhmon-sales", "report", []ros.Query{
			{Prop: "comment", Value: SalesComment},
			{Prop: "owner", Value: primary},
		})
		if err != nil {
			return nil, err
		}
		if len(rows) > 0 {
			filtered := filterSalesByMonth(rows, monthOrOwner)
			if len(filtered) > 0 {
				return filtered, nil
			}
		}

		// Percobaan 2: Server query dengan format alternatif (jika beda)
		if alt != "" && alt != primary {
			altRows, err := c.once(ctx, "mikhmon-sales", "report", []ros.Query{
				{Prop: "comment", Value: SalesComment},
				{Prop: "owner", Value: alt},
			})
			if err != nil {
				return nil, err
			}
			if len(altRows) > 0 {
				filtered := filterSalesByMonth(altRows, monthOrOwner)
				if len(filtered) > 0 {
					return filtered, nil
				}
			}
		}

		// Fallback: router mungkin menggunakan owner=*sys (v7.24+) atau custom.
		// Ambil seluruh comment=mikhmon lalu saring berdasarkan Sale.MatchesMonth.
		allRows, err := c.once(ctx, "mikhmon-sales", "report", []ros.Query{
			{Prop: "comment", Value: SalesComment},
		})
		if err != nil {
			return nil, err
		}
		return filterSalesByMonth(allRows, monthOrOwner), nil
	}

	// Bukan format bulan standar (mis. user query owner custom)
	return c.once(ctx, "mikhmon-sales", "report", []ros.Query{
		{Prop: "comment", Value: SalesComment},
		{Prop: "owner", Value: monthOrOwner},
	})
}

func filterSalesByMonth(rows []map[string]string, month string) []map[string]string {
	var out []map[string]string
	for _, r := range rows {
		name := r["name"]
		if s, err := ParseSaleName(name); err == nil && s.MatchesMonth(month) {
			out = append(out, r)
		}
	}
	return out
}

// SalesRows membaca laporan berdasarkan owner dalam bentuk typed fields.MikhmonSalesRow.
func (c *Client) SalesRows(ctx context.Context, monthOrOwner string) ([]fields.MikhmonSalesRow, error) {
	rows, err := c.Sales(ctx, monthOrOwner)
	if err != nil {
		return nil, err
	}
	return ros.DecodeAll[fields.MikhmonSalesRow](rows)
}

// AllSales membaca semua laporan penjualan Mikhmon dan mem-parse-nya ke []Sale.
// Baris yang bukan format laporan Mikhmon diabaikan secara aman.
func (c *Client) AllSales(ctx context.Context) ([]Sale, error) {
	rows, err := c.Sales(ctx, "")
	if err != nil {
		return nil, err
	}
	out := make([]Sale, 0, len(rows))
	for _, r := range rows {
		name := r["name"]
		if s, err := ParseSaleName(name); err == nil {
			out = append(out, s)
		}
	}
	return out, nil
}

// SalesByMonth membaca laporan penjualan untuk satu bulan tertentu (mis. "nov2026", "032024", atau "2024-03").
func (c *Client) SalesByMonth(ctx context.Context, month string) ([]Sale, error) {
	rows, err := c.Sales(ctx, month)
	if err != nil {
		return nil, err
	}
	out := make([]Sale, 0, len(rows))
	for _, r := range rows {
		name := r["name"]
		if s, err := ParseSaleName(name); err == nil {
			out = append(out, s)
		}
	}
	return out, nil
}

// SalesByDate membaca laporan penjualan untuk satu tanggal tertentu (mis. "2024-03-18" atau "mar/18/2024").
func (c *Client) SalesByDate(ctx context.Context, date string) ([]Sale, error) {
	var rows []map[string]string
	var err error
	if _, _, ok := MonthBuckets(date); ok {
		rows, err = c.Sales(ctx, date)
	} else {
		rows, err = c.Sales(ctx, "")
	}
	if err != nil {
		return nil, err
	}
	out := make([]Sale, 0)
	for _, r := range rows {
		name := r["name"]
		if s, err := ParseSaleName(name); err == nil && s.MatchesDate(date) {
			out = append(out, s)
		}
	}
	return out, nil
}

// SalesStats adalah indikator kinerja penjualan (KPI) Mikhmon.
type SalesStats struct {
	TotalCount  int     `json:"total_count"`
	TotalSum    float64 `json:"total_sum"`
	MonthCount  int     `json:"month_count"`
	MonthSum    float64 `json:"month_sum"`
	PrevSum     float64 `json:"prev_sum"`
	TodayCount  int     `json:"today_count"`
	TodaySum    float64 `json:"today_sum"`
	MonthBucket string  `json:"month_bucket"`
}

// SalesStats menghitung statistik KPI penjualan langsung dari router,
// identik dengan fungsi mikhmon_scripts_stats di Mikhmon v4.
func (c *Client) SalesStats(ctx context.Context) (SalesStats, error) {
	all, err := c.AllSales(ctx)
	if err != nil {
		return SalesStats{}, err
	}

	now := time.Now()
	var clockDate string
	clockRows, err := c.once(ctx, "system-clock", "list", nil)
	if err == nil && len(clockRows) > 0 {
		clockDate = clockRows[0]["date"]
	}

	var thisBucket, prevBucket, todayISO string
	if clockDate != "" {
		if t, _, ok := MonthBuckets(clockDate); ok {
			thisBucket = t
		}
		todayISO = NormalizeDateISO(clockDate)
	}
	if thisBucket == "" {
		thisBucket = strings.ToLower(now.Format("Jan2006"))
	}
	if todayISO == "" {
		todayISO = now.Format("2006-01-02")
	}
	prevTime := now.AddDate(0, -1, 0)
	prevBucket = strings.ToLower(prevTime.Format("Jan2006"))

	stats := SalesStats{
		MonthBucket: thisBucket,
	}
	for _, s := range all {
		price, _ := strconv.ParseFloat(strings.TrimSpace(s.Price), 64)
		stats.TotalCount++
		stats.TotalSum += price

		if s.MatchesMonth(thisBucket) {
			stats.MonthCount++
			stats.MonthSum += price
		}
		if s.MatchesMonth(prevBucket) {
			stats.PrevSum += price
		}
		if s.MatchesDate(todayISO) {
			stats.TodayCount++
			stats.TodaySum += price
		}
	}
	return stats, nil
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
	args, err := ros.ArgsOf(p.ToFields())
	if err != nil {
		return "", err
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

// QuickPrintRows membaca semua preset QuickPrint dalam bentuk typed fields.MikhmonQuickPrintRow.
func (c *Client) QuickPrintRows(ctx context.Context) ([]fields.MikhmonQuickPrintRow, error) {
	rows, err := c.QuickPrints(ctx)
	if err != nil {
		return nil, err
	}
	return ros.DecodeAll[fields.MikhmonQuickPrintRow](rows)
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

// VouchersTyped membaca voucher dalam bentuk typed fields.MikhmonVoucherRow.
func (c *Client) VouchersTyped(ctx context.Context, comment string) ([]fields.MikhmonVoucherRow, error) {
	rows, err := c.Vouchers(ctx, comment)
	if err != nil {
		return nil, err
	}
	return ros.DecodeAll[fields.MikhmonVoucherRow](rows)
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
