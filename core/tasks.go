package core

import (
	"context"

	"github.com/rivo/tview"
)

var (
	success   = "Success"
	running   = "Running"
	cancel    = "canceled"
)

type task struct {
	Name    string
	Status  string
	Created string
	Func    func(ctx context.Context) error
	Ctx     context.Context
	Cancel  context.CancelFunc
}

type tasks struct {
	*tview.Table
	tasks chan *task
}
