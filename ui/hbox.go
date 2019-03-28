package ui

import (
	"image/color"
	"pixelext/nodes"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
)

type HBox struct {
	nodes.BaseNode
	Padding     float64
	background  *nodes.BorderBox
	BorderWidth float64
	BorderColor color.RGBA
	oldbounds   pixel.Rect
}

func NewHBox(name string) *HBox {
	h := &HBox{
		BaseNode:    *nodes.NewBaseNode(name),
		Padding:     0,
		BorderWidth: 2,
		BorderColor: colornames.White,
		background:  nodes.NewBorderBox("bg", 1, 1),
	}
	h.Self = h
	return h
}

func (h *HBox) Init() {
	//h.background = nodes.NewBorderBox("bg", 1, 1) //pixelgl.NewCanvas(pixel.R(0, 0, 1, 1))
	h.background.SetPos(pixel.ZV)
	h.background.SetBorderColor(h.BorderColor)
	h.background.SetBorderWidth(h.BorderWidth)
	h.SetZIndex(-1)
	h.AddChild(h.background)
}

func (h *HBox) AddChild(child nodes.Node) {
	h.BaseNode.AddChild(child)

	//func (h *HBox) Update(dt float64) {
	xpos := h.Padding
	for _, child := range h.Children() {
		child.SetPos(pixel.V(xpos, h.Padding))
		xpos += child.GetBounds().W() + 2*h.Padding
	}
	h.SetBounds(pixel.R(0, 0, 0, 0))
	cb := h.GetContainerBounds()
	cb.Max = cb.Max.Add(pixel.V(h.Padding, 0))
	cb.Min = cb.Min.Sub(pixel.V(h.Padding, h.Padding))
	cb = cb.Moved(cb.Min.Scaled(-1))
	h.SetBounds(cb)
	//if cb.Size().Sub(h.oldbounds.Size()).Len() != 0 {
	h.background.SetBounds(cb)
	h.oldbounds = cb
	//}
	/*if cb.Size().Sub(h.oldbounds.Size()).Len() != 0 {
		if h.BorderWidth > 0 {
			h.background.SetBounds(cb)
			h.background.Clear(colornames.Black)
			im := imdraw.New(nil)
			im.Color = h.BorderColor
			im.Push(pixel.V(0, 0), pixel.V(0, cb.H()), pixel.V(cb.W(), cb.H()), pixel.V(cb.W(), 0), pixel.V(0, 0))
			im.Line(h.BorderWidth)
			im.Draw(h.background)
		}
		h.oldbounds = cb
	}*/
}

/*func (h *HBox) Draw(win *pixelgl.Window, mat pixel.Matrix) {
	if h.BorderWidth > 0 {
		bounds := h.background.Bounds()
		h.background.Draw(win, mat.Moved(pixel.V(bounds.W()/2, bounds.H()/2)))
	}
}*/
