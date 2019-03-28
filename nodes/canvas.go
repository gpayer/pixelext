package nodes

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Canvas struct {
	BaseNode
	canvas *pixelgl.Canvas
}

func NewCanvas(name string, w, h float64) *Canvas {
	c := &Canvas{
		BaseNode: *NewBaseNode(name),
		canvas:   pixelgl.NewCanvas(pixel.R(0, 0, w, h)),
	}
	c.Self = c
	c.SetExtraOffset(pixel.V(w/2, h/2))
	return c
}

func (c *Canvas) SetBounds(r pixel.Rect) {
	c.BaseNode.SetBounds(r)
	bounds := c.GetBounds()
	c.canvas.SetBounds(pixel.R(0, 0, bounds.W(), bounds.H()))
	c.SetExtraOffset(pixel.V(bounds.W()/2, bounds.H()/2))
}

func (c *Canvas) Draw(t pixel.Target, mat pixel.Matrix) {
	c.canvas.Draw(t, mat)
}

func (c *Canvas) Clear(color color.Color) {
	c.canvas.Clear(color)
}

func (c *Canvas) Canvas() *pixelgl.Canvas {
	return c.canvas
}
