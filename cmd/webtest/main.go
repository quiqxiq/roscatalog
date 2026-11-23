package main

import (
	"context"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	ros "example.com/roscatalog"
	"example.com/roscatalog/cmd"
	"example.com/roscatalog/rosrun"
)

//go:embed index.html
var indexHTML embed.FS

// SystemMetrics menyimpan data ringkasan resource sistem MikroTik.
type SystemMetrics struct {
	Uptime       string  `json:"uptime"`
	Version      string  `json:"version"`
	CPULoad      int     `json:"cpu_load"`
	FreeMemory   int64   `json:"free_memory"`
	TotalMemory  int64   `json:"total_memory"`
	MemoryUsedPct float64 `json:"memory_used_pct"`
	FreeHDD      int64   `json:"free_hdd"`
	TotalHDD     int64   `json:"total_hdd"`
	HDDUsedPct   float64 `json:"hdd_used_pct"`
	BoardName    string  `json:"board_name"`
	Architecture string  `json:"architecture"`
	Platform     string  `json:"platform"`
	CPU          string  `json:"cpu"`
}

// RouterConfigDTO adalah DTO konfigurasi router yang aman untuk JSON serialization.
type RouterConfigDTO struct {
	ID          string  `json:"id"`
	Address     string  `json:"address"`
	Port        int     `json:"port"`
	Username    string  `json:"username"`
	Password    string  `json:"password"`
	TLS         bool    `json:"tls"`
	Insecure    bool    `json:"insecure"`
	Poll        bool    `json:"poll"`
	StreamEvery float64 `json:"stream_every_s"`
	PollEvery   float64 `json:"poll_every_s"`
	PingTarget  string  `json:"ping_target"`
}

func toConfigDTO(cfg rosrun.Config) RouterConfigDTO {
	target := cfg.Status.Target
	if target == "" {
		target = "8.8.8.8"
	}
	return RouterConfigDTO{
		ID:          cfg.ID,
		Address:     cfg.Address,
		Port:        cfg.Port,
		Username:    cfg.Username,
		Password:    cfg.Password,
		TLS:         cfg.TLS,
		Insecure:    cfg.Insecure,
		Poll:        cfg.Poll,
		StreamEvery: cfg.StreamEvery.Seconds(),
		PollEvery:   cfg.PollEvery.Seconds(),
		PingTarget:  target,
	}
}

// RouterState adalah model state router yang dikirim ke frontend.
type RouterState struct {
	Config       RouterConfigDTO `json:"config"`
	State        string          `json:"state"` // "connected", "warning", "error", "offline"
	Connected    bool            `json:"connected"`
	Warnings     []string        `json:"warnings,omitempty"`
	LastError    string          `json:"last_error,omitempty"`
	PingTarget   string          `json:"ping_target"`
	PingLastRTT  float64         `json:"ping_last_rtt_ms"`
	PingAvgRTT   float64         `json:"ping_avg_rtt_ms"`
	PingMinRTT   float64         `json:"ping_min_rtt_ms"`
	PingMaxRTT   float64         `json:"ping_max_rtt_ms"`
	PacketLoss   float64         `json:"packet_loss"`
	PingSent     int             `json:"ping_sent"`
	PingReceived int             `json:"ping_received"`
	PingStatus   string          `json:"ping_status"`
	System       *SystemMetrics  `json:"system,omitempty"`
	LastSeen     time.Time       `json:"last_seen"`
}

// RouterHub mengelola router, Manager, background metrics, dan SSE broadcast.
type RouterHub struct {
	mu             sync.RWMutex
	mgr            *rosrun.Manager
	routers        map[string]*RouterState
	rawCfgs        map[string]rosrun.Config
	watcherCancels map[string]context.CancelFunc
	connecting     map[string]bool
	order          []string
	ctx            context.Context
	cancel         context.CancelFunc

	// SSE subscribers untuk health & ping
	healthMu   sync.Mutex
	healthSubs map[chan HealthEvent]struct{}
}

// HealthEvent disiarkan via SSE /api/events/health
type HealthEvent struct {
	RouterID     string         `json:"router_id"`
	State        string         `json:"state"`
	Connected    bool           `json:"connected"`
	Warnings     []string       `json:"warnings,omitempty"`
	LastError    string         `json:"last_error,omitempty"`
	PingTarget   string         `json:"ping_target"`
	PingLastRTT  float64        `json:"ping_last_rtt_ms"`
	PingAvgRTT   float64        `json:"ping_avg_rtt_ms"`
	PingMinRTT   float64        `json:"ping_min_rtt_ms"`
	PingMaxRTT   float64        `json:"ping_max_rtt_ms"`
	PacketLoss   float64        `json:"packet_loss"`
	PingSent     int            `json:"ping_sent"`
	PingReceived int            `json:"ping_received"`
	PingStatus   string         `json:"ping_status"`
	System       *SystemMetrics `json:"system,omitempty"`
	Timestamp    time.Time      `json:"timestamp"`
}

