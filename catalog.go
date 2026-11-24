package roscatalog

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Catalog adalah kumpulan resource yang bisa diperluas pemakai. Catalog
// bawaan berisi katalog kurasi; pemakai boleh membuat catalog sendiri dan
// menambahkan menu/command yang tidak ada di bawaan tanpa mengubah global.
type Catalog struct {
	resources []*Resource
	byID      map[string]*Resource
}

// NewCatalog membuat catalog dari resource awal. ID duplikat: yang terakhir
// menang, karena ini konstruktor; gunakan Add untuk penanganan error.
func NewCatalog(base ...*Resource) *Catalog {
	c := &Catalog{byID: map[string]*Resource{}}
	for _, r := range base {
		if r == nil || strings.TrimSpace(r.ID) == "" {
			continue
		}
		if _, exists := c.byID[r.ID]; !exists {
			c.resources = append(c.resources, r)
		}
		c.byID[r.ID] = r
	}
	return c
}

// Builtin mengembalikan catalog bawaan (resource yang di-register saat init).
func Builtin() *Catalog { return NewCatalog(registry...) }

// Add menambahkan resource setelah validasi. Tidak ada perubahan bila salah
// satu resource gagal divalidasi.
func (c *Catalog) Add(rs ...*Resource) error {
	if c == nil {
		return fmt.Errorf("%w: catalog nil", ErrInvalidResource)
	}
	for _, r := range rs {
		if err := c.checkAdd(r); err != nil {
			return err
		}
	}
	for _, r := range rs {
		c.resources = append(c.resources, r)
		c.byID[r.ID] = r
	}
	return nil
}

func (c *Catalog) checkAdd(r *Resource) error {
	if r == nil {
		return fmt.Errorf("%w: resource nil", ErrInvalidResource)
	}
	if _, exists := c.byID[r.ID]; exists {
		return fmt.Errorf("%w: %s", ErrDuplicateResource, r.ID)
	}
	return r.Validate()
}

// Replace mengganti resource dengan ID yang sama.
func (c *Catalog) Replace(r *Resource) error {
	if c == nil {
		return fmt.Errorf("%w: catalog nil", ErrInvalidResource)
	}
	if r == nil {
		return fmt.Errorf("%w: resource nil", ErrInvalidResource)
	}
	if _, exists := c.byID[r.ID]; !exists {
		return fmt.Errorf("%w: %s", ErrResourceNotFound, r.ID)
	}
	if err := r.Validate(); err != nil {
		return err
	}
	for i, old := range c.resources {
		if old.ID == r.ID {
			c.resources[i] = r
			break
		}
	}
	c.byID[r.ID] = r
	return nil
}

// Extend menambahkan seluruh isi catalog lain. Atomik: tidak ada yang
// ditambahkan bila satu ID duplikat atau tidak valid.
func (c *Catalog) Extend(other *Catalog) error {
	if other == nil {
		return nil
	}
	for _, r := range other.resources {
		if err := c.checkAdd(r); err != nil {
			return err
		}
	}
	for _, r := range other.resources {
		c.resources = append(c.resources, r)
		c.byID[r.ID] = r
	}
	return nil
}

// Find mencari resource dan mengembalikan ErrResourceNotFound bila tidak ada.
func (c *Catalog) Find(id string) (*Resource, error) {
	if c == nil || c.byID == nil {
		return nil, fmt.Errorf("%w: %s", ErrResourceNotFound, id)
	}
	r, ok := c.byID[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrResourceNotFound, id)
	}
	return r, nil
}

// Lookup mencari resource dengan pola boolean.
func (c *Catalog) Lookup(id string) (*Resource, bool) {
	if c == nil || c.byID == nil {
		return nil, false
	}
	r, ok := c.byID[id]
	return r, ok
}

// All mengembalikan semua resource dalam urutan tampil.
func (c *Catalog) All() []*Resource {
	if c == nil {
		return nil
	}
	return sortedResources(c.resources)
}

// ByGroup mengembalikan resource dalam satu grup.
func (c *Catalog) ByGroup(g Group) []*Resource {
	if c == nil {
		return nil
	}
	var out []*Resource
	for _, r := range c.All() {
		if r.Group == g {
			out = append(out, r)
		}
	}
	return out
}

// Len mengembalikan jumlah resource.
func (c *Catalog) Len() int {
	if c == nil {
		return 0
	}
	return len(c.resources)
}

