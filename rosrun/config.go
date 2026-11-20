package rosrun

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	ros "example.com/roscatalog"
)

const (
	// Port bawaan API RouterOS.
	defaultAPIPort    = 8728
	defaultAPISSLPort = 8729

	// DefaultDialTimeout membatasi dial+login bila Config.Timeout kosong.
	DefaultDialTimeout = 10 * time.Second
	// DefaultPollEvery adalah jeda poll bawaan bila Config.PollEvery kosong.
	DefaultPollEvery = 5 * time.Second
	// DefaultStreamEvery adalah 0: view interval memakai nilai katalog.
	DefaultStreamEvery time.Duration = 0
	// DefaultQueueDepth adalah nilai client.Queue bawaan.
	DefaultQueueDepth = 1000
	// DefaultMaxBuffered membatasi antrean per stream.
	DefaultMaxBuffered = 8192

	// DefaultStatusTarget adalah target ping status di sisi router.
	DefaultStatusTarget = "8.8.8.8"
	// DefaultStatusEvery adalah interval ping status bawaan.
	DefaultStatusEvery = time.Second
	// DefaultStatusTimeout adalah batas tanpa balasan ping sebelum status error.
	DefaultStatusTimeout = 15 * time.Second

	// DefaultCacheTTL menyimpan event terakhir agar subscriber baru langsung
	// mendapat snapshot terbaru.
	DefaultCacheTTL = 5 * time.Second
	// DefaultCacheMaxEntries membatasi jumlah feed yang di-cache per router.
	DefaultCacheMaxEntries = 64
)

// StatusConfig mengatur monitor status realtime per router.
//
// Monitor selalu memakai native stream /ping tanpa count/duration, terlepas
// dari Config.Stream. Ping dijalankan di koneksi feed yang sudah ada.
type StatusConfig struct {
	// Disable mematikan monitor status. Status() tetap melaporkan keadaan
	// koneksi, tetapi tanpa data ping.
	Disable bool

	// Target adalah alamat yang di-ping dari sisi router. Kosong = 8.8.8.8.
	Target string

	// Every adalah interval probe. Kosong = 1s. Batas 1s..5s.
	Every time.Duration

	// Timeout adalah lama tanpa balasan ping sukses sebelum state error.
	// Kosong = max(3*Every, 15s).
	Timeout time.Duration

	// WarnPacketLoss adalah ambang packet loss (%) untuk state warning.
	// 0 = 5%. Nilai negatif mematikan warning packet loss.
	WarnPacketLoss float64

	// SuppressTLSWarning mematikan warning "API tanpa TLS". Bawaannya
	// peringatan aktif karena kredensial API dikirim plaintext.
	SuppressTLSWarning bool
}

// CacheConfig mengatur shared feed dan snapshot cache per router.
type CacheConfig struct {
	// Disable mematikan cache; Subscribe tetap bekerja tanpa berbagi stream.
	Disable bool

	// TTL adalah umur event terakhir yang masih boleh disajikan ke subscriber
	// baru. Kosong = 5s.
	TTL time.Duration

	// MaxEntries membatasi jumlah feed per router. Kosong = 64.
	MaxEntries int
}

// Gate mengapit satu command yang benar-benar dikirim ke router. release
// dipanggil setelah command selesai (termasuk saat stream berakhir). Dipakai
// pemakai untuk limiter per router, audit, atau metrik.
//
// path adalah command API, mis. "/ip/hotspot/active/print".
type Gate func(ctx context.Context, path string) (release func(), err error)

// Config adalah parameter satu router: alamat, kredensial, dan perilaku
// koneksi. Satu Config selalu membuka dua koneksi (command dan feed stream).
type Config struct {
	// ID adalah identitas router untuk log, Status, dan Manager. Kosong =
	// alamat (host:port). Harus unik di dalam satu Manager.
	ID string

	// Address boleh "host" atau "host:port". Bila Port kosong, dipakai
	// 8728 (atau 8729 saat TLS).
	Address  string
	Port     int
	Username string
	Password string

	TLS      bool
	Insecure bool // hanya berlaku saat TLS; CHR memakai sertifikat sendiri

	// Timeout membatasi dial+login. Kosong = DefaultDialTimeout.
	Timeout time.Duration

	// Stream adalah preferensi delivery bawaan: true memakai stream, false
	// memakai poll. Plan eksplisit pada Request selalu menang. Lihat PlanFor.
	Stream bool

	// PollEvery adalah jeda poll bawaan. Kosong = DefaultPollEvery.
	// Nilai harus 5s..30s.
	PollEvery time.Duration

	// StreamEvery adalah interval stream bawaan untuk view interval/native
	// yang mendukung =interval=. Kosong = interval bawaan view. Nilai harus
	// 1s..5s.
	StreamEvery time.Duration

	// Status mengatur monitor status realtime lewat native /ping.
	Status StatusConfig

	// Cache mengatur shared feed dan snapshot cache.
	Cache CacheConfig

	// Queue adalah client.Queue (buffer sentence yang dialokasikan Listen).
	// Kosong = DefaultQueueDepth.
	Queue int

	// MaxBuffered membatasi antrean per stream. Kosong = DefaultMaxBuffered.
	MaxBuffered int

	// EmptySnapshots meneruskan Runner.EmptySnapshots.
	EmptySnapshots bool

	// Gate opsional, dipanggil sekali per command.
	Gate Gate

	// Logf opsional.
	Logf func(format string, args ...any)

	// backoff adalah jeda reconnect pertama. Tidak diekspos; test memakai
	// nilai kecil supaya cepat.
	backoff time.Duration
}

