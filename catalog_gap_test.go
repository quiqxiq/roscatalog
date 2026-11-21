package roscatalog

import "testing"

// resourceBaru is a ringan rekap resource yang ditambahkan untuk menutup
// gap dengan MikroDash (keanggotaan menjadi "terinstall" secara CATALOG).
var resourceBaru = []string{
	"app", "app-settings", "user-group", "user-settings", "device-mode", "license",
	"routerboard-settings", "disk", "wifi-ap", "wifi-channel", "wifi-configuration",
	"wifi-security", "wifi-datapath", "wifi-provisioning", "wifi-radio",
	"wifi-registration-table", "wifi-capsman", "wifi-capsman-channel",
	"wifi-capsman-configuration", "wifi-capsman-datapath", "wifi-capsman-interface",
	"wifi-capsman-provisioning", "wifi-capsman-registration-table",
	"wifi-capsman-remote-cap", "wifi-capsman-security",
	"wifi-legacy", "wifi-legacy-security-profiles", "wifi-legacy-registration-table",
	"interface-bridge-vlan", "interface-detect-internet-state", "wireguard-peers",
	"ipsec-policy", "ipsec-active-peers", "ipsec-installed-sa", "discovery-settings",
	"routing-rule", "ipv6-route", "ip-cloud",
	"capsman-manager", "capsman-radio", "interface-wifi-cap", "interface-wifi-capsman",
	"interface-wifi-capsman-remote-cap", "interface-wifi-frequency-scan",
	"console", "menu", "stub", "system-reset", "system-import", "tool-ping-status",
	"tool-sniffer",
}

// TestGapResourcesValid — seluruh resource baru punya definisi valid:
// path di bawah root, selector/action konsisten, view memakai proplist/fields.
func TestGapResourcesValid(t *testing.T) {
	for _, id := range resourceBaru {
		r, ok := Lookup(id)
		if !ok {
			t.Errorf("%s: resource tidak didefinisikan", id)
			continue
		}
		if err := r.Validate(); err != nil {
			t.Errorf("%s: validasi: %v", id, err)
		}
	}
}

// TestGapActionsPresent — write actions yang CHR ada diminta MikroDash,
// harus ada di catalog dengan selector yang tepat.
func TestGapActionsPresent(t *testing.T) {
	cases := []struct {
		id, action string
		wantSel    Selector
	}{
		{"user-group", "add", SelNone},
		{"user-group", "set", SelID},
		{"user-group", "remove", SelID},
		{"file", "add", SelNone},
		{"file", "set", SelID},
		{"system-package-update", "install", SelNone},
	}
	for _, c := range cases {
		r, ok := Lookup(c.id)
		if !ok {
			t.Errorf("%s: resource tidak didefinisikan", c.id)
			continue
		}
		a, ok := r.Action(c.action)
		if !ok {
			t.Errorf("%s: action %s tidak ada", c.id, c.action)
			continue
		}
		if a.Selector != c.wantSel {
			t.Errorf("%s/%s: selector %q, mau %q", c.id, c.action, a.Selector, c.wantSel)
		}
	}
}
