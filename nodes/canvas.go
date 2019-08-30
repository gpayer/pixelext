package nodes

import (
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Canvas struct {
	BaseNode
	canvas *pixelgl.Canvas
	size   pixel.Vec
	pixels []uint8
}

func NewCanvas(name string, w, h float64) *Canvas {
	if w < 2 {
		w = 2
	}
	if h < 2 {
		h = 2
	}
	c := &Canvas{
		BaseNode: *NewBaseNode(name),
	}
	c.Self = c
	c.size = pixel.V(w, h)
	c.canvas = pixelgl.NewCanvas(pixel.R(0, 0, c.size.X, c.size.Y))
	c.pixels = make([]uint8, int(math.Ceil(w))*int(math.Ceil(h))*4)
	return c
}

func (c *Canvas) SetSize(size pixel.Vec) {
	if size.X < 1 {
		size.X = 1
	}
	if size.Y < 1 {
		size.Y = 1
	}
	c.canvas.SetBounds(pixel.R(0, 0, math.Ceil(size.X), math.Ceil(size.Y)))
	c.size = c.canvas.Bounds().Size()
	c.pixels = make([]uint8, int(math.Ceil(c.size.X))*int(math.Ceil(c.size.Y))*4)
	SceneManager().Redraw()
}

func (c *Canvas) Size() pixel.Vec {
	return c.size
}

func (c *Canvas) Draw(win pixel.Target, mat pixel.Matrix) {
	c.canvas.Draw(win, mat)
}

func (c *Canvas) Clear(col color.Color) {
	c.canvas.Clear(col)
	var ur, ug, ub, ua uint32 = col.RGBA()
	var bcol []byte = make([]byte, 4)
	bcol[0] = byte(ur)
	bcol[1] = byte(ug)
	bcol[2] = byte(ub)
	bcol[3] = byte(ua)
	for i := 0; i < len(c.pixels); i += 4 {
		copy(c.pixels[i:i+4], bcol)
	}
	// TODO: not sure why this leads to "wrong number of pixels errors", needs to be fixed
	//if len(c.pixels) > 0 {
	//	c.canvas.SetPixels(c.pixels)
	//}
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

func (c *Canvas) DrawLine(p1, p2 pixel.Vec, col color.Color) {
	var ur, ug, ub, ua uint32 = col.RGBA()
	var bcol []byte = make([]byte, 4)
	bcol[0] = byte(ur)
	bcol[1] = byte(ug)
	bcol[2] = byte(ub)
	bcol[3] = byte(ua)

	disp := &displayDef{
		x:   int(c.Size().X),
		y:   int(c.Size().Y),
		buf: c.pixels,
	}

	gbham(int(p1.X), int(p1.Y), int(p2.X), int(p2.Y), disp, bcol)

	if len(c.pixels) > 0 {
		c.canvas.SetPixels(disp.buf)
	}
}

func (c *Canvas) DrawRect(p1, p2 pixel.Vec, col color.Color) {
	c.DrawLine(p1, pixel.V(p1.X, p2.Y), col)
	c.DrawLine(p1, pixel.V(p2.X, p1.Y), col)
	c.DrawLine(p2, pixel.V(p1.X, p2.Y), col)
	c.DrawLine(p2, pixel.V(p2.X, p1.Y), col)
}

func (c *Canvas) FillRect(p1, p2 pixel.Vec, col color.Color) {
	var starty, endy, startx, endx int
	bounds := pixel.R(p1.X, p1.Y, p2.X, p2.Y).Norm()
	startx = int(bounds.Min.X)
	starty = int(bounds.Min.Y)
	endx = int(bounds.Max.X)
	endy = int(bounds.Max.Y)

	var ur, ug, ub, ua uint32 = col.RGBA()
	var bcol []byte = make([]byte, 4)
	bcol[0] = byte(ur)
	bcol[1] = byte(ug)
	bcol[2] = byte(ub)
	bcol[3] = byte(ua)

	disp := &displayDef{
		x:   int(c.Size().X),
		y:   int(c.Size().Y),
		buf: c.pixels,
	}

	for y := starty; y <= endy; y++ {
		for x := startx; x <= endx; x++ { // TODO: replace with clever implementation
			setPixel(x, y, disp, bcol)
		}
	}
	if len(c.pixels) > 0 {
		c.canvas.SetPixels(disp.buf)
	}
}

type displayDef struct {
	x, y int
	buf  []byte
}

func setPixel(x, y int, disp *displayDef, col []byte) {
	if x >= 0 && y >= 0 && x < disp.x && y < disp.y {
		start := 4 * (disp.x*y + x)
		copy(disp.buf[start:start+4], col)
	}
}

func sgn(x int) int {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	} else {
		return 0
	}
}

func gbham(xstart, ystart, xend, yend int, disp *displayDef, col []byte) {
	var x, y, t, dx, dy, incx, incy, pdx, pdy, ddx, ddy, deltaslowdirection, deltafastdirection, err int

	dx = xend - xstart
	dy = yend - ystart

	incx = sgn(dx)
	incy = sgn(dy)
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}

	if dx > dy {
		pdx = incx
		pdy = 0
		ddx = incx
		ddy = incy
		deltaslowdirection = dy
		deltafastdirection = dx
	} else {
		pdx = 0
		pdy = incy
		ddx = incx
		ddy = incy
		deltaslowdirection = dx
		deltafastdirection = dy
	}

	x = xstart
	y = ystart
	err = deltafastdirection / 2
	setPixel(x, y, disp, col)

	for t = 0; t < deltafastdirection; t++ {
		err -= deltaslowdirection
		if err < 0 {
			err += deltafastdirection
			x += ddx
			y += ddy
		} else {
			x += pdx
			y += pdy
		}
		setPixel(x, y, disp, col)
	}
}
