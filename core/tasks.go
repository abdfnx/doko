package core

import (
	"context"

	"github.com/gdamore/tcell/v2"
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

func newTasks(ui *UI) *tasks {
	tasks := &tasks{
		Table: tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
		tasks: make(chan *task),
	}

	tasks.SetTitle("your docker tasks").SetTitleAlign(tview.AlignLeft)
	tasks.SetBorder(true)
	tasks.setEntries(ui)
	tasks.setKeybinding(ui)

	return tasks
}

func (t *tasks) name() string {
	return "tasks"
}

func (t *tasks) setKeybinding(ui *UI) {
	t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		ui.setGlobalKeybinding(event)

		return event
	})
}

func (t *tasks) entries(ui *UI) {}

func (t *tasks) setEntries(ui *UI) {
	t.entries(ui)
	table := t.Clear()

	headers := []string{
		"Name",
		"Status",
		"Created",
	}

	for i, header := range headers {
		table.SetCell(0, i, &tview.TableCell{
			Text:            header,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorWhite,
			BackgroundColor: tcell.ColorDefault,
			Attributes:      tcell.AttrBold,
		})
	}

	for i, task := range ui.state.resources.tasks {
		table.SetCell(i+1, 0, tview.NewTableCell(task.Name).
			SetTextColor(tcell.ColorLightGreen).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 1, tview.NewTableCell(task.Status).
			SetTextColor(tcell.ColorLightGreen).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 2, tview.NewTableCell(task.Created).
			SetTextColor(tcell.ColorLightGreen).
			SetMaxWidth(1).
			SetExpansion(1))

	}
}

func (t *tasks) focus(ui *UI) {
	t.SetSelectable(true, false)
	ui.app.SetFocus(t)
}

func (t *tasks) unfocus() {
	t.SetSelectable(false, false)
}

func (t *tasks) setFilterWord(word string) {}

func (t *tasks) updateEntries(ui *UI) {}
