package drivers

import (
	"fmt"
	"time"

	"github.com/casperhollemans/contract-first-booking-demo/internal/domain"
	"github.com/casperhollemans/contract-first-booking-demo/internal/service"
)

// domainCafeTime is the cafe's time zone, as wired in internal/app (UTC).
var domainCafeTime = time.UTC

// DomainDriver reaches the system through the service layer: the same use
// cases the server runs, with the fake clock, but no HTTP. Use it for rules
// and logic. (When a story stores data, it gets an in-memory repository.)
type DomainDriver struct {
	clock *fakeClock
	slots *service.SlotService
}

func NewDomainDriver() *DomainDriver {
	clock := &fakeClock{now: time.Now()}
	return &DomainDriver{clock: clock, slots: service.NewSlotService(clock)}
}

func (d *DomainDriver) SetNow(now time.Time) { d.clock.set(now) }

// CheckSlot reads date "2006-01-02" and times "15:04" in the cafe's time zone
// and asks the slot service directly.
func (d *DomainDriver) CheckSlot(date, start, end string) (bool, error) {
	s, err := time.ParseInLocation("2006-01-02 15:04", date+" "+start, domainCafeTime)
	if err != nil {
		return false, fmt.Errorf("Domain: cannot read start %q %q: %w", date, start, err)
	}
	e, err := time.ParseInLocation("2006-01-02 15:04", date+" "+end, domainCafeTime)
	if err != nil {
		return false, fmt.Errorf("Domain: cannot read end %q %q: %w", date, end, err)
	}
	return d.slots.IsSelectable(domain.Slot{Start: s, End: e}), nil
}
