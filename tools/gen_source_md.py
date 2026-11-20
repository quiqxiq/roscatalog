#!/usr/bin/env python3
"""Buat mikrotik_inspect/SOURCE.md berisi struct Go lengkap untuk semua menu.

Sumber:
  1. mikrotik_inspect/*.json        -> command, argumen, deskripsi (/console/inspect)
  2. print_fields.json              -> kolom print dari router live + completion
  3. catalog_dump.json              -> proplist/required/optional dari katalog Go

Tooling maintainer saja. Tidak dipakai saat build library.
"""

from __future__ import annotations

import argparse
import json
import re
from collections import Counter
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
INSPECT = ROOT / "mikrotik_inspect"
OUT = INSPECT / "SOURCE.md"

GENERATED = {"_index.json", "print_fields.json", "catalog_dump.json", "SOURCE.md"}

ACRONYMS = {
    "ip": "IP", "ipv4": "IPv4", "ipv6": "IPv6", "dns": "DNS", "dhcp": "DHCP",
    "dhcpv6": "DHCPv6", "ppp": "PPP", "pppoe": "PPPoE", "pptp": "PPTP",
    "l2tp": "L2TP", "sstp": "SSTP", "ovpn": "OVPN", "openvpn": "OpenVPN",
    "ipsec": "IPSec", "mac": "MAC", "vlan": "VLAN", "vrf": "VRF",
    "mtu": "MTU", "l2mtu": "L2MTU", "arp": "ARP", "rarp": "RARP",
    "cpu": "CPU", "hdd": "HDD", "ram": "RAM", "gps": "GPS", "lte": "LTE",
    "id": "ID", "nat": "NAT", "mpls": "MPLS", "ospf": "OSPF", "bgp": "BGP",
    "rip": "RIP", "isis": "ISIS", "snmp": "SNMP", "ntp": "NTP", "ssh": "SSH",
    "tls": "TLS", "ssl": "SSL", "http": "HTTP", "https": "HTTPS",
    "wifi": "WiFi", "gre": "GRE", "ipip": "IPIP", "eoip": "EOIP",
    "hw": "HW", "sw": "SW", "led": "LED", "leds": "LEDs", "lan": "LAN",
    "wan": "WAN", "vpn": "VPN", "vxlan": "VXLAN", "vpls": "VPLS",
    "vrrp": "VRRP", "bfd": "BFD", "igmp": "IGMP", "pimsm": "PIMSM",
    "sms": "SMS", "usb": "USB", "pci": "PCI", "nvme": "NVMe", "sata": "SATA",
    "rpc": "RPC", "api": "API", "tls": "TLS", "ttl": "TTL", "tx": "TX",
    "rx": "RX", "rtt": "RTT", "snmp": "SNMP", "vlan": "VLAN", "wpa": "WPA",
    "wps": "WPS", "dhcp-relay": "DHCPRelay", "caps": "CAPs", "cap": "CAP",
}

KEYWORDS = {
    "break", "case", "chan", "const", "continue", "default", "defer", "else",
    "fallthrough", "for", "func", "go", "goto", "if", "import", "interface",
    "map", "package", "range", "return", "select", "struct", "switch",
    "type", "var",
}


def load(path: Path):
    with path.open(encoding="utf-8") as f:
        return json.load(f)


def split_words(text: str) -> list[str]:
    return [w for w in re.split(r"[-_.\s]+", text) if w]


def go_field(prop: str) -> str:
    if prop == ".id":
        return "ID"
    words = split_words(prop)
    parts = []
    for w in words:
        key = w.lower()
        if key in ACRONYMS:
            parts.append(ACRONYMS[key])
        elif w and w[0].isdigit():
            parts.append("X" + w)
        else:
            parts.append(w[:1].upper() + w[1:])
    name = "".join(parts) or "Field"
    name = re.sub(r"[^A-Za-z0-9_]", "", name)
    if name in KEYWORDS:
        name += "Value"
    if name[0].isdigit():
        name = "X" + name
    return name