func applyRosStatus(r *RouterState, st rosrun.Status) {
	r.State = string(st.State)
	r.Connected = st.Connected
	r.Warnings = append([]string(nil), st.Warnings...)
	if st.Err != nil {
		r.LastError = st.Err.Error()
	} else {
		r.LastError = ""
	}
	if st.Target != "" {
		r.PingTarget = st.Target
	}
	r.PingLastRTT = float64(st.Ping.LastRTT.Microseconds()) / 1000.0
	r.PingAvgRTT = float64(st.Ping.AvgRTT.Microseconds()) / 1000.0
	r.PingMinRTT = float64(st.Ping.MinRTT.Microseconds()) / 1000.0
	r.PingMaxRTT = float64(st.Ping.MaxRTT.Microseconds()) / 1000.0
	r.PacketLoss = st.Ping.PacketLoss
	r.PingSent = st.Ping.Sent
	r.PingReceived = st.Ping.Received
	r.PingStatus = st.Ping.LastStatus
	r.LastSeen = time.Now()
}

func NewRouterHub(parentCtx context.Context) *RouterHub {
	ctx, cancel := context.WithCancel(parentCtx)
	mgr, _ := rosrun.OpenAll(ctx, nil)

	hub := &RouterHub{
		mgr:            mgr,
		routers:        make(map[string]*RouterState),
		rawCfgs:        make(map[string]rosrun.Config),
		watcherCancels: make(map[string]context.CancelFunc),
		connecting:     make(map[string]bool),
		ctx:            ctx,
		cancel:         cancel,
		healthSubs:     make(map[chan HealthEvent]struct{}),
	}

	go hub.autoReconnectLoop()
	go hub.periodicallyUpdateResources()

	return hub
}

// Close menutup seluruh sesi router dan background worker.
func (h *RouterHub) Close() {
	h.cancel()
	h.mu.Lock()
	for _, cancel := range h.watcherCancels {
		if cancel != nil {
			cancel()
		}
	}
	if h.mgr != nil {
		_ = h.mgr.Close()
	}
	h.mu.Unlock()
}

func (h *RouterHub) attachWatcher(id string, sess *rosrun.MikroTik) {
	h.mu.Lock()
	if oldCancel, exists := h.watcherCancels[id]; exists && oldCancel != nil {
		oldCancel()
	}
	watchCtx, watchCancel := context.WithCancel(h.ctx)
	h.watcherCancels[id] = watchCancel
	h.mu.Unlock()

	ch := sess.WatchStatus(watchCtx)
	go func() {
		for {
			select {
			case <-watchCtx.Done():
				return
			case st, ok := <-ch:
				if !ok {
					return
				}
				h.mu.Lock()
				r, exists := h.routers[id]
				if exists {
					applyRosStatus(r, st)
				}
				var sysCopy *SystemMetrics
				var lastErr string
				if r != nil {
					lastErr = r.LastError
					if r.System != nil {
						sc := *r.System
						sysCopy = &sc
					}
				}
				h.mu.Unlock()

				h.broadcastHealth(HealthEvent{
					RouterID:     id,
					State:        string(st.State),
					Connected:    st.Connected,
					Warnings:     st.Warnings,
					LastError:    lastErr,
					PingTarget:   st.Target,
					PingLastRTT:  float64(st.Ping.LastRTT.Microseconds()) / 1000.0,
					PingAvgRTT:   float64(st.Ping.AvgRTT.Microseconds()) / 1000.0,
					PingMinRTT:   float64(st.Ping.MinRTT.Microseconds()) / 1000.0,
					PingMaxRTT:   float64(st.Ping.MaxRTT.Microseconds()) / 1000.0,
					PacketLoss:   st.Ping.PacketLoss,
					PingSent:     st.Ping.Sent,
					PingReceived: st.Ping.Received,
					PingStatus:   st.Ping.LastStatus,
					System:       sysCopy,
					Timestamp:    time.Now(),
				})
			}
		}
	}()
}

// autoReconnectLoop mencoba menghubungkan kembali router yang offline/gagal dial.
func (h *RouterHub) autoReconnectLoop() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-h.ctx.Done():
			return
		case <-ticker.C:
			h.mu.Lock()
			var toRetry []rosrun.Config
			for id, rawCfg := range h.rawCfgs {
				if _, ok := h.mgr.Get(id); !ok {
					if !h.connecting[id] {
						h.connecting[id] = true
						toRetry = append(toRetry, rawCfg)
					}
				}
			}
			h.mu.Unlock()

			for _, cfg := range toRetry {
				go func(c rosrun.Config) {
					defer func() {
						h.mu.Lock()
						delete(h.connecting, c.ID)
						h.mu.Unlock()
					}()

					sess, err := h.mgr.Add(h.ctx, c)
					h.mu.Lock()
					r, exists := h.routers[c.ID]
					if !exists {
						h.mu.Unlock()
						if sess != nil {
							_ = sess.Close()
						}
						return
					}
					if err != nil {
						r.State = "error"
						r.Connected = false
						r.LastError = err.Error()
						h.mu.Unlock()
						return
					}
					if sess != nil {
						applyRosStatus(r, sess.Status())
						h.mu.Unlock()
						h.attachWatcher(c.ID, sess)
						log.Printf("[WebTest] Router %s (%s:%d) sekarang terhubung!", c.ID, c.Address, c.Port)
						h.broadcastHealth(HealthEvent{
							RouterID:  c.ID,
							State:     r.State,
							Connected: r.Connected,
							Timestamp: time.Now(),
						})
					} else {
						h.mu.Unlock()
					}
				}(cfg)
			}
		}
	}
}

