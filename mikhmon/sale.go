package mikhmon

import (
	"fmt"
	"strings"
)

// Sale adalah satu baris laporan penjualan Mikhmon. Name mengikuti konvensi
// tanggal-|-jam-|-user-|-harga-|-address-|-mac-|-validity-|-profile-|-comment.
type Sale struct {
	Date     string // tanggal legacy mmm/dd/yyyy (atau yyyy-mm-dd)
	Time     string // jam HH:MM:SS
	User     string
	Price    string
	Address  string
	MAC      string
	Validity string
	Profile  string
	Comment  string // comment user, mis. vc-123-11.18.26
}

// Name menyusun nama script laporan persis seperti Mikhmon.
func (s Sale) Name() string {
	return strings.Join([]string{
		s.Date, s.Time, s.User, s.Price, s.Address, s.MAC,
		s.Validity, s.Profile, s.Comment,
	}, "-|-")
}

// Source adalah properti source script, yaitu tanggal legacy.
func (s Sale) Source() string { return s.Date }

// Owner adalah bulan+tahun yang dipakai Mikhmon, mis. nov2026.
func (s Sale) Owner() string { return ownerFromDate(s.Date) }

// Validate memeriksa field yang wajib ada agar laporan bisa ditampilkan
// kembali oleh Mikhmon.
func (s Sale) Validate() error {
	for name, v := range map[string]string{
		"date": s.Date, "time": s.Time, "user": s.User,
		"price": s.Price, "validity": s.Validity, "profile": s.Profile,
	} {
		if strings.TrimSpace(v) == "" {
			return fmt.Errorf("mikhmon: sale.%s wajib diisi", name)
		}
	}
	if s.Owner() == "" {
		return fmt.Errorf("mikhmon: format sale.Date %q tidak dikenali", s.Date)
	}
	return nil
}

// ParseSaleName membalikkan Sale.Name. Comment boleh mengandung "-|-"; field
// terakhir menampung sisanya.
func ParseSaleName(name string) (Sale, error) {
	parts := strings.SplitN(name, "-|-", 9)
	if len(parts) != 9 {
		return Sale{}, fmt.Errorf("mikhmon: nama laporan %q tidak memiliki 9 field", name)
	}
	return Sale{
		Date: parts[0], Time: parts[1], User: parts[2], Price: parts[3],
		Address: parts[4], MAC: parts[5], Validity: parts[6],
		Profile: parts[7], Comment: parts[8],
	}, nil
}

// PendingSale is a sale recorded by a router when the hook was unreachable.
// It intentionally does not contain price; the backend resolves price from
// hotspot_profiles when the record is drained.
type PendingSale struct {
	Date     string
	Time     string
	User     string
	Address  string
	MAC      string
	Validity string
	Profile  string
	Comment  string
}

// Name returns the v2 fallback script name.
func (s PendingSale) Name() string {
	return strings.Join([]string{
		"v2", s.Date, s.Time, s.User, s.Address, s.MAC,
		s.Validity, s.Profile, s.Comment,
	}, "-|-")
}

// Source is the legacy date used by the sync owner scan.
func (s PendingSale) Source() string { return s.Date }

// Validate checks fields required to resolve the sale.
func (s PendingSale) Validate() error {
	for name, v := range map[string]string{
		"date": s.Date, "time": s.Time, "user": s.User, "profile": s.Profile,
	} {
		if strings.TrimSpace(v) == "" {
			return fmt.Errorf("mikhmon: pending sale.%s wajib diisi", name)
		}
	}
	if ownerFromDate(s.Date) == "" {
		return fmt.Errorf("mikhmon: format pending sale.Date %q tidak dikenali", s.Date)
	}
	return nil
}

// ParsePendingSaleName reverses PendingSale.Name. Comment may contain "-|-".
func ParsePendingSaleName(name string) (PendingSale, error) {
	parts := strings.SplitN(name, "-|-", 9)
	if len(parts) != 9 || parts[0] != "v2" {
		return PendingSale{}, fmt.Errorf("mikhmon: nama pending sale %q tidak valid", name)
	}
	return PendingSale{
		Date: parts[1], Time: parts[2], User: parts[3], Address: parts[4],
		MAC: parts[5], Validity: parts[6], Profile: parts[7], Comment: parts[8],
	}, nil
}

var monthNumbers = map[string]string{
	"01": "jan", "02": "feb", "03": "mar", "04": "apr",
	"05": "may", "06": "jun", "07": "jul", "08": "aug",
	"09": "sep", "10": "oct", "11": "nov", "12": "dec",
}

var monthNames = map[string]bool{
	"jan": true, "feb": true, "mar": true, "apr": true,
	"may": true, "jun": true, "jul": true, "aug": true,
	"sep": true, "oct": true, "nov": true, "dec": true,
}

// ownerFromDate menerima mmm/dd/yyyy (RouterOS 6) atau yyyy-mm-dd (RouterOS 7).
func ownerFromDate(date string) string {
	date = strings.TrimSpace(date)
	if len(date) >= 10 && date[4] == '-' && date[7] == '-' {
		month := monthNumbers[date[5:7]]
		if month != "" {
			return month + date[0:4]
		}
		return ""
	}
	parts := strings.Split(date, "/")
	if len(parts) == 3 {
		month := strings.ToLower(parts[0])
		if monthNames[month] && len(parts[2]) == 4 {
			return month + parts[2]
		}
	}
	return ""
}
