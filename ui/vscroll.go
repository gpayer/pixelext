package ui

import (
	"image/color"

	"github.com/gpayer/pixelext/nodes"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type scrollBar struct {
	nodes.BorderBox
	clicked bool
}

func newScrollbar(h float64) *scrollBar {
	s := &scrollBar{
		BorderBox: *nodes.NewBorderBox("scrollbar", 10, h),
	}
	s.Self = s
	return s
}

func (s *scrollBar) Update(dt float64) {
	if nodes.Events().Clicked(pixelgl.MouseButtonLeft, s) {
		s.clicked = true
	} else {
		s.clicked = false
	}
}

type VScroll struct {
	UIBase
	subscene             *nodes.SubScene
	root                 nodes.Node
	inner                UINode
	scrollbar            *scrollBar
	displayScrollbar     bool
	w, h, innerh, scroll float64
	scrolldrag           bool
	origclickpos         pixel.Vec
}

func NewVScroll(name string, w, h float64) *VScroll {
	v := &VScroll{
		UIBase:           *NewUIBase(name),
		w:                w,
		h:                h,
		scrolldrag:       false,
		displayScrollbar: false,
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
	v.scrollbar = newScrollbar(v.h)
	v.scrollbar.SetZIndex(10)
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
	v.innerh = inner.Size().Y
	switch v.halignment {
	case nodes.HAlignmentLeft:
		v.inner.SetAlignment(nodes.AlignmentCenterLeft)
	case nodes.HAlignmentRight:
		v.inner.SetAlignment(nodes.AlignmentCenterRight)
	case nodes.HAlignmentCenter:
		v.inner.SetAlignment(nodes.AlignmentCenter)
	}
	var newh float64
	if v.innerh > v.h {
		newh = v.h
	} else {
		newh = v.innerh
	}
	v.subscene.SetSize(pixel.V(v.w, newh))
	v.SetSize(pixel.V(v.w, newh))
	maxscroll := v.innerh - v.h
	if v.scroll > maxscroll {
		v.scroll = maxscroll
	} else if v.scroll < 0 {
		v.scroll = 0
	}
	v.recalcScrollbar()
}

func (v *VScroll) recalcScrollbar() {
	innerx := 0.0
	switch v.halignment {
	case nodes.HAlignmentLeft:
		innerx = -v.w / 2
	case nodes.HAlignmentRight:
		innerx = v.w / 2
	}

	if v.innerh > v.h {
		v.displayScrollbar = true
		h := v.h * (v.h / v.innerh)
		v.scrollbar.SetSize(pixel.V(10, h))
		maxmove := v.h - h
		maxscroll := v.innerh - v.h
		factor := maxmove / maxscroll
		v.scrollbar.SetPos(pixel.V(v.w/2-5, v.h/2-h/2-v.scroll*factor))
		diff := v.innerh - v.h
		v.inner.SetPos(pixel.V(innerx, -diff/2+v.scroll))
	} else {
		v.displayScrollbar = false
		v.scrollbar.Hide()
		v.inner.SetPos(pixel.V(innerx, 0))
		scrollbarh := v.h
		if scrollbarh > v.Size().Y {
			scrollbarh = v.Size().Y
		}
		v.scrollbar.SetSize(pixel.V(10, scrollbarh))
		v.scrollbar.SetPos(pixel.V(v.w/2-5, 0))
	}
	nodes.SceneManager().Redraw()
}

func (v *VScroll) Update(dt float64) {
	ev := nodes.Events()
	if v.scrolldrag {
		if ev.JustReleased(pixelgl.MouseButtonLeft) {
			v.scrolldrag = false
		} else {
			diff := v.origclickpos.Y - ev.MousePosition().Y
			if diff != 0 {
				h := v.h * (v.h / v.innerh)
				maxmove := v.h - h
				maxscroll := v.innerh - v.h
				factor := maxscroll / maxmove
				v.scroll += diff * factor
				v.SetScroll(v.scroll)
				v.origclickpos = ev.MousePosition()
			}
		}
	}
	if ev.MouseHovering(v) {
		if v.displayScrollbar {
			v.scrollbar.Show()
		}
		mousescroll := ev.MouseScroll()
		if mousescroll.Y != 0 {
			v.scroll -= mousescroll.Y * 5
			v.SetScroll(v.scroll)
		}
		if !v.scrolldrag && v.scrollbar.clicked {
			v.scrolldrag = true
			v.origclickpos = ev.MousePosition()
		}
	} else if !v.scrolldrag {
		v.scrollbar.Hide()
	}
}

func (v *VScroll) SetScroll(scroll float64) {
	maxscroll := v.innerh - v.h
	if scroll > maxscroll {
		v.scroll = maxscroll
	} else if scroll < 0 {
		v.scroll = 0
	} else {
		v.scroll = scroll
	}
	v.recalcScrollbar()
}

func (v *VScroll) Scroll() float64 {
	return v.scroll
}
