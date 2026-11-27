// Tool bantu: dump katalog roscatalog ke JSON supaya generator Python bisa
// menggabungkan proplist/aksi yang sudah dikurasi dengan data inspect + live.
package main

import (
	"encoding/json"
	"os"

	ros "github.com/quiqxiq/roscatalog"
)

type viewOut struct {
	Name     string   `json:"name"`
	Mode     string   `json:"mode"`
	Command  string   `json:"command"`
	Proplist []string `json:"proplist,omitempty"`
}

type actionOut struct {
	Name     string   `json:"name"`
	Command  string   `json:"command"`
	Selector string   `json:"selector"`
	Required []string `json:"required,omitempty"`
	Optional []string `json:"optional,omitempty"`
}

type resourceOut struct {
	ID        string      `json:"id"`
	Path      string      `json:"path"`
	Singleton bool        `json:"singleton"`
	Views     []viewOut   `json:"views"`
	Actions   []actionOut `json:"actions"`
}

func main() {
	out := make([]resourceOut, 0)
	for _, r := range ros.All() {
		item := resourceOut{
			ID: r.ID, Path: r.Path, Singleton: r.Singleton,
			Views:   make([]viewOut, 0, len(r.Views)),
			Actions: make([]actionOut, 0, len(r.Actions)),
		}
		for _, v := range r.Views {
			item.Views = append(item.Views, viewOut{
				Name: v.Name, Mode: string(v.Mode), Command: v.Command,
				Proplist: v.Proplist,
			})
		}
		for _, a := range r.Actions {
			item.Actions = append(item.Actions, actionOut{
				Name: a.Name, Command: a.Command, Selector: string(a.Selector),
				Required: a.Required, Optional: a.Optional,
			})
		}
		out = append(out, item)
	}
	blob, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("mikrotik_inspect/catalog_dump.json", append(blob, '\n'), 0o644); err != nil {
		panic(err)
	}
}
