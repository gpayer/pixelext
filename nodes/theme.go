package nodes

import (
	"image/color"

	"golang.org/x/image/colornames"
)

type Theme struct {
	Button struct {
		Enabled, Disabled, Hover, Pressed *Styles
	}
	Text     *Styles
	Slider   *Styles
	Grid     *Styles
	InputBox *Styles
}

func DefaultTheme() *Theme {
	t := &Theme{
		Text:     DefaultStyles(),
		Slider:   DefaultStyles(),
		Grid:     DefaultStyles(),
		InputBox: DefaultStyles(),
	}
	baseStyle := DefaultStyles()
	baseStyle.Border.Width = 2

	t.Button.Enabled = baseStyle.Clone()
	t.Button.Enabled.Text.Color = colornames.Black
	t.Button.Enabled.Background.Color = colornames.Steelblue
	t.Button.Disabled = baseStyle.Clone()
	t.Button.Disabled.Background.Color = colornames.Darkgray
	t.Button.Hover = baseStyle.Clone()
	t.Button.Hover.Background.Color = colornames.Lightblue
	t.Button.Hover.Text.Color = colornames.Black
	t.Button.Pressed = baseStyle.Clone()
	t.Button.Pressed.Border.Width = 5
	t.Button.Pressed.Border.Color = color.RGBA{10, 10, 10, 255}
	t.Button.Pressed.Background.Color = colornames.Darkblue
	t.Button.Pressed.Border.Color = colornames.Gray
	return t
}
