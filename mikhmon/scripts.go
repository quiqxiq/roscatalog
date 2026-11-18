package mikhmon

import (
	"errors"
	"fmt"
	"strings"

	_ "embed"
)

// ExpiryMode adalah nilai expmode Mikhmon yang disimpan di properti on-login.
type ExpiryMode string

const (
	// ExpiryNone hanya menulis metadata validity, tanpa logika expiry.
	ExpiryNone ExpiryMode = "0"
	// ExpiryRemove menghapus user saat kedaluwarsa.
	ExpiryRemove ExpiryMode = "rem"
	// ExpiryNotice menandai limit-uptime=1s dan memutus sesi saat kedaluwarsa.
	ExpiryNotice ExpiryMode = "ntf"
	// ExpiryRemoveRecord sama seperti ExpiryRemove.
	ExpiryRemoveRecord ExpiryMode = "remc"
	// ExpiryNoticeRecord sama seperti ExpiryNotice.
	ExpiryNoticeRecord ExpiryMode = "ntfc"
)

// OnLoginOptions adalah input untuk membangun properti on-login profile.
type OnLoginOptions struct {
	Mode        ExpiryMode
	Validity    string // mis. "1d"; wajib untuk mode expiry
	LockUser    bool   // kunci MAC user saat login pertama
	LockServer  bool   // kunci server hotspot saat login pertama
	Profile     string // nama profile, dipakai untuk lookup harga di backend
	RecordSales bool   // coba catat penjualan saat login pertama
	PublicURL   string // base URL Mikmon, mis. https://mikmon.example
	HookToken   string // token per-router
	VerifyTLS   bool   // check-certificate=yes
}

const lockUserScript = `; :local mac $"mac-address"; /ip hotspot user set mac-address=$mac [find where name=$user]`
const lockServerScript = `; :local mac $"mac-address"; :local srv [/ip hotspot host get [find where mac-address="$mac"] server]; /ip hotspot user set server=$srv [find where name=$user]`

const onLoginHeaderTemplate = `
:local date [ /system clock get date ];
:local time [ /system clock get time ];
:local mnames "janfebmaraprmayjunjulaugsepoctnovdec";
:local year;
:local monthname;
:local day;
:local legacy;
:if ([:pick $date 4 5] = "-") do={
  :set year [:pick $date 0 4];
  :set day [:pick $date 8 10];
  :local mnum1 [:tonum [:pick $date 5 7]];
  :set monthname [:pick $mnames (($mnum1 - 1) * 3) (($mnum1 - 1) * 3 + 3)];
} else={
  :set year [:pick $date 7 11];
  :set day [:pick $date 4 6];
  :set monthname [:pick $date 0 3];
};
:set legacy "$monthname/$day/$year";
:local saleComment $comment;
:local ucode [:pic $comment 0 2];
`

const expiryTemplate = `
/sys sch add name="$user" disable=no start-date=$date interval="{{validity}}";
:delay 2s;
:local exp [ /sys sch get [ /sys sch find where name="$user" ] next-run ];
:local expLegacy;
:local getxp [len $exp];
:if ($getxp = 8) do={
  :set expLegacy "$legacy $exp";
} else={
  :if ($getxp = 15) do={
    :local d [:pic $exp 0 6];
    :local t [:pic $exp 7 16];
    :set expLegacy "$d/$year $t";
  } else={
    :if ([:pick $exp 4 5] = "-") do={
      :local ey [:pic $exp 0 4];
      :local ed [:pic $exp 8 10];
      :local mnum2 [:tonum [:pic $exp 5 7]];
      :local emn [:pick $mnames (($mnum2 - 1) * 3) (($mnum2 - 1) * 3 + 3)];
      :local et [:pic $exp 11 [:len $exp]];
      :set expLegacy "$emn/$ed/$ey $et";
    } else={
      :set expLegacy $exp;
    }
  }
};
/ip hotspot user set comment="$expLegacy $mode" [find where name="$user"];
/sys sch remove [find where name="$user"];
`

//go:embed expire_monitor.ros
var expireMonitorScript string

// ExpireMonitorScript mengembalikan on-event scheduler Mikhmon-Expire-Monitor.
func ExpireMonitorScript() string {
	return strings.TrimSpace(expireMonitorScript)
}

