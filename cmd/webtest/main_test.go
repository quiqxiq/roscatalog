package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/quiqxiq/roscatalog/rosrun"
)

func TestParseSystemMetrics(t *testing.T) {
	sampleRow := map[string]string{
		"uptime":            "6d13:41:22",
		"version":           "7.16.1 (stable)",
		"cpu-load":          "4",
		"free-memory":       "8870129664",
		"total-memory":      "8946761728",
		"free-hdd-space":    "12345678",
		"total-hdd-space":   "67890123",
		"board-name":        "CHR",
		"architecture-name": "x86_64",
		"platform":          "MikroTik",
		"cpu":               "QEMU Virtual CPU",
	}

	metrics := parseSystemMetrics(sampleRow)
	if metrics == nil {
		t.Fatal("metrics nil")
	}

	if metrics.Uptime != "6d13:41:22" {
		t.Errorf("Uptime = %s, mau 6d13:41:22", metrics.Uptime)
	}
	if metrics.Version != "7.16.1 (stable)" {
		t.Errorf("Version = %s, mau 7.16.1 (stable)", metrics.Version)
	}
	if metrics.CPULoad != 4 {
		t.Errorf("CPULoad = %d, mau 4", metrics.CPULoad)
	}
	if metrics.TotalMemory != 8946761728 {
		t.Errorf("TotalMemory = %d", metrics.TotalMemory)
	}
	if metrics.MemoryUsedPct <= 0 || metrics.MemoryUsedPct >= 100 {
		t.Errorf("MemoryUsedPct tidak valid: %f", metrics.MemoryUsedPct)
	}
	if metrics.BoardName != "CHR" {
		t.Errorf("BoardName = %s, mau CHR", metrics.BoardName)
	}
}

