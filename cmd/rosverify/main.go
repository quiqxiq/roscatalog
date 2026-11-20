// Command rosverify memverifikasi katalog roscatalog terhadap router sungguhan
// memakai /console/inspect lewat goros/v4. Ia juga bisa dipakai untuk
// discovery menu/argumen sebelum menambah resource baru.
//
// Verifikasi:
//
//	go run ./cmd/rosverify -addr 192.168.88.1:8728 -user admin -pass secret
//	go run ./cmd/rosverify -resource hotspot-user -all
//	go run ./cmd/rosverify -json -strict > schema-chr-7.24.6.json
//
// Discovery:
//
//	go run ./cmd/rosverify -discover ip,hotspot,active
//	go run ./cmd/rosverify -discover ip,hotspot,active -verb print -completion
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	ros "example.com/roscatalog"

	routeros "github.com/quiqxiq/goros/v4"
	"github.com/quiqxiq/goros/v4/schema"
	"github.com/quiqxiq/goros/v4/transport"
)

const (
	statusOK      = "ok"
	statusWarn    = "warn"
	statusDrift   = "drift"
	statusUnknown = "unknown"
	statusError   = "error"
)

// Check adalah satu hasil verifikasi view/aksi.
type Check struct {
	Resource string   `json:"resource"`
	Kind     string   `json:"kind"` // view | action
	Name     string   `json:"name"`
	Command  string   `json:"command"`
	Path     string   `json:"path"`
	Verb     string   `json:"verb"`
	Category string   `json:"category,omitempty"`
	Status   string   `json:"status"`
	Missing  []string `json:"missing,omitempty"` // wajib tetapi tidak ditemukan
	Warn     []string `json:"warn,omitempty"`    // opsional/field tidak ditemukan
	Err      string   `json:"error,omitempty"`
	Note     string   `json:"note,omitempty"`
}

// Summary adalah ringkasan seluruh pemeriksaan.
type Summary struct {
	Resources int `json:"resources"`
	Checks    int `json:"checks"`
	OK        int `json:"ok"`
	Warn      int `json:"warn"`
	Drift     int `json:"drift"`
	Unknown   int `json:"unknown"`
	Error     int `json:"error"`
}

func main() {
	var (
		addr     = flag.String("addr", "192.168.88.1:8728", "alamat API router")
		user     = flag.String("user", "admin", "user API")
		pass     = flag.String("pass", "", "password API")
		resource = flag.String("resource", "", "batasi ke satu ID resource")
		group    = flag.String("group", "", "batasi ke satu grup (hotspot, ip, ...)")
		all      = flag.Bool("all", false, "tampilkan juga cek yang OK")
		asJSON   = flag.Bool("json", false, "keluarkan JSON")
		strict   = flag.Bool("strict", false, "exit 1 bila ada drift/error")
		optional = flag.Bool("optional", false, "ikut periksa atribut/argumen opsional (lebih berisik)")
		timeout  = flag.Duration("timeout", 15*time.Second, "timeout per panggilan inspect")

		discover   = flag.String("discover", "", "mode discovery: path comma, mis. ip,hotspot,active")
		verb       = flag.String("verb", "", "verb untuk discovery, mis. print")
		completion = flag.Bool("completion", false, "tampilkan completion names (discovery)")
	)
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	c, err := routeros.Dial(*addr, *user, *pass)
	if err != nil {
		fatalf("dial: %v", err)
	}
	defer c.Close()

	if *discover != "" {
		if err := runDiscover(ctx, c, *discover, *verb, *completion); err != nil {
			fatalf("discover: %v", err)
		}
		return
	}

	checks := verify(ctx, c, *resource, *group, *timeout, *optional)
	sum := summarize(checks)

	if *asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(struct {
			Summary Summary `json:"summary"`
			Checks  []Check `json:"checks"`
		}{sum, checks})
	} else {
		printText(checks, sum, *all)
	}

	if *strict && (sum.Drift > 0 || sum.Error > 0) {
		os.Exit(1)
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "rosverify: "+format+"\n", args...)
	os.Exit(2)
}

// splitCommand memisahkan command API menjadi path dan verb pada slash
// terakhir: "/ip/hotspot/active/print" -> "/ip/hotspot/active", "print".
func splitCommand(cmd string) (path, verb string, ok bool) {
	if cmd == "" || cmd == "/" || strings.HasSuffix(cmd, "/") {
		return "", "", false
	}
	i := strings.LastIndexByte(cmd, '/')
	switch {
	case i == 0: // command root, mis. /ping
		return "", cmd[1:], true
	case i < 0: // verb tanpa path
		return "", cmd, true
	}
	return cmd[:i], cmd[i+1:], true
}

