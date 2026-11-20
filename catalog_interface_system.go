package roscatalog

import "time"

func init() {
	// -----------------------------------------------------------------------
	// Interface
	// -----------------------------------------------------------------------
	register(
		Table("interface", GroupInterface, "Interface (semua jenis)", "/interface").
			Describe("Counter hanya ada pada view live (print stats). Interface dinamis (PPPoE) kadang tidak "+
				"mengirim .dead pada listen; lengkapi dengan list berkala.").
			Filter("name", "type", "running", "disabled").
			List("name", "type", "mtu", "actual-mtu", "running", "disabled", "comment",
				"last-link-up-time", "last-link-down-time", "link-downs").
			Live(time.Second, true, "name", "running", "rx-byte", "tx-byte", "rx-packet", "tx-packet",
				"rx-drop", "tx-drop", "tx-queue-drop", "rx-error", "tx-error").
			Hint("counter kumulatif; hitung bps di klien dari selisih, atau pakai monitor-traffic").
			Events().
			Set("name", "mtu", "comment").
			Toggle(),

		Tool("interface-monitor-traffic", GroupInterface, "Monitor traffic (bps/pps per interface)",
			"/interface/monitor-traffic",
			[]Arg{
				{Name: "interface", Required: true, Note: "nama, daftar dipisah koma, atau aggregate"},
				{Name: "once", Note: "kirim satu kali lalu selesai"},
				{Name: "duration", Note: "batasi lama streaming, mis. 10s"},
			},
			P("name", "rx-bits-per-second", "tx-bits-per-second", "rx-packets-per-second",
				"tx-packets-per-second", "rx-drops-per-second", "tx-drops-per-second",
				"rx-errors-per-second", "tx-errors-per-second", "tx-queue-drops-per-second"),
			true).
			Check().
			KeyBy("name").
			Describe("Native streaming: satu !re per interface per detik; =.proplist= didukung."),

		Table("interface-ethernet", GroupInterface, "Ethernet", "/interface/ethernet").
			Describe("Monitor memberi status link aktual dan data modul SFP (rx/tx power, suhu, tegangan).").
			Filter("name", "running"). // running juga ada di /interface
			List("name", "default-name", "mac-address", "mtu", "l2mtu", "speed",
				"auto-negotiation", "running", "disabled", "comment").
			Monitor("numbers", "name", "status", "rate", "full-duplex", "auto-negotiation",
				"tx-flow-control", "rx-flow-control", "sfp-module-present", "sfp-type",
				"sfp-vendor-name", "sfp-wavelength", "sfp-temperature", "sfp-supply-voltage",
				"sfp-tx-bias-current", "sfp-tx-power", "sfp-rx-power").KeyBy("name").
			Set("mtu", "speed", "auto-negotiation", "arp", "advertise", "tx-flow-control",
				"rx-flow-control", "comment").
			Toggle().
			Do("reset-counters", SelNumbers, nil, nil, false, "nol-kan counter interface"),

		Table("interface-vlan", GroupInterface, "VLAN", "/interface/vlan").
			Filter("name", "interface", "vlan-id"). // vlan-id angka cocok persis
			List("name", "vlan-id", "interface", "mtu", "running", "disabled", "comment").
			CRUD(P("name", "vlan-id", "interface"), P("mtu", "comment")),

		Table("interface-bridge", GroupInterface, "Bridge", "/interface/bridge").
			Filter("name").
			List("name", "protocol-mode", "vlan-filtering", "mtu", "running", "disabled", "comment").
			CRUD(P("name"), P("protocol-mode", "vlan-filtering", "comment")),

		Table("interface-bridge-port", GroupInterface, "Bridge port", "/interface/bridge/port").
			Filter("bridge", "interface"). // interface = nama port
			List("interface", "bridge", "pvid", "disabled", "comment").
			CRUD(P("interface", "bridge"), P("pvid", "comment")),

		Table("interface-bridge-host", GroupInterface, "Bridge host (tabel MAC)", "/interface/bridge/host").
			Filter("mac-address", "on-interface", "bridge", "local").
			List("mac-address", "on-interface", "bridge", "vid", "age", "local").
			MarkHeavy().
			Hint("cari MAC tertentu: ?mac-address=AA:BB:..."),

		Table("interface-list", GroupInterface, "Interface list", "/interface/list").
			List("name", "comment").
			CRUD(P("name"), P("comment")),

		Table("interface-list-member", GroupInterface, "Interface list member", "/interface/list/member").
			Filter("list", "interface").
			List("list", "interface", "disabled", "comment").
			CRUD(P("list", "interface"), P("comment")),
	)

	// -----------------------------------------------------------------------
	// System
	// -----------------------------------------------------------------------
	register(
		Single("system-resource", GroupSystem, "Resource (CPU, RAM, uptime)", "/system/resource").
			Describe("Satu baris tanpa .id. Pada view live, router mengirim baris yang sama berulang.").
			List("uptime", "version", "build-time", "architecture-name", "board-name", "platform",
				"cpu", "cpu-count", "cpu-frequency", "cpu-load", "free-memory", "total-memory",
				"free-hdd-space", "total-hdd-space").
			Live(time.Second, false, "uptime", "cpu-load", "free-memory", "total-memory",
				"free-hdd-space", "total-hdd-space"),

		Table("system-resource-cpu", GroupSystem, "CPU per core", "/system/resource/cpu").
			List("cpu", "load", "irq", "disk").
			Live(time.Second, false, "cpu", "load", "irq", "disk"),

		Single("system-health", GroupSystem, "Health (suhu, tegangan, fan)", "/system/health").
			Describe("Baseline v7 mengirim baris name/value/type. RouterOS 6 mengirim "+
				"properti langsung; definisikan resource sendiri bila field itu dibutuhkan. "+
				"Sensor yang tersedia bergantung model.").
			List("name", "value", "type").
			Live(5*time.Second, false, "name", "value", "type"),

		Single("system-routerboard", GroupSystem, "RouterBOARD", "/system/routerboard").
			List("routerboard", "model", "serial-number", "firmware-type", "factory-firmware",
				"current-firmware", "upgrade-firmware").
			Do("upgrade", SelNone, nil, nil, false, "tulis firmware baru; berlaku setelah reboot"),

		Single("system-identity", GroupSystem, "Identity", "/system/identity").
			List("name").
			Set("name"),

		Single("system-clock", GroupSystem, "Clock", "/system/clock").
			List("time", "date", "time-zone-name", "time-zone-autodetect").
			Set("time", "date", "time-zone-name", "time-zone-autodetect"),

		Table("system-package", GroupSystem, "Paket terpasang", "/system/package").
			List("name", "version", "build-time", "disabled", "scheduled").
			Toggle(),

		Single("system-package-update", GroupSystem, "Update paket", "/system/package/update").
			Describe("check-for-updates mengirim beberapa !re status lalu !done (self-terminating).").
			List("channel", "installed-version", "latest-version", "status").
			Set("channel").
			add(View{
				Name: "check", Mode: ModeNative, Command: "/system/package/update/check-for-updates",
				Proplist: P("channel", "installed-version", "latest-version", "status"), SelfTerminating: true,
			}),

		Table("system-script", GroupSystem, "Script", "/system/script").
			Describe("run memakai =number= (nama atau .id) dan tidak mengembalikan output script.").
			Filter("name", "owner").
			List("name", "owner", "policy", "last-started", "run-count", "comment").
			Named("detail", "name", "owner", "policy", "dont-require-permissions", "last-started",
				"run-count", "source", "comment").
			Hint("source bisa panjang; ambil hanya saat perlu").
			Count().
			CRUD(P("name", "source"), P("owner", "policy", "dont-require-permissions", "comment")).
			Do("run", SelNumber, nil, nil, false, "jalankan script"),

		Table("system-scheduler", GroupSystem, "Scheduler", "/system/scheduler").
			Filter("name", "disabled").
			List("name", "start-date", "start-time", "interval", "on-event", "next-run",
				"run-count", "policy", "disabled", "comment").
			CRUD(P("name", "on-event"), P("start-date", "start-time", "interval", "policy", "comment")),

		Single("system-backup", GroupSystem, "Backup", "/system/backup").
			Describe("File backup ada di /file; unduh dengan /tool/fetch (upload) atau FTP/SFTP. "+
				"load menimpa konfigurasi dan me-reboot router.").
			Do("save", SelNone, nil, P("name", "password", "dont-encrypt"), false, "").
			Do("load", SelNone, P("name"), P("password"), true, "menimpa konfigurasi"),

		Single("system-export", GroupSystem, "Export konfigurasi", "/export").
			Describe("Gunakan =file= lalu ambil file-nya. RouterOS 7 menyembunyikan secret secara default "+
				"(show-sensitive untuk menampilkan); RouterOS 6 memakai hide-sensitive.").
			Do("export", SelNone, nil, P("file", "compact", "terse", "verbose", "hide-sensitive",
				"show-sensitive"), false, ""),

		Single("system-power", GroupSystem, "Reboot / shutdown", "/system").
			Do("reboot", SelNone, nil, nil, true, "router putus sampai selesai boot").
			Do("shutdown", SelNone, nil, nil, true, "router mati; perlu akses fisik untuk menyalakan"),

		Table("log", GroupSystem, "Log router", "/log").
			Describe("API tidak punya regex, jadi filter topics/message dilakukan di klien. "+
				"follow dan listen hanya mengirim entri baru; list/tail mengirim isi buffer log dulu.").
			List("time", "topics", "message").
			MarkHeavy().
			Hint("buffer log memory; batasi dengan membaca sekali saat start").
			Tail("time", "topics", "message").
			TailOnly("time", "topics", "message").
			Events(),

		Table("file", GroupSystem, "File", "/file").
			Describe("Jangan minta properti contents lewat API untuk file besar.").
			Filter("name", "type").
			List("name", "type", "size", "creation-time").
			Count().
			Remove("hapus file di router"),

		Table("user", GroupSystem, "User router", "/user").
			Filter("name", "group", "disabled").
			List("name", "group", "address", "last-logged-in", "disabled", "comment").
			CRUD(P("name", "password", "group"), P("address", "comment")),

		Table("user-active", GroupSystem, "Sesi admin aktif", "/user/active").
			List("name", "when", "address", "via", "group", "radius").
			Events(),

		Table("ip-service", GroupSystem, "Service manajemen (ssh, api, winbox)", "/ip/service").
			Describe("Matikan service yang tidak dipakai; batasi address ke jaringan manajemen.").
			Filter("name", "disabled").
			List("name", "port", "address", "disabled").
			Set("port", "address").
			Toggle(),
	)
}
