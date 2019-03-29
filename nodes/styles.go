package nodes

import (
	"image/color"

	"golang.org/x/image/colornames"
)

type Styles struct {
	Border struct {
		Width float64
		Color color.Color
	}
	Background struct {
		Color color.Color
	}
	Text struct {
		Color color.Color
	}
	Element struct {
		EnabledColor  color.Color
		DisabledColor color.Color
	}
	Padding float64
}

func DefaultStyles() *Styles {
	s := &Styles{}
	s.Border.Width = 2
	s.Border.Color = colornames.White
	s.Background.Color = colornames.Black
	s.Text.Color = colornames.White
	s.Element.EnabledColor = colornames.Steelblue
	s.Element.DisabledColor = colornames.Darkgray
	s.Padding = 5
	return s
}
