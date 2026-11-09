// Package app is the composition root: it wires controllers, services and
// repositories into one http.Handler.
package app

import (
	"net/http"

	"github.com/casperhollemans/contract-first-booking-demo/internal/clock"
)

// New builds the whole application. Not implemented yet: the acceptance
// specs define what it must do.
func New(c clock.Clock) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", http.StatusNotImplemented)
	})
}
