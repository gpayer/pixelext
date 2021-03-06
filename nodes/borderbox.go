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
		bw := styles.Border.Width
		minx, miny := 0.0, 0.0
		maxx := bounds.W() - 1
		maxy := bounds.H() - 1
		for bw > 0 {
			b.DrawRect(pixel.V(minx, miny), pixel.V(maxx, maxy), col)
			bw -= 1
			minx += 1
			miny += 1
			maxx -= 1
			maxy -= 1
		}
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

func (b *BorderBox) KeepOnScreen() {
	globalcenter := b.LocalToGlobalPos(pixel.ZV)
	sizeh := b.Size().Scaled(.5)
	bounds := pixel.R(globalcenter.X-sizeh.X, globalcenter.Y-sizeh.Y, globalcenter.X+sizeh.X, globalcenter.Y+sizeh.Y)
	move := SceneManager().KeepOnScreen(b, bounds)
	b.SetPos(b.GetPos().Add(move))
}
