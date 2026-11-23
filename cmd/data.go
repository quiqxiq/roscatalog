package cmd

import (
	"time"

	"example.com/roscatalog/rosrun"
)

// MikrotikStreaming adalah contoh konfigurasi router dengan mode streaming (realtime).
// Sesuai config.go: Poll bernilai false (bawaan stream) untuk mendengarkan feed
// realtime (misal: interval, follow, follow-only, listen, atau monitor).
var MikrotikStreaming = rosrun.Config{
	ID:       "router-streaming-01",
	Address:  "103.139.193.128",
	Port:     1012,
	Username: "fandi1",
	Password: "001",
	TLS:      false,
	Insecure: false,
	Timeout:  10 * time.Second, // Batas waktu dial + login

	// Mode Streaming: Poll=false (default stream/realtime)
	Poll:        false,
	StreamEvery: 1 * time.Second, // Interval stream (1s..5s)

	// Batasan koneksi & concurrency
	MaxInFlight: 8,
	Grace:       30 * time.Second,
	StreamStale: 10 * time.Second,
	StreamCheck: 5 * time.Second,

	// Monitor status router realtime menggunakan native /ping
	Status: rosrun.StatusConfig{
		Disable:            false,
		Target:             "8.8.8.8",
		Every:              1 * time.Second,  // Interval probe ping (1s..5s)
		Timeout:            15 * time.Second, // Ambang batas timeout ping
		WarnPacketLoss:     5.0,              // Ambang batas % packet loss
		SuppressTLSWarning: false,
	},

	// Pengaturan cache shared feed & snapshot
	Cache: rosrun.CacheConfig{
		Disable:    false,
		TTL:        5 * time.Second,
		MaxEntries: 64,
	},

	Queue:       1000,
	MaxBuffered: 8192,
}

// MikrotikPolling adalah contoh konfigurasi router dengan mode polling berkala.
// Sesuai config.go: Poll bernilai true memaksa penarikan data secara berkala (polling)
// dari sisi aplikasi Go.
var MikrotikPolling = rosrun.Config{
	ID:       "router-polling-02",
	Address:  "103.147.134.201",
	Port:     1436,
	Username: "99barokah99",
	Password: "srbarokah99",
	TLS:      false,
	Insecure: false,
	Timeout:  10 * time.Second, // Batas waktu dial + login

	// Mode Polling: Poll=true (memaksa polling dari sisi Go)
	Poll:      true,
	PollEvery: 5 * time.Second, // Jeda antar polling (5s..30s)

	// Batasan koneksi & concurrency
	MaxInFlight: 8,
	Grace:       30 * time.Second,

	// Monitor status router realtime menggunakan native /ping
	Status: rosrun.StatusConfig{
		Disable:            false,
		Target:             "1.1.1.1",
		Every:              2 * time.Second,  // Interval probe ping (1s..5s)
		Timeout:            15 * time.Second, // Ambang batas timeout ping
		WarnPacketLoss:     5.0,              // Ambang batas % packet loss
		SuppressTLSWarning: false,
	},

	// Pengaturan cache shared feed & snapshot
	Cache: rosrun.CacheConfig{
		Disable:    false,
		TTL:        5 * time.Second,
		MaxEntries: 64,
	},

	Queue:       1000,
	MaxBuffered: 8192,
}

// MikrotikLab233 adalah konfigurasi router lab 192.168.233.1 dengan kredensial admin / r00t.
var MikrotikLab233 = rosrun.Config{
	ID:       "router-lab-233",
	Address:  "192.168.233.1",
	Port:     8728,
	Username: "admin",
	Password: "r00t",
	TLS:      false,
	Insecure: false,
	Timeout:  10 * time.Second,

	// Mode Streaming: Poll=false (default stream/realtime)
	Poll:        false,
	StreamEvery: 1 * time.Second,

	// Batasan koneksi & concurrency
	MaxInFlight: 8,
	Grace:       30 * time.Second,
	StreamStale: 10 * time.Second,
	StreamCheck: 5 * time.Second,

	// Monitor status router realtime menggunakan native /ping
	Status: rosrun.StatusConfig{
		Disable:            false,
		Target:             "8.8.8.8",
		Every:              1 * time.Second,
		Timeout:            15 * time.Second,
		WarnPacketLoss:     5.0,
		SuppressTLSWarning: false,
	},

	// Pengaturan cache shared feed & snapshot
	Cache: rosrun.CacheConfig{
		Disable:    false,
		TTL:        5 * time.Second,
		MaxEntries: 64,
	},

	Queue:       1000,
	MaxBuffered: 8192,
}

// MikrotikList berisi kumpulan konfigurasi router (streaming & polling)
// yang dapat langsung dipakai atau disesuaikan.
var MikrotikList = []rosrun.Config{
	MikrotikStreaming,
	MikrotikPolling,
	MikrotikLab233,
}
