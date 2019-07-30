package ui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gpayer/pixelext/nodes"
	"golang.org/x/image/colornames"
)

const (
	MenuStateClosed int = iota
	MenuStateOpened
	MenuStateClosing
)

type Menu struct {
	VBox
	parent  *Menu
	menubar *MenuBar
	items   []MenuItem
	state   int
}

func NewMenu(name string, parent *Menu, menubar *MenuBar) *Menu {
	m := &Menu{
		VBox:    *NewVBox(name),
		parent:  parent,
		menubar: menubar,
		state:   MenuStateClosed,
	}
	m.Self = m
	m.UISelf = m
	m.SetHAlignment(nodes.HAlignmentLeft)
	m.SetAlignment(nodes.AlignmentTopLeft)
	styles := m.GetStyles()
	styles.Border.Width = 0
	styles.Background.Color = colornames.Darkgray
	m.SetStyles(styles)
	m.SetZIndex(9999)
	return m
}

func (m *Menu) SetItems(items []MenuItem) {
	btnstyle := nodes.DefaultStyles()
	btnstyle.Border.Width = 0
	btnstyle.Background.Color = colornames.Darkgray
	hoverstyle := btnstyle.Clone()
	hoverstyle.Background.Color = colornames.Lightgray
	hoverstyle.Text.Color = colornames.Black
	clickedstyle := hoverstyle.Clone()
	clickedstyle.Background.Color = colornames.White
	disabledstyle := btnstyle.Clone()
	disabledstyle.Text.Color = colornames.Darkgray

	m.items = items
	m.RemoveChildren()
	for i, _ := range m.items {
		item := &m.items[i]
		btn := NewButton("submenu_"+item.Value, 0, 0, item.Title)
		m.menubar.value2btn[item.Value] = btn
		btn.SetButtonStyles(ButtonEnabled, btnstyle)
		btn.SetButtonStyles(ButtonHover, hoverstyle)
		btn.SetButtonStyles(ButtonPressed, clickedstyle)
		btn.SetButtonStyles(ButtonDisabled, disabledstyle)
		btn.OnClick(func() {
			m.state = MenuStateClosing
			m.menubar.onselect(item.Value)
		})
		m.AddChild(btn)
	}
}

func (m *Menu) Activate(pos pixel.Vec) {
	if m.state != MenuStateClosed {
		return
	}
	m.Show()
	m.SetPos(pos)
	m.KeepOnScreen()
	nodes.Events().SetFocus(m)
	m.state = MenuStateOpened
}

func (m *Menu) Update(dt float64) {
	if m.state == MenuStateClosing {
		m.state = MenuStateClosed
		m.Hide()
	} else {
		ev := nodes.Events()
		if !ev.MouseHovering(m) && (ev.JustPressed(pixelgl.MouseButtonLeft) || ev.JustPressed(pixelgl.MouseButtonMiddle)) {
			m.state = MenuStateClosed
			m.Hide()
		}
	}
}