// AddRouter mendaftarkan router dan mencoba membuka koneksi di latar belakang.
func (h *RouterHub) AddRouter(cfg rosrun.Config) (*RouterState, error) {
	if strings.TrimSpace(cfg.ID) == "" {
		return nil, fmt.Errorf("router ID tidak boleh kosong")
	}

	h.mu.Lock()
	if _, exists := h.routers[cfg.ID]; exists {
		h.mu.Unlock()
		return nil, fmt.Errorf("router dengan ID %q sudah ada", cfg.ID)
	}

	pingTarget := cfg.Status.Target
	if pingTarget == "" {
		pingTarget = "8.8.8.8"
	}
	state := &RouterState{
		Config:     toConfigDTO(cfg),
		State:      "connecting",
		PingTarget: pingTarget,
		LastSeen:   time.Now(),
	}

	h.routers[cfg.ID] = state
	h.rawCfgs[cfg.ID] = cfg
	h.order = append(h.order, cfg.ID)
	h.connecting[cfg.ID] = true
	cp := *state
	h.mu.Unlock()

	// Hubungkan ke MikroTik secara asynchronous agar tidak memblokir HTTP response & UI
	go func() {
		defer func() {
			h.mu.Lock()
			delete(h.connecting, cfg.ID)
			h.mu.Unlock()
		}()

		sess, err := h.mgr.Add(h.ctx, cfg)
		h.mu.Lock()
		r, exists := h.routers[cfg.ID]
		if !exists {
			h.mu.Unlock()
			if sess != nil {
				_ = sess.Close()
			}
			return
		}
		if err != nil {
			r.State = "error"
			r.Connected = false
			r.LastError = err.Error()
			h.mu.Unlock()
			h.broadcastHealth(HealthEvent{
				RouterID:  cfg.ID,
				State:     "error",
				Connected: false,
				LastError: err.Error(),
				Timestamp: time.Now(),
			})
			return
		}
		if sess != nil {
			applyRosStatus(r, sess.Status())
			h.mu.Unlock()
			h.attachWatcher(cfg.ID, sess)
			h.broadcastHealth(HealthEvent{
				RouterID:  cfg.ID,
				State:     r.State,
				Connected: r.Connected,
				Timestamp: time.Now(),
			})
		} else {
			h.mu.Unlock()
		}
	}()

	return &cp, nil
}

// UpdateRouter memperbarui konfigurasi router (re-dial di latar belakang).
func (h *RouterHub) UpdateRouter(id string, newCfg rosrun.Config) (*RouterState, error) {
	if strings.TrimSpace(newCfg.ID) == "" {
		newCfg.ID = id
	}

	h.mu.Lock()
	if _, exists := h.routers[id]; !exists {
		h.mu.Unlock()
		return nil, fmt.Errorf("router %q tidak ditemukan", id)
	}

	// Batalkan watcher lama jika ada
	if oldCancel, exists := h.watcherCancels[id]; exists && oldCancel != nil {
		oldCancel()
		delete(h.watcherCancels, id)
	}

	// Jika ID berubah, perbarui urutan & data
	if newCfg.ID != id {
		delete(h.routers, id)
		delete(h.rawCfgs, id)
		delete(h.connecting, id)
		for i, x := range h.order {
			if x == id {
				h.order[i] = newCfg.ID
				break
			}
		}
	}

	pingTarget := newCfg.Status.Target
	if pingTarget == "" {
		pingTarget = "8.8.8.8"
	}
	state := &RouterState{
		Config:     toConfigDTO(newCfg),
		State:      "connecting",
		PingTarget: pingTarget,
		LastSeen:   time.Now(),
	}

	h.routers[newCfg.ID] = state
	h.rawCfgs[newCfg.ID] = newCfg
	h.connecting[newCfg.ID] = true
	cp := *state
	h.mu.Unlock()

	// Dial ulang dengan config baru secara asinkron
	go func() {
		defer func() {
			h.mu.Lock()
			delete(h.connecting, newCfg.ID)
			h.mu.Unlock()
		}()

		_ = h.mgr.Remove(id)
		sess, err := h.mgr.Add(h.ctx, newCfg)
		h.mu.Lock()
		r, exists := h.routers[newCfg.ID]
		if !exists {
			h.mu.Unlock()
			if sess != nil {
				_ = sess.Close()
			}
			return
		}
		if err != nil {
			r.State = "error"
			r.Connected = false
			r.LastError = err.Error()
			h.mu.Unlock()
			h.broadcastHealth(HealthEvent{
				RouterID:  newCfg.ID,
				State:     "error",
				Connected: false,
				LastError: err.Error(),
				Timestamp: time.Now(),
			})
			return
		}
		if sess != nil {
			applyRosStatus(r, sess.Status())
			h.mu.Unlock()
			h.attachWatcher(newCfg.ID, sess)
			h.broadcastHealth(HealthEvent{
				RouterID:  newCfg.ID,
				State:     r.State,
				Connected: r.Connected,
				Timestamp: time.Now(),
			})
		} else {
			h.mu.Unlock()
		}
	}()

	return &cp, nil
}

