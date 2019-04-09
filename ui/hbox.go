package ui

import (
	"pixelext/nodes"

	"github.com/faiface/pixel"
)

type HBox struct {
	UIBase
	background *nodes.BorderBox
}

func NewHBox(name string) *HBox {
	h := &HBox{
		UIBase:     *NewUIBase(name),
		background: nodes.NewBorderBox("__bg", 1, 1),
	}
	h.Self = h
	h.UISelf = h
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
		uichild, ok := child.(UINode)
		if ok && child.GetName() != "__bg" {
			uichild.SetPos(pixel.V(xpos, padding))
			childbounds := uichild.Size()
			xpos += childbounds.X + 2*padding
			if childbounds.Y > maxy {
				maxy = childbounds.Y
			}
		}
	}
	size := pixel.V(xpos-2*padding, maxy+2*padding)
	h.SetSize(size)
	h.background.SetSize(size)
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
