package roscatalog

import "time"

func init() {
	register(
		Tool("tool-ping", GroupTools, "Ping", "/ping",
			[]Arg{
				{Name: "address", Required: true},
				{Name: "count", Note: "tanpa count, ping berjalan sampai /cancel"},
				{Name: "interval", Note: "jeda antar probe, mis. 1s"},
				{Name: "size"},
				{Name: "ttl"},
				{Name: "interface"},
				{Name: "src-address"},
				{Name: "do-not-fragment"},
				{Name: "routing-table", Note: "nama routing table untuk ping lewat tabel tertentu"},
			},
			P("seq", "host", "size", "ttl", "time", "status", "sent", "received",
				"packet-loss", "min-rtt", "avg-rtt", "max-rtt"),
			true).
			Describe("Satu !re per probe; field sent/received/packet-loss/min/avg/max bersifat kumulatif. "+
				"status berisi timeout bila probe gagal."),

		Tool("tool-traceroute", GroupTools, "Traceroute", "/tool/traceroute",
			[]Arg{
				{Name: "address", Required: true},
				{Name: "count"},
				{Name: "max-hops"},
				{Name: "timeout"},
				{Name: "size"},
				{Name: "protocol"},
				{Name: "port"},
				{Name: "src-address"},
				{Name: "interface"},
				{Name: "use-dns"},
			},
			P("address", "loss", "sent", "last", "avg", "best", "worst", "std-dev", "status"),
			true).
			Describe("Router mengirim ulang seluruh daftar hop tiap update."),

		Tool("tool-torch", GroupTools, "Torch (flow per sumber/tujuan)", "/tool/torch",
			[]Arg{
				{Name: "interface", Required: true},
				{Name: "src-address"},
				{Name: "dst-address"},
				{Name: "mac-protocol"},
				{Name: "protocol"},
				{Name: "port"},
				{Name: "ip-protocol"},
				{Name: "vlan-id"},
				{Name: "duration", Note: "tanpa duration, berjalan sampai /cancel"},
			},
			P("src-address", "dst-address", "mac-protocol", "ip-protocol", "src-port",
				"dst-port", "vlan-id", "tx", "rx", "tx-packets", "rx-packets"),
			true).
			Describe("Tiap update berisi seluruh flow aktif; urutkan di klien menurut tx/rx. "+
				"Berat di CPU router pada interface yang padat."),

		Tool("tool-bandwidth-test", GroupTools, "Bandwidth test", "/tool/bandwidth-test",
			[]Arg{
				{Name: "address", Required: true},
				{Name: "user"},
				{Name: "password"},
				{Name: "protocol", Note: "udp | tcp"},
				{Name: "direction", Note: "transmit | receive | both"},
				{Name: "local-udp-tx-size"},
				{Name: "remote-udp-tx-size"},
				{Name: "local-tx-speed"},
				{Name: "remote-tx-speed"},
				{Name: "tcp-connection-count"},
				{Name: "random-data"},
				{Name: "duration", Note: "tanpa duration, berjalan sampai /cancel"},
			},
			P("status", "duration", "tx-current", "tx-10-second-average", "tx-total-average",
				"rx-current", "rx-10-second-average", "rx-total-average", "lost-packets",
				"direction", "connection-count", "local-cpu-load", "remote-cpu-load"),
			true).
			Describe("Menghabiskan bandwidth dan CPU di kedua router; jangan dijalankan tanpa izin pada link produksi."),

		Tool("tool-profile", GroupTools, "Profile CPU per proses", "/tool/profile",
			[]Arg{
				{Name: "cpu", Note: "all atau nomor core"},
				{Name: "duration", Note: "tanpa duration, berjalan sampai /cancel"},
			},
			P("name", "cpu", "usage"),
			true),

		Table("tool-sniffer", GroupTools, "Sniffer (paket sniffing)", "/tool/sniffer").
			Describe("CHR-103 fields (file-limit ... streaming-server).").
			Filter("interface", "file-name", "streaming-server").
			List("id", "file-limit", "file-name", "filter-cpu", "filter-direction",
				"filter-dst-ip-address", "filter-dst-ipv6-address", "filter-dst-mac-address",
				"filter-dst-port", "filter-interface", "filter-ip-address", "filter-ip-protocol",
				"filter-ipv6-address", "filter-mac-address", "filter-mac-protocol",
				"filter-operator-between-entries", "filter-port", "filter-size",
				"filter-src-ip-address", "filter-src-ipv6-address", "filter-src-mac-address",
				"filter-src-port", "filter-stream", "filter-vlan", "max-packet-size",
				"memory-limit", "memory-scroll", "only-headers", "quick-rows", "quick-show-frame",
				"running", "streaming-enabled", "streaming-server").
			Set("file-limit", "file-name", "filter-interface", "filter-dst-ip-address",
				"filter-dst-ipv6-address", "filter-dst-mac-address", "filter-dst-port",
				"filter-ip-address", "filter-ip-protocol", "filter-ipv6-address",
				"filter-mac-address", "filter-mac-protocol", "filter-operator-between-entries",
				"filter-port", "filter-size", "filter-src-ip-address", "filter-src-ipv6-address",
				"filter-src-mac-address", "filter-src-port", "filter-stream", "filter-vlan",
				"memory-limit", "memory-scroll", "only-headers", "quick-rows", "quick-show-frame",
				"streaming-enabled", "streaming-server").
			Do("start", SelNone, nil, nil, false, "jalankan sniffer").
			Do("stop", SelNone, nil, nil, false, "hentikan sniffer").
			Do("save", SelNone, nil, P("file-name"), false, "simpan sniffer sebagai file"),

		Tool("console", GroupTools, "Console (execute)", "/execute", nil, nil, false).
			Describe("CHR-192/103 tidak punya /execute; execute console script."),

		Tool("menu", GroupTools, "Menu API", "/menu", nil, nil, false).
			Describe("CHR-192/103 tidak punya /menu."),

		Table("tool-ping-status", GroupTools, "Ping status (baris terakhir)", "/tool/ping").
			List("time", "response-time", "status", "min-rtt", "max-rtt").
			Describe("CHR-192/103 tidak punya /tool/ping; lihat tool-ping (/ping)."),

		Table("tool-netwatch", GroupTools, "Netwatch (monitor host up/down)", "/tool/netwatch").
			Describe("Perubahan status muncul sebagai perubahan properti status/since; listen layak dicoba, "+
				"fallback ke view live bila event tidak datang.").
			Filter("host", "status", "type", "disabled").
			List("host", "type", "interval", "timeout", "status", "since", "disabled", "comment").
			Live(5*time.Second, false, "host", "status", "since").
			Events().
			CRUD(P("host"), P("type", "interval", "timeout", "up-script", "down-script", "comment")),

		Tool("tool-fetch", GroupTools, "Fetch (unduh/unggah file, HTTP)", "/tool/fetch",
			[]Arg{
				{Name: "url", Required: true},
				{Name: "dst-path"},
				{Name: "mode"},
				{Name: "check-certificate"},
				{Name: "keep-result"},
				{Name: "upload"},
				{Name: "src-path"},
				{Name: "user"},
				{Name: "password"},
				{Name: "http-method"},
				{Name: "http-data"},
			},
			P("status", "downloaded", "total", "duration"),
			true).
			Describe("Mengirim beberapa !re progres lalu !done. Pakai upload=yes untuk mengirim backup/export ke server."),
	)
}
