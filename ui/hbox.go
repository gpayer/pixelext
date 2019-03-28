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
		background:  nodes.NewBorderBox("__bg", 1, 1),
	}
	h.Self = h
	return h
}

func (h *HBox) Init() {
	h.background.SetPos(pixel.ZV)
	h.background.SetBorderColor(h.BorderColor)
	h.background.SetBorderWidth(h.BorderWidth)
	h.SetZIndex(-1)
	h.AddChild(h.background)
}

func (h *HBox) AddChild(child nodes.Node) {
	h.BaseNode.AddChild(child)

	xpos := h.Padding
	maxy := 0.0
	for _, child := range h.Children() {
		if child.GetName() != "__bg" {
			child.SetPos(pixel.V(xpos, h.Padding))
			xpos += child.GetBounds().W() + 2*h.Padding
			if child.GetBounds().H() > maxy {
				maxy = child.GetBounds().H()
			}
		}
	}
	h.SetBounds(pixel.R(0, 0, 0, 0))
	cb := pixel.R(-h.Padding, -h.Padding, xpos-h.Padding, maxy+h.Padding)
	h.SetBounds(cb)
	h.background.SetBounds(cb)
	h.oldbounds = cb
}
