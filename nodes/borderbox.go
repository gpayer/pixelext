package nodes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type BorderBox struct {
	Canvas
}

func NewBorderBox(name string, w, h float64) *BorderBox {
	b := &BorderBox{
		Canvas: *NewCanvas(name, w, h),
	}
	b.Self = b
	b.redrawCanvas()
	return b
}

func (b *BorderBox) redrawCanvas() {
	styles := b.GetStyles()
	im := imdraw.New(nil)
	canvas := b.Canvas.Canvas()
	canvas.Clear(styles.Background.Color)
	if styles.Border.Width > 0 {
		bounds := canvas.Bounds()
		im.Color = styles.Border.Color
		im.Push(pixel.V(0, 0), pixel.V(bounds.W(), 0), pixel.V(bounds.W(), bounds.H()), pixel.V(0, bounds.H()))
		im.Polygon(styles.Border.Width)
		im.Draw(canvas)
	}
}

func (b *BorderBox) SetBounds(r pixel.Rect) {
	b.Canvas.SetBounds(r)
	b.redrawCanvas()
}

func (b *BorderBox) SetStyles(style *Styles) {
	b.Canvas.SetStyles(style)
	b.redrawCanvas()
}
