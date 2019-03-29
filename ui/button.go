package ui

import (
	"image/color"
	"pixelext/nodes"

	"github.com/faiface/pixel/pixelgl"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type ButtonState int

const (
	ButtonEnabled ButtonState = iota
	ButtonDisabled
	ButtonHover
	ButtonPressed
)

type Button struct {
	nodes.BaseNode
	state    ButtonState
	canvases map[ButtonState]*nodes.Canvas
	text     *nodes.Text
	enabled  bool
	onclick  func()
}

func NewButton(name string, w, h float64, text string) *Button {
	b := &Button{
		BaseNode: *nodes.NewBaseNode(name),
		state:    ButtonEnabled,
		canvases: make(map[ButtonState]*nodes.Canvas, 4),
		enabled:  true,
		onclick:  func() {},
	}
	b.Self = b
	states := []ButtonState{ButtonEnabled, ButtonDisabled, ButtonHover, ButtonPressed}
	for _, state := range states {
		b.canvases[state] = nodes.NewCanvas("", w, h)
		b.canvases[state].GetStyles().Border.Width = 1
		b.AddChild(b.canvases[state])
	}

	b.canvases[ButtonHover].GetStyles().Element.EnabledColor = colornames.Lightblue
	b.canvases[ButtonHover].GetStyles().Text.Color = colornames.Black
	b.canvases[ButtonEnabled].GetStyles().Text.Color = colornames.Black
	b.canvases[ButtonPressed].GetStyles().Border.Width = 5
	b.canvases[ButtonPressed].GetStyles().Border.Color = color.RGBA{10, 10, 10, 255}
	b.canvases[ButtonPressed].GetStyles().Element.EnabledColor = colornames.Darkblue
	b.canvases[ButtonPressed].GetStyles().Border.Color = colornames.Gray

	b.text = nodes.NewText("buttontxt", "basic")
	b.text.Printf(text)
	b.text.SetPos(pixel.V(w/2, h/2))
	b.text.SetZeroAlignment(nodes.AlignmentCenter)
	b.text.SetZIndex(10)
	b.AddChild(b.text)
	b.SetBounds(pixel.R(0, 0, w, h))
	return b
}

func (b *Button) drawCanvases() {
	for state, canvas := range b.canvases {
		styles := canvas.GetStyles()
		if state == ButtonDisabled {
			canvas.Clear(styles.Element.DisabledColor)
		} else {
			canvas.Clear(styles.Element.EnabledColor)
		}
		if styles.Border.Width > 0 {
			bounds := b.GetBounds()
			im := imdraw.New(nil)
			im.Color = styles.Border.Color
			im.Push(pixel.ZV,
				pixel.V(0, bounds.H()),
				pixel.V(bounds.W(), bounds.H()),
				pixel.V(bounds.W(), 0))
			im.Polygon(styles.Border.Width)
			im.Draw(canvas.Canvas())
		}
	}
}

func (b *Button) SetBounds(r pixel.Rect) {
	oldbounds := b.GetBounds()
	var redraw bool = false
	if oldbounds.W() != r.W() || oldbounds.H() != r.H() {
		redraw = true
	}
	b.BaseNode.SetBounds(r)

	if redraw {
		for _, canvas := range b.canvases {
			canvas.SetBounds(pixel.R(0, 0, r.W(), r.H()))
		}
		b.drawCanvases()
	}
}

func (b *Button) SetButtonStyles(state ButtonState, styles *nodes.Styles) {
	b.canvases[state].SetStyles(styles)
	b.drawCanvases()
}

func (b *Button) OnClick(fn func()) {
	b.onclick = fn
}

func (b *Button) SetEnabled(enabled bool) {
	b.enabled = enabled
}

func (b *Button) Update(dt float64) {
	if b.enabled {
		if nodes.Events().Clicked(pixelgl.MouseButtonLeft, b) {
			b.state = ButtonPressed
		} else if b.state == ButtonPressed && nodes.Events().Pressed(pixelgl.MouseButtonLeft) {
			b.state = ButtonPressed
		} else if nodes.Events().JustReleased(pixelgl.MouseButtonLeft) {
			b.state = ButtonEnabled
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

	for state, canvas := range b.canvases {
		if state == b.state {
			canvas.Show()
			b.text.SetStyles(canvas.GetStyles())
		} else {
			canvas.Hide()
		}
	}
}
