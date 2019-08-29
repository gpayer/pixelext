package ui

import (
	"math"

	"github.com/gpayer/pixelext/nodes"

	"github.com/faiface/pixel"
)

type HBox struct {
	UIBase
	background *nodes.BorderBox
	w, h       float64
}

func NewHBox(name string) *HBox {
	h := &HBox{
		UIBase:     *NewUIBase(name),
		background: nodes.NewBorderBox("__bg", 1, 1),
	}
	h.Self = h
	h.UISelf = h
	h.background.SetLocked(true)
	return h
}

func (h *HBox) Init() {
	h.background.SetPos(pixel.ZV)
	h.background.SetStyles(h.GetStyles())
	h.background.SetZIndex(-9999)
	h.AddChild(h.background)
}

func (h *HBox) recalcPositions() {
	padding := h.GetStyles().Padding
	borderWidth := h.GetStyles().Border.Width
	xpos := padding
	maxy := 0.0
	for _, child := range h.Children() {
		if child == nil {
			continue
		}
		uichild, ok := child.(UINode)
		if ok && child.GetName() != "__bg" {
			childbounds := uichild.Size()
			xpos += childbounds.X + 2*padding
			if childbounds.Y > maxy {
				maxy = childbounds.Y
			}
		}
	}
	size := pixel.V(math.Round(xpos-padding+2*borderWidth), math.Round(maxy+2*padding+2*borderWidth))
	if h.w > 0 {
		size.X = h.w
	}
	if h.h > 0 {
		size.Y = h.h
	}
	h.SetSize(size)
	h.background.SetSize(size)

	xpos = -size.X/2 + padding + borderWidth
	for _, child := range h.Children() {
		if child == nil {
			continue
		}
		uichild, ok := child.(UINode)
		if ok && child.GetName() != "__bg" {
			uichild.SetAlignment(nodes.AlignmentCenterLeft)
			uichild.SetPos(pixel.V(xpos, 0))
			childbounds := uichild.Size()
			xpos += childbounds.X + 2*padding
		}
	}
	nodes.SceneManager().Redraw()
}

func (h *HBox) AddChild(child nodes.Node) {
	h.UIBase.AddChild(child)
	h.recalcPositions()
	child.SetPosLocked(true)
}

func (h *HBox) SetStyles(styles *nodes.Styles) {
	h.UIBase.SetStyles(styles)
	h.background.SetStyles(styles)
	h.recalcPositions()
}

func (h *HBox) RemoveChild(child nodes.Node) {
	h.UIBase.RemoveChild(child)
	h.recalcPositions()
	child.SetPosLocked(false)
}

func (h *HBox) RemoveChildren() {
	for _, c := range h.Children() {
		c.SetPosLocked(false)
	}
	h.UIBase.RemoveChildren()
	h.recalcPositions()
}

func (h *HBox) ChildChanged() {
	h.recalcPositions()
}

func (h *HBox) SetFixedSize(size pixel.Vec) {
	h.w = size.X
	h.h = size.Y
	h.recalcPositions()
}
