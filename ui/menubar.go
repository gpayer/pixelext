package ui

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/gpayer/pixelext/nodes"
	"golang.org/x/image/colornames"
)

type MenuItem struct {
	Title   string
	Value   string
	SubMenu []MenuItem
	menu    *Menu
}

type MenuBar struct {
	HBox
	items     []MenuItem
	onselect  func(v string)
	value2btn map[string]*Button
}

func NewMenuBar(name string) *MenuBar {
	m := &MenuBar{
		HBox:      *NewHBox(name),
		onselect:  func(_ string) {},
		value2btn: make(map[string]*Button),
	}
	m.Self = m
	m.UISelf = m
	styles := m.GetStyles()
	styles.Border.Width = 0
	styles.Background.Color = colornames.Darkgray
	styles.Padding = 2
	m.SetStyles(styles)
	return m
}

func (m *MenuBar) SetItems(items []MenuItem) {
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

	m.RemoveChildren()
	m.items = items
	for i, _ := range items {
		item := &m.items[i]
		btn := NewButton("menubtn", 0, 0, item.Title)
		m.value2btn[item.Value] = btn
		btn.SetButtonStyles(ButtonEnabled, btnstyle)
		btn.SetButtonStyles(ButtonHover, hoverstyle)
		btn.SetButtonStyles(ButtonPressed, clickedstyle)
		btn.SetButtonStyles(ButtonDisabled, disabledstyle)
		if len(item.SubMenu) > 0 {
			menu := NewMenu("submenu_"+item.Value, nil, m)
			menu.SetItems(item.SubMenu)
			menu.Hide()
			item.menu = menu
			nodes.SceneManager().Root().AddChild(menu)
		}
		btn.OnClick(func() {
			m.onselect(item.Value)
			if item.menu != nil {
				y := btn.GetPos().Y - math.Round(m.Size().Y/2) + 2
				x := btn.GetPos().X - math.Round(btn.Size().X/2)
				globpos := btn.LocalToGlobalPos(pixel.V(x, y))
				item.menu.Activate(globpos)
			}
		})
		m.AddChild(btn)
	}
}

func (m *MenuBar) OnSelect(fn func(v string)) {
	m.onselect = fn
}

func (m *MenuBar) SetMenuItemEnabled(value string, enabled bool) {
	btn, ok := m.value2btn[value]
	if ok {
		btn.SetEnabled(enabled)
	}
}

func (m *MenuBar) SetWidth(w float64) {
	m.SetFixedSize(pixel.V(w, 0))
}
