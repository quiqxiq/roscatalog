package fields

// ========================================================================
// Mikhmon Sales (/system/script dengan comment=mikhmon)
// Resource: mikhmon-sales
// ========================================================================

// MikhmonSalesRow adalah kolom hasil print laporan penjualan Mikhmon (/system/script).
type MikhmonSalesRow struct {
	ID                     string `ros:".id,key"`
	Name                   string `ros:"name"`
	Owner                  string `ros:"owner"`
	Source                 string `ros:"source"`
	Comment                string `ros:"comment"`
	DontRequirePermissions string `ros:"dont-require-permissions,omitempty"`
	Invalid                string `ros:"invalid,omitempty"`
	LastStarted            string `ros:"last-started,omitempty"`
	Policy                 string `ros:"policy,omitempty"`
	RunCount               string `ros:"run-count,omitempty"`
}

// MikhmonSalesAdd adalah argumen tambah laporan penjualan Mikhmon ke /system/script.
type MikhmonSalesAdd struct {
	Name                   string `ros:"name,required,omitempty"`
	Source                 string `ros:"source,required,omitempty"`
	Owner                  string `ros:"owner,omitempty"`
	Comment                string `ros:"comment,omitempty"`
	DontRequirePermissions string `ros:"dont-require-permissions,omitempty"`
	Policy                 string `ros:"policy,omitempty"`
}

// MikhmonSalesSet adalah argumen ubah laporan penjualan Mikhmon di /system/script.
type MikhmonSalesSet struct {
	Numbers                string `ros:"numbers,omitempty"`
	Name                   string `ros:"name,omitempty"`
	Source                 string `ros:"source,omitempty"`
	Owner                  string `ros:"owner,omitempty"`
	Comment                string `ros:"comment,omitempty"`
	DontRequirePermissions string `ros:"dont-require-permissions,omitempty"`
	Policy                 string `ros:"policy,omitempty"`
}

// Alias konsisten antara bentuk tunggal dan jamak.
type MikhmonSaleRow = MikhmonSalesRow
type MikhmonSaleAdd = MikhmonSalesAdd
type MikhmonSaleSet = MikhmonSalesSet

// ========================================================================
// Mikhmon QuickPrint (/system/script dengan comment=QuickPrintMikhmon)
// ========================================================================

// MikhmonQuickPrintRow adalah kolom hasil print preset QuickPrint Mikhmon.
type MikhmonQuickPrintRow struct {
	ID                     string `ros:".id,key"`
	Name                   string `ros:"name"`
	Source                 string `ros:"source"`
	Comment                string `ros:"comment"`
	Owner                  string `ros:"owner,omitempty"`
	DontRequirePermissions string `ros:"dont-require-permissions,omitempty"`
	Policy                 string `ros:"policy,omitempty"`
	RunCount               string `ros:"run-count,omitempty"`
}

// MikhmonQuickPrintAdd adalah argumen tambah preset QuickPrint Mikhmon.
type MikhmonQuickPrintAdd struct {
	Name                   string `ros:"name,required,omitempty"`
	Source                 string `ros:"source,required,omitempty"`
	Comment                string `ros:"comment,omitempty"`
	Owner                  string `ros:"owner,omitempty"`
	DontRequirePermissions string `ros:"dont-require-permissions,omitempty"`
	Policy                 string `ros:"policy,omitempty"`
}

// MikhmonQuickPrintSet adalah argumen ubah preset QuickPrint Mikhmon.
type MikhmonQuickPrintSet struct {
	Numbers                string `ros:"numbers,omitempty"`
	Name                   string `ros:"name,omitempty"`
	Source                 string `ros:"source,omitempty"`
	Comment                string `ros:"comment,omitempty"`
	Owner                  string `ros:"owner,omitempty"`
	DontRequirePermissions string `ros:"dont-require-permissions,omitempty"`
	Policy                 string `ros:"policy,omitempty"`
}

// ========================================================================
// Mikhmon Hotspot Voucher (/ip/hotspot/user view "voucher" atau print detail)
// ========================================================================

