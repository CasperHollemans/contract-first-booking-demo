// Package app is the composition root: it wires controllers, services and
// repositories into one http.Handler.
package app

import (
	"net/http"
	"time"

	"github.com/casperhollemans/contract-first-booking-demo/internal/clock"
	"github.com/casperhollemans/contract-first-booking-demo/internal/controller"
	"github.com/casperhollemans/contract-first-booking-demo/internal/service"
)

// cafeTime is the cafe's time zone. The stories don't name one; the demo uses UTC.
var cafeTime = time.UTC

// New builds the whole application.
func New(c clock.Clock) http.Handler {
	slots := service.NewSlotService(c)

	mux := http.NewServeMux()
	controller.NewSlotController(slots, cafeTime).Register(mux)
	return mux
}