// DeleteRouter menghapus router dari hub dan memutuskan koneksi.
func (h *RouterHub) DeleteRouter(id string) error {
	h.mu.Lock()
	if _, exists := h.routers[id]; !exists {
		h.mu.Unlock()
		return fmt.Errorf("router %q tidak ditemukan", id)
	}

	if cancel, exists := h.watcherCancels[id]; exists && cancel != nil {
		cancel()
		delete(h.watcherCancels, id)
	}

	delete(h.routers, id)
	delete(h.rawCfgs, id)
	delete(h.connecting, id)

	for i, x := range h.order {
		if x == id {
			h.order = append(h.order[:i], h.order[i+1:]...)
			break
		}
	}
	h.mu.Unlock()

	_ = h.mgr.Remove(id)
	return nil
}

// GetAllRouters mengembalikan seluruh status router sesuai urutan pendaftaran.
func (h *RouterHub) GetAllRouters() []*RouterState {
	h.mu.Lock()
	defer h.mu.Unlock()

	out := make([]*RouterState, 0, len(h.order))
	for _, id := range h.order {
		if r, ok := h.routers[id]; ok {
			if sess, ok := h.mgr.Get(id); ok && sess != nil {
				applyRosStatus(r, sess.Status())
			}
			cp := *r
			out = append(out, &cp)
		}
	}
	return out
}

// GetRouter mengembalikan status router berdasarkan ID.
func (h *RouterHub) GetRouter(id string) (*RouterState, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	r, ok := h.routers[id]
	if !ok {
		return nil, false
	}
	if sess, ok := h.mgr.Get(id); ok && sess != nil {
		applyRosStatus(r, sess.Status())
	}
	cp := *r
	return &cp, true
}

// watchHealth mendengarkan update status & ping dari rosrun.Manager
func (h *RouterHub) watchHealth() {
	ch, err := h.mgr.WatchStatuses(h.ctx)
	if err != nil {
		return
	}

	for {
		select {
		case <-h.ctx.Done():
			return
		case st, ok := <-ch:
			if !ok {
				return
			}
			h.updateRouterHealth(st)
		}
	}
}

func (h *RouterHub) updateRouterHealth(st rosrun.RouterStatus) {
	h.mu.Lock()
	r, exists := h.routers[st.RouterID]
	if !exists {
		h.mu.Unlock()
		return
	}

	status := st.Status
	r.State = string(status.State)
	r.PingLastRTT = status.Ping.LastRTT.Seconds() * 1000
	r.PingAvgRTT = status.Ping.AvgRTT.Seconds() * 1000
	r.PingMinRTT = status.Ping.MinRTT.Seconds() * 1000
	r.PingMaxRTT = status.Ping.MaxRTT.Seconds() * 1000
	r.PacketLoss = status.Ping.PacketLoss
	r.LastSeen = time.Now()
	if status.Err != nil {
		r.LastError = status.Err.Error()
	} else if len(status.Warnings) > 0 {
		r.LastError = strings.Join(status.Warnings, ", ")
	} else {
		r.LastError = ""
	}

	var sysCopy *SystemMetrics
	if r.System != nil {
		sc := *r.System
		sysCopy = &sc
	}
	h.mu.Unlock()

	// Broadcast ke SSE
	h.broadcastHealth(HealthEvent{
		RouterID:    st.RouterID,
		State:       string(status.State),
		PingLastRTT: status.Ping.LastRTT.Seconds() * 1000,
		PingAvgRTT:  status.Ping.AvgRTT.Seconds() * 1000,
		PacketLoss:  status.Ping.PacketLoss,
		LastError:   r.LastError,
		System:      sysCopy,
		Timestamp:   time.Now(),
	})
}

func (h *RouterHub) broadcastHealth(ev HealthEvent) {
	h.healthMu.Lock()
	defer h.healthMu.Unlock()

	for ch := range h.healthSubs {
		select {
		case ch <- ev:
		default:
		}
	}
}

// periodicallyUpdateResources memperbarui resource CPU/RAM/HDD tiap 5 detik
func (h *RouterHub) periodicallyUpdateResources() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-h.ctx.Done():
			return
		case <-ticker.C:
			h.refreshAllSystemResources()
		}
	}
}

