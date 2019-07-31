package ui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gpayer/pixelext/nodes"
	"golang.org/x/image/colornames"
)

const (
	menuStateInit int = iota
	menuStateClosed
	menuStateOpened
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
		state:   menuStateInit,
	}
	m.Self = m
	m.UISelf = m
	m.SetHAlignment(nodes.HAlignmentLeft)
	m.SetAlignment(nodes.AlignmentTopLeft)
	styles := m.GetStyles()
	styles.Border.Width = 0
	styles.Background.Color = colornames.Lightgray
	m.SetStyles(styles)
	m.SetZIndex(9999)
	return m
}

func (m *Menu) SetItems(items []MenuItem) {
	btnstyle := nodes.DefaultStyles()
	btnstyle.Border.Width = 0
	btnstyle.Background.Color = colornames.Lightgray
	btnstyle.Text.Color = colornames.Black
	hoverstyle := btnstyle.Clone()
	hoverstyle.Background.Color = colornames.Darkgray
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
			m.state = menuStateClosed
			m.menubar.onselect(item.Value)
		})
		m.AddChild(btn)
	}
}

func (m *Menu) Activate(pos pixel.Vec) {
	if m.state == menuStateInit {
		m.state = menuStateClosed
		nodes.SceneManager().Root().AddChild(m)
	}
	m.Show()
	m.SetPos(pos)
	//m.KeepOnScreen()
	nodes.Events().SetFocus(m)
	m.state = menuStateOpened
}

func (m *Menu) Update(dt float64) {
	if nodes.Events().JustReleased(pixelgl.MouseButtonLeft) {
		if m.state == menuStateOpened && !nodes.Events().MouseHovering(m) {
			m.state = menuStateClosed
		}
	} else if m.state == menuStateClosed {
		// TODO: hide sub menus
		nodes.Events().SetFocus(nil)
		m.Hide()
	}
}
