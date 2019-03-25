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

// Clicked checks for mouse clicks inside of the given node
func (e *EventManager) Clicked(button pixelgl.Button, node nodeInternal) bool {
	b := node.(*BaseNode)
	if e.win.JustPressed(button) && b.bounds.Contains(b.mat.Unproject(e.win.MousePosition())) {
		return true
	}
	return false
}

func (e *EventManager) JustPressedButtons(buttons ...pixelgl.Button) []pixelgl.Button {
	var result []pixelgl.Button
	for _, b := range buttons {
		if e.win.JustPressed(b) {
			result = append(result, b)
		}
	}
	return result
}

func (e *EventManager) JustPressed(button pixelgl.Button) bool {
	return e.win.JustPressed(button)
}

func (e *EventManager) JustReleased(button pixelgl.Button) bool {
	return e.win.JustReleased(button)
}

func (e *EventManager) Pressed(button pixelgl.Button) bool {
	return e.win.Pressed(button)
}

func (e *EventManager) Repeated(button pixelgl.Button) bool {
	return e.win.Repeated(button)
}

func init() {
	events = &EventManager{}
}
