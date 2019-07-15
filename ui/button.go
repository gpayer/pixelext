package ui

import (
	"github.com/gpayer/pixelext/nodes"
	"github.com/gpayer/pixelext/services"

	"github.com/faiface/pixel/pixelgl"

	"github.com/faiface/pixel"
)

type ButtonState int

const (
	ButtonInit ButtonState = iota
	ButtonEnabled
	ButtonDisabled
	ButtonHover
	ButtonPressed
)

type Button struct {
	UIBase
	state       ButtonState
	canvases    map[ButtonState]*nodes.Canvas
	text        *Text
	enabled     bool
	w, h        float64
	textcontent string
	onclick     func()
}

func NewButton(name string, w, h float64, text string) *Button {
	b := &Button{
		UIBase:      *NewUIBase(name),
		state:       ButtonInit,
		canvases:    make(map[ButtonState]*nodes.Canvas, 4),
		enabled:     true,
		w:           w,
		h:           h,
		textcontent: text,
		onclick:     func() {},
	}
	b.Self = b
	b.UISelf = b
	states := []ButtonState{ButtonEnabled, ButtonDisabled, ButtonHover, ButtonPressed}
	for _, state := range states {
		b.canvases[state] = nodes.NewCanvas("", w, h)
		b.canvases[state].GetStyles().Border.Width = 2
		b.canvases[state].SetLocked(true)
		b.AddChild(b.canvases[state])
	}

	b.createText()
	return b
}

func (b *Button) alignText() {
	p := b.GetStyles().Padding
	size := b.Size()
	w := size.X/2 - p
	h := size.Y/2 - p
	bounds := pixel.R(-w, -h, w, h)
	AlignUINode(b.text, bounds, b.VAlignment(), b.HAlignment())
}

func (b *Button) createText() {
	styles := b.canvases[ButtonEnabled].GetStyles()
	b.text = NewText("buttontxt", styles.Text.Atlas)
	b.text.Printf(b.textcontent)
	b.text.SetLocked(true)
	b.text.SetZIndex(10)
	b.AddChild(b.text)
	b.internalSetSize()
}

func (b *Button) drawCanvases() {
	for _, canvas := range b.canvases {
		styles := canvas.GetStyles()
		canvas.Clear(styles.Background.Color)
		if styles.Border.Width > 0 {
			bounds := b.Size()
			canvas.DrawRect(pixel.ZV, pixel.V(bounds.X-1, bounds.Y-1), styles.Border.Color)
		}
	}
	nodes.SceneManager().Redraw()
}

func (b *Button) SetSize(size pixel.Vec) {
	b.w = size.X
	b.h = size.Y
	b.internalSetSize()
}

func (b *Button) GetOrigSize() pixel.Vec {
	return pixel.V(b.w, b.h)
}

func (b *Button) internalSetSize() {
	w := b.w
	h := b.h
	styles := b.canvases[ButtonEnabled].GetStyles()
	if w == 0 {
		w = b.text.Size().X + 2*styles.Padding
		b.SetVAlignment(nodes.VAlignmentCenter)
	}
	if h == 0 {
		h = b.text.Size().Y + 2*styles.Padding
		b.SetHAlignment(nodes.HAlignmentCenter)
	}
	size := pixel.V(w, h)
	b.UIBase.SetSize(size)
	for _, canvas := range b.canvases {
		canvas.SetSize(size)
	}
	b.alignText()
	b.drawCanvases()
}

func (b *Button) SetButtonStyles(state ButtonState, styles *nodes.Styles) {
	b.overrideStyles = true
	b.canvases[state].SetStyles(styles)
	b.drawCanvases()
	if state == b.state {
		b.text.OverrideStyles(styles)
	}
}

func (b *Button) SetStyles(styles *nodes.Styles) {
	oldatlas := b.canvases[ButtonEnabled].GetStyles().Text.Atlas
	b.UIBase.SetStyles(styles)
	if oldatlas != styles.Text.Atlas {
		b.RemoveChild(b.text)
		b.createText()
	}
	b.alignText()
}

func (b *Button) UpdateFromTheme(theme *nodes.Theme) {
	if b.overrideStyles {
		return
	}
	oldAtlas := b.canvases[ButtonEnabled].GetStyles().Text.Atlas
	b.canvases[ButtonEnabled].SetStyles(theme.Button.Enabled)
	b.canvases[ButtonDisabled].SetStyles(theme.Button.Disabled)
	b.canvases[ButtonHover].SetStyles(theme.Button.Hover)
	b.canvases[ButtonPressed].SetStyles(theme.Button.Pressed)

	b.drawCanvases()

	if oldAtlas != theme.Button.Enabled.Text.Atlas {
		b.RemoveChild(b.text)
		b.createText()
	}
}

func (b *Button) OnClick(fn func()) {
	b.onclick = fn
}

func (b *Button) SetEnabled(enabled bool) {
	b.enabled = enabled
}

func (b *Button) SetText(text string) {
	b.textcontent = text
	b.text.Clear()
	b.text.Printf(text)
	b.internalSetSize()
}

func (b *Button) SetVAlignment(v nodes.VerticalAlignment) {
	b.UIBase.SetVAlignment(v)
	b.alignText()
}

func (b *Button) SetHAlignment(h nodes.HorizontalAlignment) {
	b.UIBase.SetHAlignment(h)
	b.alignText()
}

func (b *Button) Update(dt float64) {
	oldstate := b.state
	if b.enabled {
		if nodes.Events().Clicked(pixelgl.MouseButtonLeft, b) {
			b.state = ButtonPressed
		} else if b.state == ButtonPressed && nodes.Events().Pressed(pixelgl.MouseButtonLeft) {
			b.state = ButtonPressed
		} else if b.state == ButtonPressed && nodes.Events().JustReleased(pixelgl.MouseButtonLeft) {
			b.state = ButtonEnabled
			clicksample := b.canvases[ButtonPressed].GetStyles().Sound.Click
			if len(clicksample) > 0 {
				services.AudioManager().PlaySample(clicksample, 1.0)
			}
			b.onclick()
		} else {
			b.state = ButtonEnabled
		}
		if b.state == ButtonEnabled && nodes.Events().MouseHovering(b) {
			b.state = ButtonHover
		}
	} else {
		b.state = ButtonDisabled
	}

	if b.state == oldstate {
		return
	}

	for state, canvas := range b.canvases {
		if state == b.state {
			canvas.Show()
			b.text.OverrideStyles(canvas.GetStyles())
		} else {
			canvas.Hide()
		}
	}
}
