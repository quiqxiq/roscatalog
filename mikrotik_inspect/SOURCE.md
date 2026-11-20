# Sumber struct RouterOS (baseline v7.9+)

> File ini dihasilkan oleh `tools/gen_source_md.py`. Jangan diedit manual — regen dari sumber di bawah.

## Sumber dan metadata

| Sumber | Keterangan |
|---|---|
| `mikrotik_inspect/*.json` | command, argumen, dan `desc` dari `/console/inspect` |
| `print_fields.json` | kolom `print` live + completion dari router 192.168.122.189:8728 (diambil 2026-11-20T01:28:20Z) |
| `catalog_dump.json` | proplist terverifikasi + required/optional dari katalog Go |

Router sumber memakai RouterOS v7.9+ (`/console/inspect` baru ada di v7.9). Command v7 relatif stabil, jadi struct di bawah dipakai sebagai baseline; perbedaan v6 ditangani terpisah lewat `Version` di katalog.

Arti tag:

- `ros:"nama"` → nama properti RouterOS yang dikirim/dibaca.
- `,key` → kunci baris (default `.id`).
- `,required` → wajib saat `add`/`set` (dari katalog).
- `,omitempty` → field kosong tidak dikirim (struct tulis).

`XxxRow` dipakai untuk proplist `print`; `XxxAdd` / `XxxSet` dipakai untuk argumen tulis.

## Ringkasan

- Entries diproses: 217
- tanpa field: 36
- punya set: 90
- punya add: 68
- dengan field: 181
- File JSON hadir: 99
- File di index tapi tidak hadir: 144
- Error baca file: 0

Catatan: `_index.json` memuat 243 entri file, tetapi hanya sebagian yang ikut ter-attach di direktori ini. Untuk path yang filenya tidak ada tetapi punya hasil `print` dari router, tetap dibuatkan struct `Row` berdasarkan data print (tanpa `Add`/`Set` karena argumennya tidak tersedia).

## Indeks menu

| Path | Field | Add | Set | Sumber |
|---|---:|---:|---:|---|
| `/` | 0 | 0 | 3 | tidak ada kolom terdeteksi |
| `/app` | 0 | 16 | 16 | tidak ada kolom terdeteksi |
| `/app/network` | 0 | 0 | 0 | tidak ada kolom terdeteksi |
| `/app/settings` | 0 | 0 | 9 | tidak ada kolom terdeteksi |
| `/caps-man` | 0 | 0 | 0 | tidak ada kolom terdeteksi |
| `/caps-man/aaa` | 0 | 0 | 5 | tidak ada kolom terdeteksi |
| `/caps-man/access-list` | 0 | 19 | 18 | tidak ada kolom terdeteksi |
| `/caps-man/actual-interface-configuration` | 0 | 0 | 57 | tidak ada kolom terdeteksi |
| `/caps-man/channel` | 0 | 12 | 12 | tidak ada kolom terdeteksi |
| `/caps-man/configuration` | 0 | 61 | 61 | tidak ada kolom terdeteksi |
| `/caps-man/datapath` | 0 | 15 | 15 | tidak ada kolom terdeteksi |
| `/caps-man/interface` | 0 | 69 | 69 | tidak ada kolom terdeteksi |
| `/caps-man/manager` | 0 | 0 | 6 | tidak ada kolom terdeteksi |
| `/caps-man/provisioning` | 0 | 14 | 13 | tidak ada kolom terdeteksi |
| `/caps-man/radio` | 0 | 0 | 0 | tidak ada kolom terdeteksi |
| `/caps-man/rates` | 0 | 9 | 9 | tidak ada kolom terdeteksi |
| `/caps-man/registration-table` | 0 | 0 | 0 | tidak ada kolom terdeteksi |
| `/caps-man/remote-cap` | 0 | 0 | 0 | tidak ada kolom terdeteksi |
| `/caps-man/security` | 0 | 13 | 13 | tidak ada kolom terdeteksi |
| `/certificate` | 0 | 15 | 15 | tidak ada kolom terdeteksi |
| `/certificate/builtin` | 18 | 0 | 0 | live 17, completion 18 |
| `/certificate/crl` | 14 | 0 | 0 | completion 14 |
| `/certificate/scep-server` | 7 | 0 | 0 | completion 7 |
| `/certificate/settings` | 5 | 0 | 0 | live 5 |
| `/console/settings` | 3 | 0 | 0 | live 3 |
| `/disk/settings` | 5 | 0 | 0 | live 5 |
| `/file` | 4 | 0 | 0 | katalog 4 |
| `/interface` | 19 | 0 | 6 | katalog 21 |
| `/interface/6to4` | 15 | 0 | 0 | completion 15 |
| `/interface/bonding` | 27 | 0 | 0 | completion 27 |
| `/interface/bridge` | 58 | 52 | 52 | live 33, completion 58, katalog 7 |
| `/interface/bridge/host` | 6 | 0 | 0 | katalog 6 |
| `/interface/bridge/port` | 5 | 4 | 4 | katalog 5 |
| `/interface/detect-internet` | 5 | 0 | 5 | live 5 |
| `/interface/eoip` | 25 | 0 | 0 | completion 25 |
| `/interface/eoipv6` | 24 | 0 | 0 | completion 24 |
| `/interface/ethernet` | 166 | 0 | 26 | live 19, completion 154, katalog 26 |
| `/interface/gre` | 16 | 0 | 0 | completion 16 |
| `/interface/gre6` | 15 | 0 | 0 | completion 15 |
| `/interface/ipip` | 17 | 0 | 0 | completion 17 |
| `/interface/ipipv6` | 15 | 0 | 0 | completion 15 |
| `/interface/l2tp-client` | 27 | 0 | 0 | completion 27 |
| `/interface/l2tp-ether` | 27 | 0 | 0 | completion 27 |
| `/interface/l2tp-server` | 12 | 0 | 0 | completion 12 |
| `/interface/list` | 7 | 5 | 5 | live 7, completion 7, katalog 2 |
| `/interface/list/member` | 4 | 3 | 3 | katalog 4 |
| `/interface/lte` | 20 | 0 | 0 | completion 20 |
| `/interface/macsec` | 42 | 9 | 9 | completion 42 |
| `/interface/macvlan` | 16 | 13 | 13 | completion 16 |
| `/interface/mesh` | 22 | 20 | 20 | completion 22 |
| `/interface/monitor-traffic` | 10 | 0 | 0 | katalog 10 |
| `/interface/ovpn-client` | 26 | 22 | 22 | completion 26 |
| `/interface/ovpn-server` | 12 | 5 | 5 | completion 12 |
| `/interface/ppp-client` | 28 | 27 | 27 | completion 28 |
| `/interface/ppp-server` | 15 | 14 | 14 | completion 15 |
| `/interface/pppoe-client` | 27 | 20 | 20 | completion 22, katalog 14 |
| `/interface/pppoe-server` | 15 | 6 | 6 | completion 15 |
| `/interface/pppoe-server/server` | 9 | 8 | 8 | katalog 9 |
| `/interface/pptp-client` | 18 | 0 | 0 | completion 18 |
| `/interface/pptp-server` | 12 | 0 | 0 | completion 12 |
| `/interface/sstp-client` | 28 | 0 | 0 | completion 28 |
| `/interface/sstp-server` | 12 | 0 | 0 | completion 12 |
| `/interface/vlan` | 19 | 14 | 14 | completion 19, katalog 7 |
| `/interface/vpls` | 25 | 0 | 0 | completion 25 |
| `/interface/vrrp` | 34 | 0 | 0 | completion 34 |
| `/interface/vxlan` | 32 | 0 | 0 | completion 32 |
| `/interface/wifi` | 141 | 134 | 134 | completion 141 |
| `/interface/wireguard` | 10 | 8 | 8 | completion 10 |
| `/interface/wireless` | 0 | 95 | 95 | tidak ada kolom terdeteksi |
| `/ip` | 0 | 0 | 0 | tidak ada kolom terdeteksi |
| `/ip/address` | 13 | 8 | 8 | live 11, completion 13, katalog 7 |
| `/ip/arp` | 13 | 7 | 7 | live 12, completion 13, katalog 7 |
| `/ip/cloud` | 3 | 0 | 0 | live 3 |
| `/ip/dhcp-client` | 34 | 17 | 17 | live 22, completion 34, katalog 10 |
| `/ip/dhcp-relay` | 12 | 11 | 11 | completion 12 |
| `/ip/dhcp-server` | 32 | 30 | 30 | live 15, completion 32, katalog 6 |
| `/ip/dhcp-server/lease` | 14 | 7 | 7 | katalog 20 |
| `/ip/dhcp-server/network` | 5 | 5 | 5 | katalog 5 |
| `/ip/dns` | 19 | 0 | 17 | live 19, katalog 5 |
| `/ip/dns/cache` | 4 | 0 | 0 | katalog 4 |
| `/ip/dns/static` | 5 | 4 | 4 | katalog 5 |
| `/ip/firewall` | 0 | 0 | 0 | tidak ada kolom terdeteksi |
| `/ip/firewall/address-list` | 7 | 4 | 4 | katalog 7 |
| `/ip/firewall/connection` | 9 | 0 | 0 | katalog 9 |
| `/ip/firewall/filter` | 19 | 18 | 18 | katalog 23 |
| `/ip/firewall/mangle` | 25 | 24 | 24 | katalog 29 |
| `/ip/firewall/nat` | 21 | 20 | 20 | katalog 25 |
| `/ip/firewall/raw` | 19 | 18 | 18 | katalog 23 |
| `/ip/hotspot` | 15 | 10 | 10 | live 14, completion 14, katalog 10 |
| `/ip/hotspot/active` | 11 | 0 | 0 | katalog 18 |
| `/ip/hotspot/cookie` | 4 | 0 | 0 | katalog 4 |
| `/ip/hotspot/host` | 10 | 0 | 0 | katalog 10 |
| `/ip/hotspot/ip-binding` | 7 | 6 | 6 | katalog 7 |
| `/ip/hotspot/profile` | 11 | 13 | 13 | katalog 11 |
| `/ip/hotspot/user` | 12 | 13 | 13 | katalog 12 |
| `/ip/hotspot/user/profile` | 17 | 16 | 16 | katalog 31 |
| `/ip/hotspot/walled-garden` | 11 | 9 | 9 | katalog 11 |
| `/ip/hotspot/walled-garden/ip` | 9 | 8 | 8 | katalog 9 |
| `/ip/ipsec` | 0 | 0 | 0 | tidak ada kolom terdeteksi |
| `/ip/kid-control` | 23 | 0 | 0 | completion 23 |
| `/ip/media` | 9 | 0 | 0 | completion 9 |
| `/ip/nat-pmp` | 1 | 0 | 0 | live 1 |
| `/ip/neighbor` | 21 | 0 | 0 | live 20, completion 21, katalog 8 |
| `/ip/packing` | 6 | 6 | 6 | completion 6 |
| `/ip/pool` | 8 | 5 | 5 | live 6, completion 8, katalog 4 |
| `/ip/pool/used` | 4 | 0 | 0 | katalog 4 |
| `/ip/proxy` | 17 | 0 | 0 | live 17 |
| `/ip/reverse-proxy` | 9 | 0 | 0 | completion 9 |
| `/ip/route` | 31 | 13 | 13 | live 15, completion 30, katalog 20 |
| `/ip/service` | 18 | 0 | 9 | live 15, completion 18, katalog 4 |
| `/ip/settings` | 23 | 0 | 0 | live 23 |
| `/ip/smb` | 5 | 0 | 0 | live 5 |
| `/ip/socks` | 7 | 0 | 0 | live 7 |
| `/ip/socksify` | 11 | 0 | 0 | completion 11 |
| `/ip/ssh` | 7 | 0 | 0 | live 7 |
| `/ip/tftp` | 12 | 0 | 0 | completion 12 |
| `/ip/traffic-flow` | 8 | 0 | 8 | live 8 |
| `/ip/upnp` | 3 | 0 | 0 | live 3 |
| `/ip/vrf` | 6 | 0 | 0 | live 5, completion 6 |
| `/ipv6/address` | 21 | 0 | 0 | live 16, completion 21 |
| `/ipv6/dhcp-client` | 33 | 0 | 0 | completion 33 |
| `/ipv6/dhcp-relay` | 11 | 0 | 0 | completion 11 |
| `/ipv6/dhcp-server` | 23 | 0 | 0 | completion 23 |
| `/ipv6/nd` | 20 | 0 | 0 | live 16, completion 20 |
| `/ipv6/neighbor` | 10 | 0 | 0 | completion 10 |
| `/ipv6/pool` | 11 | 0 | 0 | completion 11 |
| `/ipv6/route` | 30 | 0 | 0 | live 12, completion 30 |
| `/ipv6/settings` | 19 | 0 | 0 | live 19 |
| `/log` | 6 | 0 | 0 | live 5, completion 6, katalog 3 |
| `/mpls/forwarding-table` | 12 | 0 | 0 | completion 12 |
| `/mpls/interface` | 7 | 0 | 0 | completion 7 |
| `/mpls/ldp` | 14 | 0 | 0 | completion 14 |
| `/mpls/mangle` | 8 | 0 | 0 | completion 8 |
| `/mpls/settings` | 5 | 0 | 0 | live 5 |
| `/ping` | 12 | 0 | 0 | katalog 12 |
| `/port/remote-access` | 16 | 0 | 0 | completion 16 |
| `/ppp` | 0 | 0 | 0 | tidak ada kolom terdeteksi |
| `/ppp/aaa` | 5 | 0 | 5 | live 5, katalog 4 |
| `/ppp/active` | 12 | 0 | 0 | completion 12, katalog 10 |
| `/ppp/l2tp-secret` | 4 | 0 | 0 | completion 4 |
| `/ppp/profile` | 39 | 38 | 38 | live 14, completion 39, katalog 14 |
| `/ppp/secret` | 18 | 15 | 15 | completion 18, katalog 13 |
| `/queue` | 0 | 0 | 0 | tidak ada kolom terdeteksi |
| `/queue/interface` | 5 | 0 | 2 | live 5, completion 5 |
| `/queue/simple` | 43 | 26 | 25 | completion 43, katalog 18 |
| `/queue/tree` | 23 | 14 | 14 | completion 23, katalog 17 |
| `/queue/type` | 52 | 51 | 51 | live 24, completion 52 |
| `/radius` | 19 | 10 | 10 | katalog 19 |
| `/radius/incoming` | 3 | 0 | 0 | live 3, katalog 2 |
| `/routing` | 0 | 0 | 0 | tidak ada kolom terdeteksi |
| `/routing/bfd` | 0 | 0 | 0 | tidak ada kolom terdeteksi |
| `/routing/bgp` | 0 | 0 | 0 | tidak ada kolom terdeteksi |
| `/routing/fantasy` | 18 | 17 | 17 | completion 18 |
| `/routing/filter` | 0 | 0 | 0 | tidak ada kolom terdeteksi |
| `/routing/gmp` | 7 | 7 | 7 | completion 7 |
| `/routing/id` | 10 | 7 | 7 | live 7, completion 10 |
| `/routing/igmp-proxy` | 3 | 0 | 0 | live 3 |
| `/routing/nexthop` | 23 | 0 | 0 | live 18, completion 23 |
| `/routing/ospf` | 0 | 0 | 0 | tidak ada kolom terdeteksi |
| `/routing/route` | 91 | 0 | 3 | live 18, completion 91 |
| `/routing/rpki` | 11 | 0 | 0 | completion 11 |
| `/routing/rule` | 15 | 14 | 13 | completion 15 |
| `/routing/settings` | 5 | 0 | 7 | live 5 |
| `/routing/stats` | 0 | 0 | 0 | tidak ada kolom terdeteksi |
| `/routing/table` | 8 | 0 | 0 | live 5, completion 8 |
| `/safe-mode` | 2 | 0 | 0 | live 2 |
| `/snmp/community` | 13 | 0 | 0 | live 12, completion 13 |
| `/special-login` | 5 | 0 | 0 | completion 5 |
| `/system` | 0 | 0 | 0 | tidak ada kolom terdeteksi |
| `/system/backup` | 0 | 0 | 0 | tidak ada kolom terdeteksi |
| `/system/clock` | 6 | 0 | 4 | live 6, katalog 4 |
| `/system/console` | 10 | 0 | 0 | live 8, completion 10 |
| `/system/default-configuration` | 3 | 0 | 0 | live 3 |
| `/system/device-mode` | 24 | 0 | 0 | live 24 |
| `/system/health` | 11 | 0 | 1 | live 2, katalog 18 |
| `/system/history` | 11 | 0 | 0 | completion 11 |
| `/system/identity` | 1 | 0 | 1 | live 1, katalog 1 |
| `/system/keymat-provider` | 14 | 0 | 0 | completion 14 |
| `/system/leds` | 8 | 0 | 0 | completion 8 |
| `/system/license` | 2 | 0 | 0 | live 2 |
| `/system/logging` | 10 | 7 | 7 | live 9, completion 10 |
| `/system/note` | 3 | 0 | 0 | live 3 |
| `/system/ntp` | 0 | 0 | 0 | tidak ada kolom terdeteksi |
| `/system/package` | 9 | 0 | 0 | live 8, completion 9, katalog 5 |
| `/system/package/update` | 4 | 0 | 1 | katalog 4 |
| `/system/regulatory` | 1 | 0 | 0 | live 1 |
| `/system/resource` | 17 | 0 | 0 | live 17, katalog 20 |
| `/system/resource/cpu` | 4 | 0 | 0 | katalog 8 |
| `/system/routerboard` | 7 | 0 | 0 | katalog 7 |
| `/system/rtrace` | 1 | 0 | 0 | live 1 |
| `/system/scheduler` | 13 | 10 | 10 | live 12, completion 13, katalog 10 |
| `/system/script` | 10 | 7 | 7 | completion 10, katalog 14 |
| `/system/ups` | 18 | 0 | 0 | completion 18 |
| `/system/watchdog` | 6 | 0 | 0 | live 6 |
| `/task` | 12 | 0 | 0 | completion 12 |
| `/tool` | 0 | 0 | 0 | tidak ada kolom terdeteksi |
| `/tool/bandwidth-server` | 6 | 0 | 6 | live 6 |
| `/tool/bandwidth-test` | 13 | 0 | 0 | katalog 13 |
| `/tool/e-mail` | 8 | 0 | 0 | live 8 |
| `/tool/fetch` | 4 | 0 | 0 | katalog 4 |
| `/tool/graphing` | 2 | 0 | 0 | live 2 |
| `/tool/mac-server` | 1 | 0 | 0 | live 1 |
| `/tool/netwatch` | 57 | 0 | 0 | completion 57, katalog 11 |
| `/tool/profile` | 3 | 0 | 0 | katalog 3 |
| `/tool/romon` | 3 | 0 | 3 | live 3 |
| `/tool/sms` | 11 | 0 | 0 | live 11 |
| `/tool/sniffer` | 32 | 0 | 0 | live 32 |
| `/tool/torch` | 11 | 0 | 0 | katalog 11 |
| `/tool/traceroute` | 9 | 0 | 0 | katalog 9 |
| `/tool/traffic-generator` | 7 | 0 | 0 | live 7 |
| `/tool/traffic-monitor` | 10 | 9 | 9 | completion 10 |
| `/user` | 6 | 9 | 9 | katalog 6 |
| `/user/aaa` | 5 | 0 | 5 | live 5 |
| `/user/active` | 8 | 1 | 0 | live 7, completion 8, katalog 6 |
| `/user/group` | 5 | 5 | 5 | live 4, completion 5 |
| `/user/settings` | 2 | 0 | 2 | live 2 |
| `/user/ssh-keys` | 8 | 0 | 0 | completion 8 |

## Struct per menu

### `/`

- File: `_root.json`
- Sumber field: tidak ada kolom terdeteksi

Command: `beep`(3), `break`(0), `continue`(0), `convert`(7), `delay`(1), `deserialize`(4), `do`(4), `error`(1), `execute`(3), `exit`(0), `export`(8), `find`(3), `for`(6), `foreach`(4), `global`(3), `grep`(6), `if`(3), `import`(4), `jobname`(0), `len`(1), `local`(3), `lock`(1), `nothing`(0), `onerror`(3), `parse`(1), `password`(4), `pick`(3), `put`(1), `quit`(0), `range`(2), `recursive-print`(0), `redo`(0), `resolve`(4), `retry`(4), `return`(1), `rndnum`(2), `rndstr`(2), `serialize`(6), `set`(3), `time`(1), `timestamp`(1), `toarray`(2), `tobool`(1), `tocrlf`(1), `toid`(1), `toip`(1), `toip6`(1), `tolf`(1), `tonsec`(1), `tonum`(1), `tostr`(1), `totime`(1), `typeof`(1), `undo`(0), `while`(2)

```go
// MenuSet adalah argumen `//set`.
type MenuSet struct {
	Do                           string `ros:"do,omitempty"`
	Name                         string `ros:"name,omitempty"`
	Value                        string `ros:"value,omitempty"`
}
```

### `/app`

- File: `app__cmds.json`
- Sumber field: tidak ada kolom terdeteksi

Command: `add`(16), `cleanup`(1), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(22), `remove`(1), `reset`(16), `restart`(1), `set`(16), `setup`(0), `unset`(2), `update`(1)

```go
// AppAdd adalah argumen `/app/add`.
type AppAdd struct {
	AutoUpdate                   string `ros:"auto-update,omitempty"`
	ContainerCommandLines        string `ros:"container-command-lines,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Devices                      string `ros:"devices,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Environment                  string `ros:"environment,omitempty"`
	ExtraMounts                  string `ros:"extra-mounts,omitempty"`
	FirewallRedirects            string `ros:"firewall-redirects,omitempty"`
	Network                      string `ros:"network,omitempty"`
	NetworkOutgoingAccess        string `ros:"network-outgoing-access,omitempty"`
	NetworkPvid                  string `ros:"network-pvid,omitempty"`
	RequiredHWDevices            string `ros:"required-hw-devices,omitempty"`
	RequiredMounts               string `ros:"required-mounts,omitempty"`
	Secrets                      string `ros:"secrets,omitempty"`
	UseHTTPS                     string `ros:"use-https,omitempty"`
	Yaml                         string `ros:"yaml,omitempty"`
}
```

```go
// AppSet adalah argumen `/app/set`.
type AppSet struct {
	AutoUpdate                   string `ros:"auto-update,omitempty"`
	ContainerCommandLines        string `ros:"container-command-lines,omitempty"`
	Devices                      string `ros:"devices,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Environment                  string `ros:"environment,omitempty"`
	ExtraMounts                  string `ros:"extra-mounts,omitempty"`
	FirewallRedirects            string `ros:"firewall-redirects,omitempty"`
	Network                      string `ros:"network,omitempty"`
	NetworkOutgoingAccess        string `ros:"network-outgoing-access,omitempty"`
	NetworkPvid                  string `ros:"network-pvid,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	RequiredHWDevices            string `ros:"required-hw-devices,omitempty"`
	RequiredMounts               string `ros:"required-mounts,omitempty"`
	Secrets                      string `ros:"secrets,omitempty"`
	UseHTTPS                     string `ros:"use-https,omitempty"`
	Yaml                         string `ros:"yaml,omitempty"`
}
```

### `/app/network`

- File: `app_network.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command or directory (app)

Command: `find`(1), `get`(4), `print`(21), `reset`(3)

_Tidak ada struct: menu ini bukan tabel yang bisa di-print atau tidak tersedia di router uji (mis. paket tidak terpasang)._

### `/app/settings`

- File: `app_settings.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command or directory (app)

Command: `edit`(1), `export`(8), `get`(3), `print`(5), `reset`(9), `set`(9)

```go
// AppSettingsSet adalah argumen `/app/settings/set`.
type AppSettingsSet struct {
	AppStoreUrls                 string `ros:"app-store-urls,omitempty"`
	AutoUpdate                   string `ros:"auto-update,omitempty"`
	Disk                         string `ros:"disk,omitempty"`
	DownloadPath                 string `ros:"download-path,omitempty"`
	LANBridge                    string `ros:"lan-bridge,omitempty"`
	MediaPath                    string `ros:"media-path,omitempty"`
	RegistryMirrors              string `ros:"registry-mirrors,omitempty"`
	RouterIP                     string `ros:"router-ip,omitempty"`
	ShowInWebfig                 string `ros:"show-in-webfig,omitempty"`
}
```

### `/caps-man`

- File: `caps-man__cmds.json`
- Sumber field: tidak ada kolom terdeteksi

Command: `export`(8)

_Tidak ada struct: menu ini bukan tabel yang bisa di-print atau tidak tersedia di router uji (mis. paket tidak terpasang)._

### `/caps-man/aaa`

- File: `caps-man_aaa.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command or directory (caps-man)

Command: `edit`(1), `export`(8), `get`(3), `print`(5), `reset`(5), `set`(5)

```go
// CapsmanAaaSet adalah argumen `/caps-man/aaa/set`.
type CapsmanAaaSet struct {
	CalledFormat                 string `ros:"called-format,omitempty"`
	InterimUpdate                string `ros:"interim-update,omitempty"`
	MACCaching                   string `ros:"mac-caching,omitempty"`
	MACFormat                    string `ros:"mac-format,omitempty"`
	MACMode                      string `ros:"mac-mode,omitempty"`
}
```

### `/caps-man/access-list`

- File: `caps-man_access-list.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command or directory (caps-man)

Command: `add`(19), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `move`(2), `print`(22), `remove`(1), `reset`(18), `set`(18), `unset`(2)

```go
// CapsmanAccesslistAdd adalah argumen `/caps-man/access-list/add`.
type CapsmanAccesslistAdd struct {
	Action                       string `ros:"action,omitempty"`
	AllowSignalOutOfRange        string `ros:"allow-signal-out-of-range,omitempty"`
	ApTXLimit                    string `ros:"ap-tx-limit,omitempty"`
	ClientToClientForwarding     string `ros:"client-to-client-forwarding,omitempty"`
	ClientTXLimit                string `ros:"client-tx-limit,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	MACAddress                   string `ros:"mac-address,omitempty"`
	MACAddressMask               string `ros:"mac-address-mask,omitempty"`
	PlaceBefore                  string `ros:"place-before,omitempty"`
	PrivatePassphrase            string `ros:"private-passphrase,omitempty"`
	RadiusAccounting             string `ros:"radius-accounting,omitempty"`
	SignalRange                  string `ros:"signal-range,omitempty"`
	SsidRegexp                   string `ros:"ssid-regexp,omitempty"`
	Time                         string `ros:"time,omitempty"`
	VLANID                       string `ros:"vlan-id,omitempty"`
	VLANMode                     string `ros:"vlan-mode,omitempty"`
}
```

```go
// CapsmanAccesslistSet adalah argumen `/caps-man/access-list/set`.
type CapsmanAccesslistSet struct {
	Action                       string `ros:"action,omitempty"`
	AllowSignalOutOfRange        string `ros:"allow-signal-out-of-range,omitempty"`
	ApTXLimit                    string `ros:"ap-tx-limit,omitempty"`
	ClientToClientForwarding     string `ros:"client-to-client-forwarding,omitempty"`
	ClientTXLimit                string `ros:"client-tx-limit,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	MACAddress                   string `ros:"mac-address,omitempty"`
	MACAddressMask               string `ros:"mac-address-mask,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	PrivatePassphrase            string `ros:"private-passphrase,omitempty"`
	RadiusAccounting             string `ros:"radius-accounting,omitempty"`
	SignalRange                  string `ros:"signal-range,omitempty"`
	SsidRegexp                   string `ros:"ssid-regexp,omitempty"`
	Time                         string `ros:"time,omitempty"`
	VLANID                       string `ros:"vlan-id,omitempty"`
	VLANMode                     string `ros:"vlan-mode,omitempty"`
}
```

### `/caps-man/actual-interface-configuration`

- File: `caps-man_actual-interface-configuration.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command or directory (caps-man)

Command: `comment`(2), `disable`(1), `edit`(2), `enable`(1), `find`(1), `get`(4), `print`(21), `reset`(57), `set`(57), `unset`(2)

```go
// CapsmanActualinterfaceconfigurationSet adalah argumen `/caps-man/actual-interface-configuration/set`.
type CapsmanActualinterfaceconfigurationSet struct {
	ARP                          string `ros:"arp,omitempty"`
	ARPTimeout                   string `ros:"arp-timeout,omitempty"`
	ChannelBand                  string `ros:"channel.band,omitempty"`
	ChannelControlChannelWidth   string `ros:"channel.control-channel-width,omitempty"`
	ChannelExtensionChannel      string `ros:"channel.extension-channel,omitempty"`
	ChannelFrequency             string `ros:"channel.frequency,omitempty"`
	ChannelReselectInterval      string `ros:"channel.reselect-interval,omitempty"`
	ChannelSaveSelected          string `ros:"channel.save-selected,omitempty"`
	ChannelSecondaryFrequency    string `ros:"channel.secondary-frequency,omitempty"`
	ChannelSkipDfsChannels       string `ros:"channel.skip-dfs-channels,omitempty"`
	ChannelTXPower               string `ros:"channel.tx-power,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	ConfigurationCountry         string `ros:"configuration.country,omitempty"`
	ConfigurationDisconnectTimeout string `ros:"configuration.disconnect-timeout,omitempty"`
	ConfigurationDistance        string `ros:"configuration.distance,omitempty"`
	ConfigurationFrameLifetime   string `ros:"configuration.frame-lifetime,omitempty"`
	ConfigurationGuardInterval   string `ros:"configuration.guard-interval,omitempty"`
	ConfigurationHideSsid        string `ros:"configuration.hide-ssid,omitempty"`
	ConfigurationHWProtectionMode string `ros:"configuration.hw-protection-mode,omitempty"`
	ConfigurationHWRetries       string `ros:"configuration.hw-retries,omitempty"`
	ConfigurationInstallation    string `ros:"configuration.installation,omitempty"`
	ConfigurationKeepaliveFrames string `ros:"configuration.keepalive-frames,omitempty"`
	ConfigurationLoadBalancingGroup string `ros:"configuration.load-balancing-group,omitempty"`
	ConfigurationMaxStaCount     string `ros:"configuration.max-sta-count,omitempty"`
	ConfigurationMode            string `ros:"configuration.mode,omitempty"`
	ConfigurationMulticastHelper string `ros:"configuration.multicast-helper,omitempty"`
	ConfigurationRXChains        string `ros:"configuration.rx-chains,omitempty"`
	ConfigurationSsid            string `ros:"configuration.ssid,omitempty"`
	ConfigurationTXChains        string `ros:"configuration.tx-chains,omitempty"`
	DatapathBridge               string `ros:"datapath.bridge,omitempty"`
	DatapathBridgeCost           string `ros:"datapath.bridge-cost,omitempty"`
	DatapathBridgeHorizon        string `ros:"datapath.bridge-horizon,omitempty"`
	DatapathClientToClientForwarding string `ros:"datapath.client-to-client-forwarding,omitempty"`
	DatapathInterfaceList        string `ros:"datapath.interface-list,omitempty"`
	DatapathLocalForwarding      string `ros:"datapath.local-forwarding,omitempty"`
	DatapathOpenflowSwitch       string `ros:"datapath.openflow-switch,omitempty"`
	DatapathVLANID               string `ros:"datapath.vlan-id,omitempty"`
	DatapathVLANMode             string `ros:"datapath.vlan-mode,omitempty"`
	DisableRunningCheck          string `ros:"disable-running-check,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	L2MTU                        string `ros:"l2mtu,omitempty"`
	MACAddress                   string `ros:"mac-address,omitempty"`
	MasterInterface              string `ros:"master-interface,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	RadioMAC                     string `ros:"radio-mac,omitempty"`
	SecurityAuthenticationTypes  string `ros:"security.authentication-types,omitempty"`
	SecurityDisablePmkid         string `ros:"security.disable-pmkid,omitempty"`
	SecurityEapMethods           string `ros:"security.eap-methods,omitempty"`
	SecurityEapRadiusAccounting  string `ros:"security.eap-radius-accounting,omitempty"`
	SecurityEncryption           string `ros:"security.encryption,omitempty"`
	SecurityGroupEncryption      string `ros:"security.group-encryption,omitempty"`
	SecurityGroupKeyUpdate       string `ros:"security.group-key-update,omitempty"`
	SecurityPassphrase           string `ros:"security.passphrase,omitempty"`
	SecurityTLSCertificate       string `ros:"security.tls-certificate,omitempty"`
	SecurityTLSMode              string `ros:"security.tls-mode,omitempty"`
}
```

### `/caps-man/channel`

- File: `caps-man_channel.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command or directory (caps-man)

Command: `add`(12), `comment`(2), `edit`(2), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(12), `set`(12), `unset`(2)

```go
// CapsmanChannelAdd adalah argumen `/caps-man/channel/add`.
type CapsmanChannelAdd struct {
	Band                         string `ros:"band,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	ControlChannelWidth          string `ros:"control-channel-width,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	ExtensionChannel             string `ros:"extension-channel,omitempty"`
	Frequency                    string `ros:"frequency,omitempty"`
	Name                         string `ros:"name,omitempty"`
	ReselectInterval             string `ros:"reselect-interval,omitempty"`
	SaveSelected                 string `ros:"save-selected,omitempty"`
	SecondaryFrequency           string `ros:"secondary-frequency,omitempty"`
	SkipDfsChannels              string `ros:"skip-dfs-channels,omitempty"`
	TXPower                      string `ros:"tx-power,omitempty"`
}
```

```go
// CapsmanChannelSet adalah argumen `/caps-man/channel/set`.
type CapsmanChannelSet struct {
	Band                         string `ros:"band,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	ControlChannelWidth          string `ros:"control-channel-width,omitempty"`
	ExtensionChannel             string `ros:"extension-channel,omitempty"`
	Frequency                    string `ros:"frequency,omitempty"`
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	ReselectInterval             string `ros:"reselect-interval,omitempty"`
	SaveSelected                 string `ros:"save-selected,omitempty"`
	SecondaryFrequency           string `ros:"secondary-frequency,omitempty"`
	SkipDfsChannels              string `ros:"skip-dfs-channels,omitempty"`
	TXPower                      string `ros:"tx-power,omitempty"`
}
```

### `/caps-man/configuration`

- File: `caps-man_configuration.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command or directory (caps-man)

Command: `add`(61), `comment`(2), `edit`(2), `export`(8), `find`(1), `get`(4), `print`(22), `remove`(1), `reset`(61), `set`(61), `unset`(2)

```go
// CapsmanConfigurationAdd adalah argumen `/caps-man/configuration/add`.
type CapsmanConfigurationAdd struct {
	Channel                      string `ros:"channel,omitempty"`
	ChannelBand                  string `ros:"channel.band,omitempty"`
	ChannelControlChannelWidth   string `ros:"channel.control-channel-width,omitempty"`
	ChannelExtensionChannel      string `ros:"channel.extension-channel,omitempty"`
	ChannelFrequency             string `ros:"channel.frequency,omitempty"`
	ChannelReselectInterval      string `ros:"channel.reselect-interval,omitempty"`
	ChannelSaveSelected          string `ros:"channel.save-selected,omitempty"`
	ChannelSecondaryFrequency    string `ros:"channel.secondary-frequency,omitempty"`
	ChannelSkipDfsChannels       string `ros:"channel.skip-dfs-channels,omitempty"`
	ChannelTXPower               string `ros:"channel.tx-power,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Country                      string `ros:"country,omitempty"`
	Datapath                     string `ros:"datapath,omitempty"`
	DatapathARP                  string `ros:"datapath.arp,omitempty"`
	DatapathBridge               string `ros:"datapath.bridge,omitempty"`
	DatapathBridgeCost           string `ros:"datapath.bridge-cost,omitempty"`
	DatapathBridgeHorizon        string `ros:"datapath.bridge-horizon,omitempty"`
	DatapathClientToClientForwarding string `ros:"datapath.client-to-client-forwarding,omitempty"`
	DatapathInterfaceList        string `ros:"datapath.interface-list,omitempty"`
	DatapathL2MTU                string `ros:"datapath.l2mtu,omitempty"`
	DatapathLocalForwarding      string `ros:"datapath.local-forwarding,omitempty"`
	DatapathMTU                  string `ros:"datapath.mtu,omitempty"`
	DatapathOpenflowSwitch       string `ros:"datapath.openflow-switch,omitempty"`
	DatapathVLANID               string `ros:"datapath.vlan-id,omitempty"`
	DatapathVLANMode             string `ros:"datapath.vlan-mode,omitempty"`
	DisconnectTimeout            string `ros:"disconnect-timeout,omitempty"`
	Distance                     string `ros:"distance,omitempty"`
	FrameLifetime                string `ros:"frame-lifetime,omitempty"`
	GuardInterval                string `ros:"guard-interval,omitempty"`
	HideSsid                     string `ros:"hide-ssid,omitempty"`
	HWProtectionMode             string `ros:"hw-protection-mode,omitempty"`
	HWRetries                    string `ros:"hw-retries,omitempty"`
	Installation                 string `ros:"installation,omitempty"`
	KeepaliveFrames              string `ros:"keepalive-frames,omitempty"`
	LoadBalancingGroup           string `ros:"load-balancing-group,omitempty"`
	MaxStaCount                  string `ros:"max-sta-count,omitempty"`
	Mode                         string `ros:"mode,omitempty"`
	MulticastHelper              string `ros:"multicast-helper,omitempty"`
	Name                         string `ros:"name,omitempty"`
	Rates                        string `ros:"rates,omitempty"`
	RatesBasic                   string `ros:"rates.basic,omitempty"`
	RatesHtBasicMcs              string `ros:"rates.ht-basic-mcs,omitempty"`
	RatesHtSupportedMcs          string `ros:"rates.ht-supported-mcs,omitempty"`
	RatesSupported               string `ros:"rates.supported,omitempty"`
	RatesVhtBasicMcs             string `ros:"rates.vht-basic-mcs,omitempty"`
	RatesVhtSupportedMcs         string `ros:"rates.vht-supported-mcs,omitempty"`
	RXChains                     string `ros:"rx-chains,omitempty"`
	Security                     string `ros:"security,omitempty"`
	SecurityAuthenticationTypes  string `ros:"security.authentication-types,omitempty"`
	SecurityDisablePmkid         string `ros:"security.disable-pmkid,omitempty"`
	SecurityEapMethods           string `ros:"security.eap-methods,omitempty"`
	SecurityEapRadiusAccounting  string `ros:"security.eap-radius-accounting,omitempty"`
	SecurityEncryption           string `ros:"security.encryption,omitempty"`
	SecurityGroupEncryption      string `ros:"security.group-encryption,omitempty"`
	SecurityGroupKeyUpdate       string `ros:"security.group-key-update,omitempty"`
	SecurityPassphrase           string `ros:"security.passphrase,omitempty"`
	SecurityTLSCertificate       string `ros:"security.tls-certificate,omitempty"`
	SecurityTLSMode              string `ros:"security.tls-mode,omitempty"`
	Ssid                         string `ros:"ssid,omitempty"`
	TXChains                     string `ros:"tx-chains,omitempty"`
}
```

```go
// CapsmanConfigurationSet adalah argumen `/caps-man/configuration/set`.
type CapsmanConfigurationSet struct {
	Channel                      string `ros:"channel,omitempty"`
	ChannelBand                  string `ros:"channel.band,omitempty"`
	ChannelControlChannelWidth   string `ros:"channel.control-channel-width,omitempty"`
	ChannelExtensionChannel      string `ros:"channel.extension-channel,omitempty"`
	ChannelFrequency             string `ros:"channel.frequency,omitempty"`
	ChannelReselectInterval      string `ros:"channel.reselect-interval,omitempty"`
	ChannelSaveSelected          string `ros:"channel.save-selected,omitempty"`
	ChannelSecondaryFrequency    string `ros:"channel.secondary-frequency,omitempty"`
	ChannelSkipDfsChannels       string `ros:"channel.skip-dfs-channels,omitempty"`
	ChannelTXPower               string `ros:"channel.tx-power,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	Country                      string `ros:"country,omitempty"`
	Datapath                     string `ros:"datapath,omitempty"`
	DatapathARP                  string `ros:"datapath.arp,omitempty"`
	DatapathBridge               string `ros:"datapath.bridge,omitempty"`
	DatapathBridgeCost           string `ros:"datapath.bridge-cost,omitempty"`
	DatapathBridgeHorizon        string `ros:"datapath.bridge-horizon,omitempty"`
	DatapathClientToClientForwarding string `ros:"datapath.client-to-client-forwarding,omitempty"`
	DatapathInterfaceList        string `ros:"datapath.interface-list,omitempty"`
	DatapathL2MTU                string `ros:"datapath.l2mtu,omitempty"`
	DatapathLocalForwarding      string `ros:"datapath.local-forwarding,omitempty"`
	DatapathMTU                  string `ros:"datapath.mtu,omitempty"`
	DatapathOpenflowSwitch       string `ros:"datapath.openflow-switch,omitempty"`
	DatapathVLANID               string `ros:"datapath.vlan-id,omitempty"`
	DatapathVLANMode             string `ros:"datapath.vlan-mode,omitempty"`
	DisconnectTimeout            string `ros:"disconnect-timeout,omitempty"`
	Distance                     string `ros:"distance,omitempty"`
	FrameLifetime                string `ros:"frame-lifetime,omitempty"`
	GuardInterval                string `ros:"guard-interval,omitempty"`
	HideSsid                     string `ros:"hide-ssid,omitempty"`
	HWProtectionMode             string `ros:"hw-protection-mode,omitempty"`
	HWRetries                    string `ros:"hw-retries,omitempty"`
	Installation                 string `ros:"installation,omitempty"`
	KeepaliveFrames              string `ros:"keepalive-frames,omitempty"`
	LoadBalancingGroup           string `ros:"load-balancing-group,omitempty"`
	MaxStaCount                  string `ros:"max-sta-count,omitempty"`
	Mode                         string `ros:"mode,omitempty"`
	MulticastHelper              string `ros:"multicast-helper,omitempty"`
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	Rates                        string `ros:"rates,omitempty"`
	RatesBasic                   string `ros:"rates.basic,omitempty"`
	RatesHtBasicMcs              string `ros:"rates.ht-basic-mcs,omitempty"`
	RatesHtSupportedMcs          string `ros:"rates.ht-supported-mcs,omitempty"`
	RatesSupported               string `ros:"rates.supported,omitempty"`
	RatesVhtBasicMcs             string `ros:"rates.vht-basic-mcs,omitempty"`
	RatesVhtSupportedMcs         string `ros:"rates.vht-supported-mcs,omitempty"`
	RXChains                     string `ros:"rx-chains,omitempty"`
	Security                     string `ros:"security,omitempty"`
	SecurityAuthenticationTypes  string `ros:"security.authentication-types,omitempty"`
	SecurityDisablePmkid         string `ros:"security.disable-pmkid,omitempty"`
	SecurityEapMethods           string `ros:"security.eap-methods,omitempty"`
	SecurityEapRadiusAccounting  string `ros:"security.eap-radius-accounting,omitempty"`
	SecurityEncryption           string `ros:"security.encryption,omitempty"`
	SecurityGroupEncryption      string `ros:"security.group-encryption,omitempty"`
	SecurityGroupKeyUpdate       string `ros:"security.group-key-update,omitempty"`
	SecurityPassphrase           string `ros:"security.passphrase,omitempty"`
	SecurityTLSCertificate       string `ros:"security.tls-certificate,omitempty"`
	SecurityTLSMode              string `ros:"security.tls-mode,omitempty"`
	Ssid                         string `ros:"ssid,omitempty"`
	TXChains                     string `ros:"tx-chains,omitempty"`
}
```

### `/caps-man/datapath`

- File: `caps-man_datapath.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command or directory (caps-man)

Command: `add`(15), `comment`(2), `edit`(2), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(15), `set`(15), `unset`(2)

```go
// CapsmanDatapathAdd adalah argumen `/caps-man/datapath/add`.
type CapsmanDatapathAdd struct {
	ARP                          string `ros:"arp,omitempty"`
	Bridge                       string `ros:"bridge,omitempty"`
	BridgeCost                   string `ros:"bridge-cost,omitempty"`
	BridgeHorizon                string `ros:"bridge-horizon,omitempty"`
	ClientToClientForwarding     string `ros:"client-to-client-forwarding,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	InterfaceList                string `ros:"interface-list,omitempty"`
	L2MTU                        string `ros:"l2mtu,omitempty"`
	LocalForwarding              string `ros:"local-forwarding,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	Name                         string `ros:"name,omitempty"`
	OpenflowSwitch               string `ros:"openflow-switch,omitempty"`
	VLANID                       string `ros:"vlan-id,omitempty"`
	VLANMode                     string `ros:"vlan-mode,omitempty"`
}
```

```go
// CapsmanDatapathSet adalah argumen `/caps-man/datapath/set`.
type CapsmanDatapathSet struct {
	ARP                          string `ros:"arp,omitempty"`
	Bridge                       string `ros:"bridge,omitempty"`
	BridgeCost                   string `ros:"bridge-cost,omitempty"`
	BridgeHorizon                string `ros:"bridge-horizon,omitempty"`
	ClientToClientForwarding     string `ros:"client-to-client-forwarding,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	InterfaceList                string `ros:"interface-list,omitempty"`
	L2MTU                        string `ros:"l2mtu,omitempty"`
	LocalForwarding              string `ros:"local-forwarding,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	OpenflowSwitch               string `ros:"openflow-switch,omitempty"`
	VLANID                       string `ros:"vlan-id,omitempty"`
	VLANMode                     string `ros:"vlan-mode,omitempty"`
}
```

### `/caps-man/interface`

- File: `caps-man_interface.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command or directory (caps-man)

Command: `add`(69), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `hw-info`(2), `possible-channels`(2), `print`(23), `remove`(1), `reselect-channel`(2), `reset`(69), `scan`(8), `set`(69), `unset`(2)

```go
// CapsmanInterfaceAdd adalah argumen `/caps-man/interface/add`.
type CapsmanInterfaceAdd struct {
	ARP                          string `ros:"arp,omitempty"`
	ARPTimeout                   string `ros:"arp-timeout,omitempty"`
	Channel                      string `ros:"channel,omitempty"`
	ChannelBand                  string `ros:"channel.band,omitempty"`
	ChannelControlChannelWidth   string `ros:"channel.control-channel-width,omitempty"`
	ChannelExtensionChannel      string `ros:"channel.extension-channel,omitempty"`
	ChannelFrequency             string `ros:"channel.frequency,omitempty"`
	ChannelReselectInterval      string `ros:"channel.reselect-interval,omitempty"`
	ChannelSaveSelected          string `ros:"channel.save-selected,omitempty"`
	ChannelSecondaryFrequency    string `ros:"channel.secondary-frequency,omitempty"`
	ChannelSkipDfsChannels       string `ros:"channel.skip-dfs-channels,omitempty"`
	ChannelTXPower               string `ros:"channel.tx-power,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	Configuration                string `ros:"configuration,omitempty"`
	ConfigurationCountry         string `ros:"configuration.country,omitempty"`
	ConfigurationDisconnectTimeout string `ros:"configuration.disconnect-timeout,omitempty"`
	ConfigurationDistance        string `ros:"configuration.distance,omitempty"`
	ConfigurationFrameLifetime   string `ros:"configuration.frame-lifetime,omitempty"`
	ConfigurationGuardInterval   string `ros:"configuration.guard-interval,omitempty"`
	ConfigurationHideSsid        string `ros:"configuration.hide-ssid,omitempty"`
	ConfigurationHWProtectionMode string `ros:"configuration.hw-protection-mode,omitempty"`
	ConfigurationHWRetries       string `ros:"configuration.hw-retries,omitempty"`
	ConfigurationInstallation    string `ros:"configuration.installation,omitempty"`
	ConfigurationKeepaliveFrames string `ros:"configuration.keepalive-frames,omitempty"`
	ConfigurationLoadBalancingGroup string `ros:"configuration.load-balancing-group,omitempty"`
	ConfigurationMaxStaCount     string `ros:"configuration.max-sta-count,omitempty"`
	ConfigurationMode            string `ros:"configuration.mode,omitempty"`
	ConfigurationMulticastHelper string `ros:"configuration.multicast-helper,omitempty"`
	ConfigurationRXChains        string `ros:"configuration.rx-chains,omitempty"`
	ConfigurationSsid            string `ros:"configuration.ssid,omitempty"`
	ConfigurationTXChains        string `ros:"configuration.tx-chains,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Datapath                     string `ros:"datapath,omitempty"`
	DatapathBridge               string `ros:"datapath.bridge,omitempty"`
	DatapathBridgeCost           string `ros:"datapath.bridge-cost,omitempty"`
	DatapathBridgeHorizon        string `ros:"datapath.bridge-horizon,omitempty"`
	DatapathClientToClientForwarding string `ros:"datapath.client-to-client-forwarding,omitempty"`
	DatapathInterfaceList        string `ros:"datapath.interface-list,omitempty"`
	DatapathLocalForwarding      string `ros:"datapath.local-forwarding,omitempty"`
	DatapathOpenflowSwitch       string `ros:"datapath.openflow-switch,omitempty"`
	DatapathVLANID               string `ros:"datapath.vlan-id,omitempty"`
	DatapathVLANMode             string `ros:"datapath.vlan-mode,omitempty"`
	DisableRunningCheck          string `ros:"disable-running-check,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	L2MTU                        string `ros:"l2mtu,omitempty"`
	MACAddress                   string `ros:"mac-address,omitempty"`
	MasterInterface              string `ros:"master-interface,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	Name                         string `ros:"name,omitempty"`
	RadioMAC                     string `ros:"radio-mac,omitempty"`
	RadioName                    string `ros:"radio-name,omitempty"`
	Rates                        string `ros:"rates,omitempty"`
	RatesBasic                   string `ros:"rates.basic,omitempty"`
	RatesHtBasicMcs              string `ros:"rates.ht-basic-mcs,omitempty"`
	RatesHtSupportedMcs          string `ros:"rates.ht-supported-mcs,omitempty"`
	RatesSupported               string `ros:"rates.supported,omitempty"`
	RatesVhtBasicMcs             string `ros:"rates.vht-basic-mcs,omitempty"`
	RatesVhtSupportedMcs         string `ros:"rates.vht-supported-mcs,omitempty"`
	Security                     string `ros:"security,omitempty"`
	SecurityAuthenticationTypes  string `ros:"security.authentication-types,omitempty"`
	SecurityDisablePmkid         string `ros:"security.disable-pmkid,omitempty"`
	SecurityEapMethods           string `ros:"security.eap-methods,omitempty"`
	SecurityEapRadiusAccounting  string `ros:"security.eap-radius-accounting,omitempty"`
	SecurityEncryption           string `ros:"security.encryption,omitempty"`
	SecurityGroupEncryption      string `ros:"security.group-encryption,omitempty"`
	SecurityGroupKeyUpdate       string `ros:"security.group-key-update,omitempty"`
	SecurityPassphrase           string `ros:"security.passphrase,omitempty"`
	SecurityTLSCertificate       string `ros:"security.tls-certificate,omitempty"`
	SecurityTLSMode              string `ros:"security.tls-mode,omitempty"`
}
```

```go
// CapsmanInterfaceSet adalah argumen `/caps-man/interface/set`.
type CapsmanInterfaceSet struct {
	ARP                          string `ros:"arp,omitempty"`
	ARPTimeout                   string `ros:"arp-timeout,omitempty"`
	Channel                      string `ros:"channel,omitempty"`
	ChannelBand                  string `ros:"channel.band,omitempty"`
	ChannelControlChannelWidth   string `ros:"channel.control-channel-width,omitempty"`
	ChannelExtensionChannel      string `ros:"channel.extension-channel,omitempty"`
	ChannelFrequency             string `ros:"channel.frequency,omitempty"`
	ChannelReselectInterval      string `ros:"channel.reselect-interval,omitempty"`
	ChannelSaveSelected          string `ros:"channel.save-selected,omitempty"`
	ChannelSecondaryFrequency    string `ros:"channel.secondary-frequency,omitempty"`
	ChannelSkipDfsChannels       string `ros:"channel.skip-dfs-channels,omitempty"`
	ChannelTXPower               string `ros:"channel.tx-power,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	Configuration                string `ros:"configuration,omitempty"`
	ConfigurationCountry         string `ros:"configuration.country,omitempty"`
	ConfigurationDisconnectTimeout string `ros:"configuration.disconnect-timeout,omitempty"`
	ConfigurationDistance        string `ros:"configuration.distance,omitempty"`
	ConfigurationFrameLifetime   string `ros:"configuration.frame-lifetime,omitempty"`
	ConfigurationGuardInterval   string `ros:"configuration.guard-interval,omitempty"`
	ConfigurationHideSsid        string `ros:"configuration.hide-ssid,omitempty"`
	ConfigurationHWProtectionMode string `ros:"configuration.hw-protection-mode,omitempty"`
	ConfigurationHWRetries       string `ros:"configuration.hw-retries,omitempty"`
	ConfigurationInstallation    string `ros:"configuration.installation,omitempty"`
	ConfigurationKeepaliveFrames string `ros:"configuration.keepalive-frames,omitempty"`
	ConfigurationLoadBalancingGroup string `ros:"configuration.load-balancing-group,omitempty"`
	ConfigurationMaxStaCount     string `ros:"configuration.max-sta-count,omitempty"`
	ConfigurationMode            string `ros:"configuration.mode,omitempty"`
	ConfigurationMulticastHelper string `ros:"configuration.multicast-helper,omitempty"`
	ConfigurationRXChains        string `ros:"configuration.rx-chains,omitempty"`
	ConfigurationSsid            string `ros:"configuration.ssid,omitempty"`
	ConfigurationTXChains        string `ros:"configuration.tx-chains,omitempty"`
	Datapath                     string `ros:"datapath,omitempty"`
	DatapathBridge               string `ros:"datapath.bridge,omitempty"`
	DatapathBridgeCost           string `ros:"datapath.bridge-cost,omitempty"`
	DatapathBridgeHorizon        string `ros:"datapath.bridge-horizon,omitempty"`
	DatapathClientToClientForwarding string `ros:"datapath.client-to-client-forwarding,omitempty"`
	DatapathInterfaceList        string `ros:"datapath.interface-list,omitempty"`
	DatapathLocalForwarding      string `ros:"datapath.local-forwarding,omitempty"`
	DatapathOpenflowSwitch       string `ros:"datapath.openflow-switch,omitempty"`
	DatapathVLANID               string `ros:"datapath.vlan-id,omitempty"`
	DatapathVLANMode             string `ros:"datapath.vlan-mode,omitempty"`
	DisableRunningCheck          string `ros:"disable-running-check,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	L2MTU                        string `ros:"l2mtu,omitempty"`
	MACAddress                   string `ros:"mac-address,omitempty"`
	MasterInterface              string `ros:"master-interface,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	RadioMAC                     string `ros:"radio-mac,omitempty"`
	RadioName                    string `ros:"radio-name,omitempty"`
	Rates                        string `ros:"rates,omitempty"`
	RatesBasic                   string `ros:"rates.basic,omitempty"`
	RatesHtBasicMcs              string `ros:"rates.ht-basic-mcs,omitempty"`
	RatesHtSupportedMcs          string `ros:"rates.ht-supported-mcs,omitempty"`
	RatesSupported               string `ros:"rates.supported,omitempty"`
	RatesVhtBasicMcs             string `ros:"rates.vht-basic-mcs,omitempty"`
	RatesVhtSupportedMcs         string `ros:"rates.vht-supported-mcs,omitempty"`
	Security                     string `ros:"security,omitempty"`
	SecurityAuthenticationTypes  string `ros:"security.authentication-types,omitempty"`
	SecurityDisablePmkid         string `ros:"security.disable-pmkid,omitempty"`
	SecurityEapMethods           string `ros:"security.eap-methods,omitempty"`
	SecurityEapRadiusAccounting  string `ros:"security.eap-radius-accounting,omitempty"`
	SecurityEncryption           string `ros:"security.encryption,omitempty"`
	SecurityGroupEncryption      string `ros:"security.group-encryption,omitempty"`
	SecurityGroupKeyUpdate       string `ros:"security.group-key-update,omitempty"`
	SecurityPassphrase           string `ros:"security.passphrase,omitempty"`
	SecurityTLSCertificate       string `ros:"security.tls-certificate,omitempty"`
	SecurityTLSMode              string `ros:"security.tls-mode,omitempty"`
}
```

### `/caps-man/manager`

- File: `caps-man_manager.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command or directory (caps-man)

Command: `edit`(1), `export`(8), `get`(3), `print`(5), `reset`(6), `set`(6)

```go
// CapsmanManagerSet adalah argumen `/caps-man/manager/set`.
type CapsmanManagerSet struct {
	CaCertificate                string `ros:"ca-certificate,omitempty"`
	Certificate                  string `ros:"certificate,omitempty"`
	Enabled                      string `ros:"enabled,omitempty"`
	PackagePath                  string `ros:"package-path,omitempty"`
	RequirePeerCertificate       string `ros:"require-peer-certificate,omitempty"`
	UpgradePolicy                string `ros:"upgrade-policy,omitempty"`
}
```

### `/caps-man/provisioning`

- File: `caps-man_provisioning.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command or directory (caps-man)

Command: `add`(14), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `move`(2), `print`(21), `remove`(1), `reset`(13), `set`(13)

```go
// CapsmanProvisioningAdd adalah argumen `/caps-man/provisioning/add`.
type CapsmanProvisioningAdd struct {
	Action                       string `ros:"action,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	CommonNameRegexp             string `ros:"common-name-regexp,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	HWSupportedModes             string `ros:"hw-supported-modes,omitempty"`
	IdentityRegexp               string `ros:"identity-regexp,omitempty"`
	IPAddressRanges              string `ros:"ip-address-ranges,omitempty"`
	MasterConfiguration          string `ros:"master-configuration,omitempty"`
	NameFormat                   string `ros:"name-format,omitempty"`
	NamePrefix                   string `ros:"name-prefix,omitempty"`
	PlaceBefore                  string `ros:"place-before,omitempty"`
	RadioMAC                     string `ros:"radio-mac,omitempty"`
	SlaveConfigurations          string `ros:"slave-configurations,omitempty"`
}
```

```go
// CapsmanProvisioningSet adalah argumen `/caps-man/provisioning/set`.
type CapsmanProvisioningSet struct {
	Action                       string `ros:"action,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	CommonNameRegexp             string `ros:"common-name-regexp,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	HWSupportedModes             string `ros:"hw-supported-modes,omitempty"`
	IdentityRegexp               string `ros:"identity-regexp,omitempty"`
	IPAddressRanges              string `ros:"ip-address-ranges,omitempty"`
	MasterConfiguration          string `ros:"master-configuration,omitempty"`
	NameFormat                   string `ros:"name-format,omitempty"`
	NamePrefix                   string `ros:"name-prefix,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	RadioMAC                     string `ros:"radio-mac,omitempty"`
	SlaveConfigurations          string `ros:"slave-configurations,omitempty"`
}
```

### `/caps-man/radio`

- File: `caps-man_radio.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command or directory (caps-man)

Command: `find`(1), `get`(4), `hw-info`(2), `print`(21), `provision`(1)

_Tidak ada struct: menu ini bukan tabel yang bisa di-print atau tidak tersedia di router uji (mis. paket tidak terpasang)._

### `/caps-man/rates`

- File: `caps-man_rates.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command or directory (caps-man)

Command: `add`(9), `comment`(2), `edit`(2), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(9), `set`(9), `unset`(2)

```go
// CapsmanRatesAdd adalah argumen `/caps-man/rates/add`.
type CapsmanRatesAdd struct {
	Basic                        string `ros:"basic,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	HtBasicMcs                   string `ros:"ht-basic-mcs,omitempty"`
	HtSupportedMcs               string `ros:"ht-supported-mcs,omitempty"`
	Name                         string `ros:"name,omitempty"`
	Supported                    string `ros:"supported,omitempty"`
	VhtBasicMcs                  string `ros:"vht-basic-mcs,omitempty"`
	VhtSupportedMcs              string `ros:"vht-supported-mcs,omitempty"`
}
```

```go
// CapsmanRatesSet adalah argumen `/caps-man/rates/set`.
type CapsmanRatesSet struct {
	Basic                        string `ros:"basic,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	HtBasicMcs                   string `ros:"ht-basic-mcs,omitempty"`
	HtSupportedMcs               string `ros:"ht-supported-mcs,omitempty"`
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	Supported                    string `ros:"supported,omitempty"`
	VhtBasicMcs                  string `ros:"vht-basic-mcs,omitempty"`
	VhtSupportedMcs              string `ros:"vht-supported-mcs,omitempty"`
}
```

### `/caps-man/registration-table`

- File: `caps-man_registration-table.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command or directory (caps-man)

Command: `edit`(2), `find`(1), `get`(4), `print`(23), `remove`(1), `reset`(2)

_Tidak ada struct: menu ini bukan tabel yang bisa di-print atau tidak tersedia di router uji (mis. paket tidak terpasang)._

### `/caps-man/remote-cap`

- File: `caps-man_remote-cap.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command or directory (caps-man)

Command: `find`(1), `get`(4), `print`(21), `provision`(1), `remove`(1), `set-identity`(2), `upgrade`(1)

_Tidak ada struct: menu ini bukan tabel yang bisa di-print atau tidak tersedia di router uji (mis. paket tidak terpasang)._

### `/caps-man/security`

- File: `caps-man_security.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command or directory (caps-man)

Command: `add`(13), `comment`(2), `edit`(2), `export`(8), `find`(1), `get`(4), `print`(22), `remove`(1), `reset`(13), `set`(13), `unset`(2)

```go
// CapsmanSecurityAdd adalah argumen `/caps-man/security/add`.
type CapsmanSecurityAdd struct {
	AuthenticationTypes          string `ros:"authentication-types,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	DisablePmkid                 string `ros:"disable-pmkid,omitempty"`
	EapMethods                   string `ros:"eap-methods,omitempty"`
	EapRadiusAccounting          string `ros:"eap-radius-accounting,omitempty"`
	Encryption                   string `ros:"encryption,omitempty"`
	GroupEncryption              string `ros:"group-encryption,omitempty"`
	GroupKeyUpdate               string `ros:"group-key-update,omitempty"`
	Name                         string `ros:"name,omitempty"`
	Passphrase                   string `ros:"passphrase,omitempty"`
	TLSCertificate               string `ros:"tls-certificate,omitempty"`
	TLSMode                      string `ros:"tls-mode,omitempty"`
}
```

```go
// CapsmanSecuritySet adalah argumen `/caps-man/security/set`.
type CapsmanSecuritySet struct {
	AuthenticationTypes          string `ros:"authentication-types,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	DisablePmkid                 string `ros:"disable-pmkid,omitempty"`
	EapMethods                   string `ros:"eap-methods,omitempty"`
	EapRadiusAccounting          string `ros:"eap-radius-accounting,omitempty"`
	Encryption                   string `ros:"encryption,omitempty"`
	GroupEncryption              string `ros:"group-encryption,omitempty"`
	GroupKeyUpdate               string `ros:"group-key-update,omitempty"`
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	Passphrase                   string `ros:"passphrase,omitempty"`
	TLSCertificate               string `ros:"tls-certificate,omitempty"`
	TLSMode                      string `ros:"tls-mode,omitempty"`
}
```

### `/certificate`

- File: `certificate__cmds.json`
- Sumber field: tidak ada kolom terdeteksi

Command: `acme-renew`(1), `add`(15), `add-acme`(6), `add-scep`(8), `card-reinstall`(2), `card-verify`(2), `create-certificate-request`(12), `edit`(2), `enable-ssl-certificate`(14), `export`(8), `export-certificate`(4), `find`(1), `get`(4), `import`(7), `issued-revoke`(1), `print`(24), `remove`(1), `reset`(15), `scep-renew`(1), `set`(15), `sign`(14), `sign-certificate-request`(5), `unset`(2)

```go
// CertificateAdd adalah argumen `/certificate/add`.
type CertificateAdd struct {
	CommonName                   string `ros:"common-name,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Country                      string `ros:"country,omitempty"`
	DaysValid                    string `ros:"days-valid,omitempty"`
	DigestAlgorithm              string `ros:"digest-algorithm,omitempty"`
	KeySize                      string `ros:"key-size,omitempty"`
	KeyUsage                     string `ros:"key-usage,omitempty"`
	Locality                     string `ros:"locality,omitempty"`
	Name                         string `ros:"name,omitempty"`
	Organization                 string `ros:"organization,omitempty"`
	State                        string `ros:"state,omitempty"`
	SubjectAltName               string `ros:"subject-alt-name,omitempty"`
	TrustStore                   string `ros:"trust-store,omitempty"`
	Trusted                      string `ros:"trusted,omitempty"`
	Unit                         string `ros:"unit,omitempty"`
}
```

```go
// CertificateSet adalah argumen `/certificate/set`.
type CertificateSet struct {
	CommonName                   string `ros:"common-name,omitempty"`
	Country                      string `ros:"country,omitempty"`
	DaysValid                    string `ros:"days-valid,omitempty"`
	DigestAlgorithm              string `ros:"digest-algorithm,omitempty"`
	KeySize                      string `ros:"key-size,omitempty"`
	KeyUsage                     string `ros:"key-usage,omitempty"`
	Locality                     string `ros:"locality,omitempty"`
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	Organization                 string `ros:"organization,omitempty"`
	State                        string `ros:"state,omitempty"`
	SubjectAltName               string `ros:"subject-alt-name,omitempty"`
	TrustStore                   string `ros:"trust-store,omitempty"`
	Trusted                      string `ros:"trusted,omitempty"`
	Unit                         string `ros:"unit,omitempty"`
}
```

### `/certificate/builtin`

- File: `certificate_builtin.json`
- Sumber field: live 17, completion 18

Command: 

```go
// CertificateBuiltinRow adalah kolom hasil `/certificate/builtin/print`.
type CertificateBuiltinRow struct {
	ID                           string `ros:".id,key"`
	Akid                         string `ros:"akid"`
	CommonName                   string `ros:"common-name"`
	Country                      string `ros:"country"`
	DaysValid                    string `ros:"days-valid"`
	InvalidAfter                 string `ros:"invalid-after"`
	InvalidBefore                string `ros:"invalid-before"`
	Issuer                       string `ros:"issuer"`
	KeyType                      string `ros:"key-type"`
	KeyUsage                     string `ros:"key-usage"`
	Locality                     string `ros:"locality"`
	Organization                 string `ros:"organization"`
	SerialNumber                 string `ros:"serial-number"`
	Skid                         string `ros:"skid"`
	State                        string `ros:"state"`
	SubjectAltName               string `ros:"subject-alt-name"`
	Unit                         string `ros:"unit"`
	About                        string `ros:"about"`
	KeySize                      string `ros:"key-size"`
}
```

### `/certificate/crl`

- File: `certificate_crl.json`
- Sumber field: completion 14

Command: 

```go
// CertificateCrlRow adalah kolom hasil `/certificate/crl/print`.
type CertificateCrlRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Akid                         string `ros:"akid"`
	Cert                         string `ros:"cert"`
	Dynamic                      string `ros:"dynamic"`
	Expired                      string `ros:"expired"`
	Fingerprint                  string `ros:"fingerprint"`
	Invalid                      string `ros:"invalid"`
	LastUpdate                   string `ros:"last-update"`
	NextUpdate                   string `ros:"next-update"`
	Num                          string `ros:"num"`
	Revoked                      string `ros:"revoked"`
	Signature                    string `ros:"signature"`
	TrustStore                   string `ros:"trust-store"`
	Url                          string `ros:"url"`
}
```

### `/certificate/scep-server`

- File: `certificate_scep-server.json`
- Sumber field: completion 7

Command: 

```go
// CertificateScepserverRow adalah kolom hasil `/certificate/scep-server/print`.
type CertificateScepserverRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	CaCert                       string `ros:"ca-cert"`
	DaysValid                    string `ros:"days-valid"`
	Disabled                     string `ros:"disabled"`
	NextCaCert                   string `ros:"next-ca-cert"`
	Path                         string `ros:"path"`
	RequestLifetime              string `ros:"request-lifetime"`
}
```

### `/certificate/settings`

- File: `certificate_settings.json`
- Sumber field: live 5

Command: 

```go
// CertificateSettingsRow adalah kolom hasil `/certificate/settings/print`.
type CertificateSettingsRow struct {
	BuiltinTrustStore            string `ros:"builtin-trust-store"`
	CrlDownload                  string `ros:"crl-download"`
	CrlStore                     string `ros:"crl-store"`
	CrlUse                       string `ros:"crl-use"`
	CurrentDefaults              string `ros:"current-defaults"`
}
```

### `/console/settings`

- File: `console_settings.json`
- Sumber field: live 3

Command: 

```go
// ConsoleSettingsRow adalah kolom hasil `/console/settings/print`.
type ConsoleSettingsRow struct {
	LogScriptErrors              string `ros:"log-script-errors"`
	SanitizeNames                string `ros:"sanitize-names"`
	TabWidth                     string `ros:"tab-width"`
}
```

### `/disk/settings`

- File: `disk_settings.json`
- Sumber field: live 5

Command: 

```go
// DiskSettingsRow adalah kolom hasil `/disk/settings/print`.
type DiskSettingsRow struct {
	AutoMediaInterface           string `ros:"auto-media-interface"`
	AutoMediaSharing             string `ros:"auto-media-sharing"`
	AutoSmbSharing               string `ros:"auto-smb-sharing"`
	AutoSmbUser                  string `ros:"auto-smb-user"`
	DefaultMountPointTemplate    string `ros:"default-mount-point-template"`
}
```

### `/file`

- File: `(dari katalog Go)`
- Katalog: `file`
- Sumber field: katalog 4

Command: 

```go
// FileRow adalah kolom hasil `/file/print`.
type FileRow struct {
	ID                           string `ros:".id,key"`
	Name                         string `ros:"name"`
	Type                         string `ros:"type"`
	Size                         string `ros:"size"`
	CreationTime                 string `ros:"creation-time"`
}
```

### `/interface`

- File: `interface__cmds.json`
- Katalog: `interface`
- Sumber field: katalog 21

Command: `blink`(1), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `monitor-traffic`(11), `print`(24), `recursive-print`(0), `reset`(6), `reset-counters`(1), `set`(6)

```go
// InterfaceRow adalah kolom hasil `/interface/print`.
type InterfaceRow struct {
	ID                           string `ros:".id,key"`
	// string value
	Name                         string `ros:"name"`
	Type                         string `ros:"type"`
	MTU                          string `ros:"mtu"`
	ActualMTU                    string `ros:"actual-mtu"`
	Running                      string `ros:"running"`
	Disabled                     string `ros:"disabled"`
	// string value
	Comment                      string `ros:"comment"`
	LastLinkUpTime               string `ros:"last-link-up-time"`
	LastLinkDownTime             string `ros:"last-link-down-time"`
	LinkDowns                    string `ros:"link-downs"`
	RXByte                       string `ros:"rx-byte"`
	TXByte                       string `ros:"tx-byte"`
	RXPacket                     string `ros:"rx-packet"`
	TXPacket                     string `ros:"tx-packet"`
	RXDrop                       string `ros:"rx-drop"`
	TXDrop                       string `ros:"tx-drop"`
	TXQueueDrop                  string `ros:"tx-queue-drop"`
	RXError                      string `ros:"rx-error"`
	TXError                      string `ros:"tx-error"`
}
```

```go
// InterfaceSet adalah argumen `/interface/set`.
type InterfaceSet struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	L2MTU                        string `ros:"l2mtu,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
}
```

### `/interface/6to4`

- File: `interface_6to4.json`
- Sumber field: completion 15

Command: 

```go
// Interface6to4Row adalah kolom hasil `/interface/6to4/print`.
type Interface6to4Row struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	ActualMTU                    string `ros:"actual-mtu"`
	ClampTcpMss                  string `ros:"clamp-tcp-mss"`
	Comment                      string `ros:"comment"`
	CurrentRemoteAddress         string `ros:"current-remote-address"`
	Disabled                     string `ros:"disabled"`
	DontFragment                 string `ros:"dont-fragment"`
	Dscp                         string `ros:"dscp"`
	IPSecSecret                  string `ros:"ipsec-secret"`
	Keepalive                    string `ros:"keepalive"`
	LocalAddress                 string `ros:"local-address"`
	MTU                          string `ros:"mtu"`
	Name                         string `ros:"name"`
	RemoteAddress                string `ros:"remote-address"`
	Running                      string `ros:"running"`
}
```

### `/interface/bonding`

- File: `interface_bonding.json`
- Sumber field: completion 27

Command: 

```go
// InterfaceBondingRow adalah kolom hasil `/interface/bonding/print`.
type InterfaceBondingRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	ARP                          string `ros:"arp"`
	ARPInterval                  string `ros:"arp-interval"`
	ARPIPTargets                 string `ros:"arp-ip-targets"`
	ARPTimeout                   string `ros:"arp-timeout"`
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	DownDelay                    string `ros:"down-delay"`
	ForcedMACAddress             string `ros:"forced-mac-address"`
	LacpMode                     string `ros:"lacp-mode"`
	LacpRate                     string `ros:"lacp-rate"`
	LacpSystemID                 string `ros:"lacp-system-id"`
	LacpSystemPriority           string `ros:"lacp-system-priority"`
	LacpUserKey                  string `ros:"lacp-user-key"`
	LinkMonitoring               string `ros:"link-monitoring"`
	MACAddress                   string `ros:"mac-address"`
	MiiInterval                  string `ros:"mii-interval"`
	MinLinks                     string `ros:"min-links"`
	MlagID                       string `ros:"mlag-id"`
	Mode                         string `ros:"mode"`
	MTU                          string `ros:"mtu"`
	Name                         string `ros:"name"`
	Primary                      string `ros:"primary"`
	Running                      string `ros:"running"`
	Slaves                       string `ros:"slaves"`
	TransmitHashPolicy           string `ros:"transmit-hash-policy"`
	UpDelay                      string `ros:"up-delay"`
}
```

### `/interface/bridge`

- File: `interface_bridge.json`
- Katalog: `interface-bridge`
- Sumber field: live 33, completion 58, katalog 7

Command: `add`(52), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `monitor`(11), `print`(21), `remove`(1), `reset`(52), `set`(52), `unset`(2)

```go
// InterfaceBridgeRow adalah kolom hasil `/interface/bridge/print`.
type InterfaceBridgeRow struct {
	ID                           string `ros:".id,key"`
	ActualMTU                    string `ros:"actual-mtu"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	AdminMAC                     string `ros:"admin-mac"`
	// 10s..1w4d13h46m40s (time interval)
	AgeingTime                   string `ros:"ageing-time"`
	ARP                          string `ros:"arp"`
	ARPTimeout                   string `ros:"arp-timeout"`
	AutoMAC                      string `ros:"auto-mac"`
	// string value
	Comment                      string `ros:"comment"`
	DHCPSnooping                 string `ros:"dhcp-snooping"`
	DHCPv6Snooping               string `ros:"dhcpv6-snooping"`
	Disabled                     string `ros:"disabled"`
	Dynamic                      string `ros:"dynamic"`
	FastForward                  string `ros:"fast-forward"`
	// 4s..30s (time interval)
	ForwardDelay                 string `ros:"forward-delay"`
	ForwardReservedAddresses     string `ros:"forward-reserved-addresses"`
	IGMPSnooping                 string `ros:"igmp-snooping"`
	L2MTU                        string `ros:"l2mtu"`
	MACAddress                   string `ros:"mac-address"`
	Managed                      string `ros:"managed"`
	MaxLearnedEntries            string `ros:"max-learned-entries"`
	// 6s..40s (time interval)
	MaxMessageAge                string `ros:"max-message-age"`
	MlagHeartbeat                string `ros:"mlag-heartbeat"`
	MlagPeerPort                 string `ros:"mlag-peer-port"`
	MlagPriority                 string `ros:"mlag-priority"`
	MTU                          string `ros:"mtu"`
	// string value
	Name                         string `ros:"name"`
	PortCostMode                 string `ros:"port-cost-mode"`
	Priority                     string `ros:"priority"`
	ProtocolMode                 string `ros:"protocol-mode"`
	RaGuard                      string `ros:"ra-guard"`
	Running                      string `ros:"running"`
	TransmitHoldCount            string `ros:"transmit-hold-count"`
	VLANFiltering                string `ros:"vlan-filtering"`
	About                        string `ros:"about"`
	// string value
	DHCPAgentCircuitID           string `ros:"dhcp-agent-circuit-id"`
	// string value
	DHCPAgentRemoteID            string `ros:"dhcp-agent-remote-id"`
	// string value
	DHCPv6AgentCircuitID         string `ros:"dhcpv6-agent-circuit-id"`
	// string value
	DHCPv6AgentRemoteID          string `ros:"dhcpv6-agent-remote-id"`
	EtherType                    string `ros:"ether-type"`
	FrameTypes                   string `ros:"frame-types"`
	IGMPVersion                  string `ros:"igmp-version"`
	IngressFiltering             string `ros:"ingress-filtering"`
	// 1s..71w2h27m52s950ms (time interval)
	LastMemberInterval           string `ros:"last-member-interval"`
	LastMemberQueryCount         string `ros:"last-member-query-count"`
	MaxHops                      string `ros:"max-hops"`
	// 1s..71w2h27m52s950ms (time interval)
	MembershipInterval           string `ros:"membership-interval"`
	MldVersion                   string `ros:"mld-version"`
	MulticastQuerier             string `ros:"multicast-querier"`
	MulticastRouter              string `ros:"multicast-router"`
	Mvrp                         string `ros:"mvrp"`
	Pvid                         string `ros:"pvid"`
	// 1s..71w2h27m52s950ms (time interval)
	QuerierInterval              string `ros:"querier-interval"`
	QuerierUsesBridgeAddress     string `ros:"querier-uses-bridge-address"`
	// 1s..71w2h27m52s950ms (time interval)
	QueryInterval                string `ros:"query-interval"`
	// 1s..71w2h27m52s950ms (time interval)
	QueryResponseInterval        string `ros:"query-response-interval"`
	// string value
	RegionName                   string `ros:"region-name"`
	RegionRevision               string `ros:"region-revision"`
	StartupQueryCount            string `ros:"startup-query-count"`
	// 1s..71w2h27m52s950ms (time interval)
	StartupQueryInterval         string `ros:"startup-query-interval"`
}
```

```go
// InterfaceBridgeAdd adalah argumen `/interface/bridge/add`.
type InterfaceBridgeAdd struct {
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	AdminMAC                     string `ros:"admin-mac,omitempty"`
	// 10s..1w4d13h46m40s (time interval)
	AgeingTime                   string `ros:"ageing-time,omitempty"`
	ARP                          string `ros:"arp,omitempty"`
	ARPTimeout                   string `ros:"arp-timeout,omitempty"`
	AutoMAC                      string `ros:"auto-mac,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	// string value
	DHCPAgentCircuitID           string `ros:"dhcp-agent-circuit-id,omitempty"`
	// string value
	DHCPAgentRemoteID            string `ros:"dhcp-agent-remote-id,omitempty"`
	DHCPSnooping                 string `ros:"dhcp-snooping,omitempty"`
	// string value
	DHCPv6AgentCircuitID         string `ros:"dhcpv6-agent-circuit-id,omitempty"`
	// string value
	DHCPv6AgentRemoteID          string `ros:"dhcpv6-agent-remote-id,omitempty"`
	DHCPv6Snooping               string `ros:"dhcpv6-snooping,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	EtherType                    string `ros:"ether-type,omitempty"`
	FastForward                  string `ros:"fast-forward,omitempty"`
	// 4s..30s (time interval)
	ForwardDelay                 string `ros:"forward-delay,omitempty"`
	ForwardReservedAddresses     string `ros:"forward-reserved-addresses,omitempty"`
	FrameTypes                   string `ros:"frame-types,omitempty"`
	IGMPSnooping                 string `ros:"igmp-snooping,omitempty"`
	IGMPVersion                  string `ros:"igmp-version,omitempty"`
	IngressFiltering             string `ros:"ingress-filtering,omitempty"`
	// 1s..71w2h27m52s950ms (time interval)
	LastMemberInterval           string `ros:"last-member-interval,omitempty"`
	LastMemberQueryCount         string `ros:"last-member-query-count,omitempty"`
	MaxHops                      string `ros:"max-hops,omitempty"`
	MaxLearnedEntries            string `ros:"max-learned-entries,omitempty"`
	// 6s..40s (time interval)
	MaxMessageAge                string `ros:"max-message-age,omitempty"`
	// 1s..71w2h27m52s950ms (time interval)
	MembershipInterval           string `ros:"membership-interval,omitempty"`
	MlagHeartbeat                string `ros:"mlag-heartbeat,omitempty"`
	MlagPeerPort                 string `ros:"mlag-peer-port,omitempty"`
	MlagPriority                 string `ros:"mlag-priority,omitempty"`
	MldVersion                   string `ros:"mld-version,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	MulticastQuerier             string `ros:"multicast-querier,omitempty"`
	MulticastRouter              string `ros:"multicast-router,omitempty"`
	Mvrp                         string `ros:"mvrp,omitempty"`
	// string value
	Name                         string `ros:"name,required,omitempty"`
	PortCostMode                 string `ros:"port-cost-mode,omitempty"`
	Priority                     string `ros:"priority,omitempty"`
	ProtocolMode                 string `ros:"protocol-mode,omitempty"`
	Pvid                         string `ros:"pvid,omitempty"`
	// 1s..71w2h27m52s950ms (time interval)
	QuerierInterval              string `ros:"querier-interval,omitempty"`
	QuerierUsesBridgeAddress     string `ros:"querier-uses-bridge-address,omitempty"`
	// 1s..71w2h27m52s950ms (time interval)
	QueryInterval                string `ros:"query-interval,omitempty"`
	// 1s..71w2h27m52s950ms (time interval)
	QueryResponseInterval        string `ros:"query-response-interval,omitempty"`
	RaGuard                      string `ros:"ra-guard,omitempty"`
	// string value
	RegionName                   string `ros:"region-name,omitempty"`
	RegionRevision               string `ros:"region-revision,omitempty"`
	StartupQueryCount            string `ros:"startup-query-count,omitempty"`
	// 1s..71w2h27m52s950ms (time interval)
	StartupQueryInterval         string `ros:"startup-query-interval,omitempty"`
	TransmitHoldCount            string `ros:"transmit-hold-count,omitempty"`
	VLANFiltering                string `ros:"vlan-filtering,omitempty"`
}
```

```go
// InterfaceBridgeSet adalah argumen `/interface/bridge/set`.
type InterfaceBridgeSet struct {
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	AdminMAC                     string `ros:"admin-mac,omitempty"`
	// 10s..1w4d13h46m40s (time interval)
	AgeingTime                   string `ros:"ageing-time,omitempty"`
	ARP                          string `ros:"arp,omitempty"`
	ARPTimeout                   string `ros:"arp-timeout,omitempty"`
	AutoMAC                      string `ros:"auto-mac,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	// string value
	DHCPAgentCircuitID           string `ros:"dhcp-agent-circuit-id,omitempty"`
	// string value
	DHCPAgentRemoteID            string `ros:"dhcp-agent-remote-id,omitempty"`
	DHCPSnooping                 string `ros:"dhcp-snooping,omitempty"`
	// string value
	DHCPv6AgentCircuitID         string `ros:"dhcpv6-agent-circuit-id,omitempty"`
	// string value
	DHCPv6AgentRemoteID          string `ros:"dhcpv6-agent-remote-id,omitempty"`
	DHCPv6Snooping               string `ros:"dhcpv6-snooping,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	EtherType                    string `ros:"ether-type,omitempty"`
	FastForward                  string `ros:"fast-forward,omitempty"`
	// 4s..30s (time interval)
	ForwardDelay                 string `ros:"forward-delay,omitempty"`
	ForwardReservedAddresses     string `ros:"forward-reserved-addresses,omitempty"`
	FrameTypes                   string `ros:"frame-types,omitempty"`
	IGMPSnooping                 string `ros:"igmp-snooping,omitempty"`
	IGMPVersion                  string `ros:"igmp-version,omitempty"`
	IngressFiltering             string `ros:"ingress-filtering,omitempty"`
	// 1s..71w2h27m52s950ms (time interval)
	LastMemberInterval           string `ros:"last-member-interval,omitempty"`
	LastMemberQueryCount         string `ros:"last-member-query-count,omitempty"`
	MaxHops                      string `ros:"max-hops,omitempty"`
	MaxLearnedEntries            string `ros:"max-learned-entries,omitempty"`
	// 6s..40s (time interval)
	MaxMessageAge                string `ros:"max-message-age,omitempty"`
	// 1s..71w2h27m52s950ms (time interval)
	MembershipInterval           string `ros:"membership-interval,omitempty"`
	MlagHeartbeat                string `ros:"mlag-heartbeat,omitempty"`
	MlagPeerPort                 string `ros:"mlag-peer-port,omitempty"`
	MlagPriority                 string `ros:"mlag-priority,omitempty"`
	MldVersion                   string `ros:"mld-version,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	MulticastQuerier             string `ros:"multicast-querier,omitempty"`
	MulticastRouter              string `ros:"multicast-router,omitempty"`
	Mvrp                         string `ros:"mvrp,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	PortCostMode                 string `ros:"port-cost-mode,omitempty"`
	Priority                     string `ros:"priority,omitempty"`
	ProtocolMode                 string `ros:"protocol-mode,omitempty"`
	Pvid                         string `ros:"pvid,omitempty"`
	// 1s..71w2h27m52s950ms (time interval)
	QuerierInterval              string `ros:"querier-interval,omitempty"`
	QuerierUsesBridgeAddress     string `ros:"querier-uses-bridge-address,omitempty"`
	// 1s..71w2h27m52s950ms (time interval)
	QueryInterval                string `ros:"query-interval,omitempty"`
	// 1s..71w2h27m52s950ms (time interval)
	QueryResponseInterval        string `ros:"query-response-interval,omitempty"`
	RaGuard                      string `ros:"ra-guard,omitempty"`
	// string value
	RegionName                   string `ros:"region-name,omitempty"`
	RegionRevision               string `ros:"region-revision,omitempty"`
	StartupQueryCount            string `ros:"startup-query-count,omitempty"`
	// 1s..71w2h27m52s950ms (time interval)
	StartupQueryInterval         string `ros:"startup-query-interval,omitempty"`
	TransmitHoldCount            string `ros:"transmit-hold-count,omitempty"`
	VLANFiltering                string `ros:"vlan-filtering,omitempty"`
}
```

### `/interface/bridge/host`

- File: `(dari katalog Go)`
- Katalog: `interface-bridge-host`
- Sumber field: katalog 6

Command: 

```go
// InterfaceBridgeHostRow adalah kolom hasil `/interface/bridge/host/print`.
type InterfaceBridgeHostRow struct {
	ID                           string `ros:".id,key"`
	MACAddress                   string `ros:"mac-address"`
	OnInterface                  string `ros:"on-interface"`
	Bridge                       string `ros:"bridge"`
	Vid                          string `ros:"vid"`
	Age                          string `ros:"age"`
	Local                        string `ros:"local"`
}
```

### `/interface/bridge/port`

- File: `(dari katalog Go)`
- Katalog: `interface-bridge-port`
- Sumber field: katalog 5

Command: `add`(4), `set`(4)

```go
// InterfaceBridgePortRow adalah kolom hasil `/interface/bridge/port/print`.
type InterfaceBridgePortRow struct {
	ID                           string `ros:".id,key"`
	Interface                    string `ros:"interface"`
	Bridge                       string `ros:"bridge"`
	Pvid                         string `ros:"pvid"`
	Disabled                     string `ros:"disabled"`
	Comment                      string `ros:"comment"`
}
```

```go
// InterfaceBridgePortAdd adalah argumen `/interface/bridge/port/add`.
type InterfaceBridgePortAdd struct {
	Interface                    string `ros:"interface,required,omitempty"`
	Bridge                       string `ros:"bridge,required,omitempty"`
	Pvid                         string `ros:"pvid,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

```go
// InterfaceBridgePortSet adalah argumen `/interface/bridge/port/set`.
type InterfaceBridgePortSet struct {
	Interface                    string `ros:"interface,omitempty"`
	Bridge                       string `ros:"bridge,omitempty"`
	Pvid                         string `ros:"pvid,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

### `/interface/detect-internet`

- File: `interface_detect-internet.json`
- Sumber field: live 5

Command: `edit`(1), `export`(8), `get`(3), `print`(5), `reset`(5), `set`(5)

```go
// InterfaceDetectinternetRow adalah kolom hasil `/interface/detect-internet/print`.
type InterfaceDetectinternetRow struct {
	ID                           string `ros:".id,key"`
	DetectInterfaceList          string `ros:"detect-interface-list"`
	InternetInterfaceList        string `ros:"internet-interface-list"`
	LANInterfaceList             string `ros:"lan-interface-list"`
	// 1m..1d (time interval)
	RequestInterval              string `ros:"request-interval"`
	WANInterfaceList             string `ros:"wan-interface-list"`
}
```

```go
// InterfaceDetectinternetSet adalah argumen `/interface/detect-internet/set`.
type InterfaceDetectinternetSet struct {
	DetectInterfaceList          string `ros:"detect-interface-list,omitempty"`
	InternetInterfaceList        string `ros:"internet-interface-list,omitempty"`
	LANInterfaceList             string `ros:"lan-interface-list,omitempty"`
	// 1m..1d (time interval)
	RequestInterval              string `ros:"request-interval,omitempty"`
	WANInterfaceList             string `ros:"wan-interface-list,omitempty"`
}
```

### `/interface/eoip`

- File: `interface_eoip.json`
- Sumber field: completion 25

Command: 

```go
// InterfaceEOIPRow adalah kolom hasil `/interface/eoip/print`.
type InterfaceEOIPRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	ActualMTU                    string `ros:"actual-mtu"`
	AllowFastPath                string `ros:"allow-fast-path"`
	ARP                          string `ros:"arp"`
	ARPTimeout                   string `ros:"arp-timeout"`
	ClampTcpMss                  string `ros:"clamp-tcp-mss"`
	Comment                      string `ros:"comment"`
	CurrentRemoteAddress         string `ros:"current-remote-address"`
	Disabled                     string `ros:"disabled"`
	DontFragment                 string `ros:"dont-fragment"`
	Dscp                         string `ros:"dscp"`
	IPSecSecret                  string `ros:"ipsec-secret"`
	Keepalive                    string `ros:"keepalive"`
	L2MTU                        string `ros:"l2mtu"`
	LocalAddress                 string `ros:"local-address"`
	LoopProtect                  string `ros:"loop-protect"`
	LoopProtectDisableTime       string `ros:"loop-protect-disable-time"`
	LoopProtectSendInterval      string `ros:"loop-protect-send-interval"`
	LoopProtectStatus            string `ros:"loop-protect-status"`
	MACAddress                   string `ros:"mac-address"`
	MTU                          string `ros:"mtu"`
	Name                         string `ros:"name"`
	RemoteAddress                string `ros:"remote-address"`
	Running                      string `ros:"running"`
	TunnelID                     string `ros:"tunnel-id"`
}
```

### `/interface/eoipv6`

- File: `interface_eoipv6.json`
- Sumber field: completion 24

Command: 

```go
// InterfaceEoipv6Row adalah kolom hasil `/interface/eoipv6/print`.
type InterfaceEoipv6Row struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	ActualMTU                    string `ros:"actual-mtu"`
	ARP                          string `ros:"arp"`
	ARPTimeout                   string `ros:"arp-timeout"`
	ClampTcpMss                  string `ros:"clamp-tcp-mss"`
	Comment                      string `ros:"comment"`
	CurrentRemoteAddress         string `ros:"current-remote-address"`
	Disabled                     string `ros:"disabled"`
	DontFragment                 string `ros:"dont-fragment"`
	Dscp                         string `ros:"dscp"`
	IPSecSecret                  string `ros:"ipsec-secret"`
	Keepalive                    string `ros:"keepalive"`
	L2MTU                        string `ros:"l2mtu"`
	LocalAddress                 string `ros:"local-address"`
	LoopProtect                  string `ros:"loop-protect"`
	LoopProtectDisableTime       string `ros:"loop-protect-disable-time"`
	LoopProtectSendInterval      string `ros:"loop-protect-send-interval"`
	LoopProtectStatus            string `ros:"loop-protect-status"`
	MACAddress                   string `ros:"mac-address"`
	MTU                          string `ros:"mtu"`
	Name                         string `ros:"name"`
	RemoteAddress                string `ros:"remote-address"`
	Running                      string `ros:"running"`
	TunnelID                     string `ros:"tunnel-id"`
}
```

### `/interface/ethernet`

- File: `interface_ethernet.json`
- Katalog: `interface-ethernet`
- Sumber field: live 19, completion 154, katalog 26

Command: `blink`(1), `cable-test`(11), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `monitor`(11), `print`(23), `reset`(26), `reset-counters`(1), `reset-mac-address`(1), `set`(26), `unset`(2)

```go
// InterfaceEthernetRow adalah kolom hasil `/interface/ethernet/print`.
type InterfaceEthernetRow struct {
	ID                           string `ros:".id,key"`
	Advertise                    string `ros:"advertise"`
	ARP                          string `ros:"arp"`
	ARPTimeout                   string `ros:"arp-timeout"`
	AutoNegotiation              string `ros:"auto-negotiation"`
	CableSettings                string `ros:"cable-settings"`
	DefaultName                  string `ros:"default-name"`
	DisableRunningCheck          string `ros:"disable-running-check"`
	Disabled                     string `ros:"disabled"`
	LoopProtect                  string `ros:"loop-protect"`
	// time interval
	LoopProtectDisableTime       string `ros:"loop-protect-disable-time"`
	// 1s..7101w3d6h28m15s (time interval)
	LoopProtectSendInterval      string `ros:"loop-protect-send-interval"`
	LoopProtectStatus            string `ros:"loop-protect-status"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	MACAddress                   string `ros:"mac-address"`
	MTU                          string `ros:"mtu"`
	// string value
	Name                         string `ros:"name"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	OrigMACAddress               string `ros:"orig-mac-address"`
	Running                      string `ros:"running"`
	Slave                        string `ros:"slave"`
	About                        string `ros:"about"`
	ComboMode                    string `ros:"combo-mode"`
	// string value
	Comment                      string `ros:"comment"`
	CurrentLearned               string `ros:"current-learned"`
	CustomDropPacket             string `ros:"custom-drop-packet"`
	DriverRXByte                 string `ros:"driver-rx-byte"`
	DriverRXPacket               string `ros:"driver-rx-packet"`
	DriverTXByte                 string `ros:"driver-tx-byte"`
	DriverTXPacket               string `ros:"driver-tx-packet"`
	FcFecBlockCorrected          string `ros:"fc-fec-block-corrected"`
	FcFecBlockUncorrected        string `ros:"fc-fec-block-uncorrected"`
	FcFecRXBlock                 string `ros:"fc-fec-rx-block"`
	FecMode                      string `ros:"fec-mode"`
	L2MTU                        string `ros:"l2mtu"`
	MdixEnable                   string `ros:"mdix-enable"`
	NotLearned                   string `ros:"not-learned"`
	PolicyDropPacket             string `ros:"policy-drop-packet"`
	RsFecCodewords               string `ros:"rs-fec-codewords"`
	RsFecCorrected               string `ros:"rs-fec-corrected"`
	RsFecSymbolError             string `ros:"rs-fec-symbol-error"`
	RsFecUncorrected             string `ros:"rs-fec-uncorrected"`
	RXX1024X1518                 string `ros:"rx-1024-1518"`
	RXX1024Max                   string `ros:"rx-1024-max"`
	RXX128X255                   string `ros:"rx-128-255"`
	RXX1519Max                   string `ros:"rx-1519-max"`
	RXX256X511                   string `ros:"rx-256-511"`
	RXX512X1023                  string `ros:"rx-512-1023"`
	RXX64                        string `ros:"rx-64"`
	RXX65X127                    string `ros:"rx-65-127"`
	RXAlignError                 string `ros:"rx-align-error"`
	RXBroadcast                  string `ros:"rx-broadcast"`
	RXBytes                      string `ros:"rx-bytes"`
	RXCarrierError               string `ros:"rx-carrier-error"`
	RXCodeError                  string `ros:"rx-code-error"`
	RXControl                    string `ros:"rx-control"`
	RXDrop                       string `ros:"rx-drop"`
	RXErrorEvents                string `ros:"rx-error-events"`
	RXFcsError                   string `ros:"rx-fcs-error"`
	RXFlowControl                string `ros:"rx-flow-control"`
	RXFragment                   string `ros:"rx-fragment"`
	RXIPHeaderChecksumError      string `ros:"rx-ip-header-checksum-error"`
	RXJabber                     string `ros:"rx-jabber"`
	RXLengthError                string `ros:"rx-length-error"`
	RXMulticast                  string `ros:"rx-multicast"`
	RXOverflow                   string `ros:"rx-overflow"`
	RXPacket                     string `ros:"rx-packet"`
	RXPause                      string `ros:"rx-pause"`
	RXTcpChecksumError           string `ros:"rx-tcp-checksum-error"`
	RXTooLong                    string `ros:"rx-too-long"`
	RXTooShort                   string `ros:"rx-too-short"`
	RXUdpChecksumError           string `ros:"rx-udp-checksum-error"`
	RXUnicast                    string `ros:"rx-unicast"`
	RXUnknownOp                  string `ros:"rx-unknown-op"`
	SfpIgnoreRXLos               string `ros:"sfp-ignore-rx-los"`
	SfpRateSelect                string `ros:"sfp-rate-select"`
	SfpShutdownTemperature       string `ros:"sfp-shutdown-temperature"`
	Speed                        string `ros:"speed"`
	TXX1024X1518                 string `ros:"tx-1024-1518"`
	TXX1024Max                   string `ros:"tx-1024-max"`
	TXX128X255                   string `ros:"tx-128-255"`
	TXX1519Max                   string `ros:"tx-1519-max"`
	TXX256X511                   string `ros:"tx-256-511"`
	TXX512X1023                  string `ros:"tx-512-1023"`
	TXX64                        string `ros:"tx-64"`
	TXX65X127                    string `ros:"tx-65-127"`
	TXBroadcast                  string `ros:"tx-broadcast"`
	TXBytes                      string `ros:"tx-bytes"`
	TXCarrierSenseError          string `ros:"tx-carrier-sense-error"`
	TXCollision                  string `ros:"tx-collision"`
	TXControl                    string `ros:"tx-control"`
	TXDeferred                   string `ros:"tx-deferred"`
	TXDrop                       string `ros:"tx-drop"`
	TXDropByte                   string `ros:"tx-drop-byte"`
	TXDropPacket                 string `ros:"tx-drop-packet"`
	TXDropQueue0Byte             string `ros:"tx-drop-queue0-byte"`
	TXDropQueue0Packet           string `ros:"tx-drop-queue0-packet"`
	TXDropQueue1Byte             string `ros:"tx-drop-queue1-byte"`
	TXDropQueue1Packet           string `ros:"tx-drop-queue1-packet"`
	TXDropQueue2Byte             string `ros:"tx-drop-queue2-byte"`
	TXDropQueue2Packet           string `ros:"tx-drop-queue2-packet"`
	TXDropQueue3Byte             string `ros:"tx-drop-queue3-byte"`
	TXDropQueue3Packet           string `ros:"tx-drop-queue3-packet"`
	TXDropQueue4Byte             string `ros:"tx-drop-queue4-byte"`
	TXDropQueue4Packet           string `ros:"tx-drop-queue4-packet"`
	TXDropQueue5Byte             string `ros:"tx-drop-queue5-byte"`
	TXDropQueue5Packet           string `ros:"tx-drop-queue5-packet"`
	TXDropQueue6Byte             string `ros:"tx-drop-queue6-byte"`
	TXDropQueue6Packet           string `ros:"tx-drop-queue6-packet"`
	TXDropQueue7Byte             string `ros:"tx-drop-queue7-byte"`
	TXDropQueue7Packet           string `ros:"tx-drop-queue7-packet"`
	TXExcessiveCollision         string `ros:"tx-excessive-collision"`
	TXExcessiveDeferred          string `ros:"tx-excessive-deferred"`
	TXFcsError                   string `ros:"tx-fcs-error"`
	TXFlowControl                string `ros:"tx-flow-control"`
	TXFragment                   string `ros:"tx-fragment"`
	TXJabber                     string `ros:"tx-jabber"`
	TXLateCollision              string `ros:"tx-late-collision"`
	TXMulticast                  string `ros:"tx-multicast"`
	TXMultipleCollision          string `ros:"tx-multiple-collision"`
	TXPacket                     string `ros:"tx-packet"`
	TXPause                      string `ros:"tx-pause"`
	TXPauseHonored               string `ros:"tx-pause-honored"`
	TXQueueCustom0DropByte       string `ros:"tx-queue-custom0-drop-byte"`
	TXQueueCustom0DropPacket     string `ros:"tx-queue-custom0-drop-packet"`
	TXQueueCustom1DropByte       string `ros:"tx-queue-custom1-drop-byte"`
	TXQueueCustom1DropPacket     string `ros:"tx-queue-custom1-drop-packet"`
	TXQueue0Byte                 string `ros:"tx-queue0-byte"`
	TXQueue0Packet               string `ros:"tx-queue0-packet"`
	TXQueue1Byte                 string `ros:"tx-queue1-byte"`
	TXQueue1Packet               string `ros:"tx-queue1-packet"`
	TXQueue2Byte                 string `ros:"tx-queue2-byte"`
	TXQueue2Packet               string `ros:"tx-queue2-packet"`
	TXQueue3Byte                 string `ros:"tx-queue3-byte"`
	TXQueue3Packet               string `ros:"tx-queue3-packet"`
	TXQueue4Byte                 string `ros:"tx-queue4-byte"`
	TXQueue4Packet               string `ros:"tx-queue4-packet"`
	TXQueue5Byte                 string `ros:"tx-queue5-byte"`
	TXQueue5Packet               string `ros:"tx-queue5-packet"`
	TXQueue6Byte                 string `ros:"tx-queue6-byte"`
	TXQueue6Packet               string `ros:"tx-queue6-packet"`
	TXQueue7Byte                 string `ros:"tx-queue7-byte"`
	TXQueue7Packet               string `ros:"tx-queue7-packet"`
	TXRXX1024X1518               string `ros:"tx-rx-1024-1518"`
	TXRXX1024Max                 string `ros:"tx-rx-1024-max"`
	TXRXX128X255                 string `ros:"tx-rx-128-255"`
	TXRXX1519Max                 string `ros:"tx-rx-1519-max"`
	TXRXX256X511                 string `ros:"tx-rx-256-511"`
	TXRXX512X1023                string `ros:"tx-rx-512-1023"`
	TXRXX64                      string `ros:"tx-rx-64"`
	TXRXX65X127                  string `ros:"tx-rx-65-127"`
	TXSingleCollision            string `ros:"tx-single-collision"`
	TXTooLong                    string `ros:"tx-too-long"`
	TXTooShort                   string `ros:"tx-too-short"`
	TXTotalCollision             string `ros:"tx-total-collision"`
	TXUnderrun                   string `ros:"tx-underrun"`
	TXUnicast                    string `ros:"tx-unicast"`
	Status                       string `ros:"status"`
	Rate                         string `ros:"rate"`
	FullDuplex                   string `ros:"full-duplex"`
	SfpModulePresent             string `ros:"sfp-module-present"`
	SfpType                      string `ros:"sfp-type"`
	SfpVendorName                string `ros:"sfp-vendor-name"`
	SfpWavelength                string `ros:"sfp-wavelength"`
	SfpTemperature               string `ros:"sfp-temperature"`
	SfpSupplyVoltage             string `ros:"sfp-supply-voltage"`
	SfpTXBiasCurrent             string `ros:"sfp-tx-bias-current"`
	SfpTXPower                   string `ros:"sfp-tx-power"`
	SfpRXPower                   string `ros:"sfp-rx-power"`
}
```

```go
// InterfaceEthernetSet adalah argumen `/interface/ethernet/set`.
type InterfaceEthernetSet struct {
	Advertise                    string `ros:"advertise,omitempty"`
	ARP                          string `ros:"arp,omitempty"`
	ARPTimeout                   string `ros:"arp-timeout,omitempty"`
	AutoNegotiation              string `ros:"auto-negotiation,omitempty"`
	CableSettings                string `ros:"cable-settings,omitempty"`
	ComboMode                    string `ros:"combo-mode,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	DisableRunningCheck          string `ros:"disable-running-check,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	FecMode                      string `ros:"fec-mode,omitempty"`
	L2MTU                        string `ros:"l2mtu,omitempty"`
	LoopProtect                  string `ros:"loop-protect,omitempty"`
	// time interval
	LoopProtectDisableTime       string `ros:"loop-protect-disable-time,omitempty"`
	// 1s..7101w3d6h28m15s (time interval)
	LoopProtectSendInterval      string `ros:"loop-protect-send-interval,omitempty"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	MACAddress                   string `ros:"mac-address,omitempty"`
	MdixEnable                   string `ros:"mdix-enable,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	OrigMACAddress               string `ros:"orig-mac-address,omitempty"`
	RXFlowControl                string `ros:"rx-flow-control,omitempty"`
	SfpIgnoreRXLos               string `ros:"sfp-ignore-rx-los,omitempty"`
	SfpRateSelect                string `ros:"sfp-rate-select,omitempty"`
	SfpShutdownTemperature       string `ros:"sfp-shutdown-temperature,omitempty"`
	Speed                        string `ros:"speed,omitempty"`
	TXFlowControl                string `ros:"tx-flow-control,omitempty"`
}
```

### `/interface/gre`

- File: `interface_gre.json`
- Sumber field: completion 16

Command: 

```go
// InterfaceGRERow adalah kolom hasil `/interface/gre/print`.
type InterfaceGRERow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	ActualMTU                    string `ros:"actual-mtu"`
	AllowFastPath                string `ros:"allow-fast-path"`
	ClampTcpMss                  string `ros:"clamp-tcp-mss"`
	Comment                      string `ros:"comment"`
	CurrentRemoteAddress         string `ros:"current-remote-address"`
	Disabled                     string `ros:"disabled"`
	DontFragment                 string `ros:"dont-fragment"`
	Dscp                         string `ros:"dscp"`
	IPSecSecret                  string `ros:"ipsec-secret"`
	Keepalive                    string `ros:"keepalive"`
	LocalAddress                 string `ros:"local-address"`
	MTU                          string `ros:"mtu"`
	Name                         string `ros:"name"`
	RemoteAddress                string `ros:"remote-address"`
	Running                      string `ros:"running"`
}
```

### `/interface/gre6`

- File: `interface_gre6.json`
- Sumber field: completion 15

Command: 

```go
// InterfaceGre6Row adalah kolom hasil `/interface/gre6/print`.
type InterfaceGre6Row struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	ActualMTU                    string `ros:"actual-mtu"`
	ClampTcpMss                  string `ros:"clamp-tcp-mss"`
	Comment                      string `ros:"comment"`
	CurrentRemoteAddress         string `ros:"current-remote-address"`
	Disabled                     string `ros:"disabled"`
	DontFragment                 string `ros:"dont-fragment"`
	Dscp                         string `ros:"dscp"`
	IPSecSecret                  string `ros:"ipsec-secret"`
	Keepalive                    string `ros:"keepalive"`
	LocalAddress                 string `ros:"local-address"`
	MTU                          string `ros:"mtu"`
	Name                         string `ros:"name"`
	RemoteAddress                string `ros:"remote-address"`
	Running                      string `ros:"running"`
}
```

### `/interface/ipip`

- File: `interface_ipip.json`
- Sumber field: completion 17

Command: 

```go
// InterfaceIPIPRow adalah kolom hasil `/interface/ipip/print`.
type InterfaceIPIPRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	ActualMTU                    string `ros:"actual-mtu"`
	AllowFastPath                string `ros:"allow-fast-path"`
	ClampTcpMss                  string `ros:"clamp-tcp-mss"`
	Comment                      string `ros:"comment"`
	CurrentRemoteAddress         string `ros:"current-remote-address"`
	Disabled                     string `ros:"disabled"`
	DontFragment                 string `ros:"dont-fragment"`
	Dscp                         string `ros:"dscp"`
	Dynamic                      string `ros:"dynamic"`
	IPSecSecret                  string `ros:"ipsec-secret"`
	Keepalive                    string `ros:"keepalive"`
	LocalAddress                 string `ros:"local-address"`
	MTU                          string `ros:"mtu"`
	Name                         string `ros:"name"`
	RemoteAddress                string `ros:"remote-address"`
	Running                      string `ros:"running"`
}
```

### `/interface/ipipv6`

- File: `interface_ipipv6.json`
- Sumber field: completion 15

Command: 

```go
// InterfaceIpipv6Row adalah kolom hasil `/interface/ipipv6/print`.
type InterfaceIpipv6Row struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	ActualMTU                    string `ros:"actual-mtu"`
	ClampTcpMss                  string `ros:"clamp-tcp-mss"`
	Comment                      string `ros:"comment"`
	CurrentRemoteAddress         string `ros:"current-remote-address"`
	Disabled                     string `ros:"disabled"`
	DontFragment                 string `ros:"dont-fragment"`
	Dscp                         string `ros:"dscp"`
	IPSecSecret                  string `ros:"ipsec-secret"`
	Keepalive                    string `ros:"keepalive"`
	LocalAddress                 string `ros:"local-address"`
	MTU                          string `ros:"mtu"`
	Name                         string `ros:"name"`
	RemoteAddress                string `ros:"remote-address"`
	Running                      string `ros:"running"`
}
```

### `/interface/l2tp-client`

- File: `interface_l2tp-client.json`
- Sumber field: completion 27

Command: 

```go
// InterfaceL2tpclientRow adalah kolom hasil `/interface/l2tp-client/print`.
type InterfaceL2tpclientRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	AddDefaultRoute              string `ros:"add-default-route"`
	Allow                        string `ros:"allow"`
	AllowFastPath                string `ros:"allow-fast-path"`
	Comment                      string `ros:"comment"`
	ConnectTo                    string `ros:"connect-to"`
	DefaultRouteDistance         string `ros:"default-route-distance"`
	DialOnDemand                 string `ros:"dial-on-demand"`
	Disabled                     string `ros:"disabled"`
	IPSecSecret                  string `ros:"ipsec-secret"`
	KeepaliveTimeout             string `ros:"keepalive-timeout"`
	L2TPProtoVersion             string `ros:"l2tp-proto-version"`
	L2tpv3CircuitID              string `ros:"l2tpv3-circuit-id"`
	L2tpv3CookieLength           string `ros:"l2tpv3-cookie-length"`
	L2tpv3DigestHash             string `ros:"l2tpv3-digest-hash"`
	MaxMru                       string `ros:"max-mru"`
	MaxMTU                       string `ros:"max-mtu"`
	Mrru                         string `ros:"mrru"`
	Name                         string `ros:"name"`
	Password                     string `ros:"password"`
	Profile                      string `ros:"profile"`
	RandomSourcePort             string `ros:"random-source-port"`
	Running                      string `ros:"running"`
	SrcAddress                   string `ros:"src-address"`
	UseIPSec                     string `ros:"use-ipsec"`
	UsePeerDNS                   string `ros:"use-peer-dns"`
	User                         string `ros:"user"`
}
```

### `/interface/l2tp-ether`

- File: `interface_l2tp-ether.json`
- Sumber field: completion 27

Command: 

```go
// InterfaceL2tpetherRow adalah kolom hasil `/interface/l2tp-ether/print`.
type InterfaceL2tpetherRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	ActualMTU                    string `ros:"actual-mtu"`
	AllowFastPath                string `ros:"allow-fast-path"`
	CircuitID                    string `ros:"circuit-id"`
	Comment                      string `ros:"comment"`
	ConnectTo                    string `ros:"connect-to"`
	CookieLength                 string `ros:"cookie-length"`
	DigestHash                   string `ros:"digest-hash"`
	Disabled                     string `ros:"disabled"`
	Dynamic                      string `ros:"dynamic"`
	IPSecSecret                  string `ros:"ipsec-secret"`
	L2TPProtoVersion             string `ros:"l2tp-proto-version"`
	LocalAddress                 string `ros:"local-address"`
	LocalSessionID               string `ros:"local-session-id"`
	LocalTunnelID                string `ros:"local-tunnel-id"`
	MACAddress                   string `ros:"mac-address"`
	MTU                          string `ros:"mtu"`
	Name                         string `ros:"name"`
	PeerCookie                   string `ros:"peer-cookie"`
	RemoteSessionID              string `ros:"remote-session-id"`
	RemoteTunnelID               string `ros:"remote-tunnel-id"`
	Running                      string `ros:"running"`
	SendCookie                   string `ros:"send-cookie"`
	Unmanaged                    string `ros:"unmanaged"`
	UnmanagedMode                string `ros:"unmanaged-mode"`
	UseIPSec                     string `ros:"use-ipsec"`
	UseL2SpecificSublayer        string `ros:"use-l2-specific-sublayer"`
}
```

### `/interface/l2tp-server`

- File: `interface_l2tp-server.json`
- Sumber field: completion 12

Command: 

```go
// InterfaceL2tpserverRow adalah kolom hasil `/interface/l2tp-server/print`.
type InterfaceL2tpserverRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	ClientAddress                string `ros:"client-address"`
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	Dynamic                      string `ros:"dynamic"`
	Encoding                     string `ros:"encoding"`
	Mru                          string `ros:"mru"`
	MTU                          string `ros:"mtu"`
	Name                         string `ros:"name"`
	Running                      string `ros:"running"`
	Uptime                       string `ros:"uptime"`
	User                         string `ros:"user"`
}
```

### `/interface/list`

- File: `interface_list.json`
- Katalog: `interface-list`
- Sumber field: live 7, completion 7, katalog 2

Command: `add`(5), `comment`(2), `edit`(2), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(5), `set`(5)

```go
// InterfaceListRow adalah kolom hasil `/interface/list/print`.
type InterfaceListRow struct {
	ID                           string `ros:".id,key"`
	Builtin                      string `ros:"builtin"`
	// string value
	Comment                      string `ros:"comment"`
	Dynamic                      string `ros:"dynamic"`
	Exclude                      string `ros:"exclude"`
	Include                      string `ros:"include"`
	// string value
	Name                         string `ros:"name"`
	About                        string `ros:"about"`
}
```

```go
// InterfaceListAdd adalah argumen `/interface/list/add`.
type InterfaceListAdd struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Exclude                      string `ros:"exclude,omitempty"`
	Include                      string `ros:"include,omitempty"`
	// string value
	Name                         string `ros:"name,required,omitempty"`
}
```

```go
// InterfaceListSet adalah argumen `/interface/list/set`.
type InterfaceListSet struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Exclude                      string `ros:"exclude,omitempty"`
	Include                      string `ros:"include,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
}
```

### `/interface/list/member`

- File: `(dari katalog Go)`
- Katalog: `interface-list-member`
- Sumber field: katalog 4

Command: `add`(3), `set`(3)

```go
// InterfaceListMemberRow adalah kolom hasil `/interface/list/member/print`.
type InterfaceListMemberRow struct {
	ID                           string `ros:".id,key"`
	List                         string `ros:"list"`
	Interface                    string `ros:"interface"`
	Disabled                     string `ros:"disabled"`
	Comment                      string `ros:"comment"`
}
```

```go
// InterfaceListMemberAdd adalah argumen `/interface/list/member/add`.
type InterfaceListMemberAdd struct {
	List                         string `ros:"list,required,omitempty"`
	Interface                    string `ros:"interface,required,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

```go
// InterfaceListMemberSet adalah argumen `/interface/list/member/set`.
type InterfaceListMemberSet struct {
	List                         string `ros:"list,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

### `/interface/lte`

- File: `interface_lte.json`
- Sumber field: completion 20

Command: 

```go
// InterfaceLTERow adalah kolom hasil `/interface/lte/print`.
type InterfaceLTERow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	AdvertisedMTU                string `ros:"advertised-mtu"`
	AllowRoaming                 string `ros:"allow-roaming"`
	ApnProfiles                  string `ros:"apn-profiles"`
	Band                         string `ros:"band"`
	Comment                      string `ros:"comment"`
	DefaultName                  string `ros:"default-name"`
	Disabled                     string `ros:"disabled"`
	Inactive                     string `ros:"inactive"`
	Master                       string `ros:"master"`
	ModemInit                    string `ros:"modem-init"`
	MTU                          string `ros:"mtu"`
	Name                         string `ros:"name"`
	NetworkMode                  string `ros:"network-mode"`
	NrBand                       string `ros:"nr-band"`
	Operator                     string `ros:"operator"`
	Pin                          string `ros:"pin"`
	Running                      string `ros:"running"`
	SMSProtocol                  string `ros:"sms-protocol"`
	SMSRead                      string `ros:"sms-read"`
}
```

### `/interface/macsec`

- File: `interface_macsec.json`
- Sumber field: completion 42

Command: `add`(9), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `monitor`(11), `print`(23), `remove`(1), `reset`(9), `set`(9)

```go
// InterfaceMacsecRow adalah kolom hasil `/interface/macsec/print`.
type InterfaceMacsecRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	// hexadecimal string value, max length 16
	Cak                          string `ros:"cak"`
	// hexadecimal string value, max length 32
	Ckn                          string `ros:"ckn"`
	// string value
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	HWOffloaded                  string `ros:"hw-offloaded"`
	Inactive                     string `ros:"inactive"`
	Interface                    string `ros:"interface"`
	MTU                          string `ros:"mtu"`
	// string value
	Name                         string `ros:"name"`
	Profile                      string `ros:"profile"`
	Running                      string `ros:"running"`
	RXBadTag                     string `ros:"rx-bad-tag"`
	RXNoSci                      string `ros:"rx-no-sci"`
	RXNoTag                      string `ros:"rx-no-tag"`
	RXOverrun                    string `ros:"rx-overrun"`
	RXSaInvalid                  string `ros:"rx-sa-invalid"`
	RXSaNotUsingSa               string `ros:"rx-sa-not-using-sa"`
	RXSaNotValid                 string `ros:"rx-sa-not-valid"`
	RXSaOk                       string `ros:"rx-sa-ok"`
	RXSaUnusedSa                 string `ros:"rx-sa-unused-sa"`
	RXScDecryptedByte            string `ros:"rx-sc-decrypted-byte"`
	RXScDelayed                  string `ros:"rx-sc-delayed"`
	RXScInvalid                  string `ros:"rx-sc-invalid"`
	RXScLate                     string `ros:"rx-sc-late"`
	RXScNotUsingSa               string `ros:"rx-sc-not-using-sa"`
	RXScNotValid                 string `ros:"rx-sc-not-valid"`
	RXScOk                       string `ros:"rx-sc-ok"`
	RXScUnchecked                string `ros:"rx-sc-unchecked"`
	RXScUnusedSa                 string `ros:"rx-sc-unused-sa"`
	RXScValidatedByte            string `ros:"rx-sc-validated-byte"`
	RXUnknownSci                 string `ros:"rx-unknown-sci"`
	RXUntagged                   string `ros:"rx-untagged"`
	Status                       string `ros:"status"`
	TXSaEncrypted                string `ros:"tx-sa-encrypted"`
	TXSaProtected                string `ros:"tx-sa-protected"`
	TXScEncryptedByte            string `ros:"tx-sc-encrypted-byte"`
	TXScEncryptedPacket          string `ros:"tx-sc-encrypted-packet"`
	TXScProtectedByte            string `ros:"tx-sc-protected-byte"`
	TXScProtectedPacket          string `ros:"tx-sc-protected-packet"`
	TXTooLong                    string `ros:"tx-too-long"`
	TXUntagged                   string `ros:"tx-untagged"`
}
```

```go
// InterfaceMacsecAdd adalah argumen `/interface/macsec/add`.
type InterfaceMacsecAdd struct {
	// hexadecimal string value, max length 16
	Cak                          string `ros:"cak,omitempty"`
	// hexadecimal string value, max length 32
	Ckn                          string `ros:"ckn,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Profile                      string `ros:"profile,omitempty"`
}
```

```go
// InterfaceMacsecSet adalah argumen `/interface/macsec/set`.
type InterfaceMacsecSet struct {
	// hexadecimal string value, max length 16
	Cak                          string `ros:"cak,omitempty"`
	// hexadecimal string value, max length 32
	Ckn                          string `ros:"ckn,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	Profile                      string `ros:"profile,omitempty"`
}
```

### `/interface/macvlan`

- File: `interface_macvlan.json`
- Sumber field: completion 16

Command: `add`(13), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(13), `set`(13)

```go
// InterfaceMacvlanRow adalah kolom hasil `/interface/macvlan/print`.
type InterfaceMacvlanRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	ARP                          string `ros:"arp"`
	ARPTimeout                   string `ros:"arp-timeout"`
	// string value
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	Interface                    string `ros:"interface"`
	L2MTU                        string `ros:"l2mtu"`
	LoopProtect                  string `ros:"loop-protect"`
	// time interval
	LoopProtectDisableTime       string `ros:"loop-protect-disable-time"`
	// 1s..7101w3d6h28m15s (time interval)
	LoopProtectSendInterval      string `ros:"loop-protect-send-interval"`
	LoopProtectStatus            string `ros:"loop-protect-status"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	MACAddress                   string `ros:"mac-address"`
	Mode                         string `ros:"mode"`
	MTU                          string `ros:"mtu"`
	// string value
	Name                         string `ros:"name"`
	Running                      string `ros:"running"`
}
```

```go
// InterfaceMacvlanAdd adalah argumen `/interface/macvlan/add`.
type InterfaceMacvlanAdd struct {
	ARP                          string `ros:"arp,omitempty"`
	ARPTimeout                   string `ros:"arp-timeout,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	LoopProtect                  string `ros:"loop-protect,omitempty"`
	// time interval
	LoopProtectDisableTime       string `ros:"loop-protect-disable-time,omitempty"`
	// 1s..7101w3d6h28m15s (time interval)
	LoopProtectSendInterval      string `ros:"loop-protect-send-interval,omitempty"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	MACAddress                   string `ros:"mac-address,omitempty"`
	Mode                         string `ros:"mode,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
}
```

```go
// InterfaceMacvlanSet adalah argumen `/interface/macvlan/set`.
type InterfaceMacvlanSet struct {
	ARP                          string `ros:"arp,omitempty"`
	ARPTimeout                   string `ros:"arp-timeout,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	LoopProtect                  string `ros:"loop-protect,omitempty"`
	// time interval
	LoopProtectDisableTime       string `ros:"loop-protect-disable-time,omitempty"`
	// 1s..7101w3d6h28m15s (time interval)
	LoopProtectSendInterval      string `ros:"loop-protect-send-interval,omitempty"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	MACAddress                   string `ros:"mac-address,omitempty"`
	Mode                         string `ros:"mode,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
}
```

### `/interface/mesh`

- File: `interface_mesh.json`
- Sumber field: completion 22

Command: `add`(20), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(20), `set`(20), `traceroute`(8)

```go
// InterfaceMeshRow adalah kolom hasil `/interface/mesh/print`.
type InterfaceMeshRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	AdminMAC                     string `ros:"admin-mac"`
	ARP                          string `ros:"arp"`
	ARPTimeout                   string `ros:"arp-timeout"`
	AutoMAC                      string `ros:"auto-mac"`
	// string value
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	HwmpDefaultHoplimit          string `ros:"hwmp-default-hoplimit"`
	// 0s..1h (time interval)
	HwmpPrepLifetime             string `ros:"hwmp-prep-lifetime"`
	HwmpPreqDestinationOnly      string `ros:"hwmp-preq-destination-only"`
	HwmpPreqReplyAndForward      string `ros:"hwmp-preq-reply-and-forward"`
	HwmpPreqRetries              string `ros:"hwmp-preq-retries"`
	// 1s..30s (time interval)
	HwmpPreqWaitingTime          string `ros:"hwmp-preq-waiting-time"`
	// 0s..30m (time interval)
	HwmpRannInterval             string `ros:"hwmp-rann-interval"`
	// 0s..1h (time interval)
	HwmpRannLifetime             string `ros:"hwmp-rann-lifetime"`
	HwmpRannPropagationDelay     string `ros:"hwmp-rann-propagation-delay"`
	MACAddress                   string `ros:"mac-address"`
	MeshPortal                   string `ros:"mesh-portal"`
	MTU                          string `ros:"mtu"`
	// string value
	Name                         string `ros:"name"`
	ReoptimizePaths              string `ros:"reoptimize-paths"`
	Running                      string `ros:"running"`
}
```

```go
// InterfaceMeshAdd adalah argumen `/interface/mesh/add`.
type InterfaceMeshAdd struct {
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	AdminMAC                     string `ros:"admin-mac,omitempty"`
	ARP                          string `ros:"arp,omitempty"`
	ARPTimeout                   string `ros:"arp-timeout,omitempty"`
	AutoMAC                      string `ros:"auto-mac,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	HwmpDefaultHoplimit          string `ros:"hwmp-default-hoplimit,omitempty"`
	// 0s..1h (time interval)
	HwmpPrepLifetime             string `ros:"hwmp-prep-lifetime,omitempty"`
	HwmpPreqDestinationOnly      string `ros:"hwmp-preq-destination-only,omitempty"`
	HwmpPreqReplyAndForward      string `ros:"hwmp-preq-reply-and-forward,omitempty"`
	HwmpPreqRetries              string `ros:"hwmp-preq-retries,omitempty"`
	// 1s..30s (time interval)
	HwmpPreqWaitingTime          string `ros:"hwmp-preq-waiting-time,omitempty"`
	// 0s..30m (time interval)
	HwmpRannInterval             string `ros:"hwmp-rann-interval,omitempty"`
	// 0s..1h (time interval)
	HwmpRannLifetime             string `ros:"hwmp-rann-lifetime,omitempty"`
	HwmpRannPropagationDelay     string `ros:"hwmp-rann-propagation-delay,omitempty"`
	MeshPortal                   string `ros:"mesh-portal,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	ReoptimizePaths              string `ros:"reoptimize-paths,omitempty"`
}
```

```go
// InterfaceMeshSet adalah argumen `/interface/mesh/set`.
type InterfaceMeshSet struct {
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	AdminMAC                     string `ros:"admin-mac,omitempty"`
	ARP                          string `ros:"arp,omitempty"`
	ARPTimeout                   string `ros:"arp-timeout,omitempty"`
	AutoMAC                      string `ros:"auto-mac,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	HwmpDefaultHoplimit          string `ros:"hwmp-default-hoplimit,omitempty"`
	// 0s..1h (time interval)
	HwmpPrepLifetime             string `ros:"hwmp-prep-lifetime,omitempty"`
	HwmpPreqDestinationOnly      string `ros:"hwmp-preq-destination-only,omitempty"`
	HwmpPreqReplyAndForward      string `ros:"hwmp-preq-reply-and-forward,omitempty"`
	HwmpPreqRetries              string `ros:"hwmp-preq-retries,omitempty"`
	// 1s..30s (time interval)
	HwmpPreqWaitingTime          string `ros:"hwmp-preq-waiting-time,omitempty"`
	// 0s..30m (time interval)
	HwmpRannInterval             string `ros:"hwmp-rann-interval,omitempty"`
	// 0s..1h (time interval)
	HwmpRannLifetime             string `ros:"hwmp-rann-lifetime,omitempty"`
	HwmpRannPropagationDelay     string `ros:"hwmp-rann-propagation-delay,omitempty"`
	MeshPortal                   string `ros:"mesh-portal,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	ReoptimizePaths              string `ros:"reoptimize-paths,omitempty"`
}
```

### `/interface/monitor-traffic`

- File: `(dari katalog Go)`
- Katalog: `interface-monitor-traffic`
- Sumber field: katalog 10

Command: 

```go
// InterfaceMonitortrafficRow adalah kolom hasil `/interface/monitor-traffic/print`.
type InterfaceMonitortrafficRow struct {
	Name                         string `ros:"name"`
	RXBitsPerSecond              string `ros:"rx-bits-per-second"`
	TXBitsPerSecond              string `ros:"tx-bits-per-second"`
	RXPacketsPerSecond           string `ros:"rx-packets-per-second"`
	TXPacketsPerSecond           string `ros:"tx-packets-per-second"`
	RXDropsPerSecond             string `ros:"rx-drops-per-second"`
	TXDropsPerSecond             string `ros:"tx-drops-per-second"`
	RXErrorsPerSecond            string `ros:"rx-errors-per-second"`
	TXErrorsPerSecond            string `ros:"tx-errors-per-second"`
	TXQueueDropsPerSecond        string `ros:"tx-queue-drops-per-second"`
}
```

### `/interface/ovpn-client`

- File: `interface_ovpn-client.json`
- Sumber field: completion 26

Command: `add`(22), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `import-ovpn-configuration`(14), `monitor`(11), `print`(22), `remove`(1), `reset`(22), `set`(22)

```go
// InterfaceOvpnclientRow adalah kolom hasil `/interface/ovpn-client/print`.
type InterfaceOvpnclientRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	AddDefaultRoute              string `ros:"add-default-route"`
	Auth                         string `ros:"auth"`
	Certificate                  string `ros:"certificate"`
	Cipher                       string `ros:"cipher"`
	// string value
	Comment                      string `ros:"comment"`
	ConnectTo                    string `ros:"connect-to"`
	Disabled                     string `ros:"disabled"`
	DisconnectNotify             string `ros:"disconnect-notify"`
	HWCrypto                     string `ros:"hw-crypto"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	MACAddress                   string `ros:"mac-address"`
	MaxMTU                       string `ros:"max-mtu"`
	Mode                         string `ros:"mode"`
	// string value
	Name                         string `ros:"name"`
	// string value, max length 1000
	Password                     string `ros:"password"`
	Port                         string `ros:"port"`
	Profile                      string `ros:"profile"`
	Protocol                     string `ros:"protocol"`
	RouteNopull                  string `ros:"route-nopull"`
	Running                      string `ros:"running"`
	TLSAuth                      string `ros:"tls-auth"`
	TLSCrypt                     string `ros:"tls-crypt"`
	TLSVersion                   string `ros:"tls-version"`
	UsePeerDNS                   string `ros:"use-peer-dns"`
	// string value
	User                         string `ros:"user"`
	VerifyServerCertificate      string `ros:"verify-server-certificate"`
}
```

```go
// InterfaceOvpnclientAdd adalah argumen `/interface/ovpn-client/add`.
type InterfaceOvpnclientAdd struct {
	AddDefaultRoute              string `ros:"add-default-route,omitempty"`
	Auth                         string `ros:"auth,omitempty"`
	Certificate                  string `ros:"certificate,omitempty"`
	Cipher                       string `ros:"cipher,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	ConnectTo                    string `ros:"connect-to,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	DisconnectNotify             string `ros:"disconnect-notify,omitempty"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	MACAddress                   string `ros:"mac-address,omitempty"`
	MaxMTU                       string `ros:"max-mtu,omitempty"`
	Mode                         string `ros:"mode,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	// string value, max length 1000
	Password                     string `ros:"password,omitempty"`
	Port                         string `ros:"port,omitempty"`
	Profile                      string `ros:"profile,omitempty"`
	Protocol                     string `ros:"protocol,omitempty"`
	RouteNopull                  string `ros:"route-nopull,omitempty"`
	TLSVersion                   string `ros:"tls-version,omitempty"`
	UsePeerDNS                   string `ros:"use-peer-dns,omitempty"`
	// string value
	User                         string `ros:"user,omitempty"`
	VerifyServerCertificate      string `ros:"verify-server-certificate,omitempty"`
}
```

```go
// InterfaceOvpnclientSet adalah argumen `/interface/ovpn-client/set`.
type InterfaceOvpnclientSet struct {
	AddDefaultRoute              string `ros:"add-default-route,omitempty"`
	Auth                         string `ros:"auth,omitempty"`
	Certificate                  string `ros:"certificate,omitempty"`
	Cipher                       string `ros:"cipher,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	ConnectTo                    string `ros:"connect-to,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	DisconnectNotify             string `ros:"disconnect-notify,omitempty"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	MACAddress                   string `ros:"mac-address,omitempty"`
	MaxMTU                       string `ros:"max-mtu,omitempty"`
	Mode                         string `ros:"mode,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	// string value, max length 1000
	Password                     string `ros:"password,omitempty"`
	Port                         string `ros:"port,omitempty"`
	Profile                      string `ros:"profile,omitempty"`
	Protocol                     string `ros:"protocol,omitempty"`
	RouteNopull                  string `ros:"route-nopull,omitempty"`
	TLSVersion                   string `ros:"tls-version,omitempty"`
	UsePeerDNS                   string `ros:"use-peer-dns,omitempty"`
	// string value
	User                         string `ros:"user,omitempty"`
	VerifyServerCertificate      string `ros:"verify-server-certificate,omitempty"`
}
```

### `/interface/ovpn-server`

- File: `interface_ovpn-server.json`
- Sumber field: completion 12

Command: `add`(5), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `monitor`(11), `print`(21), `remove`(1), `reset`(5), `set`(5)

```go
// InterfaceOvpnserverRow adalah kolom hasil `/interface/ovpn-server/print`.
type InterfaceOvpnserverRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	ClientAddress                string `ros:"client-address"`
	// string value
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	Dynamic                      string `ros:"dynamic"`
	Encoding                     string `ros:"encoding"`
	HWCrypto                     string `ros:"hw-crypto"`
	MTU                          string `ros:"mtu"`
	// string value
	Name                         string `ros:"name"`
	Running                      string `ros:"running"`
	Uptime                       string `ros:"uptime"`
	// string value
	User                         string `ros:"user"`
}
```

```go
// InterfaceOvpnserverAdd adalah argumen `/interface/ovpn-server/add`.
type InterfaceOvpnserverAdd struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	// string value
	User                         string `ros:"user,omitempty"`
}
```

```go
// InterfaceOvpnserverSet adalah argumen `/interface/ovpn-server/set`.
type InterfaceOvpnserverSet struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	// string value
	User                         string `ros:"user,omitempty"`
}
```

### `/interface/ppp-client`

- File: `interface_ppp-client.json`
- Sumber field: completion 28

Command: `add`(27), `at-chat`(3), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `firmware-upgrade`(14), `get`(4), `info`(12), `monitor`(11), `print`(22), `remove`(1), `reset`(27), `scan`(6), `set`(27), `unset`(2)

```go
// InterfacePppclientRow adalah kolom hasil `/interface/ppp-client/print`.
type InterfacePppclientRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	AddDefaultRoute              string `ros:"add-default-route"`
	// pap|chap|mschap1|mschap2[,Allow*]
	Allow                        string `ros:"allow"`
	// string value
	Apn                          string `ros:"apn"`
	// string value
	Comment                      string `ros:"comment"`
	DataChannel                  string `ros:"data-channel"`
	DefaultRouteDistance         string `ros:"default-route-distance"`
	// string value
	DialCommand                  string `ros:"dial-command"`
	DialOnDemand                 string `ros:"dial-on-demand"`
	Disabled                     string `ros:"disabled"`
	InfoChannel                  string `ros:"info-channel"`
	KeepaliveTimeout             string `ros:"keepalive-timeout"`
	MaxMru                       string `ros:"max-mru"`
	MaxMTU                       string `ros:"max-mtu"`
	// string value
	ModemInit                    string `ros:"modem-init"`
	Mrru                         string `ros:"mrru"`
	// string value
	Name                         string `ros:"name"`
	NetworkMode                  string `ros:"network-mode"`
	NullModem                    string `ros:"null-modem"`
	// string value
	Password                     string `ros:"password"`
	// string value
	Phone                        string `ros:"phone"`
	// string value
	Pin                          string `ros:"pin"`
	Port                         string `ros:"port"`
	Profile                      string `ros:"profile"`
	// A.B.C.D (IP address)
	RemoteAddress                string `ros:"remote-address"`
	Running                      string `ros:"running"`
	UsePeerDNS                   string `ros:"use-peer-dns"`
	// string value
	User                         string `ros:"user"`
}
```

```go
// InterfacePppclientAdd adalah argumen `/interface/ppp-client/add`.
type InterfacePppclientAdd struct {
	AddDefaultRoute              string `ros:"add-default-route,omitempty"`
	// pap|chap|mschap1|mschap2[,Allow*]
	Allow                        string `ros:"allow,omitempty"`
	// string value
	Apn                          string `ros:"apn,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	DataChannel                  string `ros:"data-channel,omitempty"`
	DefaultRouteDistance         string `ros:"default-route-distance,omitempty"`
	// string value
	DialCommand                  string `ros:"dial-command,omitempty"`
	DialOnDemand                 string `ros:"dial-on-demand,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	InfoChannel                  string `ros:"info-channel,omitempty"`
	KeepaliveTimeout             string `ros:"keepalive-timeout,omitempty"`
	MaxMru                       string `ros:"max-mru,omitempty"`
	MaxMTU                       string `ros:"max-mtu,omitempty"`
	// string value
	ModemInit                    string `ros:"modem-init,omitempty"`
	Mrru                         string `ros:"mrru,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	NetworkMode                  string `ros:"network-mode,omitempty"`
	NullModem                    string `ros:"null-modem,omitempty"`
	// string value
	Password                     string `ros:"password,omitempty"`
	// string value
	Phone                        string `ros:"phone,omitempty"`
	// string value
	Pin                          string `ros:"pin,omitempty"`
	Port                         string `ros:"port,omitempty"`
	Profile                      string `ros:"profile,omitempty"`
	// A.B.C.D (IP address)
	RemoteAddress                string `ros:"remote-address,omitempty"`
	UsePeerDNS                   string `ros:"use-peer-dns,omitempty"`
	// string value
	User                         string `ros:"user,omitempty"`
}
```

```go
// InterfacePppclientSet adalah argumen `/interface/ppp-client/set`.
type InterfacePppclientSet struct {
	AddDefaultRoute              string `ros:"add-default-route,omitempty"`
	// pap|chap|mschap1|mschap2[,Allow*]
	Allow                        string `ros:"allow,omitempty"`
	// string value
	Apn                          string `ros:"apn,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	DataChannel                  string `ros:"data-channel,omitempty"`
	DefaultRouteDistance         string `ros:"default-route-distance,omitempty"`
	// string value
	DialCommand                  string `ros:"dial-command,omitempty"`
	DialOnDemand                 string `ros:"dial-on-demand,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	InfoChannel                  string `ros:"info-channel,omitempty"`
	KeepaliveTimeout             string `ros:"keepalive-timeout,omitempty"`
	MaxMru                       string `ros:"max-mru,omitempty"`
	MaxMTU                       string `ros:"max-mtu,omitempty"`
	// string value
	ModemInit                    string `ros:"modem-init,omitempty"`
	Mrru                         string `ros:"mrru,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	NetworkMode                  string `ros:"network-mode,omitempty"`
	NullModem                    string `ros:"null-modem,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	// string value
	Password                     string `ros:"password,omitempty"`
	// string value
	Phone                        string `ros:"phone,omitempty"`
	// string value
	Pin                          string `ros:"pin,omitempty"`
	Port                         string `ros:"port,omitempty"`
	Profile                      string `ros:"profile,omitempty"`
	// A.B.C.D (IP address)
	RemoteAddress                string `ros:"remote-address,omitempty"`
	UsePeerDNS                   string `ros:"use-peer-dns,omitempty"`
	// string value
	User                         string `ros:"user,omitempty"`
}
```

### `/interface/ppp-server`

- File: `interface_ppp-server.json`
- Sumber field: completion 15

Command: `add`(14), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `monitor`(11), `print`(21), `remove`(1), `reset`(14), `set`(14)

```go
// InterfacePppserverRow adalah kolom hasil `/interface/ppp-server/print`.
type InterfacePppserverRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	// pap|chap|mschap1|mschap2[,Authentication*]
	Authentication               string `ros:"authentication"`
	// string value
	Comment                      string `ros:"comment"`
	DataChannel                  string `ros:"data-channel"`
	Disabled                     string `ros:"disabled"`
	MaxMru                       string `ros:"max-mru"`
	MaxMTU                       string `ros:"max-mtu"`
	// string value
	ModemInit                    string `ros:"modem-init"`
	Mrru                         string `ros:"mrru"`
	// string value
	Name                         string `ros:"name"`
	NullModem                    string `ros:"null-modem"`
	Port                         string `ros:"port"`
	Profile                      string `ros:"profile"`
	RingCount                    string `ros:"ring-count"`
	Running                      string `ros:"running"`
}
```

```go
// InterfacePppserverAdd adalah argumen `/interface/ppp-server/add`.
type InterfacePppserverAdd struct {
	// pap|chap|mschap1|mschap2[,Authentication*]
	Authentication               string `ros:"authentication,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	DataChannel                  string `ros:"data-channel,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	MaxMru                       string `ros:"max-mru,omitempty"`
	MaxMTU                       string `ros:"max-mtu,omitempty"`
	// string value
	ModemInit                    string `ros:"modem-init,omitempty"`
	Mrru                         string `ros:"mrru,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	NullModem                    string `ros:"null-modem,omitempty"`
	Port                         string `ros:"port,omitempty"`
	Profile                      string `ros:"profile,omitempty"`
	RingCount                    string `ros:"ring-count,omitempty"`
}
```

```go
// InterfacePppserverSet adalah argumen `/interface/ppp-server/set`.
type InterfacePppserverSet struct {
	// pap|chap|mschap1|mschap2[,Authentication*]
	Authentication               string `ros:"authentication,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	DataChannel                  string `ros:"data-channel,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	MaxMru                       string `ros:"max-mru,omitempty"`
	MaxMTU                       string `ros:"max-mtu,omitempty"`
	// string value
	ModemInit                    string `ros:"modem-init,omitempty"`
	Mrru                         string `ros:"mrru,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	NullModem                    string `ros:"null-modem,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	Port                         string `ros:"port,omitempty"`
	Profile                      string `ros:"profile,omitempty"`
	RingCount                    string `ros:"ring-count,omitempty"`
}
```

### `/interface/pppoe-client`

- File: `interface_pppoe-client.json`
- Katalog: `pppoe-client`
- Sumber field: completion 22, katalog 14

Command: `add`(20), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `monitor`(11), `print`(22), `remove`(1), `reset`(20), `scan`(6), `set`(20), `unset`(2)

```go
// InterfacePppoeclientRow adalah kolom hasil `/interface/pppoe-client/print`.
type InterfacePppoeclientRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	// string value
	AcName                       string `ros:"ac-name"`
	AddDefaultRoute              string `ros:"add-default-route"`
	// pap|chap|mschap1|mschap2[,Allow*]
	Allow                        string `ros:"allow"`
	// string value
	Comment                      string `ros:"comment"`
	DefaultRouteDistance         string `ros:"default-route-distance"`
	DialOnDemand                 string `ros:"dial-on-demand"`
	Disabled                     string `ros:"disabled"`
	// hexadecimal string value
	HostUniq                     string `ros:"host-uniq"`
	Interface                    string `ros:"interface"`
	Invalid                      string `ros:"invalid"`
	KeepaliveTimeout             string `ros:"keepalive-timeout"`
	MaxMru                       string `ros:"max-mru"`
	MaxMTU                       string `ros:"max-mtu"`
	Mrru                         string `ros:"mrru"`
	// string value
	Name                         string `ros:"name"`
	// string value
	Password                     string `ros:"password"`
	Profile                      string `ros:"profile"`
	Running                      string `ros:"running"`
	// string value
	ServiceName                  string `ros:"service-name"`
	UsePeerDNS                   string `ros:"use-peer-dns"`
	// string value
	User                         string `ros:"user"`
	Status                       string `ros:"status"`
	Uptime                       string `ros:"uptime"`
	Encoding                     string `ros:"encoding"`
	LocalAddress                 string `ros:"local-address"`
	RemoteAddress                string `ros:"remote-address"`
}
```

```go
// InterfacePppoeclientAdd adalah argumen `/interface/pppoe-client/add`.
type InterfacePppoeclientAdd struct {
	// string value
	AcName                       string `ros:"ac-name,omitempty"`
	AddDefaultRoute              string `ros:"add-default-route,omitempty"`
	// pap|chap|mschap1|mschap2[,Allow*]
	Allow                        string `ros:"allow,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	DefaultRouteDistance         string `ros:"default-route-distance,omitempty"`
	DialOnDemand                 string `ros:"dial-on-demand,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	// hexadecimal string value
	HostUniq                     string `ros:"host-uniq,omitempty"`
	Interface                    string `ros:"interface,required,omitempty"`
	KeepaliveTimeout             string `ros:"keepalive-timeout,omitempty"`
	MaxMru                       string `ros:"max-mru,omitempty"`
	MaxMTU                       string `ros:"max-mtu,omitempty"`
	Mrru                         string `ros:"mrru,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	// string value
	Password                     string `ros:"password,required,omitempty"`
	Profile                      string `ros:"profile,omitempty"`
	// string value
	ServiceName                  string `ros:"service-name,omitempty"`
	UsePeerDNS                   string `ros:"use-peer-dns,omitempty"`
	// string value
	User                         string `ros:"user,required,omitempty"`
}
```

```go
// InterfacePppoeclientSet adalah argumen `/interface/pppoe-client/set`.
type InterfacePppoeclientSet struct {
	// string value
	AcName                       string `ros:"ac-name,omitempty"`
	AddDefaultRoute              string `ros:"add-default-route,omitempty"`
	// pap|chap|mschap1|mschap2[,Allow*]
	Allow                        string `ros:"allow,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	DefaultRouteDistance         string `ros:"default-route-distance,omitempty"`
	DialOnDemand                 string `ros:"dial-on-demand,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	// hexadecimal string value
	HostUniq                     string `ros:"host-uniq,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	KeepaliveTimeout             string `ros:"keepalive-timeout,omitempty"`
	MaxMru                       string `ros:"max-mru,omitempty"`
	MaxMTU                       string `ros:"max-mtu,omitempty"`
	Mrru                         string `ros:"mrru,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	// string value
	Password                     string `ros:"password,omitempty"`
	Profile                      string `ros:"profile,omitempty"`
	// string value
	ServiceName                  string `ros:"service-name,omitempty"`
	UsePeerDNS                   string `ros:"use-peer-dns,omitempty"`
	// string value
	User                         string `ros:"user,omitempty"`
}
```

### `/interface/pppoe-server`

- File: `interface_pppoe-server.json`
- Sumber field: completion 15

Command: `add`(6), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `monitor`(11), `print`(21), `remove`(1), `reset`(6), `set`(6)

```go
// InterfacePppoeserverRow adalah kolom hasil `/interface/pppoe-server/print`.
type InterfacePppoeserverRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	// string value
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	Dynamic                      string `ros:"dynamic"`
	Encoding                     string `ros:"encoding"`
	Interface                    string `ros:"interface"`
	Mru                          string `ros:"mru"`
	MTU                          string `ros:"mtu"`
	// string value
	Name                         string `ros:"name"`
	RemoteAddress                string `ros:"remote-address"`
	Running                      string `ros:"running"`
	// string value
	Service                      string `ros:"service"`
	Uptime                       string `ros:"uptime"`
	// string value
	User                         string `ros:"user"`
	Vid                          string `ros:"vid"`
}
```

```go
// InterfacePppoeserverAdd adalah argumen `/interface/pppoe-server/add`.
type InterfacePppoeserverAdd struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	// string value
	Service                      string `ros:"service,omitempty"`
	// string value
	User                         string `ros:"user,omitempty"`
}
```

```go
// InterfacePppoeserverSet adalah argumen `/interface/pppoe-server/set`.
type InterfacePppoeserverSet struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	// string value
	Service                      string `ros:"service,omitempty"`
	// string value
	User                         string `ros:"user,omitempty"`
}
```

### `/interface/pppoe-server/server`

- File: `(dari katalog Go)`
- Katalog: `pppoe-server`
- Sumber field: katalog 9

Command: `add`(8), `set`(8)

```go
// InterfacePppoeserverServerRow adalah kolom hasil `/interface/pppoe-server/server/print`.
type InterfacePppoeserverServerRow struct {
	ID                           string `ros:".id,key"`
	ServiceName                  string `ros:"service-name"`
	Interface                    string `ros:"interface"`
	DefaultProfile               string `ros:"default-profile"`
	Authentication               string `ros:"authentication"`
	KeepaliveTimeout             string `ros:"keepalive-timeout"`
	MaxSessions                  string `ros:"max-sessions"`
	OneSessionPerHost            string `ros:"one-session-per-host"`
	Disabled                     string `ros:"disabled"`
	Comment                      string `ros:"comment"`
}
```

```go
// InterfacePppoeserverServerAdd adalah argumen `/interface/pppoe-server/server/add`.
type InterfacePppoeserverServerAdd struct {
	ServiceName                  string `ros:"service-name,required,omitempty"`
	Interface                    string `ros:"interface,required,omitempty"`
	DefaultProfile               string `ros:"default-profile,omitempty"`
	Authentication               string `ros:"authentication,omitempty"`
	KeepaliveTimeout             string `ros:"keepalive-timeout,omitempty"`
	MaxSessions                  string `ros:"max-sessions,omitempty"`
	OneSessionPerHost            string `ros:"one-session-per-host,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

```go
// InterfacePppoeserverServerSet adalah argumen `/interface/pppoe-server/server/set`.
type InterfacePppoeserverServerSet struct {
	ServiceName                  string `ros:"service-name,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	DefaultProfile               string `ros:"default-profile,omitempty"`
	Authentication               string `ros:"authentication,omitempty"`
	KeepaliveTimeout             string `ros:"keepalive-timeout,omitempty"`
	MaxSessions                  string `ros:"max-sessions,omitempty"`
	OneSessionPerHost            string `ros:"one-session-per-host,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

### `/interface/pptp-client`

- File: `interface_pptp-client.json`
- Sumber field: completion 18

Command: 

```go
// InterfacePptpclientRow adalah kolom hasil `/interface/pptp-client/print`.
type InterfacePptpclientRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	AddDefaultRoute              string `ros:"add-default-route"`
	Allow                        string `ros:"allow"`
	Comment                      string `ros:"comment"`
	ConnectTo                    string `ros:"connect-to"`
	DefaultRouteDistance         string `ros:"default-route-distance"`
	DialOnDemand                 string `ros:"dial-on-demand"`
	Disabled                     string `ros:"disabled"`
	KeepaliveTimeout             string `ros:"keepalive-timeout"`
	MaxMru                       string `ros:"max-mru"`
	MaxMTU                       string `ros:"max-mtu"`
	Mrru                         string `ros:"mrru"`
	Name                         string `ros:"name"`
	Password                     string `ros:"password"`
	Profile                      string `ros:"profile"`
	Running                      string `ros:"running"`
	UsePeerDNS                   string `ros:"use-peer-dns"`
	User                         string `ros:"user"`
}
```

### `/interface/pptp-server`

- File: `interface_pptp-server.json`
- Sumber field: completion 12

Command: 

```go
// InterfacePptpserverRow adalah kolom hasil `/interface/pptp-server/print`.
type InterfacePptpserverRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	ClientAddress                string `ros:"client-address"`
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	Dynamic                      string `ros:"dynamic"`
	Encoding                     string `ros:"encoding"`
	Mru                          string `ros:"mru"`
	MTU                          string `ros:"mtu"`
	Name                         string `ros:"name"`
	Running                      string `ros:"running"`
	Uptime                       string `ros:"uptime"`
	User                         string `ros:"user"`
}
```

### `/interface/sstp-client`

- File: `interface_sstp-client.json`
- Sumber field: completion 28

Command: 

```go
// InterfaceSstpclientRow adalah kolom hasil `/interface/sstp-client/print`.
type InterfaceSstpclientRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	AddDefaultRoute              string `ros:"add-default-route"`
	AddSni                       string `ros:"add-sni"`
	Authentication               string `ros:"authentication"`
	Certificate                  string `ros:"certificate"`
	Ciphers                      string `ros:"ciphers"`
	Comment                      string `ros:"comment"`
	ConnectTo                    string `ros:"connect-to"`
	DefaultRouteDistance         string `ros:"default-route-distance"`
	DialOnDemand                 string `ros:"dial-on-demand"`
	Disabled                     string `ros:"disabled"`
	HTTPProxy                    string `ros:"http-proxy"`
	HWCrypto                     string `ros:"hw-crypto"`
	KeepaliveTimeout             string `ros:"keepalive-timeout"`
	MaxMru                       string `ros:"max-mru"`
	MaxMTU                       string `ros:"max-mtu"`
	Mrru                         string `ros:"mrru"`
	Name                         string `ros:"name"`
	Password                     string `ros:"password"`
	Pfs                          string `ros:"pfs"`
	Port                         string `ros:"port"`
	Profile                      string `ros:"profile"`
	ProxyPort                    string `ros:"proxy-port"`
	Running                      string `ros:"running"`
	TLSVersion                   string `ros:"tls-version"`
	User                         string `ros:"user"`
	VerifyServerAddressFromCertificate string `ros:"verify-server-address-from-certificate"`
	VerifyServerCertificate      string `ros:"verify-server-certificate"`
}
```

### `/interface/sstp-server`

- File: `interface_sstp-server.json`
- Sumber field: completion 12

Command: 

```go
// InterfaceSstpserverRow adalah kolom hasil `/interface/sstp-server/print`.
type InterfaceSstpserverRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	ClientAddress                string `ros:"client-address"`
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	Dynamic                      string `ros:"dynamic"`
	Encoding                     string `ros:"encoding"`
	HWCrypto                     string `ros:"hw-crypto"`
	MTU                          string `ros:"mtu"`
	Name                         string `ros:"name"`
	Running                      string `ros:"running"`
	Uptime                       string `ros:"uptime"`
	User                         string `ros:"user"`
}
```

### `/interface/vlan`

- File: `interface_vlan.json`
- Katalog: `interface-vlan`
- Sumber field: completion 19, katalog 7

Command: `add`(14), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(14), `set`(14)

```go
// InterfaceVLANRow adalah kolom hasil `/interface/vlan/print`.
type InterfaceVLANRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	ARP                          string `ros:"arp"`
	ARPTimeout                   string `ros:"arp-timeout"`
	// string value
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	HWOffloaded                  string `ros:"hw-offloaded"`
	Interface                    string `ros:"interface"`
	L2MTU                        string `ros:"l2mtu"`
	LoopProtect                  string `ros:"loop-protect"`
	// time interval
	LoopProtectDisableTime       string `ros:"loop-protect-disable-time"`
	// 1s..7101w3d6h28m15s (time interval)
	LoopProtectSendInterval      string `ros:"loop-protect-send-interval"`
	LoopProtectStatus            string `ros:"loop-protect-status"`
	MACAddress                   string `ros:"mac-address"`
	MTU                          string `ros:"mtu"`
	Mvrp                         string `ros:"mvrp"`
	// string value
	Name                         string `ros:"name"`
	Running                      string `ros:"running"`
	UseServiceTag                string `ros:"use-service-tag"`
	VLANID                       string `ros:"vlan-id"`
}
```

```go
// InterfaceVLANAdd adalah argumen `/interface/vlan/add`.
type InterfaceVLANAdd struct {
	ARP                          string `ros:"arp,omitempty"`
	ARPTimeout                   string `ros:"arp-timeout,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Interface                    string `ros:"interface,required,omitempty"`
	LoopProtect                  string `ros:"loop-protect,omitempty"`
	// time interval
	LoopProtectDisableTime       string `ros:"loop-protect-disable-time,omitempty"`
	// 1s..7101w3d6h28m15s (time interval)
	LoopProtectSendInterval      string `ros:"loop-protect-send-interval,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	Mvrp                         string `ros:"mvrp,omitempty"`
	// string value
	Name                         string `ros:"name,required,omitempty"`
	UseServiceTag                string `ros:"use-service-tag,omitempty"`
	VLANID                       string `ros:"vlan-id,required,omitempty"`
}
```

```go
// InterfaceVLANSet adalah argumen `/interface/vlan/set`.
type InterfaceVLANSet struct {
	ARP                          string `ros:"arp,omitempty"`
	ARPTimeout                   string `ros:"arp-timeout,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	LoopProtect                  string `ros:"loop-protect,omitempty"`
	// time interval
	LoopProtectDisableTime       string `ros:"loop-protect-disable-time,omitempty"`
	// 1s..7101w3d6h28m15s (time interval)
	LoopProtectSendInterval      string `ros:"loop-protect-send-interval,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	Mvrp                         string `ros:"mvrp,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	UseServiceTag                string `ros:"use-service-tag,omitempty"`
	VLANID                       string `ros:"vlan-id,omitempty"`
}
```

### `/interface/vpls`

- File: `interface_vpls.json`
- Sumber field: completion 25

Command: 

```go
// InterfaceVPLSRow adalah kolom hasil `/interface/vpls/print`.
type InterfaceVPLSRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	ARP                          string `ros:"arp"`
	ARPTimeout                   string `ros:"arp-timeout"`
	BGPSignaled                  string `ros:"bgp-signaled"`
	BGPVPLS                      string `ros:"bgp-vpls"`
	BGPVPLSPrfx                  string `ros:"bgp-vpls-prfx"`
	Bridge                       string `ros:"bridge"`
	BridgeCost                   string `ros:"bridge-cost"`
	BridgeHorizon                string `ros:"bridge-horizon"`
	BridgePvid                   string `ros:"bridge-pvid"`
	CiscoBGPSignaled             string `ros:"cisco-bgp-signaled"`
	CiscoStaticID                string `ros:"cisco-static-id"`
	Comment                      string `ros:"comment"`
	DisableRunningCheck          string `ros:"disable-running-check"`
	Disabled                     string `ros:"disabled"`
	Dynamic                      string `ros:"dynamic"`
	MACAddress                   string `ros:"mac-address"`
	MTU                          string `ros:"mtu"`
	Name                         string `ros:"name"`
	Peer                         string `ros:"peer"`
	PwControlWord                string `ros:"pw-control-word"`
	PwL2MTU                      string `ros:"pw-l2mtu"`
	PwType                       string `ros:"pw-type"`
	Running                      string `ros:"running"`
	VPLSID                       string `ros:"vpls-id"`
}
```

### `/interface/vrrp`

- File: `interface_vrrp.json`
- Sumber field: completion 34

Command: 

```go
// InterfaceVRRPRow adalah kolom hasil `/interface/vrrp/print`.
type InterfaceVRRPRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	ARP                          string `ros:"arp"`
	ARPTimeout                   string `ros:"arp-timeout"`
	Authentication               string `ros:"authentication"`
	Backup                       string `ros:"backup"`
	Comment                      string `ros:"comment"`
	ConnectionTrackingMode       string `ros:"connection-tracking-mode"`
	ConnectionTrackingPort       string `ros:"connection-tracking-port"`
	Disabled                     string `ros:"disabled"`
	Failure                      string `ros:"failure"`
	GroupAuthority               string `ros:"group-authority"`
	GroupMaster                  string `ros:"group-master"`
	GrpAuthority                 string `ros:"grp-authority"`
	GrpMember                    string `ros:"grp-member"`
	Interface                    string `ros:"interface"`
	Interval                     string `ros:"interval"`
	Invalid                      string `ros:"invalid"`
	MACAddress                   string `ros:"mac-address"`
	Master                       string `ros:"master"`
	MTU                          string `ros:"mtu"`
	Name                         string `ros:"name"`
	OnBackup                     string `ros:"on-backup"`
	OnFail                       string `ros:"on-fail"`
	OnMaster                     string `ros:"on-master"`
	Password                     string `ros:"password"`
	PreemptionMode               string `ros:"preemption-mode"`
	Priority                     string `ros:"priority"`
	RemoteAddress                string `ros:"remote-address"`
	Running                      string `ros:"running"`
	SyncConnectionTracking       string `ros:"sync-connection-tracking"`
	V3ChecksumAsV2               string `ros:"v3-checksum-as-v2"`
	V3Protocol                   string `ros:"v3-protocol"`
	Version                      string `ros:"version"`
	Vrid                         string `ros:"vrid"`
}
```

### `/interface/vxlan`

- File: `interface_vxlan.json`
- Sumber field: completion 32

Command: 

```go
// InterfaceVXLANRow adalah kolom hasil `/interface/vxlan/print`.
type InterfaceVXLANRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	AllowFastPath                string `ros:"allow-fast-path"`
	ARP                          string `ros:"arp"`
	ARPTimeout                   string `ros:"arp-timeout"`
	Bridge                       string `ros:"bridge"`
	BridgePvid                   string `ros:"bridge-pvid"`
	Checksum                     string `ros:"checksum"`
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	DontFragment                 string `ros:"dont-fragment"`
	Group                        string `ros:"group"`
	HW                           string `ros:"hw"`
	HWOffloaded                  string `ros:"hw-offloaded"`
	Interface                    string `ros:"interface"`
	L2MTU                        string `ros:"l2mtu"`
	Learning                     string `ros:"learning"`
	LocalAddress                 string `ros:"local-address"`
	LoopProtect                  string `ros:"loop-protect"`
	LoopProtectDisableTime       string `ros:"loop-protect-disable-time"`
	LoopProtectSendInterval      string `ros:"loop-protect-send-interval"`
	LoopProtectStatus            string `ros:"loop-protect-status"`
	MACAddress                   string `ros:"mac-address"`
	MaxFdbSize                   string `ros:"max-fdb-size"`
	MTU                          string `ros:"mtu"`
	Name                         string `ros:"name"`
	Port                         string `ros:"port"`
	RemCsum                      string `ros:"rem-csum"`
	Running                      string `ros:"running"`
	TTL                          string `ros:"ttl"`
	Vni                          string `ros:"vni"`
	VtepVRF                      string `ros:"vtep-vrf"`
	VtepsIPVersion               string `ros:"vteps-ip-version"`
}
```

### `/interface/wifi`

- File: `interface_wifi.json`
- Sumber field: completion 141

Command: `add`(134), `comment`(2), `devel`(6), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(2), `flat-snoop`(7), `frequency-scan`(9), `get`(5), `liberate`(1), `monitor`(11), `print`(23), `remove`(1), `reset`(134), `reset-mac-address`(1), `roam`(3), `scan`(9), `set`(134), `sniffer`(12), `spectral-history`(9), `spectral-scan`(12), `unset`(2), `wps-client`(12), `wps-push-button`(1)

```go
// InterfaceWiFiRow adalah kolom hasil `/interface/wifi/print`.
type InterfaceWiFiRow struct {
	ID                           string `ros:".id,key"`
	Aaa                          string `ros:"aaa"`
	// string value
	AaaCalledFormat              string `ros:"aaa.called-format"`
	// string value
	AaaCallingFormat             string `ros:"aaa.calling-format"`
	AaaInterimUpdate             string `ros:"aaa.interim-update"`
	AaaMACCaching                string `ros:"aaa.mac-caching"`
	// string value
	AaaNasIdentifier             string `ros:"aaa.nas-identifier"`
	// string value
	AaaPasswordFormat            string `ros:"aaa.password-format"`
	// string value
	AaaUsernameFormat            string `ros:"aaa.username-format"`
	About                        string `ros:"about"`
	ARP                          string `ros:"arp"`
	ARPTimeout                   string `ros:"arp-timeout"`
	Bound                        string `ros:"bound"`
	CAP                          string `ros:"cap"`
	Channel                      string `ros:"channel"`
	ChannelAfc                   string `ros:"channel.afc"`
	ChannelBand                  string `ros:"channel.band"`
	ChannelDeprioritizeUniiX3X4  string `ros:"channel.deprioritize-unii-3-4"`
	ChannelFrequency             string `ros:"channel.frequency"`
	ChannelPreamblePuncturing    string `ros:"channel.preamble-puncturing"`
	// 1s..42w6d (time interval)
	ChannelReselectInterval      string `ros:"channel.reselect-interval"`
	ChannelReselectTime          string `ros:"channel.reselect-time"`
	ChannelSecondaryFrequency    string `ros:"channel.secondary-frequency"`
	ChannelSkipDfsChannels       string `ros:"channel.skip-dfs-channels"`
	ChannelWidth                 string `ros:"channel.width"`
	// string value
	Comment                      string `ros:"comment"`
	Configuration                string `ros:"configuration"`
	ConfigurationAntennaGain     string `ros:"configuration.antenna-gain"`
	// 100ms..1s (time interval)
	ConfigurationBeaconInterval  string `ros:"configuration.beacon-interval"`
	// 0|1|2|3|4|5|6|7[,ConfigurationChains*]
	ConfigurationChains          string `ros:"configuration.chains"`
	ConfigurationCountry         string `ros:"configuration.country"`
	ConfigurationDistance        string `ros:"configuration.distance"`
	ConfigurationDtimPeriod      string `ros:"configuration.dtim-period"`
	ConfigurationHideSsid        string `ros:"configuration.hide-ssid"`
	ConfigurationHWProtectionMode string `ros:"configuration.hw-protection-mode"`
	ConfigurationInstallation    string `ros:"configuration.installation"`
	ConfigurationManager         string `ros:"configuration.manager"`
	ConfigurationMaxClients      string `ros:"configuration.max-clients"`
	ConfigurationMode            string `ros:"configuration.mode"`
	ConfigurationMulticastEnhance string `ros:"configuration.multicast-enhance"`
	ConfigurationQosClassifier   string `ros:"configuration.qos-classifier"`
	// string value, max length 32
	ConfigurationSsid            string `ros:"configuration.ssid"`
	ConfigurationStationRoaming  string `ros:"configuration.station-roaming"`
	// 0|1|2|3|4|5|6|7[,ConfigurationTxChains*]
	ConfigurationTXChains        string `ros:"configuration.tx-chains"`
	ConfigurationTXPower         string `ros:"configuration.tx-power"`
	Datapath                     string `ros:"datapath"`
	DatapathBridge               string `ros:"datapath.bridge"`
	DatapathBridgeCost           string `ros:"datapath.bridge-cost"`
	DatapathBridgeHorizon        string `ros:"datapath.bridge-horizon"`
	DatapathClientIsolation      string `ros:"datapath.client-isolation"`
	DatapathInterfaceList        string `ros:"datapath.interface-list"`
	DatapathTrafficProcessing    string `ros:"datapath.traffic-processing"`
	DatapathVLANID               string `ros:"datapath.vlan-id"`
	DefaultName                  string `ros:"default-name"`
	DisableRunningCheck          string `ros:"disable-running-check"`
	Disabled                     string `ros:"disabled"`
	Dynamic                      string `ros:"dynamic"`
	Inactive                     string `ros:"inactive"`
	Interworking                 string `ros:"interworking"`
	// string value, min length 3, max length 3
	InterworkingX3gppInfo        string `ros:"interworking.3gpp-info"`
	// hexadecimal string value
	InterworkingX3gppInfoRaw     string `ros:"interworking.3gpp-info-raw"`
	InterworkingAuthenticationTypes string `ros:"interworking.authentication-types"`
	InterworkingConnectionCapabilities string `ros:"interworking.connection-capabilities"`
	// string value
	InterworkingDomainNames      string `ros:"interworking.domain-names"`
	InterworkingEsr              string `ros:"interworking.esr"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	InterworkingHessid           string `ros:"interworking.hessid"`
	InterworkingHotspot20        string `ros:"interworking.hotspot20"`
	InterworkingHotspot20Dgaf    string `ros:"interworking.hotspot20-dgaf"`
	InterworkingInternet         string `ros:"interworking.internet"`
	InterworkingIPv4Availability string `ros:"interworking.ipv4-availability"`
	InterworkingIPv6Availability string `ros:"interworking.ipv6-availability"`
	InterworkingNetworkType      string `ros:"interworking.network-type"`
	InterworkingOperationalClasses string `ros:"interworking.operational-classes"`
	// string value
	InterworkingOperatorNames    string `ros:"interworking.operator-names"`
	// string value
	InterworkingRealms           string `ros:"interworking.realms"`
	// hexadecimal string value
	InterworkingRealmsRaw        string `ros:"interworking.realms-raw"`
	// string value
	InterworkingRoamingOis       string `ros:"interworking.roaming-ois"`
	InterworkingUesa             string `ros:"interworking.uesa"`
	InterworkingVenue            string `ros:"interworking.venue"`
	// string value
	InterworkingVenueNames       string `ros:"interworking.venue-names"`
	InterworkingWANAtCapacity    string `ros:"interworking.wan-at-capacity"`
	InterworkingWANDownlink      string `ros:"interworking.wan-downlink"`
	InterworkingWANDownlinkLoad  string `ros:"interworking.wan-downlink-load"`
	InterworkingWANMeasurementDuration string `ros:"interworking.wan-measurement-duration"`
	InterworkingWANStatus        string `ros:"interworking.wan-status"`
	InterworkingWANSymmetric     string `ros:"interworking.wan-symmetric"`
	InterworkingWANUplink        string `ros:"interworking.wan-uplink"`
	InterworkingWANUplinkLoad    string `ros:"interworking.wan-uplink-load"`
	L2MTU                        string `ros:"l2mtu"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	MACAddress                   string `ros:"mac-address"`
	Master                       string `ros:"master"`
	MasterInterface              string `ros:"master-interface"`
	MldInterface                 string `ros:"mld-interface"`
	// string value
	MldName                      string `ros:"mld-name"`
	MTU                          string `ros:"mtu"`
	// string value
	Name                         string `ros:"name"`
	Network                      string `ros:"network"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	RadioMAC                     string `ros:"radio-mac"`
	Running                      string `ros:"running"`
	Security                     string `ros:"security"`
	// wpa-psk|wpa2-psk|wpa2-psk-sha2|wpa-eap|wpa2-eap|wpa3-psk|owe|wpa3-eap|wpa3-eap-192[,SecurityAuthenticationTypes*]
	SecurityAuthenticationTypes  string `ros:"security.authentication-types"`
	SecurityBeaconProtection     string `ros:"security.beacon-protection"`
	// string value
	SecurityConnectGroup         string `ros:"security.connect-group"`
	SecurityConnectPriority      string `ros:"security.connect-priority"`
	// 19|20|21[,SecurityDhGroups*]
	SecurityDhGroups             string `ros:"security.dh-groups"`
	SecurityDisablePmkid         string `ros:"security.disable-pmkid"`
	SecurityEapAccounting        string `ros:"security.eap-accounting"`
	// string value
	SecurityEapAnonymousIdentity string `ros:"security.eap-anonymous-identity"`
	SecurityEapCertificateMode   string `ros:"security.eap-certificate-mode"`
	SecurityEapMethods           string `ros:"security.eap-methods"`
	// string value
	SecurityEapPassword          string `ros:"security.eap-password"`
	SecurityEapTLSCertificate    string `ros:"security.eap-tls-certificate"`
	// string value
	SecurityEapUsername          string `ros:"security.eap-username"`
	// tkip|ccmp|gcmp|ccmp-256|gcmp-256[,SecurityEncryption*]
	SecurityEncryption           string `ros:"security.encryption"`
	SecurityFt                   string `ros:"security.ft"`
	SecurityFtMobilityDomain     string `ros:"security.ft-mobility-domain"`
	// string value, min length 2, max length 96
	SecurityFtNasIdentifier      string `ros:"security.ft-nas-identifier"`
	SecurityFtOverDs             string `ros:"security.ft-over-ds"`
	SecurityFtPreserveVlanid     string `ros:"security.ft-preserve-vlanid"`
	// 1s..6w3d12h15m (time interval)
	SecurityFtR0KeyLifetime      string `ros:"security.ft-r0-key-lifetime"`
	// 0s..1m10s (time interval)
	SecurityFtReassociationDeadline string `ros:"security.ft-reassociation-deadline"`
	SecurityGroupEncryption      string `ros:"security.group-encryption"`
	// 30s..7101w3d6h28m15s (time interval)
	SecurityGroupKeyUpdate       string `ros:"security.group-key-update"`
	SecurityManagementEncryption string `ros:"security.management-encryption"`
	SecurityManagementProtection string `ros:"security.management-protection"`
	SecurityMultiPassphraseGroup string `ros:"security.multi-passphrase-group"`
	SecurityOweTransitionInterface string `ros:"security.owe-transition-interface"`
	// string value
	SecurityPassphrase           string `ros:"security.passphrase"`
	SecuritySaeAntiCloggingThreshold string `ros:"security.sae-anti-clogging-threshold"`
	SecuritySaeMaxFailureRate    string `ros:"security.sae-max-failure-rate"`
	SecuritySaePwe               string `ros:"security.sae-pwe"`
	SecurityWPS                  string `ros:"security.wps"`
	Steering                     string `ros:"steering"`
	SteeringX2gProbeDelay        string `ros:"steering.2g-probe-delay"`
	// string value
	SteeringNeighborGroup        string `ros:"steering.neighbor-group"`
	SteeringRrm                  string `ros:"steering.rrm"`
	SteeringTransitionRequestCount string `ros:"steering.transition-request-count"`
	// time interval
	SteeringTransitionRequestPeriod string `ros:"steering.transition-request-period"`
	SteeringTransitionThreshold  string `ros:"steering.transition-threshold"`
	// time interval
	SteeringTransitionThresholdTime string `ros:"steering.transition-threshold-time"`
	SteeringTransitionTime       string `ros:"steering.transition-time"`
	SteeringWnm                  string `ros:"steering.wnm"`
}
```

```go
// InterfaceWiFiAdd adalah argumen `/interface/wifi/add`.
type InterfaceWiFiAdd struct {
	Aaa                          string `ros:"aaa,omitempty"`
	// string value
	AaaCalledFormat              string `ros:"aaa.called-format,omitempty"`
	// string value
	AaaCallingFormat             string `ros:"aaa.calling-format,omitempty"`
	AaaInterimUpdate             string `ros:"aaa.interim-update,omitempty"`
	AaaMACCaching                string `ros:"aaa.mac-caching,omitempty"`
	// string value
	AaaNasIdentifier             string `ros:"aaa.nas-identifier,omitempty"`
	// string value
	AaaPasswordFormat            string `ros:"aaa.password-format,omitempty"`
	// string value
	AaaUsernameFormat            string `ros:"aaa.username-format,omitempty"`
	ARP                          string `ros:"arp,omitempty"`
	ARPTimeout                   string `ros:"arp-timeout,omitempty"`
	Channel                      string `ros:"channel,omitempty"`
	ChannelAfc                   string `ros:"channel.afc,omitempty"`
	ChannelBand                  string `ros:"channel.band,omitempty"`
	ChannelDeprioritizeUniiX3X4  string `ros:"channel.deprioritize-unii-3-4,omitempty"`
	ChannelFrequency             string `ros:"channel.frequency,omitempty"`
	ChannelPreamblePuncturing    string `ros:"channel.preamble-puncturing,omitempty"`
	// 1s..42w6d (time interval)
	ChannelReselectInterval      string `ros:"channel.reselect-interval,omitempty"`
	ChannelReselectTime          string `ros:"channel.reselect-time,omitempty"`
	ChannelSecondaryFrequency    string `ros:"channel.secondary-frequency,omitempty"`
	ChannelSkipDfsChannels       string `ros:"channel.skip-dfs-channels,omitempty"`
	ChannelWidth                 string `ros:"channel.width,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Configuration                string `ros:"configuration,omitempty"`
	ConfigurationAntennaGain     string `ros:"configuration.antenna-gain,omitempty"`
	// 100ms..1s (time interval)
	ConfigurationBeaconInterval  string `ros:"configuration.beacon-interval,omitempty"`
	// 0|1|2|3|4|5|6|7[,ConfigurationChains*]
	ConfigurationChains          string `ros:"configuration.chains,omitempty"`
	ConfigurationCountry         string `ros:"configuration.country,omitempty"`
	ConfigurationDistance        string `ros:"configuration.distance,omitempty"`
	ConfigurationDtimPeriod      string `ros:"configuration.dtim-period,omitempty"`
	ConfigurationHideSsid        string `ros:"configuration.hide-ssid,omitempty"`
	ConfigurationHWProtectionMode string `ros:"configuration.hw-protection-mode,omitempty"`
	ConfigurationInstallation    string `ros:"configuration.installation,omitempty"`
	ConfigurationManager         string `ros:"configuration.manager,omitempty"`
	ConfigurationMaxClients      string `ros:"configuration.max-clients,omitempty"`
	ConfigurationMode            string `ros:"configuration.mode,omitempty"`
	ConfigurationMulticastEnhance string `ros:"configuration.multicast-enhance,omitempty"`
	ConfigurationQosClassifier   string `ros:"configuration.qos-classifier,omitempty"`
	// string value, max length 32
	ConfigurationSsid            string `ros:"configuration.ssid,omitempty"`
	ConfigurationStationRoaming  string `ros:"configuration.station-roaming,omitempty"`
	// 0|1|2|3|4|5|6|7[,ConfigurationTxChains*]
	ConfigurationTXChains        string `ros:"configuration.tx-chains,omitempty"`
	ConfigurationTXPower         string `ros:"configuration.tx-power,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Datapath                     string `ros:"datapath,omitempty"`
	DatapathBridge               string `ros:"datapath.bridge,omitempty"`
	DatapathBridgeCost           string `ros:"datapath.bridge-cost,omitempty"`
	DatapathBridgeHorizon        string `ros:"datapath.bridge-horizon,omitempty"`
	DatapathClientIsolation      string `ros:"datapath.client-isolation,omitempty"`
	DatapathInterfaceList        string `ros:"datapath.interface-list,omitempty"`
	DatapathOpenflowSwitch       string `ros:"datapath.openflow-switch,omitempty"`
	DatapathTrafficProcessing    string `ros:"datapath.traffic-processing,omitempty"`
	DatapathVLANID               string `ros:"datapath.vlan-id,omitempty"`
	DisableRunningCheck          string `ros:"disable-running-check,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Interworking                 string `ros:"interworking,omitempty"`
	// string value, min length 3, max length 3
	InterworkingX3gppInfo        string `ros:"interworking.3gpp-info,omitempty"`
	// hexadecimal string value
	InterworkingX3gppInfoRaw     string `ros:"interworking.3gpp-info-raw,omitempty"`
	InterworkingAuthenticationTypes string `ros:"interworking.authentication-types,omitempty"`
	InterworkingConnectionCapabilities string `ros:"interworking.connection-capabilities,omitempty"`
	// string value
	InterworkingDomainNames      string `ros:"interworking.domain-names,omitempty"`
	InterworkingEsr              string `ros:"interworking.esr,omitempty"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	InterworkingHessid           string `ros:"interworking.hessid,omitempty"`
	InterworkingHotspot20        string `ros:"interworking.hotspot20,omitempty"`
	InterworkingHotspot20Dgaf    string `ros:"interworking.hotspot20-dgaf,omitempty"`
	InterworkingInternet         string `ros:"interworking.internet,omitempty"`
	InterworkingIPv4Availability string `ros:"interworking.ipv4-availability,omitempty"`
	InterworkingIPv6Availability string `ros:"interworking.ipv6-availability,omitempty"`
	InterworkingNetworkType      string `ros:"interworking.network-type,omitempty"`
	InterworkingOperationalClasses string `ros:"interworking.operational-classes,omitempty"`
	// string value
	InterworkingOperatorNames    string `ros:"interworking.operator-names,omitempty"`
	// string value
	InterworkingRealms           string `ros:"interworking.realms,omitempty"`
	// hexadecimal string value
	InterworkingRealmsRaw        string `ros:"interworking.realms-raw,omitempty"`
	// string value
	InterworkingRoamingOis       string `ros:"interworking.roaming-ois,omitempty"`
	InterworkingUesa             string `ros:"interworking.uesa,omitempty"`
	InterworkingVenue            string `ros:"interworking.venue,omitempty"`
	// string value
	InterworkingVenueNames       string `ros:"interworking.venue-names,omitempty"`
	InterworkingWANAtCapacity    string `ros:"interworking.wan-at-capacity,omitempty"`
	InterworkingWANDownlink      string `ros:"interworking.wan-downlink,omitempty"`
	InterworkingWANDownlinkLoad  string `ros:"interworking.wan-downlink-load,omitempty"`
	InterworkingWANMeasurementDuration string `ros:"interworking.wan-measurement-duration,omitempty"`
	InterworkingWANStatus        string `ros:"interworking.wan-status,omitempty"`
	InterworkingWANSymmetric     string `ros:"interworking.wan-symmetric,omitempty"`
	InterworkingWANUplink        string `ros:"interworking.wan-uplink,omitempty"`
	InterworkingWANUplinkLoad    string `ros:"interworking.wan-uplink-load,omitempty"`
	L2MTU                        string `ros:"l2mtu,omitempty"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	MACAddress                   string `ros:"mac-address,omitempty"`
	MasterInterface              string `ros:"master-interface,omitempty"`
	MldInterface                 string `ros:"mld-interface,omitempty"`
	// string value
	MldName                      string `ros:"mld-name,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	RadioMAC                     string `ros:"radio-mac,omitempty"`
	Security                     string `ros:"security,omitempty"`
	// wpa-psk|wpa2-psk|wpa2-psk-sha2|wpa-eap|wpa2-eap|wpa3-psk|owe|wpa3-eap|wpa3-eap-192[,SecurityAuthenticationTypes*]
	SecurityAuthenticationTypes  string `ros:"security.authentication-types,omitempty"`
	SecurityBeaconProtection     string `ros:"security.beacon-protection,omitempty"`
	// string value
	SecurityConnectGroup         string `ros:"security.connect-group,omitempty"`
	SecurityConnectPriority      string `ros:"security.connect-priority,omitempty"`
	// 19|20|21[,SecurityDhGroups*]
	SecurityDhGroups             string `ros:"security.dh-groups,omitempty"`
	SecurityDisablePmkid         string `ros:"security.disable-pmkid,omitempty"`
	SecurityEapAccounting        string `ros:"security.eap-accounting,omitempty"`
	// string value
	SecurityEapAnonymousIdentity string `ros:"security.eap-anonymous-identity,omitempty"`
	SecurityEapCertificateMode   string `ros:"security.eap-certificate-mode,omitempty"`
	SecurityEapMethods           string `ros:"security.eap-methods,omitempty"`
	// string value
	SecurityEapPassword          string `ros:"security.eap-password,omitempty"`
	SecurityEapTLSCertificate    string `ros:"security.eap-tls-certificate,omitempty"`
	// string value
	SecurityEapUsername          string `ros:"security.eap-username,omitempty"`
	// tkip|ccmp|gcmp|ccmp-256|gcmp-256[,SecurityEncryption*]
	SecurityEncryption           string `ros:"security.encryption,omitempty"`
	SecurityFt                   string `ros:"security.ft,omitempty"`
	SecurityFtMobilityDomain     string `ros:"security.ft-mobility-domain,omitempty"`
	// string value, min length 2, max length 96
	SecurityFtNasIdentifier      string `ros:"security.ft-nas-identifier,omitempty"`
	SecurityFtOverDs             string `ros:"security.ft-over-ds,omitempty"`
	SecurityFtPreserveVlanid     string `ros:"security.ft-preserve-vlanid,omitempty"`
	// 1s..6w3d12h15m (time interval)
	SecurityFtR0KeyLifetime      string `ros:"security.ft-r0-key-lifetime,omitempty"`
	// 0s..1m10s (time interval)
	SecurityFtReassociationDeadline string `ros:"security.ft-reassociation-deadline,omitempty"`
	SecurityGroupEncryption      string `ros:"security.group-encryption,omitempty"`
	// 30s..7101w3d6h28m15s (time interval)
	SecurityGroupKeyUpdate       string `ros:"security.group-key-update,omitempty"`
	SecurityManagementEncryption string `ros:"security.management-encryption,omitempty"`
	SecurityManagementProtection string `ros:"security.management-protection,omitempty"`
	SecurityMultiPassphraseGroup string `ros:"security.multi-passphrase-group,omitempty"`
	SecurityOweTransitionInterface string `ros:"security.owe-transition-interface,omitempty"`
	// string value
	SecurityPassphrase           string `ros:"security.passphrase,omitempty"`
	SecuritySaeAntiCloggingThreshold string `ros:"security.sae-anti-clogging-threshold,omitempty"`
	SecuritySaeMaxFailureRate    string `ros:"security.sae-max-failure-rate,omitempty"`
	SecuritySaePwe               string `ros:"security.sae-pwe,omitempty"`
	SecurityWPS                  string `ros:"security.wps,omitempty"`
	Steering                     string `ros:"steering,omitempty"`
	SteeringX2gProbeDelay        string `ros:"steering.2g-probe-delay,omitempty"`
	// string value
	SteeringNeighborGroup        string `ros:"steering.neighbor-group,omitempty"`
	SteeringRrm                  string `ros:"steering.rrm,omitempty"`
	SteeringTransitionRequestCount string `ros:"steering.transition-request-count,omitempty"`
	// time interval
	SteeringTransitionRequestPeriod string `ros:"steering.transition-request-period,omitempty"`
	SteeringTransitionThreshold  string `ros:"steering.transition-threshold,omitempty"`
	// time interval
	SteeringTransitionThresholdTime string `ros:"steering.transition-threshold-time,omitempty"`
	SteeringTransitionTime       string `ros:"steering.transition-time,omitempty"`
	SteeringWnm                  string `ros:"steering.wnm,omitempty"`
}
```

```go
// InterfaceWiFiSet adalah argumen `/interface/wifi/set`.
type InterfaceWiFiSet struct {
	Aaa                          string `ros:"aaa,omitempty"`
	// string value
	AaaCalledFormat              string `ros:"aaa.called-format,omitempty"`
	// string value
	AaaCallingFormat             string `ros:"aaa.calling-format,omitempty"`
	AaaInterimUpdate             string `ros:"aaa.interim-update,omitempty"`
	AaaMACCaching                string `ros:"aaa.mac-caching,omitempty"`
	// string value
	AaaNasIdentifier             string `ros:"aaa.nas-identifier,omitempty"`
	// string value
	AaaPasswordFormat            string `ros:"aaa.password-format,omitempty"`
	// string value
	AaaUsernameFormat            string `ros:"aaa.username-format,omitempty"`
	ARP                          string `ros:"arp,omitempty"`
	ARPTimeout                   string `ros:"arp-timeout,omitempty"`
	Channel                      string `ros:"channel,omitempty"`
	ChannelAfc                   string `ros:"channel.afc,omitempty"`
	ChannelBand                  string `ros:"channel.band,omitempty"`
	ChannelDeprioritizeUniiX3X4  string `ros:"channel.deprioritize-unii-3-4,omitempty"`
	ChannelFrequency             string `ros:"channel.frequency,omitempty"`
	ChannelPreamblePuncturing    string `ros:"channel.preamble-puncturing,omitempty"`
	// 1s..42w6d (time interval)
	ChannelReselectInterval      string `ros:"channel.reselect-interval,omitempty"`
	ChannelReselectTime          string `ros:"channel.reselect-time,omitempty"`
	ChannelSecondaryFrequency    string `ros:"channel.secondary-frequency,omitempty"`
	ChannelSkipDfsChannels       string `ros:"channel.skip-dfs-channels,omitempty"`
	ChannelWidth                 string `ros:"channel.width,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Configuration                string `ros:"configuration,omitempty"`
	ConfigurationAntennaGain     string `ros:"configuration.antenna-gain,omitempty"`
	// 100ms..1s (time interval)
	ConfigurationBeaconInterval  string `ros:"configuration.beacon-interval,omitempty"`
	// 0|1|2|3|4|5|6|7[,ConfigurationChains*]
	ConfigurationChains          string `ros:"configuration.chains,omitempty"`
	ConfigurationCountry         string `ros:"configuration.country,omitempty"`
	ConfigurationDistance        string `ros:"configuration.distance,omitempty"`
	ConfigurationDtimPeriod      string `ros:"configuration.dtim-period,omitempty"`
	ConfigurationHideSsid        string `ros:"configuration.hide-ssid,omitempty"`
	ConfigurationHWProtectionMode string `ros:"configuration.hw-protection-mode,omitempty"`
	ConfigurationInstallation    string `ros:"configuration.installation,omitempty"`
	ConfigurationManager         string `ros:"configuration.manager,omitempty"`
	ConfigurationMaxClients      string `ros:"configuration.max-clients,omitempty"`
	ConfigurationMode            string `ros:"configuration.mode,omitempty"`
	ConfigurationMulticastEnhance string `ros:"configuration.multicast-enhance,omitempty"`
	ConfigurationQosClassifier   string `ros:"configuration.qos-classifier,omitempty"`
	// string value, max length 32
	ConfigurationSsid            string `ros:"configuration.ssid,omitempty"`
	ConfigurationStationRoaming  string `ros:"configuration.station-roaming,omitempty"`
	// 0|1|2|3|4|5|6|7[,ConfigurationTxChains*]
	ConfigurationTXChains        string `ros:"configuration.tx-chains,omitempty"`
	ConfigurationTXPower         string `ros:"configuration.tx-power,omitempty"`
	Datapath                     string `ros:"datapath,omitempty"`
	DatapathBridge               string `ros:"datapath.bridge,omitempty"`
	DatapathBridgeCost           string `ros:"datapath.bridge-cost,omitempty"`
	DatapathBridgeHorizon        string `ros:"datapath.bridge-horizon,omitempty"`
	DatapathClientIsolation      string `ros:"datapath.client-isolation,omitempty"`
	DatapathInterfaceList        string `ros:"datapath.interface-list,omitempty"`
	DatapathOpenflowSwitch       string `ros:"datapath.openflow-switch,omitempty"`
	DatapathTrafficProcessing    string `ros:"datapath.traffic-processing,omitempty"`
	DatapathVLANID               string `ros:"datapath.vlan-id,omitempty"`
	DisableRunningCheck          string `ros:"disable-running-check,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Interworking                 string `ros:"interworking,omitempty"`
	// string value, min length 3, max length 3
	InterworkingX3gppInfo        string `ros:"interworking.3gpp-info,omitempty"`
	// hexadecimal string value
	InterworkingX3gppInfoRaw     string `ros:"interworking.3gpp-info-raw,omitempty"`
	InterworkingAuthenticationTypes string `ros:"interworking.authentication-types,omitempty"`
	InterworkingConnectionCapabilities string `ros:"interworking.connection-capabilities,omitempty"`
	// string value
	InterworkingDomainNames      string `ros:"interworking.domain-names,omitempty"`
	InterworkingEsr              string `ros:"interworking.esr,omitempty"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	InterworkingHessid           string `ros:"interworking.hessid,omitempty"`
	InterworkingHotspot20        string `ros:"interworking.hotspot20,omitempty"`
	InterworkingHotspot20Dgaf    string `ros:"interworking.hotspot20-dgaf,omitempty"`
	InterworkingInternet         string `ros:"interworking.internet,omitempty"`
	InterworkingIPv4Availability string `ros:"interworking.ipv4-availability,omitempty"`
	InterworkingIPv6Availability string `ros:"interworking.ipv6-availability,omitempty"`
	InterworkingNetworkType      string `ros:"interworking.network-type,omitempty"`
	InterworkingOperationalClasses string `ros:"interworking.operational-classes,omitempty"`
	// string value
	InterworkingOperatorNames    string `ros:"interworking.operator-names,omitempty"`
	// string value
	InterworkingRealms           string `ros:"interworking.realms,omitempty"`
	// hexadecimal string value
	InterworkingRealmsRaw        string `ros:"interworking.realms-raw,omitempty"`
	// string value
	InterworkingRoamingOis       string `ros:"interworking.roaming-ois,omitempty"`
	InterworkingUesa             string `ros:"interworking.uesa,omitempty"`
	InterworkingVenue            string `ros:"interworking.venue,omitempty"`
	// string value
	InterworkingVenueNames       string `ros:"interworking.venue-names,omitempty"`
	InterworkingWANAtCapacity    string `ros:"interworking.wan-at-capacity,omitempty"`
	InterworkingWANDownlink      string `ros:"interworking.wan-downlink,omitempty"`
	InterworkingWANDownlinkLoad  string `ros:"interworking.wan-downlink-load,omitempty"`
	InterworkingWANMeasurementDuration string `ros:"interworking.wan-measurement-duration,omitempty"`
	InterworkingWANStatus        string `ros:"interworking.wan-status,omitempty"`
	InterworkingWANSymmetric     string `ros:"interworking.wan-symmetric,omitempty"`
	InterworkingWANUplink        string `ros:"interworking.wan-uplink,omitempty"`
	InterworkingWANUplinkLoad    string `ros:"interworking.wan-uplink-load,omitempty"`
	L2MTU                        string `ros:"l2mtu,omitempty"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	MACAddress                   string `ros:"mac-address,omitempty"`
	MasterInterface              string `ros:"master-interface,omitempty"`
	MldInterface                 string `ros:"mld-interface,omitempty"`
	// string value
	MldName                      string `ros:"mld-name,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	RadioMAC                     string `ros:"radio-mac,omitempty"`
	Security                     string `ros:"security,omitempty"`
	// wpa-psk|wpa2-psk|wpa2-psk-sha2|wpa-eap|wpa2-eap|wpa3-psk|owe|wpa3-eap|wpa3-eap-192[,SecurityAuthenticationTypes*]
	SecurityAuthenticationTypes  string `ros:"security.authentication-types,omitempty"`
	SecurityBeaconProtection     string `ros:"security.beacon-protection,omitempty"`
	// string value
	SecurityConnectGroup         string `ros:"security.connect-group,omitempty"`
	SecurityConnectPriority      string `ros:"security.connect-priority,omitempty"`
	// 19|20|21[,SecurityDhGroups*]
	SecurityDhGroups             string `ros:"security.dh-groups,omitempty"`
	SecurityDisablePmkid         string `ros:"security.disable-pmkid,omitempty"`
	SecurityEapAccounting        string `ros:"security.eap-accounting,omitempty"`
	// string value
	SecurityEapAnonymousIdentity string `ros:"security.eap-anonymous-identity,omitempty"`
	SecurityEapCertificateMode   string `ros:"security.eap-certificate-mode,omitempty"`
	SecurityEapMethods           string `ros:"security.eap-methods,omitempty"`
	// string value
	SecurityEapPassword          string `ros:"security.eap-password,omitempty"`
	SecurityEapTLSCertificate    string `ros:"security.eap-tls-certificate,omitempty"`
	// string value
	SecurityEapUsername          string `ros:"security.eap-username,omitempty"`
	// tkip|ccmp|gcmp|ccmp-256|gcmp-256[,SecurityEncryption*]
	SecurityEncryption           string `ros:"security.encryption,omitempty"`
	SecurityFt                   string `ros:"security.ft,omitempty"`
	SecurityFtMobilityDomain     string `ros:"security.ft-mobility-domain,omitempty"`
	// string value, min length 2, max length 96
	SecurityFtNasIdentifier      string `ros:"security.ft-nas-identifier,omitempty"`
	SecurityFtOverDs             string `ros:"security.ft-over-ds,omitempty"`
	SecurityFtPreserveVlanid     string `ros:"security.ft-preserve-vlanid,omitempty"`
	// 1s..6w3d12h15m (time interval)
	SecurityFtR0KeyLifetime      string `ros:"security.ft-r0-key-lifetime,omitempty"`
	// 0s..1m10s (time interval)
	SecurityFtReassociationDeadline string `ros:"security.ft-reassociation-deadline,omitempty"`
	SecurityGroupEncryption      string `ros:"security.group-encryption,omitempty"`
	// 30s..7101w3d6h28m15s (time interval)
	SecurityGroupKeyUpdate       string `ros:"security.group-key-update,omitempty"`
	SecurityManagementEncryption string `ros:"security.management-encryption,omitempty"`
	SecurityManagementProtection string `ros:"security.management-protection,omitempty"`
	SecurityMultiPassphraseGroup string `ros:"security.multi-passphrase-group,omitempty"`
	SecurityOweTransitionInterface string `ros:"security.owe-transition-interface,omitempty"`
	// string value
	SecurityPassphrase           string `ros:"security.passphrase,omitempty"`
	SecuritySaeAntiCloggingThreshold string `ros:"security.sae-anti-clogging-threshold,omitempty"`
	SecuritySaeMaxFailureRate    string `ros:"security.sae-max-failure-rate,omitempty"`
	SecuritySaePwe               string `ros:"security.sae-pwe,omitempty"`
	SecurityWPS                  string `ros:"security.wps,omitempty"`
	Steering                     string `ros:"steering,omitempty"`
	SteeringX2gProbeDelay        string `ros:"steering.2g-probe-delay,omitempty"`
	// string value
	SteeringNeighborGroup        string `ros:"steering.neighbor-group,omitempty"`
	SteeringRrm                  string `ros:"steering.rrm,omitempty"`
	SteeringTransitionRequestCount string `ros:"steering.transition-request-count,omitempty"`
	// time interval
	SteeringTransitionRequestPeriod string `ros:"steering.transition-request-period,omitempty"`
	SteeringTransitionThreshold  string `ros:"steering.transition-threshold,omitempty"`
	// time interval
	SteeringTransitionThresholdTime string `ros:"steering.transition-threshold-time,omitempty"`
	SteeringTransitionTime       string `ros:"steering.transition-time,omitempty"`
	SteeringWnm                  string `ros:"steering.wnm,omitempty"`
}
```

### `/interface/wireguard`

- File: `interface_wireguard.json`
- Sumber field: completion 10

Command: `add`(8), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(8), `set`(8), `wg-export`(3), `wg-import`(3)

```go
// InterfaceWireguardRow adalah kolom hasil `/interface/wireguard/print`.
type InterfaceWireguardRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	// string value
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	ListenPort                   string `ros:"listen-port"`
	MTU                          string `ros:"mtu"`
	// string value
	Name                         string `ros:"name"`
	// string value
	PrivateKey                   string `ros:"private-key"`
	PublicKey                    string `ros:"public-key"`
	Running                      string `ros:"running"`
	VRF                          string `ros:"vrf"`
}
```

```go
// InterfaceWireguardAdd adalah argumen `/interface/wireguard/add`.
type InterfaceWireguardAdd struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	ListenPort                   string `ros:"listen-port,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	// string value
	PrivateKey                   string `ros:"private-key,omitempty"`
	VRF                          string `ros:"vrf,omitempty"`
}
```

```go
// InterfaceWireguardSet adalah argumen `/interface/wireguard/set`.
type InterfaceWireguardSet struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	ListenPort                   string `ros:"listen-port,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	// string value
	PrivateKey                   string `ros:"private-key,omitempty"`
	VRF                          string `ros:"vrf,omitempty"`
}
```

### `/interface/wireless`

- File: `interface_wireless.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command or directory (wireless)

Command: `add`(95), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `frequency-monitor`(6), `get`(4), `monitor`(11), `print`(23), `remove`(1), `reset`(95), `reset-configuration`(1), `reset-mac-address`(1), `scan`(10), `set`(95), `setup-repeater`(15), `spectral-history`(10), `spectral-scan`(12), `wps-client`(14), `wps-push-button`(1)

```go
// InterfaceWirelessAdd adalah argumen `/interface/wireless/add`.
type InterfaceWirelessAdd struct {
	AdaptiveNoiseImmunity        string `ros:"adaptive-noise-immunity,omitempty"`
	AllowSharedkey               string `ros:"allow-sharedkey,omitempty"`
	AmpduPriorities              string `ros:"ampdu-priorities,omitempty"`
	AmsduLimit                   string `ros:"amsdu-limit,omitempty"`
	AmsduThreshold               string `ros:"amsdu-threshold,omitempty"`
	AntennaGain                  string `ros:"antenna-gain,omitempty"`
	AntennaMode                  string `ros:"antenna-mode,omitempty"`
	Area                         string `ros:"area,omitempty"`
	ARP                          string `ros:"arp,omitempty"`
	ARPTimeout                   string `ros:"arp-timeout,omitempty"`
	Band                         string `ros:"band,omitempty"`
	BasicRatesAg                 string `ros:"basic-rates-a/g,omitempty"`
	BasicRatesB                  string `ros:"basic-rates-b,omitempty"`
	BridgeMode                   string `ros:"bridge-mode,omitempty"`
	BurstTime                    string `ros:"burst-time,omitempty"`
	ChannelWidth                 string `ros:"channel-width,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	Compression                  string `ros:"compression,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Country                      string `ros:"country,omitempty"`
	DefaultApTXLimit             string `ros:"default-ap-tx-limit,omitempty"`
	DefaultAuthentication        string `ros:"default-authentication,omitempty"`
	DefaultClientTXLimit         string `ros:"default-client-tx-limit,omitempty"`
	DefaultForwarding            string `ros:"default-forwarding,omitempty"`
	DisableRunningCheck          string `ros:"disable-running-check,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	DisconnectTimeout            string `ros:"disconnect-timeout,omitempty"`
	Distance                     string `ros:"distance,omitempty"`
	FrameLifetime                string `ros:"frame-lifetime,omitempty"`
	Frequency                    string `ros:"frequency,omitempty"`
	FrequencyMode                string `ros:"frequency-mode,omitempty"`
	FrequencyOffset              string `ros:"frequency-offset,omitempty"`
	GuardInterval                string `ros:"guard-interval,omitempty"`
	HideSsid                     string `ros:"hide-ssid,omitempty"`
	HtBasicMcs                   string `ros:"ht-basic-mcs,omitempty"`
	HtSupportedMcs               string `ros:"ht-supported-mcs,omitempty"`
	HWFragmentationThreshold     string `ros:"hw-fragmentation-threshold,omitempty"`
	HWProtectionMode             string `ros:"hw-protection-mode,omitempty"`
	HWProtectionThreshold        string `ros:"hw-protection-threshold,omitempty"`
	HWRetries                    string `ros:"hw-retries,omitempty"`
	Installation                 string `ros:"installation,omitempty"`
	InterworkingProfile          string `ros:"interworking-profile,omitempty"`
	KeepaliveFrames              string `ros:"keepalive-frames,omitempty"`
	L2MTU                        string `ros:"l2mtu,omitempty"`
	MACAddress                   string `ros:"mac-address,omitempty"`
	MasterInterface              string `ros:"master-interface,omitempty"`
	MaxStationCount              string `ros:"max-station-count,omitempty"`
	Mode                         string `ros:"mode,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	MulticastBuffering           string `ros:"multicast-buffering,omitempty"`
	MulticastHelper              string `ros:"multicast-helper,omitempty"`
	Name                         string `ros:"name,omitempty"`
	NoiseFloorThreshold          string `ros:"noise-floor-threshold,omitempty"`
	Nv2CellRadius                string `ros:"nv2-cell-radius,omitempty"`
	Nv2DownlinkRatio             string `ros:"nv2-downlink-ratio,omitempty"`
	Nv2Mode                      string `ros:"nv2-mode,omitempty"`
	Nv2NoiseFloorOffset          string `ros:"nv2-noise-floor-offset,omitempty"`
	Nv2PresharedKey              string `ros:"nv2-preshared-key,omitempty"`
	Nv2Qos                       string `ros:"nv2-qos,omitempty"`
	Nv2QueueCount                string `ros:"nv2-queue-count,omitempty"`
	Nv2Security                  string `ros:"nv2-security,omitempty"`
	Nv2SyncSecret                string `ros:"nv2-sync-secret,omitempty"`
	OnFailRetryTime              string `ros:"on-fail-retry-time,omitempty"`
	PreambleMode                 string `ros:"preamble-mode,omitempty"`
	PrismCardtype                string `ros:"prism-cardtype,omitempty"`
	RadioName                    string `ros:"radio-name,omitempty"`
	RateSelection                string `ros:"rate-selection,omitempty"`
	RateSet                      string `ros:"rate-set,omitempty"`
	RXChains                     string `ros:"rx-chains,omitempty"`
	ScanList                     string `ros:"scan-list,omitempty"`
	SecondaryFrequency           string `ros:"secondary-frequency,omitempty"`
	SecurityProfile              string `ros:"security-profile,omitempty"`
	SkipDfsChannels              string `ros:"skip-dfs-channels,omitempty"`
	Ssid                         string `ros:"ssid,omitempty"`
	StationBridgeCloneMAC        string `ros:"station-bridge-clone-mac,omitempty"`
	StationRoaming               string `ros:"station-roaming,omitempty"`
	SupportedRatesAg             string `ros:"supported-rates-a/g,omitempty"`
	SupportedRatesB              string `ros:"supported-rates-b,omitempty"`
	TdmaPeriodSize               string `ros:"tdma-period-size,omitempty"`
	TXChains                     string `ros:"tx-chains,omitempty"`
	TXPower                      string `ros:"tx-power,omitempty"`
	TXPowerMode                  string `ros:"tx-power-mode,omitempty"`
	UpdateStatsInterval          string `ros:"update-stats-interval,omitempty"`
	VhtBasicMcs                  string `ros:"vht-basic-mcs,omitempty"`
	VhtSupportedMcs              string `ros:"vht-supported-mcs,omitempty"`
	VLANID                       string `ros:"vlan-id,omitempty"`
	VLANMode                     string `ros:"vlan-mode,omitempty"`
	WdsCostRange                 string `ros:"wds-cost-range,omitempty"`
	WdsDefaultBridge             string `ros:"wds-default-bridge,omitempty"`
	WdsDefaultCost               string `ros:"wds-default-cost,omitempty"`
	WdsIgnoreSsid                string `ros:"wds-ignore-ssid,omitempty"`
	WdsMode                      string `ros:"wds-mode,omitempty"`
	WirelessProtocol             string `ros:"wireless-protocol,omitempty"`
	WmmSupport                   string `ros:"wmm-support,omitempty"`
	WPSMode                      string `ros:"wps-mode,omitempty"`
}
```

```go
// InterfaceWirelessSet adalah argumen `/interface/wireless/set`.
type InterfaceWirelessSet struct {
	AdaptiveNoiseImmunity        string `ros:"adaptive-noise-immunity,omitempty"`
	AllowSharedkey               string `ros:"allow-sharedkey,omitempty"`
	AmpduPriorities              string `ros:"ampdu-priorities,omitempty"`
	AmsduLimit                   string `ros:"amsdu-limit,omitempty"`
	AmsduThreshold               string `ros:"amsdu-threshold,omitempty"`
	AntennaGain                  string `ros:"antenna-gain,omitempty"`
	AntennaMode                  string `ros:"antenna-mode,omitempty"`
	Area                         string `ros:"area,omitempty"`
	ARP                          string `ros:"arp,omitempty"`
	ARPTimeout                   string `ros:"arp-timeout,omitempty"`
	Band                         string `ros:"band,omitempty"`
	BasicRatesAg                 string `ros:"basic-rates-a/g,omitempty"`
	BasicRatesB                  string `ros:"basic-rates-b,omitempty"`
	BridgeMode                   string `ros:"bridge-mode,omitempty"`
	BurstTime                    string `ros:"burst-time,omitempty"`
	ChannelWidth                 string `ros:"channel-width,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	Compression                  string `ros:"compression,omitempty"`
	Country                      string `ros:"country,omitempty"`
	DefaultApTXLimit             string `ros:"default-ap-tx-limit,omitempty"`
	DefaultAuthentication        string `ros:"default-authentication,omitempty"`
	DefaultClientTXLimit         string `ros:"default-client-tx-limit,omitempty"`
	DefaultForwarding            string `ros:"default-forwarding,omitempty"`
	DisableRunningCheck          string `ros:"disable-running-check,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	DisconnectTimeout            string `ros:"disconnect-timeout,omitempty"`
	Distance                     string `ros:"distance,omitempty"`
	FrameLifetime                string `ros:"frame-lifetime,omitempty"`
	Frequency                    string `ros:"frequency,omitempty"`
	FrequencyMode                string `ros:"frequency-mode,omitempty"`
	FrequencyOffset              string `ros:"frequency-offset,omitempty"`
	GuardInterval                string `ros:"guard-interval,omitempty"`
	HideSsid                     string `ros:"hide-ssid,omitempty"`
	HtBasicMcs                   string `ros:"ht-basic-mcs,omitempty"`
	HtSupportedMcs               string `ros:"ht-supported-mcs,omitempty"`
	HWFragmentationThreshold     string `ros:"hw-fragmentation-threshold,omitempty"`
	HWProtectionMode             string `ros:"hw-protection-mode,omitempty"`
	HWProtectionThreshold        string `ros:"hw-protection-threshold,omitempty"`
	HWRetries                    string `ros:"hw-retries,omitempty"`
	Installation                 string `ros:"installation,omitempty"`
	InterworkingProfile          string `ros:"interworking-profile,omitempty"`
	KeepaliveFrames              string `ros:"keepalive-frames,omitempty"`
	L2MTU                        string `ros:"l2mtu,omitempty"`
	MACAddress                   string `ros:"mac-address,omitempty"`
	MasterInterface              string `ros:"master-interface,omitempty"`
	MaxStationCount              string `ros:"max-station-count,omitempty"`
	Mode                         string `ros:"mode,omitempty"`
	MTU                          string `ros:"mtu,omitempty"`
	MulticastBuffering           string `ros:"multicast-buffering,omitempty"`
	MulticastHelper              string `ros:"multicast-helper,omitempty"`
	Name                         string `ros:"name,omitempty"`
	NoiseFloorThreshold          string `ros:"noise-floor-threshold,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	Nv2CellRadius                string `ros:"nv2-cell-radius,omitempty"`
	Nv2DownlinkRatio             string `ros:"nv2-downlink-ratio,omitempty"`
	Nv2Mode                      string `ros:"nv2-mode,omitempty"`
	Nv2NoiseFloorOffset          string `ros:"nv2-noise-floor-offset,omitempty"`
	Nv2PresharedKey              string `ros:"nv2-preshared-key,omitempty"`
	Nv2Qos                       string `ros:"nv2-qos,omitempty"`
	Nv2QueueCount                string `ros:"nv2-queue-count,omitempty"`
	Nv2Security                  string `ros:"nv2-security,omitempty"`
	Nv2SyncSecret                string `ros:"nv2-sync-secret,omitempty"`
	OnFailRetryTime              string `ros:"on-fail-retry-time,omitempty"`
	PreambleMode                 string `ros:"preamble-mode,omitempty"`
	PrismCardtype                string `ros:"prism-cardtype,omitempty"`
	RadioName                    string `ros:"radio-name,omitempty"`
	RateSelection                string `ros:"rate-selection,omitempty"`
	RateSet                      string `ros:"rate-set,omitempty"`
	RXChains                     string `ros:"rx-chains,omitempty"`
	ScanList                     string `ros:"scan-list,omitempty"`
	SecondaryFrequency           string `ros:"secondary-frequency,omitempty"`
	SecurityProfile              string `ros:"security-profile,omitempty"`
	SkipDfsChannels              string `ros:"skip-dfs-channels,omitempty"`
	Ssid                         string `ros:"ssid,omitempty"`
	StationBridgeCloneMAC        string `ros:"station-bridge-clone-mac,omitempty"`
	StationRoaming               string `ros:"station-roaming,omitempty"`
	SupportedRatesAg             string `ros:"supported-rates-a/g,omitempty"`
	SupportedRatesB              string `ros:"supported-rates-b,omitempty"`
	TdmaPeriodSize               string `ros:"tdma-period-size,omitempty"`
	TXChains                     string `ros:"tx-chains,omitempty"`
	TXPower                      string `ros:"tx-power,omitempty"`
	TXPowerMode                  string `ros:"tx-power-mode,omitempty"`
	UpdateStatsInterval          string `ros:"update-stats-interval,omitempty"`
	VhtBasicMcs                  string `ros:"vht-basic-mcs,omitempty"`
	VhtSupportedMcs              string `ros:"vht-supported-mcs,omitempty"`
	VLANID                       string `ros:"vlan-id,omitempty"`
	VLANMode                     string `ros:"vlan-mode,omitempty"`
	WdsCostRange                 string `ros:"wds-cost-range,omitempty"`
	WdsDefaultBridge             string `ros:"wds-default-bridge,omitempty"`
	WdsDefaultCost               string `ros:"wds-default-cost,omitempty"`
	WdsIgnoreSsid                string `ros:"wds-ignore-ssid,omitempty"`
	WdsMode                      string `ros:"wds-mode,omitempty"`
	WirelessProtocol             string `ros:"wireless-protocol,omitempty"`
	WmmSupport                   string `ros:"wmm-support,omitempty"`
	WPSMode                      string `ros:"wps-mode,omitempty"`
}
```

### `/ip`

- File: `ip__cmds.json`
- Sumber field: tidak ada kolom terdeteksi

Command: `export`(8), `recursive-print`(0)

_Tidak ada struct: menu ini bukan tabel yang bisa di-print atau tidak tersedia di router uji (mis. paket tidak terpasang)._

### `/ip/address`

- File: `ip_address.json`
- Katalog: `ip-address`
- Sumber field: live 11, completion 13, katalog 7

Command: `add`(8), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(8), `set`(8)

```go
// IPAddressRow adalah kolom hasil `/ip/address/print`.
type IPAddressRow struct {
	ID                           string `ros:".id,key"`
	ActualInterface              string `ros:"actual-interface"`
	// A.B.C.D (IP address)
	Address                      string `ros:"address"`
	// string value
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	Dynamic                      string `ros:"dynamic"`
	Interface                    string `ros:"interface"`
	Invalid                      string `ros:"invalid"`
	// A.B.C.D (IP address)
	Network                      string `ros:"network"`
	Slave                        string `ros:"slave"`
	VRF                          string `ros:"vrf"`
	About                        string `ros:"about"`
	// A.B.C.D (IP address)
	Broadcast                    string `ros:"broadcast"`
	// A.B.C.D (IP address)
	Netmask                      string `ros:"netmask"`
}
```

```go
// IPAddressAdd adalah argumen `/ip/address/add`.
type IPAddressAdd struct {
	// A.B.C.D (IP address)
	Address                      string `ros:"address,required,omitempty"`
	// A.B.C.D (IP address)
	Broadcast                    string `ros:"broadcast,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Interface                    string `ros:"interface,required,omitempty"`
	// A.B.C.D (IP address)
	Netmask                      string `ros:"netmask,omitempty"`
	// A.B.C.D (IP address)
	Network                      string `ros:"network,omitempty"`
}
```

```go
// IPAddressSet adalah argumen `/ip/address/set`.
type IPAddressSet struct {
	// A.B.C.D (IP address)
	Address                      string `ros:"address,omitempty"`
	// A.B.C.D (IP address)
	Broadcast                    string `ros:"broadcast,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	// A.B.C.D (IP address)
	Netmask                      string `ros:"netmask,omitempty"`
	// A.B.C.D (IP address)
	Network                      string `ros:"network,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
}
```

### `/ip/arp`

- File: `ip_arp.json`
- Katalog: `ip-arp`
- Sumber field: live 12, completion 13, katalog 7

Command: `add`(7), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(7), `set`(7)

```go
// IPARPRow adalah kolom hasil `/ip/arp/print`.
type IPARPRow struct {
	ID                           string `ros:".id,key"`
	// A.B.C.D (IP address)
	Address                      string `ros:"address"`
	Complete                     string `ros:"complete"`
	DHCP                         string `ros:"dhcp"`
	Disabled                     string `ros:"disabled"`
	Dynamic                      string `ros:"dynamic"`
	Interface                    string `ros:"interface"`
	Invalid                      string `ros:"invalid"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	MACAddress                   string `ros:"mac-address"`
	Published                    string `ros:"published"`
	Status                       string `ros:"status"`
	VRF                          string `ros:"vrf"`
	About                        string `ros:"about"`
	// string value
	Comment                      string `ros:"comment"`
}
```

```go
// IPARPAdd adalah argumen `/ip/arp/add`.
type IPARPAdd struct {
	// A.B.C.D (IP address)
	Address                      string `ros:"address,required,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Interface                    string `ros:"interface,required,omitempty"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	MACAddress                   string `ros:"mac-address,required,omitempty"`
	Published                    string `ros:"published,omitempty"`
}
```

```go
// IPARPSet adalah argumen `/ip/arp/set`.
type IPARPSet struct {
	// A.B.C.D (IP address)
	Address                      string `ros:"address,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	MACAddress                   string `ros:"mac-address,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	Published                    string `ros:"published,omitempty"`
}
```

### `/ip/cloud`

- File: `ip_cloud.json`
- Sumber field: live 3

Command: 

```go
// IPCloudRow adalah kolom hasil `/ip/cloud/print`.
type IPCloudRow struct {
	DdnsEnabled                  string `ros:"ddns-enabled"`
	DdnsUpdateInterval           string `ros:"ddns-update-interval"`
	UpdateTime                   string `ros:"update-time"`
}
```

### `/ip/dhcp-client`

- File: `ip_dhcp-client.json`
- Katalog: `dhcp-client`
- Sumber field: live 22, completion 34, katalog 10

Command: `add`(17), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(21), `release`(1), `remove`(1), `renew`(1), `reset`(17), `set`(17), `unset`(2)

```go
// IPDhcpclientRow adalah kolom hasil `/ip/dhcp-client/print`.
type IPDhcpclientRow struct {
	ID                           string `ros:".id,key"`
	AddDefaultRoute              string `ros:"add-default-route"`
	Address                      string `ros:"address"`
	AllowReconfigure             string `ros:"allow-reconfigure"`
	CheckGateway                 string `ros:"check-gateway"`
	// string value
	Comment                      string `ros:"comment"`
	DefaultRouteDistance         string `ros:"default-route-distance"`
	DefaultRouteTables           string `ros:"default-route-tables"`
	DHCPOptions                  string `ros:"dhcp-options"`
	DHCPServer                   string `ros:"dhcp-server"`
	Disabled                     string `ros:"disabled"`
	Dynamic                      string `ros:"dynamic"`
	ExpiresAfter                 string `ros:"expires-after"`
	Gateway                      string `ros:"gateway"`
	Interface                    string `ros:"interface"`
	Invalid                      string `ros:"invalid"`
	// string value
	Name                         string `ros:"name"`
	PrimaryDNS                   string `ros:"primary-dns"`
	Status                       string `ros:"status"`
	UseBroadcast                 string `ros:"use-broadcast"`
	UsePeerDNS                   string `ros:"use-peer-dns"`
	UsePeerNTP                   string `ros:"use-peer-ntp"`
	About                        string `ros:"about"`
	CAPsManagers                 string `ros:"caps-managers"`
	CustomHostnameSuffix         string `ros:"custom-hostname-suffix"`
	CustomSourceMACAddress       string `ros:"custom-source-mac-address"`
	Dscp                         string `ros:"dscp"`
	Netmask                      string `ros:"netmask"`
	PrimaryNTP                   string `ros:"primary-ntp"`
	ReconfigureKey               string `ros:"reconfigure-key"`
	ReconfigureLastCounter       string `ros:"reconfigure-last-counter"`
	// string value
	Script                       string `ros:"script"`
	SecondaryDNS                 string `ros:"secondary-dns"`
	SecondaryNTP                 string `ros:"secondary-ntp"`
	VLANPriority                 string `ros:"vlan-priority"`
}
```

```go
// IPDhcpclientAdd adalah argumen `/ip/dhcp-client/add`.
type IPDhcpclientAdd struct {
	AddDefaultRoute              string `ros:"add-default-route,omitempty"`
	AllowReconfigure             string `ros:"allow-reconfigure,omitempty"`
	CheckGateway                 string `ros:"check-gateway,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	DefaultRouteDistance         string `ros:"default-route-distance,omitempty"`
	DefaultRouteTables           string `ros:"default-route-tables,omitempty"`
	DHCPOptions                  string `ros:"dhcp-options,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Dscp                         string `ros:"dscp,omitempty"`
	Interface                    string `ros:"interface,required,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	// string value
	Script                       string `ros:"script,omitempty"`
	UseBroadcast                 string `ros:"use-broadcast,omitempty"`
	UsePeerDNS                   string `ros:"use-peer-dns,omitempty"`
	UsePeerNTP                   string `ros:"use-peer-ntp,omitempty"`
	VLANPriority                 string `ros:"vlan-priority,omitempty"`
}
```

```go
// IPDhcpclientSet adalah argumen `/ip/dhcp-client/set`.
type IPDhcpclientSet struct {
	AddDefaultRoute              string `ros:"add-default-route,omitempty"`
	AllowReconfigure             string `ros:"allow-reconfigure,omitempty"`
	CheckGateway                 string `ros:"check-gateway,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	DefaultRouteDistance         string `ros:"default-route-distance,omitempty"`
	DefaultRouteTables           string `ros:"default-route-tables,omitempty"`
	DHCPOptions                  string `ros:"dhcp-options,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Dscp                         string `ros:"dscp,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	// string value
	Script                       string `ros:"script,omitempty"`
	UseBroadcast                 string `ros:"use-broadcast,omitempty"`
	UsePeerDNS                   string `ros:"use-peer-dns,omitempty"`
	UsePeerNTP                   string `ros:"use-peer-ntp,omitempty"`
	VLANPriority                 string `ros:"vlan-priority,omitempty"`
}
```

### `/ip/dhcp-relay`

- File: `ip_dhcp-relay.json`
- Sumber field: completion 12

Command: `add`(11), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `monitor`(11), `print`(21), `remove`(1), `reset`(11), `reset-counters`(0), `set`(11), `unset`(2)

```go
// IPDHCPRelayRow adalah kolom hasil `/ip/dhcp-relay/print`.
type IPDHCPRelayRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	AddRelayInfo                 string `ros:"add-relay-info"`
	DelayThreshold               string `ros:"delay-threshold"`
	// A.B.C.D (IP address)
	DHCPServer                   string `ros:"dhcp-server"`
	DHCPServerVRF                string `ros:"dhcp-server-vrf"`
	Disabled                     string `ros:"disabled"`
	Interface                    string `ros:"interface"`
	Invalid                      string `ros:"invalid"`
	// A.B.C.D (IP address)
	LocalAddress                 string `ros:"local-address"`
	LocalAddressAsSrcIP          string `ros:"local-address-as-src-ip"`
	// string value
	Name                         string `ros:"name"`
	// string value
	RelayInfoRemoteID            string `ros:"relay-info-remote-id"`
}
```

```go
// IPDHCPRelayAdd adalah argumen `/ip/dhcp-relay/add`.
type IPDHCPRelayAdd struct {
	AddRelayInfo                 string `ros:"add-relay-info,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	DelayThreshold               string `ros:"delay-threshold,omitempty"`
	// A.B.C.D (IP address)
	DHCPServer                   string `ros:"dhcp-server,omitempty"`
	DHCPServerVRF                string `ros:"dhcp-server-vrf,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	// A.B.C.D (IP address)
	LocalAddress                 string `ros:"local-address,omitempty"`
	LocalAddressAsSrcIP          string `ros:"local-address-as-src-ip,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	// string value
	RelayInfoRemoteID            string `ros:"relay-info-remote-id,omitempty"`
}
```

```go
// IPDHCPRelaySet adalah argumen `/ip/dhcp-relay/set`.
type IPDHCPRelaySet struct {
	AddRelayInfo                 string `ros:"add-relay-info,omitempty"`
	DelayThreshold               string `ros:"delay-threshold,omitempty"`
	// A.B.C.D (IP address)
	DHCPServer                   string `ros:"dhcp-server,omitempty"`
	DHCPServerVRF                string `ros:"dhcp-server-vrf,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	// A.B.C.D (IP address)
	LocalAddress                 string `ros:"local-address,omitempty"`
	LocalAddressAsSrcIP          string `ros:"local-address-as-src-ip,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	// string value
	RelayInfoRemoteID            string `ros:"relay-info-remote-id,omitempty"`
}
```

### `/ip/dhcp-server`

- File: `ip_dhcp-server.json`
- Katalog: `dhcp-server`
- Sumber field: live 15, completion 32, katalog 6

Command: `add`(30), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(30), `set`(30), `setup`(0)

```go
// IPDhcpserverRow adalah kolom hasil `/ip/dhcp-server/print`.
type IPDhcpserverRow struct {
	ID                           string `ros:".id,key"`
	// string value, min length 1, max length 64
	AddDNSEntriesSuffix          string `ros:"add-dns-entries-suffix"`
	// string value
	AddressLists                 string `ros:"address-lists"`
	AddressPool                  string `ros:"address-pool"`
	Disabled                     string `ros:"disabled"`
	Dynamic                      string `ros:"dynamic"`
	// client-mac|client-id|opt-82[,DynamicLeaseIdentifiers*]
	DynamicLeaseIdentifiers      string `ros:"dynamic-lease-identifiers"`
	Interface                    string `ros:"interface"`
	Invalid                      string `ros:"invalid"`
	// string value
	LeaseScript                  string `ros:"lease-script"`
	// 1s..7101w3d6h28m15s (time interval)
	LeaseTime                    string `ros:"lease-time"`
	// string value
	Name                         string `ros:"name"`
	SupportBroadbandTr101        string `ros:"support-broadband-tr101"`
	UseRadius                    string `ros:"use-radius"`
	UseReconfigure               string `ros:"use-reconfigure"`
	About                        string `ros:"about"`
	AddARP                       string `ros:"add-arp"`
	AddDNSEntries                string `ros:"add-dns-entries"`
	AllowDualStackQueue          string `ros:"allow-dual-stack-queue"`
	AlwaysBroadcast              string `ros:"always-broadcast"`
	Authoritative                string `ros:"authoritative"`
	BootpLeaseTime               string `ros:"bootp-lease-time"`
	BootpSupport                 string `ros:"bootp-support"`
	ClientMACLimit               string `ros:"client-mac-limit"`
	// string value
	Comment                      string `ros:"comment"`
	ConflictDetection            string `ros:"conflict-detection"`
	DelayThreshold               string `ros:"delay-threshold"`
	DHCPOptionSet                string `ros:"dhcp-option-set"`
	InsertQueueBefore            string `ros:"insert-queue-before"`
	ParentQueue                  string `ros:"parent-queue"`
	// A.B.C.D (IP address)
	Relay                        string `ros:"relay"`
	// A.B.C.D (IP address)
	ServerAddress                string `ros:"server-address"`
	UseFramedAsClassless         string `ros:"use-framed-as-classless"`
}
```

```go
// IPDhcpserverAdd adalah argumen `/ip/dhcp-server/add`.
type IPDhcpserverAdd struct {
	AddARP                       string `ros:"add-arp,omitempty"`
	AddDNSEntries                string `ros:"add-dns-entries,omitempty"`
	// string value, min length 1, max length 64
	AddDNSEntriesSuffix          string `ros:"add-dns-entries-suffix,omitempty"`
	// string value
	AddressLists                 string `ros:"address-lists,omitempty"`
	AddressPool                  string `ros:"address-pool,omitempty"`
	AllowDualStackQueue          string `ros:"allow-dual-stack-queue,omitempty"`
	AlwaysBroadcast              string `ros:"always-broadcast,omitempty"`
	Authoritative                string `ros:"authoritative,omitempty"`
	BootpLeaseTime               string `ros:"bootp-lease-time,omitempty"`
	BootpSupport                 string `ros:"bootp-support,omitempty"`
	ClientMACLimit               string `ros:"client-mac-limit,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	ConflictDetection            string `ros:"conflict-detection,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	DelayThreshold               string `ros:"delay-threshold,omitempty"`
	DHCPOptionSet                string `ros:"dhcp-option-set,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	// client-mac|client-id|opt-82[,DynamicLeaseIdentifiers*]
	DynamicLeaseIdentifiers      string `ros:"dynamic-lease-identifiers,omitempty"`
	InsertQueueBefore            string `ros:"insert-queue-before,omitempty"`
	Interface                    string `ros:"interface,required,omitempty"`
	// string value
	LeaseScript                  string `ros:"lease-script,omitempty"`
	// 1s..7101w3d6h28m15s (time interval)
	LeaseTime                    string `ros:"lease-time,omitempty"`
	// string value
	Name                         string `ros:"name,required,omitempty"`
	ParentQueue                  string `ros:"parent-queue,omitempty"`
	// A.B.C.D (IP address)
	Relay                        string `ros:"relay,omitempty"`
	// A.B.C.D (IP address)
	ServerAddress                string `ros:"server-address,omitempty"`
	SupportBroadbandTr101        string `ros:"support-broadband-tr101,omitempty"`
	UseFramedAsClassless         string `ros:"use-framed-as-classless,omitempty"`
	UseRadius                    string `ros:"use-radius,omitempty"`
	UseReconfigure               string `ros:"use-reconfigure,omitempty"`
}
```

```go
// IPDhcpserverSet adalah argumen `/ip/dhcp-server/set`.
type IPDhcpserverSet struct {
	AddARP                       string `ros:"add-arp,omitempty"`
	AddDNSEntries                string `ros:"add-dns-entries,omitempty"`
	// string value, min length 1, max length 64
	AddDNSEntriesSuffix          string `ros:"add-dns-entries-suffix,omitempty"`
	// string value
	AddressLists                 string `ros:"address-lists,omitempty"`
	AddressPool                  string `ros:"address-pool,omitempty"`
	AllowDualStackQueue          string `ros:"allow-dual-stack-queue,omitempty"`
	AlwaysBroadcast              string `ros:"always-broadcast,omitempty"`
	Authoritative                string `ros:"authoritative,omitempty"`
	BootpLeaseTime               string `ros:"bootp-lease-time,omitempty"`
	BootpSupport                 string `ros:"bootp-support,omitempty"`
	ClientMACLimit               string `ros:"client-mac-limit,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	ConflictDetection            string `ros:"conflict-detection,omitempty"`
	DelayThreshold               string `ros:"delay-threshold,omitempty"`
	DHCPOptionSet                string `ros:"dhcp-option-set,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	// client-mac|client-id|opt-82[,DynamicLeaseIdentifiers*]
	DynamicLeaseIdentifiers      string `ros:"dynamic-lease-identifiers,omitempty"`
	InsertQueueBefore            string `ros:"insert-queue-before,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	// string value
	LeaseScript                  string `ros:"lease-script,omitempty"`
	// 1s..7101w3d6h28m15s (time interval)
	LeaseTime                    string `ros:"lease-time,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	ParentQueue                  string `ros:"parent-queue,omitempty"`
	// A.B.C.D (IP address)
	Relay                        string `ros:"relay,omitempty"`
	// A.B.C.D (IP address)
	ServerAddress                string `ros:"server-address,omitempty"`
	SupportBroadbandTr101        string `ros:"support-broadband-tr101,omitempty"`
	UseFramedAsClassless         string `ros:"use-framed-as-classless,omitempty"`
	UseRadius                    string `ros:"use-radius,omitempty"`
	UseReconfigure               string `ros:"use-reconfigure,omitempty"`
}
```

### `/ip/dhcp-server/lease`

- File: `(dari katalog Go)`
- Katalog: `dhcp-lease`
- Sumber field: katalog 20

Command: `add`(7), `set`(7)

```go
// IPDhcpserverLeaseRow adalah kolom hasil `/ip/dhcp-server/lease/print`.
type IPDhcpserverLeaseRow struct {
	ID                           string `ros:".id,key"`
	Address                      string `ros:"address"`
	MACAddress                   string `ros:"mac-address"`
	ClientID                     string `ros:"client-id"`
	HostName                     string `ros:"host-name"`
	Status                       string `ros:"status"`
	Server                       string `ros:"server"`
	ActiveAddress                string `ros:"active-address"`
	ActiveMACAddress             string `ros:"active-mac-address"`
	ExpiresAfter                 string `ros:"expires-after"`
	LastSeen                     string `ros:"last-seen"`
	Dynamic                      string `ros:"dynamic"`
	Disabled                     string `ros:"disabled"`
	BlockAccess                  string `ros:"block-access"`
	Comment                      string `ros:"comment"`
}
```

```go
// IPDhcpserverLeaseAdd adalah argumen `/ip/dhcp-server/lease/add`.
type IPDhcpserverLeaseAdd struct {
	Address                      string `ros:"address,required,omitempty"`
	MACAddress                   string `ros:"mac-address,required,omitempty"`
	Server                       string `ros:"server,omitempty"`
	ClientID                     string `ros:"client-id,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	BlockAccess                  string `ros:"block-access,omitempty"`
	RateLimit                    string `ros:"rate-limit,omitempty"`
}
```

```go
// IPDhcpserverLeaseSet adalah argumen `/ip/dhcp-server/lease/set`.
type IPDhcpserverLeaseSet struct {
	Address                      string `ros:"address,omitempty"`
	MACAddress                   string `ros:"mac-address,omitempty"`
	Server                       string `ros:"server,omitempty"`
	ClientID                     string `ros:"client-id,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	BlockAccess                  string `ros:"block-access,omitempty"`
	RateLimit                    string `ros:"rate-limit,omitempty"`
}
```

### `/ip/dhcp-server/network`

- File: `(dari katalog Go)`
- Katalog: `dhcp-server-network`
- Sumber field: katalog 5

Command: `add`(5), `set`(5)

```go
// IPDhcpserverNetworkRow adalah kolom hasil `/ip/dhcp-server/network/print`.
type IPDhcpserverNetworkRow struct {
	ID                           string `ros:".id,key"`
	Address                      string `ros:"address"`
	Gateway                      string `ros:"gateway"`
	DNSServer                    string `ros:"dns-server"`
	Domain                       string `ros:"domain"`
	Comment                      string `ros:"comment"`
}
```

```go
// IPDhcpserverNetworkAdd adalah argumen `/ip/dhcp-server/network/add`.
type IPDhcpserverNetworkAdd struct {
	Address                      string `ros:"address,required,omitempty"`
	Gateway                      string `ros:"gateway,omitempty"`
	DNSServer                    string `ros:"dns-server,omitempty"`
	Domain                       string `ros:"domain,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

```go
// IPDhcpserverNetworkSet adalah argumen `/ip/dhcp-server/network/set`.
type IPDhcpserverNetworkSet struct {
	Address                      string `ros:"address,omitempty"`
	Gateway                      string `ros:"gateway,omitempty"`
	DNSServer                    string `ros:"dns-server,omitempty"`
	Domain                       string `ros:"domain,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

### `/ip/dns`

- File: `ip_dns.json`
- Katalog: `ip-dns`
- Sumber field: live 19, katalog 5

Command: `edit`(1), `export`(8), `get`(3), `print`(5), `reset`(17), `set`(17)

```go
// IPDNSRow adalah kolom hasil `/ip/dns/print`.
type IPDNSRow struct {
	ID                           string `ros:".id,key"`
	// time interval
	AddressListExtraTime         string `ros:"address-list-extra-time"`
	AllowRemoteRequests          string `ros:"allow-remote-requests"`
	// time interval
	CacheMaxTTL                  string `ros:"cache-max-ttl"`
	CacheSize                    string `ros:"cache-size"`
	CacheUsed                    string `ros:"cache-used"`
	DohMaxConcurrentQueries      string `ros:"doh-max-concurrent-queries"`
	DohMaxServerConnections      string `ros:"doh-max-server-connections"`
	// time interval
	DohTimeout                   string `ros:"doh-timeout"`
	DynamicServers               string `ros:"dynamic-servers"`
	MaxConcurrentQueries         string `ros:"max-concurrent-queries"`
	MaxConcurrentTcpSessions     string `ros:"max-concurrent-tcp-sessions"`
	MaxUdpPacketSize             string `ros:"max-udp-packet-size"`
	MdnsRepeatIfaces             string `ros:"mdns-repeat-ifaces"`
	// time interval
	QueryServerTimeout           string `ros:"query-server-timeout"`
	// time interval
	QueryTotalTimeout            string `ros:"query-total-timeout"`
	Servers                      string `ros:"servers"`
	// string value
	UseDohServer                 string `ros:"use-doh-server"`
	VerifyDohCert                string `ros:"verify-doh-cert"`
	VRF                          string `ros:"vrf"`
}
```

```go
// IPDNSSet adalah argumen `/ip/dns/set`.
type IPDNSSet struct {
	// time interval
	AddressListExtraTime         string `ros:"address-list-extra-time,omitempty"`
	AllowRemoteRequests          string `ros:"allow-remote-requests,omitempty"`
	// time interval
	CacheMaxTTL                  string `ros:"cache-max-ttl,omitempty"`
	CacheSize                    string `ros:"cache-size,omitempty"`
	DohMaxConcurrentQueries      string `ros:"doh-max-concurrent-queries,omitempty"`
	DohMaxServerConnections      string `ros:"doh-max-server-connections,omitempty"`
	// time interval
	DohTimeout                   string `ros:"doh-timeout,omitempty"`
	MaxConcurrentQueries         string `ros:"max-concurrent-queries,omitempty"`
	MaxConcurrentTcpSessions     string `ros:"max-concurrent-tcp-sessions,omitempty"`
	MaxUdpPacketSize             string `ros:"max-udp-packet-size,omitempty"`
	MdnsRepeatIfaces             string `ros:"mdns-repeat-ifaces,omitempty"`
	// time interval
	QueryServerTimeout           string `ros:"query-server-timeout,omitempty"`
	// time interval
	QueryTotalTimeout            string `ros:"query-total-timeout,omitempty"`
	Servers                      string `ros:"servers,omitempty"`
	// string value
	UseDohServer                 string `ros:"use-doh-server,omitempty"`
	VerifyDohCert                string `ros:"verify-doh-cert,omitempty"`
	VRF                          string `ros:"vrf,omitempty"`
}
```

### `/ip/dns/cache`

- File: `(dari katalog Go)`
- Katalog: `ip-dns-cache`
- Sumber field: katalog 4

Command: 

```go
// IPDNSCacheRow adalah kolom hasil `/ip/dns/cache/print`.
type IPDNSCacheRow struct {
	ID                           string `ros:".id,key"`
	Name                         string `ros:"name"`
	Type                         string `ros:"type"`
	Data                         string `ros:"data"`
	TTL                          string `ros:"ttl"`
}
```

### `/ip/dns/static`

- File: `(dari katalog Go)`
- Katalog: `ip-dns-static`
- Sumber field: katalog 5

Command: `add`(4), `set`(4)

```go
// IPDNSStaticRow adalah kolom hasil `/ip/dns/static/print`.
type IPDNSStaticRow struct {
	ID                           string `ros:".id,key"`
	Name                         string `ros:"name"`
	Address                      string `ros:"address"`
	TTL                          string `ros:"ttl"`
	Disabled                     string `ros:"disabled"`
	Comment                      string `ros:"comment"`
}
```

```go
// IPDNSStaticAdd adalah argumen `/ip/dns/static/add`.
type IPDNSStaticAdd struct {
	Name                         string `ros:"name,required,omitempty"`
	Address                      string `ros:"address,required,omitempty"`
	TTL                          string `ros:"ttl,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

```go
// IPDNSStaticSet adalah argumen `/ip/dns/static/set`.
type IPDNSStaticSet struct {
	Name                         string `ros:"name,omitempty"`
	Address                      string `ros:"address,omitempty"`
	TTL                          string `ros:"ttl,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

### `/ip/firewall`

- File: `ip_firewall.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command

Command: `export`(8)

_Tidak ada struct: menu ini bukan tabel yang bisa di-print atau tidak tersedia di router uji (mis. paket tidak terpasang)._

### `/ip/firewall/address-list`

- File: `(dari katalog Go)`
- Katalog: `fw-address-list`
- Sumber field: katalog 7

Command: `add`(4), `set`(4)

```go
// IPFirewallAddresslistRow adalah kolom hasil `/ip/firewall/address-list/print`.
type IPFirewallAddresslistRow struct {
	ID                           string `ros:".id,key"`
	List                         string `ros:"list"`
	Address                      string `ros:"address"`
	Timeout                      string `ros:"timeout"`
	CreationTime                 string `ros:"creation-time"`
	Dynamic                      string `ros:"dynamic"`
	Disabled                     string `ros:"disabled"`
	Comment                      string `ros:"comment"`
}
```

```go
// IPFirewallAddresslistAdd adalah argumen `/ip/firewall/address-list/add`.
type IPFirewallAddresslistAdd struct {
	List                         string `ros:"list,required,omitempty"`
	Address                      string `ros:"address,required,omitempty"`
	Timeout                      string `ros:"timeout,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

```go
// IPFirewallAddresslistSet adalah argumen `/ip/firewall/address-list/set`.
type IPFirewallAddresslistSet struct {
	List                         string `ros:"list,omitempty"`
	Address                      string `ros:"address,omitempty"`
	Timeout                      string `ros:"timeout,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

### `/ip/firewall/connection`

- File: `(dari katalog Go)`
- Katalog: `fw-connection`
- Sumber field: katalog 9

Command: 

```go
// IPFirewallConnectionRow adalah kolom hasil `/ip/firewall/connection/print`.
type IPFirewallConnectionRow struct {
	ID                           string `ros:".id,key"`
	Protocol                     string `ros:"protocol"`
	SrcAddress                   string `ros:"src-address"`
	DstAddress                   string `ros:"dst-address"`
	ReplySrcAddress              string `ros:"reply-src-address"`
	ReplyDstAddress              string `ros:"reply-dst-address"`
	TcpState                     string `ros:"tcp-state"`
	Timeout                      string `ros:"timeout"`
	OrigBytes                    string `ros:"orig-bytes"`
	ReplBytes                    string `ros:"repl-bytes"`
}
```

### `/ip/firewall/filter`

- File: `(dari katalog Go)`
- Katalog: `fw-filter`
- Sumber field: katalog 23

Command: `add`(18), `set`(18)

```go
// IPFirewallFilterRow adalah kolom hasil `/ip/firewall/filter/print`.
type IPFirewallFilterRow struct {
	ID                           string `ros:".id,key"`
	Chain                        string `ros:"chain"`
	Action                       string `ros:"action"`
	Protocol                     string `ros:"protocol"`
	SrcAddress                   string `ros:"src-address"`
	DstAddress                   string `ros:"dst-address"`
	SrcPort                      string `ros:"src-port"`
	DstPort                      string `ros:"dst-port"`
	InInterface                  string `ros:"in-interface"`
	OutInterface                 string `ros:"out-interface"`
	InInterfaceList              string `ros:"in-interface-list"`
	OutInterfaceList             string `ros:"out-interface-list"`
	SrcAddressList               string `ros:"src-address-list"`
	DstAddressList               string `ros:"dst-address-list"`
	ConnectionState              string `ros:"connection-state"`
	Log                          string `ros:"log"`
	Disabled                     string `ros:"disabled"`
	Comment                      string `ros:"comment"`
	Bytes                        string `ros:"bytes"`
	Packets                      string `ros:"packets"`
}
```

```go
// IPFirewallFilterAdd adalah argumen `/ip/firewall/filter/add`.
type IPFirewallFilterAdd struct {
	Chain                        string `ros:"chain,required,omitempty"`
	Action                       string `ros:"action,required,omitempty"`
	Protocol                     string `ros:"protocol,omitempty"`
	SrcAddress                   string `ros:"src-address,omitempty"`
	DstAddress                   string `ros:"dst-address,omitempty"`
	SrcPort                      string `ros:"src-port,omitempty"`
	DstPort                      string `ros:"dst-port,omitempty"`
	InInterface                  string `ros:"in-interface,omitempty"`
	OutInterface                 string `ros:"out-interface,omitempty"`
	InInterfaceList              string `ros:"in-interface-list,omitempty"`
	OutInterfaceList             string `ros:"out-interface-list,omitempty"`
	SrcAddressList               string `ros:"src-address-list,omitempty"`
	DstAddressList               string `ros:"dst-address-list,omitempty"`
	ConnectionState              string `ros:"connection-state,omitempty"`
	Log                          string `ros:"log,omitempty"`
	LogPrefix                    string `ros:"log-prefix,omitempty"`
	PlaceBefore                  string `ros:"place-before,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

```go
// IPFirewallFilterSet adalah argumen `/ip/firewall/filter/set`.
type IPFirewallFilterSet struct {
	Chain                        string `ros:"chain,omitempty"`
	Action                       string `ros:"action,omitempty"`
	Protocol                     string `ros:"protocol,omitempty"`
	SrcAddress                   string `ros:"src-address,omitempty"`
	DstAddress                   string `ros:"dst-address,omitempty"`
	SrcPort                      string `ros:"src-port,omitempty"`
	DstPort                      string `ros:"dst-port,omitempty"`
	InInterface                  string `ros:"in-interface,omitempty"`
	OutInterface                 string `ros:"out-interface,omitempty"`
	InInterfaceList              string `ros:"in-interface-list,omitempty"`
	OutInterfaceList             string `ros:"out-interface-list,omitempty"`
	SrcAddressList               string `ros:"src-address-list,omitempty"`
	DstAddressList               string `ros:"dst-address-list,omitempty"`
	ConnectionState              string `ros:"connection-state,omitempty"`
	Log                          string `ros:"log,omitempty"`
	LogPrefix                    string `ros:"log-prefix,omitempty"`
	PlaceBefore                  string `ros:"place-before,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

### `/ip/firewall/mangle`

- File: `(dari katalog Go)`
- Katalog: `fw-mangle`
- Sumber field: katalog 29

Command: `add`(24), `set`(24)

```go
// IPFirewallMangleRow adalah kolom hasil `/ip/firewall/mangle/print`.
type IPFirewallMangleRow struct {
	ID                           string `ros:".id,key"`
	Chain                        string `ros:"chain"`
	Action                       string `ros:"action"`
	Protocol                     string `ros:"protocol"`
	SrcAddress                   string `ros:"src-address"`
	DstAddress                   string `ros:"dst-address"`
	SrcPort                      string `ros:"src-port"`
	DstPort                      string `ros:"dst-port"`
	InInterface                  string `ros:"in-interface"`
	OutInterface                 string `ros:"out-interface"`
	InInterfaceList              string `ros:"in-interface-list"`
	OutInterfaceList             string `ros:"out-interface-list"`
	SrcAddressList               string `ros:"src-address-list"`
	DstAddressList               string `ros:"dst-address-list"`
	ConnectionState              string `ros:"connection-state"`
	ConnectionMark               string `ros:"connection-mark"`
	PacketMark                   string `ros:"packet-mark"`
	NewConnectionMark            string `ros:"new-connection-mark"`
	NewPacketMark                string `ros:"new-packet-mark"`
	NewRoutingMark               string `ros:"new-routing-mark"`
	Passthrough                  string `ros:"passthrough"`
	Log                          string `ros:"log"`
	Disabled                     string `ros:"disabled"`
	Comment                      string `ros:"comment"`
	Bytes                        string `ros:"bytes"`
	Packets                      string `ros:"packets"`
}
```

```go
// IPFirewallMangleAdd adalah argumen `/ip/firewall/mangle/add`.
type IPFirewallMangleAdd struct {
	Chain                        string `ros:"chain,required,omitempty"`
	Action                       string `ros:"action,required,omitempty"`
	Protocol                     string `ros:"protocol,omitempty"`
	SrcAddress                   string `ros:"src-address,omitempty"`
	DstAddress                   string `ros:"dst-address,omitempty"`
	SrcPort                      string `ros:"src-port,omitempty"`
	DstPort                      string `ros:"dst-port,omitempty"`
	InInterface                  string `ros:"in-interface,omitempty"`
	OutInterface                 string `ros:"out-interface,omitempty"`
	InInterfaceList              string `ros:"in-interface-list,omitempty"`
	OutInterfaceList             string `ros:"out-interface-list,omitempty"`
	SrcAddressList               string `ros:"src-address-list,omitempty"`
	DstAddressList               string `ros:"dst-address-list,omitempty"`
	ConnectionState              string `ros:"connection-state,omitempty"`
	ConnectionMark               string `ros:"connection-mark,omitempty"`
	PacketMark                   string `ros:"packet-mark,omitempty"`
	NewConnectionMark            string `ros:"new-connection-mark,omitempty"`
	NewPacketMark                string `ros:"new-packet-mark,omitempty"`
	NewRoutingMark               string `ros:"new-routing-mark,omitempty"`
	Passthrough                  string `ros:"passthrough,omitempty"`
	Log                          string `ros:"log,omitempty"`
	LogPrefix                    string `ros:"log-prefix,omitempty"`
	PlaceBefore                  string `ros:"place-before,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

```go
// IPFirewallMangleSet adalah argumen `/ip/firewall/mangle/set`.
type IPFirewallMangleSet struct {
	Chain                        string `ros:"chain,omitempty"`
	Action                       string `ros:"action,omitempty"`
	Protocol                     string `ros:"protocol,omitempty"`
	SrcAddress                   string `ros:"src-address,omitempty"`
	DstAddress                   string `ros:"dst-address,omitempty"`
	SrcPort                      string `ros:"src-port,omitempty"`
	DstPort                      string `ros:"dst-port,omitempty"`
	InInterface                  string `ros:"in-interface,omitempty"`
	OutInterface                 string `ros:"out-interface,omitempty"`
	InInterfaceList              string `ros:"in-interface-list,omitempty"`
	OutInterfaceList             string `ros:"out-interface-list,omitempty"`
	SrcAddressList               string `ros:"src-address-list,omitempty"`
	DstAddressList               string `ros:"dst-address-list,omitempty"`
	ConnectionState              string `ros:"connection-state,omitempty"`
	ConnectionMark               string `ros:"connection-mark,omitempty"`
	PacketMark                   string `ros:"packet-mark,omitempty"`
	NewConnectionMark            string `ros:"new-connection-mark,omitempty"`
	NewPacketMark                string `ros:"new-packet-mark,omitempty"`
	NewRoutingMark               string `ros:"new-routing-mark,omitempty"`
	Passthrough                  string `ros:"passthrough,omitempty"`
	Log                          string `ros:"log,omitempty"`
	LogPrefix                    string `ros:"log-prefix,omitempty"`
	PlaceBefore                  string `ros:"place-before,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

### `/ip/firewall/nat`

- File: `(dari katalog Go)`
- Katalog: `fw-nat`
- Sumber field: katalog 25

Command: `add`(20), `set`(20)

```go
// IPFirewallNATRow adalah kolom hasil `/ip/firewall/nat/print`.
type IPFirewallNATRow struct {
	ID                           string `ros:".id,key"`
	Chain                        string `ros:"chain"`
	Action                       string `ros:"action"`
	Protocol                     string `ros:"protocol"`
	SrcAddress                   string `ros:"src-address"`
	DstAddress                   string `ros:"dst-address"`
	SrcPort                      string `ros:"src-port"`
	DstPort                      string `ros:"dst-port"`
	InInterface                  string `ros:"in-interface"`
	OutInterface                 string `ros:"out-interface"`
	InInterfaceList              string `ros:"in-interface-list"`
	OutInterfaceList             string `ros:"out-interface-list"`
	SrcAddressList               string `ros:"src-address-list"`
	DstAddressList               string `ros:"dst-address-list"`
	ConnectionState              string `ros:"connection-state"`
	ToAddresses                  string `ros:"to-addresses"`
	ToPorts                      string `ros:"to-ports"`
	Log                          string `ros:"log"`
	Disabled                     string `ros:"disabled"`
	Comment                      string `ros:"comment"`
	Bytes                        string `ros:"bytes"`
	Packets                      string `ros:"packets"`
}
```

```go
// IPFirewallNATAdd adalah argumen `/ip/firewall/nat/add`.
type IPFirewallNATAdd struct {
	Chain                        string `ros:"chain,required,omitempty"`
	Action                       string `ros:"action,required,omitempty"`
	Protocol                     string `ros:"protocol,omitempty"`
	SrcAddress                   string `ros:"src-address,omitempty"`
	DstAddress                   string `ros:"dst-address,omitempty"`
	SrcPort                      string `ros:"src-port,omitempty"`
	DstPort                      string `ros:"dst-port,omitempty"`
	InInterface                  string `ros:"in-interface,omitempty"`
	OutInterface                 string `ros:"out-interface,omitempty"`
	InInterfaceList              string `ros:"in-interface-list,omitempty"`
	OutInterfaceList             string `ros:"out-interface-list,omitempty"`
	SrcAddressList               string `ros:"src-address-list,omitempty"`
	DstAddressList               string `ros:"dst-address-list,omitempty"`
	ConnectionState              string `ros:"connection-state,omitempty"`
	ToAddresses                  string `ros:"to-addresses,omitempty"`
	ToPorts                      string `ros:"to-ports,omitempty"`
	Log                          string `ros:"log,omitempty"`
	LogPrefix                    string `ros:"log-prefix,omitempty"`
	PlaceBefore                  string `ros:"place-before,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

```go
// IPFirewallNATSet adalah argumen `/ip/firewall/nat/set`.
type IPFirewallNATSet struct {
	Chain                        string `ros:"chain,omitempty"`
	Action                       string `ros:"action,omitempty"`
	Protocol                     string `ros:"protocol,omitempty"`
	SrcAddress                   string `ros:"src-address,omitempty"`
	DstAddress                   string `ros:"dst-address,omitempty"`
	SrcPort                      string `ros:"src-port,omitempty"`
	DstPort                      string `ros:"dst-port,omitempty"`
	InInterface                  string `ros:"in-interface,omitempty"`
	OutInterface                 string `ros:"out-interface,omitempty"`
	InInterfaceList              string `ros:"in-interface-list,omitempty"`
	OutInterfaceList             string `ros:"out-interface-list,omitempty"`
	SrcAddressList               string `ros:"src-address-list,omitempty"`
	DstAddressList               string `ros:"dst-address-list,omitempty"`
	ConnectionState              string `ros:"connection-state,omitempty"`
	ToAddresses                  string `ros:"to-addresses,omitempty"`
	ToPorts                      string `ros:"to-ports,omitempty"`
	Log                          string `ros:"log,omitempty"`
	LogPrefix                    string `ros:"log-prefix,omitempty"`
	PlaceBefore                  string `ros:"place-before,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

### `/ip/firewall/raw`

- File: `(dari katalog Go)`
- Katalog: `fw-raw`
- Sumber field: katalog 23

Command: `add`(18), `set`(18)

```go
// IPFirewallRawRow adalah kolom hasil `/ip/firewall/raw/print`.
type IPFirewallRawRow struct {
	ID                           string `ros:".id,key"`
	Chain                        string `ros:"chain"`
	Action                       string `ros:"action"`
	Protocol                     string `ros:"protocol"`
	SrcAddress                   string `ros:"src-address"`
	DstAddress                   string `ros:"dst-address"`
	SrcPort                      string `ros:"src-port"`
	DstPort                      string `ros:"dst-port"`
	InInterface                  string `ros:"in-interface"`
	OutInterface                 string `ros:"out-interface"`
	InInterfaceList              string `ros:"in-interface-list"`
	OutInterfaceList             string `ros:"out-interface-list"`
	SrcAddressList               string `ros:"src-address-list"`
	DstAddressList               string `ros:"dst-address-list"`
	ConnectionState              string `ros:"connection-state"`
	Log                          string `ros:"log"`
	Disabled                     string `ros:"disabled"`
	Comment                      string `ros:"comment"`
	Bytes                        string `ros:"bytes"`
	Packets                      string `ros:"packets"`
}
```

```go
// IPFirewallRawAdd adalah argumen `/ip/firewall/raw/add`.
type IPFirewallRawAdd struct {
	Chain                        string `ros:"chain,required,omitempty"`
	Action                       string `ros:"action,required,omitempty"`
	Protocol                     string `ros:"protocol,omitempty"`
	SrcAddress                   string `ros:"src-address,omitempty"`
	DstAddress                   string `ros:"dst-address,omitempty"`
	SrcPort                      string `ros:"src-port,omitempty"`
	DstPort                      string `ros:"dst-port,omitempty"`
	InInterface                  string `ros:"in-interface,omitempty"`
	OutInterface                 string `ros:"out-interface,omitempty"`
	InInterfaceList              string `ros:"in-interface-list,omitempty"`
	OutInterfaceList             string `ros:"out-interface-list,omitempty"`
	SrcAddressList               string `ros:"src-address-list,omitempty"`
	DstAddressList               string `ros:"dst-address-list,omitempty"`
	ConnectionState              string `ros:"connection-state,omitempty"`
	Log                          string `ros:"log,omitempty"`
	LogPrefix                    string `ros:"log-prefix,omitempty"`
	PlaceBefore                  string `ros:"place-before,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

```go
// IPFirewallRawSet adalah argumen `/ip/firewall/raw/set`.
type IPFirewallRawSet struct {
	Chain                        string `ros:"chain,omitempty"`
	Action                       string `ros:"action,omitempty"`
	Protocol                     string `ros:"protocol,omitempty"`
	SrcAddress                   string `ros:"src-address,omitempty"`
	DstAddress                   string `ros:"dst-address,omitempty"`
	SrcPort                      string `ros:"src-port,omitempty"`
	DstPort                      string `ros:"dst-port,omitempty"`
	InInterface                  string `ros:"in-interface,omitempty"`
	OutInterface                 string `ros:"out-interface,omitempty"`
	InInterfaceList              string `ros:"in-interface-list,omitempty"`
	OutInterfaceList             string `ros:"out-interface-list,omitempty"`
	SrcAddressList               string `ros:"src-address-list,omitempty"`
	DstAddressList               string `ros:"dst-address-list,omitempty"`
	ConnectionState              string `ros:"connection-state,omitempty"`
	Log                          string `ros:"log,omitempty"`
	LogPrefix                    string `ros:"log-prefix,omitempty"`
	PlaceBefore                  string `ros:"place-before,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

### `/ip/hotspot`

- File: `ip_hotspot.json`
- Katalog: `hotspot-server`
- Sumber field: live 14, completion 14, katalog 10

Command: `add`(10), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(10), `reset-html`(1), `set`(10), `setup`(0)

```go
// IPHotspotRow adalah kolom hasil `/ip/hotspot/print`.
type IPHotspotRow struct {
	ID                           string `ros:".id,key"`
	HTTPS                        string `ros:"HTTPS"`
	AddressPool                  string `ros:"address-pool"`
	AddressesPerMAC              string `ros:"addresses-per-mac"`
	Disabled                     string `ros:"disabled"`
	IdleTimeout                  string `ros:"idle-timeout"`
	Interface                    string `ros:"interface"`
	Invalid                      string `ros:"invalid"`
	IPOfDNSName                  string `ros:"ip-of-dns-name"`
	KeepaliveTimeout             string `ros:"keepalive-timeout"`
	LoginTimeout                 string `ros:"login-timeout"`
	// string value
	Name                         string `ros:"name"`
	Profile                      string `ros:"profile"`
	ProxyStatus                  string `ros:"proxy-status"`
	About                        string `ros:"about"`
	Comment                      string `ros:"comment"`
}
```

```go
// IPHotspotAdd adalah argumen `/ip/hotspot/add`.
type IPHotspotAdd struct {
	AddressPool                  string `ros:"address-pool,omitempty"`
	AddressesPerMAC              string `ros:"addresses-per-mac,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	IdleTimeout                  string `ros:"idle-timeout,omitempty"`
	Interface                    string `ros:"interface,required,omitempty"`
	KeepaliveTimeout             string `ros:"keepalive-timeout,omitempty"`
	LoginTimeout                 string `ros:"login-timeout,omitempty"`
	// string value
	Name                         string `ros:"name,required,omitempty"`
	Profile                      string `ros:"profile,omitempty"`
}
```

```go
// IPHotspotSet adalah argumen `/ip/hotspot/set`.
type IPHotspotSet struct {
	AddressPool                  string `ros:"address-pool,omitempty"`
	AddressesPerMAC              string `ros:"addresses-per-mac,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	IdleTimeout                  string `ros:"idle-timeout,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	KeepaliveTimeout             string `ros:"keepalive-timeout,omitempty"`
	LoginTimeout                 string `ros:"login-timeout,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	Profile                      string `ros:"profile,omitempty"`
}
```

### `/ip/hotspot/active`

- File: `(dari katalog Go)`
- Katalog: `hotspot-active`
- Sumber field: katalog 18

Command: 

```go
// IPHotspotActiveRow adalah kolom hasil `/ip/hotspot/active/print`.
type IPHotspotActiveRow struct {
	ID                           string `ros:".id,key"`
	User                         string `ros:"user"`
	Address                      string `ros:"address"`
	MACAddress                   string `ros:"mac-address"`
	Server                       string `ros:"server"`
	LoginBy                      string `ros:"login-by"`
	Uptime                       string `ros:"uptime"`
	SessionTimeLeft              string `ros:"session-time-left"`
	IdleTime                     string `ros:"idle-time"`
	BytesIn                      string `ros:"bytes-in"`
	BytesOut                     string `ros:"bytes-out"`
	Radius                       string `ros:"radius"`
}
```

### `/ip/hotspot/cookie`

- File: `(dari katalog Go)`
- Katalog: `hotspot-cookie`
- Sumber field: katalog 4

Command: 

```go
// IPHotspotCookieRow adalah kolom hasil `/ip/hotspot/cookie/print`.
type IPHotspotCookieRow struct {
	ID                           string `ros:".id,key"`
	User                         string `ros:"user"`
	MACAddress                   string `ros:"mac-address"`
	Domain                       string `ros:"domain"`
	ExpiresIn                    string `ros:"expires-in"`
}
```

### `/ip/hotspot/host`

- File: `(dari katalog Go)`
- Katalog: `hotspot-host`
- Sumber field: katalog 10

Command: 

```go
// IPHotspotHostRow adalah kolom hasil `/ip/hotspot/host/print`.
type IPHotspotHostRow struct {
	ID                           string `ros:".id,key"`
	MACAddress                   string `ros:"mac-address"`
	Address                      string `ros:"address"`
	ToAddress                    string `ros:"to-address"`
	Server                       string `ros:"server"`
	Uptime                       string `ros:"uptime"`
	IdleTime                     string `ros:"idle-time"`
	FoundBy                      string `ros:"found-by"`
	Authorized                   string `ros:"authorized"`
	Bypassed                     string `ros:"bypassed"`
	Comment                      string `ros:"comment"`
}
```

### `/ip/hotspot/ip-binding`

- File: `(dari katalog Go)`
- Katalog: `hotspot-ip-binding`
- Sumber field: katalog 7

Command: `add`(6), `set`(6)

```go
// IPHotspotIpbindingRow adalah kolom hasil `/ip/hotspot/ip-binding/print`.
type IPHotspotIpbindingRow struct {
	ID                           string `ros:".id,key"`
	MACAddress                   string `ros:"mac-address"`
	Address                      string `ros:"address"`
	ToAddress                    string `ros:"to-address"`
	Server                       string `ros:"server"`
	Type                         string `ros:"type"`
	Disabled                     string `ros:"disabled"`
	Comment                      string `ros:"comment"`
}
```

```go
// IPHotspotIpbindingAdd adalah argumen `/ip/hotspot/ip-binding/add`.
type IPHotspotIpbindingAdd struct {
	MACAddress                   string `ros:"mac-address,omitempty"`
	Address                      string `ros:"address,omitempty"`
	ToAddress                    string `ros:"to-address,omitempty"`
	Server                       string `ros:"server,omitempty"`
	Type                         string `ros:"type,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

```go
// IPHotspotIpbindingSet adalah argumen `/ip/hotspot/ip-binding/set`.
type IPHotspotIpbindingSet struct {
	MACAddress                   string `ros:"mac-address,omitempty"`
	Address                      string `ros:"address,omitempty"`
	ToAddress                    string `ros:"to-address,omitempty"`
	Server                       string `ros:"server,omitempty"`
	Type                         string `ros:"type,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

### `/ip/hotspot/profile`

- File: `(dari katalog Go)`
- Katalog: `hotspot-server-profile`
- Sumber field: katalog 11

Command: `add`(13), `set`(13)

```go
// IPHotspotProfileRow adalah kolom hasil `/ip/hotspot/profile/print`.
type IPHotspotProfileRow struct {
	ID                           string `ros:".id,key"`
	Name                         string `ros:"name"`
	HotspotAddress               string `ros:"hotspot-address"`
	DNSName                      string `ros:"dns-name"`
	HtmlDirectory                string `ros:"html-directory"`
	LoginBy                      string `ros:"login-by"`
	UseRadius                    string `ros:"use-radius"`
	RateLimit                    string `ros:"rate-limit"`
	HTTPCookieLifetime           string `ros:"http-cookie-lifetime"`
	InstallHotspotQueue          string `ros:"install-hotspot-queue"`
	MACAuthMode                  string `ros:"mac-auth-mode"`
	Comment                      string `ros:"comment"`
}
```

```go
// IPHotspotProfileAdd adalah argumen `/ip/hotspot/profile/add`.
type IPHotspotProfileAdd struct {
	Name                         string `ros:"name,required,omitempty"`
	HotspotAddress               string `ros:"hotspot-address,omitempty"`
	DNSName                      string `ros:"dns-name,omitempty"`
	HtmlDirectory                string `ros:"html-directory,omitempty"`
	LoginBy                      string `ros:"login-by,omitempty"`
	UseRadius                    string `ros:"use-radius,omitempty"`
	RateLimit                    string `ros:"rate-limit,omitempty"`
	HTTPCookieLifetime           string `ros:"http-cookie-lifetime,omitempty"`
	InstallHotspotQueue          string `ros:"install-hotspot-queue,omitempty"`
	MACAuthMode                  string `ros:"mac-auth-mode,omitempty"`
	RadiusMACFormat              string `ros:"radius-mac-format,omitempty"`
	SmtpServer                   string `ros:"smtp-server,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

```go
// IPHotspotProfileSet adalah argumen `/ip/hotspot/profile/set`.
type IPHotspotProfileSet struct {
	Name                         string `ros:"name,omitempty"`
	HotspotAddress               string `ros:"hotspot-address,omitempty"`
	DNSName                      string `ros:"dns-name,omitempty"`
	HtmlDirectory                string `ros:"html-directory,omitempty"`
	LoginBy                      string `ros:"login-by,omitempty"`
	UseRadius                    string `ros:"use-radius,omitempty"`
	RateLimit                    string `ros:"rate-limit,omitempty"`
	HTTPCookieLifetime           string `ros:"http-cookie-lifetime,omitempty"`
	InstallHotspotQueue          string `ros:"install-hotspot-queue,omitempty"`
	MACAuthMode                  string `ros:"mac-auth-mode,omitempty"`
	RadiusMACFormat              string `ros:"radius-mac-format,omitempty"`
	SmtpServer                   string `ros:"smtp-server,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

### `/ip/hotspot/user`

- File: `(dari katalog Go)`
- Katalog: `hotspot-user`
- Sumber field: katalog 12

Command: `add`(13), `set`(13)

```go
// IPHotspotUserRow adalah kolom hasil `/ip/hotspot/user/print`.
type IPHotspotUserRow struct {
	ID                           string `ros:".id,key"`
	Name                         string `ros:"name"`
	Profile                      string `ros:"profile"`
	Server                       string `ros:"server"`
	MACAddress                   string `ros:"mac-address"`
	Address                      string `ros:"address"`
	LimitUptime                  string `ros:"limit-uptime"`
	LimitBytesTotal              string `ros:"limit-bytes-total"`
	Uptime                       string `ros:"uptime"`
	BytesIn                      string `ros:"bytes-in"`
	BytesOut                     string `ros:"bytes-out"`
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
}
```

```go
// IPHotspotUserAdd adalah argumen `/ip/hotspot/user/add`.
type IPHotspotUserAdd struct {
	Name                         string `ros:"name,required,omitempty"`
	Password                     string `ros:"password,omitempty"`
	Profile                      string `ros:"profile,omitempty"`
	Server                       string `ros:"server,omitempty"`
	Address                      string `ros:"address,omitempty"`
	MACAddress                   string `ros:"mac-address,omitempty"`
	LimitUptime                  string `ros:"limit-uptime,omitempty"`
	LimitBytesIn                 string `ros:"limit-bytes-in,omitempty"`
	LimitBytesOut                string `ros:"limit-bytes-out,omitempty"`
	LimitBytesTotal              string `ros:"limit-bytes-total,omitempty"`
	Email                        string `ros:"email,omitempty"`
	Routes                       string `ros:"routes,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

```go
// IPHotspotUserSet adalah argumen `/ip/hotspot/user/set`.
type IPHotspotUserSet struct {
	Name                         string `ros:"name,omitempty"`
	Password                     string `ros:"password,omitempty"`
	Profile                      string `ros:"profile,omitempty"`
	Server                       string `ros:"server,omitempty"`
	Address                      string `ros:"address,omitempty"`
	MACAddress                   string `ros:"mac-address,omitempty"`
	LimitUptime                  string `ros:"limit-uptime,omitempty"`
	LimitBytesIn                 string `ros:"limit-bytes-in,omitempty"`
	LimitBytesOut                string `ros:"limit-bytes-out,omitempty"`
	LimitBytesTotal              string `ros:"limit-bytes-total,omitempty"`
	Email                        string `ros:"email,omitempty"`
	Routes                       string `ros:"routes,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

### `/ip/hotspot/user/profile`

- File: `(dari katalog Go)`
- Katalog: `hotspot-user-profile`
- Sumber field: katalog 31

Command: `add`(16), `set`(16)

```go
// IPHotspotUserProfileRow adalah kolom hasil `/ip/hotspot/user/profile/print`.
type IPHotspotUserProfileRow struct {
	ID                           string `ros:".id,key"`
	Name                         string `ros:"name"`
	RateLimit                    string `ros:"rate-limit"`
	SharedUsers                  string `ros:"shared-users"`
	SessionTimeout               string `ros:"session-timeout"`
	IdleTimeout                  string `ros:"idle-timeout"`
	KeepaliveTimeout             string `ros:"keepalive-timeout"`
	StatusAutorefresh            string `ros:"status-autorefresh"`
	AddressPool                  string `ros:"address-pool"`
	AddressList                  string `ros:"address-list"`
	ParentQueue                  string `ros:"parent-queue"`
	MACCookieTimeout             string `ros:"mac-cookie-timeout"`
	AddMACCookie                 string `ros:"add-mac-cookie"`
	TransparentProxy             string `ros:"transparent-proxy"`
	Comment                      string `ros:"comment"`
	OnLogin                      string `ros:"on-login"`
	OnLogout                     string `ros:"on-logout"`
	InsertQueueBefore            string `ros:"insert-queue-before"`
}
```

```go
// IPHotspotUserProfileAdd adalah argumen `/ip/hotspot/user/profile/add`.
type IPHotspotUserProfileAdd struct {
	Name                         string `ros:"name,required,omitempty"`
	RateLimit                    string `ros:"rate-limit,omitempty"`
	SharedUsers                  string `ros:"shared-users,omitempty"`
	SessionTimeout               string `ros:"session-timeout,omitempty"`
	IdleTimeout                  string `ros:"idle-timeout,omitempty"`
	KeepaliveTimeout             string `ros:"keepalive-timeout,omitempty"`
	StatusAutorefresh            string `ros:"status-autorefresh,omitempty"`
	AddressPool                  string `ros:"address-pool,omitempty"`
	AddressList                  string `ros:"address-list,omitempty"`
	ParentQueue                  string `ros:"parent-queue,omitempty"`
	MACCookieTimeout             string `ros:"mac-cookie-timeout,omitempty"`
	AddMACCookie                 string `ros:"add-mac-cookie,omitempty"`
	TransparentProxy             string `ros:"transparent-proxy,omitempty"`
	OnLogin                      string `ros:"on-login,omitempty"`
	OnLogout                     string `ros:"on-logout,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

```go
// IPHotspotUserProfileSet adalah argumen `/ip/hotspot/user/profile/set`.
type IPHotspotUserProfileSet struct {
	Name                         string `ros:"name,omitempty"`
	RateLimit                    string `ros:"rate-limit,omitempty"`
	SharedUsers                  string `ros:"shared-users,omitempty"`
	SessionTimeout               string `ros:"session-timeout,omitempty"`
	IdleTimeout                  string `ros:"idle-timeout,omitempty"`
	KeepaliveTimeout             string `ros:"keepalive-timeout,omitempty"`
	StatusAutorefresh            string `ros:"status-autorefresh,omitempty"`
	AddressPool                  string `ros:"address-pool,omitempty"`
	AddressList                  string `ros:"address-list,omitempty"`
	ParentQueue                  string `ros:"parent-queue,omitempty"`
	MACCookieTimeout             string `ros:"mac-cookie-timeout,omitempty"`
	AddMACCookie                 string `ros:"add-mac-cookie,omitempty"`
	TransparentProxy             string `ros:"transparent-proxy,omitempty"`
	OnLogin                      string `ros:"on-login,omitempty"`
	OnLogout                     string `ros:"on-logout,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

### `/ip/hotspot/walled-garden`

- File: `(dari katalog Go)`
- Katalog: `hotspot-walled-garden`
- Sumber field: katalog 11

Command: `add`(9), `set`(9)

```go
// IPHotspotWalledgardenRow adalah kolom hasil `/ip/hotspot/walled-garden/print`.
type IPHotspotWalledgardenRow struct {
	ID                           string `ros:".id,key"`
	Action                       string `ros:"action"`
	Server                       string `ros:"server"`
	SrcAddress                   string `ros:"src-address"`
	DstAddress                   string `ros:"dst-address"`
	DstHost                      string `ros:"dst-host"`
	DstPort                      string `ros:"dst-port"`
	Path                         string `ros:"path"`
	Method                       string `ros:"method"`
	Hits                         string `ros:"hits"`
	Disabled                     string `ros:"disabled"`
	Comment                      string `ros:"comment"`
}
```

```go
// IPHotspotWalledgardenAdd adalah argumen `/ip/hotspot/walled-garden/add`.
type IPHotspotWalledgardenAdd struct {
	Action                       string `ros:"action,omitempty"`
	Server                       string `ros:"server,omitempty"`
	SrcAddress                   string `ros:"src-address,omitempty"`
	DstAddress                   string `ros:"dst-address,omitempty"`
	DstHost                      string `ros:"dst-host,omitempty"`
	DstPort                      string `ros:"dst-port,omitempty"`
	Path                         string `ros:"path,omitempty"`
	Method                       string `ros:"method,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

```go
// IPHotspotWalledgardenSet adalah argumen `/ip/hotspot/walled-garden/set`.
type IPHotspotWalledgardenSet struct {
	Action                       string `ros:"action,omitempty"`
	Server                       string `ros:"server,omitempty"`
	SrcAddress                   string `ros:"src-address,omitempty"`
	DstAddress                   string `ros:"dst-address,omitempty"`
	DstHost                      string `ros:"dst-host,omitempty"`
	DstPort                      string `ros:"dst-port,omitempty"`
	Path                         string `ros:"path,omitempty"`
	Method                       string `ros:"method,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

### `/ip/hotspot/walled-garden/ip`

- File: `(dari katalog Go)`
- Katalog: `hotspot-walled-garden-ip`
- Sumber field: katalog 9

Command: `add`(8), `set`(8)

```go
// IPHotspotWalledgardenIPRow adalah kolom hasil `/ip/hotspot/walled-garden/ip/print`.
type IPHotspotWalledgardenIPRow struct {
	ID                           string `ros:".id,key"`
	Action                       string `ros:"action"`
	Server                       string `ros:"server"`
	SrcAddress                   string `ros:"src-address"`
	DstAddress                   string `ros:"dst-address"`
	DstHost                      string `ros:"dst-host"`
	Protocol                     string `ros:"protocol"`
	DstPort                      string `ros:"dst-port"`
	Disabled                     string `ros:"disabled"`
	Comment                      string `ros:"comment"`
}
```

```go
// IPHotspotWalledgardenIPAdd adalah argumen `/ip/hotspot/walled-garden/ip/add`.
type IPHotspotWalledgardenIPAdd struct {
	Action                       string `ros:"action,omitempty"`
	Server                       string `ros:"server,omitempty"`
	SrcAddress                   string `ros:"src-address,omitempty"`
	DstAddress                   string `ros:"dst-address,omitempty"`
	DstHost                      string `ros:"dst-host,omitempty"`
	Protocol                     string `ros:"protocol,omitempty"`
	DstPort                      string `ros:"dst-port,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

```go
// IPHotspotWalledgardenIPSet adalah argumen `/ip/hotspot/walled-garden/ip/set`.
type IPHotspotWalledgardenIPSet struct {
	Action                       string `ros:"action,omitempty"`
	Server                       string `ros:"server,omitempty"`
	SrcAddress                   string `ros:"src-address,omitempty"`
	DstAddress                   string `ros:"dst-address,omitempty"`
	DstHost                      string `ros:"dst-host,omitempty"`
	Protocol                     string `ros:"protocol,omitempty"`
	DstPort                      string `ros:"dst-port,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

### `/ip/ipsec`

- File: `ip_ipsec.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command

Command: `export`(8)

_Tidak ada struct: menu ini bukan tabel yang bisa di-print atau tidak tersedia di router uji (mis. paket tidak terpasang)._

### `/ip/kid-control`

- File: `ip_kid-control.json`
- Sumber field: completion 23

Command: 

```go
// IPKidcontrolRow adalah kolom hasil `/ip/kid-control/print`.
type IPKidcontrolRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Blocked                      string `ros:"blocked"`
	Disabled                     string `ros:"disabled"`
	Fri                          string `ros:"fri"`
	Mon                          string `ros:"mon"`
	Name                         string `ros:"name"`
	PauseLeft                    string `ros:"pause-left"`
	PauseTill                    string `ros:"pause-till"`
	Paused                       string `ros:"paused"`
	RateLimit                    string `ros:"rate-limit"`
	RateLimited                  string `ros:"rate-limited"`
	Sat                          string `ros:"sat"`
	Sun                          string `ros:"sun"`
	Thu                          string `ros:"thu"`
	Tue                          string `ros:"tue"`
	TurFri                       string `ros:"tur-fri"`
	TurMon                       string `ros:"tur-mon"`
	TurSat                       string `ros:"tur-sat"`
	TurSun                       string `ros:"tur-sun"`
	TurThu                       string `ros:"tur-thu"`
	TurTue                       string `ros:"tur-tue"`
	TurWed                       string `ros:"tur-wed"`
	Wed                          string `ros:"wed"`
}
```

### `/ip/media`

- File: `ip_media.json`
- Sumber field: completion 9

Command: 

```go
// IPMediaRow adalah kolom hasil `/ip/media/print`.
type IPMediaRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	AllowedHostname              string `ros:"allowed-hostname"`
	AllowedIP                    string `ros:"allowed-ip"`
	Disabled                     string `ros:"disabled"`
	Dynamic                      string `ros:"dynamic"`
	FriendlyName                 string `ros:"friendly-name"`
	Interface                    string `ros:"interface"`
	Path                         string `ros:"path"`
	Status                       string `ros:"status"`
}
```

### `/ip/nat-pmp`

- File: `ip_nat-pmp.json`
- Sumber field: live 1

Command: 

```go
// IPNatpmpRow adalah kolom hasil `/ip/nat-pmp/print`.
type IPNatpmpRow struct {
	Enabled                      string `ros:"enabled"`
}
```

### `/ip/neighbor`

- File: `ip_neighbor.json`
- Katalog: `ip-neighbor`
- Sumber field: live 20, completion 21, katalog 8

Command: `export`(8), `find`(1), `get`(4), `print`(22), `unset`(2)

```go
// IPNeighborRow adalah kolom hasil `/ip/neighbor/print`.
type IPNeighborRow struct {
	ID                           string `ros:".id,key"`
	Address                      string `ros:"address"`
	Address4                     string `ros:"address4"`
	Address6                     string `ros:"address6"`
	Age                          string `ros:"age"`
	Board                        string `ros:"board"`
	DiscoveredBy                 string `ros:"discovered-by"`
	Identity                     string `ros:"identity"`
	Interface                    string `ros:"interface"`
	InterfaceName                string `ros:"interface-name"`
	IPv6                         string `ros:"ipv6"`
	MACAddress                   string `ros:"mac-address"`
	Platform                     string `ros:"platform"`
	SoftwareID                   string `ros:"software-id"`
	SystemCAPs                   string `ros:"system-caps"`
	SystemCAPsEnabled            string `ros:"system-caps-enabled"`
	SystemDescription            string `ros:"system-description"`
	Unpack                       string `ros:"unpack"`
	Uptime                       string `ros:"uptime"`
	Version                      string `ros:"version"`
	About                        string `ros:"about"`
	Running                      string `ros:"running"`
}
```

### `/ip/packing`

- File: `ip_packing.json`
- Sumber field: completion 6

Command: `add`(6), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(6), `set`(6)

```go
// IPPackingRow adalah kolom hasil `/ip/packing/print`.
type IPPackingRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	AggregatedSize               string `ros:"aggregated-size"`
	Disabled                     string `ros:"disabled"`
	Interface                    string `ros:"interface"`
	Packing                      string `ros:"packing"`
	Unpacking                    string `ros:"unpacking"`
}
```

```go
// IPPackingAdd adalah argumen `/ip/packing/add`.
type IPPackingAdd struct {
	AggregatedSize               string `ros:"aggregated-size,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	Packing                      string `ros:"packing,omitempty"`
	Unpacking                    string `ros:"unpacking,omitempty"`
}
```

```go
// IPPackingSet adalah argumen `/ip/packing/set`.
type IPPackingSet struct {
	AggregatedSize               string `ros:"aggregated-size,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	Packing                      string `ros:"packing,omitempty"`
	Unpacking                    string `ros:"unpacking,omitempty"`
}
```

### `/ip/pool`

- File: `ip_pool.json`
- Katalog: `ip-pool`
- Sumber field: live 6, completion 8, katalog 4

Command: `add`(5), `comment`(2), `edit`(2), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(5), `set`(5)

```go
// IPPoolRow adalah kolom hasil `/ip/pool/print`.
type IPPoolRow struct {
	ID                           string `ros:".id,key"`
	Available                    string `ros:"available"`
	// string value
	Name                         string `ros:"name"`
	// A.B.C.D[-A.B.C.D |0..32 |/A.B.C.D ] (IP address range)
	Ranges                       string `ros:"ranges"`
	Total                        string `ros:"total"`
	Used                         string `ros:"used"`
	About                        string `ros:"about"`
	// string value
	Comment                      string `ros:"comment"`
	NextPool                     string `ros:"next-pool"`
}
```

```go
// IPPoolAdd adalah argumen `/ip/pool/add`.
type IPPoolAdd struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	// string value
	Name                         string `ros:"name,required,omitempty"`
	NextPool                     string `ros:"next-pool,omitempty"`
	// A.B.C.D[-A.B.C.D |0..32 |/A.B.C.D ] (IP address range)
	Ranges                       string `ros:"ranges,required,omitempty"`
}
```

```go
// IPPoolSet adalah argumen `/ip/pool/set`.
type IPPoolSet struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	NextPool                     string `ros:"next-pool,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	// A.B.C.D[-A.B.C.D |0..32 |/A.B.C.D ] (IP address range)
	Ranges                       string `ros:"ranges,omitempty"`
}
```

### `/ip/pool/used`

- File: `(dari katalog Go)`
- Katalog: `ip-pool-used`
- Sumber field: katalog 4

Command: 

```go
// IPPoolUsedRow adalah kolom hasil `/ip/pool/used/print`.
type IPPoolUsedRow struct {
	ID                           string `ros:".id,key"`
	Pool                         string `ros:"pool"`
	Address                      string `ros:"address"`
	Owner                        string `ros:"owner"`
	Info                         string `ros:"info"`
}
```

### `/ip/proxy`

- File: `ip_proxy.json`
- Sumber field: live 17

Command: 

```go
// IPProxyRow adalah kolom hasil `/ip/proxy/print`.
type IPProxyRow struct {
	AlwaysFromCache              string `ros:"always-from-cache"`
	Anonymous                    string `ros:"anonymous"`
	CacheAdministrator           string `ros:"cache-administrator"`
	CacheHitDscp                 string `ros:"cache-hit-dscp"`
	CacheOnDisk                  string `ros:"cache-on-disk"`
	CachePath                    string `ros:"cache-path"`
	Enabled                      string `ros:"enabled"`
	MaxCacheObjectSize           string `ros:"max-cache-object-size"`
	MaxCacheSize                 string `ros:"max-cache-size"`
	MaxClientConnections         string `ros:"max-client-connections"`
	MaxFreshTime                 string `ros:"max-fresh-time"`
	MaxServerConnections         string `ros:"max-server-connections"`
	ParentProxy                  string `ros:"parent-proxy"`
	ParentProxyPort              string `ros:"parent-proxy-port"`
	Port                         string `ros:"port"`
	SerializeConnections         string `ros:"serialize-connections"`
	SrcAddress                   string `ros:"src-address"`
}
```

### `/ip/reverse-proxy`

- File: `ip_reverse-proxy.json`
- Sumber field: completion 9

Command: 

```go
// IPReverseproxyRow adalah kolom hasil `/ip/reverse-proxy/print`.
type IPReverseproxyRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Certificate                  string `ros:"certificate"`
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	Dynamic                      string `ros:"dynamic"`
	IPAddress                    string `ros:"ip-address"`
	Port                         string `ros:"port"`
	Sni                          string `ros:"sni"`
	VRF                          string `ros:"vrf"`
}
```

### `/ip/route`

- File: `ip_route.json`
- Katalog: `ip-route`
- Sumber field: live 15, completion 30, katalog 20

Command: `add`(13), `check`(12), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(22), `remove`(1), `reset`(13), `set`(13), `unset`(2)

```go
// IPRouteRow adalah kolom hasil `/ip/route/print`.
type IPRouteRow struct {
	ID                           string `ros:".id,key"`
	Active                       string `ros:"active"`
	Connect                      string `ros:"connect"`
	DHCP                         string `ros:"dhcp"`
	Distance                     string `ros:"distance"`
	DstAddress                   string `ros:"dst-address"`
	Dynamic                      string `ros:"dynamic"`
	Gateway                      string `ros:"gateway"`
	ImmediateGw                  string `ros:"immediate-gw"`
	Inactive                     string `ros:"inactive"`
	LocalAddress                 string `ros:"local-address"`
	RoutingTable                 string `ros:"routing-table"`
	Scope                        string `ros:"scope"`
	TargetScope                  string `ros:"target-scope"`
	VRFInterface                 string `ros:"vrf-interface"`
	About                        string `ros:"about"`
	BGP                          string `ros:"bgp"`
	BGPMPLSVPN                   string `ros:"bgp-mpls-vpn"`
	Blackhole                    string `ros:"blackhole"`
	CheckGateway                 string `ros:"check-gateway"`
	// string value
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	Ecmp                         string `ros:"ecmp"`
	HWOffloaded                  string `ros:"hw-offloaded"`
	IsIs                         string `ros:"is-is"`
	Modem                        string `ros:"modem"`
	OSPF                         string `ros:"ospf"`
	PrefSrc                      string `ros:"pref-src"`
	RIP                          string `ros:"rip"`
	Static                       string `ros:"static"`
	VPN                          string `ros:"vpn"`
	RoutingMark                  string `ros:"routing-mark"`
}
```

```go
// IPRouteAdd adalah argumen `/ip/route/add`.
type IPRouteAdd struct {
	Blackhole                    string `ros:"blackhole,omitempty"`
	CheckGateway                 string `ros:"check-gateway,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Distance                     string `ros:"distance,omitempty"`
	DstAddress                   string `ros:"dst-address,required,omitempty"`
	Gateway                      string `ros:"gateway,required,omitempty"`
	PrefSrc                      string `ros:"pref-src,omitempty"`
	RoutingTable                 string `ros:"routing-table,omitempty"`
	Scope                        string `ros:"scope,omitempty"`
	TargetScope                  string `ros:"target-scope,omitempty"`
	VRFInterface                 string `ros:"vrf-interface,omitempty"`
}
```

```go
// IPRouteSet adalah argumen `/ip/route/set`.
type IPRouteSet struct {
	Blackhole                    string `ros:"blackhole,omitempty"`
	CheckGateway                 string `ros:"check-gateway,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Distance                     string `ros:"distance,omitempty"`
	DstAddress                   string `ros:"dst-address,omitempty"`
	Gateway                      string `ros:"gateway,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	PrefSrc                      string `ros:"pref-src,omitempty"`
	RoutingTable                 string `ros:"routing-table,omitempty"`
	Scope                        string `ros:"scope,omitempty"`
	TargetScope                  string `ros:"target-scope,omitempty"`
	VRFInterface                 string `ros:"vrf-interface,omitempty"`
}
```

### `/ip/service`

- File: `ip_service.json`
- Katalog: `ip-service`
- Sumber field: live 15, completion 18, katalog 4

Command: `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(21), `reset`(9), `set`(9)

```go
// IPServiceRow adalah kolom hasil `/ip/service/print`.
type IPServiceRow struct {
	ID                           string `ros:".id,key"`
	AvailableFrom                string `ros:"available-from"`
	Certificate                  string `ros:"certificate"`
	Connection                   string `ros:"connection"`
	Disabled                     string `ros:"disabled"`
	Dynamic                      string `ros:"dynamic"`
	Invalid                      string `ros:"invalid"`
	Local                        string `ros:"local"`
	MaxSessions                  string `ros:"max-sessions"`
	Name                         string `ros:"name"`
	Port                         string `ros:"port"`
	Proto                        string `ros:"proto"`
	Remote                       string `ros:"remote"`
	TLSVersion                   string `ros:"tls-version"`
	VRF                          string `ros:"vrf"`
	About                        string `ros:"about"`
	Address                      string `ros:"address"`
	Container                    string `ros:"container"`
	Netns                        string `ros:"netns"`
}
```

```go
// IPServiceSet adalah argumen `/ip/service/set`.
type IPServiceSet struct {
	Address                      string `ros:"address,omitempty"`
	AvailableFrom                string `ros:"available-from,omitempty"`
	Certificate                  string `ros:"certificate,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	MaxSessions                  string `ros:"max-sessions,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	Port                         string `ros:"port,omitempty"`
	TLSVersion                   string `ros:"tls-version,omitempty"`
	VRF                          string `ros:"vrf,omitempty"`
}
```

### `/ip/settings`

- File: `ip_settings.json`
- Sumber field: live 23

Command: 

```go
// IPSettingsRow adalah kolom hasil `/ip/settings/print`.
type IPSettingsRow struct {
	AcceptRedirects              string `ros:"accept-redirects"`
	AcceptSourceRoute            string `ros:"accept-source-route"`
	AllowFastPath                string `ros:"allow-fast-path"`
	ARPTimeout                   string `ros:"arp-timeout"`
	IcmpErrorsUseInboundInterfaceAddress string `ros:"icmp-errors-use-inbound-interface-address"`
	IcmpRateLimit                string `ros:"icmp-rate-limit"`
	IcmpRateMask                 string `ros:"icmp-rate-mask"`
	IPForward                    string `ros:"ip-forward"`
	IPv4FastPathActive           string `ros:"ipv4-fast-path-active"`
	IPv4FastPathBytes            string `ros:"ipv4-fast-path-bytes"`
	IPv4FastPathPackets          string `ros:"ipv4-fast-path-packets"`
	IPv4FasttrackActive          string `ros:"ipv4-fasttrack-active"`
	IPv4FasttrackBytes           string `ros:"ipv4-fasttrack-bytes"`
	IPv4FasttrackPackets         string `ros:"ipv4-fasttrack-packets"`
	IPv4FragmentTime             string `ros:"ipv4-fragment-time"`
	IPv4HighFragmentThresh       string `ros:"ipv4-high-fragment-thresh"`
	IPv4MultipathHashPolicy      string `ros:"ipv4-multipath-hash-policy"`
	MaxNeighborEntries           string `ros:"max-neighbor-entries"`
	RpFilter                     string `ros:"rp-filter"`
	SecureRedirects              string `ros:"secure-redirects"`
	SendRedirects                string `ros:"send-redirects"`
	TcpSyncookies                string `ros:"tcp-syncookies"`
	TcpTimestamps                string `ros:"tcp-timestamps"`
}
```

### `/ip/smb`

- File: `ip_smb.json`
- Sumber field: live 5

Command: 

```go
// IPSmbRow adalah kolom hasil `/ip/smb/print`.
type IPSmbRow struct {
	Comment                      string `ros:"comment"`
	Domain                       string `ros:"domain"`
	Enabled                      string `ros:"enabled"`
	Interfaces                   string `ros:"interfaces"`
	Status                       string `ros:"status"`
}
```

### `/ip/socks`

- File: `ip_socks.json`
- Sumber field: live 7

Command: 

```go
// IPSocksRow adalah kolom hasil `/ip/socks/print`.
type IPSocksRow struct {
	AuthMethod                   string `ros:"auth-method"`
	ConnectionIdleTimeout        string `ros:"connection-idle-timeout"`
	Enabled                      string `ros:"enabled"`
	MaxConnections               string `ros:"max-connections"`
	Port                         string `ros:"port"`
	Version                      string `ros:"version"`
	VRF                          string `ros:"vrf"`
}
```

### `/ip/socksify`

- File: `ip_socksify.json`
- Sumber field: completion 11

Command: 

```go
// IPSocksifyRow adalah kolom hasil `/ip/socksify/print`.
type IPSocksifyRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Comment                      string `ros:"comment"`
	ConnectionTimeout            string `ros:"connection-timeout"`
	Disabled                     string `ros:"disabled"`
	Dynamic                      string `ros:"dynamic"`
	Name                         string `ros:"name"`
	Port                         string `ros:"port"`
	Socks5Password               string `ros:"socks5-password"`
	Socks5Port                   string `ros:"socks5-port"`
	Socks5Server                 string `ros:"socks5-server"`
	Socks5User                   string `ros:"socks5-user"`
}
```

### `/ip/ssh`

- File: `ip_ssh.json`
- Sumber field: live 7

Command: 

```go
// IPSSHRow adalah kolom hasil `/ip/ssh/print`.
type IPSSHRow struct {
	Ciphers                      string `ros:"ciphers"`
	ForwardingEnabled            string `ros:"forwarding-enabled"`
	HostKeySize                  string `ros:"host-key-size"`
	HostKeyType                  string `ros:"host-key-type"`
	PasswordAuthentication       string `ros:"password-authentication"`
	PublickeyAuthenticationOptions string `ros:"publickey-authentication-options"`
	StrongCrypto                 string `ros:"strong-crypto"`
}
```

### `/ip/tftp`

- File: `ip_tftp.json`
- Sumber field: completion 12

Command: 

```go
// IPTftpRow adalah kolom hasil `/ip/tftp/print`.
type IPTftpRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Allow                        string `ros:"allow"`
	AllowOverwrite               string `ros:"allow-overwrite"`
	AllowRollover                string `ros:"allow-rollover"`
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	Hits                         string `ros:"hits"`
	IPAddresses                  string `ros:"ip-addresses"`
	ReadOnly                     string `ros:"read-only"`
	ReadingWindowSize            string `ros:"reading-window-size"`
	RealFilename                 string `ros:"real-filename"`
	ReqFilename                  string `ros:"req-filename"`
}
```

### `/ip/traffic-flow`

- File: `ip_traffic-flow.json`
- Sumber field: live 8

Command: `edit`(1), `export`(8), `get`(3), `monitor`(10), `print`(5), `reset`(8), `set`(8)

```go
// IPTrafficflowRow adalah kolom hasil `/ip/traffic-flow/print`.
type IPTrafficflowRow struct {
	ID                           string `ros:".id,key"`
	// 1s..7101w3d6h28m15s (time interval)
	ActiveFlowTimeout            string `ros:"active-flow-timeout"`
	CacheEntries                 string `ros:"cache-entries"`
	Enabled                      string `ros:"enabled"`
	// 1s..7101w3d6h28m15s (time interval)
	InactiveFlowTimeout          string `ros:"inactive-flow-timeout"`
	Interfaces                   string `ros:"interfaces"`
	PacketSampling               string `ros:"packet-sampling"`
	SamplingInterval             string `ros:"sampling-interval"`
	SamplingSpace                string `ros:"sampling-space"`
}
```

```go
// IPTrafficflowSet adalah argumen `/ip/traffic-flow/set`.
type IPTrafficflowSet struct {
	// 1s..7101w3d6h28m15s (time interval)
	ActiveFlowTimeout            string `ros:"active-flow-timeout,omitempty"`
	CacheEntries                 string `ros:"cache-entries,omitempty"`
	Enabled                      string `ros:"enabled,omitempty"`
	// 1s..7101w3d6h28m15s (time interval)
	InactiveFlowTimeout          string `ros:"inactive-flow-timeout,omitempty"`
	Interfaces                   string `ros:"interfaces,omitempty"`
	PacketSampling               string `ros:"packet-sampling,omitempty"`
	SamplingInterval             string `ros:"sampling-interval,omitempty"`
	SamplingSpace                string `ros:"sampling-space,omitempty"`
}
```

### `/ip/upnp`

- File: `ip_upnp.json`
- Sumber field: live 3

Command: 

```go
// IPUpnpRow adalah kolom hasil `/ip/upnp/print`.
type IPUpnpRow struct {
	AllowDisableExternalInterface string `ros:"allow-disable-external-interface"`
	Enabled                      string `ros:"enabled"`
	ShowDummyRule                string `ros:"show-dummy-rule"`
}
```

### `/ip/vrf`

- File: `ip_vrf.json`
- Sumber field: live 5, completion 6

Command: 

```go
// IPVRFRow adalah kolom hasil `/ip/vrf/print`.
type IPVRFRow struct {
	ID                           string `ros:".id,key"`
	Builtin                      string `ros:"builtin"`
	Disabled                     string `ros:"disabled"`
	Interfaces                   string `ros:"interfaces"`
	Name                         string `ros:"name"`
	About                        string `ros:"about"`
	Comment                      string `ros:"comment"`
}
```

### `/ipv6/address`

- File: `ipv6_address.json`
- Sumber field: live 16, completion 21

Command: 

```go
// IPv6AddressRow adalah kolom hasil `/ipv6/address/print`.
type IPv6AddressRow struct {
	ID                           string `ros:".id,key"`
	ActualInterface              string `ros:"actual-interface"`
	Address                      string `ros:"address"`
	Advertise                    string `ros:"advertise"`
	AutoLinkLocal                string `ros:"auto-link-local"`
	Deprecated                   string `ros:"deprecated"`
	Disabled                     string `ros:"disabled"`
	Dynamic                      string `ros:"dynamic"`
	EuiX64                       string `ros:"eui-64"`
	FromPool                     string `ros:"from-pool"`
	Interface                    string `ros:"interface"`
	Invalid                      string `ros:"invalid"`
	LinkLocal                    string `ros:"link-local"`
	NoDad                        string `ros:"no-dad"`
	Slave                        string `ros:"slave"`
	VRF                          string `ros:"vrf"`
	About                        string `ros:"about"`
	Comment                      string `ros:"comment"`
	FromPoolPolicy               string `ros:"from-pool-policy"`
	Global                       string `ros:"global"`
	Preferred                    string `ros:"preferred"`
	Valid                        string `ros:"valid"`
}
```

### `/ipv6/dhcp-client`

- File: `ipv6_dhcp-client.json`
- Sumber field: completion 33

Command: 

```go
// IPv6DhcpclientRow adalah kolom hasil `/ipv6/dhcp-client/print`.
type IPv6DhcpclientRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	AcceptPrefixWithoutAddress   string `ros:"accept-prefix-without-address"`
	AddDefaultRoute              string `ros:"add-default-route"`
	Address                      string `ros:"address"`
	AllowReconfigure             string `ros:"allow-reconfigure"`
	CheckGateway                 string `ros:"check-gateway"`
	Comment                      string `ros:"comment"`
	CustomDuid                   string `ros:"custom-duid"`
	CustomIanaID                 string `ros:"custom-iana-id"`
	CustomIapdID                 string `ros:"custom-iapd-id"`
	DefaultRouteDistance         string `ros:"default-route-distance"`
	DefaultRouteTables           string `ros:"default-route-tables"`
	DHCPOptions                  string `ros:"dhcp-options"`
	DHCPServerV6                 string `ros:"dhcp-server-v6"`
	Disabled                     string `ros:"disabled"`
	Duid                         string `ros:"duid"`
	Dynamic                      string `ros:"dynamic"`
	Interface                    string `ros:"interface"`
	Invalid                      string `ros:"invalid"`
	PoolName                     string `ros:"pool-name"`
	PoolPrefixLength             string `ros:"pool-prefix-length"`
	Prefix                       string `ros:"prefix"`
	PrefixAddressLists           string `ros:"prefix-address-lists"`
	PrefixHint                   string `ros:"prefix-hint"`
	RapidCommit                  string `ros:"rapid-commit"`
	ReconfigureKey               string `ros:"reconfigure-key"`
	ReconfigureLastCounter       string `ros:"reconfigure-last-counter"`
	Request                      string `ros:"request"`
	Script                       string `ros:"script"`
	Status                       string `ros:"status"`
	UseInterfaceDuid             string `ros:"use-interface-duid"`
	UsePeerDNS                   string `ros:"use-peer-dns"`
	ValidateServerDuid           string `ros:"validate-server-duid"`
}
```

### `/ipv6/dhcp-relay`

- File: `ipv6_dhcp-relay.json`
- Sumber field: completion 11

Command: 

```go
// IPv6DHCPRelayRow adalah kolom hasil `/ipv6/dhcp-relay/print`.
type IPv6DHCPRelayRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Comment                      string `ros:"comment"`
	DelayThreshold               string `ros:"delay-threshold"`
	DHCPOptions                  string `ros:"dhcp-options"`
	DHCPServer                   string `ros:"dhcp-server"`
	Disabled                     string `ros:"disabled"`
	Interface                    string `ros:"interface"`
	Invalid                      string `ros:"invalid"`
	LinkAddress                  string `ros:"link-address"`
	Name                         string `ros:"name"`
	StoreRelayedBindings         string `ros:"store-relayed-bindings"`
}
```

### `/ipv6/dhcp-server`

- File: `ipv6_dhcp-server.json`
- Sumber field: completion 23

Command: 

```go
// IPv6DhcpserverRow adalah kolom hasil `/ipv6/dhcp-server/print`.
type IPv6DhcpserverRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	AddressLists                 string `ros:"address-lists"`
	AddressPool                  string `ros:"address-pool"`
	AllowDualStackQueue          string `ros:"allow-dual-stack-queue"`
	BindingScript                string `ros:"binding-script"`
	Comment                      string `ros:"comment"`
	DHCPOption                   string `ros:"dhcp-option"`
	Disabled                     string `ros:"disabled"`
	Duid                         string `ros:"duid"`
	Dynamic                      string `ros:"dynamic"`
	IgnoreIaNaBindings           string `ros:"ignore-ia-na-bindings"`
	InsertQueueBefore            string `ros:"insert-queue-before"`
	Interface                    string `ros:"interface"`
	Invalid                      string `ros:"invalid"`
	LeaseTime                    string `ros:"lease-time"`
	Name                         string `ros:"name"`
	ParentQueue                  string `ros:"parent-queue"`
	Preference                   string `ros:"preference"`
	PrefixPool                   string `ros:"prefix-pool"`
	RapidCommit                  string `ros:"rapid-commit"`
	RouteDistance                string `ros:"route-distance"`
	UseRadius                    string `ros:"use-radius"`
	UseReconfigure               string `ros:"use-reconfigure"`
}
```

### `/ipv6/nd`

- File: `ipv6_nd.json`
- Sumber field: live 16, completion 20

Command: 

```go
// IPv6NdRow adalah kolom hasil `/ipv6/nd/print`.
type IPv6NdRow struct {
	ID                           string `ros:".id,key"`
	AdvertiseDNS                 string `ros:"advertise-dns"`
	AdvertiseMACAddress          string `ros:"advertise-mac-address"`
	Default                      string `ros:"default"`
	Disabled                     string `ros:"disabled"`
	HopLimit                     string `ros:"hop-limit"`
	Interface                    string `ros:"interface"`
	ManagedAddressConfiguration  string `ros:"managed-address-configuration"`
	MTU                          string `ros:"mtu"`
	OtherConfiguration           string `ros:"other-configuration"`
	RaDelay                      string `ros:"ra-delay"`
	RaInterval                   string `ros:"ra-interval"`
	RaLifetime                   string `ros:"ra-lifetime"`
	RaPreference                 string `ros:"ra-preference"`
	ReachableTime                string `ros:"reachable-time"`
	RetransmitInterval           string `ros:"retransmit-interval"`
	About                        string `ros:"about"`
	Comment                      string `ros:"comment"`
	DNS                          string `ros:"dns"`
	Invalid                      string `ros:"invalid"`
	Pref64                       string `ros:"pref64"`
}
```

### `/ipv6/neighbor`

- File: `ipv6_neighbor.json`
- Sumber field: completion 10

Command: 

```go
// IPv6NeighborRow adalah kolom hasil `/ipv6/neighbor/print`.
type IPv6NeighborRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Address                      string `ros:"address"`
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	Dynamic                      string `ros:"dynamic"`
	Interface                    string `ros:"interface"`
	MACAddress                   string `ros:"mac-address"`
	Router                       string `ros:"router"`
	Status                       string `ros:"status"`
	VRF                          string `ros:"vrf"`
}
```

### `/ipv6/pool`

- File: `ipv6_pool.json`
- Sumber field: completion 11

Command: 

```go
// IPv6PoolRow adalah kolom hasil `/ipv6/pool/print`.
type IPv6PoolRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	ActualPrefix                 string `ros:"actual-prefix"`
	Comment                      string `ros:"comment"`
	Dynamic                      string `ros:"dynamic"`
	FromPool                     string `ros:"from-pool"`
	Inactive                     string `ros:"inactive"`
	Name                         string `ros:"name"`
	PreferredLifetime            string `ros:"preferred-lifetime"`
	Prefix                       string `ros:"prefix"`
	PrefixLength                 string `ros:"prefix-length"`
	ValidLifetime                string `ros:"valid-lifetime"`
}
```

### `/ipv6/route`

- File: `ipv6_route.json`
- Sumber field: live 12, completion 30

Command: 

```go
// IPv6RouteRow adalah kolom hasil `/ipv6/route/print`.
type IPv6RouteRow struct {
	ID                           string `ros:".id,key"`
	Active                       string `ros:"active"`
	Connect                      string `ros:"connect"`
	Distance                     string `ros:"distance"`
	DstAddress                   string `ros:"dst-address"`
	Dynamic                      string `ros:"dynamic"`
	Gateway                      string `ros:"gateway"`
	ImmediateGw                  string `ros:"immediate-gw"`
	Inactive                     string `ros:"inactive"`
	RoutingTable                 string `ros:"routing-table"`
	Scope                        string `ros:"scope"`
	TargetScope                  string `ros:"target-scope"`
	About                        string `ros:"about"`
	BGP                          string `ros:"bgp"`
	BGPMPLSVPN                   string `ros:"bgp-mpls-vpn"`
	Blackhole                    string `ros:"blackhole"`
	CheckGateway                 string `ros:"check-gateway"`
	Comment                      string `ros:"comment"`
	DHCP                         string `ros:"dhcp"`
	Disabled                     string `ros:"disabled"`
	Ecmp                         string `ros:"ecmp"`
	HWOffloaded                  string `ros:"hw-offloaded"`
	IsIs                         string `ros:"is-is"`
	Modem                        string `ros:"modem"`
	OSPF                         string `ros:"ospf"`
	PrefSrc                      string `ros:"pref-src"`
	RIP                          string `ros:"rip"`
	Slaac                        string `ros:"slaac"`
	Static                       string `ros:"static"`
	VPN                          string `ros:"vpn"`
	VRFInterface                 string `ros:"vrf-interface"`
}
```

### `/ipv6/settings`

- File: `ipv6_settings.json`
- Sumber field: live 19

Command: 

```go
// IPv6SettingsRow adalah kolom hasil `/ipv6/settings/print`.
type IPv6SettingsRow struct {
	AcceptRedirects              string `ros:"accept-redirects"`
	AcceptRouterAdvertisements   string `ros:"accept-router-advertisements"`
	AcceptRouterAdvertisementsOn string `ros:"accept-router-advertisements-on"`
	AllowFastPath                string `ros:"allow-fast-path"`
	DisableIPv6                  string `ros:"disable-ipv6"`
	DisableLinkLocalAddress      string `ros:"disable-link-local-address"`
	Forward                      string `ros:"forward"`
	IPv6FastPathActive           string `ros:"ipv6-fast-path-active"`
	IPv6FastPathBytes            string `ros:"ipv6-fast-path-bytes"`
	IPv6FastPathPackets          string `ros:"ipv6-fast-path-packets"`
	IPv6FasttrackActive          string `ros:"ipv6-fasttrack-active"`
	IPv6FasttrackBytes           string `ros:"ipv6-fasttrack-bytes"`
	IPv6FasttrackPackets         string `ros:"ipv6-fasttrack-packets"`
	MaxNeighborEntries           string `ros:"max-neighbor-entries"`
	MinNeighborEntries           string `ros:"min-neighbor-entries"`
	MultipathHashPolicy          string `ros:"multipath-hash-policy"`
	SoftMaxNeighborEntries       string `ros:"soft-max-neighbor-entries"`
	StaleNeighborDetectInterval  string `ros:"stale-neighbor-detect-interval"`
	StaleNeighborTimeout         string `ros:"stale-neighbor-timeout"`
}
```

### `/log`

- File: `log.json`
- Katalog: `log`
- Sumber field: live 5, completion 6, katalog 3

Command: `debug`(1), `error`(1), `find`(1), `get`(4), `info`(1), `print`(21), `warning`(1)

```go
// LogRow adalah kolom hasil `/log/print`.
type LogRow struct {
	ID                           string `ros:".id,key"`
	ExtraInfo                    string `ros:"extra-info"`
	Message                      string `ros:"message"`
	Time                         string `ros:"time"`
	Topics                       string `ros:"topics"`
	About                        string `ros:"about"`
	Buffer                       string `ros:"buffer"`
}
```

### `/mpls/forwarding-table`

- File: `mpls_forwarding-table.json`
- Sumber field: completion 12

Command: 

```go
// MPLSForwardingtableRow adalah kolom hasil `/mpls/forwarding-table/print`.
type MPLSForwardingtableRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Label                        string `ros:"label"`
	Ldp                          string `ros:"ldp"`
	Nexthops                     string `ros:"nexthops"`
	Prefix                       string `ros:"prefix"`
	TeSender                     string `ros:"te-sender"`
	TeSession                    string `ros:"te-session"`
	TrafficEng                   string `ros:"traffic-eng"`
	Type                         string `ros:"type"`
	VPLS                         string `ros:"vpls"`
	VPN                          string `ros:"vpn"`
	VRF                          string `ros:"vrf"`
}
```

### `/mpls/interface`

- File: `mpls_interface.json`
- Sumber field: completion 7

Command: 

```go
// MPLSInterfaceRow adalah kolom hasil `/mpls/interface/print`.
type MPLSInterfaceRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Builtin                      string `ros:"builtin"`
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	Input                        string `ros:"input"`
	Interface                    string `ros:"interface"`
	MPLSMTU                      string `ros:"mpls-mtu"`
}
```

### `/mpls/ldp`

- File: `mpls_ldp.json`
- Sumber field: completion 14

Command: 

```go
// MPLSLdpRow adalah kolom hasil `/mpls/ldp/print`.
type MPLSLdpRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Afi                          string `ros:"afi"`
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	DistributeForDefault         string `ros:"distribute-for-default"`
	HopLimit                     string `ros:"hop-limit"`
	Inactive                     string `ros:"inactive"`
	LoopDetect                   string `ros:"loop-detect"`
	LsrID                        string `ros:"lsr-id"`
	PathVectorLimit              string `ros:"path-vector-limit"`
	PreferredAfi                 string `ros:"preferred-afi"`
	TransportAddresses           string `ros:"transport-addresses"`
	UseExplicitNull              string `ros:"use-explicit-null"`
	VRF                          string `ros:"vrf"`
}
```

### `/mpls/mangle`

- File: `mpls_mangle.json`
- Sumber field: completion 8

Command: 

```go
// MPLSMangleRow adalah kolom hasil `/mpls/mangle/print`.
type MPLSMangleRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Chain                        string `ros:"chain"`
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	Exp                          string `ros:"exp"`
	Packets                      string `ros:"packets"`
	SetExp                       string `ros:"set-exp"`
	SetMark                      string `ros:"set-mark"`
}
```

### `/mpls/settings`

- File: `mpls_settings.json`
- Sumber field: live 5

Command: 

```go
// MPLSSettingsRow adalah kolom hasil `/mpls/settings/print`.
type MPLSSettingsRow struct {
	AllowFastPath                string `ros:"allow-fast-path"`
	DynamicLabelRange            string `ros:"dynamic-label-range"`
	MPLSFastPathBytes            string `ros:"mpls-fast-path-bytes"`
	MPLSFastPathPackets          string `ros:"mpls-fast-path-packets"`
	PropagateTTL                 string `ros:"propagate-ttl"`
}
```

### `/ping`

- File: `(dari katalog Go)`
- Katalog: `tool-ping`
- Sumber field: katalog 12

Command: 

```go
// PingRow adalah kolom hasil `/ping/print`.
type PingRow struct {
	Seq                          string `ros:"seq"`
	Host                         string `ros:"host"`
	Size                         string `ros:"size"`
	TTL                          string `ros:"ttl"`
	Time                         string `ros:"time"`
	Status                       string `ros:"status"`
	Sent                         string `ros:"sent"`
	Received                     string `ros:"received"`
	PacketLoss                   string `ros:"packet-loss"`
	MinRTT                       string `ros:"min-rtt"`
	AvgRTT                       string `ros:"avg-rtt"`
	MaxRTT                       string `ros:"max-rtt"`
}
```

### `/port/remote-access`

- File: `port_remote-access.json`
- Sumber field: completion 16

Command: 

```go
// PortRemoteaccessRow adalah kolom hasil `/port/remote-access/print`.
type PortRemoteaccessRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Active                       string `ros:"active"`
	ActivePeer                   string `ros:"active-peer"`
	ActivePeerPort               string `ros:"active-peer-port"`
	Busy                         string `ros:"busy"`
	Channel                      string `ros:"channel"`
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	Inactive                     string `ros:"inactive"`
	IPPort                       string `ros:"ip-port"`
	LocalAddress                 string `ros:"local-address"`
	LogFile                      string `ros:"log-file"`
	LoggingActive                string `ros:"logging-active"`
	Port                         string `ros:"port"`
	Protocol                     string `ros:"protocol"`
	RemoteAddresses              string `ros:"remote-addresses"`
}
```

### `/ppp`

- File: `ppp__cmds.json`
- Sumber field: tidak ada kolom terdeteksi

Command: `export`(8)

_Tidak ada struct: menu ini bukan tabel yang bisa di-print atau tidak tersedia di router uji (mis. paket tidak terpasang)._

### `/ppp/aaa`

- File: `ppp_aaa.json`
- Katalog: `ppp-aaa`
- Sumber field: live 5, katalog 4

Command: `edit`(1), `export`(8), `get`(3), `print`(5), `reset`(5), `set`(5)

```go
// PPPAaaRow adalah kolom hasil `/ppp/aaa/print`.
type PPPAaaRow struct {
	ID                           string `ros:".id,key"`
	Accounting                   string `ros:"accounting"`
	EnableIPv6Accounting         string `ros:"enable-ipv6-accounting"`
	// time interval
	InterimUpdate                string `ros:"interim-update"`
	UseCircuitIDInNasPortID      string `ros:"use-circuit-id-in-nas-port-id"`
	UseRadius                    string `ros:"use-radius"`
}
```

```go
// PPPAaaSet adalah argumen `/ppp/aaa/set`.
type PPPAaaSet struct {
	Accounting                   string `ros:"accounting,omitempty"`
	EnableIPv6Accounting         string `ros:"enable-ipv6-accounting,omitempty"`
	// time interval
	InterimUpdate                string `ros:"interim-update,omitempty"`
	UseCircuitIDInNasPortID      string `ros:"use-circuit-id-in-nas-port-id,omitempty"`
	UseRadius                    string `ros:"use-radius,omitempty"`
}
```

### `/ppp/active`

- File: `ppp_active.json`
- Katalog: `ppp-active`
- Sumber field: completion 12, katalog 10

Command: `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(2)

```go
// PPPActiveRow adalah kolom hasil `/ppp/active/print`.
type PPPActiveRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Address                      string `ros:"address"`
	CallerID                     string `ros:"caller-id"`
	// string value
	Comment                      string `ros:"comment"`
	Encoding                     string `ros:"encoding"`
	LimitBytesIn                 string `ros:"limit-bytes-in"`
	LimitBytesOut                string `ros:"limit-bytes-out"`
	Name                         string `ros:"name"`
	Radius                       string `ros:"radius"`
	Service                      string `ros:"service"`
	SessionID                    string `ros:"session-id"`
	Uptime                       string `ros:"uptime"`
}
```

### `/ppp/l2tp-secret`

- File: `ppp_l2tp-secret.json`
- Sumber field: completion 4

Command: 

```go
// PPPL2tpsecretRow adalah kolom hasil `/ppp/l2tp-secret/print`.
type PPPL2tpsecretRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Address                      string `ros:"address"`
	Comment                      string `ros:"comment"`
	Secret                       string `ros:"secret"`
}
```

### `/ppp/profile`

- File: `ppp_profile.json`
- Katalog: `ppp-profile`
- Sumber field: live 14, completion 39, katalog 14

Command: `add`(38), `comment`(2), `edit`(2), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(38), `set`(38), `unset`(2)

```go
// PPPProfileRow adalah kolom hasil `/ppp/profile/print`.
type PPPProfileRow struct {
	ID                           string `ros:".id,key"`
	AddressList                  string `ros:"address-list"`
	BridgeLearning               string `ros:"bridge-learning"`
	ChangeTcpMss                 string `ros:"change-tcp-mss"`
	Default                      string `ros:"default"`
	// string value
	Name                         string `ros:"name"`
	// string value
	OnDown                       string `ros:"on-down"`
	// string value
	OnUp                         string `ros:"on-up"`
	OnlyOne                      string `ros:"only-one"`
	UseCompression               string `ros:"use-compression"`
	UseEncryption                string `ros:"use-encryption"`
	UseIPv6                      string `ros:"use-ipv6"`
	UseMPLS                      string `ros:"use-mpls"`
	UseUpnp                      string `ros:"use-upnp"`
	About                        string `ros:"about"`
	Bridge                       string `ros:"bridge"`
	BridgeHorizon                string `ros:"bridge-horizon"`
	BridgePathCost               string `ros:"bridge-path-cost"`
	BridgePortPriority           string `ros:"bridge-port-priority"`
	BridgePortTrusted            string `ros:"bridge-port-trusted"`
	BridgePortVid                string `ros:"bridge-port-vid"`
	// string value
	Comment                      string `ros:"comment"`
	// time interval
	DHCPv6LeaseTime              string `ros:"dhcpv6-lease-time"`
	DHCPv6PdPool                 string `ros:"dhcpv6-pd-pool"`
	DHCPv6UseRadius              string `ros:"dhcpv6-use-radius"`
	// A.B.C.D (IP address)
	DNSServer                    string `ros:"dns-server"`
	// time interval
	IdleTimeout                  string `ros:"idle-timeout"`
	IncomingFilter               string `ros:"incoming-filter"`
	InsertQueueBefore            string `ros:"insert-queue-before"`
	InterfaceList                string `ros:"interface-list"`
	LocalAddress                 string `ros:"local-address"`
	OutgoingFilter               string `ros:"outgoing-filter"`
	ParentQueue                  string `ros:"parent-queue"`
	QueueType                    string `ros:"queue-type"`
	// string value
	RateLimit                    string `ros:"rate-limit"`
	RemoteAddress                string `ros:"remote-address"`
	RemoteIPv6PrefixPool         string `ros:"remote-ipv6-prefix-pool"`
	RemoteIPv6PrefixReuse        string `ros:"remote-ipv6-prefix-reuse"`
	// time interval
	SessionTimeout               string `ros:"session-timeout"`
	// A.B.C.D (IP address)
	WinsServer                   string `ros:"wins-server"`
}
```

```go
// PPPProfileAdd adalah argumen `/ppp/profile/add`.
type PPPProfileAdd struct {
	AddressList                  string `ros:"address-list,omitempty"`
	Bridge                       string `ros:"bridge,omitempty"`
	BridgeHorizon                string `ros:"bridge-horizon,omitempty"`
	BridgeLearning               string `ros:"bridge-learning,omitempty"`
	BridgePathCost               string `ros:"bridge-path-cost,omitempty"`
	BridgePortPriority           string `ros:"bridge-port-priority,omitempty"`
	BridgePortTrusted            string `ros:"bridge-port-trusted,omitempty"`
	BridgePortVid                string `ros:"bridge-port-vid,omitempty"`
	ChangeTcpMss                 string `ros:"change-tcp-mss,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	// time interval
	DHCPv6LeaseTime              string `ros:"dhcpv6-lease-time,omitempty"`
	DHCPv6PdPool                 string `ros:"dhcpv6-pd-pool,omitempty"`
	DHCPv6UseRadius              string `ros:"dhcpv6-use-radius,omitempty"`
	// A.B.C.D (IP address)
	DNSServer                    string `ros:"dns-server,omitempty"`
	// time interval
	IdleTimeout                  string `ros:"idle-timeout,omitempty"`
	IncomingFilter               string `ros:"incoming-filter,omitempty"`
	InsertQueueBefore            string `ros:"insert-queue-before,omitempty"`
	InterfaceList                string `ros:"interface-list,omitempty"`
	LocalAddress                 string `ros:"local-address,omitempty"`
	// string value
	Name                         string `ros:"name,required,omitempty"`
	// string value
	OnDown                       string `ros:"on-down,omitempty"`
	// string value
	OnUp                         string `ros:"on-up,omitempty"`
	OnlyOne                      string `ros:"only-one,omitempty"`
	OutgoingFilter               string `ros:"outgoing-filter,omitempty"`
	ParentQueue                  string `ros:"parent-queue,omitempty"`
	QueueType                    string `ros:"queue-type,omitempty"`
	// string value
	RateLimit                    string `ros:"rate-limit,omitempty"`
	RemoteAddress                string `ros:"remote-address,omitempty"`
	RemoteIPv6PrefixPool         string `ros:"remote-ipv6-prefix-pool,omitempty"`
	RemoteIPv6PrefixReuse        string `ros:"remote-ipv6-prefix-reuse,omitempty"`
	// time interval
	SessionTimeout               string `ros:"session-timeout,omitempty"`
	UseCompression               string `ros:"use-compression,omitempty"`
	UseEncryption                string `ros:"use-encryption,omitempty"`
	UseIPv6                      string `ros:"use-ipv6,omitempty"`
	UseMPLS                      string `ros:"use-mpls,omitempty"`
	UseUpnp                      string `ros:"use-upnp,omitempty"`
	// A.B.C.D (IP address)
	WinsServer                   string `ros:"wins-server,omitempty"`
}
```

```go
// PPPProfileSet adalah argumen `/ppp/profile/set`.
type PPPProfileSet struct {
	AddressList                  string `ros:"address-list,omitempty"`
	Bridge                       string `ros:"bridge,omitempty"`
	BridgeHorizon                string `ros:"bridge-horizon,omitempty"`
	BridgeLearning               string `ros:"bridge-learning,omitempty"`
	BridgePathCost               string `ros:"bridge-path-cost,omitempty"`
	BridgePortPriority           string `ros:"bridge-port-priority,omitempty"`
	BridgePortTrusted            string `ros:"bridge-port-trusted,omitempty"`
	BridgePortVid                string `ros:"bridge-port-vid,omitempty"`
	ChangeTcpMss                 string `ros:"change-tcp-mss,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	// time interval
	DHCPv6LeaseTime              string `ros:"dhcpv6-lease-time,omitempty"`
	DHCPv6PdPool                 string `ros:"dhcpv6-pd-pool,omitempty"`
	DHCPv6UseRadius              string `ros:"dhcpv6-use-radius,omitempty"`
	// A.B.C.D (IP address)
	DNSServer                    string `ros:"dns-server,omitempty"`
	// time interval
	IdleTimeout                  string `ros:"idle-timeout,omitempty"`
	IncomingFilter               string `ros:"incoming-filter,omitempty"`
	InsertQueueBefore            string `ros:"insert-queue-before,omitempty"`
	InterfaceList                string `ros:"interface-list,omitempty"`
	LocalAddress                 string `ros:"local-address,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	// string value
	OnDown                       string `ros:"on-down,omitempty"`
	// string value
	OnUp                         string `ros:"on-up,omitempty"`
	OnlyOne                      string `ros:"only-one,omitempty"`
	OutgoingFilter               string `ros:"outgoing-filter,omitempty"`
	ParentQueue                  string `ros:"parent-queue,omitempty"`
	QueueType                    string `ros:"queue-type,omitempty"`
	// string value
	RateLimit                    string `ros:"rate-limit,omitempty"`
	RemoteAddress                string `ros:"remote-address,omitempty"`
	RemoteIPv6PrefixPool         string `ros:"remote-ipv6-prefix-pool,omitempty"`
	RemoteIPv6PrefixReuse        string `ros:"remote-ipv6-prefix-reuse,omitempty"`
	// time interval
	SessionTimeout               string `ros:"session-timeout,omitempty"`
	UseCompression               string `ros:"use-compression,omitempty"`
	UseEncryption                string `ros:"use-encryption,omitempty"`
	UseIPv6                      string `ros:"use-ipv6,omitempty"`
	UseMPLS                      string `ros:"use-mpls,omitempty"`
	UseUpnp                      string `ros:"use-upnp,omitempty"`
	// A.B.C.D (IP address)
	WinsServer                   string `ros:"wins-server,omitempty"`
}
```

### `/ppp/secret`

- File: `ppp_secret.json`
- Katalog: `ppp-secret`
- Sumber field: completion 18, katalog 13

Command: `add`(15), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(22), `remove`(1), `reset`(15), `set`(15), `unset`(2)

```go
// PPPSecretRow adalah kolom hasil `/ppp/secret/print`.
type PPPSecretRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	// string value
	CallerID                     string `ros:"caller-id"`
	// string value
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	// string value
	IPv6Routes                   string `ros:"ipv6-routes"`
	LastCallerID                 string `ros:"last-caller-id"`
	LastDisconnectReason         string `ros:"last-disconnect-reason"`
	LastLoggedOut                string `ros:"last-logged-out"`
	LimitBytesIn                 string `ros:"limit-bytes-in"`
	LimitBytesOut                string `ros:"limit-bytes-out"`
	// A.B.C.D (IP address)
	LocalAddress                 string `ros:"local-address"`
	// string value, min length 1
	Name                         string `ros:"name"`
	// string value
	Password                     string `ros:"password"`
	Profile                      string `ros:"profile"`
	// A.B.C.D (IP address)
	RemoteAddress                string `ros:"remote-address"`
	// IPv6/0..128 (IPv6 prefix)
	RemoteIPv6Prefix             string `ros:"remote-ipv6-prefix"`
	// string value
	Routes                       string `ros:"routes"`
	Service                      string `ros:"service"`
}
```

```go
// PPPSecretAdd adalah argumen `/ppp/secret/add`.
type PPPSecretAdd struct {
	// string value
	CallerID                     string `ros:"caller-id,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	// string value
	IPv6Routes                   string `ros:"ipv6-routes,omitempty"`
	LimitBytesIn                 string `ros:"limit-bytes-in,omitempty"`
	LimitBytesOut                string `ros:"limit-bytes-out,omitempty"`
	// A.B.C.D (IP address)
	LocalAddress                 string `ros:"local-address,omitempty"`
	// string value, min length 1
	Name                         string `ros:"name,required,omitempty"`
	// string value
	Password                     string `ros:"password,required,omitempty"`
	Profile                      string `ros:"profile,omitempty"`
	// A.B.C.D (IP address)
	RemoteAddress                string `ros:"remote-address,omitempty"`
	// IPv6/0..128 (IPv6 prefix)
	RemoteIPv6Prefix             string `ros:"remote-ipv6-prefix,omitempty"`
	// string value
	Routes                       string `ros:"routes,omitempty"`
	Service                      string `ros:"service,omitempty"`
}
```

```go
// PPPSecretSet adalah argumen `/ppp/secret/set`.
type PPPSecretSet struct {
	// string value
	CallerID                     string `ros:"caller-id,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	// string value
	IPv6Routes                   string `ros:"ipv6-routes,omitempty"`
	LimitBytesIn                 string `ros:"limit-bytes-in,omitempty"`
	LimitBytesOut                string `ros:"limit-bytes-out,omitempty"`
	// A.B.C.D (IP address)
	LocalAddress                 string `ros:"local-address,omitempty"`
	// string value, min length 1
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	// string value
	Password                     string `ros:"password,omitempty"`
	Profile                      string `ros:"profile,omitempty"`
	// A.B.C.D (IP address)
	RemoteAddress                string `ros:"remote-address,omitempty"`
	// IPv6/0..128 (IPv6 prefix)
	RemoteIPv6Prefix             string `ros:"remote-ipv6-prefix,omitempty"`
	// string value
	Routes                       string `ros:"routes,omitempty"`
	Service                      string `ros:"service,omitempty"`
}
```

### `/queue`

- File: `queue__cmds.json`
- Sumber field: tidak ada kolom terdeteksi

Command: `export`(8), `monitor`(10)

_Tidak ada struct: menu ini bukan tabel yang bisa di-print atau tidak tersedia di router uji (mis. paket tidak terpasang)._

### `/queue/interface`

- File: `queue_interface.json`
- Sumber field: live 5, completion 5

Command: `edit`(2), `export`(8), `find`(1), `get`(4), `print`(21), `reset`(2), `set`(2)

```go
// QueueInterfaceRow adalah kolom hasil `/queue/interface/print`.
type QueueInterfaceRow struct {
	ID                           string `ros:".id,key"`
	ActiveQueue                  string `ros:"active-queue"`
	DefaultQueue                 string `ros:"default-queue"`
	Interface                    string `ros:"interface"`
	Queue                        string `ros:"queue"`
	About                        string `ros:"about"`
}
```

```go
// QueueInterfaceSet adalah argumen `/queue/interface/set`.
type QueueInterfaceSet struct {
	Numbers                      string `ros:"numbers,omitempty"`
	Queue                        string `ros:"queue,omitempty"`
}
```

### `/queue/simple`

- File: `queue_simple.json`
- Katalog: `queue-simple`
- Sumber field: completion 43, katalog 18

Command: `add`(26), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `move`(2), `print`(26), `remove`(1), `reset`(25), `reset-counters`(1), `reset-counters-all`(0), `set`(25), `unset`(2)

```go
// QueueSimpleRow adalah kolom hasil `/queue/simple/print`.
type QueueSimpleRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	BucketSize                   string `ros:"bucket-size"`
	BurstLimit                   string `ros:"burst-limit"`
	BurstThreshold               string `ros:"burst-threshold"`
	// time interval
	BurstTime                    string `ros:"burst-time"`
	Bytes                        string `ros:"bytes"`
	// string value
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	Dropped                      string `ros:"dropped"`
	Dst                          string `ros:"dst"`
	Dynamic                      string `ros:"dynamic"`
	Invalid                      string `ros:"invalid"`
	LimitAt                      string `ros:"limit-at"`
	MaxLimit                     string `ros:"max-limit"`
	// string value
	Name                         string `ros:"name"`
	PacketMarks                  string `ros:"packet-marks"`
	PacketRate                   string `ros:"packet-rate"`
	Packets                      string `ros:"packets"`
	Parent                       string `ros:"parent"`
	PcqQueues                    string `ros:"pcq-queues"`
	Priority                     string `ros:"priority"`
	Queue                        string `ros:"queue"`
	QueuedBytes                  string `ros:"queued-bytes"`
	QueuedPackets                string `ros:"queued-packets"`
	Rate                         string `ros:"rate"`
	Target                       string `ros:"target"`
	Time                         string `ros:"time"`
	TotalBucketSize              string `ros:"total-bucket-size"`
	TotalBurstLimit              string `ros:"total-burst-limit"`
	TotalBurstThreshold          string `ros:"total-burst-threshold"`
	// time interval
	TotalBurstTime               string `ros:"total-burst-time"`
	TotalBytes                   string `ros:"total-bytes"`
	TotalDropped                 string `ros:"total-dropped"`
	TotalLimitAt                 string `ros:"total-limit-at"`
	TotalMaxLimit                string `ros:"total-max-limit"`
	TotalPacketRate              string `ros:"total-packet-rate"`
	TotalPackets                 string `ros:"total-packets"`
	TotalPcqQueues               string `ros:"total-pcq-queues"`
	TotalPriority                string `ros:"total-priority"`
	TotalQueue                   string `ros:"total-queue"`
	TotalQueuedBytes             string `ros:"total-queued-bytes"`
	TotalQueuedPackets           string `ros:"total-queued-packets"`
	TotalRate                    string `ros:"total-rate"`
}
```

```go
// QueueSimpleAdd adalah argumen `/queue/simple/add`.
type QueueSimpleAdd struct {
	BucketSize                   string `ros:"bucket-size,omitempty"`
	BurstLimit                   string `ros:"burst-limit,omitempty"`
	BurstThreshold               string `ros:"burst-threshold,omitempty"`
	// time interval
	BurstTime                    string `ros:"burst-time,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Dst                          string `ros:"dst,omitempty"`
	LimitAt                      string `ros:"limit-at,omitempty"`
	MaxLimit                     string `ros:"max-limit,omitempty"`
	// string value
	Name                         string `ros:"name,required,omitempty"`
	PacketMarks                  string `ros:"packet-marks,omitempty"`
	Parent                       string `ros:"parent,omitempty"`
	PlaceBefore                  string `ros:"place-before,omitempty"`
	Priority                     string `ros:"priority,omitempty"`
	Queue                        string `ros:"queue,omitempty"`
	Target                       string `ros:"target,required,omitempty"`
	Time                         string `ros:"time,omitempty"`
	TotalBucketSize              string `ros:"total-bucket-size,omitempty"`
	TotalBurstLimit              string `ros:"total-burst-limit,omitempty"`
	TotalBurstThreshold          string `ros:"total-burst-threshold,omitempty"`
	// time interval
	TotalBurstTime               string `ros:"total-burst-time,omitempty"`
	TotalLimitAt                 string `ros:"total-limit-at,omitempty"`
	TotalMaxLimit                string `ros:"total-max-limit,omitempty"`
	TotalPriority                string `ros:"total-priority,omitempty"`
	TotalQueue                   string `ros:"total-queue,omitempty"`
}
```

```go
// QueueSimpleSet adalah argumen `/queue/simple/set`.
type QueueSimpleSet struct {
	BucketSize                   string `ros:"bucket-size,omitempty"`
	BurstLimit                   string `ros:"burst-limit,omitempty"`
	BurstThreshold               string `ros:"burst-threshold,omitempty"`
	// time interval
	BurstTime                    string `ros:"burst-time,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Dst                          string `ros:"dst,omitempty"`
	LimitAt                      string `ros:"limit-at,omitempty"`
	MaxLimit                     string `ros:"max-limit,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	PacketMarks                  string `ros:"packet-marks,omitempty"`
	Parent                       string `ros:"parent,omitempty"`
	Priority                     string `ros:"priority,omitempty"`
	Queue                        string `ros:"queue,omitempty"`
	Target                       string `ros:"target,omitempty"`
	Time                         string `ros:"time,omitempty"`
	TotalBucketSize              string `ros:"total-bucket-size,omitempty"`
	TotalBurstLimit              string `ros:"total-burst-limit,omitempty"`
	TotalBurstThreshold          string `ros:"total-burst-threshold,omitempty"`
	// time interval
	TotalBurstTime               string `ros:"total-burst-time,omitempty"`
	TotalLimitAt                 string `ros:"total-limit-at,omitempty"`
	TotalMaxLimit                string `ros:"total-max-limit,omitempty"`
	TotalPriority                string `ros:"total-priority,omitempty"`
	TotalQueue                   string `ros:"total-queue,omitempty"`
}
```

### `/queue/tree`

- File: `queue_tree.json`
- Katalog: `queue-tree`
- Sumber field: completion 23, katalog 17

Command: `add`(14), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(26), `remove`(1), `reset`(14), `reset-counters`(1), `reset-counters-all`(0), `set`(14)

```go
// QueueTreeRow adalah kolom hasil `/queue/tree/print`.
type QueueTreeRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	BucketSize                   string `ros:"bucket-size"`
	BurstLimit                   string `ros:"burst-limit"`
	BurstThreshold               string `ros:"burst-threshold"`
	// time interval
	BurstTime                    string `ros:"burst-time"`
	Bytes                        string `ros:"bytes"`
	// string value
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	Dropped                      string `ros:"dropped"`
	Invalid                      string `ros:"invalid"`
	LimitAt                      string `ros:"limit-at"`
	MaxLimit                     string `ros:"max-limit"`
	// string value
	Name                         string `ros:"name"`
	PacketMark                   string `ros:"packet-mark"`
	PacketRate                   string `ros:"packet-rate"`
	Packets                      string `ros:"packets"`
	Parent                       string `ros:"parent"`
	PcqQueues                    string `ros:"pcq-queues"`
	Priority                     string `ros:"priority"`
	Queue                        string `ros:"queue"`
	QueuedBytes                  string `ros:"queued-bytes"`
	QueuedPackets                string `ros:"queued-packets"`
	Rate                         string `ros:"rate"`
}
```

```go
// QueueTreeAdd adalah argumen `/queue/tree/add`.
type QueueTreeAdd struct {
	BucketSize                   string `ros:"bucket-size,omitempty"`
	BurstLimit                   string `ros:"burst-limit,omitempty"`
	BurstThreshold               string `ros:"burst-threshold,omitempty"`
	// time interval
	BurstTime                    string `ros:"burst-time,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	LimitAt                      string `ros:"limit-at,omitempty"`
	MaxLimit                     string `ros:"max-limit,omitempty"`
	// string value
	Name                         string `ros:"name,required,omitempty"`
	PacketMark                   string `ros:"packet-mark,omitempty"`
	Parent                       string `ros:"parent,required,omitempty"`
	Priority                     string `ros:"priority,omitempty"`
	Queue                        string `ros:"queue,omitempty"`
}
```

```go
// QueueTreeSet adalah argumen `/queue/tree/set`.
type QueueTreeSet struct {
	BucketSize                   string `ros:"bucket-size,omitempty"`
	BurstLimit                   string `ros:"burst-limit,omitempty"`
	BurstThreshold               string `ros:"burst-threshold,omitempty"`
	// time interval
	BurstTime                    string `ros:"burst-time,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	LimitAt                      string `ros:"limit-at,omitempty"`
	MaxLimit                     string `ros:"max-limit,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	PacketMark                   string `ros:"packet-mark,omitempty"`
	Parent                       string `ros:"parent,omitempty"`
	Priority                     string `ros:"priority,omitempty"`
	Queue                        string `ros:"queue,omitempty"`
}
```

### `/queue/type`

- File: `queue_type.json`
- Sumber field: live 24, completion 52

Command: `add`(51), `edit`(2), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(51), `set`(51)

```go
// QueueTypeRow adalah kolom hasil `/queue/type/print`.
type QueueTypeRow struct {
	ID                           string `ros:".id,key"`
	Default                      string `ros:"default"`
	Kind                         string `ros:"kind"`
	MqPfifoLimit                 string `ros:"mq-pfifo-limit"`
	// string value
	Name                         string `ros:"name"`
	PcqBurstRate                 string `ros:"pcq-burst-rate"`
	PcqBurstThreshold            string `ros:"pcq-burst-threshold"`
	// 1s..7101w3d6h28m15s (time interval)
	PcqBurstTime                 string `ros:"pcq-burst-time"`
	// src-address|dst-address|src-port|dst-port[,PcqClassifier*]
	PcqClassifier                string `ros:"pcq-classifier"`
	PcqDstAddressMask            string `ros:"pcq-dst-address-mask"`
	PcqDstAddress6Mask           string `ros:"pcq-dst-address6-mask"`
	PcqLimit                     string `ros:"pcq-limit"`
	PcqRate                      string `ros:"pcq-rate"`
	PcqSrcAddressMask            string `ros:"pcq-src-address-mask"`
	PcqSrcAddress6Mask           string `ros:"pcq-src-address6-mask"`
	PcqTotalLimit                string `ros:"pcq-total-limit"`
	PfifoLimit                   string `ros:"pfifo-limit"`
	RedAvgPacket                 string `ros:"red-avg-packet"`
	RedBurst                     string `ros:"red-burst"`
	RedLimit                     string `ros:"red-limit"`
	RedMaxThreshold              string `ros:"red-max-threshold"`
	RedMinThreshold              string `ros:"red-min-threshold"`
	SfqAllot                     string `ros:"sfq-allot"`
	SfqPerturb                   string `ros:"sfq-perturb"`
	About                        string `ros:"about"`
	BfifoLimit                   string `ros:"bfifo-limit"`
	CakeAckFilter                string `ros:"cake-ack-filter"`
	CakeAtm                      string `ros:"cake-atm"`
	CakeAutorateIngress          string `ros:"cake-autorate-ingress"`
	CakeBandwidth                string `ros:"cake-bandwidth"`
	CakeDiffserv                 string `ros:"cake-diffserv"`
	CakeFlowmode                 string `ros:"cake-flowmode"`
	CakeMemlimit                 string `ros:"cake-memlimit"`
	CakeMpu                      string `ros:"cake-mpu"`
	CakeNAT                      string `ros:"cake-nat"`
	CakeOverhead                 string `ros:"cake-overhead"`
	CakeOverheadScheme           string `ros:"cake-overhead-scheme"`
	// time interval
	CakeRTT                      string `ros:"cake-rtt"`
	CakeRTTScheme                string `ros:"cake-rtt-scheme"`
	CakeWash                     string `ros:"cake-wash"`
	// time interval
	CodelCeThreshold             string `ros:"codel-ce-threshold"`
	CodelEcn                     string `ros:"codel-ecn"`
	// time interval
	CodelInterval                string `ros:"codel-interval"`
	CodelLimit                   string `ros:"codel-limit"`
	// time interval
	CodelTarget                  string `ros:"codel-target"`
	// time interval
	FqCodelCeThreshold           string `ros:"fq-codel-ce-threshold"`
	FqCodelEcn                   string `ros:"fq-codel-ecn"`
	FqCodelFlows                 string `ros:"fq-codel-flows"`
	// time interval
	FqCodelInterval              string `ros:"fq-codel-interval"`
	FqCodelLimit                 string `ros:"fq-codel-limit"`
	FqCodelMemlimit              string `ros:"fq-codel-memlimit"`
	FqCodelQuantum               string `ros:"fq-codel-quantum"`
	// time interval
	FqCodelTarget                string `ros:"fq-codel-target"`
}
```

```go
// QueueTypeAdd adalah argumen `/queue/type/add`.
type QueueTypeAdd struct {
	BfifoLimit                   string `ros:"bfifo-limit,omitempty"`
	CakeAckFilter                string `ros:"cake-ack-filter,omitempty"`
	CakeAtm                      string `ros:"cake-atm,omitempty"`
	CakeAutorateIngress          string `ros:"cake-autorate-ingress,omitempty"`
	CakeBandwidth                string `ros:"cake-bandwidth,omitempty"`
	CakeDiffserv                 string `ros:"cake-diffserv,omitempty"`
	CakeFlowmode                 string `ros:"cake-flowmode,omitempty"`
	CakeMemlimit                 string `ros:"cake-memlimit,omitempty"`
	CakeMpu                      string `ros:"cake-mpu,omitempty"`
	CakeNAT                      string `ros:"cake-nat,omitempty"`
	CakeOverhead                 string `ros:"cake-overhead,omitempty"`
	CakeOverheadScheme           string `ros:"cake-overhead-scheme,omitempty"`
	// time interval
	CakeRTT                      string `ros:"cake-rtt,omitempty"`
	CakeRTTScheme                string `ros:"cake-rtt-scheme,omitempty"`
	CakeWash                     string `ros:"cake-wash,omitempty"`
	// time interval
	CodelCeThreshold             string `ros:"codel-ce-threshold,omitempty"`
	CodelEcn                     string `ros:"codel-ecn,omitempty"`
	// time interval
	CodelInterval                string `ros:"codel-interval,omitempty"`
	CodelLimit                   string `ros:"codel-limit,omitempty"`
	// time interval
	CodelTarget                  string `ros:"codel-target,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	// time interval
	FqCodelCeThreshold           string `ros:"fq-codel-ce-threshold,omitempty"`
	FqCodelEcn                   string `ros:"fq-codel-ecn,omitempty"`
	FqCodelFlows                 string `ros:"fq-codel-flows,omitempty"`
	// time interval
	FqCodelInterval              string `ros:"fq-codel-interval,omitempty"`
	FqCodelLimit                 string `ros:"fq-codel-limit,omitempty"`
	FqCodelMemlimit              string `ros:"fq-codel-memlimit,omitempty"`
	FqCodelQuantum               string `ros:"fq-codel-quantum,omitempty"`
	// time interval
	FqCodelTarget                string `ros:"fq-codel-target,omitempty"`
	Kind                         string `ros:"kind,omitempty"`
	MqPfifoLimit                 string `ros:"mq-pfifo-limit,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	PcqBurstRate                 string `ros:"pcq-burst-rate,omitempty"`
	PcqBurstThreshold            string `ros:"pcq-burst-threshold,omitempty"`
	// 1s..7101w3d6h28m15s (time interval)
	PcqBurstTime                 string `ros:"pcq-burst-time,omitempty"`
	// src-address|dst-address|src-port|dst-port[,PcqClassifier*]
	PcqClassifier                string `ros:"pcq-classifier,omitempty"`
	PcqDstAddressMask            string `ros:"pcq-dst-address-mask,omitempty"`
	PcqDstAddress6Mask           string `ros:"pcq-dst-address6-mask,omitempty"`
	PcqLimit                     string `ros:"pcq-limit,omitempty"`
	PcqRate                      string `ros:"pcq-rate,omitempty"`
	PcqSrcAddressMask            string `ros:"pcq-src-address-mask,omitempty"`
	PcqSrcAddress6Mask           string `ros:"pcq-src-address6-mask,omitempty"`
	PcqTotalLimit                string `ros:"pcq-total-limit,omitempty"`
	PfifoLimit                   string `ros:"pfifo-limit,omitempty"`
	RedAvgPacket                 string `ros:"red-avg-packet,omitempty"`
	RedBurst                     string `ros:"red-burst,omitempty"`
	RedLimit                     string `ros:"red-limit,omitempty"`
	RedMaxThreshold              string `ros:"red-max-threshold,omitempty"`
	RedMinThreshold              string `ros:"red-min-threshold,omitempty"`
	SfqAllot                     string `ros:"sfq-allot,omitempty"`
	SfqPerturb                   string `ros:"sfq-perturb,omitempty"`
}
```

```go
// QueueTypeSet adalah argumen `/queue/type/set`.
type QueueTypeSet struct {
	BfifoLimit                   string `ros:"bfifo-limit,omitempty"`
	CakeAckFilter                string `ros:"cake-ack-filter,omitempty"`
	CakeAtm                      string `ros:"cake-atm,omitempty"`
	CakeAutorateIngress          string `ros:"cake-autorate-ingress,omitempty"`
	CakeBandwidth                string `ros:"cake-bandwidth,omitempty"`
	CakeDiffserv                 string `ros:"cake-diffserv,omitempty"`
	CakeFlowmode                 string `ros:"cake-flowmode,omitempty"`
	CakeMemlimit                 string `ros:"cake-memlimit,omitempty"`
	CakeMpu                      string `ros:"cake-mpu,omitempty"`
	CakeNAT                      string `ros:"cake-nat,omitempty"`
	CakeOverhead                 string `ros:"cake-overhead,omitempty"`
	CakeOverheadScheme           string `ros:"cake-overhead-scheme,omitempty"`
	// time interval
	CakeRTT                      string `ros:"cake-rtt,omitempty"`
	CakeRTTScheme                string `ros:"cake-rtt-scheme,omitempty"`
	CakeWash                     string `ros:"cake-wash,omitempty"`
	// time interval
	CodelCeThreshold             string `ros:"codel-ce-threshold,omitempty"`
	CodelEcn                     string `ros:"codel-ecn,omitempty"`
	// time interval
	CodelInterval                string `ros:"codel-interval,omitempty"`
	CodelLimit                   string `ros:"codel-limit,omitempty"`
	// time interval
	CodelTarget                  string `ros:"codel-target,omitempty"`
	// time interval
	FqCodelCeThreshold           string `ros:"fq-codel-ce-threshold,omitempty"`
	FqCodelEcn                   string `ros:"fq-codel-ecn,omitempty"`
	FqCodelFlows                 string `ros:"fq-codel-flows,omitempty"`
	// time interval
	FqCodelInterval              string `ros:"fq-codel-interval,omitempty"`
	FqCodelLimit                 string `ros:"fq-codel-limit,omitempty"`
	FqCodelMemlimit              string `ros:"fq-codel-memlimit,omitempty"`
	FqCodelQuantum               string `ros:"fq-codel-quantum,omitempty"`
	// time interval
	FqCodelTarget                string `ros:"fq-codel-target,omitempty"`
	Kind                         string `ros:"kind,omitempty"`
	MqPfifoLimit                 string `ros:"mq-pfifo-limit,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	PcqBurstRate                 string `ros:"pcq-burst-rate,omitempty"`
	PcqBurstThreshold            string `ros:"pcq-burst-threshold,omitempty"`
	// 1s..7101w3d6h28m15s (time interval)
	PcqBurstTime                 string `ros:"pcq-burst-time,omitempty"`
	// src-address|dst-address|src-port|dst-port[,PcqClassifier*]
	PcqClassifier                string `ros:"pcq-classifier,omitempty"`
	PcqDstAddressMask            string `ros:"pcq-dst-address-mask,omitempty"`
	PcqDstAddress6Mask           string `ros:"pcq-dst-address6-mask,omitempty"`
	PcqLimit                     string `ros:"pcq-limit,omitempty"`
	PcqRate                      string `ros:"pcq-rate,omitempty"`
	PcqSrcAddressMask            string `ros:"pcq-src-address-mask,omitempty"`
	PcqSrcAddress6Mask           string `ros:"pcq-src-address6-mask,omitempty"`
	PcqTotalLimit                string `ros:"pcq-total-limit,omitempty"`
	PfifoLimit                   string `ros:"pfifo-limit,omitempty"`
	RedAvgPacket                 string `ros:"red-avg-packet,omitempty"`
	RedBurst                     string `ros:"red-burst,omitempty"`
	RedLimit                     string `ros:"red-limit,omitempty"`
	RedMaxThreshold              string `ros:"red-max-threshold,omitempty"`
	RedMinThreshold              string `ros:"red-min-threshold,omitempty"`
	SfqAllot                     string `ros:"sfq-allot,omitempty"`
	SfqPerturb                   string `ros:"sfq-perturb,omitempty"`
}
```

### `/radius`

- File: `(dari katalog Go)`
- Katalog: `radius`
- Sumber field: katalog 19

Command: `add`(10), `set`(10)

```go
// RadiusRow adalah kolom hasil `/radius/print`.
type RadiusRow struct {
	ID                           string `ros:".id,key"`
	Service                      string `ros:"service"`
	Address                      string `ros:"address"`
	AuthenticationPort           string `ros:"authentication-port"`
	AccountingPort               string `ros:"accounting-port"`
	Timeout                      string `ros:"timeout"`
	Realm                        string `ros:"realm"`
	SrcAddress                   string `ros:"src-address"`
	AccountingBackup             string `ros:"accounting-backup"`
	Status                       string `ros:"status"`
	Disabled                     string `ros:"disabled"`
	Comment                      string `ros:"comment"`
	Pending                      string `ros:"pending"`
	Requests                     string `ros:"requests"`
	Accepts                      string `ros:"accepts"`
	Rejects                      string `ros:"rejects"`
	Resends                      string `ros:"resends"`
	Timeouts                     string `ros:"timeouts"`
	BadReplies                   string `ros:"bad-replies"`
	LastRequestRTT               string `ros:"last-request-rtt"`
}
```

```go
// RadiusAdd adalah argumen `/radius/add`.
type RadiusAdd struct {
	Service                      string `ros:"service,required,omitempty"`
	Address                      string `ros:"address,required,omitempty"`
	Secret                       string `ros:"secret,required,omitempty"`
	AuthenticationPort           string `ros:"authentication-port,omitempty"`
	AccountingPort               string `ros:"accounting-port,omitempty"`
	Timeout                      string `ros:"timeout,omitempty"`
	Realm                        string `ros:"realm,omitempty"`
	SrcAddress                   string `ros:"src-address,omitempty"`
	AccountingBackup             string `ros:"accounting-backup,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

```go
// RadiusSet adalah argumen `/radius/set`.
type RadiusSet struct {
	Service                      string `ros:"service,omitempty"`
	Address                      string `ros:"address,omitempty"`
	Secret                       string `ros:"secret,omitempty"`
	AuthenticationPort           string `ros:"authentication-port,omitempty"`
	AccountingPort               string `ros:"accounting-port,omitempty"`
	Timeout                      string `ros:"timeout,omitempty"`
	Realm                        string `ros:"realm,omitempty"`
	SrcAddress                   string `ros:"src-address,omitempty"`
	AccountingBackup             string `ros:"accounting-backup,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
}
```

### `/radius/incoming`

- File: `radius_incoming.json`
- Katalog: `radius-incoming`
- Sumber field: live 3, katalog 2

Command: 

```go
// RadiusIncomingRow adalah kolom hasil `/radius/incoming/print`.
type RadiusIncomingRow struct {
	Accept                       string `ros:"accept"`
	Port                         string `ros:"port"`
	VRF                          string `ros:"vrf"`
}
```

### `/routing`

- File: `routing__cmds.json`
- Sumber field: tidak ada kolom terdeteksi

Command: `discourse`(5), `export`(8)

_Tidak ada struct: menu ini bukan tabel yang bisa di-print atau tidak tersedia di router uji (mis. paket tidak terpasang)._

### `/routing/bfd`

- File: `routing_bfd.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command

Command: `export`(8)

_Tidak ada struct: menu ini bukan tabel yang bisa di-print atau tidak tersedia di router uji (mis. paket tidak terpasang)._

### `/routing/bgp`

- File: `routing_bgp.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command

Command: `export`(8)

_Tidak ada struct: menu ini bukan tabel yang bisa di-print atau tidak tersedia di router uji (mis. paket tidak terpasang)._

### `/routing/fantasy`

- File: `routing_fantasy.json`
- Sumber field: completion 18

Command: `add`(17), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(20), `remove`(1), `reset`(17), `set`(17)

```go
// RoutingFantasyRow adalah kolom hasil `/routing/fantasy/print`.
type RoutingFantasyRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	// string value
	Comment                      string `ros:"comment"`
	Count                        string `ros:"count"`
	DealerID                     string `ros:"dealer-id"`
	Disabled                     string `ros:"disabled"`
	DstAddress                   string `ros:"dst-address"`
	Gateway                      string `ros:"gateway"`
	InstanceID                   string `ros:"instance-id"`
	Invalid                      string `ros:"invalid"`
	// string value
	Name                         string `ros:"name"`
	Offset                       string `ros:"offset"`
	PrefixLength                 string `ros:"prefix-length"`
	PrivOffs                     string `ros:"priv-offs"`
	PrivSize                     string `ros:"priv-size"`
	Scope                        string `ros:"scope"`
	// string value
	Seed                         string `ros:"seed"`
	TargetScope                  string `ros:"target-scope"`
	UseHold                      string `ros:"use-hold"`
}
```

```go
// RoutingFantasyAdd adalah argumen `/routing/fantasy/add`.
type RoutingFantasyAdd struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Count                        string `ros:"count,omitempty"`
	DealerID                     string `ros:"dealer-id,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	DstAddress                   string `ros:"dst-address,omitempty"`
	Gateway                      string `ros:"gateway,omitempty"`
	InstanceID                   string `ros:"instance-id,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Offset                       string `ros:"offset,omitempty"`
	PrefixLength                 string `ros:"prefix-length,omitempty"`
	PrivOffs                     string `ros:"priv-offs,omitempty"`
	PrivSize                     string `ros:"priv-size,omitempty"`
	Scope                        string `ros:"scope,omitempty"`
	// string value
	Seed                         string `ros:"seed,omitempty"`
	TargetScope                  string `ros:"target-scope,omitempty"`
	UseHold                      string `ros:"use-hold,omitempty"`
}
```

```go
// RoutingFantasySet adalah argumen `/routing/fantasy/set`.
type RoutingFantasySet struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Count                        string `ros:"count,omitempty"`
	DealerID                     string `ros:"dealer-id,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	DstAddress                   string `ros:"dst-address,omitempty"`
	Gateway                      string `ros:"gateway,omitempty"`
	InstanceID                   string `ros:"instance-id,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	Offset                       string `ros:"offset,omitempty"`
	PrefixLength                 string `ros:"prefix-length,omitempty"`
	PrivOffs                     string `ros:"priv-offs,omitempty"`
	PrivSize                     string `ros:"priv-size,omitempty"`
	Scope                        string `ros:"scope,omitempty"`
	// string value
	Seed                         string `ros:"seed,omitempty"`
	TargetScope                  string `ros:"target-scope,omitempty"`
	UseHold                      string `ros:"use-hold,omitempty"`
}
```

### `/routing/filter`

- File: `routing_filter.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command

Command: `export`(8), `filter-wizard`(54), `sync`(1), `test-as-path-regexp`(3)

_Tidak ada struct: menu ini bukan tabel yang bisa di-print atau tidak tersedia di router uji (mis. paket tidak terpasang)._

### `/routing/gmp`

- File: `routing_gmp.json`
- Sumber field: completion 7

Command: `add`(7), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(20), `remove`(1), `reset`(7), `set`(7), `unset`(2)

```go
// RoutingGmpRow adalah kolom hasil `/routing/gmp/print`.
type RoutingGmpRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	// string value
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	Exclude                      string `ros:"exclude"`
	Groups                       string `ros:"groups"`
	Interfaces                   string `ros:"interfaces"`
	Sources                      string `ros:"sources"`
}
```

```go
// RoutingGmpAdd adalah argumen `/routing/gmp/add`.
type RoutingGmpAdd struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Exclude                      string `ros:"exclude,omitempty"`
	Groups                       string `ros:"groups,omitempty"`
	Interfaces                   string `ros:"interfaces,omitempty"`
	Sources                      string `ros:"sources,omitempty"`
}
```

```go
// RoutingGmpSet adalah argumen `/routing/gmp/set`.
type RoutingGmpSet struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Exclude                      string `ros:"exclude,omitempty"`
	Groups                       string `ros:"groups,omitempty"`
	Interfaces                   string `ros:"interfaces,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	Sources                      string `ros:"sources,omitempty"`
}
```

### `/routing/id`

- File: `routing_id.json`
- Sumber field: live 7, completion 10

Command: `add`(7), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(7), `set`(7), `unset`(2)

```go
// RoutingIDRow adalah kolom hasil `/routing/id/print`.
type RoutingIDRow struct {
	ID                           string `ros:".id,key"`
	Dynamic                      string `ros:"dynamic"`
	DynamicID                    string `ros:"dynamic-id"`
	Inactive                     string `ros:"inactive"`
	// string value
	Name                         string `ros:"name"`
	// any|only-static|only-loopback|only-vrf|only-active|lowest[,SelectDynamicId*]
	SelectDynamicID              string `ros:"select-dynamic-id"`
	SelectFromVRF                string `ros:"select-from-vrf"`
	About                        string `ros:"about"`
	// string value
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	// A.B.C.D (IP address)
	ID2                          string `ros:"id"`
}
```

```go
// RoutingIDAdd adalah argumen `/routing/id/add`.
type RoutingIDAdd struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	// A.B.C.D (IP address)
	ID                           string `ros:"id,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	// any|only-static|only-loopback|only-vrf|only-active|lowest[,SelectDynamicId*]
	SelectDynamicID              string `ros:"select-dynamic-id,omitempty"`
	SelectFromVRF                string `ros:"select-from-vrf,omitempty"`
}
```

```go
// RoutingIDSet adalah argumen `/routing/id/set`.
type RoutingIDSet struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	// A.B.C.D (IP address)
	ID                           string `ros:"id,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	// any|only-static|only-loopback|only-vrf|only-active|lowest[,SelectDynamicId*]
	SelectDynamicID              string `ros:"select-dynamic-id,omitempty"`
	SelectFromVRF                string `ros:"select-from-vrf,omitempty"`
}
```

### `/routing/igmp-proxy`

- File: `routing_igmp-proxy.json`
- Sumber field: live 3

Command: 

```go
// RoutingIgmpproxyRow adalah kolom hasil `/routing/igmp-proxy/print`.
type RoutingIgmpproxyRow struct {
	QueryInterval                string `ros:"query-interval"`
	QueryResponseInterval        string `ros:"query-response-interval"`
	QuickLeave                   string `ros:"quick-leave"`
}
```

### `/routing/nexthop`

- File: `routing_nexthop.json`
- Sumber field: live 18, completion 23

Command: 

```go
// RoutingNexthopRow adalah kolom hasil `/routing/nexthop/print`.
type RoutingNexthopRow struct {
	ID                           string `ros:".id,key"`
	Address                      string `ros:"address"`
	Afi                          string `ros:"afi"`
	BGPVPN                       string `ros:"bgp-vpn"`
	CheckGateway                 string `ros:"check-gateway"`
	GwCheckOk                    string `ros:"gw-check-ok"`
	ImmediateGwAddress           string `ros:"immediate-gw.address"`
	ImmediateGwBlackhole         string `ros:"immediate-gw.blackhole"`
	ImmediateGwFlapCount         string `ros:"immediate-gw.flap-count"`
	ImmediateGwInterfaceIdx      string `ros:"immediate-gw.interface-idx"`
	ImmediateGwMPLSLabel         string `ros:"immediate-gw.mpls-label"`
	ImmediateGwMPLSPeerID        string `ros:"immediate-gw.mpls-peer-id"`
	ImmediateGwWeight            string `ros:"immediate-gw.weight"`
	InterfaceOk                  string `ros:"interface-ok"`
	Reachable                    string `ros:"reachable"`
	Scope                        string `ros:"scope"`
	TargetScope                  string `ros:"target-scope"`
	Unresolved                   string `ros:"unresolved"`
	About                        string `ros:"about"`
	GwState                      string `ros:"gw-state"`
	Hops                         string `ros:"hops"`
	Links                        string `ros:"links"`
	PrevReachable                string `ros:"prev-reachable"`
	Routes                       string `ros:"routes"`
}
```

### `/routing/ospf`

- File: `routing_ospf.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command

Command: `export`(8)

_Tidak ada struct: menu ini bukan tabel yang bisa di-print atau tidak tersedia di router uji (mis. paket tidak terpasang)._

### `/routing/route`

- File: `routing_route.json`
- Sumber field: live 18, completion 91

Command: `comment`(2), `disable`(1), `edit`(2), `enable`(1), `find`(1), `get`(4), `print`(21), `reset`(3), `set`(3), `unset`(2)

```go
// RoutingRouteRow adalah kolom hasil `/routing/route/print`.
type RoutingRouteRow struct {
	ID                           string `ros:".id,key"`
	Active                       string `ros:"active"`
	Afi                          string `ros:"afi"`
	BelongsTo                    string `ros:"belongs-to"`
	Connect                      string `ros:"connect"`
	Contribution                 string `ros:"contribution"`
	DebugFwpPtr                  string `ros:"debug.fwp-ptr"`
	DHCP                         string `ros:"dhcp"`
	Distance                     string `ros:"distance"`
	DstAddress                   string `ros:"dst-address"`
	Gateway                      string `ros:"gateway"`
	ImmediateGw                  string `ros:"immediate-gw"`
	LocalAddress                 string `ros:"local-address"`
	NexthopID                    string `ros:"nexthop-id"`
	RoutingTable                 string `ros:"routing-table"`
	Scope                        string `ros:"scope"`
	TargetScope                  string `ros:"target-scope"`
	VRFInterface                 string `ros:"vrf-interface"`
	About                        string `ros:"about"`
	BGP                          string `ros:"bgp"`
	BGPMPLSVPN                   string `ros:"bgp-mpls-vpn"`
	BGPNet                       string `ros:"bgp-net"`
	BGPAggregator                string `ros:"bgp.aggregator"`
	BGPAsPath                    string `ros:"bgp.as-path"`
	BGPAtomicAggregate           string `ros:"bgp.atomic-aggregate"`
	BGPClusterList               string `ros:"bgp.cluster-list"`
	BGPCommunities               string `ros:"bgp.communities"`
	BGPEsi                       string `ros:"bgp.esi"`
	BGPExtCommunities            string `ros:"bgp.ext-communities"`
	BGPIgpMetric                 string `ros:"bgp.igp-metric"`
	BGPLargeCommunities          string `ros:"bgp.large-communities"`
	BGPLocalPref                 string `ros:"bgp.local-pref"`
	BGPMed                       string `ros:"bgp.med"`
	BGPNexthop                   string `ros:"bgp.nexthop"`
	BGPOrigin                    string `ros:"bgp.origin"`
	BGPOriginatorID              string `ros:"bgp.originator-id"`
	BGPOtc                       string `ros:"bgp.otc"`
	BGPOutNexthop                string `ros:"bgp.out-nexthop"`
	BGPPmsi                      string `ros:"bgp.pmsi"`
	BGPSession                   string `ros:"bgp.session"`
	BGPUnknown                   string `ros:"bgp.unknown"`
	BGPWeight                    string `ros:"bgp.weight"`
	Blackhole                    string `ros:"blackhole"`
	CheckGateway                 string `ros:"check-gateway"`
	// string value
	Comment                      string `ros:"comment"`
	DebugInputPtr                string `ros:"debug.input.ptr"`
	DebugInputRaw                string `ros:"debug.input.raw"`
	DebugInputRefs               string `ros:"debug.input.refs"`
	DebugInputSize               string `ros:"debug.input.size"`
	DebugMergePtr                string `ros:"debug.merge.ptr"`
	DebugMergeRaw                string `ros:"debug.merge.raw"`
	DebugMergeRefs               string `ros:"debug.merge.refs"`
	DebugMergeSize               string `ros:"debug.merge.size"`
	DebugRoutePtr                string `ros:"debug.route-ptr"`
	Disabled                     string `ros:"disabled"`
	Ecmp                         string `ros:"ecmp"`
	Evpn                         string `ros:"evpn"`
	Filtered                     string `ros:"filtered"`
	Hoplimit                     string `ros:"hoplimit"`
	HWOffloaded                  string `ros:"hw-offloaded"`
	ISIS                         string `ros:"isis"`
	ISISMetric                   string `ros:"isis.metric"`
	ISISType                     string `ros:"isis.type"`
	Label                        string `ros:"label"`
	LdpAddress                   string `ros:"ldp-address"`
	LdpMapping                   string `ros:"ldp-mapping"`
	LdpLabel                     string `ros:"ldp.label"`
	LdpPeerID                    string `ros:"ldp.peer-id"`
	Modem                        string `ros:"modem"`
	MPLSInLabel                  string `ros:"mpls.in-label"`
	MPLSLabels                   string `ros:"mpls.labels"`
	MPLSOutLabel                 string `ros:"mpls.out-label"`
	MTU                          string `ros:"mtu"`
	OSPF                         string `ros:"ospf"`
	OSPFMetric                   string `ros:"ospf.metric"`
	OSPFTag                      string `ros:"ospf.tag"`
	OSPFType                     string `ros:"ospf.type"`
	PrefSrc                      string `ros:"pref-src"`
	RIP                          string `ros:"rip"`
	RIPMetric                    string `ros:"rip.metric"`
	RIPTag                       string `ros:"rip.tag"`
	RouteCost                    string `ros:"route-cost"`
	RouteDistinguisher           string `ros:"route-distinguisher"`
	Rpki                         string `ros:"rpki"`
	Slaac                        string `ros:"slaac"`
	Static                       string `ros:"static"`
	TotalCost                    string `ros:"total-cost"`
	Unreachable                  string `ros:"unreachable"`
	VeBlockOffset                string `ros:"ve-block-offset"`
	VeBlockSize                  string `ros:"ve-block-size"`
	VeID                         string `ros:"ve-id"`
	VPN                          string `ros:"vpn"`
}
```

```go
// RoutingRouteSet adalah argumen `/routing/route/set`.
type RoutingRouteSet struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
}
```

### `/routing/rpki`

- File: `routing_rpki.json`
- Sumber field: completion 11

Command: 

```go
// RoutingRpkiRow adalah kolom hasil `/routing/rpki/print`.
type RoutingRpkiRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Address                      string `ros:"address"`
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	ExpireInterval               string `ros:"expire-interval"`
	Group                        string `ros:"group"`
	Port                         string `ros:"port"`
	Preference                   string `ros:"preference"`
	RefreshInterval              string `ros:"refresh-interval"`
	RetryInterval                string `ros:"retry-interval"`
	VRF                          string `ros:"vrf"`
}
```

### `/routing/rule`

- File: `routing_rule.json`
- Sumber field: completion 15

Command: `add`(14), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `move`(2), `print`(21), `remove`(1), `reset`(13), `set`(13), `unset`(2)

```go
// RoutingRuleRow adalah kolom hasil `/routing/rule/print`.
type RoutingRuleRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Action                       string `ros:"action"`
	Chain                        string `ros:"chain"`
	Comment                      string `ros:"comment"`
	Default                      string `ros:"default"`
	Disabled                     string `ros:"disabled"`
	DstAddress                   string `ros:"dst-address"`
	Inactive                     string `ros:"inactive"`
	Interface                    string `ros:"interface"`
	MinPrefix                    string `ros:"min-prefix"`
	Realm                        string `ros:"realm"`
	RoutingMark                  string `ros:"routing-mark"`
	SrcAddress                   string `ros:"src-address"`
	Table                        string `ros:"table"`
	VRF                          string `ros:"vrf"`
}
```

```go
// RoutingRuleAdd adalah argumen `/routing/rule/add`.
type RoutingRuleAdd struct {
	Action                       string `ros:"action,omitempty"`
	Chain                        string `ros:"chain,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	DstAddress                   string `ros:"dst-address,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	MinPrefix                    string `ros:"min-prefix,omitempty"`
	PlaceBefore                  string `ros:"place-before,omitempty"`
	Realm                        string `ros:"realm,omitempty"`
	RoutingMark                  string `ros:"routing-mark,omitempty"`
	SrcAddress                   string `ros:"src-address,omitempty"`
	Table                        string `ros:"table,omitempty"`
	VRF                          string `ros:"vrf,omitempty"`
}
```

```go
// RoutingRuleSet adalah argumen `/routing/rule/set`.
type RoutingRuleSet struct {
	Action                       string `ros:"action,omitempty"`
	Chain                        string `ros:"chain,omitempty"`
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	DstAddress                   string `ros:"dst-address,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	MinPrefix                    string `ros:"min-prefix,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	Realm                        string `ros:"realm,omitempty"`
	RoutingMark                  string `ros:"routing-mark,omitempty"`
	SrcAddress                   string `ros:"src-address,omitempty"`
	Table                        string `ros:"table,omitempty"`
	VRF                          string `ros:"vrf,omitempty"`
}
```

### `/routing/settings`

- File: `routing_settings.json`
- Sumber field: live 5

Command: `edit`(1), `export`(8), `get`(3), `print`(5), `reset`(7), `set`(7)

```go
// RoutingSettingsRow adalah kolom hasil `/routing/settings/print`.
type RoutingSettingsRow struct {
	ID                           string `ros:".id,key"`
	CheckGatewayPingCount        string `ros:"check-gateway-ping-count"`
	// 100ms..10m (time interval)
	CheckGatewayPingInterval     string `ros:"check-gateway-ping-interval"`
	// 1ms..1s (time interval)
	CheckGatewayPingTimeout      string `ros:"check-gateway-ping-timeout"`
	PolicyRules                  string `ros:"policy-rules"`
	SingleProcess                string `ros:"single-process"`
}
```

```go
// RoutingSettingsSet adalah argumen `/routing/settings/set`.
type RoutingSettingsSet struct {
	CheckGatewayPingCount        string `ros:"check-gateway-ping-count,omitempty"`
	// 100ms..10m (time interval)
	CheckGatewayPingInterval     string `ros:"check-gateway-ping-interval,omitempty"`
	// 1ms..1s (time interval)
	CheckGatewayPingTimeout      string `ros:"check-gateway-ping-timeout,omitempty"`
	ConnectedInChain             string `ros:"connected-in-chain,omitempty"`
	DynamicInChain               string `ros:"dynamic-in-chain,omitempty"`
	PolicyRules                  string `ros:"policy-rules,omitempty"`
	SingleProcess                string `ros:"single-process,omitempty"`
}
```

### `/routing/stats`

- File: `routing_stats.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command

Command: 

_Tidak ada struct: menu ini bukan tabel yang bisa di-print atau tidak tersedia di router uji (mis. paket tidak terpasang)._

### `/routing/table`

- File: `routing_table.json`
- Sumber field: live 5, completion 8

Command: 

```go
// RoutingTableRow adalah kolom hasil `/routing/table/print`.
type RoutingTableRow struct {
	ID                           string `ros:".id,key"`
	Dynamic                      string `ros:"dynamic"`
	Fib                          string `ros:"fib"`
	Invalid                      string `ros:"invalid"`
	Name                         string `ros:"name"`
	About                        string `ros:"about"`
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	Used                         string `ros:"used"`
}
```

### `/safe-mode`

- File: `safe-mode.json`
- Sumber field: live 2

Command: 

```go
// SafemodeRow adalah kolom hasil `/safe-mode/print`.
type SafemodeRow struct {
	Current                      string `ros:"current"`
	Enabled                      string `ros:"enabled"`
}
```

### `/snmp/community`

- File: `snmp_community.json`
- Sumber field: live 12, completion 13

Command: 

```go
// SNMPCommunityRow adalah kolom hasil `/snmp/community/print`.
type SNMPCommunityRow struct {
	ID                           string `ros:".id,key"`
	Addresses                    string `ros:"addresses"`
	AuthenticationPassword       string `ros:"authentication-password"`
	AuthenticationProtocol       string `ros:"authentication-protocol"`
	Default                      string `ros:"default"`
	Disabled                     string `ros:"disabled"`
	EncryptionPassword           string `ros:"encryption-password"`
	EncryptionProtocol           string `ros:"encryption-protocol"`
	Name                         string `ros:"name"`
	ReadAccess                   string `ros:"read-access"`
	Security                     string `ros:"security"`
	WriteAccess                  string `ros:"write-access"`
	About                        string `ros:"about"`
	Comment                      string `ros:"comment"`
}
```

### `/special-login`

- File: `special-login.json`
- Sumber field: completion 5

Command: 

```go
// SpecialloginRow adalah kolom hasil `/special-login/print`.
type SpecialloginRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Channel                      string `ros:"channel"`
	Disabled                     string `ros:"disabled"`
	Port                         string `ros:"port"`
	User                         string `ros:"user"`
}
```

### `/system`

- File: `system__cmds.json`
- Katalog: `system-power`
- Sumber field: tidak ada kolom terdeteksi

Command: `check-disk`(0), `check-installation`(9), `export`(8), `reboot`(0), `recursive-print`(0), `reset-configuration`(8), `serial-terminal`(3), `shell`(0), `shutdown`(0), `ssh`(7), `ssh-exec`(9), `sup-output`(10), `telnet`(3)

_Tidak ada struct: menu ini bukan tabel yang bisa di-print atau tidak tersedia di router uji (mis. paket tidak terpasang)._

### `/system/backup`

- File: `system_backup.json`
- Katalog: `system-backup`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command

Command: `load`(3), `save`(4)

_Tidak ada struct: menu ini bukan tabel yang bisa di-print atau tidak tersedia di router uji (mis. paket tidak terpasang)._

### `/system/clock`

- File: `system_clock.json`
- Katalog: `system-clock`
- Sumber field: live 6, katalog 4

Command: `edit`(1), `export`(8), `get`(3), `print`(5), `reset`(4), `set`(4)

```go
// SystemClockRow adalah kolom hasil `/system/clock/print`.
type SystemClockRow struct {
	ID                           string `ros:".id,key"`
	Date                         string `ros:"date"`
	DstActive                    string `ros:"dst-active"`
	GmtOffset                    string `ros:"gmt-offset"`
	Time                         string `ros:"time"`
	TimeZoneAutodetect           string `ros:"time-zone-autodetect"`
	TimeZoneName                 string `ros:"time-zone-name"`
}
```

```go
// SystemClockSet adalah argumen `/system/clock/set`.
type SystemClockSet struct {
	Date                         string `ros:"date,omitempty"`
	Time                         string `ros:"time,omitempty"`
	TimeZoneAutodetect           string `ros:"time-zone-autodetect,omitempty"`
	TimeZoneName                 string `ros:"time-zone-name,omitempty"`
}
```

### `/system/console`

- File: `system_console.json`
- Sumber field: live 8, completion 10

Command: 

```go
// SystemConsoleRow adalah kolom hasil `/system/console/print`.
type SystemConsoleRow struct {
	ID                           string `ros:".id,key"`
	Channel                      string `ros:"channel"`
	Default                      string `ros:"default"`
	Disabled                     string `ros:"disabled"`
	Free                         string `ros:"free"`
	Port                         string `ros:"port"`
	Term                         string `ros:"term"`
	Vcno                         string `ros:"vcno"`
	About                        string `ros:"about"`
	Used                         string `ros:"used"`
	Wedged                       string `ros:"wedged"`
}
```

### `/system/default-configuration`

- File: `system_default-configuration.json`
- Sumber field: live 3

Command: 

```go
// SystemDefaultconfigurationRow adalah kolom hasil `/system/default-configuration/print`.
type SystemDefaultconfigurationRow struct {
	CAPsModeScript               string `ros:"caps-mode-script"`
	CustomScript                 string `ros:"custom-script"`
	Script                       string `ros:"script"`
}
```

### `/system/device-mode`

- File: `system_device-mode.json`
- Sumber field: live 24

Command: 

```go
// SystemDevicemodeRow adalah kolom hasil `/system/device-mode/print`.
type SystemDevicemodeRow struct {
	AllowedVersions              string `ros:"allowed-versions"`
	AttemptCount                 string `ros:"attempt-count"`
	BandwidthTest                string `ros:"bandwidth-test"`
	Container                    string `ros:"container"`
	Email                        string `ros:"email"`
	Fetch                        string `ros:"fetch"`
	Flagged                      string `ros:"flagged"`
	FlaggingEnabled              string `ros:"flagging-enabled"`
	Hotspot                      string `ros:"hotspot"`
	InstallAnyVersion            string `ros:"install-any-version"`
	IPSec                        string `ros:"ipsec"`
	L2TP                         string `ros:"l2tp"`
	Mode                         string `ros:"mode"`
	Partitions                   string `ros:"partitions"`
	PPTP                         string `ros:"pptp"`
	Proxy                        string `ros:"proxy"`
	Romon                        string `ros:"romon"`
	Routerboard                  string `ros:"routerboard"`
	Scheduler                    string `ros:"scheduler"`
	Smb                          string `ros:"smb"`
	Sniffer                      string `ros:"sniffer"`
	Socks                        string `ros:"socks"`
	TrafficGen                   string `ros:"traffic-gen"`
	Zerotier                     string `ros:"zerotier"`
}
```

### `/system/health`

- File: `system_health.json`
- Katalog: `system-health`
- Sumber field: live 2, katalog 18

Command: `edit`(1), `export`(8), `get`(3), `print`(6), `reset`(1), `set`(1)

```go
// SystemHealthRow adalah kolom hasil `/system/health/print`.
type SystemHealthRow struct {
	ID                           string `ros:".id,key"`
	State                        string `ros:"state"`
	StateAfterReboot             string `ros:"state-after-reboot"`
	Name                         string `ros:"name"`
	Value                        string `ros:"value"`
	Type                         string `ros:"type"`
	Voltage                      string `ros:"voltage"`
	Temperature                  string `ros:"temperature"`
	CPUTemperature               string `ros:"cpu-temperature"`
	PowerConsumption             string `ros:"power-consumption"`
	Fan1Speed                    string `ros:"fan1-speed"`
	Fan2Speed                    string `ros:"fan2-speed"`
}
```

```go
// SystemHealthSet adalah argumen `/system/health/set`.
type SystemHealthSet struct {
	StateAfterReboot             string `ros:"state-after-reboot,omitempty"`
}
```

### `/system/history`

- File: `system_history.json`
- Sumber field: completion 11

Command: `find`(1), `get`(4), `print`(20)

```go
// SystemHistoryRow adalah kolom hasil `/system/history/print`.
type SystemHistoryRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Action                       string `ros:"action"`
	By                           string `ros:"by"`
	FloatingUndo                 string `ros:"floating-undo"`
	Policy                       string `ros:"policy"`
	Redo                         string `ros:"redo"`
	Redoable                     string `ros:"redoable"`
	Time                         string `ros:"time"`
	Trace                        string `ros:"trace"`
	Undo                         string `ros:"undo"`
	Undoable                     string `ros:"undoable"`
}
```

### `/system/identity`

- File: `system_identity.json`
- Katalog: `system-identity`
- Sumber field: live 1, katalog 1

Command: `edit`(1), `export`(8), `get`(3), `print`(6), `reset`(1), `set`(1)

```go
// SystemIdentityRow adalah kolom hasil `/system/identity/print`.
type SystemIdentityRow struct {
	ID                           string `ros:".id,key"`
	// string value
	Name                         string `ros:"name"`
}
```

```go
// SystemIdentitySet adalah argumen `/system/identity/set`.
type SystemIdentitySet struct {
	// string value
	Name                         string `ros:"name,omitempty"`
}
```

### `/system/keymat-provider`

- File: `system_keymat-provider.json`
- Sumber field: completion 14

Command: 

```go
// SystemKeymatproviderRow adalah kolom hasil `/system/keymat-provider/print`.
type SystemKeymatproviderRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Disabled                     string `ros:"disabled"`
	HkdfHash                     string `ros:"hkdf-hash"`
	HkdfSecret                   string `ros:"hkdf-secret"`
	KeySize                      string `ros:"key-size"`
	Name                         string `ros:"name"`
	QkdAddress                   string `ros:"qkd-address"`
	QkdCacheSize                 string `ros:"qkd-cache-size"`
	QkdCacheState                string `ros:"qkd-cache-state"`
	QkdCertificate               string `ros:"qkd-certificate"`
	QkdKmeID                     string `ros:"qkd-kme-id"`
	QkdPeerSaeID                 string `ros:"qkd-peer-sae-id"`
	QkdTotalKeysReceived         string `ros:"qkd-total-keys-received"`
	Type                         string `ros:"type"`
}
```

### `/system/leds`

- File: `system_leds.json`
- Sumber field: completion 8

Command: 

```go
// SystemLEDsRow adalah kolom hasil `/system/leds/print`.
type SystemLEDsRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Color                        string `ros:"color"`
	Default                      string `ros:"default"`
	Disabled                     string `ros:"disabled"`
	Interface                    string `ros:"interface"`
	LEDs                         string `ros:"leds"`
	ModemSignalThreshold         string `ros:"modem-signal-threshold"`
	Type                         string `ros:"type"`
}
```

### `/system/license`

- File: `system_license.json`
- Sumber field: live 2

Command: `export`(8), `generate-new-id`(0), `get`(3), `print`(6), `renew`(12)

```go
// SystemLicenseRow adalah kolom hasil `/system/license/print`.
type SystemLicenseRow struct {
	Level                        string `ros:"level"`
	SystemID                     string `ros:"system-id"`
}
```

### `/system/logging`

- File: `system_logging.json`
- Sumber field: live 9, completion 10

Command: `add`(7), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(7), `set`(7)

```go
// SystemLoggingRow adalah kolom hasil `/system/logging/print`.
type SystemLoggingRow struct {
	ID                           string `ros:".id,key"`
	Action                       string `ros:"action"`
	Default                      string `ros:"default"`
	Disabled                     string `ros:"disabled"`
	Invalid                      string `ros:"invalid"`
	Managed                      string `ros:"managed"`
	// string value
	Prefix                       string `ros:"prefix"`
	// string value
	Regex                        string `ros:"regex"`
	Topics                       string `ros:"topics"`
	About                        string `ros:"about"`
	// string value
	Comment                      string `ros:"comment"`
}
```

```go
// SystemLoggingAdd adalah argumen `/system/logging/add`.
type SystemLoggingAdd struct {
	Action                       string `ros:"action,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	// string value
	Prefix                       string `ros:"prefix,omitempty"`
	// string value
	Regex                        string `ros:"regex,omitempty"`
	Topics                       string `ros:"topics,omitempty"`
}
```

```go
// SystemLoggingSet adalah argumen `/system/logging/set`.
type SystemLoggingSet struct {
	Action                       string `ros:"action,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	// string value
	Prefix                       string `ros:"prefix,omitempty"`
	// string value
	Regex                        string `ros:"regex,omitempty"`
	Topics                       string `ros:"topics,omitempty"`
}
```

### `/system/note`

- File: `system_note.json`
- Sumber field: live 3

Command: 

```go
// SystemNoteRow adalah kolom hasil `/system/note/print`.
type SystemNoteRow struct {
	Note                         string `ros:"note"`
	ShowAtCliLogin               string `ros:"show-at-cli-login"`
	ShowAtLogin                  string `ros:"show-at-login"`
}
```

### `/system/ntp`

- File: `system_ntp.json`
- Sumber field: tidak ada kolom terdeteksi
- Catatan print live: from RouterOS device: no such command

Command: `export`(8), `monitor-peers`(5)

_Tidak ada struct: menu ini bukan tabel yang bisa di-print atau tidak tersedia di router uji (mis. paket tidak terpasang)._

### `/system/package`

- File: `system_package.json`
- Katalog: `system-package`
- Sumber field: live 8, completion 9, katalog 5

Command: `apply-changes`(1), `disable`(1), `downgrade`(0), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(21), `uninstall`(1), `unschedule`(1)

```go
// SystemPackageRow adalah kolom hasil `/system/package/print`.
type SystemPackageRow struct {
	ID                           string `ros:".id,key"`
	Available                    string `ros:"available"`
	BuildTime                    string `ros:"build-time"`
	Disabled                     string `ros:"disabled"`
	Name                         string `ros:"name"`
	Scheduled                    string `ros:"scheduled"`
	Size                         string `ros:"size"`
	Version                      string `ros:"version"`
	About                        string `ros:"about"`
	Bundle                       string `ros:"bundle"`
}
```

### `/system/package/update`

- File: `(dari katalog Go)`
- Katalog: `system-package-update`
- Sumber field: katalog 4

Command: `set`(1)

```go
// SystemPackageUpdateRow adalah kolom hasil `/system/package/update/print`.
type SystemPackageUpdateRow struct {
	ID                           string `ros:".id,key"`
	Channel                      string `ros:"channel"`
	InstalledVersion             string `ros:"installed-version"`
	LatestVersion                string `ros:"latest-version"`
	Status                       string `ros:"status"`
}
```

```go
// SystemPackageUpdateSet adalah argumen `/system/package/update/set`.
type SystemPackageUpdateSet struct {
	Channel                      string `ros:"channel,omitempty"`
}
```

### `/system/regulatory`

- File: `system_regulatory.json`
- Sumber field: live 1

Command: 

```go
// SystemRegulatoryRow adalah kolom hasil `/system/regulatory/print`.
type SystemRegulatoryRow struct {
	Fcc                          string `ros:"fcc"`
}
```

### `/system/resource`

- File: `system_resource.json`
- Katalog: `system-resource`
- Sumber field: live 17, katalog 20

Command: `export`(8), `get`(3), `monitor`(10), `print`(6)

```go
// SystemResourceRow adalah kolom hasil `/system/resource/print`.
type SystemResourceRow struct {
	ArchitectureName             string `ros:"architecture-name"`
	BoardName                    string `ros:"board-name"`
	BuildTime                    string `ros:"build-time"`
	CPU                          string `ros:"cpu"`
	CPUCount                     string `ros:"cpu-count"`
	CPUFrequency                 string `ros:"cpu-frequency"`
	CPULoad                      string `ros:"cpu-load"`
	FreeHDDSpace                 string `ros:"free-hdd-space"`
	FreeMemory                   string `ros:"free-memory"`
	MinimumVersion               string `ros:"minimum-version"`
	Platform                     string `ros:"platform"`
	TotalHDDSpace                string `ros:"total-hdd-space"`
	TotalMemory                  string `ros:"total-memory"`
	Uptime                       string `ros:"uptime"`
	Version                      string `ros:"version"`
	WriteSectSinceReboot         string `ros:"write-sect-since-reboot"`
	WriteSectTotal               string `ros:"write-sect-total"`
}
```

### `/system/resource/cpu`

- File: `(dari katalog Go)`
- Katalog: `system-resource-cpu`
- Sumber field: katalog 8

Command: 

```go
// SystemResourceCPURow adalah kolom hasil `/system/resource/cpu/print`.
type SystemResourceCPURow struct {
	ID                           string `ros:".id,key"`
	CPU                          string `ros:"cpu"`
	Load                         string `ros:"load"`
	Irq                          string `ros:"irq"`
	Disk                         string `ros:"disk"`
}
```

### `/system/routerboard`

- File: `(dari katalog Go)`
- Katalog: `system-routerboard`
- Sumber field: katalog 7

Command: 

```go
// SystemRouterboardRow adalah kolom hasil `/system/routerboard/print`.
type SystemRouterboardRow struct {
	Routerboard                  string `ros:"routerboard"`
	Model                        string `ros:"model"`
	SerialNumber                 string `ros:"serial-number"`
	FirmwareType                 string `ros:"firmware-type"`
	FactoryFirmware              string `ros:"factory-firmware"`
	CurrentFirmware              string `ros:"current-firmware"`
	UpgradeFirmware              string `ros:"upgrade-firmware"`
}
```

### `/system/rtrace`

- File: `system_rtrace.json`
- Sumber field: live 1

Command: 

```go
// SystemRtraceRow adalah kolom hasil `/system/rtrace/print`.
type SystemRtraceRow struct {
	State                        string `ros:"state"`
}
```

### `/system/scheduler`

- File: `system_scheduler.json`
- Katalog: `system-scheduler`
- Sumber field: live 12, completion 13, katalog 10

Command: `add`(10), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(10), `set`(10), `unset`(2)

```go
// SystemSchedulerRow adalah kolom hasil `/system/scheduler/print`.
type SystemSchedulerRow struct {
	ID                           string `ros:".id,key"`
	// string value
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	// time interval
	Interval                     string `ros:"interval"`
	// string value
	Name                         string `ros:"name"`
	NextRun                      string `ros:"next-run"`
	// string value
	OnEvent                      string `ros:"on-event"`
	Owner                        string `ros:"owner"`
	Policy                       string `ros:"policy"`
	RunCount                     string `ros:"run-count"`
	StartDate                    string `ros:"start-date"`
	StartTime                    string `ros:"start-time"`
	About                        string `ros:"about"`
	Days                         string `ros:"days"`
}
```

```go
// SystemSchedulerAdd adalah argumen `/system/scheduler/add`.
type SystemSchedulerAdd struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Days                         string `ros:"days,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	// time interval
	Interval                     string `ros:"interval,omitempty"`
	// string value
	Name                         string `ros:"name,required,omitempty"`
	// string value
	OnEvent                      string `ros:"on-event,required,omitempty"`
	Policy                       string `ros:"policy,omitempty"`
	StartDate                    string `ros:"start-date,omitempty"`
	StartTime                    string `ros:"start-time,omitempty"`
}
```

```go
// SystemSchedulerSet adalah argumen `/system/scheduler/set`.
type SystemSchedulerSet struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Days                         string `ros:"days,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	// time interval
	Interval                     string `ros:"interval,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	// string value
	OnEvent                      string `ros:"on-event,omitempty"`
	Policy                       string `ros:"policy,omitempty"`
	StartDate                    string `ros:"start-date,omitempty"`
	StartTime                    string `ros:"start-time,omitempty"`
}
```

### `/system/script`

- File: `system_script.json`
- Katalog: `system-script`
- Sumber field: completion 10, katalog 14

Command: `add`(7), `comment`(2), `edit`(2), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(7), `run`(2), `set`(7), `unset`(2)

```go
// SystemScriptRow adalah kolom hasil `/system/script/print`.
type SystemScriptRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	// string value
	Comment                      string `ros:"comment"`
	DontRequirePermissions       string `ros:"dont-require-permissions"`
	Invalid                      string `ros:"invalid"`
	LastStarted                  string `ros:"last-started"`
	// string value
	Name                         string `ros:"name"`
	// string value
	Owner                        string `ros:"owner"`
	Policy                       string `ros:"policy"`
	RunCount                     string `ros:"run-count"`
	// string value
	Source                       string `ros:"source"`
}
```

```go
// SystemScriptAdd adalah argumen `/system/script/add`.
type SystemScriptAdd struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	DontRequirePermissions       string `ros:"dont-require-permissions,omitempty"`
	// string value
	Name                         string `ros:"name,required,omitempty"`
	// string value
	Owner                        string `ros:"owner,omitempty"`
	Policy                       string `ros:"policy,omitempty"`
	// string value
	Source                       string `ros:"source,required,omitempty"`
}
```

```go
// SystemScriptSet adalah argumen `/system/script/set`.
type SystemScriptSet struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	DontRequirePermissions       string `ros:"dont-require-permissions,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	// string value
	Owner                        string `ros:"owner,omitempty"`
	Policy                       string `ros:"policy,omitempty"`
	// string value
	Source                       string `ros:"source,omitempty"`
}
```

### `/system/ups`

- File: `system_ups.json`
- Sumber field: completion 18

Command: 

```go
// SystemUpsRow adalah kolom hasil `/system/ups/print`.
type SystemUpsRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	AlarmSetting                 string `ros:"alarm-setting"`
	CheckCapabilities            string `ros:"check-capabilities"`
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	Invalid                      string `ros:"invalid"`
	Load                         string `ros:"load"`
	ManufactureDate              string `ros:"manufacture-date"`
	MinRuntime                   string `ros:"min-runtime"`
	Model                        string `ros:"model"`
	Name                         string `ros:"name"`
	NominalBatteryVoltage        string `ros:"nominal-battery-voltage"`
	OfflineAfter                 string `ros:"offline-after"`
	OfflineTime                  string `ros:"offline-time"`
	OnLine                       string `ros:"on-line"`
	Port                         string `ros:"port"`
	Serial                       string `ros:"serial"`
	Version                      string `ros:"version"`
}
```

### `/system/watchdog`

- File: `system_watchdog.json`
- Sumber field: live 6

Command: 

```go
// SystemWatchdogRow adalah kolom hasil `/system/watchdog/print`.
type SystemWatchdogRow struct {
	AutoSendSupout               string `ros:"auto-send-supout"`
	AutomaticSupout              string `ros:"automatic-supout"`
	PingStartAfterBoot           string `ros:"ping-start-after-boot"`
	PingTimeout                  string `ros:"ping-timeout"`
	WatchAddress                 string `ros:"watch-address"`
	WatchdogTimer                string `ros:"watchdog-timer"`
}
```

### `/task`

- File: `task.json`
- Sumber field: completion 12

Command: 

```go
// TaskRow adalah kolom hasil `/task/print`.
type TaskRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Append                       string `ros:"append"`
	Autosave                     string `ros:"autosave"`
	Current                      string `ros:"current"`
	FileName                     string `ros:"file-name"`
	Name                         string `ros:"name"`
	SaveInterval                 string `ros:"save-interval"`
	Source                       string `ros:"source"`
	TaskID                       string `ros:"task-id"`
	TermPid                      string `ros:"term-pid"`
	Terminated                   string `ros:"terminated"`
	User                         string `ros:"user"`
}
```

### `/tool`

- File: `tool__cmds.json`
- Sumber field: tidak ada kolom terdeteksi

Command: `bandwidth-test`(21), `dns-update`(7), `export`(8), `fetch`(31), `flood-ping`(14), `ip-scan`(7), `mac-scan`(6), `mac-telnet`(2), `ping`(15), `ping-speed`(14), `profile`(4), `recursive-print`(0), `snmp-get`(15), `snmp-walk`(14), `speed-test`(15), `torch`(16), `traceroute`(18), `wol`(3)

_Tidak ada struct: menu ini bukan tabel yang bisa di-print atau tidak tersedia di router uji (mis. paket tidak terpasang)._

### `/tool/bandwidth-server`

- File: `tool_bandwidth-server.json`
- Sumber field: live 6

Command: `edit`(1), `export`(8), `get`(3), `print`(5), `reset`(6), `set`(6)

```go
// ToolBandwidthserverRow adalah kolom hasil `/tool/bandwidth-server/print`.
type ToolBandwidthserverRow struct {
	ID                           string `ros:".id,key"`
	AllocateUdpPortsFrom         string `ros:"allocate-udp-ports-from"`
	// A.B.C.D/M (IP prefix)
	AllowedAddresses4            string `ros:"allowed-addresses4"`
	// IPv6/0..128 (IPv6 prefix)
	AllowedAddresses6            string `ros:"allowed-addresses6"`
	Authenticate                 string `ros:"authenticate"`
	Enabled                      string `ros:"enabled"`
	MaxSessions                  string `ros:"max-sessions"`
}
```

```go
// ToolBandwidthserverSet adalah argumen `/tool/bandwidth-server/set`.
type ToolBandwidthserverSet struct {
	AllocateUdpPortsFrom         string `ros:"allocate-udp-ports-from,omitempty"`
	// A.B.C.D/M (IP prefix)
	AllowedAddresses4            string `ros:"allowed-addresses4,omitempty"`
	// IPv6/0..128 (IPv6 prefix)
	AllowedAddresses6            string `ros:"allowed-addresses6,omitempty"`
	Authenticate                 string `ros:"authenticate,omitempty"`
	Enabled                      string `ros:"enabled,omitempty"`
	MaxSessions                  string `ros:"max-sessions,omitempty"`
}
```

### `/tool/bandwidth-test`

- File: `(dari katalog Go)`
- Katalog: `tool-bandwidth-test`
- Sumber field: katalog 13

Command: 

```go
// ToolBandwidthtestRow adalah kolom hasil `/tool/bandwidth-test/print`.
type ToolBandwidthtestRow struct {
	Status                       string `ros:"status"`
	Duration                     string `ros:"duration"`
	TXCurrent                    string `ros:"tx-current"`
	TXX10SecondAverage           string `ros:"tx-10-second-average"`
	TXTotalAverage               string `ros:"tx-total-average"`
	RXCurrent                    string `ros:"rx-current"`
	RXX10SecondAverage           string `ros:"rx-10-second-average"`
	RXTotalAverage               string `ros:"rx-total-average"`
	LostPackets                  string `ros:"lost-packets"`
	Direction                    string `ros:"direction"`
	ConnectionCount              string `ros:"connection-count"`
	LocalCPULoad                 string `ros:"local-cpu-load"`
	RemoteCPULoad                string `ros:"remote-cpu-load"`
}
```

### `/tool/e-mail`

- File: `tool_e-mail.json`
- Sumber field: live 8

Command: 

```go
// ToolEmailRow adalah kolom hasil `/tool/e-mail/print`.
type ToolEmailRow struct {
	CertificateVerification      string `ros:"certificate-verification"`
	From                         string `ros:"from"`
	Password                     string `ros:"password"`
	Port                         string `ros:"port"`
	Server                       string `ros:"server"`
	TLS                          string `ros:"tls"`
	User                         string `ros:"user"`
	VRF                          string `ros:"vrf"`
}
```

### `/tool/fetch`

- File: `(dari katalog Go)`
- Katalog: `tool-fetch`
- Sumber field: katalog 4

Command: 

```go
// ToolFetchRow adalah kolom hasil `/tool/fetch/print`.
type ToolFetchRow struct {
	Status                       string `ros:"status"`
	Downloaded                   string `ros:"downloaded"`
	Total                        string `ros:"total"`
	Duration                     string `ros:"duration"`
}
```

### `/tool/graphing`

- File: `tool_graphing.json`
- Sumber field: live 2

Command: 

```go
// ToolGraphingRow adalah kolom hasil `/tool/graphing/print`.
type ToolGraphingRow struct {
	PageRefresh                  string `ros:"page-refresh"`
	StoreEvery                   string `ros:"store-every"`
}
```

### `/tool/mac-server`

- File: `tool_mac-server.json`
- Sumber field: live 1

Command: 

```go
// ToolMacserverRow adalah kolom hasil `/tool/mac-server/print`.
type ToolMacserverRow struct {
	AllowedInterfaceList         string `ros:"allowed-interface-list"`
}
```

### `/tool/netwatch`

- File: `tool_netwatch.json`
- Katalog: `tool-netwatch`
- Sumber field: completion 57, katalog 11

Command: 

```go
// ToolNetwatchRow adalah kolom hasil `/tool/netwatch/print`.
type ToolNetwatchRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	AcceptIcmpTimeExceeded       string `ros:"accept-icmp-time-exceeded"`
	Certificate                  string `ros:"certificate"`
	CheckCertificate             string `ros:"check-certificate"`
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	DNSServer                    string `ros:"dns-server"`
	DoneTests                    string `ros:"done-tests"`
	DownScript                   string `ros:"down-script"`
	EarlyFailureDetection        string `ros:"early-failure-detection"`
	EarlySuccessDetection        string `ros:"early-success-detection"`
	FailedTests                  string `ros:"failed-tests"`
	Host                         string `ros:"host"`
	HTTPCodes                    string `ros:"http-codes"`
	HTTPRespTime                 string `ros:"http-resp-time"`
	HTTPStatusCode               string `ros:"http-status-code"`
	IgnoreInitialDown            string `ros:"ignore-initial-down"`
	IgnoreInitialUp              string `ros:"ignore-initial-up"`
	Interval                     string `ros:"interval"`
	IP                           string `ros:"ip"`
	Ip6                          string `ros:"ip6"`
	LossCount                    string `ros:"loss-count"`
	LossPercent                  string `ros:"loss-percent"`
	MailServers                  string `ros:"mail-servers"`
	Name                         string `ros:"name"`
	NameServers                  string `ros:"name-servers"`
	PacketCount                  string `ros:"packet-count"`
	PacketInterval               string `ros:"packet-interval"`
	PacketSize                   string `ros:"packet-size"`
	Port                         string `ros:"port"`
	RecordType                   string `ros:"record-type"`
	ResponseCount                string `ros:"response-count"`
	RTTAvg                       string `ros:"rtt-avg"`
	RTTJitter                    string `ros:"rtt-jitter"`
	RTTMax                       string `ros:"rtt-max"`
	RTTMin                       string `ros:"rtt-min"`
	RTTStdev                     string `ros:"rtt-stdev"`
	SentCount                    string `ros:"sent-count"`
	Since                        string `ros:"since"`
	SrcAddress                   string `ros:"src-address"`
	StartDelay                   string `ros:"start-delay"`
	StartupDelay                 string `ros:"startup-delay"`
	Status                       string `ros:"status"`
	TcpConnectTime               string `ros:"tcp-connect-time"`
	TestScript                   string `ros:"test-script"`
	ThrAvg                       string `ros:"thr-avg"`
	ThrHTTPTime                  string `ros:"thr-http-time"`
	ThrJitter                    string `ros:"thr-jitter"`
	ThrLossCount                 string `ros:"thr-loss-count"`
	ThrLossPercent               string `ros:"thr-loss-percent"`
	ThrMax                       string `ros:"thr-max"`
	ThrStdev                     string `ros:"thr-stdev"`
	ThrTcpConnTime               string `ros:"thr-tcp-conn-time"`
	Timeout                      string `ros:"timeout"`
	TTL                          string `ros:"ttl"`
	Type                         string `ros:"type"`
	UpScript                     string `ros:"up-script"`
}
```

### `/tool/profile`

- File: `(dari katalog Go)`
- Katalog: `tool-profile`
- Sumber field: katalog 3

Command: 

```go
// ToolProfileRow adalah kolom hasil `/tool/profile/print`.
type ToolProfileRow struct {
	Name                         string `ros:"name"`
	CPU                          string `ros:"cpu"`
	Usage                        string `ros:"usage"`
}
```

### `/tool/romon`

- File: `tool_romon.json`
- Sumber field: live 3

Command: `discover`(5), `edit`(1), `export`(8), `get`(3), `ping`(7), `print`(6), `reset`(3), `set`(3), `ssh`(4)

```go
// ToolRomonRow adalah kolom hasil `/tool/romon/print`.
type ToolRomonRow struct {
	ID                           string `ros:".id,key"`
	Enabled                      string `ros:"enabled"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	ID2                          string `ros:"id"`
	// string value
	Secrets                      string `ros:"secrets"`
}
```

```go
// ToolRomonSet adalah argumen `/tool/romon/set`.
type ToolRomonSet struct {
	Enabled                      string `ros:"enabled,omitempty"`
	// AB[:|-|.]CD[:|-|.]EF[:|-|.]GH[:|-|.]IJ[:|-|.]KL (MAC address)
	ID                           string `ros:"id,omitempty"`
	// string value
	Secrets                      string `ros:"secrets,omitempty"`
}
```

### `/tool/sms`

- File: `tool_sms.json`
- Sumber field: live 11

Command: 

```go
// ToolSMSRow adalah kolom hasil `/tool/sms/print`.
type ToolSMSRow struct {
	AllowedNumber                string `ros:"allowed-number"`
	Channel                      string `ros:"channel"`
	LastUssd                     string `ros:"last-ussd"`
	Polling                      string `ros:"polling"`
	Port                         string `ros:"port"`
	ReceiveEnabled               string `ros:"receive-enabled"`
	RemoveSentSMSAfterSend       string `ros:"remove-sent-sms-after-send"`
	Secret                       string `ros:"secret"`
	SimPin                       string `ros:"sim-pin"`
	SMSStorage                   string `ros:"sms-storage"`
	Status                       string `ros:"status"`
}
```

### `/tool/sniffer`

- File: `tool_sniffer.json`
- Sumber field: live 32

Command: 

```go
// ToolSnifferRow adalah kolom hasil `/tool/sniffer/print`.
type ToolSnifferRow struct {
	FileLimit                    string `ros:"file-limit"`
	FileName                     string `ros:"file-name"`
	FilterCPU                    string `ros:"filter-cpu"`
	FilterDirection              string `ros:"filter-direction"`
	FilterDstIPAddress           string `ros:"filter-dst-ip-address"`
	FilterDstIPv6Address         string `ros:"filter-dst-ipv6-address"`
	FilterDstMACAddress          string `ros:"filter-dst-mac-address"`
	FilterDstPort                string `ros:"filter-dst-port"`
	FilterInterface              string `ros:"filter-interface"`
	FilterIPAddress              string `ros:"filter-ip-address"`
	FilterIPProtocol             string `ros:"filter-ip-protocol"`
	FilterIPv6Address            string `ros:"filter-ipv6-address"`
	FilterMACAddress             string `ros:"filter-mac-address"`
	FilterMACProtocol            string `ros:"filter-mac-protocol"`
	FilterOperatorBetweenEntries string `ros:"filter-operator-between-entries"`
	FilterPort                   string `ros:"filter-port"`
	FilterSize                   string `ros:"filter-size"`
	FilterSrcIPAddress           string `ros:"filter-src-ip-address"`
	FilterSrcIPv6Address         string `ros:"filter-src-ipv6-address"`
	FilterSrcMACAddress          string `ros:"filter-src-mac-address"`
	FilterSrcPort                string `ros:"filter-src-port"`
	FilterStream                 string `ros:"filter-stream"`
	FilterVLAN                   string `ros:"filter-vlan"`
	MaxPacketSize                string `ros:"max-packet-size"`
	MemoryLimit                  string `ros:"memory-limit"`
	MemoryScroll                 string `ros:"memory-scroll"`
	OnlyHeaders                  string `ros:"only-headers"`
	QuickRows                    string `ros:"quick-rows"`
	QuickShowFrame               string `ros:"quick-show-frame"`
	Running                      string `ros:"running"`
	StreamingEnabled             string `ros:"streaming-enabled"`
	StreamingServer              string `ros:"streaming-server"`
}
```

### `/tool/torch`

- File: `(dari katalog Go)`
- Katalog: `tool-torch`
- Sumber field: katalog 11

Command: 

```go
// ToolTorchRow adalah kolom hasil `/tool/torch/print`.
type ToolTorchRow struct {
	SrcAddress                   string `ros:"src-address"`
	DstAddress                   string `ros:"dst-address"`
	MACProtocol                  string `ros:"mac-protocol"`
	IPProtocol                   string `ros:"ip-protocol"`
	SrcPort                      string `ros:"src-port"`
	DstPort                      string `ros:"dst-port"`
	VLANID                       string `ros:"vlan-id"`
	TX                           string `ros:"tx"`
	RX                           string `ros:"rx"`
	TXPackets                    string `ros:"tx-packets"`
	RXPackets                    string `ros:"rx-packets"`
}
```

### `/tool/traceroute`

- File: `(dari katalog Go)`
- Katalog: `tool-traceroute`
- Sumber field: katalog 9

Command: 

```go
// ToolTracerouteRow adalah kolom hasil `/tool/traceroute/print`.
type ToolTracerouteRow struct {
	Address                      string `ros:"address"`
	Loss                         string `ros:"loss"`
	Sent                         string `ros:"sent"`
	Last                         string `ros:"last"`
	Avg                          string `ros:"avg"`
	Best                         string `ros:"best"`
	Worst                        string `ros:"worst"`
	StdDev                       string `ros:"std-dev"`
	Status                       string `ros:"status"`
}
```

### `/tool/traffic-generator`

- File: `tool_traffic-generator.json`
- Sumber field: live 7

Command: 

```go
// ToolTrafficgeneratorRow adalah kolom hasil `/tool/traffic-generator/print`.
type ToolTrafficgeneratorRow struct {
	LatencyDistributionMax       string `ros:"latency-distribution-max"`
	LatencyDistributionMeasureInterval string `ros:"latency-distribution-measure-interval"`
	LatencyDistributionSamples   string `ros:"latency-distribution-samples"`
	MeasureOutOfOrder            string `ros:"measure-out-of-order"`
	Running                      string `ros:"running"`
	StatsSamplesToKeep           string `ros:"stats-samples-to-keep"`
	TestID                       string `ros:"test-id"`
}
```

### `/tool/traffic-monitor`

- File: `tool_traffic-monitor.json`
- Sumber field: completion 10

Command: `add`(9), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(9), `set`(9)

```go
// ToolTrafficmonitorRow adalah kolom hasil `/tool/traffic-monitor/print`.
type ToolTrafficmonitorRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	// string value
	Comment                      string `ros:"comment"`
	Disabled                     string `ros:"disabled"`
	Interface                    string `ros:"interface"`
	Invalid                      string `ros:"invalid"`
	// string value
	Name                         string `ros:"name"`
	// string value
	OnEvent                      string `ros:"on-event"`
	Threshold                    string `ros:"threshold"`
	Traffic                      string `ros:"traffic"`
	Trigger                      string `ros:"trigger"`
}
```

```go
// ToolTrafficmonitorAdd adalah argumen `/tool/traffic-monitor/add`.
type ToolTrafficmonitorAdd struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	// string value
	OnEvent                      string `ros:"on-event,omitempty"`
	Threshold                    string `ros:"threshold,omitempty"`
	Traffic                      string `ros:"traffic,omitempty"`
	Trigger                      string `ros:"trigger,omitempty"`
}
```

```go
// ToolTrafficmonitorSet adalah argumen `/tool/traffic-monitor/set`.
type ToolTrafficmonitorSet struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Interface                    string `ros:"interface,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	// string value
	OnEvent                      string `ros:"on-event,omitempty"`
	Threshold                    string `ros:"threshold,omitempty"`
	Traffic                      string `ros:"traffic,omitempty"`
	Trigger                      string `ros:"trigger,omitempty"`
}
```

### `/user`

- File: `user__cmds.json`
- Katalog: `user`
- Sumber field: katalog 6

Command: `add`(9), `comment`(2), `disable`(1), `edit`(2), `enable`(1), `expire-password`(1), `export`(8), `find`(1), `get`(4), `print`(22), `recursive-print`(0), `remove`(1), `reset`(9), `set`(9)

```go
// UserRow adalah kolom hasil `/user/print`.
type UserRow struct {
	ID                           string `ros:".id,key"`
	// string value
	Name                         string `ros:"name"`
	Group                        string `ros:"group"`
	Address                      string `ros:"address"`
	LastLoggedIn                 string `ros:"last-logged-in"`
	Disabled                     string `ros:"disabled"`
	// string value
	Comment                      string `ros:"comment"`
}
```

```go
// UserAdd adalah argumen `/user/add`.
type UserAdd struct {
	Address                      string `ros:"address,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Group                        string `ros:"group,required,omitempty"`
	InactivityPolicy             string `ros:"inactivity-policy,omitempty"`
	InactivityTimeout            string `ros:"inactivity-timeout,omitempty"`
	// string value
	Name                         string `ros:"name,required,omitempty"`
	Password                     string `ros:"password,required,omitempty"`
}
```

```go
// UserSet adalah argumen `/user/set`.
type UserSet struct {
	Address                      string `ros:"address,omitempty"`
	// string value
	Comment                      string `ros:"comment,omitempty"`
	Disabled                     string `ros:"disabled,omitempty"`
	Group                        string `ros:"group,omitempty"`
	InactivityPolicy             string `ros:"inactivity-policy,omitempty"`
	InactivityTimeout            string `ros:"inactivity-timeout,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	Password                     string `ros:"password,omitempty"`
}
```

### `/user/aaa`

- File: `user_aaa.json`
- Sumber field: live 5

Command: `edit`(1), `export`(8), `get`(3), `print`(5), `reset`(5), `set`(5)

```go
// UserAaaRow adalah kolom hasil `/user/aaa/print`.
type UserAaaRow struct {
	ID                           string `ros:".id,key"`
	Accounting                   string `ros:"accounting"`
	DefaultGroup                 string `ros:"default-group"`
	ExcludeGroups                string `ros:"exclude-groups"`
	// time interval
	InterimUpdate                string `ros:"interim-update"`
	UseRadius                    string `ros:"use-radius"`
}
```

```go
// UserAaaSet adalah argumen `/user/aaa/set`.
type UserAaaSet struct {
	Accounting                   string `ros:"accounting,omitempty"`
	DefaultGroup                 string `ros:"default-group,omitempty"`
	ExcludeGroups                string `ros:"exclude-groups,omitempty"`
	// time interval
	InterimUpdate                string `ros:"interim-update,omitempty"`
	UseRadius                    string `ros:"use-radius,omitempty"`
}
```

### `/user/active`

- File: `user_active.json`
- Katalog: `user-active`
- Sumber field: live 7, completion 8, katalog 6

Command: `add`(1), `find`(1), `get`(4), `print`(21), `remove`(1)

```go
// UserActiveRow adalah kolom hasil `/user/active/print`.
type UserActiveRow struct {
	ID                           string `ros:".id,key"`
	Address                      string `ros:"address"`
	Group                        string `ros:"group"`
	Name                         string `ros:"name"`
	Radius                       string `ros:"radius"`
	Via                          string `ros:"via"`
	When                         string `ros:"when"`
	About                        string `ros:"about"`
	ByRomon                      string `ros:"by-romon"`
}
```

```go
// UserActiveAdd adalah argumen `/user/active/add`.
type UserActiveAdd struct {
	CopyFrom                     string `ros:"copy-from,omitempty"`
}
```

### `/user/group`

- File: `user_group.json`
- Sumber field: live 4, completion 5

Command: `add`(5), `comment`(2), `edit`(2), `export`(8), `find`(1), `get`(4), `print`(21), `remove`(1), `reset`(5), `set`(5)

```go
// UserGroupRow adalah kolom hasil `/user/group/print`.
type UserGroupRow struct {
	ID                           string `ros:".id,key"`
	// string value
	Name                         string `ros:"name"`
	Policy                       string `ros:"policy"`
	Skin                         string `ros:"skin"`
	About                        string `ros:"about"`
	// string value
	Comment                      string `ros:"comment"`
}
```

```go
// UserGroupAdd adalah argumen `/user/group/add`.
type UserGroupAdd struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	CopyFrom                     string `ros:"copy-from,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Policy                       string `ros:"policy,omitempty"`
	Skin                         string `ros:"skin,omitempty"`
}
```

```go
// UserGroupSet adalah argumen `/user/group/set`.
type UserGroupSet struct {
	// string value
	Comment                      string `ros:"comment,omitempty"`
	// string value
	Name                         string `ros:"name,omitempty"`
	Numbers                      string `ros:"numbers,omitempty"`
	Policy                       string `ros:"policy,omitempty"`
	Skin                         string `ros:"skin,omitempty"`
}
```

### `/user/settings`

- File: `user_settings.json`
- Sumber field: live 2

Command: `edit`(1), `export`(8), `get`(3), `print`(5), `reset`(2), `set`(2)

```go
// UserSettingsRow adalah kolom hasil `/user/settings/print`.
type UserSettingsRow struct {
	ID                           string `ros:".id,key"`
	MinimumCategories            string `ros:"minimum-categories"`
	MinimumPasswordLength        string `ros:"minimum-password-length"`
}
```

```go
// UserSettingsSet adalah argumen `/user/settings/set`.
type UserSettingsSet struct {
	MinimumCategories            string `ros:"minimum-categories,omitempty"`
	MinimumPasswordLength        string `ros:"minimum-password-length,omitempty"`
}
```

### `/user/ssh-keys`

- File: `user_ssh-keys.json`
- Sumber field: completion 8

Command: 

```go
// UserSshkeysRow adalah kolom hasil `/user/ssh-keys/print`.
type UserSshkeysRow struct {
	ID                           string `ros:".id,key"`
	About                        string `ros:"about"`
	Bits                         string `ros:"bits"`
	Comment                      string `ros:"comment"`
	Fingerprint                  string `ros:"fingerprint"`
	Info                         string `ros:"info"`
	Key                          string `ros:"key"`
	KeyType                      string `ros:"key-type"`
	User                         string `ros:"user"`
}
```
