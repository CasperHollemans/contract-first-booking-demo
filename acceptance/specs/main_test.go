package specs

import (
	"os"
	"testing"

	"github.com/casperhollemans/contract-first-booking-demo/acceptance/drivers"
)

func TestMain(m *testing.M) {
	code := m.Run()
	drivers.CloseBrowser()
	os.Exit(code)
}
