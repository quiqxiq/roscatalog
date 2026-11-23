package mikhmon

import (
	"fmt"
	"strings"

	"example.com/roscatalog/fields"
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

// MonthText mengembalikan bucket bulan teks (mis. mar2024 atau nov2026).
func (s Sale) MonthText() string {
	t, _, _ := MonthBuckets(s.Date)
	return t
}

// MonthNumeric mengembalikan bucket bulan numerik (mis. 032024 atau 112026).
func (s Sale) MonthNumeric() string {
	_, n, _ := MonthBuckets(s.Date)
	return n
}

// MatchesMonth memeriksa apakah tanggal penjualan berada di bulan yang dicari.
// Parameter month dapat berupa teks (mar2024), numerik (032024), atau ISO (2024-03).
func (s Sale) MatchesMonth(month string) bool {
	tBucket, nBucket, ok := MonthBuckets(month)
	if !ok {
		return false
	}
	myT, myN, myOK := MonthBuckets(s.Date)
	if !myOK {
		return false
	}
	return myT == tBucket || myN == nBucket
}

// MatchesDate memeriksa apakah tanggal penjualan cocok dengan tanggal yang dicari.
// Menerima format ISO (2024-03-18) maupun legacy Mikhmon (mar/18/2024 atau 03/18/2024).
func (s Sale) MatchesDate(date string) bool {
	date = strings.TrimSpace(date)
	if strings.TrimSpace(s.Date) == date {
		return true
	}
	normA := NormalizeDateISO(s.Date)
	normB := NormalizeDateISO(date)
	return normA != "" && normA == normB
}

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

// ToFields mengonversi Sale ke fields.MikhmonSalesAdd untuk disimpan ke /system/script.
func (s Sale) ToFields() fields.MikhmonSalesAdd {
	return fields.MikhmonSalesAdd{
		Name:    s.Name(),
		Source:  s.Source(),
		Owner:   s.Owner(),
		Comment: SalesComment,
	}
}

// SaleFromRow mem-parse fields.MikhmonSalesRow menjadi domain Sale.
func SaleFromRow(row fields.MikhmonSalesRow) (Sale, error) {
	return ParseSaleName(row.Name)
}

var monthNumbers = map[string]string{
	"01": "jan", "02": "feb", "03": "mar", "04": "apr",
	"05": "may", "06": "jun", "07": "jul", "08": "aug",
	"09": "sep", "10": "oct", "11": "nov", "12": "dec",
}

var monthNameToNum = map[string]string{
	"jan": "01", "feb": "02", "mar": "03", "apr": "04",
	"may": "05", "jun": "06", "jul": "07", "aug": "08",
	"sep": "09", "oct": "10", "nov": "11", "dec": "12",
}

var monthNames = map[string]bool{
	"jan": true, "feb": true, "mar": true, "apr": true,
	"may": true, "jun": true, "jul": true, "aug": true,
	"sep": true, "oct": true, "nov": true, "dec": true,
}

// MonthBuckets mengonversi tanggal, bulan teks (mar2024), numerik (032024),
// atau ISO (2024-03) menjadi pasangan text bucket ("mar2024") dan numeric bucket ("032024").
func MonthBuckets(val string) (textBucket, numericBucket string, ok bool) {
	val = strings.TrimSpace(val)
	if val == "" {
		return "", "", false
	}
	// Format ISO date: yyyy-mm-dd
	if len(val) >= 10 && val[4] == '-' && val[7] == '-' {
		year := val[0:4]
		mm := val[5:7]
		if mText, ok := monthNumbers[mm]; ok && isDigits(year) {
			return mText + year, mm + year, true
		}
		return "", "", false
	}
	// Format ISO month: yyyy-mm
	if len(val) == 7 && val[4] == '-' {
		year := val[0:4]
		mm := val[5:7]
		if mText, ok := monthNumbers[mm]; ok && isDigits(year) {
			return mText + year, mm + year, true
		}
		return "", "", false
	}
	// Format slash: mmm/dd/yyyy, mm/dd/yyyy, mmm/yyyy, mm/yyyy
	if strings.Contains(val, "/") {
		parts := strings.Split(val, "/")
		if len(parts) == 3 {
			year := parts[2]
			if len(year) != 4 || !isDigits(year) {
				return "", "", false
			}
			m := strings.ToLower(parts[0])
			if monthNames[m] {
				return m + year, monthNameToNum[m] + year, true
			}
			if len(m) == 1 {
				m = "0" + m
			}
			if mText, ok := monthNumbers[m]; ok {
				return mText + year, m + year, true
			}
			return "", "", false
		}
		if len(parts) == 2 {
			year := parts[1]
			if len(year) != 4 || !isDigits(year) {
				return "", "", false
			}
			m := strings.ToLower(parts[0])
			if monthNames[m] {
				return m + year, monthNameToNum[m] + year, true
			}
			if len(m) == 1 {
				m = "0" + m
			}
			if mText, ok := monthNumbers[m]; ok {
				return mText + year, m + year, true
			}
			return "", "", false
		}
	}
	// Format text bucket: mmmYYYY (mis. mar2024, nov2026) -> len 7
	if len(val) == 7 {
		m := strings.ToLower(val[0:3])
		year := val[3:7]
		if monthNames[m] && isDigits(year) {
			return m + year, monthNameToNum[m] + year, true
		}
	}
	// Format numeric bucket: MMYYYY (mis. 032024, 112026) -> len 6
	if len(val) == 6 && isDigits(val) {
		mm := val[0:2]
		year := val[2:6]
		if mText, ok := monthNumbers[mm]; ok {
			return mText + year, val, true
		}
	}
	return "", "", false
}

func isDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return len(s) > 0
}

// NormalizeDateISO mengonversi tanggal (yyyy-mm-dd atau mmm/dd/yyyy atau mm/dd/yyyy)
// menjadi format ISO yyyy-mm-dd. Mengembalikan string kosong jika tidak valid.
func NormalizeDateISO(date string) string {
	date = strings.TrimSpace(date)
	if len(date) >= 10 && date[4] == '-' && date[7] == '-' {
		mm := date[5:7]
		if _, ok := monthNumbers[mm]; ok {
			return date[0:10]
		}
		return ""
	}
	parts := strings.Split(date, "/")
	if len(parts) == 3 {
		year := parts[2]
		if len(year) != 4 || !isDigits(year) {
			return ""
		}
		day := parts[1]
		if len(day) == 1 {
			day = "0" + day
		}
		m := strings.ToLower(parts[0])
		var mm string
		if monthNames[m] {
			mm = monthNameToNum[m]
		} else {
			if len(m) == 1 {
				m = "0" + m
			}
			if _, ok := monthNumbers[m]; ok {
				mm = m
			} else {
				return ""
			}
		}
		return fmt.Sprintf("%s-%s-%s", year, mm, day)
	}
	return ""
}

// ownerFromDate menerima mmm/dd/yyyy (RouterOS 6) atau yyyy-mm-dd (RouterOS 7).
func ownerFromDate(date string) string {
	t, _, ok := MonthBuckets(date)
	if !ok {
		return ""
	}
	return t
}
