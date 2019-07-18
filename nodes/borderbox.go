package nodes

import (
	"github.com/faiface/pixel"
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
	canvas := b.Canvas.Canvas()
	b.Clear(styles.Background.Color)
	if styles.Border.Width > 0 {
		bounds := canvas.Bounds()
		col := styles.Border.Color
		//bw := styles.Border.Width / 2
		maxx := bounds.W() - 1
		maxy := bounds.H() - 1
		b.DrawRect(pixel.V(0, 0), pixel.V(maxx, maxy), col)
	}
	SceneManager().Redraw()
}

func (b *BorderBox) SetSize(size pixel.Vec) {
	b.Canvas.SetSize(size)
	b.redrawCanvas()
}

func (b *BorderBox) SetStyles(styles *Styles) {
	b.Canvas.SetStyles(styles)
	b.redrawCanvas()
}

func (b *BorderBox) Mount() {
	b.Canvas.Mount()
	b.redrawCanvas()
}

func (b *BorderBox) Unmount() {
	b.Canvas.Unmount()
}
