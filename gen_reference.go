//go:build ignore

// Jalankan: go run gen_reference.go > REFERENCE.md
package main

import (
	"fmt"
	"os"
	"strings"

	ros "example.com/roscatalog"
)

var groupTitle = map[ros.Group]string{
	ros.GroupHotspot:   "Hotspot dan Mikhmon",
	ros.GroupPPP:       "PPP, PPPoE dan RADIUS",
	ros.GroupIP:        "IP (address, route, pool, ARP, DNS, DHCP)",
	ros.GroupFirewall:  "Firewall",
	ros.GroupQueue:     "Queue",
	ros.GroupInterface: "Interface",
	ros.GroupSystem:    "Sistem dan operasional",
	ros.GroupTools:     "Troubleshooting (tools)",
}

func code(s string) string { return "`" + s + "`" }

func codeList(items []string) string {
	if len(items) == 0 {
		return "-"
	}
	out := make([]string, len(items))
	for i, s := range items {
		out[i] = code(s)
	}
	return strings.Join(out, ", ")
}

func deliveries(r *ros.Resource, v ros.View) string {
	ds := r.Deliveries(v.Name, v.Version)
	if len(ds) < 2 && v.Mode == ros.ModeListen {
		ds = r.Deliveries(v.Name, ros.V7) // list bergantung versi
	}
	out := make([]string, len(ds))
	for i, d := range ds {
		out[i] = string(d)
	}
	return strings.Join(out, ", ")
}

func cell(s string) string { return strings.ReplaceAll(s, "|", "\\|") }

func fields(v ros.View) string {
	if v.Mode == ros.ModeListen {
		return "diambil dari view list lalu disaring di klien"
	}
	if len(v.Proplist) == 0 {
		return "-"
	}
	if v.ServerProplist {
		return "server: " + codeList(v.ProplistWords())
	}
	return "klien: " + codeList(v.Proplist)
}

func extras(v ros.View) string {
	var parts []string
	for _, f := range v.Flags {
		parts = append(parts, code("="+f+"="))
	}
	if v.Mode == ros.ModeInterval {
		parts = append(parts, code("=interval="+ros.FormatTime(v.Interval)))
	}
	for _, a := range v.Args {
		s := code(a.Name)
		if a.Required {
			s += "*"
		}
		parts = append(parts, s)
	}
	if len(parts) == 0 {
		return "-"
	}
	return strings.Join(parts, " ")
}

