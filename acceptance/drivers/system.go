package drivers

import (
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/casperhollemans/contract-first-booking-demo/internal/app"
)

// fakeClock is the only fake in the acceptance tests: everything else is the
// real system, reached over a real port.
type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *fakeClock) set(now time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = now
}

// system is one running instance of the real server, isolated per test.
type system struct {
	clock   *fakeClock
	baseURL string
}

func startSystem(t *testing.T) *system {
	t.Helper()
	clock := &fakeClock{now: time.Now()}
	srv := httptest.NewServer(app.New(clock))
	t.Cleanup(srv.Close)
	return &system{clock: clock, baseURL: srv.URL}
}

func (s *system) SetNow(now time.Time) { s.clock.set(now) }
