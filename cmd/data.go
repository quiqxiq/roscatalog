package cmd

import (
	"os"
	"strconv"
	"time"

	"example.com/roscatalog/rosrun"
)

// Konfigurasi contoh untuk cmd/multirunner, cmd/webtest, dan cmd/mikhmon-live.
//
// Jangan menaruh kredensial produksi di file ini. Semua nilai bisa di-override
// lewat environment variable, sehingga file tetap aman di-commit:
//
//	ROSCATALOG_ROUTER1_ADDR / _PORT / _USER / _PASS / _TLS
//	ROSCATALOG_ROUTER2_ADDR / _PORT / _USER / _PASS / _TLS
//	ROSCATALOG_ROUTER3_ADDR / _PORT / _USER / _PASS / _TLS
//
// Contoh:
//
//	ROSCATALOG_ROUTER1_ADDR=192.168.88.1:8728 \
//	ROSCATALOG_ROUTER1_USER=admin \
//	ROSCATALOG_ROUTER1_PASS=secret \
//	go run ./cmd/multicheck
type routerEnv struct {
	addr string
	port string
	user string
	pass string
	tls  string
}

func fromEnv(def rosrun.Config, env routerEnv) rosrun.Config {
	if v := os.Getenv(env.addr); v != "" {
		def.Address = v
	}
	if v := os.Getenv(env.port); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			def.Port = p
		}
	}
	if v := os.Getenv(env.user); v != "" {
		def.Username = v
	}
	if v := os.Getenv(env.pass); v != "" {
		def.Password = v
	}
	if v := os.Getenv(env.tls); v != "" {
		def.TLS = v == "1" || v == "true" || v == "yes"
	}
	return def
}

// MikrotikStreaming adalah contoh konfigurasi router dengan mode streaming
// (realtime). Poll=false adalah bawaan: interval, follow, listen, atau command
// native dipilih sesuai view.
var MikrotikStreaming = fromEnv(rosrun.Config{
	ID:       "router-streaming-01",
	Address:  "192.168.88.1:8728",
	Username: "admin",
	Timeout:  10 * time.Second,

	StreamEvery: time.Second, // interval stream (1s..5s)

	MaxInFlight: 8,
	Grace:       30 * time.Second,
	StreamStale: 10 * time.Second,
	StreamCheck: 5 * time.Second,

	Status: rosrun.StatusConfig{
		Target:         "8.8.8.8",
		Every:          time.Second,
		Timeout:        15 * time.Second,
		WarnPacketLoss: 5.0,
	},
	Cache: rosrun.CacheConfig{
		TTL:        5 * time.Second,
		MaxEntries: 64,
	},

	Queue:       1000,
	MaxBuffered: 8192,
}, routerEnv{
	addr: "ROSCATALOG_ROUTER1_ADDR", port: "ROSCATALOG_ROUTER1_PORT",
	user: "ROSCATALOG_ROUTER1_USER", pass: "ROSCATALOG_ROUTER1_PASS",
	tls: "ROSCATALOG_ROUTER1_TLS",
})

// MikrotikPolling adalah contoh konfigurasi router dengan polling dari Go.
var MikrotikPolling = fromEnv(rosrun.Config{
	ID:       "router-polling-02",
	Address:  "192.168.88.2:8728",
	Username: "admin",
	Timeout:  10 * time.Second,

	Poll:      true,
	PollEvery: 5 * time.Second, // jeda antar polling (5s..30s)

	MaxInFlight: 8,
	Grace:       30 * time.Second,

	Status: rosrun.StatusConfig{
		Target:         "1.1.1.1",
		Every:          2 * time.Second,
		Timeout:        15 * time.Second,
		WarnPacketLoss: 5.0,
	},
	Cache: rosrun.CacheConfig{
		TTL:        5 * time.Second,
		MaxEntries: 64,
	},

	Queue:       1000,
	MaxBuffered: 8192,
}, routerEnv{
	addr: "ROSCATALOG_ROUTER2_ADDR", port: "ROSCATALOG_ROUTER2_PORT",
	user: "ROSCATALOG_ROUTER2_USER", pass: "ROSCATALOG_ROUTER2_PASS",
	tls: "ROSCATALOG_ROUTER2_TLS",
})

// MikrotikLab233 adalah contoh konfigurasi router lab.
var MikrotikLab233 = fromEnv(rosrun.Config{
	ID:       "router-lab-233",
	Address:  "192.168.233.1:8728",
	Username: "admin",
	Timeout:  10 * time.Second,

	StreamEvery: time.Second,

	MaxInFlight: 8,
	Grace:       30 * time.Second,
	StreamStale: 10 * time.Second,
	StreamCheck: 5 * time.Second,

	Status: rosrun.StatusConfig{
		Target:         "8.8.8.8",
		Every:          time.Second,
		Timeout:        15 * time.Second,
		WarnPacketLoss: 5.0,
	},
	Cache: rosrun.CacheConfig{
		TTL:        5 * time.Second,
		MaxEntries: 64,
	},

	Queue:       1000,
	MaxBuffered: 8192,
}, routerEnv{
	addr: "ROSCATALOG_ROUTER3_ADDR", port: "ROSCATALOG_ROUTER3_PORT",
	user: "ROSCATALOG_ROUTER3_USER", pass: "ROSCATALOG_ROUTER3_PASS",
	tls: "ROSCATALOG_ROUTER3_TLS",
})

// MikrotikList berisi kumpulan konfigurasi router contoh yang dapat langsung
// dipakai atau disesuaikan. Kredensial diambil dari environment bila diisi.
var MikrotikList = []rosrun.Config{
	MikrotikStreaming,
	MikrotikPolling,
	MikrotikLab233,
}
