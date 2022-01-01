package core

import (
	"strconv"

	"github.com/abdfnx/doko/shared"
	"github.com/abdfnx/doko/docker"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
)

type searchImageResult struct {
	Name        string
	Stars       string
	Official    string
	Description string
}

type searchImageResults struct {
	keyword            string
	searchImageResults []*searchImageResult
	*tview.Table
}

func newSearchImageResults(ui *UI, keyword string) *searchImageResults {
	searchImageResults := &searchImageResults{
		keyword: keyword,
		Table:   tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
	}

	searchImageResults.SetTitle("search result").SetTitleAlign(tview.AlignLeft)
	searchImageResults.SetBorder(true)
	searchImageResults.setEntries(ui)
	searchImageResults.setKeybinding(ui)

	return searchImageResults
}

func newSearchInputField(ui *UI) {
	viewName := "searchImageInput"
	searchInput := tview.NewInputField().SetLabel("Image")
	searchInput.SetLabelWidth(6)
	searchInput.SetBorder(true)

	searchInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			if searchInput.GetText() == "" {
				ui.message("please input some text", "OK", ui.currentPanel().name(), func() {})
				return
			}

			ui.pages.AddAndSwitchToPage("searchImageResults", ui.modal(newSearchImageResults(ui, searchInput.GetText()), 100, 50), true).ShowPage("main")
		}
	})

	closeSearchInput := func() {
		currentPanel := ui.state.panels.panel[ui.state.panels.currentPanel]
		ui.closeAndSwitchPanel(viewName, currentPanel.name())
	}

	searchInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
			case tcell.KeyEsc:
				closeSearchInput()
		}

		return event
	})

	ui.pages.AddAndSwitchToPage(viewName, ui.modal(searchInput, 80, 3), true).ShowPage("main")
}

func (s *searchImageResults) name() string {
	return "searchImageResults"
}

func (s *searchImageResults) pullImage(ui *UI) {
	currentPanel := ui.state.panels.panel[ui.state.panels.currentPanel]
	ui.pullImage(s.selectedSearchImageResult().Name, s.name(), currentPanel.name())
}

func (s *searchImageResults) setKeybinding(ui *UI) {
	s.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
			case tcell.KeyEsc:
				s.closePanel(ui)
			case tcell.KeyEnter:
				s.pullImage(ui)
		}

		switch event.Rune() {
			case 'p':
				s.pullImage(ui)
			case 'q':
				s.closePanel(ui)
		}

		return event
	})
}

func (s *searchImageResults) entries(ui *UI) {
	images, err := docker.Client.SearchImage(s.keyword)

	if err != nil {
		ui.message("error", "OK", s.name(), func() {})
		return
	}

	if len(images) == 0 {
		ui.message("no image found", "OK", s.name(), func() {})
		return
	}

	s.searchImageResults = make([]*searchImageResult, 0)

	var official string

	for _, image := range images {
		if image.IsOfficial {
			official = "[OK]"
		}

		s.searchImageResults = append(s.searchImageResults, &searchImageResult{
			Name:        image.Name,
			Stars:       strconv.Itoa(image.StarCount),
			Official:    official,
			Description: shared.CutNewline(image.Description),
		})
	}

}

func (s *searchImageResults) setEntries(ui *UI) {
	s.entries(ui)
	table := s.Clear()

	headers := []string{
		"Name",
		"Star",
		"Official",
		"Description",
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

	for i, image := range s.searchImageResults {
		table.SetCell(i+1, 0, tview.NewTableCell(image.Name).
			SetTextColor(tcell.ColorLightYellow).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 1, tview.NewTableCell(image.Stars).
			SetTextColor(tcell.ColorLightYellow))

		table.SetCell(i+1, 2, tview.NewTableCell(image.Official).
			SetTextColor(tcell.ColorLightYellow))

		table.SetCell(i+1, 3, tview.NewTableCell(image.Description).
			SetTextColor(tcell.ColorLightYellow).
			SetMaxWidth(1).
			SetExpansion(1))

	}
}

func (s *searchImageResults) updateEntries(ui *UI) {
	go ui.app.QueueUpdateDraw(func() {
		s.setEntries(ui)
	})
}

func (s *searchImageResults) focus(ui *UI) {
	s.SetSelectable(true, false)
	ui.app.SetFocus(s)
}

func (s *searchImageResults) unfocus() {
	s.SetSelectable(false, false)
}

func (s *searchImageResults) closePanel(ui *UI) {
	currentPanel := ui.state.panels.panel[ui.state.panels.currentPanel]
	ui.closeAndSwitchPanel(s.name(), currentPanel.name())
}

func (s *searchImageResults) selectedSearchImageResult() *searchImageResult {
	row, _ := s.GetSelection()

	if len(s.searchImageResults) == 0 || row-1 < 0 {
		return nil
	}

	return s.searchImageResults[row-1]
}
