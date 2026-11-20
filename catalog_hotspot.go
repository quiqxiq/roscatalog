package roscatalog

import "time"

func init() {
	register(
		// -------------------------------------------------------------------
		// Server hotspot
		// -------------------------------------------------------------------
		Table("hotspot-server", GroupHotspot, "Hotspot server", "/ip/hotspot").
			Check().
			Filter("name", "interface", "profile", "disabled").
			List("name", "interface", "address-pool", "profile", "idle-timeout",
				"keepalive-timeout", "login-timeout", "addresses-per-mac", "disabled", "comment").
			CRUD(P("name", "interface"), P("address-pool", "profile", "idle-timeout",
				"keepalive-timeout", "login-timeout", "addresses-per-mac", "comment")),

		Table("hotspot-server-profile", GroupHotspot, "Hotspot server profile", "/ip/hotspot/profile").
			Filter("name").
			List("name", "hotspot-address", "dns-name", "html-directory", "login-by",
				"use-radius", "rate-limit", "http-cookie-lifetime", "install-hotspot-queue",
				"mac-auth-mode", "comment").
			CRUD(P("name"), P("hotspot-address", "dns-name", "html-directory", "login-by",
				"use-radius", "rate-limit", "http-cookie-lifetime", "install-hotspot-queue",
				"mac-auth-mode", "radius-mac-format", "smtp-server", "comment")),

		// -------------------------------------------------------------------
		// User profile (dipakai Mikhmon untuk paket voucher)
		// -------------------------------------------------------------------
		Table("hotspot-user-profile", GroupHotspot, "Hotspot user profile (paket voucher)", "/ip/hotspot/user/profile").
			Filter("name").
			List("name", "rate-limit", "shared-users", "session-timeout", "idle-timeout",
				"keepalive-timeout", "status-autorefresh", "address-pool", "address-list",
				"parent-queue", "mac-cookie-timeout", "add-mac-cookie", "transparent-proxy", "comment").
			Hint("tanpa on-login/on-logout supaya ringan").
			Named("detail", "name", "rate-limit", "shared-users", "session-timeout", "idle-timeout",
				"keepalive-timeout", "status-autorefresh", "address-pool", "address-list",
				"parent-queue", "mac-cookie-timeout", "add-mac-cookie", "transparent-proxy",
				"on-login", "on-logout", "insert-queue-before", "comment").
			Hint("on-login memuat script expire/harga Mikhmon; ambil hanya saat perlu").
			CRUD(P("name"), P("rate-limit", "shared-users", "session-timeout", "idle-timeout",
				"keepalive-timeout", "status-autorefresh", "address-pool", "address-list",
				"parent-queue", "mac-cookie-timeout", "add-mac-cookie", "transparent-proxy",
				"on-login", "on-logout", "comment")),

		// -------------------------------------------------------------------
		// User (voucher)
		// -------------------------------------------------------------------
		Table("hotspot-user", GroupHotspot, "Hotspot user (voucher)", "/ip/hotspot/user").
			Check().
			Describe("Query ?comment=... hanya cocok persis; pencarian prefix (mis. batch vc-) dilakukan di sisi klien.").
			Filter("name", "profile", "server", "mac-address", "address", "comment", "disabled").
			List("name", "profile", "server", "mac-address", "address", "limit-uptime",
				"limit-bytes-total", "uptime", "bytes-in", "bytes-out", "comment", "disabled").
			MarkHeavy().
			Hint("ribuan voucher: pakai Count atau filter profile/comment").
			Named("voucher", "name", "password", "profile", "limit-uptime", "limit-bytes-total", "comment").
			Hint("berisi password; untuk cetak voucher saja").
			Count().
			Events().
			Hint("sinkronisasi add/remove/ubah user; perubahan counter (uptime, bytes) tidak dijamin muncul lewat listen").
			CRUD(P("name"), P("password", "profile", "server", "address", "mac-address",
				"limit-uptime", "limit-bytes-in", "limit-bytes-out", "limit-bytes-total", "email", "routes", "comment")).
			Do("reset-counters", SelNumbers, nil, nil, false, "nol-kan uptime dan bytes user"),

		// -------------------------------------------------------------------
		// Active, host, binding, cookie
		// -------------------------------------------------------------------
		Table("hotspot-active", GroupHotspot, "Hotspot active (sesi login)", "/ip/hotspot/active").
			Check().
			Describe("Listen cocok untuk login/logout; uptime dan bytes berubah terus dan perubahan counter tidak dijamin muncul di listen, jadi ambil lewat view live.").
			Filter("user", "server", "address", "mac-address", "login-by").
			List("user", "address", "mac-address", "server", "login-by", "uptime",
				"session-time-left", "idle-time", "bytes-in", "bytes-out", "radius").
			Live(2*time.Second, true, "user", "address", "uptime", "session-time-left",
				"idle-time", "bytes-in", "bytes-out").
			Hint("snapshot penuh tiap siklus; naikkan interval bila sesi ribuan").
			Count().
			Events().
			Remove("kick: memutus sesi login user"),

		Table("hotspot-host", GroupHotspot, "Hotspot host", "/ip/hotspot/host").
			Check().
			Filter("mac-address", "address", "server").
			List("mac-address", "address", "to-address", "server", "uptime", "idle-time",
				"found-by", "authorized", "bypassed", "comment").
			MarkHeavy().
			Remove("").
			Do("make-binding", SelNumbers, nil, P("comment", "type"), false,
				"type: regular | bypassed | blocked"),

		Table("hotspot-ip-binding", GroupHotspot, "Hotspot IP binding (bypass/blokir)", "/ip/hotspot/ip-binding").
			Check().
			Filter("mac-address", "address", "type", "server").
			List("mac-address", "address", "to-address", "server", "type", "disabled", "comment").
			CRUD(nil, P("mac-address", "address", "to-address", "server", "type", "comment")),

		Table("hotspot-cookie", GroupHotspot, "Hotspot cookie (mac-cookie)", "/ip/hotspot/cookie").
			Check().
			Filter("user", "mac-address").
			List("user", "mac-address", "domain", "expires-in").
			Remove("hapus cookie agar user harus login ulang"),

		// -------------------------------------------------------------------
		// Walled garden
		// -------------------------------------------------------------------
		Table("hotspot-walled-garden", GroupHotspot, "Walled garden (host/path)", "/ip/hotspot/walled-garden").
			List("action", "server", "src-address", "dst-address", "dst-host", "dst-port",
				"path", "method", "hits", "disabled", "comment").
			CRUD(nil, P("action", "server", "src-address", "dst-address", "dst-host", "dst-port",
				"path", "method", "comment")).
			Do("reset-counters", SelNumbers, nil, nil, false, "").
			Do("reset-counters-all", SelNone, nil, nil, false, ""),

		Table("hotspot-walled-garden-ip", GroupHotspot, "Walled garden IP", "/ip/hotspot/walled-garden/ip").
			List("action", "server", "src-address", "dst-address", "dst-host", "protocol",
				"dst-port", "disabled", "comment").
			CRUD(nil, P("action", "server", "src-address", "dst-address", "dst-host", "protocol",
				"dst-port", "comment")),

		// -------------------------------------------------------------------
		// Laporan penjualan Mikhmon (disimpan sebagai /system/script)
		// -------------------------------------------------------------------
		Table("mikhmon-sales", GroupHotspot, "Laporan penjualan Mikhmon (/system/script)", "/system/script").
			Describe("Konvensi Mikhmon, bukan perilaku RouterOS: name berformat "+
				"tanggal-|-jam-|-user-|-harga-|-ip-|-mac-|-validity-|-profile-|-comment, "+
				"owner berisi bulan+tahun (mis. nov2026), source berisi tanggal, comment=mikhmon. "+
				"Verifikasi format terhadap versi Mikhmon yang dipakai.").
			Filter("owner", "comment", "name").
			Named("report", "name", "owner", "source", "comment").
			MarkHeavy().
			Hint("selalu filter ?owner=<bulantahun>").
			Count().
			Add(P("name", "source"), P("owner", "comment")).
			Remove("hapus laporan per .id; hapus per bulan = print ?owner= lalu remove satu per satu"),
	)
}
