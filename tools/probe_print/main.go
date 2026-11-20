// Probe: ambil kolom hasil /path/print dari router untuk semua menu yang
// tercatat di mikrotik_inspect/_index.json.
//
// Dipakai sekali untuk mengisi mikrotik_inspect/print_fields.json.
// Bukan bagian dari library; hanya tooling maintainer.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/quiqxiq/goros/v4"
	"github.com/quiqxiq/goros/v4/schema"
	"github.com/quiqxiq/goros/v4/transport"

	ros "example.com/roscatalog"
	"example.com/roscatalog/rosrun"
)

type indexFile struct {
	Source string `json:"source"`
	Files  []struct {
		File  string `json:"file"`
		Path  string `json:"path"`
		Scope string `json:"scope"`
	} `json:"files"`
}

type pathResult struct {
	Fields        []string `json:"fields,omitempty"`     // dari live print
	Completion    []string `json:"completion,omitempty"` // dari inspect completion
	Rows          int      `json:"rows"`
	Error         string   `json:"error,omitempty"`
	CompletionErr string   `json:"completion_error,omitempty"`
}

var validName = regexp.MustCompile(`^\.?[A-Za-z][A-Za-z0-9._-]*$`)

func cleanNames(names []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, n := range names {
		if !validName.MatchString(n) || seen[n] {
			continue
		}
		seen[n] = true
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

type output struct {
	Router     string                 `json:"router"`
	CapturedAt string                 `json:"captured_at"`
	Source     string                 `json:"source"`
	Paths      map[string]*pathResult `json:"paths"`
}

func main() {
	addr := flag.String("addr", "192.168.122.189:8728", "alamat API router")
	user := flag.String("user", "admin", "user API")
	pass := flag.String("pass", "r00t", "password API")
	index := flag.String("index", "mikrotik_inspect/_index.json", "file indeks inspect")
	out := flag.String("out", "mikrotik_inspect/print_fields.json", "file hasil")
	perCmd := flag.Duration("timeout", 4*time.Second, "batas waktu per print")
	delay := flag.Duration("delay", 15*time.Millisecond, "jeda antar command")
	flag.Parse()

	raw, err := os.ReadFile(*index)
	if err != nil {
		panic(err)
	}
	var idx indexFile
	if err := json.Unmarshal(raw, &idx); err != nil {
		panic(err)
	}

	paths := []string{}
	for _, f := range idx.Files {
		if f.Scope == "menu" && f.Path != "" {
			paths = append(paths, f.Path)
		}
	}
	sort.Strings(paths)

	c, err := routeros.Dial(*addr, *user, *pass)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	run := rosrun.New(c)

	res := output{
		Router:     *addr,
		CapturedAt: time.Now().UTC().Format(time.RFC3339),
		Source:     "live /<path>/print (tanpa .proplist=)",
		Paths:      make(map[string]*pathResult, len(paths)),
	}

	for i, p := range paths {
		pr := &pathResult{}
		res.Paths[p] = pr
		sentence := []string{p + "/print"}
		req := ros.Request{
			Delivery: ros.DeliveryOnce,
			Sentence: sentence,
		}
		ctx, cancel := context.WithTimeout(context.Background(), *perCmd)
		rows, err := run.Once(ctx, req)
		cancel()
		if err != nil {
			pr.Error = err.Error()
		} else {
			pr.Rows = len(rows)
			seen := map[string]bool{}
			for _, row := range rows {
				for k := range row {
					seen[k] = true
				}
			}
			for k := range seen {
				pr.Fields = append(pr.Fields, k)
			}
			sort.Strings(pr.Fields)
		}

		// Completion proplist: nama field tetap didapat walau tabel kosong.
		ictx, icancel := context.WithTimeout(context.Background(), 10*time.Second)
		pathComma := strings.ReplaceAll(p, "/", ",") + ",print,proplist"
		nodes, cerr := c.InspectNodes(ictx, transport.InspectCompletion, pathComma)
		icancel()
		if cerr != nil {
			pr.CompletionErr = cerr.Error()
		} else {
			pr.Completion = cleanNames(schema.ExtractCompletionNames(nodes))
		}
		if (i+1)%25 == 0 {
			fmt.Fprintf(os.Stderr, "[%d/%d] %s\n", i+1, len(paths), p)
		}
		time.Sleep(*delay)
	}

	blob, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(*out, append(blob, '\n'), 0o644); err != nil {
		panic(err)
	}
	ok, failed, empty, completed := 0, 0, 0, 0
	for _, pr := range res.Paths {
		if len(pr.Completion) > 0 {
			completed++
		}
		switch {
		case pr.Error != "":
			failed++
		case len(pr.Fields) == 0:
			empty++
		default:
			ok++
		}
	}
	fmt.Printf("selesai: %d path, %d sukses, %d kosong, %d error; %d punya completion -> %s\n",
		len(paths), ok, empty, failed, completed, *out)
}
