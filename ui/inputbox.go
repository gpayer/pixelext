package ui

import (
	"pixelext/nodes"

	"github.com/faiface/pixel/pixelgl"

	"github.com/faiface/pixel"
)

type InputBox struct {
	UIBase
	background    *nodes.BorderBox
	content       string
	text          *Text
	cursor        *nodes.Canvas
	cursorpos     int
	blinkduration float64
	totaltime     float64
	disabled      bool
	focused       bool
	cursorshown   bool
	onchange      func(string)
	onenter       func(string)
}

func NewInputBox(name, atlasname string, w, h float64) *InputBox {
	i := &InputBox{
		UIBase:        *NewUIBase(name),
		background:    nodes.NewBorderBox("bbox", w, h),
		text:          NewText("txt", atlasname),
		cursor:        nodes.NewCanvas("cursor", 1, h-6),
		content:       "",
		blinkduration: 0.5,
		cursorpos:     0,
		totaltime:     0,
		focused:       false,
		disabled:      false,
		onchange:      func(_ string) {},
		onenter:       func(_ string) {},
	}
	i.Self = i
	i.UISelf = i
	i.cursor.Clear(i.GetStyles().Text.Color)
	i.SetSize(pixel.V(w, h))
	return i
}

func (i *InputBox) Init() {
	padding := i.GetStyles().Padding
	size := i.Size()
	i.background.SetZIndex(-1)
	i.AddChild(i.background)
	i.text.SetPos(pixel.V(-size.X/2+padding, 0))
	i.AddChild(i.text)
	i.cursor.SetZIndex(1)
	i.SetPos(pixel.V(-size.X/2+padding+i.text.Size().X+2, 0))
	i.AddChild(i.cursor)
	i.cursor.Hide()
}

func (i *InputBox) Update(dt float64) {
	if i.focused {
		if nodes.Events().JustPressed(pixelgl.MouseButtonLeft) && !nodes.Events().Clicked(pixelgl.MouseButtonLeft, i) {
			i.focused = false
			i.cursor.Hide()
			i.cursorshown = false
			return
		}
		// TODO: key input
		i.totaltime += dt
		if i.totaltime > 0.5 {
			i.totaltime = 0
			if i.cursorshown {
				i.cursor.Show()
			} else {
				i.cursor.Hide()
			}
			i.cursorshown = !i.cursorshown
		}
	} else {
		if nodes.Events().Clicked(pixelgl.MouseButtonLeft, i) {
			i.focused = true
			i.totaltime = 0
			i.cursor.Show()
			i.cursorshown = true
			// TODO: calculate click location and set cursor accordingly
		}
	}
}
