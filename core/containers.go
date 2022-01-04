package core

import (
	"time"
	"strings"

	"github.com/abdfnx/doko/log"
	"github.com/abdfnx/doko/shared"
	"github.com/abdfnx/doko/docker"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/docker/docker/api/types"
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

func newContainers(ui *UI) *containers {
	containers := &containers{
		Table: tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
	}

	containers.SetTitle("container list").SetTitleAlign(tview.AlignLeft)
	containers.SetBorder(true)
	containers.setEntries(ui)
	containers.setKeybinding(ui)

	return containers
}

func (c *containers) name() string {
	return "containers"
}

func (c *containers) setKeybinding(ui *UI) {
	c.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		ui.setGlobalKeybinding(event)
		switch event.Key() {
			case tcell.KeyEnter:
				ui.inspectContainer()
			case tcell.KeyCtrlE:
				ui.attachContainerForm()
			case tcell.KeyCtrlL:
				ui.tailContainerLog()
			case tcell.KeyCtrlK:
				ui.killContainer()
			case tcell.KeyCtrlR:
				c.setEntries(ui)
		}

		switch event.Rune() {
			case 'd':
				ui.removeContainer()
			case 'r':
				ui.renameContainerForm()
			case 't':
				ui.startContainer()
			case 's':
				ui.stopContainer()
			case 'e':
				ui.exportContainerForm()
			case 'c':
				ui.commitContainerForm()
		}

		return event
	})
}

func (c *containers) entries(ui *UI) {
	containers, err := docker.Client.Containers(types.ContainerListOptions{All: true})

	if err != nil {
		logger.Logger.Error(err)
		return
	}

	ui.state.resources.containers = make([]*container, 0)

	for _, con := range containers {
		if strings.Index(con.Names[0][1:], c.filterWord) == -1 {
			continue
		}

		ui.state.resources.containers = append(ui.state.resources.containers, &container{
			ID:      con.ID[:12],
			Image:   con.Image,
			Name:    con.Names[0][1:],
			Status:  con.Status,
			Created: shared.ParseDateToString(con.Created),
			Port:    shared.ParsePortToString(con.Ports),
		})
	}
}

func (c *containers) setEntries(ui *UI) {
	c.entries(ui)
	table := c.Clear()

	headers := []string{
		"ID",
		"Name",
		"Image",
		"Status",
		"Created",
		"Port",
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

	for i, container := range ui.state.resources.containers {
		table.SetCell(i+1, 0, tview.NewTableCell(container.ID).
			SetTextColor(tcell.ColorLightGreen).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 1, tview.NewTableCell(container.Name).
			SetTextColor(tcell.ColorLightGreen).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 2, tview.NewTableCell(container.Image).
			SetTextColor(tcell.ColorLightGreen).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 3, tview.NewTableCell(container.Status).
			SetTextColor(tcell.ColorLightGreen).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 4, tview.NewTableCell(container.Created).
			SetTextColor(tcell.ColorLightGreen).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 5, tview.NewTableCell(container.Port).
			SetTextColor(tcell.ColorLightGreen).
			SetMaxWidth(1).
			SetExpansion(1))
	}
}

func (c *containers) focus(ui *UI) {
	c.SetSelectable(true, false)
	ui.app.SetFocus(c)
}

func (c *containers) unfocus() {
	c.SetSelectable(false, false)
}

func (c *containers) updateEntries(ui *UI) {
	go ui.app.QueueUpdateDraw(func() {
		c.setEntries(ui)
	})
}

func (c *containers) setFilterWord(word string) {
	c.filterWord = word
}

func (c *containers) monitoringContainers(ui *UI) {
	logger.Logger.Info("start monitoring containers")
	ticker := time.NewTicker(5 * time.Second)

LOOP:
	for {
		select {
			case <-ticker.C:
				c.updateEntries(ui)
			case <-ui.state.stopChans["container"]:
				ticker.Stop()
				break LOOP
		}
	}

	logger.Logger.Info("stop monitoring containers")
}
