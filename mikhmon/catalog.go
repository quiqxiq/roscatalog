// Package mikhmon adalah overlay katalog dan workflow khusus Mikhmon v4.
//
// Paket ini tidak menduplikasi definisi path/verb RouterOS. Sumber
// kebenarannya tetap katalog root roscatalog; paket ini hanya:
//
//   - menyaring resource yang benar-benar dipakai Mikhmon,
//   - menyediakan builder script on-login dan Expire Monitor,
//   - menyediakan format laporan penjualan dan QuickPrint,
//   - membungkus aksi umum menjadi Client yang bisa dijalankan lewat rosrun.
//
// Eksekusi tetap diserahkan ke Executor (biasanya *rosrun.Runner), sehingga
// paket ini tidak bergantung langsung pada goros.
package mikhmon

import (
	"fmt"

	ros "example.com/roscatalog"
)

const (
	// ExpireMonitorName adalah nama scheduler yang dipasang Mikhmon.
	ExpireMonitorName = "Mikhmon-Expire-Monitor"
	// ExpireMonitorComment adalah comment scheduler tersebut.
	ExpireMonitorComment = "Mikhmon Expire Monitor"
	// SalesComment menandai script laporan penjualan Mikhmon.
	SalesComment = "mikhmon"
	// QuickPrintComment menandai preset QuickPrint Mikhmon.
	QuickPrintComment = "QuickPrintMikhmon"
	// ExpiredLimitUptime adalah penanda user yang sudah kedaluwarsa.
	ExpiredLimitUptime = "1s"
)

// ids adalah resource root yang dipakai Mikhmon v4. Urutannya mengikuti peta
// di MIKHMON.md; resource lain sengaja tidak diekspos.
var ids = []string{
	// Dashboard dan monitoring.
	"system-resource",
	"system-routerboard",
	"system-health",
	"system-clock",
	"system-identity",
	"interface",
	"interface-monitor-traffic",
	"log",

	// Hotspot: server, user, profile, sesi, host, binding, cookie.
	"hotspot-server",
	"hotspot-user",
	"hotspot-user-profile",
	"hotspot-active",
	"hotspot-host",
	"hotspot-ip-binding",
	"hotspot-cookie",

	// DHCP, pool, NAT, queue.
	"dhcp-lease",
	"ip-pool",
	"fw-nat",
	"queue-simple",

	// Scheduler, script, dan laporan penjualan Mikhmon.
	"system-scheduler",
	"system-script",
	"mikhmon-sales",
}

var idSet = func() map[string]bool {
	m := make(map[string]bool, len(ids))
	for _, id := range ids {
		m[id] = true
	}
	return m
}()

// IDs mengembalikan salinan daftar resource yang dipakai Mikhmon.
func IDs() []string {
	return append([]string(nil), ids...)
}

// IsResource melaporkan apakah id termasuk subset Mikhmon.
func IsResource(id string) bool {
	return idSet[id]
}

// Resources mengembalikan resource katalog root yang termasuk subset Mikhmon.
// Resource yang tidak ada di root dilewati, sehingga aman bila katalog root
// berubah.
func Resources() []*ros.Resource {
	out := make([]*ros.Resource, 0, len(ids))
	for _, id := range ids {
		if r, ok := ros.Lookup(id); ok {
			out = append(out, r)
		}
	}
	return out
}

// Lookup mencari resource, tetapi hanya di subset Mikhmon.
func Lookup(id string) (*ros.Resource, bool) {
	if !idSet[id] {
		return nil, false
	}
	return ros.Lookup(id)
}

func resource(id string) (*ros.Resource, error) {
	r, ok := Lookup(id)
	if !ok {
		return nil, fmt.Errorf("mikhmon: resource %q tidak ada di subset Mikhmon", id)
	}
	return r, nil
}