func (h *RouterHub) refreshAllSystemResources() {
	ids := h.mgr.IDs()
	for _, id := range ids {
		go func(routerID string) {
			sess, ok := h.mgr.Get(routerID)
			if !ok || sess == nil {
				return
			}
			ctx, cancel := context.WithTimeout(h.ctx, 3*time.Second)
			defer cancel()

			res, ok := ros.Lookup("system-resource")
			if !ok {
				return
			}
			req, err := res.Build("list", ros.Once(), ros.Params{})
			if err != nil {
				return
			}
			rows, err := sess.Once(ctx, req)
			if err != nil || len(rows) == 0 {
				return
			}

			sys := parseSystemMetrics(rows[0])

			h.mu.Lock()
			var ev HealthEvent
			if r, found := h.routers[routerID]; found {
				r.System = sys
				r.LastSeen = time.Now()
				ev = HealthEvent{
					RouterID:     routerID,
					State:        r.State,
					Connected:    r.Connected,
					Warnings:     r.Warnings,
					LastError:    r.LastError,
					PingTarget:   r.PingTarget,
					PingLastRTT:  r.PingLastRTT,
					PingAvgRTT:   r.PingAvgRTT,
					PingMinRTT:   r.PingMinRTT,
					PingMaxRTT:   r.PingMaxRTT,
					PacketLoss:   r.PacketLoss,
					PingSent:     r.PingSent,
					PingReceived: r.PingReceived,
					PingStatus:   r.PingStatus,
					System:       sys,
					Timestamp:    time.Now(),
				}
			}
			h.mu.Unlock()

			if ev.RouterID != "" {
				h.broadcastHealth(ev)
			}
		}(id)
	}
}

func parseSystemMetrics(row map[string]string) *SystemMetrics {
	cpuLoad, _ := strconv.Atoi(row["cpu-load"])
	freeMem, _ := strconv.ParseInt(row["free-memory"], 10, 64)
	totalMem, _ := strconv.ParseInt(row["total-memory"], 10, 64)
	freeHDD, _ := strconv.ParseInt(row["free-hdd-space"], 10, 64)
	totalHDD, _ := strconv.ParseInt(row["total-hdd-space"], 10, 64)

	var memPct float64
	if totalMem > 0 {
		memPct = float64(totalMem-freeMem) / float64(totalMem) * 100.0
	}
	var hddPct float64
	if totalHDD > 0 {
		hddPct = float64(totalHDD-freeHDD) / float64(totalHDD) * 100.0
	}

	return &SystemMetrics{
		Uptime:        row["uptime"],
		Version:       row["version"],
		CPULoad:       cpuLoad,
		FreeMemory:    freeMem,
		TotalMemory:   totalMem,
		MemoryUsedPct: memPct,
		FreeHDD:       freeHDD,
		TotalHDD:      totalHDD,
		HDDUsedPct:    hddPct,
		BoardName:     row["board-name"],
		Architecture:  row["architecture-name"],
		Platform:      row["platform"],
		CPU:           row["cpu"],
	}
}

// --- HTTP SERVER & HANDLERS ---

type Server struct {
	hub *RouterHub
	mux *http.ServeMux
}

