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
	// ExpiryNone hanya menulis metadata harga/validity, tanpa logika expiry.
	ExpiryNone ExpiryMode = "0"
	// ExpiryRemove menghapus user saat kedaluwarsa.
	ExpiryRemove ExpiryMode = "rem"
	// ExpiryNotice menandai limit-uptime=1s dan memutus sesi saat kedaluwarsa.
	ExpiryNotice ExpiryMode = "ntf"
	// ExpiryRemoveRecord sama seperti ExpiryRemove, plus catatan penjualan.
	ExpiryRemoveRecord ExpiryMode = "remc"
	// ExpiryNoticeRecord sama seperti ExpiryNotice, plus catatan penjualan.
	ExpiryNoticeRecord ExpiryMode = "ntfc"
)

// OnLoginOptions adalah input untuk membangun properti on-login profile.
type OnLoginOptions struct {
	Mode         ExpiryMode
	Price        string // harga modal, mis. "5000" (default "0")
	Validity     string // mis. "1d"; wajib untuk mode expiry
	SellingPrice string // harga jual, mis. "7000" (default "0")
	LockUser     bool   // kunci MAC user saat login pertama
	LockServer   bool   // kunci server hotspot saat login pertama
	Profile      string // nama profile, dipakai di nama script laporan
	RecordSales  bool   // catat penjualan (otomatis true untuk remc/ntfc)
}

const lockUserScript = `; [:local mac $"mac-address"; /ip hotspot user set mac-address=$mac [find where name=$user]]`
const lockServerScript = `; [:local mac $"mac-address"; :local srv [/ip hotspot host get [find where mac-address="$mac"] server]; /ip hotspot user set server=$srv [find where name=$user]]`

const onLoginLogicTemplate = `
:local date [ /system clock get date ];
:local time [ /system clock get time ];
:local mnames "janfebmaraprmayjunjulaugsepoctnovdec";
:local year;
:local monthname;
:local day;
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
:local legacy "$monthname/$day/$year";
:local comment [ /ip hotspot user get [ /ip hotspot user find where name="$user" ] comment ];
:local ucode [:pick $comment 0 2];
:if ($ucode = "vc" or $ucode = "up" or $comment = "") do={
  /sys sch add name="$user" disable=no start-date=$date interval="{{validity}}";
  :delay 2s;
  :local exp [ /sys sch get [ /sys sch find where name="$user" ] next-run ];
  :local expLegacy;
  :local getxp [len $exp];
  :if ($getxp = 8) do={
    :set expLegacy "$legacy $exp";
  } else={
    :if ($getxp = 15) do={
      :local d [:pick $exp 0 6];
      :local t [:pick $exp 7 16];
      :set expLegacy "$d/$year $t";
    } else={
      :if ([:pick $exp 4 5] = "-") do={
        :local ey [:pick $exp 0 4];
        :local ed [:pick $exp 8 10];
        :local mnum2 [:tonum [:pick $exp 5 7]];
        :local emn [:pick $mnames (($mnum2 - 1) * 3) (($mnum2 - 1) * 3 + 3)];
        :local et [:pick $exp 11 [:len $exp]];
        :set expLegacy "$emn/$ed/$ey $et";
      } else={
        :set expLegacy $exp;
      }
    }
  };
  /ip hotspot user set comment="$expLegacy" [find where name="$user"];
  /sys sch remove [find where name="$user"];
`

const saleRecordTemplate = `; :local mac $"mac-address"; :local time [/system clock get time ]; :local sname "$legacy-|-$time-|-$user-|-{{price}}-|-$address-|-$mac-|-{{validity}}-|-{{profile}}-|-$comment"; /system script add name=$sname owner="$monthname$year" source=$legacy comment=mikhmon; /system script set [find name=$sname] owner="$monthname$year"`

const profileMonitorTemplate = `
:local normdate do={
  :local mnames "janfebmaraprmayjunjulaugsepoctnovdec";
  :local v $d;
  :if ([:pick $v 4 5] = "-") do={
    :local mn [:tonum [:pick $v 5 7]];
    :local m [:pick $mnames (($mn - 1) * 3) (($mn - 1) * 3 + 3)];
    :local rest "";
    :if ([:len $v] > 10) do={:set rest [:pick $v 10 [:len $v]];};
    :local mm [:pick $v 8 10];
    :local yy [:pick $v 0 4];
    :set v "$m/$mm/$yy$rest";
  };
  :return $v;
};
:local dateint do={
  :local mnames "janfebmaraprmayjunjulaugsepoctnovdec";
  :local days [:pick $d 4 6];
  :local month [:pick $d 0 3];
  :local year [:pick $d 7 11];
  :local monthint ([:find $mnames $month] / 3);
  :local monthnum ($monthint + 1);
  :if ([:len $monthnum] = 1) do={
    :local zero "0";
    :return [:tonum "$year$zero$monthnum$days"];
  } else={
    :return [:tonum "$year$monthnum$days"];
  };
};
:local timeint do={
  :local hours [:pick $t 0 2];
  :local minutes [:pick $t 3 5];
  :return ($hours * 60 + $minutes);
};
:local date [/system clock get date];
:local time [/system clock get time];
:local todaydate [$normdate d=$date];
:local today [$dateint d=$todaydate];
:local curtime [$timeint t=$time];
:foreach i in [/ip hotspot user find where profile="{{profile_name}}"] do={
  :local comment [/ip hotspot user get $i comment];
  :local name [/ip hotspot user get $i name];
  :local norm [$normdate d=$comment];
  :if ([:pick $norm 3] = "/" and [:pick $norm 6] = "/" and [:len $norm] >= 19) do={
    :local gettime [:pick $norm 12 20];
    :local expd [$dateint d=$norm];
    :local expt [$timeint t=$gettime];
    :if ($expd < $today or ($expd = $today and $expt < $curtime)) do={
      [ /ip hotspot user {{mode}} $i ];
      [ /ip hotspot active remove [find where user=$name] ];
    };
  };
};
`