// Validate memeriksa parameter wajib. Nilai opsional yang kosong diisi
// denganDefaults saat Open.
func (c Config) Validate() error {
	if strings.TrimSpace(c.Address) == "" {
		return fmt.Errorf("rosrun: address is required")
	}
	if strings.TrimSpace(c.Username) == "" {
		return fmt.Errorf("rosrun: username is required")
	}
	if c.Port < 0 || c.Port > 65535 {
		return fmt.Errorf("rosrun: port %d out of range", c.Port)
	}
	if c.Timeout < 0 {
		return fmt.Errorf("rosrun: timeout must not be negative")
	}
	if c.PollEvery < 0 {
		return fmt.Errorf("rosrun: poll interval must not be negative")
	}
	if c.PollEvery > 0 {
		if c.PollEvery < ros.MinPollEvery {
			return fmt.Errorf("rosrun: poll interval %s below minimum %s", c.PollEvery, ros.MinPollEvery)
		}
		if c.PollEvery > ros.MaxPollEvery {
			return fmt.Errorf("rosrun: poll interval %s above maximum %s", c.PollEvery, ros.MaxPollEvery)
		}
	}
	if c.StreamEvery < 0 {
		return fmt.Errorf("rosrun: stream interval must not be negative")
	}
	if c.StreamEvery > 0 {
		if c.StreamEvery < ros.MinPrintInterval {
			return fmt.Errorf("rosrun: stream interval %s below minimum %s", c.StreamEvery, ros.MinPrintInterval)
		}
		if c.StreamEvery > ros.MaxPrintInterval {
			return fmt.Errorf("rosrun: stream interval %s above maximum %s", c.StreamEvery, ros.MaxPrintInterval)
		}
	}
	if c.Status.Every < 0 {
		return fmt.Errorf("rosrun: status interval must not be negative")
	}
	if c.Status.Every > 0 {
		if c.Status.Every < ros.MinPrintInterval {
			return fmt.Errorf("rosrun: status interval %s below minimum %s", c.Status.Every, ros.MinPrintInterval)
		}
		if c.Status.Every > ros.MaxPrintInterval {
			return fmt.Errorf("rosrun: status interval %s above maximum %s", c.Status.Every, ros.MaxPrintInterval)
		}
	}
	if c.Status.Timeout < 0 {
		return fmt.Errorf("rosrun: status timeout must not be negative")
	}
	if c.Status.Timeout > 0 && c.Status.Every > 0 && c.Status.Timeout < c.Status.Every {
		return fmt.Errorf("rosrun: status timeout %s below status interval %s", c.Status.Timeout, c.Status.Every)
	}
	if c.Cache.TTL < 0 {
		return fmt.Errorf("rosrun: cache ttl must not be negative")
	}
	if c.Cache.MaxEntries < 0 {
		return fmt.Errorf("rosrun: cache max entries must not be negative")
	}
	if c.Queue < 0 || c.MaxBuffered < 0 {
		return fmt.Errorf("rosrun: queue and max buffered must not be negative")
	}
	return nil
}

// withDefaults mengembalikan salinan Config dengan nilai opsional terisi.
func (c Config) withDefaults() Config {
	if c.Timeout == 0 {
		c.Timeout = DefaultDialTimeout
	}
	if c.PollEvery == 0 {
		c.PollEvery = DefaultPollEvery
	}
	if c.StreamEvery == 0 {
		c.StreamEvery = DefaultStreamEvery
	}
	if c.Status.Target == "" {
		c.Status.Target = DefaultStatusTarget
	}
	if c.Status.Every == 0 {
		c.Status.Every = DefaultStatusEvery
	}
	if c.Status.Timeout == 0 {
		c.Status.Timeout = 3 * c.Status.Every
		if c.Status.Timeout < DefaultStatusTimeout {
			c.Status.Timeout = DefaultStatusTimeout
		}
	}
	if c.Status.WarnPacketLoss == 0 {
		c.Status.WarnPacketLoss = 5
	}
	if c.Cache.TTL == 0 {
		c.Cache.TTL = DefaultCacheTTL
	}
	if c.Cache.MaxEntries == 0 {
		c.Cache.MaxEntries = DefaultCacheMaxEntries
	}
	if c.Queue == 0 {
		c.Queue = DefaultQueueDepth
	}
	if c.MaxBuffered == 0 {
		c.MaxBuffered = DefaultMaxBuffered
	}
	return c
}

// addr mengembalikan alamat "host:port" siap dial.
func (c Config) addr() string {
	host := strings.TrimSpace(c.Address)
	if h, p, err := net.SplitHostPort(host); err == nil && h != "" && p != "" {
		return host
	}

	port := c.Port
	if port == 0 {
		if c.TLS {
			port = defaultAPISSLPort
		} else {
			port = defaultAPIPort
		}
	}
	return net.JoinHostPort(host, strconv.Itoa(port))
}
