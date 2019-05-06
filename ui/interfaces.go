package ui

import (
	"github.com/gpayer/pixelext/nodes"

	"github.com/faiface/pixel"
)

type UINode interface {
	nodes.Node
	Size() pixel.Vec
	SetAlignment(a nodes.Alignment)
}
