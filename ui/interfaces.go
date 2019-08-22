package ui

import (
	"github.com/gpayer/pixelext/nodes"

	"github.com/faiface/pixel"
)

type UINode interface {
	nodes.Node
	GetOrigPos() pixel.Vec
	Size() pixel.Vec
	Alignment() nodes.Alignment
	SetAlignment(a nodes.Alignment)
	OverrideStyles(styles *nodes.Styles)
	HAlignment() nodes.HorizontalAlignment
	SetHAlignment(h nodes.HorizontalAlignment)
	VAlignment() nodes.VerticalAlignment
	SetVAlignment(v nodes.VerticalAlignment)
	IsStylesOverridden() bool
	SetStylesOverridden(overridden bool)
}
