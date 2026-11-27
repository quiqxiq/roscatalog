package roscatalog

func init() {
	register(
		Table("disk", GroupSystem, "Storage disk", "/disk").
			Filter("slot", "fs").
			List("slot", "fs", "mounted", "free", "size"),

		Table("app", GroupSystem, "App runtime", "/app").
			Describe("Statis per build router; format field bisa berubah antar versi. "+
				"CHR-192/103: field kosong; CHR lama lebih lengkap dari MikroDash.").
			List("name", "disabled", "running", "status", "ui-url", "interface").
			Set("auto-update", "disabled", "environment", "network", "use-https", "yaml").
			Toggle(),

		Single("app-settings", GroupSystem, "App settings", "/app/settings").
			Describe("CHR-103 fields: auto-update, registry-mirrors, show-in-webfig; CHR lain memberi 0 field.").
			List("disk", "lan-bridge", "router-ip",
				"auto-update", "registry-mirrors", "show-in-webfig").
			Set("disk", "lan-bridge", "router-ip",
				"auto-update", "registry-mirrors", "show-in-webfig"),

		Table("user-group", GroupSystem, "User groups", "/user/group").
			Filter("name").
			List("name", "policy", "skin", "comment").
			Add(P("name"), P("policy", "skin", "comment")).
			Set(P("name", "policy", "skin", "comment")...).
			Remove("hapus group"),

		Single("user-settings", GroupSystem, "User settings", "/user/settings").
			List("minimum-categories", "minimum-password-length"),

		Single("device-mode", GroupSystem, "Device mode", "/system/device-mode").
			List("allowed-versions", "attempt-count", "bandwidth-test", "container", "email",
				"fetch", "flagged", "flagging-enabled", "hotspot", "install-any-version", "ipsec",
				"l2tp", "mode", "partitions", "pptp", "proxy", "romon", "routerboard", "scheduler",
				"smb", "socks", "sniffer", "traffic-gen", "zerotier"),

		Single("license", GroupSystem, "License", "/system/license").
			List("features", "nlevel", "software-id", "level", "system-id"),

		Single("routerboard-settings", GroupSystem, "RouterBOARD settings", "/system/routerboard/settings").
			List("board-identifier", "bootloader-protocol", "cpu-mode", "channels"),

		// CHR-192/103 tidak punya menu ini; field dari MikroDash/server.
		Single("system-reset", GroupSystem, "Reset konfigurasi", "/system/reset-configuration").
			DoAt("/system/reset-configuration", "reset-configuration", SelNone, nil, nil, true,
				"reset ke konfigurasi standar, lalu reboot otomatis"),

		Single("system-import", GroupSystem, "Import konfigurasi", "/import").
			DoAt("/import", "import", SelNone, P("file-name"), P("verbose", "dry-run"), true,
				"import config dari file (putus sampai selesai)"),
	)
}