// Validate memeriksa seluruh catalog, termasuk ID duplikat.
func (c *Catalog) Validate() error {
	if c == nil {
		return fmt.Errorf("%w: catalog nil", ErrInvalidResource)
	}
	var errs []error
	seen := map[string]bool{}
	for _, r := range c.resources {
		if r == nil {
			errs = append(errs, fmt.Errorf("%w: resource nil", ErrInvalidResource))
			continue
		}
		if seen[r.ID] {
			errs = append(errs, fmt.Errorf("%w: %s", ErrDuplicateResource, r.ID))
		}
		seen[r.ID] = true
		if err := r.Validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// Validate memeriksa konsistensi satu resource sebelum dikirim ke router.
func (r *Resource) Validate() error {
	if r == nil {
		return fmt.Errorf("%w: resource nil", ErrInvalidResource)
	}
	var errs []error
	if strings.TrimSpace(r.ID) == "" {
		errs = append(errs, fmt.Errorf("%w: ID kosong", ErrInvalidResource))
	}
	if strings.TrimSpace(r.Path) == "" || !strings.HasPrefix(r.Path, "/") {
		errs = append(errs, fmt.Errorf("%w: %s: path harus diawali '/'", ErrInvalidResource, r.ID))
	}
	if strings.ContainsAny(r.Path, " \t\r\n") {
		errs = append(errs, fmt.Errorf("%w: %s: path mengandung spasi", ErrInvalidResource, r.ID))
	}
	if len(r.Views) == 0 && len(r.Actions) == 0 {
		errs = append(errs, fmt.Errorf("%w: %s: tidak punya view/action", ErrInvalidResource, r.ID))
	}

	views := map[string]bool{}
	for _, v := range r.Views {
		if strings.TrimSpace(v.Name) == "" {
			errs = append(errs, fmt.Errorf("%w: %s: view tanpa nama", ErrInvalidResource, r.ID))
			continue
		}
		if views[v.Name] {
			errs = append(errs, fmt.Errorf("%w: %s: view %q duplikat", ErrInvalidResource, r.ID, v.Name))
		}
		views[v.Name] = true
		if v.Mode == "" {
			errs = append(errs, fmt.Errorf("%w: %s/%s: mode kosong", ErrInvalidResource, r.ID, v.Name))
		}
		switch v.Mode {
		case ModeOneShot, ModeNative, ModeInterval, ModeFollow, ModeListen:
		default:
			errs = append(errs, fmt.Errorf("%w: %s/%s: mode tidak dikenal %q", ErrInvalidResource, r.ID, v.Name, v.Mode))
		}
		if v.Command == "" {
			errs = append(errs, fmt.Errorf("%w: %s/%s: command kosong", ErrInvalidResource, r.ID, v.Name))
		} else if v.Command != r.Path && !strings.HasPrefix(v.Command, r.Path+"/") {
			errs = append(errs, fmt.Errorf("%w: %s/%s: command %q tidak di bawah %q",
				ErrInvalidResource, r.ID, v.Name, v.Command, r.Path))
		}
	}

	acts := map[string]bool{}
	for _, a := range r.Actions {
		if strings.TrimSpace(a.Name) == "" {
			errs = append(errs, fmt.Errorf("%w: %s: action tanpa nama", ErrInvalidResource, r.ID))
			continue
		}
		if acts[a.Name] {
			errs = append(errs, fmt.Errorf("%w: %s: action %q duplikat", ErrInvalidResource, r.ID, a.Name))
		}
		acts[a.Name] = true
		want := r.Path + "/" + a.Name
		if a.Command == "" {
			errs = append(errs, fmt.Errorf("%w: %s/%s: command kosong", ErrInvalidResource, r.ID, a.Name))
		} else if a.Command != want {
			errs = append(errs, fmt.Errorf("%w: %s/%s: command %q, want %q",
				ErrInvalidResource, r.ID, a.Name, a.Command, want))
		}
	}
	return errors.Join(errs...)
}

func sortedResources(in []*Resource) []*Resource {
	out := append([]*Resource(nil), in...)
	rank := make(map[Group]int, len(GroupOrder))
	for i, g := range GroupOrder {
		rank[g] = i
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Group != out[j].Group {
			return rank[out[i].Group] < rank[out[j].Group]
		}
		return out[i].ID < out[j].ID
	})
	return out
}
