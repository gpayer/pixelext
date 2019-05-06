package ui

import (
	"github.com/gpayer/pixelext/nodes"

	"github.com/faiface/pixel"
)

type ButtonGroup struct {
	UIBase
	hbox     *HBox
	buttons  []*Button
	current  *Button
	h        float64
	onchange func(v string)
}

func NewButtonGroup(name string, h float64) *ButtonGroup {
	g := &ButtonGroup{
		UIBase:   *NewUIBase(name),
		hbox:     NewHBox("hbox"),
		h:        h,
		onchange: func(_ string) {},
	}
	g.Self = g
	g.UISelf = g
	g.AddChild(g.hbox)
	g.hbox.SetPos(pixel.ZV)
	return g
}

func (g *ButtonGroup) AddButton(caption, value string, w float64) {
	btn := NewButton("btn"+value, w, g.h, caption)
	if g.current != nil {
		btn.SetEnabled(true)
	}
	g.buttons = append(g.buttons, btn)
	if len(g.buttons) == 1 {
		g.current = btn
		btn.SetEnabled(false)
	}
	btn.OnClick(func() {
		if g.current != nil {
			g.current.SetEnabled(true)
		}
		g.current = btn
		btn.SetEnabled(false)
		g.onchange(value)
	})
	g.hbox.AddChild(btn)
	g.UISelf.SetSize(g.hbox.Size())
	nodes.SceneManager().Redraw()
}

func (g *ButtonGroup) SetAlignment(a nodes.Alignment) {
	g.UIBase.SetAlignment(a)
	g.hbox.SetAlignment(a)
	nodes.SceneManager().Redraw()
}

func (g *ButtonGroup) OnChange(fn func(v string)) {
	g.onchange = fn
}

func (g *ButtonGroup) SetStyles(styles *nodes.Styles) {
	g.UIBase.SetStyles(styles)
	g.hbox.SetStyles(styles)
	nodes.SceneManager().Redraw()
}