func verify(ctx context.Context, c *routeros.Client, resource, group string, timeout time.Duration, optional bool) []Check {
	var out []Check
	for _, r := range ros.All() {
		if resource != "" && r.ID != resource {
			continue
		}
		if group != "" && string(r.Group) != group {
			continue
		}

		for _, v := range r.Views {
			out = append(out, checkView(ctx, c, r, v, timeout, optional))
		}
		for _, a := range r.Actions {
			out = append(out, checkAction(ctx, c, r, a, timeout, optional))
		}
	}
	return out
}

func inspect(ctx context.Context, c *routeros.Client, path, verb string, timeout time.Duration) (*schema.CommandSchema, error) {
	ictx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.Inspect(ictx, path, verb)
}

func checkView(ctx context.Context, c *routeros.Client, r *ros.Resource, v ros.View, timeout time.Duration, optional bool) Check {
	chk := Check{Resource: r.ID, Kind: "view", Name: v.Name, Command: v.Command}
	path, verb, ok := splitCommand(v.Command)
	if !ok {
		chk.Status = statusError
		chk.Err = "command tidak bisa dipecah menjadi path/verb"
		return chk
	}
	chk.Path, chk.Verb = path, verb

	sch, err := inspect(ctx, c, path, verb, timeout)
	if err != nil {
		chk.Status = statusError
		chk.Err = err.Error()
		return chk
	}
	chk.Category = string(sch.Category)
	if sch.Category == schema.CategoryUnknown {
		chk.Status = statusUnknown
		chk.Note = unknownNote(v.Mode, path)
		return chk
	}

	var required, warns []string
	if v.Mode == ros.ModeNative {
		for _, a := range v.Args {
			if a.Required || optional {
				required = append(required, a.Name)
			}
		}
	} else {
		for _, p := range v.Proplist {
			if p != ".id" {
				warns = append(warns, p)
			}
		}
	}
	for _, f := range v.Flags {
		warns = append(warns, "flag:"+f)
	}

	found := 0
	for _, n := range required {
		if sch.HasAttribute(n) {
			found++
		} else {
			chk.Missing = append(chk.Missing, n)
		}
	}
	for _, n := range warns {
		name := strings.TrimPrefix(n, "flag:")
		if sch.HasAttribute(name) {
			found++
		} else {
			chk.Warn = append(chk.Warn, n)
		}
	}

	expected := len(required) + len(warns)
	switch {
	case found == 0 && expected > 0:
		// Discovery kosong/incomplete: jangan klaim drift. Ini sering terjadi
		// pada menu singleton (system/resource, system/health, ip/dns) yang
		// tidak mengembalikan field lewat completion.
		chk.Missing, chk.Warn = nil, nil
		chk.Status = statusUnknown
		chk.Note = "atribut/field tidak ditemukan; discovery /console/inspect kosong atau tidak lengkap"
	case len(chk.Missing) > 0:
		chk.Status = statusDrift
	case len(chk.Warn) > 0:
		chk.Status = statusWarn
	default:
		chk.Status = statusOK
	}
	return chk
}

func checkAction(ctx context.Context, c *routeros.Client, r *ros.Resource, a ros.Action, timeout time.Duration, optional bool) Check {
	chk := Check{Resource: r.ID, Kind: "action", Name: a.Name, Command: a.Command}
	path, verb, ok := splitCommand(a.Command)
	if !ok {
		chk.Status = statusError
		chk.Err = "command tidak bisa dipecah menjadi path/verb"
		return chk
	}
	chk.Path, chk.Verb = path, verb

	sch, err := inspect(ctx, c, path, verb, timeout)
	if err != nil {
		chk.Status = statusError
		chk.Err = err.Error()
		return chk
	}
	chk.Category = string(sch.Category)
	if sch.Category == schema.CategoryUnknown {
		chk.Status = statusUnknown
		chk.Note = unknownNote("", path)
		return chk
	}

	found := 0
	for _, n := range a.Required {
		if sch.HasAttribute(n) {
			found++
		} else {
			chk.Missing = append(chk.Missing, n)
		}
	}
	if optional {
		for _, n := range a.Optional {
			if sch.HasAttribute(n) {
				found++
			} else {
				chk.Warn = append(chk.Warn, n)
			}
		}
	}

	expected := len(a.Required)
	if optional {
		expected += len(a.Optional)
	}
	switch {
	case found == 0 && expected > 0:
		chk.Missing, chk.Warn = nil, nil
		chk.Status = statusUnknown
		chk.Note = "atribut tidak ditemukan; discovery /console/inspect kosong atau tidak lengkap"
	case len(chk.Missing) > 0:
		chk.Status = statusDrift
	case len(chk.Warn) > 0:
		chk.Status = statusWarn
	default:
		chk.Status = statusOK
	}
	return chk
}