//go:embed expire_monitor.ros
var expireMonitorScript string

// ExpireMonitorScript mengembalikan on-event scheduler global Mikhmon-Expire-Monitor.
func ExpireMonitorScript() string {
	return strings.TrimSpace(expireMonitorScript)
}

// ProfileMonitorScript menghasilkan on-event scheduler Monitor Profile untuk satu
// user profile tertentu, persis seperti mikhmon_v3_bgservice_script di Mikhmon v4.
func ProfileMonitorScript(profile string, mode ExpiryMode) (string, error) {
	if err := checkInlineValue("profile", profile); err != nil {
		return "", err
	}
	var modeCmd string
	switch mode {
	case ExpiryRemove, ExpiryRemoveRecord:
		modeCmd = "remove"
	case ExpiryNotice, ExpiryNoticeRecord:
		modeCmd = "set limit-uptime=1s"
	case ExpiryNone:
		return "", nil
	default:
		return "", fmt.Errorf("mikhmon: expmode %q tidak dikenal", mode)
	}
	s := strings.ReplaceAll(profileMonitorTemplate, "{{profile_name}}", profile)
	s = strings.ReplaceAll(s, "{{mode}}", modeCmd)
	return OneLine(s), nil
}

// OneLine merapatkan semua whitespace menjadi satu spasi.
func OneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// OnLoginScript membangun properti on-login satu baris persis seperti
// mikhmon_v3_onlogin_script di Mikhmon v4.
func OnLoginScript(o OnLoginOptions) (string, error) {
	switch o.Mode {
	case ExpiryNone, ExpiryRemove, ExpiryNotice, ExpiryRemoveRecord, ExpiryNoticeRecord:
	default:
		return "", fmt.Errorf("mikhmon: expmode %q tidak dikenal", o.Mode)
	}
	price := o.Price
	if price == "" {
		price = "0"
	}
	selling := o.SellingPrice
	if selling == "" {
		selling = "0"
	}
	if err := checkInlineValue("price", price); err != nil {
		return "", err
	}
	if err := checkInlineValue("selling price", selling); err != nil {
		return "", err
	}
	if err := checkInlineValue("validity", o.Validity); err != nil {
		return "", err
	}
	if err := checkInlineValue("profile", o.Profile); err != nil {
		return "", err
	}

	getlock := "Disable"
	if o.LockUser {
		getlock = "Enable"
	}
	locks := lockScript(o)

	if o.Mode == ExpiryNone {
		put := fmt.Sprintf(`:put (",,%s,,,noexp,%s,")`, price, getlock)
		return OneLine(put + locks), nil
	}

	if o.Validity == "" {
		return "", errors.New("mikhmon: validity wajib diisi untuk mode expiry")
	}

	recordSales := o.Mode == ExpiryRemoveRecord || o.Mode == ExpiryNoticeRecord || o.RecordSales
	if recordSales && strings.TrimSpace(o.Profile) == "" {
		return "", errors.New("mikhmon: profile wajib diisi untuk record penjualan")
	}

	put := fmt.Sprintf(`:put (",%s,%s,%s,%s,,%s,")`, o.Mode, price, o.Validity, selling, getlock)

	logic := strings.ReplaceAll(onLoginLogicTemplate, "{{validity}}", o.Validity)

	record := ""
	if recordSales {
		record = strings.NewReplacer(
			"{{price}}", price,
			"{{validity}}", o.Validity,
			"{{profile}}", o.Profile,
		).Replace(saleRecordTemplate)
	}

	out := put + `; {` + logic + record + locks + `}}`
	return OneLine(out), nil
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

// OnLoginMeta adalah metadata Mikhmon yang diparsing dari script on-login yang ada.
type OnLoginMeta struct {
	Mode        string
	Price       string
	Validity    string
	Selling     string
	LockUser    bool
	LockServer  bool
	RecordSales bool
}

// ParseOnLoginScript mengurai tuple :put metadata dari script on-login Mikhmon.
func ParseOnLoginScript(script string) OnLoginMeta {
	var meta OnLoginMeta
	if idx := strings.Index(script, `:put (`); idx >= 0 {
		rest := script[idx+len(`:put (`):]
		if end := strings.Index(rest, `)`); end >= 0 {
			tuple := strings.Trim(rest[:end], `"`)
			parts := strings.Split(tuple, ",")
			if len(parts) >= 7 {
				meta.Mode = strings.TrimSpace(parts[1])
				meta.Price = strings.TrimSpace(parts[2])
				meta.Validity = strings.TrimSpace(parts[3])
				meta.Selling = strings.TrimSpace(parts[4])
				meta.LockUser = strings.EqualFold(strings.TrimSpace(parts[6]), "Enable")
				if len(parts) >= 8 {
					meta.LockServer = strings.EqualFold(strings.TrimSpace(parts[7]), "Enable")
				}
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
	meta.RecordSales = meta.Mode == string(ExpiryRemoveRecord) ||
		meta.Mode == string(ExpiryNoticeRecord) ||
		strings.Contains(script, "comment=mikhmon")
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
