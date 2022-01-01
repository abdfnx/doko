package core

import (
	"github.com/rivo/tview"
)

type network struct {
	ID         string
	Name       string
	Driver     string
	Scope      string
	containers string
}

type networks struct {
	*tview.Table
	filterWord string
}