// MikhmonVoucherRow adalah kolom hasil cetak voucher Mikhmon (termasuk password).
type MikhmonVoucherRow struct {
	ID              string `ros:".id,key"`
	Name            string `ros:"name"`
	Password        string `ros:"password"`
	Profile         string `ros:"profile"`
	LimitUptime     string `ros:"limit-uptime"`
	LimitBytesTotal string `ros:"limit-bytes-total"`
	Comment         string `ros:"comment"`
	Server          string `ros:"server,omitempty"`
	Address         string `ros:"address,omitempty"`
	MACAddress      string `ros:"mac-address,omitempty"`
	Disabled        string `ros:"disabled,omitempty"`
	Uptime          string `ros:"uptime,omitempty"`
	BytesIn         string `ros:"bytes-in,omitempty"`
	BytesOut        string `ros:"bytes-out,omitempty"`
}

// MikhmonUserAdd adalah argumen tambah hotspot user/voucher sesuai Mikhmon v4.
type MikhmonUserAdd struct {
	Name            string `ros:"name,required,omitempty"`
	Password        string `ros:"password,omitempty"`
	Profile         string `ros:"profile,omitempty"`
	Server          string `ros:"server,omitempty"`
	Comment         string `ros:"comment,omitempty"`
	Address         string `ros:"address,omitempty"`
	MACAddress      string `ros:"mac-address,omitempty"`
	LimitUptime     string `ros:"limit-uptime,omitempty"`
	LimitBytesIn    string `ros:"limit-bytes-in,omitempty"`
	LimitBytesOut   string `ros:"limit-bytes-out,omitempty"`
	LimitBytesTotal string `ros:"limit-bytes-total,omitempty"`
	Email           string `ros:"email,omitempty"`
	Disabled        string `ros:"disabled,omitempty"`
	Routes          string `ros:"routes,omitempty"`
}

// MikhmonUserSet adalah argumen ubah hotspot user/voucher sesuai Mikhmon v4.
type MikhmonUserSet struct {
	Numbers         string `ros:"numbers,omitempty"`
	Name            string `ros:"name,omitempty"`
	Password        string `ros:"password,omitempty"`
	Profile         string `ros:"profile,omitempty"`
	Server          string `ros:"server,omitempty"`
	Comment         string `ros:"comment,omitempty"`
	Address         string `ros:"address,omitempty"`
	MACAddress      string `ros:"mac-address,omitempty"`
	LimitUptime     string `ros:"limit-uptime,omitempty"`
	LimitBytesIn    string `ros:"limit-bytes-in,omitempty"`
	LimitBytesOut   string `ros:"limit-bytes-out,omitempty"`
	LimitBytesTotal string `ros:"limit-bytes-total,omitempty"`
	Email           string `ros:"email,omitempty"`
	Disabled        string `ros:"disabled,omitempty"`
	Routes          string `ros:"routes,omitempty"`
}

// ========================================================================
// Mikhmon Hotspot User Profile (/ip/hotspot/user/profile)
// ========================================================================

// MikhmonProfileRow adalah kolom hasil detail user profile pada Mikhmon v4.
type MikhmonProfileRow struct {
	ID                string `ros:".id,key"`
	Name              string `ros:"name"`
	AddressPool       string `ros:"address-pool"`
	RateLimit         string `ros:"rate-limit"`
	SharedUsers       string `ros:"shared-users"`
	StatusAutorefresh string `ros:"status-autorefresh"`
	ParentQueue       string `ros:"parent-queue"`
	OnLogin           string `ros:"on-login"`
	OnLogout          string `ros:"on-logout,omitempty"`
	Comment           string `ros:"comment,omitempty"`
	SessionTimeout    string `ros:"session-timeout,omitempty"`
	IdleTimeout       string `ros:"idle-timeout,omitempty"`
	KeepaliveTimeout  string `ros:"keepalive-timeout,omitempty"`
	TransparentProxy  string `ros:"transparent-proxy,omitempty"`
	InsertQueueBefore string `ros:"insert-queue-before,omitempty"`
}

