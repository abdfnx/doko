package core

import (
	"fmt"
	"time"
	"strings"

	"github.com/abdfnx/doko/log"
	"github.com/abdfnx/doko/shared"
	"github.com/abdfnx/doko/docker"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/docker/docker/api/types"
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

func newNetworks(ui *UI) *networks {
	networks := &networks{
		Table: tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
	}

	networks.SetTitle("network list").SetTitleAlign(tview.AlignLeft)
	networks.SetBorder(true)
	networks.setEntries(ui)
	networks.setKeybinding(ui)

	return networks
}

func (n *networks) name() string {
	return "networks"
}

func (n *networks) setKeybinding(ui *UI) {
	n.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		ui.setGlobalKeybinding(event)

		switch event.Key() {
			case tcell.KeyEnter:
				ui.inspectNetwork()
			case tcell.KeyCtrlR:
				n.setEntries(ui)
		}

		switch event.Rune() {
			case 'd':
				ui.removeNetwork()
		}

		return event
	})
}

func (n *networks) entries(ui *UI) {
	networks, err := docker.Client.Networks(types.NetworkListOptions{})

	if err != nil {
		logger.Logger.Error(err)
		return
	}

	keys := make([]string, 0, len(networks))
	tmpMap := make(map[string]*network)

	for _, net := range networks {
		if strings.Index(net.Name, n.filterWord) == -1 {
			continue
		}

		var containers string

		net, err := docker.Client.InspectNetwork(net.ID)

		if err != nil {
			logger.Logger.Error(err)
			continue
		}

		for _, endpoint := range net.Containers {
			containers += fmt.Sprintf("%s ", endpoint.Name)
		}

		tmpMap[net.ID[:12]] = &network{
			ID:         net.ID,
			Name:       net.Name,
			Driver:     net.Driver,
			Scope:      net.Scope,
			containers: containers,
		}

		keys = append(keys, net.ID[:12])

	}

	ui.state.resources.networks = make([]*network, 0)

	for _, key := range shared.SortKeys(keys) {
		ui.state.resources.networks = append(ui.state.resources.networks, tmpMap[key])
	}
}

func (n *networks) setEntries(ui *UI) {
	n.entries(ui)
	table := n.Clear()

	headers := []string{
		"ID",
		"Name",
		"Driver",
		"Scope",
		"Containers",
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

	for i, network := range ui.state.resources.networks {
		table.SetCell(i+1, 0, tview.NewTableCell(network.ID).
			SetTextColor(tcell.ColorLightSkyBlue).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 1, tview.NewTableCell(network.Name).
			SetTextColor(tcell.ColorLightSkyBlue).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 2, tview.NewTableCell(network.Driver).
			SetTextColor(tcell.ColorLightSkyBlue).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 3, tview.NewTableCell(network.Scope).
			SetTextColor(tcell.ColorLightSkyBlue).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 4, tview.NewTableCell(network.containers).
			SetTextColor(tcell.ColorLightSkyBlue).
			SetMaxWidth(1).
			SetExpansion(1))
	}
}

func (n *networks) focus(ui *UI) {
	n.SetSelectable(true, false)
	ui.app.SetFocus(n)
}

func (n *networks) unfocus() {
	n.SetSelectable(false, false)
}

func (n *networks) updateEntries(ui *UI) {
	go ui.app.QueueUpdateDraw(func() {
		n.setEntries(ui)
	})
}

func (n *networks) setFilterWord(word string) {
	n.filterWord = word
}

func (n *networks) monitoringNetworks(ui *UI) {
	logger.Logger.Info("start monitoring networks")
	ticker := time.NewTicker(5 * time.Second)

LOOP:
	for {
		select {
			case <-ticker.C:
				n.updateEntries(ui)
			case <-ui.state.stopChans["network"]:
				ticker.Stop()
				break LOOP
		}
	}

	logger.Logger.Info("stop monitoring networks")
}
