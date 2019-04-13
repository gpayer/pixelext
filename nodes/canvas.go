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
	return c
}

func (c *Canvas) SetSize(size pixel.Vec) {
	c.canvas.SetBounds(pixel.R(0, 0, size.X, size.Y))
	SceneManager().Redraw()
}

func (c *Canvas) Draw(win *pixelgl.Window, mat pixel.Matrix) {
	c.canvas.Draw(win, mat)
}

func (c *Canvas) Clear(color color.Color) {
	c.canvas.Clear(color)
	SceneManager().Redraw()
}

func (c *Canvas) Canvas() *pixelgl.Canvas {
	return c.canvas
}
