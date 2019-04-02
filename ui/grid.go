package ui

import (
	"pixelext/nodes"

	"github.com/faiface/pixel"
)

type Grid struct {
	nodes.BaseNode
	bbox *nodes.BorderBox
	cols int
}

func NewGrid(name string, cols int) *Grid {
	g := &Grid{
		BaseNode: *nodes.NewBaseNode(name),
		cols:     cols,
	}
	g.Self = g
	g.SetZeroAlignment(nodes.AlignmentTopLeft)
	g.bbox = nodes.NewBorderBox("__bbox", g.GetStyles().Padding, g.GetStyles().Padding)
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
		if cb.W() > maxx[curcol] {
			maxx[curcol] = cb.W()
		}
		if cb.H() > maxy[row] {
			maxy[row] = cb.H()
		}
		curcol++
		if curcol == g.cols {
			curcol = 0
			row++
		}
	}

	row = 0
	curcol = 0
	x := styles.Padding
	y := -styles.Padding
	for _, child := range g.Children() {
		if child.GetName() == "__bbox" {
			continue
		}
		//child.SetZeroAlignment(nodes.AlignmentTopLeft)
		child.SetPos(pixel.V(x, -y))
		x += maxx[curcol] + styles.Padding
		curcol++
		if curcol == g.cols {
			curcol = 0
			y -= (maxy[row] + styles.Padding)
			x = styles.Padding
			row++
		}
	}
	b = pixel.R(0, 0, sumSlice(maxx)+float64(len(maxx)+1)*styles.Padding, sumSlice(maxy)+float64(len(maxy)+1)*styles.Padding).Norm()
	g.bbox.SetBounds(b)
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