func TestServerEndpoints(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hub := NewRouterHub(ctx)
	defer hub.Close()
	server := NewServer(hub)

	// 1. GET / (index.html)
	t.Run("GET /", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("GET / returned status %d, expected 200", w.Code)
		}
		body := w.Body.String()
		if !strings.Contains(body, "ROSCATALOG MULTI-ROUTER TESTER") {
			t.Errorf("index.html does not contain expected brand title")
		}
	})

	// 2. GET /api/catalog
	t.Run("GET /api/catalog", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/catalog", nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("GET /api/catalog returned status %d", w.Code)
		}
		var list []map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}
		if len(list) == 0 {
			t.Fatalf("catalog list empty")
		}
		foundHotspot := false
		for _, item := range list {
			if item["id"] == "hotspot-active" {
				foundHotspot = true
				break
			}
		}
		if !foundHotspot {
			t.Errorf("catalog does not contain 'hotspot-active'")
		}
	})

	// 2.b GET /api/catalog/build
	t.Run("GET /api/catalog/build", func(t *testing.T) {
		// A. Stream build (menggunakan view 'live' yang mendukung mode interval)
		reqStream := httptest.NewRequest(http.MethodGet, "/api/catalog/build?resource=hotspot-active&view=live&delivery=stream_interval&proplist=user,address", nil)
		wStream := httptest.NewRecorder()
		server.ServeHTTP(wStream, reqStream)
		if wStream.Code != http.StatusOK {
			t.Fatalf("GET /api/catalog/build (stream) returned %d: %s", wStream.Code, wStream.Body.String())
		}
		var resStream map[string]any
		_ = json.Unmarshal(wStream.Body.Bytes(), &resStream)
		if resStream["tag"] != "STREAM" {
			t.Errorf("expected tag STREAM, got %v", resStream["tag"])
		}
		cmdStr, _ := resStream["command"].(string)
		if !strings.Contains(cmdStr, "/ip/hotspot/active/print") || !strings.Contains(cmdStr, "=interval=") || !strings.Contains(cmdStr, "user,address") {
			t.Errorf("unexpected command: %s", cmdStr)
		}

		// B. Poll build (menggunakan view 'list')
		reqPoll := httptest.NewRequest(http.MethodGet, "/api/catalog/build?resource=hotspot-active&view=list&delivery=poll&proplist=user,address", nil)
		wPoll := httptest.NewRecorder()
		server.ServeHTTP(wPoll, reqPoll)
		if wPoll.Code != http.StatusOK {
			t.Fatalf("GET /api/catalog/build (poll) returned %d: %s", wPoll.Code, wPoll.Body.String())
		}
		var resPoll map[string]any
		_ = json.Unmarshal(wPoll.Body.Bytes(), &resPoll)
		if resPoll["tag"] != "POLL" {
			t.Errorf("expected tag POLL, got %v", resPoll["tag"])
		}
		cmdPollStr, _ := resPoll["command"].(string)
		if !strings.Contains(cmdPollStr, "/ip/hotspot/active/print") || strings.Contains(cmdPollStr, "=interval=") || !strings.Contains(cmdPollStr, "user,address") {
			t.Errorf("poll command should not contain =interval=: %s", cmdPollStr)
		}
	})

	// 3. Router CRUD lifecycle
	t.Run("Router CRUD", func(t *testing.T) {
		// A. POST /api/routers (Add dummy router)
		addPayload := map[string]any{
			"id":             "test-router-unit-01",
			"address":        "127.0.0.1",
			"port":           9999, // Unused port, should register as error/offline
			"username":       "testuser",
			"password":       "testpass",
			"poll":           false,
			"stream_every_s": 1.0,
			"poll_every_s":   5.0,
			"ping_target":    "8.8.8.8",
			"tls":            false,
			"insecure":       false,
		}
		b, _ := json.Marshal(addPayload)
		req := httptest.NewRequest(http.MethodPost, "/api/routers", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("POST /api/routers returned status %d: %s", w.Code, w.Body.String())
		}

		var created RouterState
		if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
			t.Fatalf("Unmarshal created router error: %v", err)
		}
		if created.Config.ID != "test-router-unit-01" {
			t.Errorf("created ID = %s, expected test-router-unit-01", created.Config.ID)
		}

		// B. GET /api/routers
		reqList := httptest.NewRequest(http.MethodGet, "/api/routers", nil)
		wList := httptest.NewRecorder()
		server.ServeHTTP(wList, reqList)
		if wList.Code != http.StatusOK {
			t.Fatalf("GET /api/routers returned status %d", wList.Code)
		}
		var routers []*RouterState
		_ = json.Unmarshal(wList.Body.Bytes(), &routers)
		if len(routers) != 1 {
			t.Fatalf("expected 1 router, got %d", len(routers))
		}
		if routers[0].Config.Password != "" {
			t.Errorf("GET /api/routers leaked raw password")
		}
		if !routers[0].Config.HasPassword {
			t.Errorf("expected has_password=true")
		}

		// C. GET /api/routers/{id}
		reqGet := httptest.NewRequest(http.MethodGet, "/api/routers/test-router-unit-01", nil)
		wGet := httptest.NewRecorder()
		server.ServeHTTP(wGet, reqGet)
		if wGet.Code != http.StatusOK {
			t.Fatalf("GET /api/routers/{id} returned status %d", wGet.Code)
		}

		// D. PUT /api/routers/{id} (Update address & username)
		updatePayload := map[string]any{
			"address":  "127.0.0.2",
			"port":     9998,
			"username": "updated_user",
		}
		bUp, _ := json.Marshal(updatePayload)
		reqPut := httptest.NewRequest(http.MethodPut, "/api/routers/test-router-unit-01", bytes.NewReader(bUp))
		reqPut.Header.Set("Content-Type", "application/json")
		wPut := httptest.NewRecorder()
		server.ServeHTTP(wPut, reqPut)
		if wPut.Code != http.StatusOK {
			t.Fatalf("PUT /api/routers/{id} returned status %d", wPut.Code)
		}
		var updated RouterState
		_ = json.Unmarshal(wPut.Body.Bytes(), &updated)
		if updated.Config.Username != "updated_user" {
			t.Errorf("updated username = %s, expected updated_user", updated.Config.Username)
		}
		if updated.Config.Address != "127.0.0.2" {
			t.Errorf("updated address = %s, expected 127.0.0.2", updated.Config.Address)
		}

		// E. DELETE /api/routers/{id}
		reqDel := httptest.NewRequest(http.MethodDelete, "/api/routers/test-router-unit-01", nil)
		wDel := httptest.NewRecorder()
		server.ServeHTTP(wDel, reqDel)
		if wDel.Code != http.StatusOK {
			t.Fatalf("DELETE /api/routers/{id} returned status %d", wDel.Code)
		}

		// Verify deletion
		reqGetAfter := httptest.NewRequest(http.MethodGet, "/api/routers/test-router-unit-01", nil)
		wGetAfter := httptest.NewRecorder()
		server.ServeHTTP(wGetAfter, reqGetAfter)
		if wGetAfter.Code != http.StatusNotFound {
			t.Errorf("expected 404 after deletion, got %d", wGetAfter.Code)
		}
	})
}

func TestHealthSSEEndpoint(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hub := NewRouterHub(ctx)
	defer hub.Close()
	server := NewServer(hub)

	// Add dummy router to verify initial SSE event
	_, _ = hub.AddRouter(rosrun.Config{
		ID:      "sse-router-01",
		Address: "127.0.0.1",
		Port:    9997,
	})

	sseCtx, sseCancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer sseCancel()

	req := httptest.NewRequest(http.MethodGet, "/api/events/health", nil).WithContext(sseCtx)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)
	body := w.Body.String()

	if !strings.Contains(body, "event: health") {
		t.Errorf("SSE health output does not contain 'event: health'")
	}
	if !strings.Contains(body, "sse-router-01") {
		t.Errorf("SSE health output does not contain router id 'sse-router-01'")
	}
}

func TestServerTokenAuth(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hub := NewRouterHub(ctx)
	defer hub.Close()
	server := NewServerWithToken(hub, "secret-tok")

	req := httptest.NewRequest(http.MethodGet, "/api/routers", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated GET /api/routers = %d, want 401", w.Code)
	}

	reqAuth := httptest.NewRequest(http.MethodGet, "/api/routers", nil)
	reqAuth.Header.Set("Authorization", "Bearer secret-tok")
	wAuth := httptest.NewRecorder()
	server.ServeHTTP(wAuth, reqAuth)
	if wAuth.Code != http.StatusOK {
		t.Fatalf("authenticated GET /api/routers = %d, want 200", wAuth.Code)
	}
}