// MikhmonProfileAdd adalah argumen tambah profile hotspot pada Mikhmon v4.
type MikhmonProfileAdd struct {
	Name              string `ros:"name,required,omitempty"`
	AddressPool       string `ros:"address-pool,omitempty"`
	RateLimit         string `ros:"rate-limit,omitempty"`
	SharedUsers       string `ros:"shared-users,omitempty"`
	StatusAutorefresh string `ros:"status-autorefresh,omitempty"`
	ParentQueue       string `ros:"parent-queue,omitempty"`
	OnLogin           string `ros:"on-login,omitempty"`
	OnLogout          string `ros:"on-logout,omitempty"`
	Comment           string `ros:"comment,omitempty"`
	TransparentProxy  string `ros:"transparent-proxy,omitempty"`
}

// MikhmonProfileSet adalah argumen ubah profile hotspot pada Mikhmon v4.
type MikhmonProfileSet struct {
	Numbers           string `ros:"numbers,omitempty"`
	Name              string `ros:"name,omitempty"`
	AddressPool       string `ros:"address-pool,omitempty"`
	RateLimit         string `ros:"rate-limit,omitempty"`
	SharedUsers       string `ros:"shared-users,omitempty"`
	StatusAutorefresh string `ros:"status-autorefresh,omitempty"`
	ParentQueue       string `ros:"parent-queue,omitempty"`
	OnLogin           string `ros:"on-login,omitempty"`
	OnLogout          string `ros:"on-logout,omitempty"`
	Comment           string `ros:"comment,omitempty"`
	TransparentProxy  string `ros:"transparent-proxy,omitempty"`
}

// ========================================================================
// Mikhmon Scheduler (/system/scheduler: Expire Monitor & Profile Monitor)
// ========================================================================

// MikhmonSchedulerRow adalah baris scheduler yang dipantau atau dibuat oleh Mikhmon v4.
type MikhmonSchedulerRow struct {
	ID        string `ros:".id,key"`
	Name      string `ros:"name"`
	StartTime string `ros:"start-time"`
	Interval  string `ros:"interval"`
	OnEvent   string `ros:"on-event"`
	Comment   string `ros:"comment"`
	Disabled  string `ros:"disabled"`
	NextRun   string `ros:"next-run,omitempty"`
	RunCount  string `ros:"run-count,omitempty"`
	StartDate string `ros:"start-date,omitempty"`
}

// MikhmonSchedulerAdd adalah argumen tambah scheduler Mikhmon.
type MikhmonSchedulerAdd struct {
	Name      string `ros:"name,required,omitempty"`
	StartTime string `ros:"start-time,omitempty"`
	Interval  string `ros:"interval,omitempty"`
	OnEvent   string `ros:"on-event,required,omitempty"`
	Comment   string `ros:"comment,omitempty"`
	Disabled  string `ros:"disabled,omitempty"`
	StartDate string `ros:"start-date,omitempty"`
}

// MikhmonSchedulerSet adalah argumen ubah scheduler Mikhmon.
type MikhmonSchedulerSet struct {
	Numbers   string `ros:"numbers,omitempty"`
	Name      string `ros:"name,omitempty"`
	StartTime string `ros:"start-time,omitempty"`
	Interval  string `ros:"interval,omitempty"`
	OnEvent   string `ros:"on-event,omitempty"`
	Comment   string `ros:"comment,omitempty"`
	Disabled  string `ros:"disabled,omitempty"`
	StartDate string `ros:"start-date,omitempty"`
}

// ========================================================================
// Mikhmon Hotspot Active Session (/ip/hotspot/active)
// ========================================================================

// MikhmonActiveRow adalah kolom hasil sesi hotspot aktif pada dashboard Mikhmon.
type MikhmonActiveRow struct {
	ID               string `ros:".id,key"`
	Server           string `ros:"server"`
	User             string `ros:"user"`
	Address          string `ros:"address"`
	MACAddress       string `ros:"mac-address"`
	Uptime           string `ros:"uptime"`
	BytesIn          string `ros:"bytes-in"`
	BytesOut         string `ros:"bytes-out"`
	SessionTimeLeft  string `ros:"session-time-left,omitempty"`
	KeepaliveTimeout string `ros:"keepalive-timeout,omitempty"`
	LoginBy          string `ros:"login-by,omitempty"`
	Comment          string `ros:"comment,omitempty"`
}
