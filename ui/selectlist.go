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
	content     interface{}
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
	onselect                            func(v string, content interface{})
	onunselect                          func(v string, content interface{})
	multiselect                         bool
	btnPool                             []*Button
	btnCurrent                          []*Button
}

func NewSelectList(name string, w, h float64, multiselect bool) *SelectList {
	defaultStyles := nodes.DefaultStyles()

	s := &SelectList{
		UIBase:      *NewUIBase(name),
		background:  nodes.NewBorderBox("background", w, h),
		vscroll:     NewVScroll("vscroll", w-2*defaultStyles.Padding-2*defaultStyles.Border.Width, h-2*defaultStyles.Border.Width),
		list:        NewVBox("list"),
		value2entry: make(map[string]*selectListEntry),
		onselect:    func(_ string, _ interface{}) {},
		onunselect:  func(_ string, _ interface{}) {},
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
	s.vscroll.SetAlignment(nodes.AlignmentTopCenter)
	size := s.Size()
	s.vscroll.SetPos(pixel.V(0, size.Y/2-s.GetStyles().Border.Width))
	s.AddChild(s.vscroll)
	st := s.list.GetStyles()
	st.Border.Width = 0
	st.Padding = 2
	s.list.OverrideStyles(st)
	s.list.SetHAlignment(nodes.HAlignmentLeft)
	s.vscroll.SetInner(s.list)
}

func (s *SelectList) OnSelect(onselect func(v string, content interface{})) {
	s.onselect = onselect
}

func (s *SelectList) OnUnselect(onunselect func(v string, content interface{})) {
	s.onunselect = onunselect
}

func (s *SelectList) AddEntry(text string, value string, content interface{}) {
	btnW := s.Size().X - 2*s.GetStyles().Padding
	var btn *Button
	if len(s.btnPool) > 0 {
		btn = s.btnPool[len(s.btnPool)-1]
		s.btnPool = s.btnPool[:len(s.btnPool)-1]
		btn.SetSize(pixel.V(btnW, 0))
		btn.SetText(text)
	} else {
		btn = NewButton("btn", btnW, 0, text)
	}
	s.btnCurrent = append(s.btnCurrent, btn)
	btn.SetHAlignment(nodes.HAlignmentLeft)
	btn.OnClick(func() {
		e := s.value2entry[value]
		if e.selected {
			if s.multiselect {
				s.onunselect(value, e.content)
				e.btn.SetButtonStyles(ButtonEnabled, s.enabledStyle)
			}
		} else {
			if !s.multiselect {
				s.UnselectAll()
			}
			s.onselect(value, e.content)
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
		content:  content,
	}
	s.entries = append(s.entries, newentry)
	s.value2entry[value] = newentry
}

func (s *SelectList) Clear() {
	s.entries = s.entries[:0]
	s.value2entry = make(map[string]*selectListEntry)
	s.btnPool = append(s.btnPool, s.btnCurrent...)
	s.btnCurrent = s.btnCurrent[:0]
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

type SelectListSelection struct {
	value   string
	content interface{}
}

func (s *SelectList) Selected() []SelectListSelection {
	var sel []SelectListSelection
	for _, e := range s.entries {
		sel = append(sel, SelectListSelection{
			value:   e.value,
			content: e.content,
		})
	}
	return sel
}

func (s *SelectList) SetSize(size pixel.Vec) {
	s.UIBase.SetSize(size)
	s.background.SetSize(size)
	st := s.GetStyles()
	s.vscroll.SetSize(size.Sub(pixel.V(2*st.Padding+2*st.Border.Width, 2*st.Border.Width)))
	s.vscroll.SetPos(pixel.V(0, size.Y/2-st.Border.Width))
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
