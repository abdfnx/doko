package core

import (
	"context"

	"github.com/abdfnx/doko/shared"

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

func (ui *UI) imagePanel() *images {
	for _, panel := range ui.state.panels.panel {
		if panel.name() == "images" {
			return panel.(*images)
		}
	}

	return nil
}

func (ui *UI) containerPanel() *containers {
	for _, panel := range ui.state.panels.panel {
		if panel.name() == "containers" {
			return panel.(*containers)
		}
	}

	return nil
}

func (ui *UI) volumePanel() *volumes {
	for _, panel := range ui.state.panels.panel {
		if panel.name() == "volumes" {
			return panel.(*volumes)
		}
	}

	return nil
}

func (ui *UI) networkPanel() *networks {
	for _, panel := range ui.state.panels.panel {
		if panel.name() == "networks" {
			return panel.(*networks)
		}
	}

	return nil
}

func (ui *UI) taskPanel() *tasks {
	for _, panel := range ui.state.panels.panel {
		if panel.name() == "tasks" {
			return panel.(*tasks)
		}
	}

	return nil
}

func (ui *UI) monitoringTask() {
	shared.Logger.Info("start monitoring task")
LOOP:
	for {
		select {
			case task := <-ui.taskPanel().tasks:
				go func() {
					if err := task.Func(task.Ctx); err != nil {
						task.Status = err.Error()
					} else {
						task.Status = success
					}

					ui.updateTask()
				}()

			case <-ui.state.stopChans["task"]:
				shared.Logger.Info("stop monitoring task")
				break LOOP
			}
	}
}

func (ui *UI) startTask(taskName string, f func(ctx context.Context) error) {
	ctx, cancel := context.WithCancel(context.Background())

	task := &task{
		Name:    taskName,
		Status:  running,
		Created: shared.DateNow(),
		Func:    f,
		Ctx:     ctx,
		Cancel:  cancel,
	}

	ui.state.resources.tasks = append(ui.state.resources.tasks, task)
	ui.updateTask()
	ui.taskPanel().tasks <- task
}

func (ui *UI) cancelTask() {
	taskPanel := ui.taskPanel()
	row, _ := taskPanel.GetSelection()

	task := ui.state.resources.tasks[row-1]

	if task.Status == running {
		task.Cancel()
		task.Status = cancel
		ui.updateTask()
	}
}

func (ui *UI) updateTask() {
	go ui.app.QueueUpdateDraw(func() {
		ui.taskPanel().setEntries(ui)
	})
}

func (ui *UI) initPanels() {
	tasks := newTasks(ui)
	images := newImages(ui)
	containers := newContainers(ui)
	volumes := newVolumes(ui)
	networks := newNetworks(ui)
	info := newInfo()
	nav := newNavigate()

	ui.state.panels.panel = append(ui.state.panels.panel, tasks)
	ui.state.panels.panel = append(ui.state.panels.panel, images)
	ui.state.panels.panel = append(ui.state.panels.panel, containers)
	ui.state.panels.panel = append(ui.state.panels.panel, volumes)
	ui.state.panels.panel = append(ui.state.panels.panel, networks)
	ui.state.navigate = nav

	grid := tview.NewGrid().SetRows(2, 0, 0, 0, 0, 0, 2).
		AddItem(info, 0, 0, 1, 1, 0, 0, true).
		AddItem(tasks, 1, 0, 1, 1, 0, 0, true).
		AddItem(images, 2, 0, 1, 1, 0, 0, true).
		AddItem(containers, 3, 0, 1, 1, 0, 0, true).
		AddItem(volumes, 4, 0, 1, 1, 0, 0, true).
		AddItem(networks, 5, 0, 1, 1, 0, 0, true).
		AddItem(nav, 6, 0, 1, 1, 0, 0, true)

	ui.pages = tview.NewPages().
		AddAndSwitchToPage("main", grid, true)

	ui.app.SetRoot(ui.pages, true)
	ui.switchPanel("images")
}

func (ui *UI) startMonitoring() {
	stop := make(chan int, 1)
	ui.state.stopChans["task"] = stop
	ui.state.stopChans["image"] = stop
	ui.state.stopChans["volume"] = stop
	ui.state.stopChans["network"] = stop
	ui.state.stopChans["container"] = stop

	go ui.monitoringTask()
	go ui.imagePanel().monitoringImages(ui)
	go ui.networkPanel().monitoringNetworks(ui)
	go ui.volumePanel().monitoringVolumes(ui)
	go ui.containerPanel().monitoringContainers(ui)
}

func (ui *UI) stopMonitoring() {
	ui.state.stopChans["task"] <- 1
	ui.state.stopChans["image"] <- 1
	ui.state.stopChans["volume"] <- 1
	ui.state.stopChans["network"] <- 1
	ui.state.stopChans["container"] <- 1
}

// Start start application
func (ui *UI) Start() error {
	ui.initPanels()
	ui.startMonitoring()
	if err := ui.app.Run(); err != nil {
		ui.app.Stop()
		return err
	}

	return nil
}

// Stop stop application
func (ui *UI) Stop() error {
	ui.stopMonitoring()
	ui.app.Stop()
	return nil
}

func (ui *UI) selectedImage() *image {
	row, _ := ui.imagePanel().GetSelection()
	if len(ui.state.resources.images) == 0 {
		return nil
	}

	if row-1 < 0 {
		return nil
	}

	return ui.state.resources.images[row-1]
}

func (ui *UI) selectedContainer() *container {
	row, _ := ui.containerPanel().GetSelection()
	if len(ui.state.resources.containers) == 0 {
		return nil
	}

	if row-1 < 0 {
		return nil
	}

	return ui.state.resources.containers[row-1]
}

func (ui *UI) selectedVolume() *volume {
	row, _ := ui.volumePanel().GetSelection()
	if len(ui.state.resources.volumes) == 0 {
		return nil
	}

	if row-1 < 0 {
		return nil
	}

	return ui.state.resources.volumes[row-1]
}

func (ui *UI) selectedNetwork() *network {
	row, _ := ui.networkPanel().GetSelection()
	if len(ui.state.resources.networks) == 0 {
		return nil
	}

	if row-1 < 0 {
		return nil
	}

	return ui.state.resources.networks[row-1]
}

func (ui *UI) message(message, doneLabel, page string, doneFunc func()) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{doneLabel}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			ui.closeAndSwitchPanel("modal", page)
			if buttonLabel == doneLabel {
				doneFunc()
			}
		})

	ui.pages.AddAndSwitchToPage("modal", ui.modal(modal, 80, 29), true).ShowPage("main")
}

func (ui *UI) confirm(message, doneLabel, page string, doneFunc func()) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{doneLabel, "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			ui.closeAndSwitchPanel("modal", page)
			if buttonLabel == doneLabel {
				doneFunc()
			}
		})

	ui.pages.AddAndSwitchToPage("modal", ui.modal(modal, 80, 29), true).ShowPage("main")
}

func (ui *UI) switchPanel(panelName string) {
	for i, panel := range ui.state.panels.panel {
		if panel.name() == panelName {
			ui.state.navigate.update(panelName)
			panel.focus(ui)
			ui.state.panels.currentPanel = i
		} else {
			panel.unfocus()
		}
	}
}

func (ui *UI) closeAndSwitchPanel(removePanel, switchPanel string) {
	ui.pages.RemovePage(removePanel).ShowPage("main")
	ui.switchPanel(switchPanel)
}

func (ui *UI) modal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func (ui *UI) currentPanel() panel {
	return ui.state.panels.panel[ui.state.panels.currentPanel]
}
