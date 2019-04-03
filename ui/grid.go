package ui

import (
	"math"
	"pixelext/nodes"

	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
)

type Grid struct {
	nodes.BaseNode
	bbox *nodes.Canvas
	cols int
}

func NewGrid(name string, cols int) *Grid {
	g := &Grid{
		BaseNode: *nodes.NewBaseNode(name),
		cols:     cols,
	}
	g.Self = g
	//g.SetZeroAlignment(nodes.AlignmentTopLeft)
	g.bbox = nodes.NewCanvas("__bbox", 1, 1)
	g.bbox.SetZIndex(-10)
	g.bbox.SetPos(pixel.V(0, 0))
	g.bbox.SetZeroAlignment(nodes.AlignmentTopLeft)
	g.AddChild(g.bbox)
	return g
}

func (g *Grid) recalcPositions() {
	styles := g.GetStyles()
	maxx := make([]float64, 1)
	maxy := make([]float64, 1)
	row := 0
	curcol := 0
	var b pixel.Rect

	for _, child := range g.Children() {
		if child.GetName() == "__bbox" {
			continue
		}
		cb := child.GetContainerBounds()
		if curcol >= len(maxx) {
			maxx = append(maxx, 0.0)
		}
		if row >= len(maxy) {
			maxy = append(maxy, 0.0)
		}
		curw := cb.W() + 2*child.GetStyles().Padding
		if curw > maxx[curcol] {
			maxx[curcol] = curw
		}
		curh := cb.H() + 2*child.GetStyles().Padding
		if curh > maxy[row] {
			maxy[row] = curh
		}
		curcol++
		if curcol == g.cols {
			curcol = 0
			row++
		}
	}

	row = 0
	curcol = 0
	x := 0.0
	y := 0.0
	for _, child := range g.Children() {
		if child.GetName() == "__bbox" {
			continue
		}
		x += child.GetStyles().Padding
		y += child.GetStyles().Padding
		child.SetZeroAlignment(nodes.AlignmentTopLeft)
		child.SetPos(pixel.V(x, -y))
		x += maxx[curcol] - child.GetStyles().Padding
		y -= child.GetStyles().Padding
		curcol++
		if curcol == g.cols {
			curcol = 0
			y += maxy[row]
			x = 0
			row++
		}
	}
	b = pixel.R(0, 0, math.Round(sumSlice(maxx)), math.Round(sumSlice(maxy))).Norm()

	g.bbox.SetBounds(b)
	im := imdraw.New(nil)
	im.Color = styles.Border.Color
	im.Push(b.Min, pixel.V(b.Max.X, b.Min.Y), b.Max, pixel.V(b.Min.X, b.Max.Y))
	im.Polygon(styles.Border.Width)
	x = maxx[0]
	for i := 1; i < len(maxx); i++ {
		im.Push(pixel.V(x, b.Min.Y), pixel.V(x, b.Max.Y))
		im.Line(1)
		x += maxx[i]
	}
	y = b.Max.Y - maxy[0]
	for i := 1; i < len(maxy); i++ {
		im.Push(pixel.V(b.Min.X, y), pixel.V(b.Max.X, y))
		im.Line(1)
		y -= maxy[i]
	}
	g.bbox.Clear(styles.Background.Color)
	im.Draw(g.bbox.Canvas())

	g.SetBounds(b)
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
	if len(g.Children()) > 1 {
		g.recalcPositions()
	}
}

func (g *Grid) SetStyles(styles *nodes.Styles) {
	g.BaseNode.SetStyles(styles)
	g.bbox.SetStyles(styles)
}
