package ui

import (
	"pixelext/nodes"

	"github.com/faiface/pixel"
)

type UIBase struct {
	nodes.BaseNode
	alignment nodes.Alignment
	UISelf    UINode
	size      pixel.Vec
	origpos   pixel.Vec
}

func NewUIBase(name string) *UIBase {
	ui := &UIBase{
		BaseNode:  *nodes.NewBaseNode(name),
		alignment: nodes.AlignmentCenter,
	}
	ui.Self = ui
	ui.UISelf = ui
	return ui
}

func (ui *UIBase) SetOrigin(origin pixel.Vec) {
	ui.BaseNode.SetOrigin(origin)
}

func (ui *UIBase) SetPos(pos pixel.Vec) {
	ui.origpos = pos
	whalf := ui.size.X / 2
	hhalf := ui.size.Y / 2
	switch ui.alignment {
	case nodes.AlignmentBottomLeft:
		pos = pos.Add(pixel.V(whalf, hhalf))
	case nodes.AlignmentCenterLeft:
		pos = pos.Add(pixel.V(whalf, 0))
	case nodes.AlignmentTopLeft:
		pos = pos.Add(pixel.V(whalf, -hhalf))
	case nodes.AlignmentBottomCenter:
		pos = pos.Add(pixel.V(0, hhalf))
	case nodes.AlignmentCenter:
	case nodes.AlignmentTopCenter:
		pos = pos.Add(pixel.V(0, -hhalf))
	case nodes.AlignmentBottomRight:
		pos = pos.Add(pixel.V(-whalf, hhalf))
	case nodes.AlignmentCenterRight:
		pos = pos.Add(pixel.V(-whalf, 0))
	case nodes.AlignmentTopRight:
		pos = pos.Add(pixel.V(-whalf, -hhalf))
	default:
	}
	ui.BaseNode.SetPos(pos)
}

// Contains returns whether the given point (in local coordinates) lies within the
// boundaries of this UI element
func (ui *UIBase) Contains(point pixel.Vec) bool {
	// TODO
	return false
}

func (ui *UIBase) SetSize(size pixel.Vec) {
	ui.size = size
	ui.UISelf.SetPos(ui.origpos)
}

func (ui *UIBase) Size() pixel.Vec {
	return ui.size
}

func (ui *UIBase) SetAlignment(a nodes.Alignment) {
	ui.alignment = a
	ui.UISelf.SetPos(ui.origpos)
}
