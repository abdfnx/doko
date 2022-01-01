package core

import (
	"github.com/rivo/tview"
)

type image struct {
	ID      string
	Repo    string
	Tag     string
	Created string
	Size    string
}

type images struct {
	*tview.Table
	filterWord string
}