func unknownNote(m ros.Mode, path string) string {
	if m == ros.ModeListen {
		return "listen adalah verb API-only; /console/inspect adalah skema console"
	}
	if path == "" {
		return "command root; /console/inspect tidak menyediakan skema untuknya"
	}
	return "inspect tidak bisa resolve path/verb ini (mungkin v6, menu dinamis, atau tidak ada di perangkat)"
}

func summarize(checks []Check) Summary {
	var s Summary
	s.Checks = len(checks)
	res := map[string]bool{}
	for _, c := range checks {
		res[c.Resource] = true
		switch c.Status {
		case statusOK:
			s.OK++
		case statusWarn:
			s.Warn++
		case statusDrift:
			s.Drift++
		case statusUnknown:
			s.Unknown++
		case statusError:
			s.Error++
		}
	}
	s.Resources = len(res)
	return s
}

func printText(checks []Check, sum Summary, all bool) {
	for _, c := range checks {
		if !all && c.Status == statusOK {
			continue
		}
		line := fmt.Sprintf("%-7s %s/%s/%s  %s", strings.ToUpper(c.Status), c.Resource, c.Kind, c.Name, c.Command)
		var detail []string
		if c.Category != "" {
			detail = append(detail, "category="+c.Category)
		}
		if len(c.Missing) > 0 {
			detail = append(detail, "wajib-tidak-ada="+strings.Join(c.Missing, ","))
		}
		if len(c.Warn) > 0 {
			detail = append(detail, "tidak-ada="+strings.Join(c.Warn, ","))
		}
		if c.Note != "" {
			detail = append(detail, c.Note)
		}
		if c.Err != "" {
			detail = append(detail, "err="+c.Err)
		}
		if len(detail) > 0 {
			line += "  [" + strings.Join(detail, "; ") + "]"
		}
		fmt.Println(line)
	}
	fmt.Printf("\nringkasan: %d resource, %d cek — ok=%d warn=%d drift=%d unknown=%d error=%d\n",
		sum.Resources, sum.Checks, sum.OK, sum.Warn, sum.Drift, sum.Unknown, sum.Error)
}

func runDiscover(ctx context.Context, c *routeros.Client, path, verb string, completion bool) error {
	path = normalizePath(path)
	if path == "" {
		return fmt.Errorf("path kosong")
	}

	ictx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	children, err := c.InspectNodes(ictx, transport.InspectChild, path)
	if err != nil {
		return fmt.Errorf("child %s: %w", path, err)
	}
	fmt.Printf("=== child %s (%d node)\n", path, len(children))
	for _, n := range children {
		fmt.Printf("  %-28s type=%-4s node-type=%-4s\n", n.Name, n.Type, n.NodeType)
	}
	if verb == "" {
		return nil
	}

	cmdPath := path + "," + verb
	children, err = c.InspectNodes(ictx, transport.InspectChild, cmdPath)
	if err != nil {
		return fmt.Errorf("child %s: %w", cmdPath, err)
	}
	fmt.Printf("\n=== child %s (%d node)\n", cmdPath, len(children))
	for _, n := range children {
		fmt.Printf("  %-28s type=%-4s node-type=%-4s\n", n.Name, n.Type, n.NodeType)
	}

	if completion {
		fieldArg := ""
		if verb == "print" {
			fieldArg = ",proplist"
		} else if verb == "get" {
			fieldArg = ",value-name"
		}
		rows, err := c.InspectNodes(ictx, transport.InspectCompletion, cmdPath+fieldArg)
		if err != nil {
			return fmt.Errorf("completion %s: %w", cmdPath+fieldArg, err)
		}
		names := schema.ExtractCompletionNames(rows)
		sort.Strings(names)
		fmt.Printf("\n=== completion %s (%d nama)\n", cmdPath+fieldArg, len(names))
		for _, n := range names {
			fmt.Printf("  %s\n", n)
		}
	}
	return nil
}

func normalizePath(s string) string {
	s = strings.ReplaceAll(s, "/", ",")
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return strings.Join(out, ",")
}
