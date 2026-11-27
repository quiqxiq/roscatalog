package roscatalog

func init() {
	register(
		Table("ppp-secret", GroupPPP, "PPP secret (akun PPPoE)", "/ppp/secret").
			Check().
			Describe("password sengaja tidak masuk proplist list; hanya dikirim saat add/set.").
			Filter("name", "service", "profile", "caller-id", "disabled").
			List("name", "service", "profile", "caller-id", "local-address", "remote-address",
				"limit-bytes-in", "limit-bytes-out", "last-logged-out", "last-caller-id",
				"last-disconnect-reason", "disabled", "comment").
			MarkHeavy().
			Hint("ribuan pelanggan: filter name/profile atau Count").
			Count().
			Events().
			CRUD(P("name", "password"), P("service", "profile", "caller-id", "local-address",
				"remote-address", "routes", "limit-bytes-in", "limit-bytes-out", "comment")),

		Table("ppp-active", GroupPPP, "PPP active (sesi PPPoE aktif)", "/ppp/active").
			Check().
			Describe("Listen cocok untuk deteksi pelanggan online/offline; uptime dihitung klien dari waktu event.").
			Filter("name", "service", "caller-id", "address"). // caller-id = MAC pada PPPoE
			List("name", "service", "caller-id", "address", "uptime", "encoding",
				"session-id", "limit-bytes-in", "limit-bytes-out", "radius").
			Count().
			Events().
			Remove("kick: memutus sesi PPPoE (client biasanya reconnect otomatis)"),

		Table("ppp-profile", GroupPPP, "PPP profile (paket)", "/ppp/profile").
			Check().
			Filter("name").
			List("name", "local-address", "remote-address", "rate-limit", "parent-queue",
				"address-list", "dns-server", "only-one", "use-compression", "use-encryption",
				"session-timeout", "idle-timeout", "change-tcp-mss", "comment").
			Add(P("name"), P("local-address", "remote-address", "rate-limit", "parent-queue",
				"address-list", "dns-server", "only-one", "session-timeout", "idle-timeout",
				"change-tcp-mss", "on-up", "on-down", "comment")).
			Set("name", "local-address", "remote-address", "rate-limit", "parent-queue",
				"address-list", "dns-server", "only-one", "session-timeout", "idle-timeout",
				"change-tcp-mss", "on-up", "on-down", "comment").
			Remove(""),

		Single("ppp-aaa", GroupPPP, "PPP AAA (RADIUS untuk PPP)", "/ppp/aaa").
			Check().
			List("use-radius", "accounting", "interim-update", "use-circuit-id-in-nas-port-id").
			Set("use-radius", "accounting", "interim-update", "use-circuit-id-in-nas-port-id",
				"enable-ipv6-accounting"),

		Table("pppoe-server", GroupPPP, "PPPoE server (service)", "/interface/pppoe-server/server").
			List("service-name", "interface", "default-profile", "authentication",
				"keepalive-timeout", "max-sessions", "one-session-per-host", "disabled", "comment").
			CRUD(P("service-name", "interface"), P("default-profile", "authentication",
				"keepalive-timeout", "max-sessions", "one-session-per-host", "comment")),

		Table("pppoe-client", GroupInterface, "PPPoE client (uplink)", "/interface/pppoe-client").
			Describe("Field monitor mengikuti pppoe-server/monitor yang terverifikasi, ditambah service-name dan ac-name; validasi di router target.").
			List("name", "interface", "user", "service-name", "running", "disabled", "comment").
			Monitor("numbers", "status", "uptime", "encoding", "service-name", "ac-name",
				"local-address", "remote-address").
			CRUD(P("interface", "user", "password"), P("name", "service-name", "add-default-route",
				"use-peer-dns", "comment")),

		// -------------------------------------------------------------------
		// RADIUS
		// -------------------------------------------------------------------
		Table("radius", GroupPPP, "RADIUS client", "/radius").
			Check().
			Describe("secret tidak masuk proplist list. Monitor menunjukkan counter request/accept/reject per server.").
			List("service", "address", "authentication-port", "accounting-port", "timeout",
				"realm", "src-address", "accounting-backup", "status", "disabled", "comment").
			Monitor("numbers", "pending", "requests", "accepts", "rejects", "resends",
				"timeouts", "bad-replies", "last-request-rtt").
			CRUD(P("service", "address", "secret"), P("authentication-port", "accounting-port",
				"timeout", "realm", "src-address", "accounting-backup", "comment")).
			Do("reset-counters", SelNumbers, nil, nil, false, "nol-kan counter monitor"),

		Single("radius-incoming", GroupPPP, "RADIUS incoming (CoA / Disconnect)", "/radius/incoming").
			Describe("Aktifkan accept=yes agar FreeRADIUS bisa mengirim Disconnect-Request/CoA (port 3799).").
			List("accept", "port").
			Set("accept", "port"),
	)
}
