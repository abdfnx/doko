package core

import (
	"github.com/rivo/tview"
)

// UI struct have all `doko` panels
type UI struct {
	app   *tview.Application
	pages *tview.Pages
}