func main() {
	w := os.Stdout
	all := ros.All()

	views, actions, verified := 0, 0, 0
	modes := map[ros.Mode]int{}
	for _, r := range all {
		views += len(r.Views)
		actions += len(r.Actions)
		if r.Verified {
			verified++
		}
		for _, v := range r.Views {
			modes[v.Mode]++
		}
	}

	fmt.Fprintf(w, "# Katalog command RouterOS untuk Go\n\n")
	fmt.Fprintf(w, "Dibuat otomatis dari paket `roscatalog` (`go run gen_reference.go`). "+
		"%d resource, %d view baca, %d aksi tulis. Dokumen ini dan kode tidak bisa berbeda karena berasal dari data yang sama.\n\n",
		len(all), views, actions)

	fmt.Fprintf(w, "## Cara membaca\n\n")
	fmt.Fprintf(w, "| Mode | Pemicu di API | Isi tiap pesan | Snapshot awal | Berhenti |\n|---|---|---|---|---|\n")
	fmt.Fprintf(w, "| `oneshot` | `/path/print` | Semua baris, lalu `!done` | Ya | Otomatis |\n")
	fmt.Fprintf(w, "| `native` | command `monitor`, `monitor-traffic`, `ping`, `torch`, ... | Satu `!re` per update (nilai terkini) | Tidak relevan | `once`, `duration`, `count`, atau `/cancel` |\n")
	fmt.Fprintf(w, "| `interval` | `/path/print` + `=interval=` (+ `=stats=`) | Snapshot penuh berulang | Ya, tiap siklus | `/cancel` |\n")
	fmt.Fprintf(w, "| `follow` | `/path/print` + `=follow=` atau `=follow-only=` | Entri baru (`follow` kirim entri lama dulu) | `follow`: ya; `follow-only`: tidak | `/cancel` |\n")
	fmt.Fprintf(w, "| `listen` | `/path/listen` | Hanya item yang berubah/muncul; hilang = `=.dead=true` | Tidak | `/cancel` |\n\n")
	fmt.Fprintf(w, "Jumlah view per mode: oneshot %d, native %d, interval %d, follow %d, listen %d.\n\n",
		modes[ros.ModeOneShot], modes[ros.ModeNative], modes[ros.ModeInterval], modes[ros.ModeFollow], modes[ros.ModeListen])

	fmt.Fprintf(w, "### Stream atau poll\n\n")
	fmt.Fprintf(w, "Setiap view bisa dikirim lewat `Plan`: **stream** (router yang mengirim terus) atau **poll** (Go mengulang one-shot tiap N ms). `Resource.Build(view, versi, plan, params)` menghasilkan `Request` untuk `rosrun`.\n\n")
	fmt.Fprintf(w, "| Mode view | Stream | Poll |\n|---|---|---|\n")
	fmt.Fprintf(w, "| `oneshot` | tidak ada | `print` tiap Every |\n")
	fmt.Fprintf(w, "| `interval` | `print =interval=<Every>` (+`=stats=`) | `print` (+`=stats=`) tiap Every, tanpa interval |\n")
	fmt.Fprintf(w, "| `native` | command tanpa `once` (opsional `=interval=`) | command + `=once=` tiap Every (hanya bila command punya `once`) |\n")
	fmt.Fprintf(w, "| `follow` | `print =follow=` / `=follow-only=` | `print` tiap Every + diff (`follow-only`: hasil pertama hanya baseline) |\n")
	fmt.Fprintf(w, "| `listen` | `/listen` (opsional Seed: print awal) | `print` view list tiap Every + diff |\n\n")
	fmt.Fprintf(w, "Pilih **stream** untuk data yang sering berubah dengan latensi rendah. Pilih **poll** bila feed banyak, router lemah, jeda perlu diatur per klien/tenant, atau agar tidak menahan stream terbuka di router. Batas proyek: interval stream %s..%s, jeda poll %s..%s.\n\n", ros.MinPrintInterval, ros.MaxPrintInterval, ros.MinPollEvery, ros.MaxPollEvery)

	fmt.Fprintf(w, "Keterangan kolom dan tanda:\n\n")
	fmt.Fprintf(w, "- **Pengiriman** = pilihan `Plan` yang valid untuk view itu; yang pertama adalah bawaan.\n")
	fmt.Fprintf(w, "- **Proplist `server:`** dikirim sebagai `=.proplist=` (sudah termasuk `.id` untuk tabel). `rosrun` tetap menyaring ulang kolom di sisi klien dan, bila router menolak `.proplist` saat start, mengulang sekali tanpa `.proplist` (`Request.Fallback`).\n")
	fmt.Fprintf(w, "- Tanda `*` pada argumen berarti wajib. Nilai kosong (mis. `once`) dikirim sebagai `=once=`.\n")
	fmt.Fprintf(w, "- Aksi: kolom *Selector* menunjukkan argumen pemilih item (`=.id=`, `=numbers=`, `=number=`). Aksi **destruktif** mengubah atau menghapus data/layanan.\n")
	fmt.Fprintf(w, "- ✅ = nama properti dicocokkan dengan CLI reference RouterOS 7.26. Resource tanpa tanda berasal dari pengetahuan umum RouterOS 6.4x/7.x: validasi di router target dengan `NeverSeen()`.\n")
	fmt.Fprintf(w, "- ⚠ = tabel berpotensi besar: pakai filter, `count-only`, atau proplist sempit.\n")
	fmt.Fprintf(w, "- `v6`/`v7` = view khusus mayor tersebut. Properti kosong (mis. `comment` yang belum diisi) memang tidak dikirim router.\n\n")

	fmt.Fprintf(w, "Terverifikasi: %d dari %d resource (hotspot, PPP, RADIUS, DHCP lease, monitor-traffic).\n\n", verified, len(all))

	fmt.Fprintf(w, "### Aturan penanganan di Go (goros/v4)\n\n")
	fmt.Fprintf(w, "- Library mengelola `.tag` sendiri; jangan kirim `Params.Tag`. Batalkan stream lewat `ListenReply.Cancel()` (`/cancel =tag=`); `!trap` kategori 2 sesudah cancel bukan error.\n")
	fmt.Fprintf(w, "- Jangan pakai varian `*Context` library dengan ctx yang bisa dibatalkan: ia memanggil `c.r.Cancel()` yang mematikan pembaca bersama semua command. `rosrun` memakai `RunArgs`/`ListenArgs`.\n")
	fmt.Fprintf(w, "- Channel stream library tidak berpenyangga (`Client.Queue` bawaan 0): konsumen lambat memblokir semua balasan di koneksi itu. `rosrun` menguras channel ke antrean sendiri (batas `MaxBuffered`).\n")
	fmt.Fprintf(w, "- `listen` tidak punya snapshot awal: `Plan.WithSeed()` memulai listen dulu, baru `print`, lalu menggabungkan per `.id`. Perubahan parsial digabung dengan state sebelumnya.\n")
	fmt.Fprintf(w, "- `interval` tidak menandai batas siklus. `rosrun` memakai `.section` bila ada, kunci yang berulang, atau jeda tanpa baris (Every/5, 20ms..300ms). Tabel kosong tidak mengirim apa pun; `Runner.EmptySnapshots` mengirim snapshot kosong setelah 3x Every tanpa baris (inferensi).\n")
	fmt.Fprintf(w, "- API tidak punya `where` dan regex `~`: gunakan query word `?prop=value` (persis) dan filter prefix/regex di klien.\n")
	fmt.Fprintf(w, "- Deteksi kemampuan menu lewat `InspectChild(...)` (`/console/inspect`).\n\n")

	for _, g := range ros.GroupOrder {
		rs := ros.ByGroup(g)
		if len(rs) == 0 {
			continue
		}
		fmt.Fprintf(w, "---\n\n## %s\n\n", groupTitle[g])
		for _, r := range rs {
			mark := ""
			if r.Verified {
				mark = " ✅"
			}
			fmt.Fprintf(w, "### %s: %s%s\n\n", r.ID, r.Title, mark)
			fmt.Fprintf(w, "Path: %s", code(r.Path))
			if r.Singleton {
				fmt.Fprintf(w, " (tanpa baris)")
			}
			fmt.Fprintf(w, "\n\n")
			if len(r.Filters) > 0 {
				fmt.Fprintf(w, "Filter server `?prop=nilai` (persis): %s\n\n", codeList(r.Filters))
			}
			if r.Notes != "" {
				fmt.Fprintf(w, "%s\n\n", r.Notes)
			}
			if len(r.Views) > 0 {
				fmt.Fprintf(w, "| View | Mode | Pengiriman | Command | Flag / argumen | Proplist |\n|---|---|---|---|---|---|\n")
				for _, v := range r.Views {
					name := v.Name
					if v.Version != ros.AnyVersion {
						name += " (" + string(v.Version) + ")"
					}
					if v.Heavy {
						name += " ⚠"
					}
					note := ""
					if v.Note != "" {
						note = "<br>" + cell(v.Note)
					}
					fmt.Fprintf(w, "| %s | %s | %s | %s | %s | %s%s |\n",
						name, v.Mode, deliveries(r, v), code(v.Command), cell(extras(v)), cell(fields(v)), note)
				}
				fmt.Fprintf(w, "\n")
			}
			if len(r.Actions) > 0 {
				fmt.Fprintf(w, "| Aksi | Command | Selector | Wajib | Opsional |\n|---|---|---|---|---|\n")
				for _, a := range r.Actions {
					name := a.Name
					if a.Destructive {
						name += " (destruktif)"
					}
					sel := "-"
					if a.Selector != ros.SelNone {
						sel = code("=" + string(a.Selector) + "=")
					}
					opt := codeList(a.Optional)
					if a.Name == "set" && len(a.Optional) > 8 {
						opt = codeList(a.Optional[:8]) + fmt.Sprintf(", ... (+%d)", len(a.Optional)-8)
					}
					if a.Note != "" {
						opt += "<br>" + cell(a.Note)
					}
					fmt.Fprintf(w, "| %s | %s | %s | %s | %s |\n",
						name, code(a.Command), sel, codeList(a.Required), opt)
				}
				fmt.Fprintf(w, "\n")
			}
		}
	}
}
