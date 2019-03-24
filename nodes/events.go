package nodes

import (
	"github.com/faiface/pixel/pixelgl"
)

var events *EventManager

func Events() *EventManager {
	return events
}

type EventManager struct {
	win *pixelgl.Window
}

func (e *EventManager) SetWin(win *pixelgl.Window) {
	e.win = win
}

func (e *EventManager) Clicked(button pixelgl.Button, node nodeInternal) bool {
	b := node.(*BaseNode)
	if e.win.JustPressed(button) && b.bounds.Contains(b.mat.Unproject(e.win.MousePosition())) {
		return true
	}
	return false
}

func init() {
	events = &EventManager{}
}
