package drivers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

// APIDriver reaches the system through its JSON HTTP API.
type APIDriver struct {
	*system
}

func NewAPIDriver(t *testing.T) *APIDriver {
	return &APIDriver{system: startSystem(t)}
}

// CheckSlot calls GET /slots/check?date=&start=&end=
// and expects 200 {"selectable": true|false}.
func (d *APIDriver) CheckSlot(date, start, end string) (bool, error) {
	q := url.Values{"date": {date}, "start": {start}, "end": {end}}
	resp, err := http.Get(d.baseURL + "/slots/check?" + q.Encode())
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("GET /slots/check: expected 200, got %d", resp.StatusCode)
	}
	var body struct {
		Selectable *bool `json:"selectable"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return false, fmt.Errorf("GET /slots/check: invalid JSON: %w", err)
	}
	if body.Selectable == nil {
		return false, fmt.Errorf(`GET /slots/check: response has no "selectable" field`)
	}
	return *body.Selectable, nil
}
