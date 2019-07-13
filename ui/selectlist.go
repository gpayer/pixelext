package ui

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/gpayer/pixelext/nodes"
	"golang.org/x/image/colornames"
)

type selectListEntry struct {
	text, value string
	selected    bool
	btn         *Button
}

type SelectList struct {
	UIBase
	background                          *nodes.BorderBox
	vscroll                             *VScroll
	list                                *VBox
	baseStyle, enabledStyle, hoverStyle *nodes.Styles
	pressedStyle, selectedStyle         *nodes.Styles
	entries                             []*selectListEntry
	value2entry                         map[string]*selectListEntry
	onselect                            func(v string)
	onunselect                          func(v string)
	multiselect                         bool
}

func NewSelectList(name string, w, h float64, multiselect bool) *SelectList {
	defaultStyles := nodes.DefaultStyles()

	s := &SelectList{
		UIBase:      *NewUIBase(name),
		background:  nodes.NewBorderBox("background", w, h),
		vscroll:     NewVScroll("vscroll", w-2*defaultStyles.Padding-2*defaultStyles.Border.Width, h-2*defaultStyles.Border.Width),
		list:        NewVBox("list"),
		value2entry: make(map[string]*selectListEntry),
		onselect:    func(_ string) {},
		onunselect:  func(_ string) {},
		multiselect: multiselect,
	}
	s.Self = s
	s.UISelf = s
	s.UIBase.SetSize(pixel.V(w, h))

	s.baseStyle = nodes.DefaultStyles()
	s.baseStyle.Border.Width = 0
	s.baseStyle.Padding = 0

	s.enabledStyle = s.baseStyle.Clone()
	s.enabledStyle.Background.Color = colornames.Black

	s.selectedStyle = s.baseStyle.Clone()
	s.selectedStyle.Background.Color = colornames.Gray
	s.selectedStyle.Text.Color = colornames.Black

	s.hoverStyle = s.baseStyle.Clone()
	s.hoverStyle.Background.Color = color.RGBA{64, 64, 64, 255}

	s.pressedStyle = s.baseStyle.Clone()
	s.pressedStyle.Background.Color = colornames.White
	s.pressedStyle.Text.Color = colornames.Black
	return s
}

func (s *SelectList) Init() {
	s.background.SetLocked(true)
	s.AddChild(s.background)
	s.vscroll.SetLocked(true)
	s.vscroll.SetHAlignment(nodes.HAlignmentLeft)
	s.AddChild(s.vscroll)
	st := s.list.GetStyles()
	st.Border.Width = 0
	st.Padding = 2
	s.list.OverrideStyles(st)
	s.list.SetHAlignment(nodes.HAlignmentLeft)
	s.vscroll.SetInner(s.list)
}

func (s *SelectList) OnSelect(onselect func(v string)) {
	s.onselect = onselect
}

func (s *SelectList) OnUnselect(onunselect func(v string)) {
	s.onunselect = onunselect
}

func (s *SelectList) AddEntry(text string, value string) {
	btnW := s.Size().X - 2*s.GetStyles().Padding
	btn := NewButton("btn", btnW, 0, text)
	btn.SetHAlignment(nodes.HAlignmentLeft)
	btn.OnClick(func() {
		e := s.value2entry[value]
		if e.selected {
			if s.multiselect {
				s.onunselect(value)
				e.btn.SetButtonStyles(ButtonEnabled, s.enabledStyle)
			}
		} else {
			if !s.multiselect {
				s.UnselectAll()
			}
			s.onselect(value)
			e.btn.SetButtonStyles(ButtonEnabled, s.selectedStyle)
		}
		e.selected = !e.selected
	})

	btn.SetButtonStyles(ButtonEnabled, s.enabledStyle)

	btn.SetButtonStyles(ButtonHover, s.hoverStyle)

	btn.SetButtonStyles(ButtonPressed, s.pressedStyle)
	btn.SetStyles(s.baseStyle)

	s.list.AddChild(btn)
	s.vscroll.SetInner(s.list)
	newentry := &selectListEntry{
		text:     text,
		value:    value,
		selected: false,
		btn:      btn,
	}
	s.entries = append(s.entries, newentry)
	s.value2entry[value] = newentry
}

func (s *SelectList) Clear() {
	s.entries = s.entries[:0]
	s.value2entry = make(map[string]*selectListEntry)
	s.list.RemoveChildren()
	s.vscroll.SetInner(s.list)
}

func (s *SelectList) UnselectAll() {
	for _, e := range s.entries {
		if e.selected {
			e.selected = false
			e.btn.SetButtonStyles(ButtonEnabled, s.enabledStyle)
		}
	}
}

func (s *SelectList) Selected() []string {
	var sel []string
	for _, e := range s.entries {
		sel = append(sel, e.value)
	}
	return sel
}

func (s *SelectList) SetSize(size pixel.Vec) {
	s.UIBase.SetSize(size)
	s.background.SetSize(size)
	st := s.GetStyles()
	s.vscroll.SetSize(size.Sub(pixel.V(2*st.Padding+2*st.Border.Width, 2*st.Border.Width)))
}

func (s *SelectList) SetStyles(styles *nodes.Styles) {
	s.UIBase.SetStyles(styles)
	s.SetSize(s.Size())
}

func (s *SelectList) MultiSelect() bool {
	return s.multiselect
}

func (s *SelectList) SetMultiSelect(ms bool) {
	s.multiselect = ms
	if !s.multiselect {
		s.UnselectAll()
	}
}
