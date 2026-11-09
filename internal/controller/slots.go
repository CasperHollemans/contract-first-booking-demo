// Package controller is the HTTP interface: it turns requests into service
// calls and service answers into JSON or HTML.
package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/casperhollemans/contract-first-booking-demo/internal/domain"
	"github.com/casperhollemans/contract-first-booking-demo/internal/service"
)

type SlotController struct {
	slots *service.SlotService
	loc   *time.Location // the cafe's time zone
}

func NewSlotController(slots *service.SlotService, cafeTime *time.Location) *SlotController {
	return &SlotController{slots: slots, loc: cafeTime}
}

func (c *SlotController) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /slots/check", c.checkJSON)
}

// checkJSON: GET /slots/check?date=2026-12-09&start=10:00&end=11:00
// → 200 {"selectable": bool}, or 400 if the input can't be read.
func (c *SlotController) checkJSON(w http.ResponseWriter, r *http.Request) {
	slot, err := parseSlot(r.URL.Query(), c.loc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Selectable bool `json:"selectable"`
	}{c.slots.IsSelectable(slot)})
}

var errBadSlot = errors.New(`expected date=YYYY-MM-DD, start=HH:MM and end=HH:MM`)

// parseSlot reads date, start and end. It only checks the format: whether the
// slot follows the booking rules is up to the domain.
func parseSlot(q url.Values, loc *time.Location) (domain.Slot, error) {
	start, err := time.ParseInLocation("2006-01-02 15:04", q.Get("date")+" "+q.Get("start"), loc)
	if err != nil {
		return domain.Slot{}, errBadSlot
	}
	end, err := time.ParseInLocation("2006-01-02 15:04", q.Get("date")+" "+q.Get("end"), loc)
	if err != nil {
		return domain.Slot{}, errBadSlot
	}
	return domain.Slot{Start: start, End: end}, nil
}