func NewServer(hub *RouterHub) *Server {
	s := &Server{
		hub: hub,
		mux: http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// CORS Headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/", s.handleIndex)
	s.mux.HandleFunc("/api/routers", s.handleRouters)
	s.mux.HandleFunc("/api/routers/", s.handleRouterByID)
	s.mux.HandleFunc("/api/system-resource", s.handleSystemResource)
	s.mux.HandleFunc("/api/catalog", s.handleCatalog)
	s.mux.HandleFunc("/api/catalog/build", s.handleCatalogBuild)
	s.mux.HandleFunc("/api/events/health", s.handleHealthSSE)
	s.mux.HandleFunc("/api/events/stream", s.handleStreamSSE)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data, err := indexHTML.ReadFile("index.html")
	if err != nil {
		http.Error(w, "index.html tidak ditemukan: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (s *Server) handleRouters(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		routers := s.hub.GetAllRouters()
		writeJSON(w, http.StatusOK, routers)

	case http.MethodPost:
		var req struct {
			ID          string  `json:"id"`
			Address     string  `json:"address"`
			Port        int     `json:"port"`
			Username    string  `json:"username"`
			Password    string  `json:"password"`
			TLS         bool    `json:"tls"`
			Insecure    bool    `json:"insecure"`
			Poll        bool    `json:"poll"`
			StreamEvery float64 `json:"stream_every_s"`
			PollEvery   float64 `json:"poll_every_s"`
			PingTarget  string  `json:"ping_target"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON payload tidak valid: " + err.Error()})
			return
		}

		if req.Port <= 0 {
			if req.TLS {
				req.Port = 8729
			} else {
				req.Port = 8728
			}
		}
		if req.PingTarget == "" {
			req.PingTarget = "8.8.8.8"
		}
		streamEvery := 1 * time.Second
		if req.StreamEvery > 0 {
			streamEvery = time.Duration(req.StreamEvery * float64(time.Second))
		}
		pollEvery := 5 * time.Second
		if req.PollEvery > 0 {
			pollEvery = time.Duration(req.PollEvery * float64(time.Second))
		}

		cfg := rosrun.Config{
			ID:          req.ID,
			Address:     req.Address,
			Port:        req.Port,
			Username:    req.Username,
			Password:    req.Password,
			TLS:         req.TLS,
			Insecure:    req.Insecure,
			Timeout:     10 * time.Second,
			Poll:        req.Poll,
			StreamEvery: streamEvery,
			PollEvery:   pollEvery,
			Status: rosrun.StatusConfig{
				Disable: false,
				Target:  req.PingTarget,
				Every:   1 * time.Second,
				Timeout: 15 * time.Second,
			},
		}

		st, err := s.hub.AddRouter(cfg)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, st)

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Metode tidak didukung"})
	}
}

func (s *Server) handleRouterByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/routers/")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID router diperlukan"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		st, ok := s.hub.GetRouter(id)
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Router tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusOK, st)

	case http.MethodPut:
		var req struct {
			Address     string  `json:"address"`
			Port        int     `json:"port"`
			Username    string  `json:"username"`
			Password    string  `json:"password"`
			TLS         bool    `json:"tls"`
			Insecure    bool    `json:"insecure"`
			Poll        bool    `json:"poll"`
			StreamEvery float64 `json:"stream_every_s"`
			PollEvery   float64 `json:"poll_every_s"`
			PingTarget  string  `json:"ping_target"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON payload tidak valid: " + err.Error()})
			return
		}

		old, ok := s.hub.GetRouter(id)
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Router tidak ditemukan"})
			return
		}

		s.hub.mu.RLock()
		rawCfg, hasRaw := s.hub.rawCfgs[id]
		s.hub.mu.RUnlock()

		cfg := rawCfg
		if !hasRaw {
			cfg = rosrun.Config{
				ID:       id,
				Address:  old.Config.Address,
				Port:     old.Config.Port,
				Username: old.Config.Username,
				Password: old.Config.Password,
				Timeout:  10 * time.Second,
			}
		}
		if req.Address != "" {
			cfg.Address = req.Address
		}
		if req.Port > 0 {
			cfg.Port = req.Port
		}
		if req.Username != "" {
			cfg.Username = req.Username
		}
		if req.Password != "" {
			cfg.Password = req.Password
		}
		cfg.TLS = req.TLS
		cfg.Insecure = req.Insecure
		cfg.Poll = req.Poll
		if req.StreamEvery > 0 {
			cfg.StreamEvery = time.Duration(req.StreamEvery * float64(time.Second))
		}
		if req.PollEvery > 0 {
			cfg.PollEvery = time.Duration(req.PollEvery * float64(time.Second))
		}
		if req.PingTarget != "" {
			cfg.Status.Target = req.PingTarget
		}

		st, err := s.hub.UpdateRouter(id, cfg)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, st)

	case http.MethodDelete:
		if err := s.hub.DeleteRouter(id); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted", "id": id})

	default:
		http.Error(w, "Metode tidak didukung", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleSystemResource(w http.ResponseWriter, r *http.Request) {
	routerID := r.URL.Query().Get("router")
	if routerID == "" {
		http.Error(w, "Query parameter 'router' diperlukan", http.StatusBadRequest)
		return
	}

	sess, ok := s.hub.mgr.Get(routerID)
	if !ok || sess == nil {
		http.Error(w, "Router tidak aktif atau tidak ditemukan", http.StatusNotFound)
		return
	}

	res, ok := ros.Lookup("system-resource")
	if !ok {
		http.Error(w, "Resource 'system-resource' tidak ditemukan", http.StatusInternalServerError)
		return
	}

	req, err := res.Build("list", ros.Once(), ros.Params{})
	if err != nil {
		http.Error(w, "Gagal build request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rows, err := sess.Once(ctx, req)
	if err != nil {
		http.Error(w, "Gagal mengambil system-resource: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(rows) == 0 {
		http.Error(w, "Tidak ada data resource dari router", http.StatusNoContent)
		return
	}

	metrics := parseSystemMetrics(rows[0])
	writeJSON(w, http.StatusOK, metrics)
}

func (s *Server) handleCatalog(w http.ResponseWriter, r *http.Request) {
	all := ros.Builtin().All()

	type ViewSummary struct {
		Name string `json:"name"`
		Note string `json:"note"`
	}

	type ResourceSummary struct {
		ID    string        `json:"id"`
		Path  string        `json:"path"`
		Title string        `json:"title"`
		Group string        `json:"group"`
		Views []ViewSummary `json:"views"`
		Notes string        `json:"notes"`
	}

	resList := make([]ResourceSummary, 0, len(all))
	for _, res := range all {
		views := make([]ViewSummary, 0, len(res.Views))
		for _, v := range res.Views {
			views = append(views, ViewSummary{
				Name: v.Name,
				Note: v.Note,
			})
		}
		resList = append(resList, ResourceSummary{
			ID:    res.ID,
			Path:  res.Path,
			Title: res.Title,
			Group: string(res.Group),
			Views: views,
			Notes: res.Notes,
		})
	}

	writeJSON(w, http.StatusOK, resList)
}

// handleCatalogBuild menyusun command MikroTik dari resource & view katalog tanpa mengeksekusinya.
func (s *Server) handleCatalogBuild(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Metode tidak didukung"})
		return
	}

	q := r.URL.Query()
	resID := q.Get("resource")
	view := q.Get("view")
	delivery := q.Get("delivery") // "stream_interval", "stream_follow", "poll", "once", "auto"
	intervalStr := q.Get("interval")
	proplistStr := q.Get("proplist")
	routerID := q.Get("router")

	if resID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Parameter 'resource' wajib diisi"})
		return
	}

	res, ok := ros.Lookup(resID)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Resource %q tidak ditemukan di katalog", resID)})
		return
	}

	if view == "" {
		if len(res.Views) > 0 {
			view = res.Views[0].Name
		} else {
			view = "list"
		}
	}

	var interval time.Duration = 1 * time.Second
	if intervalStr != "" {
		if d, err := time.ParseDuration(intervalStr); err == nil && d > 0 {
			interval = d
		}
	}

	var proplist []string
	if proplistStr != "" {
		for _, p := range strings.Split(proplistStr, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				proplist = append(proplist, p)
			}
		}
	}

	params := ros.Params{Proplist: proplist}

	var plan ros.Plan
	var tag string

	if delivery == "auto" && routerID != "" {
		if sess, ok := s.hub.mgr.Get(routerID); ok && sess != nil {
			plan = sess.PlanFor(res, view)
		} else {
			s.hub.mu.RLock()
			rawCfg, hasRaw := s.hub.rawCfgs[routerID]
			s.hub.mu.RUnlock()
			if hasRaw && rawCfg.Poll {
				plan = ros.Poll(ros.MinPollEvery)
			} else {
				plan = ros.StreamEvery(interval)
			}
		}
	} else {
		switch delivery {
		case "stream_follow":
			plan = ros.Stream()
		case "poll":
			pollInterval := interval
			if pollInterval < ros.MinPollEvery {
				pollInterval = ros.MinPollEvery
			}
			plan = ros.Poll(pollInterval)
		case "once":
			plan = ros.Once()
		default: // "stream_interval"
			plan = ros.StreamEvery(interval)
		}
	}

	req, err := res.Build(view, plan, params)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Build failed: " + err.Error()})
		return
	}

	if req.Delivery == ros.DeliveryPoll || req.Delivery == ros.DeliveryOnce {
		tag = "POLL"
	} else {
		tag = "STREAM"
	}

	sentenceStr := strings.Join(req.Sentence, " ")
	formatted := fmt.Sprintf("[%s] %s", tag, sentenceStr)

	writeJSON(w, http.StatusOK, map[string]any{
		"tag":       tag,
		"sentence":  req.Sentence,
		"command":   sentenceStr,
		"formatted": formatted,
		"resource":  resID,
		"view":      view,
		"delivery":  string(req.Delivery),
	})
}

func (s *Server) handleHealthSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming tidak didukung pada browser ini", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := make(chan HealthEvent, 16)
	s.hub.healthMu.Lock()
	s.hub.healthSubs[ch] = struct{}{}
	s.hub.healthMu.Unlock()

	defer func() {
		s.hub.healthMu.Lock()
		delete(s.hub.healthSubs, ch)
		s.hub.healthMu.Unlock()
	}()

	// Kirim inisial state segera
	for _, rState := range s.hub.GetAllRouters() {
		ev := HealthEvent{
			RouterID:     rState.Config.ID,
			State:        rState.State,
			Connected:    rState.Connected,
			Warnings:     rState.Warnings,
			LastError:    rState.LastError,
			PingTarget:   rState.PingTarget,
			PingLastRTT:  rState.PingLastRTT,
			PingAvgRTT:   rState.PingAvgRTT,
			PingMinRTT:   rState.PingMinRTT,
			PingMaxRTT:   rState.PingMaxRTT,
			PacketLoss:   rState.PacketLoss,
			PingSent:     rState.PingSent,
			PingReceived: rState.PingReceived,
			PingStatus:   rState.PingStatus,
			System:       rState.System,
			Timestamp:    time.Now(),
		}
		b, _ := json.Marshal(ev)
		fmt.Fprintf(w, "event: health\ndata: %s\n\n", b)
	}
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case ev := <-ch:
			b, err := json.Marshal(ev)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "event: health\ndata: %s\n\n", b)
			flusher.Flush()
		}
	}
}

