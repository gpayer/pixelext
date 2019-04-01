package ui

import "pixelext/nodes"

type ButtonGroup struct {
	nodes.BaseNode
	hbox     *HBox
	buttons  []*Button
	current  *Button
	h        float64
	onchange func(v string)
}

func NewButtonGroup(name string, h float64) *ButtonGroup {
	g := &ButtonGroup{
		BaseNode: *nodes.NewBaseNode(name),
		hbox:     NewHBox("hbox"),
		h:        h,
		onchange: func(_ string) {},
	}
	g.Self = g
	g.AddChild(g.hbox)
	return g
}

func (g *ButtonGroup) AddButton(caption, value string, w float64) {
	btn := NewButton("", w, g.h, caption)
	if g.current != nil {
		btn.SetEnabled(false)
	}
	g.buttons = append(g.buttons, btn)
	btn.OnClick(func() {
		if g.current != nil {
			g.current.SetEnabled(true)
		}
		g.current = btn
		btn.SetEnabled(false)
		g.onchange(value)
	})
	g.hbox.AddChild(btn)
}

func (g *ButtonGroup) SetZeroAlignment(a nodes.Alignment) {
	g.BaseNode.SetZeroAlignment(a)
	g.hbox.SetZeroAlignment(a)
}

func (g *ButtonGroup) OnChange(fn func(v string)) {
	g.onchange = fn
}
