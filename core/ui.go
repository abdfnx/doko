package core

import (
	"github.com/rivo/tview"
)

type panels struct {
	currentPanel int
	panel        []panel
}

type state struct {
	panels    panels
	navigate  *navigate
}

// UI struct have all `doko` panels
type UI struct {
	app   *tview.Application
	pages *tview.Pages
	state *state
}
