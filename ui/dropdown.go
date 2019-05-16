package ui

import (
	"image/color"
	"math"

	"github.com/gpayer/pixelext/nodes"

	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
)

type dropDownState int

const (
	dropDownClosed dropDownState = iota
	dropDownOpening
	dropDownOpened
)

type DropDown struct {
	UIBase
	hdropdown  float64
	state      dropDownState
	current    string
	atlasname  string
	values     map[string]string
	value      *Text
	background *nodes.BorderBox
	sub        *nodes.SubScene
	btn        *nodes.Canvas
	dropdown   *nodes.BorderBox
	vscroll    *VScroll
	list       *VBox
	onchange   func(value string, text string)
}

func NewDropDown(name, atlasname string, w, h, hdropdown float64) *DropDown {
	d := &DropDown{
		UIBase:    *NewUIBase(name),
		hdropdown: hdropdown,
		current:   "",
		atlasname: atlasname,
		values:    make(map[string]string),
		onchange:  func(_ string, _ string) {},
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
	d.background.SetZIndex(-1)
	d.AddChild(d.background)

	subscenewidth := size.X - 15 - styles.Padding
	d.sub = nodes.NewSubScene("subscene", subscenewidth, size.Y-4)
	d.sub.SetPos(pixel.V(-5, 0))
	d.AddChild(d.sub)

	subroot := nodes.NewBaseNode("subroot")
	d.sub.SetRoot(subroot)

	d.value = NewText("value", d.atlasname)
	d.value.Printf("---")
	d.value.SetAlignment(nodes.AlignmentCenterLeft)
	d.value.SetPos(pixel.V(-subscenewidth/2+2, 0))
	d.value.SetZIndex(10)
	subroot.AddChild(d.value)

	d.btn = nodes.NewCanvas("btn", 15, size.Y)
	im := imdraw.New(nil)
	im.Push(pixel.V(4, math.Round(size.Y/2)+3), pixel.V(8, math.Round(size.Y/2)-3), pixel.V(12, math.Round(size.Y/2)+3))
	im.Line(1)
	im.Draw(d.btn.Canvas())
	d.btn.SetPos(pixel.V(size.X/2-8, 0))
	d.AddChild(d.btn)

	d.dropdown = nodes.NewBorderBox("dropdown", size.X, d.hdropdown+4)
	d.dropdown.Hide()
	d.AddChild(d.dropdown)

	d.vscroll = NewVScroll("vscroll", size.X, d.hdropdown)
	d.dropdown.AddChild(d.vscroll)

	d.list = NewVBox("dropdownvbox")
	listStyles := d.list.GetStyles()
	listStyles.Border.Width = 0
	listStyles.Padding = 2

	d.createDropdown()
}

func (d *DropDown) initValue(text, value string) {
	btn := NewButton("btn", 0, 0, text)
	btn.OnClick(func() {
		d.onchange(value, text)
		d.value.Clear()
		d.value.Printf("%s", text)
		d.state = dropDownClosed
		d.vscroll.SetScroll(0)
		d.current = value
	})

	baseStyle := nodes.DefaultStyles()
	baseStyle.Text.Atlas = d.atlasname
	baseStyle.Border.Width = 0
	baseStyle.Padding = 0

	enabledStyle := baseStyle.Clone()
	enabledStyle.Background.Color = colornames.Black
	btn.SetButtonStyles(ButtonEnabled, enabledStyle)

	hoverStyle := baseStyle.Clone()
	hoverStyle.Background.Color = color.RGBA{64, 64, 64, 255}
	btn.SetButtonStyles(ButtonHover, hoverStyle)

	pressedStyle := baseStyle.Clone()
	pressedStyle.Background.Color = colornames.White
	pressedStyle.Text.Color = colornames.Black
	btn.SetButtonStyles(ButtonPressed, pressedStyle)
	btn.SetStyles(baseStyle)

	d.list.AddChild(btn)
}

func (d *DropDown) SetSize(size pixel.Vec) {
	d.UIBase.SetSize(size)
	if d.Initialized() {
		styles := d.GetStyles()
		d.background.SetSize(size)
		subscenewidth := size.X - 15 - styles.Padding
		d.sub.SetSize(pixel.V(subscenewidth, size.Y-4))
		d.value.SetPos(pixel.V(-subscenewidth/2+2, 0))
	}
}

func (d *DropDown) createDropdown() {
	d.list.RemoveChildren()
	for val, txt := range d.values {
		d.initValue(txt, val)
	}
	d.vscroll.SetInner(d.list)
	size := d.vscroll.Size().Add(pixel.V(2, 2))
	d.dropdown.SetSize(size)
	d.dropdown.SetPos(pixel.V(0, -d.Size().Y/2-size.Y/2))
}

func (d *DropDown) AddValue(text, value string) {
	d.values[value] = text
	if d.Initialized() {
		d.createDropdown()
	}
}

func (d *DropDown) SetValue(value string) {
	text, ok := d.values[value]
	if ok {
		d.value.Clear()
		d.value.Printf(text)
		d.current = value
	}
}

func (d *DropDown) RemoveValue(value string) {
	delete(d.values, value)
	if d.current == value {
		d.value.Clear()
		d.value.Printf("---")
		d.current = ""
	}
	if d.Initialized() {
		d.createDropdown()
	}
}

func (d *DropDown) ChangeValue(value, text string) {
	_, ok := d.values[value]
	if ok {
		d.values[value] = text
		if d.current == value {
			d.value.Clear()
			d.value.Printf(text)
		}
		if d.Initialized() {
			d.createDropdown()
		}
	} else {
		d.AddValue(text, value)
	}
}

func (d *DropDown) Clear() {
	d.values = make(map[string]string, 0)
	d.value.Clear()
	d.value.Printf("---")
	d.current = ""
	if d.Initialized() {
		d.createDropdown()
	}
}

func (d *DropDown) OnChange(fn func(string, string)) {
	d.onchange = fn
}

func (d *DropDown) Update(dt float64) {
	if nodes.Events().Clicked(pixelgl.MouseButtonLeft, d) {
		if d.state == dropDownClosed {
			d.dropdown.Show()
			d.state = dropDownOpening
		}
	} else if nodes.Events().JustReleased(pixelgl.MouseButtonLeft) {
		if d.state == dropDownOpening {
			d.state = dropDownOpened
		} else if d.state == dropDownOpened {
			d.vscroll.SetScroll(0)
			d.state = dropDownClosed
		}
	} else if d.state == dropDownClosed {
		d.dropdown.Hide()
	}
}
