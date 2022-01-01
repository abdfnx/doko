package core

import (
	"github.com/rivo/tview"
)

type panels struct {
	currentPanel int
	panel        []panel
}

type resources struct {
	images     []*image
	networks   []*network
	volumes    []*volume
	containers []*container
	tasks      []*task
}

type state struct {
	panels          panels
	navigate        *navigate
	resources resources
	stopChans 		map[string]chan int
}

// UI struct have all `doko` panels
type UI struct {
	app   *tview.Application
	pages *tview.Pages
	state *state
}

func newState() *state {
	return &state{
		stopChans: make(map[string]chan int),
	}
}

func New() *UI {
	return &UI{
		app:   tview.NewApplication(),
		state: newState(),
	}
}
