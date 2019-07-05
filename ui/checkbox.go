package ui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gpayer/pixelext/nodes"
	"golang.org/x/image/colornames"
)

type CheckBox struct {
	UIBase
	box             *nodes.Canvas
	oldstate, state bool
	onchange        func(v bool)
}

func NewCheckBox(name string, w, h float64, state bool) *CheckBox {
	c := &CheckBox{
		UIBase:   *NewUIBase(name),
		state:    state,
		oldstate: false,
		onchange: func(v bool) {},
	}
	c.Self = c
	c.UISelf = c
	c.UIBase.SetSize(pixel.V(w, h))
	return c
}

func (c *CheckBox) drawFalse() {
	size := c.Size()
	c.box.Clear(colornames.Black)
	c.box.DrawRect(pixel.ZV, pixel.V(size.X-1, size.Y-1), colornames.White)
}

func (c *CheckBox) drawTrue() {
	size := c.Size()
	c.box.Clear(colornames.Skyblue)
	c.box.DrawRect(pixel.ZV, pixel.V(size.X-1, size.Y-1), colornames.White)
}

func (c *CheckBox) SetSize(size pixel.Vec) {
	c.UIBase.SetSize(size)
	c.box.SetSize(size)
}

func (c *CheckBox) State() bool {
	return c.state
}

func (c *CheckBox) SetState(state bool) {
	c.state = state
}

func (c *CheckBox) OnChange(fn func(v bool)) {
	c.onchange = fn
}

func (c *CheckBox) Init() {
	size := c.Size()
	c.box = nodes.NewCanvas("box", size.X, size.Y)
	c.box.SetLocked(true)
	c.AddChild(c.box)
	c.drawFalse()
}

func (c *CheckBox) Update(dt float64) {
	if nodes.Events().Clicked(pixelgl.MouseButtonLeft, c) {
		c.state = !c.state
		c.onchange(c.state)
	}
	if c.oldstate != c.state {
		if c.state {
			c.drawTrue()
		} else {
			c.drawFalse()
		}
		c.oldstate = c.state
	}
}
