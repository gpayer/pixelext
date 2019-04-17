package ui

import (
	"pixelext/nodes"

	"github.com/faiface/pixel"
)

type VScroll struct {
	UIBase
	subscene             *nodes.SubScene
	root                 nodes.Node
	inner                UINode
	scrollbar            *nodes.BorderBox
	w, h, innerh, scroll float64
}

func NewVScroll(name string, w, h float64) *VScroll {
	v := &VScroll{
		UIBase: *NewUIBase(name),
		w:      w,
		h:      h,
	}
	if w == 0 {
		w = 1
	}
	v.subscene = nodes.NewSubScene("scoll", w, h)
	v.Self = v
	v.UISelf = v
	return v
}

func (v *VScroll) Init() {
	v.root = nodes.NewBaseNode("root")
	v.subscene.SetRoot(v.root)
	v.scrollbar = nodes.NewBorderBox("scrollbar", 10, v.h) // TODO: event handler, new derived element necessary
	v.scrollbar.SetPos(pixel.V(v.w/2-5, 0))
	v.scrollbar.Hide()
	v.AddChild(v.scrollbar)
}

func (v *VScroll) SetInner(inner UINode) {
	// TODO: remove previous child
	v.root.AddChild(inner)
	v.inner = inner
	v.w = inner.Size().X
	v.scroll = 0
	v.innerh = inner.Size().Y
	v.recalcScrollbar()
	nodes.SceneManager().Redraw()
}

func (v *VScroll) recalcScrollbar() {
	if v.innerh > v.h {
		//maxscroll := v.innerh - v.h
		v.scrollbar.SetSize(pixel.V(10, v.h*(v.h/v.innerh)))
		// TODO: calc scrollbar position
	} else {
		v.scrollbar.SetSize(pixel.V(10, v.h))
		v.scrollbar.SetPos(pixel.V(v.w/2-5, 0))
	}
}
