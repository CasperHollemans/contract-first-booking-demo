// Package clock is the port for "now", so tests can control time.
package clock

import "time"

type Clock interface {
	Now() time.Time
}

// System is the real clock, used in production.
type System struct{}

func (System) Now() time.Time { return time.Now() }
