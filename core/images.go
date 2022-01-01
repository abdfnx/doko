package core

import (
	"time"
	"strings"

	"github.com/abdfnx/doko/shared"
	"github.com/abdfnx/doko/docker"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/docker/docker/api/types"
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

func newImages(ui *UI) *images {
	images := &images{
		Table: tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
	}

	images.SetTitle("image list").SetTitleAlign(tview.AlignLeft)
	images.SetBorder(true)
	images.setEntries(ui)
	images.setKeybinding(ui)

	return images
}

func (i *images) name() string {
	return "images"
}

func (i *images) setKeybinding(ui *UI) {
	i.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		ui.setGlobalKeybinding(event)

		switch event.Key() {
			case tcell.KeyEnter:
				ui.inspectImage()
			case tcell.KeyCtrlL:
				ui.loadImageForm()
			case tcell.KeyCtrlR:
				i.setEntries(ui)
		}

		switch event.Rune() {
			case 'c':
				ui.createContainerForm()
			case 'p':
				ui.pullImageForm()
			case 'd':
				ui.removeImage()
			case 'i':
				ui.importImageForm()
			case 's':
				ui.saveImageForm()
			case 'f':
				newSearchInputField(ui)
		}

		return event
	})
}

func (i *images) entries(ui *UI) {
	images, err := docker.Client.Images(types.ImageListOptions{})

	if err != nil {
		return
	}

	ui.state.resources.images = make([]*image, 0)

	for _, imgInfo := range images {
		for _, repoTag := range imgInfo.RepoTags {
			repo, tag := shared.ParseRepoTag(repoTag)
			if strings.Index(repo, i.filterWord) == -1 {
				continue
			}

			ui.state.resources.images = append(ui.state.resources.images, &image{
				ID:      imgInfo.ID[7:19],
				Repo:    repo,
				Tag:     tag,
				Created: shared.ParseDateToString(imgInfo.Created),
				Size:    shared.ParseSizeToString(imgInfo.Size),
			})
		}
	}
}

func (i *images) setEntries(ui *UI) {
	i.entries(ui)

	table := i.Clear()

	headers := []string{
		"ID",
		"Repo",
		"Tag",
		"Created",
		"Size",
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

	for i, image := range ui.state.resources.images {
		table.SetCell(i+1, 0, tview.NewTableCell(image.ID).
			SetTextColor(tcell.ColorLightYellow).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 1, tview.NewTableCell(image.Repo).
			SetTextColor(tcell.ColorLightYellow).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 2, tview.NewTableCell(image.Tag).
			SetTextColor(tcell.ColorLightYellow).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 3, tview.NewTableCell(image.Created).
			SetTextColor(tcell.ColorLightYellow).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 4, tview.NewTableCell(image.Size).
			SetTextColor(tcell.ColorLightYellow).
			SetMaxWidth(1).
			SetExpansion(1))
	}
}

func (i *images) updateEntries(ui *UI) {
	go ui.app.QueueUpdateDraw(func() {
		i.setEntries(ui)
	})
}

func (i *images) focus(ui *UI) {
	i.SetSelectable(true, false)
	ui.app.SetFocus(i)
}

func (i *images) unfocus() {
	i.SetSelectable(false, false)
}

func (i *images) setFilterWord(word string) {
	i.filterWord = word
}

func (i *images) monitoringImages(ui *UI) {
	shared.Logger.Info("start monitoring images")
	ticker := time.NewTicker(5 * time.Second)

LOOP:
	for {
		select {
			case <-ticker.C:
				i.updateEntries(ui)
			case <-ui.state.stopChans["image"]:
				ticker.Stop()
				break LOOP
		}
	}

	shared.Logger.Info("stop monitoring images")
}
