package ui

import (
	"image/color"
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
	v.subscene = nodes.NewSubScene("scoll", w, 1)
	v.AddChild(v.subscene)
	v.Self = v
	v.UISelf = v
	return v
}

func (v *VScroll) Init() {
	v.root = nodes.NewBaseNode("root")
	v.subscene.SetRoot(v.root)
	v.scrollbar = nodes.NewBorderBox("scrollbar", 10, v.h) // TODO: event handler, new derived element necessary
	scrollbarstyle := v.scrollbar.GetStyles()
	scrollbarstyle.Border.Width = 0
	scrollbarstyle.Background.Color = color.RGBA{128, 128, 128, 128}
	v.scrollbar.SetPos(pixel.V(v.w/2-5, 0))
	v.scrollbar.Hide()
	v.AddChild(v.scrollbar)
}

func (v *VScroll) SetInner(inner UINode) {
	v.root.RemoveChildren()
	v.root.AddChild(inner)
	v.inner = inner
	v.w = inner.Size().X
	v.scroll = 0
	v.innerh = inner.Size().Y
	var newh float64
	if v.innerh > v.h {
		newh = v.h
	} else {
		newh = v.innerh
	}
	v.subscene.SetSize(pixel.V(v.w, newh))
	v.SetSize(pixel.V(v.w, newh))
	v.recalcScrollbar()
}

func (v *VScroll) recalcScrollbar() {
	if v.innerh > v.h {
		h := v.h * (v.h / v.innerh)
		v.scrollbar.SetSize(pixel.V(10, h))
		maxmove := v.h - h
		maxscroll := v.innerh - v.h
		factor := maxmove / maxscroll
		v.scrollbar.SetPos(pixel.V(v.w/2-5, v.h/2-h/2-v.scroll*factor))
		diff := v.innerh - v.h
		v.inner.SetPos(pixel.V(0, -diff/2+v.scroll))
	} else {
		v.scrollbar.SetSize(pixel.V(10, v.h))
		v.scrollbar.SetPos(pixel.V(v.w/2-5, 0))
	}
	nodes.SceneManager().Redraw()
}

func (v *VScroll) Update(dt float64) {
	if nodes.Events().MouseHovering(v) {
		v.scrollbar.Show()
		mousescroll := nodes.Events().MouseScroll()
		if mousescroll.Y != 0 {
			maxscroll := v.innerh - v.h
			v.scroll -= mousescroll.Y * 5
			if v.scroll < 0 {
				v.scroll = 0
			} else if v.scroll > maxscroll {
				v.scroll = maxscroll
			}
			v.recalcScrollbar()
		}
	} else {
		v.scrollbar.Hide()
	}
}

func (v *VScroll) SetScroll(scroll float64) {
	maxscroll := v.innerh - v.h
	if scroll > maxscroll {
		v.scroll = scroll
	} else if scroll < 0 {
		v.scroll = 0
	} else {
		v.scroll = scroll
	}
	v.recalcScrollbar()
}
