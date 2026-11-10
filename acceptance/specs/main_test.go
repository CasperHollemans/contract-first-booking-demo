package specs

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/casperhollemans/contract-first-booking-demo/acceptance/drivers"
	"github.com/casperhollemans/contract-first-booking-demo/acceptance/dsl"
)

func TestMain(m *testing.M) {
	code := m.Run()
	drivers.CloseBrowser()

	// Story guard: every story proves at least one spec through the UI.
	if missing := dsl.StoriesWithoutUI(); len(missing) > 0 {
		fmt.Fprintf(os.Stderr, "FAIL story guard: %s ran specs but none through the UI: add UI to its happy-path spec\n",
			strings.Join(missing, ", "))
		code = 1
	}
	os.Exit(code)
}
