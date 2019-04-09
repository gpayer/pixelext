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
	ui.BaseNode.SetPos(pos)
}

func (ui *UIBase) Contains(point pixel.Vec) bool {
	return false
}

func (ui *UIBase) SetSize(size pixel.Vec) {
	ui.size = size
}

func (ui *UIBase) Size() pixel.Vec {
	return ui.size
}

func (ui *UIBase) SetAlignment(a nodes.Alignment) {
	oldpos := ui.UISelf.GetPos()
	ui.alignment = a
	ui.UISelf.SetPos(oldpos)
}
