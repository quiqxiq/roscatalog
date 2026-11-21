package mikhmon

import (
	"fmt"
	"strings"

	"example.com/roscatalog/fields"
)

// QuickPrintPreset adalah preset cetak voucher Mikhmon. Source-nya adalah
// string internal Mikhmon, bukan command API.
type QuickPrintPreset struct {
	Package      string
	Server       string
	UserMode     string // vc | up
	UserLength   string
	Prefix       string
	Char         string
	Profile      string
	TimeLimit    string
	DataLimit    string
	Comment      string
	Validity     string
	Price        string
	SellingPrice string
	UserLock     string
}

// Name adalah nama script di /system/script.
func (p QuickPrintPreset) Name() string {
	pkg := strings.TrimSpace(p.Package)
	if pkg == "" {
		pkg = "Preset"
	}
	return "Quick_Print_" + strings.Join(strings.Fields(pkg), "-")
}

// Source menyusun string internal QuickPrint. SellingPrice kosong memakai
// Price, sama seperti post_quick_print.php.
func (p QuickPrintPreset) Source() string {
	price := p.Price
	if price == "" {
		price = "0"
	}
	selling := p.SellingPrice
	if selling == "" {
		selling = price
	}
	parts := []string{
		p.Package, p.Server, p.UserMode, p.UserLength, p.Prefix, p.Char,
		p.Profile, p.TimeLimit, p.DataLimit, p.Comment, p.Validity,
		price + "_" + selling, p.UserLock,
	}
	return "#" + strings.Join(parts, "#")
}

// ToFields mengonversi QuickPrintPreset ke fields.MikhmonQuickPrintAdd.
func (p QuickPrintPreset) ToFields() fields.MikhmonQuickPrintAdd {
	return fields.MikhmonQuickPrintAdd{
		Name:    p.Name(),
		Source:  p.Source(),
		Comment: QuickPrintComment,
	}
}

// PresetFromRow mem-parse fields.MikhmonQuickPrintRow menjadi domain QuickPrintPreset.
func PresetFromRow(row fields.MikhmonQuickPrintRow) (QuickPrintPreset, error) {
	return ParseQuickPrintSource(row.Source)
}

// ParseQuickPrintSource membalikkan QuickPrintPreset.Source.
func ParseQuickPrintSource(source string) (QuickPrintPreset, error) {
	if !strings.HasPrefix(source, "#") {
		return QuickPrintPreset{}, fmt.Errorf("mikhmon: source QuickPrint harus diawali '#'")
	}
	parts := strings.Split(source[1:], "#")
	if len(parts) != 13 {
		return QuickPrintPreset{}, fmt.Errorf("mikhmon: source QuickPrint memiliki %d field, butuh 13", len(parts))
	}
	price, selling, ok := strings.Cut(parts[11], "_")
	if !ok {
		return QuickPrintPreset{}, fmt.Errorf("mikhmon: field harga QuickPrint %q tidak berisi price_sprice", parts[11])
	}
	return QuickPrintPreset{
		Package: parts[0], Server: parts[1], UserMode: parts[2],
		UserLength: parts[3], Prefix: parts[4], Char: parts[5],
		Profile: parts[6], TimeLimit: parts[7], DataLimit: parts[8],
		Comment: parts[9], Validity: parts[10], Price: price,
		SellingPrice: selling, UserLock: parts[12],
	}, nil
}
