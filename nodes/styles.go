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
	Foreground struct {
		Color color.Color
	}
	Text struct {
		Color   color.Color
		Atlas   string
		OffsetY float64
	}
	Padding       float64
	RoundToPixels bool
}

func DefaultStyles() *Styles {
	s := &Styles{}
	s.Border.Width = 2
	s.Border.Color = colornames.White
	s.Background.Color = colornames.Black
	s.Foreground.Color = colornames.Steelblue
	s.Text.Color = colornames.White
	s.Text.Atlas = "basic"
	s.Text.OffsetY = 0
	s.Padding = 5
	s.RoundToPixels = true
	return s
}

func (s *Styles) Clone() *Styles {
	clonedStyle := &Styles{}
	*clonedStyle = *s
	return clonedStyle
}
