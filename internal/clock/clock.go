// Package clock is the port for "now", so tests can control time.
package clock

import "time"

type Clock interface {
	Now() time.Time
}
