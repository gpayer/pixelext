package ui

import (
	"math"
	"pixelext/nodes"

	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
)

type DropDown struct {
	UIBase
	hdropdown       float64
	cleared, opened bool
	current         string
	atlasname       string
	values          map[string]string
	value           *Text
	background      *nodes.BorderBox
	btn             *nodes.Canvas
	dropdown        *nodes.BorderBox
	vscroll         *VScroll
	list            *VBox
}

func NewDropDown(name, atlasname string, w, h, hdropdown float64) *DropDown {
	d := &DropDown{
		UIBase:    *NewUIBase(name),
		cleared:   true,
		opened:    false,
		hdropdown: hdropdown,
		current:   "",
		atlasname: atlasname,
		values:    make(map[string]string),
	}
	d.Self = d
	d.UISelf = d
	d.SetSize(pixel.V(w, h))
	return d
}

func (d *DropDown) Init() {
	size := d.Size()
	styles := d.GetStyles()
	d.background = nodes.NewBorderBox("background", size.X, size.Y)
	d.AddChild(d.background)
	d.value = NewText("value", d.atlasname)
	d.value.Printf("---")
	d.value.SetAlignment(nodes.AlignmentCenterLeft)
	d.value.SetPos(pixel.V(-size.X/2+styles.Padding, 0))
	d.value.SetZIndex(10)
	d.AddChild(d.value)

	d.btn = nodes.NewCanvas("btn", 15, size.Y)
	im := imdraw.New(nil)
	im.Push(pixel.V(4, math.Round(size.Y/2)+3), pixel.V(8, math.Round(size.Y/2)-3), pixel.V(12, math.Round(size.Y/2)+3))
	im.Line(1)
	im.Draw(d.btn.Canvas())
	d.btn.SetPos(pixel.V(size.X/2-8, 0))
	d.AddChild(d.btn)

	d.dropdown = nodes.NewBorderBox("dropdown", size.X, d.hdropdown+4)
	d.dropdown.SetPos(pixel.V(0, size.Y/2-d.hdropdown/2))
	d.dropdown.Hide()
	d.AddChild(d.dropdown)

	d.vscroll = NewVScroll("vscroll", size.X, d.hdropdown)
	d.dropdown.AddChild(d.vscroll)

	d.list = NewVBox("dropdownvbox")
	d.vscroll.SetInner(d.list)
}

func (d *DropDown) AddValue(text, value string) {
	btn := NewButton("btn", 0, 0, text)
	btn.OnClick(func() {

	})
	d.list.AddChild(btn)
	d.vscroll.SetInner(d.list)
}
