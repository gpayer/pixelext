package nodes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var events *EventManager

func Events() *EventManager {
	return events
}

type EventManager struct {
	win                *pixelgl.Window
	focused            Node
	handledButtons     map[pixelgl.Button]bool
	mouseScrollHandled bool
}

func (e *EventManager) reset() {
	for b, _ := range e.handledButtons {
		e.handledButtons[b] = false
	}
	e.mouseScrollHandled = false
}

func (e *EventManager) setWin(win *pixelgl.Window) {
	e.win = win
}

// Clicked checks for mouse clicks inside of the given node
func (e *EventManager) Clicked(button pixelgl.Button, node Node) bool {
	if !e.IsButtonHandled(button) && e.win.JustPressed(button) && e.MouseHovering(node) {
		e.handledButtons[button] = true
		return true
	}
	return false
}

func (e *EventManager) MouseHovering(node Node) bool {
	if node.Contains(node._getLastMat().Unproject(e.win.MousePosition())) {
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

func (e *EventManager) MouseScroll() pixel.Vec {
	e.mouseScrollHandled = true
	return e.win.MouseScroll()
}

func (e *EventManager) MousePosition() pixel.Vec {
	return e.win.MousePosition()
}

func (e *EventManager) LocalMousePosition(node Node) pixel.Vec {
	return node._getLastMat().Unproject(e.win.MousePosition()).Add(node.GetExtraOffset())
}

func (e *EventManager) MousePreviousPosition() pixel.Vec {
	return e.win.MousePreviousPosition()
}

func (e *EventManager) Typed() string {
	return e.win.Typed()
}

// SetFocus is used to set the currently focused Node. Pass nil to implicitly give every
// Node the input focus.
func (e *EventManager) SetFocus(node Node) {
	e.focused = node
}

// IsFocused can be used by a Node to check whether it has the input focus and can use keyboard inputs.
// Special case nil: in this case every Node implicitly has the input focus.
func (e *EventManager) IsFocused(node Node) bool {
	if e.focused == nil {
		return true
	}
	return e.focused == node
}

func (e *EventManager) IsButtonHandled(b pixelgl.Button) bool {
	handled, ok := e.handledButtons[b]
	if ok {
		return handled
	} else {
		return false
	}
}

func (e *EventManager) IsMouseScrollHandled() bool {
	return e.mouseScrollHandled
}

func init() {
	events = &EventManager{
		handledButtons: make(map[pixelgl.Button]bool, 0),
	}
}
