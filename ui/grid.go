package ui

import (
	"math"

	"github.com/gpayer/pixelext/nodes"

	"github.com/faiface/pixel"
)

type Grid struct {
	UIBase
	bbox       *nodes.Canvas
	cols       int
	uichildren []UINode
}

func NewGrid(name string, cols int) *Grid {
	g := &Grid{
		UIBase: *NewUIBase(name),
		cols:   cols,
	}
	g.Self = g
	g.UISelf = g
	//g.SetZeroAlignment(nodes.AlignmentTopLeft)
	g.bbox = nodes.NewCanvas("__bbox", 1, 1)
	g.bbox.SetZIndex(-10)
	g.bbox.SetPos(pixel.V(0, 0))
	g.AddChild(g.bbox)
	return g
}

func (g *Grid) recalcPositions() {
	styles := g.GetStyles()
	maxx := make([]float64, 1)
	maxy := make([]float64, 1)
	row := 0
	curcol := 0
	var b pixel.Vec

	for _, uichild := range g.uichildren {
		uichild.SetAlignment(nodes.AlignmentTopLeft)
		cb := uichild.Size()
		if curcol >= len(maxx) {
			maxx = append(maxx, 0.0)
		}
		if row >= len(maxy) {
			maxy = append(maxy, 0.0)
		}
		curw := cb.X + 2*uichild.GetStyles().Padding
		if curw > maxx[curcol] {
			maxx[curcol] = curw
		}
		curh := cb.Y + 2*uichild.GetStyles().Padding
		if curh > maxy[row] {
			maxy[row] = curh
		}
		curcol++
		if curcol == g.cols {
			curcol = 0
			row++
		}
	}

	b = pixel.V(math.Round(sumSlice(maxx)), math.Round(sumSlice(maxy)))
	g.bbox.SetSize(b)
	g.bbox.Clear(styles.Background.Color)

	row = 0
	curcol = 0
	x := math.Round(-b.X / 2)
	y := math.Round(-b.Y / 2)
	for _, child := range g.uichildren {
		x += child.GetStyles().Padding
		y += child.GetStyles().Padding
		child.SetPos(pixel.V(x, -y))
		x += maxx[curcol] - child.GetStyles().Padding
		y -= child.GetStyles().Padding
		curcol++
		if curcol == g.cols {
			curcol = 0
			y += maxy[row]
			x = math.Round(-b.X / 2)
			row++
		}
	}

	if styles.Border.Width > 0 {
		g.bbox.DrawRect(pixel.V(0, 0), pixel.V(b.X-1, b.Y-1), styles.Border.Color)
		x = maxx[0]
		for i := 1; i < len(maxx); i++ {
			g.bbox.DrawLine(pixel.V(x, 0), pixel.V(x, b.Y), styles.Border.Color)
			x += maxx[i]
		}
		y = b.Y - maxy[0]
		for i := 1; i < len(maxy); i++ {
			g.bbox.DrawLine(pixel.V(0, y), pixel.V(b.X, y), styles.Border.Color)
			y -= maxy[i]
		}
	}

	g.UISelf.SetSize(b)
	nodes.SceneManager().Redraw()
}

func sumSlice(sl []float64) float64 {
	sum := 0.0
	for _, v := range sl {
		sum += v
	}
	return sum
}

func (g *Grid) AddChild(child nodes.Node) {
	g.BaseNode.AddChild(child)
	uichild, ok := child.(UINode)
	if ok {
		g.uichildren = append(g.uichildren, uichild)
	}
	if len(g.uichildren) > 1 {
		g.recalcPositions()
	}
}

func (g *Grid) SetStyles(styles *nodes.Styles) {
	g.BaseNode.SetStyles(styles)
	g.bbox.SetStyles(styles)
}

func (g *Grid) RemoveChild(child nodes.Node) {
	if child == g.bbox {
		return
	}
	g.UIBase.RemoveChild(child)
	uichild, ok := child.(UINode)
	if ok {
		for i, uic := range g.uichildren {
			if uichild == uic {
				l := len(g.uichildren)
				g.uichildren[i] = nil
				copy(g.uichildren[i:l-1], g.uichildren[i+1:l])
				g.uichildren = g.uichildren[:l-1]
				break
			}
		}
	}

	g.recalcPositions()
}

func (g *Grid) RemoveChildren() {
	g.UIBase.RemoveChildren()
	g.uichildren = make([]UINode, 0)
	g.AddChild(g.bbox)
	g.recalcPositions()
}

func (g *Grid) ChildChanged() {
	g.recalcPositions()
}

func (g *Grid) Cols() int {
	return g.cols
}

func (g *Grid) SetCols(cols int) {
	if cols > 0 {
		g.cols = cols
		g.recalcPositions()
	}
}
