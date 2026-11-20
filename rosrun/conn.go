package rosrun

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/quiqxiq/goros/v4"
)

// managedConn adalah satu *routeros.Client yang bisa diganti saat reconnect.
// Semua pembacaan client/runner/errC lewat method agar tidak pernah melihat
// pointer setengah jadi.
type managedConn struct {
	name string

	mu     sync.RWMutex
	client *routeros.Client
	errC   <-chan error
	runner *Runner

	down atomic.Bool

	errMu   sync.Mutex
	lastErr error
}

func newManagedConn(name string, client *routeros.Client, errC <-chan error, run *Runner) *managedConn {
	c := &managedConn{name: name}
	c.set(client, errC, run)
	return c
}

func (c *managedConn) get() *routeros.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.client
}

func (c *managedConn) run() *Runner {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.runner
}

func (c *managedConn) errChan() <-chan error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.errC
}

func (c *managedConn) set(client *routeros.Client, errC <-chan error, run *Runner) {
	c.mu.Lock()
	c.client, c.errC, c.runner = client, errC, run
	c.mu.Unlock()
}

func (c *managedConn) markDown(err error) {
	c.down.Store(true)
	if err != nil {
		c.errMu.Lock()
		c.lastErr = err
		c.errMu.Unlock()
	}
}

func (c *managedConn) markUp() {
	c.down.Store(false)
	c.errMu.Lock()
	c.lastErr = nil
	c.errMu.Unlock()
}

func (c *managedConn) isDown() bool { return c.down.Load() }

func (c *managedConn) err() error {
	c.errMu.Lock()
	defer c.errMu.Unlock()
	return c.lastErr
}

// dial membuka satu koneksi login+async. Setiap koneksi punya loop async
// sendiri; jangan pernah memberi ctx berumur pendek ke AsyncContext karena
// pembatalannya menutup pembaca bersama untuk seluruh koneksi.
func dial(ctx context.Context, cfg Config, name string) (*managedConn, error) {
	dialCtx := ctx
	if cfg.Timeout > 0 {
		var cancel context.CancelFunc
		dialCtx, cancel = context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()
	}

	address := cfg.addr()

	var client *routeros.Client
	var err error
	if cfg.TLS {
		client, err = routeros.DialTLSContext(dialCtx, address, cfg.Username, cfg.Password, &tls.Config{
			InsecureSkipVerify: cfg.Insecure,
		})
	} else {
		client, err = routeros.DialContext(dialCtx, address, cfg.Username, cfg.Password)
	}
	if err != nil {
		return nil, fmt.Errorf("rosrun: dial %s (%s): %w", address, name, err)
	}

	// Queue dibaca Listen saat mengalokasikan channel. Tanpanya, satu stream
	// yang konsumennya lambat memblokir dispatcher bersama.
	client.Queue = cfg.Queue

	errC := client.AsyncContext(context.Background())
	conn := newManagedConn(name, client, errC, newRunner(client, cfg))
	return conn, nil
}

// dialRetry mencoba login beberapa kali untuk gangguan transien. RouterOS 6
// kadang menutup koneksi dengan EOF bila sesi API dibuat beruntun.
func dialRetry(ctx context.Context, cfg Config, name string) (*managedConn, error) {
	var last error
	for attempt := 0; attempt < 3; attempt++ {
		conn, err := dial(ctx, cfg, name)
		if err == nil {
			return conn, nil
		}
		last = err
		if !isTransientDialErr(err) || attempt == 2 {
			break
		}
		wait := time.Duration(150*(1<<attempt)) * time.Millisecond
		if !sleepContext(ctx, wait) {
			break
		}
	}
	return nil, last
}

func isTransientDialErr(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
		return true
	}
	var ne net.Error
	if errors.As(err, &ne) && ne.Timeout() {
		return true
	}
	msg := strings.ToLower(err.Error())
	for _, s := range []string{"connection reset", "broken pipe", "connection refused", "eof", "i/o timeout"} {
		if strings.Contains(msg, s) {
			return true
		}
	}
	return false
}

// supervise menunggu koneksi mati lalu menggantinya. Satu supervisor per
// koneksi (cmd dan feed).
func (m *MikroTik) supervise(conn *managedConn) {
	defer m.wg.Done()

	for {
		select {
		case <-m.stopCh:
			return
		case err, ok := <-conn.errChan():
			if !ok {
				// Loop async berhenti tanpa error eksplisit (mis. Close).
				select {
				case <-m.stopCh:
					return
				default:
				}
			}
			conn.markDown(err)
			m.publishStatus()
		}

		if !m.reconnect(conn) {
			return
		}
	}
}

// reconnect melakukan dial ulang dengan backoff sampai berhasil atau sesi
// ditutup. Saat berhasil, runner dibangun ulang di atas client baru dan epoch
// dinaikkan supaya Watch tahu ada jeda data.
func (m *MikroTik) reconnect(conn *managedConn) bool {
	base := m.cfg.backoff
	if base <= 0 {
		base = time.Second
	}
	const maxBackoff = 10 * time.Second
	backoff := base

	for {
		select {
		case <-m.stopCh:
			return false
		default:
		}

		fresh, err := dial(m.runCtx, m.cfg, conn.name)
		if err != nil {
			conn.markDown(err)
			m.publishStatus()
			if m.cfg.Logf != nil {
				m.cfg.Logf("rosrun: %s reconnect failed: %v (retry in %s)", conn.name, err, backoff)
			}
			wait := backoff
			if wait > 20*time.Millisecond {
				wait += time.Duration(rand.Int64N(int64(wait / 4)))
			}
			select {
			case <-m.stopCh:
				return false
			case <-time.After(wait):
			}
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}

		conn.set(fresh.get(), fresh.errChan(), fresh.run())
		conn.markUp()
		m.reconnects.Add(1)
		m.epochs.bump()
		m.publishStatus()
		if m.cfg.Logf != nil {
			m.cfg.Logf("rosrun: %s reconnected (epoch %d)", conn.name, m.epochs.current())
		}
		return true
	}
}
