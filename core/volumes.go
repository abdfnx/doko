package core

import (
	"github.com/rivo/tview"
)

type volume struct {
	Name       string
	MountPoint string
	Driver     string
	Created    string
}

type volumes struct {
	*tview.Table
	filterWord string
}
