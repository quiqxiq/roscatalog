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
	cancel context.CancelFunc

	down atomic.Bool

	errMu   sync.Mutex
	lastErr error
}

func newManagedConn(name string, client *routeros.Client, errC <-chan error, run *Runner, cancel context.CancelFunc) *managedConn {
	c := &managedConn{name: name}
	c.set(client, errC, run, cancel)
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

func (c *managedConn) set(client *routeros.Client, errC <-chan error, run *Runner, cancel context.CancelFunc) (*routeros.Client, context.CancelFunc) {
	c.mu.Lock()
	oldClient, oldCancel := c.client, c.cancel
	c.client, c.errC, c.runner, c.cancel = client, errC, run, cancel
	c.mu.Unlock()
	return oldClient, oldCancel
}

func (c *managedConn) replace(fresh *managedConn) {
	if fresh == nil {
		return
	}
	fresh.mu.Lock()
	client, errC, runner, cancel := fresh.client, fresh.errC, fresh.runner, fresh.cancel
	fresh.client, fresh.errC, fresh.runner, fresh.cancel = nil, nil, nil, nil
	fresh.mu.Unlock()

	oldClient, oldCancel := c.set(client, errC, runner, cancel)
	if oldClient != nil {
		_ = oldClient.Close()
	}
	if oldCancel != nil {
		oldCancel()
	}
}

func (c *managedConn) close() error {
	c.mu.Lock()
	client, cancel := c.client, c.cancel
	c.client, c.errC, c.runner, c.cancel = nil, nil, nil, nil
	c.mu.Unlock()
	var err error
	if client != nil {
		err = client.Close()
	}
	if cancel != nil {
		cancel()
	}
	return err
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
func dial(ctx context.Context, cfg Config, name string, lim *limiter) (*managedConn, error) {
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

	asyncCtx, asyncCancel := context.WithCancel(context.Background())
	errC := client.AsyncContext(asyncCtx)
	conn := newManagedConn(name, client, errC, newRunner(client, cfg, lim), asyncCancel)
	return conn, nil
}

// dialRetry mencoba login beberapa kali untuk gangguan transien. RouterOS 6
// kadang menutup koneksi dengan EOF bila sesi API dibuat beruntun.
func dialRetry(ctx context.Context, cfg Config, name string, lim *limiter) (*managedConn, error) {
	var last error
	for attempt := 0; attempt < 3; attempt++ {
		conn, err := dial(ctx, cfg, name, lim)
		if err == nil {
			return conn, nil
		}
		last = err
		if isLoginError(err) || !isTransientDialErr(err) || attempt == 2 {
			// Kredensial salah tidak membaik dengan retry cepat; setiap
			// percobaan menulis login gagal ke log router.
			break
		}
		wait := time.Duration(150*(1<<attempt)) * time.Millisecond
		if !sleepContext(ctx, wait) {
			break
		}
	}
	return nil, last
}

// isLoginError melaporkan penolakan kredensial, yang tidak membaik dengan
// retry cepat dan justru membanjiri log login router.
func isLoginError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	for _, s := range []string{
		"invalid user name or password", "cannot log in", "login failed",
		"not authorized", "invalid username or password",
	} {
		if strings.Contains(msg, s) {
			return true
		}
	}
	return false
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

// supervise menunggu koneksi mati lalu menggantinya. Ada tepat satu
// supervisor untuk satu-satunya koneksi sesi.
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
			if m.limiter != nil {
				m.limiter.streamsGone()
			}
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

		fresh, err := dial(m.runCtx, m.cfg, conn.name, m.limiter)
		if err != nil {
			conn.markDown(err)
			if m.limiter != nil {
				m.limiter.streamsGone()
			}
			m.publishStatus()
			if m.cfg.Logf != nil {
				m.cfg.Logf("rosrun: %s reconnect failed: %v (retry in %s)", conn.name, err, backoff)
			}
			limit := maxBackoff
			if isLoginError(err) {
				// Auth backoff: jangan mencoba tiap detik dengan password salah.
				if backoff < 30*time.Second {
					backoff = 30 * time.Second
				}
				limit = 5 * time.Minute
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
			if backoff > limit {
				backoff = limit
			}
			continue
		}

		select {
		case <-m.stopCh:
			_ = fresh.close()
			return false
		default:
		}

		conn.replace(fresh)
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
