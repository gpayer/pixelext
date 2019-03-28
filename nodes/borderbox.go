package nodes

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"

	"golang.org/x/image/colornames"
)

type BorderBox struct {
	Canvas
	bordercolor, backgroundcolor color.Color
	borderwidth                  float64
}

func NewBorderBox(name string, w, h float64) *BorderBox {
	b := &BorderBox{
		Canvas:          *NewCanvas(name, w, h),
		bordercolor:     colornames.White,
		backgroundcolor: colornames.Black,
		borderwidth:     2,
	}
	b.Self = b
	b.redrawCanvas()
	return b
}

func (b *BorderBox) redrawCanvas() {
	im := imdraw.New(nil)
	canvas := b.Canvas.Canvas()
	canvas.Clear(b.backgroundcolor)
	if b.borderwidth > 0 {
		bounds := canvas.Bounds()
		im.Color = b.bordercolor
		im.Push(pixel.V(0, 0), pixel.V(bounds.W(), 0), pixel.V(bounds.W(), bounds.H()), pixel.V(0, bounds.H()))
		im.Polygon(b.borderwidth)
		im.Draw(canvas)
	}
}

func (b *BorderBox) SetBounds(r pixel.Rect) {
	b.Canvas.SetBounds(r)
	b.redrawCanvas()
}

func (b *BorderBox) SetBorderColor(col color.Color) {
	b.bordercolor = col
	b.redrawCanvas()
}

func (b *BorderBox) SetBackgroundColor(col color.Color) {
	b.backgroundcolor = col
	b.redrawCanvas()
}

func (b *BorderBox) SetBorderWidth(w float64) {
	b.borderwidth = w
	b.redrawCanvas()
}
