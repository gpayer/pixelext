package ui

import (
	"pixelext/nodes"

	"github.com/faiface/pixel"
)

type HBox struct {
	nodes.BaseNode
	background *nodes.BorderBox
}

func NewHBox(name string) *HBox {
	h := &HBox{
		BaseNode:   *nodes.NewBaseNode(name),
		background: nodes.NewBorderBox("__bg", 1, 1),
	}
	h.Self = h
	return h
}

func (h *HBox) Init() {
	h.background.SetPos(pixel.ZV)
	h.background.SetStyles(h.GetStyles())
	h.SetZIndex(-1)
	h.AddChild(h.background)
}

func (h *HBox) recalcPositions() {
	padding := h.GetStyles().Padding
	xpos := padding
	maxy := 0.0
	for _, child := range h.Children() {
		if child.GetName() != "__bg" {
			child.SetPos(pixel.V(xpos, padding))
			xpos += child.GetBounds().W() + 2*padding
			if child.GetBounds().H() > maxy {
				maxy = child.GetBounds().H()
			}
		}
	}
	h.SetBounds(pixel.R(0, 0, 0, 0))
	cb := pixel.R(-padding, -padding, xpos-padding, maxy+padding)
	h.SetBounds(cb)
	h.background.SetBounds(cb)
}

func (h *HBox) AddChild(child nodes.Node) {
	h.BaseNode.AddChild(child)
	h.recalcPositions()
}

func (h *HBox) SetStyles(styles *nodes.Styles) {
	h.BaseNode.SetStyles(styles)
	h.background.SetStyles(styles)
	h.recalcPositions()
}
