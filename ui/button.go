package ui

import (
	"image/color"

	"github.com/gpayer/pixelext/nodes"
	"github.com/gpayer/pixelext/services"

	"github.com/faiface/pixel/pixelgl"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
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
		b.AddChild(b.canvases[state])
	}

	b.canvases[ButtonHover].GetStyles().Background.Color = colornames.Lightblue
	b.canvases[ButtonHover].GetStyles().Text.Color = colornames.Black
	b.canvases[ButtonEnabled].GetStyles().Text.Color = colornames.Black
	b.canvases[ButtonPressed].GetStyles().Border.Width = 5
	b.canvases[ButtonPressed].GetStyles().Border.Color = color.RGBA{10, 10, 10, 255}
	b.canvases[ButtonPressed].GetStyles().Background.Color = colornames.Darkblue
	b.canvases[ButtonPressed].GetStyles().Border.Color = colornames.Gray

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
	w := b.w
	h := b.h
	b.text = NewText("buttontxt", styles.Text.Atlas)
	b.text.Printf(b.textcontent)
	if w == 0 {
		w = b.text.Size().X + 2*styles.Padding
		b.SetVAlignment(nodes.VAlignmentCenter)
	}
	if h == 0 {
		h = b.text.Size().Y + 2*styles.Padding
		b.SetHAlignment(nodes.HAlignmentCenter)
	}
	b.text.SetZIndex(10)
	b.AddChild(b.text)
	b.SetSize(pixel.V(w, h))
}

func (b *Button) drawCanvases() {
	for _, canvas := range b.canvases {
		styles := canvas.GetStyles()
		canvas.Clear(styles.Background.Color)
		if styles.Border.Width > 0 {
			bounds := b.Size()
			im := imdraw.New(nil)
			im.Color = styles.Border.Color
			im.Push(pixel.ZV,
				pixel.V(0, bounds.Y),
				pixel.V(bounds.X, bounds.Y),
				pixel.V(bounds.X, 0))
			im.Polygon(styles.Border.Width)
			im.Draw(canvas.Canvas())
		}
	}
	nodes.SceneManager().Redraw()
}

func (b *Button) SetSize(size pixel.Vec) {
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
			b.text.SetStyles(canvas.GetStyles())
		} else {
			canvas.Hide()
		}
	}
}