// handleStreamSSE mengeksekusi tes live Stream / Poll / Once pada resource katalog
func (s *Server) handleStreamSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming tidak didukung pada browser ini", http.StatusBadRequest)
		return
	}

	q := r.URL.Query()
	routerID := q.Get("router")
	resID := q.Get("resource")
	view := q.Get("view")
	delivery := q.Get("delivery") // "stream_follow", "stream_interval", "poll", "once"
	intervalStr := q.Get("interval")
	proplistStr := q.Get("proplist")

	if routerID == "" || resID == "" {
		http.Error(w, "Parameter 'router' dan 'resource' wajib diisi", http.StatusBadRequest)
		return
	}

	sess, ok := s.hub.mgr.Get(routerID)
	if !ok || sess == nil {
		http.Error(w, fmt.Sprintf("Router %q tidak aktif atau tidak ditemukan", routerID), http.StatusNotFound)
		return
	}

	res, ok := ros.Lookup(resID)
	if !ok {
		http.Error(w, fmt.Sprintf("Resource %q tidak ditemukan di katalog", resID), http.StatusNotFound)
		return
	}

	if view == "" {
		if len(res.Views) > 0 {
			view = res.Views[0].Name
		} else {
			view = "list"
		}
	}

	var interval time.Duration = 1 * time.Second
	if intervalStr != "" {
		if d, err := time.ParseDuration(intervalStr); err == nil && d > 0 {
			interval = d
		}
	}

	var proplist []string
	if proplistStr != "" {
		for _, p := range strings.Split(proplistStr, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				proplist = append(proplist, p)
			}
		}
	}

	params := ros.Params{Proplist: proplist}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	sendSSE := func(eventType string, data any) {
		b, err := json.Marshal(data)
		if err != nil {
			return
		}
		fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, b)
		flusher.Flush()
	}

	// 1. One-shot / Once print
	if delivery == "once" {
		req, err := res.Build(view, ros.Once(), params)
		if err != nil {
			sendSSE("error", map[string]string{"message": "Build failed: " + err.Error()})
			return
		}
		sendSSE("status", map[string]string{"message": fmt.Sprintf("Running Once print on %s (%s)...", res.Path, view)})

		rows, err := sess.Once(r.Context(), req)
		if err != nil {
			sendSSE("error", map[string]string{"message": err.Error()})
			return
		}
		sendSSE("snapshot", map[string]any{
			"count": len(rows),
			"rows":  rows,
		})
		sendSSE("done", map[string]string{"message": "Completed snapshot"})
		return
	}

	// 2. Stream (Interval / Follow) atau Poll (Ticker)
	var plan ros.Plan
	switch delivery {
	case "stream_follow":
		plan = ros.Stream()
	case "poll":
		pollInterval := interval
		if pollInterval < ros.MinPollEvery {
			pollInterval = ros.MinPollEvery
		}
		plan = ros.Poll(pollInterval)
	default: // "stream_interval"
		plan = ros.StreamEvery(interval)
	}

	req, err := res.Build(view, plan, params)
	if err != nil {
		sendSSE("error", map[string]string{"message": "Build failed: " + err.Error()})
		return
	}

	sendSSE("status", map[string]string{
		"message":  fmt.Sprintf("Started %s on %s (%s)", delivery, res.Path, view),
		"sentence": strings.Join(req.Sentence, " "),
	})

	evCh, err := sess.Watch(r.Context(), req)
	if err != nil {
		sendSSE("error", map[string]string{"message": "Watch failed: " + err.Error()})
		return
	}

	for {
		select {
		case <-r.Context().Done():
			// Browser menutup koneksi atau klik stop -> cancel dikirim otomatis ke router
			return
		case ev, ok := <-evCh:
			if !ok {
				sendSSE("done", map[string]string{"message": "Stream selesai dari sisi router"})
				return
			}

			switch ev.Kind {
			case rosrun.KindRow, rosrun.KindAdded, rosrun.KindChanged, rosrun.KindRemoved:
				sendSSE("row", map[string]any{
					"kind": string(ev.Kind),
					"row":  ev.Row,
					"time": time.Now().Format("15:04:05.000"),
				})
			case rosrun.KindSnapshot:
				sendSSE("snapshot", map[string]any{
					"rows":  ev.Rows,
					"count": len(ev.Rows),
					"time":  time.Now().Format("15:04:05.000"),
				})
			case rosrun.KindCount:
				sendSSE("row", map[string]any{
					"kind":  "count",
					"count": ev.Count,
					"time":  time.Now().Format("15:04:05.000"),
				})
			case rosrun.KindError:
				errMsg := "stream error"
				if ev.Err != nil {
					errMsg = ev.Err.Error()
				}
				sendSSE("error", map[string]any{
					"error": errMsg,
					"time":  time.Now().Format("15:04:05.000"),
				})
			case rosrun.KindReconnecting:
				sendSSE("status", map[string]string{
					"message": "Reconnecting ke router...",
				})
			}
		}
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func main() {
	port := flag.Int("port", 8080, "Port HTTP server dashboard")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hub := NewRouterHub(ctx)
	defer hub.Close()

	// Seed router dari cmd.MikrotikList (data.go) secara asinkron
	go func() {
		log.Printf("[WebTest] Menyemai %d router dari data.go...", len(cmd.MikrotikList))
		for _, cfg := range cmd.MikrotikList {
			st, err := hub.AddRouter(cfg)
			if err != nil {
				log.Printf("[WebTest] Router %s (%s:%d): %v", cfg.ID, cfg.Address, cfg.Port, err)
			} else {
				log.Printf("[WebTest] Router %s (%s:%d) terdaftar, state: %s", cfg.ID, cfg.Address, cfg.Port, st.State)
			}
		}
	}()

	server := NewServer(hub)
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("[WebTest] Dashboard roscatalog berjalan di http://localhost%s", addr)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: server,
	}

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server HTTP berhenti dengan error: %v", err)
	}
}
