// Command server runs the booking app with the real clock.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/casperhollemans/contract-first-booking-demo/internal/app"
	"github.com/casperhollemans/contract-first-booking-demo/internal/clock"
)

func main() {
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}
	log.Printf("booking demo listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, app.New(clock.System{})))
}
