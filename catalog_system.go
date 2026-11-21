package roscatalog

func init() {
	register(
		Table("disk", GroupSystem, "Storage disk", "/disk").
			Filter("slot", "fs").
			List("slot", "fs", "mounted", "free", "size"),

		Single("app", GroupSystem, "App runtime", "/app").
			Describe("Statis per build router; format field bisa berubah antar versi.").
			List("id", "name", "disabled", "running", "status", "ui-url", "interface"),

		Table("app-settings", GroupSystem, "App settings", "/app/settings").
			Describe("CHR-103 fields: auto-update, registry-mirrors, show-in-webfig; CHR lain memberi 0 field.").
			List("disk", "lan-bridge", "assumed-lan-bridge", "router-ip", "assumed-router-ip",
				"auto-update", "registry-mirrors", "show-in-webfig"),

		Table("user-group", GroupSystem, "User groups", "/user/group").
			Filter("name", "group", "disabled").
			List("id", "name", "policy", "skin", "comment").
			Add(P("name", "policy", "skin", "comment"), nil).
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
	)
}
