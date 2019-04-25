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
	sub           *nodes.SubScene
	text          *Text
	cursor        *nodes.Canvas
	cursorpos     int
	textscroll    float64
	blinkduration float64
	totaltime     float64
	disabled      bool
	focused       bool
	cursorshown   bool
	maxlen        int
	onchange      func(string)
	onenter       func(string)
}

func NewInputBox(name, atlasname string, w, h float64) *InputBox {
	i := &InputBox{
		UIBase:        *NewUIBase(name),
		background:    nodes.NewBorderBox("bbox", w, h),
		sub:           nodes.NewSubScene("sub", w-8, h-6),
		text:          NewText("txt", atlasname),
		cursor:        nodes.NewCanvas("cursor", 2, h-6),
		content:       "",
		textscroll:    0,
		blinkduration: 0.5,
		cursorpos:     0,
		totaltime:     0,
		focused:       false,
		disabled:      false,
		maxlen:        255,
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
	//padding := i.GetStyles().Padding
	//size := i.Size()
	i.background.SetZIndex(-1)
	i.AddChild(i.background)
	i.sub.SetPos(pixel.V(0, 0))
	i.AddChild(i.sub)
	subroot := nodes.NewBaseNode("subroot")
	i.sub.SetRoot(subroot)
	i.text.SetAlignment(nodes.AlignmentCenterLeft)
	subroot.AddChild(i.text)
	i.cursor.SetZIndex(1)
	subroot.AddChild(i.cursor)
	i.cursor.Hide()
	i.recalc()
}

func (i *InputBox) recalc() {
	if i.cursorpos > len(i.content) {
		i.cursorpos = len(i.content)
	} else if i.cursorpos < 0 {
		i.cursorpos = 0
	}
	//padding := i.GetStyles().Padding
	size := i.Size()
	subwhalf := (size.X - 21) / 2
	i.text.Clear()
	i.text.Printf(i.content)
	leftcursorcontent := i.content[0:i.cursorpos]
	cursoroffset := i.text.Text().BoundsOf(leftcursorcontent).W() + 2
	cursorposx := -subwhalf - i.textscroll + cursoroffset
	if cursorposx < -subwhalf || cursorposx > subwhalf {
		i.textscroll = -subwhalf + cursoroffset
		if i.textscroll < 0 {
			i.textscroll = 0
		}
		cursorposx = -subwhalf - i.textscroll + cursoroffset
	}
	//fmt.Printf("cursorposx: %f, textscroll: %f\n", cursorposx, i.textscroll)
	i.text.SetPos(pixel.V(-subwhalf-i.textscroll, 0))
	i.cursor.SetPos(pixel.V(cursorposx, 0))
}

func (i *InputBox) Update(dt float64) {
	ev := nodes.Events()
	if i.focused {
		if ev.JustPressed(pixelgl.MouseButtonLeft) && !ev.Clicked(pixelgl.MouseButtonLeft, i) {
			i.focused = false
			i.cursor.Hide()
			i.cursorshown = false
			return
		}
		if ev.JustPressed(pixelgl.KeyBackspace) || ev.Repeated(pixelgl.KeyBackspace) {
			if i.cursorpos > 0 {
				if i.cursorpos == 1 {
					i.content = i.content[1:]
				} else if i.cursorpos < len(i.content) {
					i.content = i.content[:i.cursorpos-1] + i.content[i.cursorpos:]
				} else {
					i.content = i.content[:i.cursorpos-1]
				}
				i.cursorpos--
				i.recalc()
			}
		} else if ev.JustPressed(pixelgl.KeyLeft) || ev.Repeated(pixelgl.KeyLeft) {
			i.cursorpos--
			i.recalc()
		} else if ev.JustPressed(pixelgl.KeyRight) || ev.Repeated(pixelgl.KeyRight) {
			i.cursorpos++
			i.recalc()
		} else { // TODO: more cursor keys
			typed := ev.Typed()
			if len(typed) > 0 && len(i.content) < i.maxlen {
				if i.cursorpos == 0 {
					i.content = typed + i.content
				} else if i.cursorpos < len(i.content) {
					i.content = i.content[:i.cursorpos] + typed + i.content[i.cursorpos:]
				} else {
					i.content += typed
				}
				i.cursorpos += len(typed)
				if len(i.content) > i.maxlen {
					i.content = i.content[:i.maxlen]
				}
				i.recalc()
			}
		}
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
		if ev.Clicked(pixelgl.MouseButtonLeft, i) {
			if !i.focused {
				i.focused = true
				i.totaltime = 0
				i.cursor.Show()
				i.cursorshown = true
			}
			// TODO: calculate click location and set cursor accordingly
		}
	}
}
