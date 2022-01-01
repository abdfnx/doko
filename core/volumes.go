package core

import (
	"time"
	"strings"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"

	"github.com/abdfnx/doko/shared"
	"github.com/abdfnx/doko/docker"
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

var replacer = strings.NewReplacer("T", " ", "Z", "")

func newVolumes(ui *UI) *volumes {
	volumes := &volumes{
		Table: tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
	}

	volumes.SetTitle("volume list").SetTitleAlign(tview.AlignLeft)
	volumes.SetBorder(true)
	volumes.setEntries(ui)
	volumes.setKeybinding(ui)

	return volumes
}

func (v *volumes) name() string {
	return "volumes"
}

func (v *volumes) setKeybinding(ui *UI) {
	v.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		ui.setGlobalKeybinding(event)
		switch event.Key() {
			case tcell.KeyEnter:
				ui.inspectVolume()
			case tcell.KeyCtrlR:
				v.setEntries(ui)
		}

		switch event.Rune() {
			case 'd':
				ui.removeVolume()
			case 'c':
				ui.createVolumeForm()
		}

		return event
	})
}

func (v *volumes) entries(ui *UI) {
	volumes, err := docker.Client.Volumes()
	if err != nil {
		shared.Logger.Error(err)
		return
	}

	keys := make([]string, 0, len(volumes))
	tmpMap := make(map[string]*volume)

	for _, vo := range volumes {
		if strings.Index(vo.Name, v.filterWord) == -1 {
			continue
		}

		tmpMap[vo.Name] = &volume{
			Name:       vo.Name,
			MountPoint: vo.Mountpoint,
			Driver:     vo.Driver,
			Created:    replacer.Replace(vo.CreatedAt),
		}

		keys = append(keys, vo.Name)
	}

	ui.state.resources.volumes = make([]*volume, 0)

	for _, key := range shared.SortKeys(keys) {
		ui.state.resources.volumes = append(ui.state.resources.volumes, tmpMap[key])
	}
}

func (v *volumes) setEntries(ui *UI) {
	v.entries(ui)
	table := v.Clear()

	headers := []string{
		"Name",
		"MountPoint",
		"Driver",
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

	for i, network := range ui.state.resources.volumes {
		table.SetCell(i+1, 0, tview.NewTableCell(network.Name).
			SetTextColor(tcell.ColorLightPink).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 1, tview.NewTableCell(network.MountPoint).
			SetTextColor(tcell.ColorLightPink).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 2, tview.NewTableCell(network.Driver).
			SetTextColor(tcell.ColorLightPink).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 3, tview.NewTableCell(network.Created).
			SetTextColor(tcell.ColorLightPink).
			SetMaxWidth(1).
			SetExpansion(1))
	}
}

func (v *volumes) focus(ui *UI) {
	v.SetSelectable(true, false)
	ui.app.SetFocus(v)
}

func (v *volumes) unfocus() {
	v.SetSelectable(false, false)
}

func (v *volumes) updateEntries(ui *UI) {
	go ui.app.QueueUpdateDraw(func() {
		v.setEntries(ui)
	})
}

func (v *volumes) setFilterWord(word string) {
	v.filterWord = word
}

func (v *volumes) monitoringVolumes(ui *UI) {
	shared.Logger.Info("start monitoring volumes")
	ticker := time.NewTicker(5 * time.Second)

LOOP:
	for {
		select {
			case <-ticker.C:
				v.updateEntries(ui)
			case <-ui.state.stopChans["volume"]:
				ticker.Stop()
				break LOOP
		}
	}

	shared.Logger.Info("stop monitoring volumes")
}
