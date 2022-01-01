package core

import (
	"github.com/rivo/tview"
)

type container struct {
	ID      string
	Name    string
	Image   string
	Status  string
	Created string
	Port    string
}

type containers struct {
	*tview.Table
	filterWord string
}
