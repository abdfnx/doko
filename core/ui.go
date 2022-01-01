package core

import (
	"github.com/rivo/tview"
)

type panels struct {
	currentPanel int
	panel        []panel
}

type dockerResources struct {
	images     []*image
	networks   []*network
	volumes    []*volume
}

type state struct {
	panels          panels
	navigate        *navigate
	dockerResources dockerResources
	stopChans 		map[string]chan int
}

// UI struct have all `doko` panels
type UI struct {
	app   *tview.Application
	pages *tview.Pages
	state *state
}
