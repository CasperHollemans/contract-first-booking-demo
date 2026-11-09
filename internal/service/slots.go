// Package service holds the use cases. It applies domain rules using ports
// such as the clock (and, later, repositories).
package service

import (
	"github.com/casperhollemans/contract-first-booking-demo/internal/clock"
	"github.com/casperhollemans/contract-first-booking-demo/internal/domain"
)

type SlotService struct {
	clock clock.Clock
}

func NewSlotService(c clock.Clock) *SlotService {
	return &SlotService{clock: c}
}

// IsSelectable reports whether the slot can be selected right now.
func (s *SlotService) IsSelectable(slot domain.Slot) bool {
	return slot.SelectableAt(s.clock.Now())
}