def go_type(path: str, suffix: str, used: set[str]) -> str:
    segs = [s for s in path.split("/") if s]
    parts = []
    for s in segs:
        key = s.lower()
        if key in ACRONYMS:
            parts.append(ACRONYMS[key])
        else:
            parts.append(s[:1].upper() + s[1:])
    base = re.sub(r"[^A-Za-z0-9_]", "", "".join(parts)) or "Menu"
    name = base + suffix
    if name in used:
        i = 2
        while f"{name}{i}" in used:
            i += 1
        name = f"{name}{i}"
    used.add(name)
    return name


def cmds_of(data: dict) -> dict:
    return {
        k: v for k, v in data.items()
        if isinstance(v, dict) and v.get("_type") == "cmd"
    }


def args_of(cmd: dict) -> list[str]:
    return [
        k for k, v in cmd.items()
        if k != "_type" and isinstance(v, dict) and v.get("_type") == "arg"
    ]


def desc_of(cmd: dict, prop: str) -> str:
    v = cmd.get(prop)
    if isinstance(v, dict):
        return v.get("desc", "")
    return ""


def clean_desc(s: str) -> str:
    s = re.sub(r"\s+", " ", s or "").strip()
    return s.replace("`", "'").replace("*/", "* /")


def schema_of(path: str, data: dict, cat: dict | None, probe: dict) -> tuple[list[str], dict]:
    """Gabungkan kolom print dari live + completion + katalog."""
    live = set(probe.get("fields", []))
    comp = set(probe.get("completion", []))
    catprops: list[str] = []
    if cat:
        for v in cat.get("views", []):
            if v.get("name") in {"list", "live", "detail", "monitor", "run", "status"}:
                catprops.extend(v.get("proplist", []))

    cmds = cmds_of(data)
    is_table = (
        "remove" in cmds or "set" in cmds or "add" in cmds or ".id" in live
        or ((bool(comp) or cat is not None) and not (cat and cat.get("singleton")))
    )

    props: list[str] = []
    if is_table:
        props.append(".id")
    seen = set(props)
    for group in (sorted(live), sorted(comp), catprops):
        for p in group:
            if p in seen or not p:
                continue
            seen.add(p)
            props.append(p)

    # Deskripsi property: cari dari add/set/print.
    descs: dict[str, str] = {}
    for name in ("add", "set", "print", "get", "reset"):
        cmd = cmds.get(name)
        if not cmd:
            continue
        for p in props:
            if p in descs:
                continue
            d = desc_of(cmd, p)
            if d:
                descs[p] = clean_desc(d)
    return props, {"live": len(live), "comp": len(comp), "cat": len(catprops),
                   "table": is_table, "descs": descs}


def required_of(cat: dict | None) -> set[str]:
    if not cat:
        return set()
    for a in cat.get("actions", []):
        if a.get("name") in {"add", "set"}:
            return set(a.get("required", []))
    return set()


def emit_struct(lines: list[str], name: str, props: list[str], descs: dict,
                required: set[str], omitempty: bool, note: str) -> None:
    lines.append("```go")
    lines.append(f"// {name} {note}".rstrip())
    lines.append(f"type {name} struct {{")
    used_go: set[str] = set()
    for p in props:
        field = go_field(p)
        if field in used_go:
            i = 2
            while f"{field}{i}" in used_go:
                i += 1
            field = f"{field}{i}"
        used_go.add(field)
        tag = f'ros:"{p}'
        if p == ".id":
            tag += ",key"
        if p in required:
            tag += ",required"
        if omitempty:
            tag += ",omitempty"
        tag += '"'
        comment = descs.get(p, "")
        if comment:
            lines.append(f"\t// {comment}")
        lines.append(f"\t{field:<28} string `{tag}`")
    lines.append("}")
    lines.append("```")
    lines.append("")


