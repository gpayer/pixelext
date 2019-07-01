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

func (c *Canvas) Size() pixel.Vec {
	return c.canvas.Bounds().Size()
}

func (c *Canvas) Draw(win pixel.Target, mat pixel.Matrix) {
	c.canvas.Draw(win, mat)
}

func (c *Canvas) Clear(color color.Color) {
	c.canvas.Clear(color)
	SceneManager().Redraw()
}

func (c *Canvas) Canvas() *pixelgl.Canvas {
	return c.canvas
}

func (c *Canvas) GetCanvas() *pixelgl.Canvas {
	return c.canvas
}

func (c *Canvas) SetUniform(name string, value interface{}) {
	c.canvas.SetUniform(name, value)
}

func (c *Canvas) SetFragmentShader(src string) {
	c.canvas.SetFragmentShader(src)
}

func (c *Canvas) Contains(point pixel.Vec) bool {
	size := c.canvas.Bounds().Size().Scaled(.5)
	bounds := pixel.R(-size.X, -size.Y, size.X, size.Y)
	return bounds.Contains(point)
}
