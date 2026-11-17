package roscatalog

import "time"

// ruleProps adalah properti umum rule firewall.
var ruleProps = []string{
	"chain", "action", "protocol", "src-address", "dst-address", "src-port", "dst-port",
	"in-interface", "out-interface", "in-interface-list", "out-interface-list",
	"src-address-list", "dst-address-list", "connection-state",
}

func concat(parts ...[]string) []string {
	var out []string
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

// fwRule membuat resource rule firewall (filter/nat/mangle/raw).
func fwRule(id, title, path string, extra []string) *Resource {
	list := concat(ruleProps, extra, P("log", "disabled", "comment"))
	opt := concat(ruleProps[2:], extra, P("log", "log-prefix", "place-before", "comment"))
	return Table(id, GroupFirewall, title, path).
		Describe("Counter bytes/packets hanya muncul pada view live (print stats). "+
			"place-before pada add berisi *ID rule yang akan ditaruh di depannya.").
		Filter("chain", "action", "disabled", "dynamic").
		List(list...).
		Live(time.Second, true, "chain", "action", "bytes", "packets", "disabled", "comment").
		Hint("pantau hit rule; cocokkan baris dengan .id").
		Count().
		CRUD(P("chain", "action"), opt).
		Move().
		ResetCounters()
}

func init() {
	register(
		fwRule("fw-filter", "Firewall filter", "/ip/firewall/filter", nil),
		fwRule("fw-nat", "Firewall NAT", "/ip/firewall/nat", P("to-addresses", "to-ports")),
		fwRule("fw-mangle", "Firewall mangle", "/ip/firewall/mangle",
			P("connection-mark", "packet-mark", "new-connection-mark", "new-packet-mark",
				"new-routing-mark", "passthrough")),
		fwRule("fw-raw", "Firewall raw", "/ip/firewall/raw", nil),

		Table("fw-address-list", GroupFirewall, "Firewall address-list (isolir, blacklist)", "/ip/firewall/address-list").
			Describe("timeout pada print adalah sisa waktu; entri berakhir otomatis dan hilang sendiri.").
			Filter("list", "address", "dynamic", "disabled").
			List("list", "address", "timeout", "creation-time", "dynamic", "disabled", "comment").
			MarkHeavy().
			Hint("selalu filter ?list=<nama>").
			Count().
			Events().
			CRUD(P("list", "address"), P("timeout", "comment")),

		Table("fw-connection", GroupFirewall, "Connection tracking", "/ip/firewall/connection").
			Describe("API tidak punya regex/prefix, jadi satu IP tidak bisa difilter di sisi router "+
				"(src-address berisi ip:port). Pakai Count untuk dashboard; untuk detail ambil proplist "+
								"sempit lalu filter di klien.").
			Filter("protocol", "tcp-state"). // src-address/dst-address berisi ip:port: hanya cocok persis
			List("protocol", "src-address", "dst-address", "reply-src-address", "reply-dst-address",
				"tcp-state", "timeout", "orig-bytes", "repl-bytes").
			MarkHeavy().
			Hint("tabel bisa puluhan ribu baris; jangan dipakai tanpa Count/proplist").
			Count().
			Remove("putus koneksi per .id"),
	)

	// -----------------------------------------------------------------------
	// Queue
	// -----------------------------------------------------------------------
	register(
		Table("queue-simple", GroupQueue, "Simple queue", "/queue/simple").
			Describe("rate dan packet-rate berformat upload/download (tx/rx) dalam bit/s.").
			Filter("name", "target", "parent", "disabled", "dynamic").
			List("name", "target", "max-limit", "limit-at", "parent", "priority", "queue",
				"disabled", "comment").
			Live(time.Second, true, "name", "target", "rate", "packet-rate", "bytes", "packets",
				"dropped", "queued-bytes", "queued-packets").
			Hint("pantau throughput per queue").
			Count().
			CRUD(P("name", "target"), P("max-limit", "limit-at", "burst-limit", "burst-threshold",
				"burst-time", "parent", "priority", "queue", "place-before", "comment")).
			Move().
			ResetCounters(),

		Table("queue-tree", GroupQueue, "Queue tree", "/queue/tree").
			Filter("name", "parent", "packet-mark", "disabled").
			List("name", "parent", "packet-mark", "limit-at", "max-limit", "priority", "queue",
				"disabled", "comment").
			Live(time.Second, true, "name", "parent", "rate", "packet-rate", "bytes", "packets",
				"dropped", "queued-packets").
			CRUD(P("name", "parent"), P("packet-mark", "limit-at", "max-limit", "priority",
				"queue", "comment")).
			ResetCounters(),
	)
}
