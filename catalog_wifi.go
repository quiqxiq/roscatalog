package roscatalog

func init() {
	// Wi-Fi v7 (interface/wifi): data CHR 7.24 memberi menu, CHR 7.24.6 dengan
	// paket wifi memuat proplist — MikroDash (target v7.25+) definisikan isi
	// field ini agar catatannya lengkap untuk kolektor time-series.
	register(
		Table("wifi-ap", GroupInterface, "Wi-Fi AP (SRS)", "/interface/wifi").
			Describe("Statis per build router; set diaktifkan lewat enable/disable pada interface").
			Filter("name", "ssid", "mode").
			List("name", "default-name", "disabled", "running", "master-interface",
				"radio-mac", "mac-address", "cap", "ssid").
			Set("ssid", "mode", "country", "security", "channel", "datapath", "comment"),

		Table("wifi-channel", GroupInterface, "Wi-Fi channel", "/interface/wifi/channel").
			List("name", "band", "frequency", "width", "secondary-frequency",
				"skip-dfs-channels", "disabled", "comment"),

		Table("wifi-configuration", GroupInterface, "Wi-Fi configuration", "/interface/wifi/configuration").
			List("name", "ssid", "mode", "country", "hide-ssid", "security", "channel",
				"datapath", "manager", "disabled", "comment").
			Set("ssid", "mode", "country", "security", "channel", "datapath", "manager", "comment"),

		Table("wifi-security", GroupInterface, "Wi-Fi security", "/interface/wifi/security").
			List("name", "authentication-types", "wps", "ft", "ft-over-ds",
				"connect-priority", "disabled", "comment").
			Set("authentication-types", "wps", "ft", "ft-over-ds", "connect-priority", "comment"),

		Table("wifi-datapath", GroupInterface, "Wi-Fi datapath", "/interface/wifi/datapath").
			List("name", "bridge", "vlan-id", "client-isolation", "local-forwarding"),

		Table("wifi-provisioning", GroupInterface, "Wi-Fi provisioning", "/interface/wifi/provisioning").
			List("supported-bands", "action", "master-configuration", "slave-configurations"),

		Table("wifi-radio", GroupInterface, "Wi-Fi radio", "/interface/wifi/radio").
			List("name", "radio-mac", "interface", "cap", "disabled"),

		Table("wifi-registration-table", GroupInterface, "Wi-Fi clients", "/interface/wifi/registration-table").
			List("mac-address", "interface", "ssid", "signal", "uptime"),

		// Legacy Wi-Fi (v7 downgrade): CHR 7.24 tidak punya menu ini; field dari
		// bawaan wi-Fi package MikroDash.
		Table("wifi-legacy", GroupInterface, "Legacy Wi-Fi", "/interface/wireless").
			List("name", "default-name", "disabled", "running", "ssid", "mode",
				"band", "frequency", "channel-width", "security", "comment"),
		Table("wifi-legacy-security-profiles", GroupInterface, "Legacy Wi-Fi security", "/interface/wireless/security-profiles").
			List("name", "mode", "authentication-types", "default"),
		Table("wifi-legacy-registration-table", GroupInterface, "Legacy Wi-Fi clients", "/interface/wireless/registration-table").
			List("mac-address", "interface", "ssid", "signal-strength", "uptime"),

		// Caps-Man v7 (/caps-man): CHR 7.24 error di beberapa submenu (rates,
		// registration-table) — field dari MikroDash.
		Table("wifi-capsman", GroupInterface, "CapsMan", "/caps-man").
			List("name", "interfaces", "provisioning", "priority", "ds-properties",
				"manager-configuration", "registration-table", "security", "datapath",
				"channel", "interface", "remote-cap").
			Set("priority", "ds-properties", "manager-configuration", "security", "datapath"),

		Table("wifi-capsman-channel", GroupInterface, "CapsMan channel", "/caps-man/channel").
			List("name", "band", "frequency", "width", "extension-channel",
				"secondary-frequency", "scan-list", "center-freq-strs", "skip-dfs-channels",
				"disabled", "comment").
			Set("frequency", "width", "comment"),

		Table("wifi-capsman-configuration", GroupInterface, "CapsMan configuration", "/caps-man/configuration").
			List("name", "ssid", "mode", "country", "hide-ssid", "security", "channel",
				"datapath", "comment").
			Set("ssid", "mode", "country", "security", "channel", "datapath", "comment"),

		Table("wifi-capsman-datapath", GroupInterface, "CapsMan datapath", "/caps-man/datapath").
			List("name", "bridge", "vlan-id", "vlan-mode", "client-to-client-forwarding",
				"auto-startup", "stp-mode", "l2mtu"),

		Table("wifi-capsman-interface", GroupInterface, "CapsMan management interface", "/caps-man/interface").
			List("name", "mac-address", "master-interface", "radio-mac", "configuration",
				"mtu", "ap-clients"),

		Table("wifi-capsman-provisioning", GroupInterface, "CapsMan provisioning", "/caps-man/provisioning").
			List("name", "hw-supported-modes", "action", "master-configuration",
				"slave-configurations"),

		Table("wifi-capsman-registration-table", GroupInterface, "CapsMan clients", "/caps-man/registration-table").
			List("mac-address", "interface", "ssid", "signal", "uptime"),

		Table("wifi-capsman-remote-cap", GroupInterface, "CapsMan remote cap", "/caps-man/remote-cap").
			List("identity", "address", "board-name", "state"),

		Table("wifi-capsman-security", GroupInterface, "CapsMan security", "/caps-man/security").
			List("name", "authentication-types", "encryption", "group-encryption",
				"eap-methods", "wpa-version", "encryption-protocols", "use-dot1x"),

		// CHR-192/103 tidak punya menu-menu ini; field dari MikroDash/capsman.go.
		Table("capsman-manager", GroupInterface, "CapsMan manager (v7)", "/caps-man/manager").
			List("name", "interfaces", "priority", "ds-properties", "manager-configuration",
				"registration-table", "security", "datapath", "channel", "interface", "remote-cap"),

		Table("capsman-radio", GroupInterface, "CapsMan radio (v7)", "/caps-man/radio").
			List("name", "radio-mac", "interface", "channel", "usage", "comment"),

		Table("interface-wifi-cap", GroupInterface, "Wi-Fi cap (v7)", "/interface/wifi/cap").
			List("name", "is-ready", "interfaces", "setup-time", "mode",
				"tsf-clock-level", "is-pending", "strength"),

		Table("interface-wifi-capsman", GroupInterface, "Wi-Fi capsman (v7)", "/interface/wifi/capsman").
			List("name", "interfaces", "priority", "ds-properties", "manager-configuration",
				"registration-table", "security", "datapath", "channel", "interface", "remote-cap", "rates"),

		Table("interface-wifi-capsman-remote-cap", GroupInterface, "Wi-Fi capsman remote cap (v7)",
			"/interface/wifi/capsman/remote-cap").
			List("identity", "address", "board-name", "state"),

		Tool("interface-wifi-frequency-scan", GroupInterface, "Wi-Fi frequency scan (v7)",
			"/interface/wifi/frequency-scan", nil, nil, false).
			Describe("Native streaming; CHR-192/103 tidak punya menu ini."),
	)
}
