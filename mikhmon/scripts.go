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
	Price        string // harga modal; kosong menjadi "0"
	Validity     string // mis. "1d"
	SellingPrice string // harga jual; kosong menjadi "0"
	LockUser     bool   // kunci MAC user saat login pertama
	LockServer   bool   // kunci server hotspot saat login pertama
	Profile      string // nama profile, dipakai di nama script laporan
}

const lockUserScript = `; :local mac $"mac-address"; /ip hotspot user set mac-address=$mac [find where name=$user]`
const lockServerScript = `; :local mac $"mac-address"; :local srv [/ip hotspot host get [find where mac-address="$mac"] server]; /ip hotspot user set server=$srv [find where name=$user]`

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
:local ucode [:pic $comment 0 2];
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

const saleRecordTemplate = `; :local mac $"mac-address"; :local time [/system clock get time ]; :local sname "$legacy-|-$time-|-$user-|-{{price}}-|-$address-|-$mac-|-{{validity}}-|-{{profile}}-|-$comment"; /system script add name=$sname owner="$monthname$year" source=$legacy comment=mikhmon; /system script set [find name=$sname] owner="$monthname$year"`

//go:embed expire_monitor.ros
var expireMonitorScript string

// ExpireMonitorScript mengembalikan on-event scheduler Mikhmon-Expire-Monitor.
// String ini disalin dari assets/js/func.js (setExpMon) dan sengaja satu baris
// agar aman dikirim sebagai satu parameter API.
func ExpireMonitorScript() string {
	return strings.TrimSpace(expireMonitorScript)
}

// OneLine merapatkan semua whitespace menjadi satu spasi. Dipakai untuk
// memastikan script RouterOS tidak mengandung newline saat dikirim lewat API.
func OneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// OnLoginScript membangun properti on-login satu baris dengan perilaku yang
// sama seperti mikhmon_onlogin_script() dari Mikhmon v4:
//
//   - metadata harga/validity ditulis lewat :put;
//   - scheduler sementara dipakai untuk menghitung next-run;
//   - comment user menyimpan tanggal kedaluwarsa dan mode N/X;
//   - mode remc/ntfc menambahkan script laporan penjualan;
//   - lock MAC/server diterapkan saat expiry pertama kali dihitung.
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

	put := fmt.Sprintf(`:put (",%s,%s,%s,%s,,%s,%s,")`,
		o.Mode, price, o.Validity, selling, enabledFlag(o.LockUser), enabledFlag(o.LockServer))

	if o.Mode == ExpiryNone {
		return OneLine(put + lockScript(o)), nil
	}
	if o.Validity == "" {
		return "", errors.New("mikhmon: validity wajib diisi untuk mode expiry")
	}
	if (o.Mode == ExpiryRemoveRecord || o.Mode == ExpiryNoticeRecord) && o.Profile == "" {
		return "", errors.New("mikhmon: profile wajib diisi untuk mode remc/ntfc")
	}
	if err := checkInlineValue("profile", o.Profile); err != nil {
		return "", err
	}

	mode := "X"
	if o.Mode == ExpiryNotice || o.Mode == ExpiryNoticeRecord {
		mode = "N"
	}

	logic := strings.ReplaceAll(onLoginLogicTemplate, "{{validity}}", o.Validity)
	record := ""
	if o.Mode == ExpiryRemoveRecord || o.Mode == ExpiryNoticeRecord {
		record = strings.NewReplacer(
			"{{price}}", price,
			"{{validity}}", o.Validity,
			"{{profile}}", o.Profile,
		).Replace(saleRecordTemplate)
	}

	out := put + `; :local mode "` + mode + `"; {` + logic + record + lockScript(o) + `}}`
	return OneLine(out), nil
}

func enabledFlag(v bool) string {
	if v {
		return "Enable"
	}
	return "Disable"
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
