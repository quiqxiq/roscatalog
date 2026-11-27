package roscatalog

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// File ini berisi helper opsional untuk memakai struct Go sebagai daftar field.
// Struct tidak pernah diwajibkan: Params.Raw mengembalikan semua kolom tanpa
// struct, dan Params.Proplist menerima daftar nama apa adanya. Helper di sini
// hanya kemudahan untuk pemakai yang ingin tipe Go.

// rosTag memecah tag `ros:"nama,opsi"` menjadi nama properti dan himpunan opsi.
// Nama kosong fallback ke bentuk kebab-case dari nama field Go.
func rosTag(f reflect.StructField) (string, map[string]bool, bool) {
	tag, ok := f.Tag.Lookup("ros")
	if !ok {
		return kebabName(f.Name), nil, true
	}
	if tag == "-" {
		return "", nil, false
	}
	parts := strings.Split(tag, ",")
	name := parts[0]
	opts := make(map[string]bool, len(parts)-1)
	for _, p := range parts[1:] {
		opts[p] = true
	}
	if name == "" {
		name = kebabName(f.Name)
	}
	return name, opts, true
}

// kebabName mengubah FieldName menjadi field-name. Dipakai hanya sebagai
// fallback bila tag ros tidak mencantumkan nama.
func kebabName(s string) string {
	var b strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 && (unicode.IsLower(runes[i-1]) || unicode.IsDigit(runes[i-1]) ||
				(i+1 < len(runes) && unicode.IsLower(runes[i+1]) && unicode.IsUpper(runes[i-1]))) {
				b.WriteByte('-')
			}
			b.WriteRune(unicode.ToLower(r))
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func indirectType(t reflect.Type) reflect.Type {
	for t != nil && t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

// FieldNames mengembalikan nama properti RouterOS dari tag `ros:` pada T,
// dalam urutan deklarasi. Struct bisa bertingkat (anonymous embed).
func FieldNames[T any]() []string {
	t := indirectType(reflect.TypeOf((*T)(nil)).Elem())
	if t == nil || t.Kind() != reflect.Struct {
		return nil
	}
	out := []string{}
	seen := map[string]bool{}
	var walk func(reflect.Type)
	walk = func(t reflect.Type) {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.PkgPath != "" { // field unexported
				continue
			}
			ft := indirectType(f.Type)
			if f.Anonymous && ft != nil && ft.Kind() == reflect.Struct {
				walk(ft)
				continue
			}
			name, _, use := rosTag(f)
			if !use || name == "" || seen[name] {
				continue
			}
			seen[name] = true
			out = append(out, name)
		}
	}
	walk(t)
	return out
}

// Decode mengisi struct T dari satu baris hasil print. Field yang tidak
// dikirim router dibiarkan kosong — itu normal (properti kosong tidak dikirim).
func Decode[T any](row map[string]string) (T, error) {
	var out T
	v := reflect.ValueOf(&out).Elem()
	if err := decodeInto(v, row); err != nil {
		return out, err
	}
	return out, nil
}

// DecodeAll mengubah seluruh baris menjadi struct.
func DecodeAll[T any](rows []map[string]string) ([]T, error) {
	out := make([]T, 0, len(rows))
	for _, row := range rows {
		v, err := Decode[T](row)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}

func decodeInto(v reflect.Value, row map[string]string) error {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("Decode: butuh struct, dapat %s", v.Type())
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fv := v.Field(i)
		if f.PkgPath != "" || !fv.CanSet() {
			continue
		}
		ft := indirectType(f.Type)
		if f.Anonymous && ft != nil && ft.Kind() == reflect.Struct {
			if err := decodeInto(fv, row); err != nil {
				return err
			}
			continue
		}
		name, _, use := rosTag(f)
		if !use || name == "" {
			continue
		}
		val, ok := row[name]
		if !ok || fv.Kind() != reflect.String {
			continue
		}
		fv.SetString(val)
	}
	return nil
}

// ArgsOf mengubah struct Add/Set menjadi argumen untuk Action.SentenceWith.
// Field `ros:",omitempty"` yang kosong dilewati; field `,required` yang kosong
// menghasilkan error. Field bertanda `ros:"-"` diabaikan.
func ArgsOf[T any](value T) (map[string]string, error) {
	v := reflect.ValueOf(value)
	if !v.IsValid() {
		return nil, nil
	}
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil, nil
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ArgsOf: butuh struct, dapat %s", v.Type())
	}
	out := map[string]string{}
	if err := argsOfValue(v, out); err != nil {
		return nil, err
	}
	return out, nil
}

func argsOfValue(v reflect.Value, out map[string]string) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fv := v.Field(i)
		if f.PkgPath != "" || !fv.IsValid() {
			continue
		}
		ft := indirectType(f.Type)
		if f.Anonymous && ft != nil && ft.Kind() == reflect.Struct {
			if err := argsOfValue(fv, out); err != nil {
				return err
			}
			continue
		}
		if fv.Kind() == reflect.Pointer && fv.IsNil() {
			continue
		}
		for fv.Kind() == reflect.Pointer {
			fv = fv.Elem()
		}
		name, opts, use := rosTag(f)
		if !use || name == "" {
			continue
		}
		var s string
		var empty bool
		switch fv.Kind() {
		case reflect.String:
			s = fv.String()
			empty = s == ""
		case reflect.Bool:
			if fv.Bool() {
				s = "yes"
			} else {
				s = "no"
				empty = true
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			empty = fv.IsZero()
			s = fmt.Sprint(fv.Interface())
		default:
			continue
		}
		if empty && opts["required"] {
			return fmt.Errorf("%w: field %q is required", ErrMissingArgument, f.Name)
		}
		if empty && opts["omitempty"] {
			continue
		}
		out[name] = s
	}
	return nil
}
