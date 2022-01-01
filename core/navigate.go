package core

import (
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
)

type navigate struct {
	*tview.TextView
	keybindings map[string]string
}

func newNavigate() *navigate {
	return &navigate{
		TextView: tview.NewTextView().SetTextColor(tcell.ColorYellow),
		keybindings: map[string]string{
			"images":     " p: pull image, i: import image, s: save image, Ctrl+l: load image, f: find image, /: filter d: delete image,\n c: create new container, Enter: inspect image, Ctrl+r: refresh images list",
			"containers": " e: export container, c: commit container, /: filter, Ctrl+e: exec container cmd, t: start container, s: stop container,\n Ctrl+k: kill container, d: delete container, Enter: inspect container, Ctrl+r: refresh container list, Ctrl+l: show container logs",
			"networks":   " d: delete network, Enter: inspect network, /: filter",
			"volumes":    " c: create volume, d: delete volume\n /: filter, Enter: inspect volume, Ctrl+r: refresh volume list",
		},
	}
}

func (n *navigate) update(panel string) {
	n.SetText(n.keybindings[panel])
}