def main() -> None:
    ap = argparse.ArgumentParser()
    ap.add_argument("--probe", default="print_fields.json", help="file hasil probe (relatif terhadap mikrotik_inspect/)")
    ap.add_argument("--out", default="SOURCE.md", help="file keluaran (relatif terhadap mikrotik_inspect/)")
    args = ap.parse_args()

    index = load(INSPECT / "_index.json")
    probe = load(INSPECT / args.probe)["paths"]
    catalog = {r["path"]: r for r in load(INSPECT / "catalog_dump.json")}

    by_file = {f["file"]: f for f in index["files"]}
    actual = sorted(
        p.name for p in INSPECT.glob("*.json") if p.name not in GENERATED
    )
    parse_errors: list[str] = []
    entries: list[tuple[dict, dict]] = []
    seen_paths: set[str] = set()

    # 1) Semua file JSON yang benar-benar ada di direktori.
    for name in actual:
        f = by_file.get(name, {"file": name, "path": "/" + name[:-5], "scope": "menu"})
        try:
            entries.append((f, load(INSPECT / name)))
            seen_paths.add(f["path"])
        except Exception as e:  # noqa: BLE001
            parse_errors.append(f"{name}: {e}")

    # 2) Path dari index yang filenya tidak ikut ter-attach tetapi punya data
    #    print dari router: tetap dibuatkan struct Row dari hasil print.
    for f in index["files"]:
        path = f["path"]
        if f["file"] in actual or path in seen_paths or f["scope"] != "menu":
            continue
        pr = probe.get(path, {})
        if pr.get("fields") or pr.get("completion"):
            entries.append((f, {}))
            seen_paths.add(path)

    # 3) Resource katalog yang pathnya tidak tercakup file inspect / data print.
    for path, cat in catalog.items():
        if path in seen_paths:
            continue
        pseudo: dict = {}
        for a in cat.get("actions", []):
            if a.get("name") in {"add", "set"}:
                names = a.get("required", []) + a.get("optional", [])
                pseudo[a["name"]] = {
                    "_type": "cmd",
                    **{p: {"_type": "arg"} for p in names},
                }
        if pseudo or any(v.get("proplist") for v in cat.get("views", [])):
            entries.append((
                {"file": "(dari katalog Go)", "path": path, "scope": "menu"},
                pseudo,
            ))
            seen_paths.add(path)

    entries.sort(key=lambda item: item[0]["path"])
    data_by_path = entries
    missing_count = sum(
        1 for f in index["files"] if f["file"] not in actual
    )

    lines: list[str] = []
    lines.append("# Sumber struct RouterOS (baseline v7.9+)")
    lines.append("")
    lines.append(
        "> File ini dihasilkan oleh `tools/gen_source_md.py`. "
        "Jangan diedit manual — regen dari sumber di bawah."
    )
    lines.append("")
    lines.append("## Sumber dan metadata")
    lines.append("")
    lines.append("| Sumber | Keterangan |")
    lines.append("|---|---|")
    lines.append("| `mikrotik_inspect/*.json` | command, argumen, dan `desc` dari `/console/inspect` |")
    probe_meta = load(INSPECT / args.probe)
    lines.append(
        f"| `{args.probe}` | kolom `print` live + completion dari router "
        f"{probe_meta.get('router', '?')} (diambil {probe_meta.get('captured_at', '?')}) |"
    )
    lines.append("| `catalog_dump.json` | proplist terverifikasi + required/optional dari katalog Go |")
    lines.append("")
    lines.append(
        "Router sumber memakai RouterOS v7.9+ (`/console/inspect` baru ada di v7.9). "
        "Command v7 relatif stabil, jadi struct di bawah dipakai sebagai baseline; "
        "perbedaan v6 ditangani terpisah lewat `Version` di katalog."
    )
    lines.append("")
    lines.append("Arti tag:")
    lines.append("")
    lines.append("- `ros:\"nama\"` → nama properti RouterOS yang dikirim/dibaca.")
    lines.append("- `,key` → kunci baris (default `.id`).")
    lines.append("- `,required` → wajib saat `add`/`set` (dari katalog).")
    lines.append("- `,omitempty` → field kosong tidak dikirim (struct tulis).")
    lines.append("")
    lines.append("`XxxRow` dipakai untuk proplist `print`; `XxxAdd` / `XxxSet` dipakai untuk argumen tulis.")
    lines.append("")

    used: set[str] = set()
    stats = Counter()
    index_rows: list[str] = []
    sections: list[str] = []

    for f, data in data_by_path:
        path = f["path"]
        cat = catalog.get(path)
        pr = probe.get(path, {})
        props, meta = schema_of(path, data, cat, pr)
        cmds = cmds_of(data)
        add_args = args_of(cmds.get("add", {})) if "add" in cmds else []
        set_args = args_of(cmds.get("set", {})) if "set" in cmds else []
        req = required_of(cat)

        real_fields = [p for p in props if p != ".id"]
        has_real = bool(real_fields)
        if has_real:
            stats["dengan field"] += 1
        else:
            stats["tanpa field"] += 1
        if add_args:
            stats["punya add"] += 1
        if set_args:
            stats["punya set"] += 1

        src = []
        if meta["live"]:
            src.append(f"live {meta['live']}")
        if meta["comp"]:
            src.append(f"completion {meta['comp']}")
        if meta["cat"]:
            src.append(f"katalog {meta['cat']}")
        source_label = ", ".join(src) if src else "tidak ada kolom terdeteksi"
        index_rows.append(
            f"| `{path}` | {len(real_fields)} | {len(add_args)} | {len(set_args)} | {source_label} |"
        )

        s = [f"### `{path}`", ""]
        s.append(f"- File: `{f['file']}`")
        if cat:
            s.append(f"- Katalog: `{cat['id']}`")
        s.append(f"- Sumber field: {source_label}")
        if pr.get("error"):
            s.append(f"- Catatan print live: {clean_desc(pr['error'])}")
        if pr.get("completion_error"):
            s.append(f"- Catatan completion: {clean_desc(pr['completion_error'])}")
        s.append("")
        s.append(
            "Command: "
            + ", ".join(f"`{c}`({len(args_of(v))})" for c, v in sorted(cmds.items()))
            or "Command: -"
        )
        s.append("")

        if has_real:
            rown = go_type(path, "Row", used)
            emit_struct(
                s, rown, props, meta["descs"], set(),
                False, "adalah kolom hasil `/" + path.strip("/") + "/print`.",
            )
        if add_args:
            addn = go_type(path, "Add", used)
            emit_struct(
                s, addn, add_args, meta["descs"], req,
                True, "adalah argumen `/" + path.strip("/") + "/add`.",
            )
        if set_args:
            setn = go_type(path, "Set", used)
            emit_struct(
                s, setn, set_args, meta["descs"], set(),
                True, "adalah argumen `/" + path.strip("/") + "/set`.",
            )
        if not has_real and not add_args and not set_args:
            s.append("_Tidak ada struct: menu ini bukan tabel yang bisa di-print "
                     "atau tidak tersedia di router uji (mis. paket tidak terpasang)._")
            s.append("")
        sections.extend(s)

    lines.append("## Ringkasan")
    lines.append("")
    lines.append(f"- Entries diproses: {len(data_by_path)}")
    for k, v in stats.items():
        lines.append(f"- {k}: {v}")
    lines.append(f"- File JSON hadir: {len(actual)}")
    lines.append(f"- File di index tapi tidak hadir: {missing_count}")
    lines.append(f"- Error baca file: {len(parse_errors)}")
    lines.append("")
    lines.append(
        "Catatan: `_index.json` memuat 243 entri file, tetapi hanya sebagian "
        "yang ikut ter-attach di direktori ini. Untuk path yang filenya tidak "
        "ada tetapi punya hasil `print` dari router, tetap dibuatkan struct "
        "`Row` berdasarkan data print (tanpa `Add`/`Set` karena argumennya "
        "tidak tersedia)."
    )
    lines.append("")
    if parse_errors:
        for e in parse_errors:
            lines.append(f"- {e}")
        lines.append("")

    lines.append("## Indeks menu")
    lines.append("")
    lines.append("| Path | Field | Add | Set | Sumber |")
    lines.append("|---|---:|---:|---:|---|")
    lines.extend(index_rows)
    lines.append("")
    lines.append("## Struct per menu")
    lines.append("")
    lines.extend(sections)

    target = INSPECT / args.out
    target.write_text("\n".join(lines), encoding="utf-8")
    print(f"ditulis: {target} ({len(lines)} baris, {len(data_by_path)} menu)")


if __name__ == "__main__":
    main()
