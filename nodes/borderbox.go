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
		bw := styles.Border.Width / 2
		im.Push(pixel.V(0, bw), pixel.V(bounds.W(), bw))
		im.Line(styles.Border.Width)
		im.Push(pixel.V(bounds.W()-bw, 0), pixel.V(bounds.W()-bw, bounds.H()))
		im.Line(styles.Border.Width)
		im.Push(pixel.V(bounds.W(), bounds.H()-bw), pixel.V(0, bounds.H()-bw))
		im.Line(styles.Border.Width)
		im.Push(pixel.V(bw, bounds.H()), pixel.V(bw, 0))
		im.Line(styles.Border.Width)
		im.Draw(canvas)
	}
}

func (b *BorderBox) SetSize(size pixel.Vec) {
	b.Canvas.SetSize(size)
	b.redrawCanvas()
}

func (b *BorderBox) SetStyles(style *Styles) {
	b.Canvas.SetStyles(style)
	b.redrawCanvas()
}