// OneLine merapatkan semua whitespace menjadi satu spasi.
func OneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// OnLoginScript membangun properti on-login satu baris. Harga dan selling
// price tidak ditulis ke script; backend membacanya dari hotspot_profiles.
func OnLoginScript(o OnLoginOptions) (string, error) {
	switch o.Mode {
	case ExpiryNone, ExpiryRemove, ExpiryNotice, ExpiryRemoveRecord, ExpiryNoticeRecord:
	default:
		return "", fmt.Errorf("mikhmon: expmode %q tidak dikenal", o.Mode)
	}
	if err := checkInlineValue("validity", o.Validity); err != nil {
		return "", err
	}
	if err := checkInlineValue("profile", o.Profile); err != nil {
		return "", err
	}
	if o.RecordSales && strings.TrimSpace(o.Profile) == "" {
		return "", errors.New("mikhmon: profile wajib diisi untuk record penjualan")
	}
	if o.PublicURL != "" {
		if err := checkInlineValue("public URL", o.PublicURL); err != nil {
			return "", err
		}
	}
	if o.HookToken != "" {
		if err := checkInlineValue("hook token", o.HookToken); err != nil {
			return "", err
		}
	}

	mode := "0"
	expiry := ""
	if o.Mode != ExpiryNone {
		if o.Validity == "" {
			return "", errors.New("mikhmon: validity wajib diisi untuk mode expiry")
		}
		mode = "X"
		if o.Mode == ExpiryNotice || o.Mode == ExpiryNoticeRecord {
			mode = "N"
		}
		expiry = strings.ReplaceAll(expiryTemplate, "{{validity}}", o.Validity)
	}

	sale, err := saleScript(o)
	if err != nil {
		return "", err
	}
	locks := lockScript(o)
	if expiry == "" && sale == "" && locks == "" {
		return "", nil
	}
	out := onLoginHeaderTemplate +
		`; :local mode "` + mode + `"; { :if ($ucode = "vc" or $ucode = "up" or $comment = "") do={ ` +
		expiry + sale + locks + ` }}`
	return OneLine(out), nil
}

func saleScript(o OnLoginOptions) (string, error) {
	if !o.RecordSales {
		return "", nil
	}
	payload := `:local mac $"mac-address"; :local payload "v2-|-$date-|-$time-|-$user-|-$address-|-$mac-|-` +
		o.Validity + `-|-` + o.Profile + `-|-$saleComment"; `
	fallback := `/system script add name=$payload owner="$monthname$year" source=$date comment=mikhmon; /system script set [find name=$payload] owner="$monthname$year"`
	base := strings.TrimRight(strings.TrimSpace(o.PublicURL), "/")
	if base == "" || strings.TrimSpace(o.HookToken) == "" {
		return payload + fallback + `; `, nil
	}
	cert := "no"
	if o.VerifyTLS {
		cert = "yes"
	}
	fetchMode := "https"
	if strings.HasPrefix(strings.ToLower(base), "http://") {
		fetchMode = "http"
	}
	fetch := fmt.Sprintf(
		`:do { /tool fetch url="%s/api/v1/hook/sales" mode=%s check-certificate=%s http-method=post http-data=$payload http-header-field="Authorization: Bearer %s,Content-Type: text/plain" output=none } on-error={ %s }; `,
		base, fetchMode, cert, o.HookToken, fallback)
	return payload + fetch, nil
}

func lockScript(o OnLoginOptions) string {
	var b strings.Builder
	if o.LockUser {
		b.WriteString(lockUserScript)
	}
	if o.LockServer {
		b.WriteString(lockServerScript)
	}
	return b.String()
}

func checkInlineValue(name, v string) error {
	if strings.ContainsAny(v, "\"\r\n") {
		return fmt.Errorf("mikhmon: %s mengandung karakter yang tidak aman untuk script inline", name)
	}
	return nil
}

// OnLoginMeta is the Mikhmon metadata parsed from an existing on-login script.
type OnLoginMeta struct {
	Mode        string
	Price       string
	Validity    string
	Selling     string
	LockUser    bool
	LockServer  bool
	RecordSales bool
}

// ParseOnLoginScript extracts the legacy :put tuple when it exists. New
// scripts no longer contain price; their profile rows keep that metadata.
func ParseOnLoginScript(script string) OnLoginMeta {
	var meta OnLoginMeta
	if idx := strings.Index(script, `:put (`); idx >= 0 {
		rest := script[idx+len(`:put (`):]
		if end := strings.Index(rest, `)`); end >= 0 {
			tuple := strings.Trim(rest[:end], `"`)
			parts := strings.Split(tuple, ",")
			if len(parts) >= 8 {
				meta.Mode = strings.TrimSpace(parts[1])
				meta.Price = strings.TrimSpace(parts[2])
				meta.Validity = strings.TrimSpace(parts[3])
				meta.Selling = strings.TrimSpace(parts[4])
				meta.LockUser = strings.EqualFold(strings.TrimSpace(parts[6]), "Enable")
				meta.LockServer = strings.EqualFold(strings.TrimSpace(parts[7]), "Enable")
			}
		}
	}
	if meta.Mode == "" {
		meta.Mode = between(script, `:local mode "`, `"`)
	}
	if meta.Validity == "" {
		meta.Validity = between(script, `interval="`, `"`)
	}
	if !meta.LockUser {
		meta.LockUser = strings.Contains(script, "mac-address=$mac")
	}
	if !meta.LockServer {
		meta.LockServer = strings.Contains(script, "server=$srv")
	}
	meta.RecordSales = strings.Contains(script, "/tool fetch") || strings.Contains(script, "comment=mikhmon")
	return meta
}

func between(s, start, end string) string {
	i := strings.Index(s, start)
	if i < 0 {
		return ""
	}
	rest := s[i+len(start):]
	j := strings.Index(rest, end)
	if j < 0 {
		return ""
	}
	return rest[:j]
}
