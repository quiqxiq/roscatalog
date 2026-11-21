package roscatalog

import "time"

func init() {
	register(
		Table("ip-address", GroupIP, "IP address", "/ip/address").
			Filter("interface", "address", "disabled", "dynamic").
			List("address", "network", "interface", "disabled", "dynamic", "invalid", "comment").
			Events().
			CRUD(P("address", "interface"), P("network", "comment")),

		Table("ip-route", GroupIP, "Routing table (static/dinamis)", "/ip/route").
			Describe("Baseline v7 memakai routing-table. RouterOS 6 memakai routing-mark; "+
				"definisikan resource sendiri bila field itu dibutuhkan. "+
				"Rute default aktif: ?dst-address=0.0.0.0/0 lalu cek active=true.").
			Filter("dst-address", "gateway", "distance", "active", "dynamic", "disabled").
			List("dst-address", "gateway", "distance", "scope", "target-scope", "routing-table",
				"active", "dynamic", "disabled", "comment").
			Count().
			Events().
			CRUD(P("dst-address", "gateway"), P("distance", "scope", "target-scope",
				"routing-table", "check-gateway", "comment")),

		Table("ip-pool", GroupIP, "IP pool", "/ip/pool").
			Filter("name").
			List("name", "ranges", "next-pool", "comment").
			CRUD(P("name", "ranges"), P("next-pool", "comment")),

		Table("ip-pool-used", GroupIP, "IP pool terpakai", "/ip/pool/used").
			Filter("pool", "address", "owner").
			List("pool", "address", "owner", "info").
			MarkHeavy().
			Hint("selalu filter ?pool=<nama>").
			Count(),

		Table("ip-arp", GroupIP, "ARP", "/ip/arp").
			Describe("RouterOS 7 menambah properti status (reachable/stale/...) pada print.").
			Filter("address", "mac-address", "interface", "dynamic"). // dynamic=true utk entri otomatis
			List("address", "mac-address", "interface", "dynamic", "disabled", "published", "comment").
			MarkHeavy().
			Count().
			CRUD(P("address", "mac-address", "interface"), P("comment")),

		Table("ip-neighbor", GroupIP, "Neighbor discovery (MNDP/CDP/LLDP)", "/ip/neighbor").
			Filter("interface", "mac-address", "identity").
			List("interface", "address", "mac-address", "identity", "platform", "version", "board", "age").
			Events(),

		Single("ip-dns", GroupIP, "DNS (pengaturan)", "/ip/dns").
			List("servers", "dynamic-servers", "allow-remote-requests", "cache-size", "cache-used").
			Set("servers", "allow-remote-requests", "cache-size", "cache-max-ttl"),

		Table("ip-dns-cache", GroupIP, "DNS cache", "/ip/dns/cache").
			Filter("name", "type").
			List("name", "type", "data", "ttl").
			MarkHeavy().
			Count().
			Do("flush", SelNone, nil, nil, false, "kosongkan seluruh cache DNS"),

		Table("ip-dns-static", GroupIP, "DNS static", "/ip/dns/static").
			Filter("name", "address", "disabled").
			List("name", "address", "ttl", "disabled", "comment").
			CRUD(P("name", "address"), P("ttl", "comment")),

		// -------------------------------------------------------------------
		// DHCP
		// -------------------------------------------------------------------
		Table("dhcp-server", GroupIP, "DHCP server", "/ip/dhcp-server").
			List("name", "interface", "lease-time", "address-pool", "disabled", "comment").
			CRUD(P("name", "interface"), P("lease-time", "address-pool", "authoritative",
				"use-radius", "comment")),

		Table("dhcp-server-network", GroupIP, "DHCP server network", "/ip/dhcp-server/network").
			List("address", "gateway", "dns-server", "domain", "comment").
			CRUD(P("address"), P("gateway", "dns-server", "domain", "comment")),

		Table("dhcp-lease", GroupIP, "DHCP lease", "/ip/dhcp-server/lease").
			Check().
			Describe("active-address dan active-mac-address diisi untuk lease yang sedang terpakai; "+
				"address/mac-address adalah konfigurasi lease statis.").
			Filter("address", "mac-address", "server", "status", "host-name", "comment", "disabled").
			List("address", "mac-address", "client-id", "host-name", "status", "server",
				"active-address", "active-mac-address", "expires-after", "last-seen",
				"dynamic", "disabled", "block-access", "comment").
			MarkHeavy().
			Count().
			Events().
			Hint("lease baru/berubah/hilang; perubahan expires-after dan last-seen tidak dijamin muncul lewat listen").
			Live(5*time.Second, false, "address", "mac-address", "host-name", "status",
				"expires-after", "last-seen").
			CRUD(P("address", "mac-address"), P("server", "client-id", "comment", "block-access", "rate-limit")).
			Do("make-static", SelNumbers, nil, nil, false, "ubah lease dinamis jadi statis").
			Do("check-status", SelNumbers, nil, nil, false, "bebaskan lease berstatus conflict"),

		Table("dhcp-client", GroupIP, "DHCP client", "/ip/dhcp-client").
			Describe("release/renew memakai =numbers=.").
			List("interface", "status", "address", "gateway", "dhcp-server", "expires-after",
				"add-default-route", "use-peer-dns", "disabled", "comment").
			Events().
			CRUD(P("interface"), P("add-default-route", "use-peer-dns", "use-peer-ntp",
				"default-route-distance", "comment")).
			Do("release", SelNumbers, nil, nil, false, "lepas binding lalu restart client").
			Do("renew", SelNumbers, nil, nil, false, "perbarui lease"),

		Table("ipsec-policy", GroupIP, "IPsec policy", "/ip/ipsec/policy").
			Filter("peer", "dst-address", "action", "protocol").
			List("id", "dst-address", "protocol", "action", "peer", "disabled", "template"),

		Table("ipsec-active-peers", GroupIP, "IPsec active peers", "/ip/ipsec/active-peers").
			List("id", "identity", "address", "state", "dpd-state", "phase1-lifetime",
				"phase2-lifetime"),

		Table("ipsec-installed-sa", GroupIP, "IPsec installed SA", "/ip/ipsec/installed-sa").
			List("id", "name", "security-profile", "phase1-lifetime", "phase2-lifetime",
				"dont-send-deltas", "lifetime-strict-bytes-kb"),

		Single("discovery-settings", GroupIP, "neighbor discovery settings", "/ip/neighbor/discovery-settings").
			Describe("CHR 7.24.6: CHR mengirim kolom kosong; field dari MikroDash.").
			List("name", "interface", "address", "mac-address", "identity", "platform"),

		Table("routing-rule", GroupIP, "Routing rule", "/routing/rule").
			List("id", "table", "disabled"),

		Table("ipv6-route", GroupIP, "IPv6 route", "/ipv6/route").
			Filter("dst-address", "gateway", "distance", "routing-table", "scope", "target-scope").
			List("id", "active", "connect", "distance", "dst-address", "dynamic", "gateway",
				"immediate-gw", "inactive", "routing-table", "scope", "target-scope").
			CRUD(P("dst-address", "gateway"), P("distance", "routing-table", "scope", "target-scope",
				"immediate-gw", "comment")),

		Single("ip-cloud", GroupIP, "IP cloud", "/ip/cloud").
			List("ddns-enabled", "ddns-update-interval", "update-time"),
	)
}
