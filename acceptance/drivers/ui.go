package drivers

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// UI contract (data-testid attributes on the page served at "/"):
//
//	date         <input type="date">   value "2006-01-02"
//	start, end   <input type="time">   value "15:04"
//	check        button that submits the form
//	slot-result  element that appears with data-selectable="true" | "false"
const uiTimeout = 5 * time.Second

var (
	browserOnce   sync.Once
	browserCtx    context.Context
	browserCancel context.CancelFunc
)

// browser starts one headless Chrome shared by all UI tests in the run.
func browser() context.Context {
	browserOnce.Do(func() {
		allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), chromedp.DefaultExecAllocatorOptions[:]...)
		ctx, cancel := chromedp.NewContext(allocCtx)
		browserCtx = ctx
		browserCancel = func() { cancel(); allocCancel() }
	})
	return browserCtx
}

// CloseBrowser shuts down the shared browser; call it from TestMain.
func CloseBrowser() {
	if browserCancel != nil {
		browserCancel()
	}
}

// UIDriver reaches the system through the web page in headless Chrome.
type UIDriver struct {
	*system
	tab context.Context
}

func NewUIDriver(t *testing.T) *UIDriver {
	tab, err := newTab(t)
	if err != nil {
		t.Fatalf("UI: cannot open a browser tab: %v", err)
	}
	return &UIDriver{system: startSystem(t), tab: tab}
}

// newTab opens a tab up front, so per-call timeouts don't close it.
func newTab(t *testing.T) (context.Context, error) {
	tab, cancel := chromedp.NewContext(browser())
	t.Cleanup(cancel)
	return tab, chromedp.Do(tab)
}

func (d *UIDriver) CheckSlot(date, start, end string) (bool, error) {
	ctx, cancel := context.WithTimeout(d.tab, uiTimeout)
	defer cancel()

	result := chromedp.CSS(`[data-testid="slot-result"][data-selectable]`)
	err := chromedp.Do(ctx,
		chromedp.Navigate(d.baseURL+"/"),
		chromedp.SetValue(chromedp.CSS(`[data-testid="date"]`), date),
		chromedp.SetValue(chromedp.CSS(`[data-testid="start"]`), start),
		chromedp.SetValue(chromedp.CSS(`[data-testid="end"]`), end),
		chromedp.Click(chromedp.CSS(`[data-testid="check"]`)),
		chromedp.WaitVisible(result),
	)
	if err != nil {
		return false, fmt.Errorf("UI: checking slot via the page: %w", err)
	}
	attr, err := chromedp.Run(ctx, chromedp.AttributeValue(result, "data-selectable"))
	if err != nil {
		return false, fmt.Errorf("UI: reading slot-result: %w", err)
	}
	selectable := attr.Value
	switch selectable {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf(`UI: slot-result has data-selectable=%q, expected "true" or "false"`, selectable)
	}
}
