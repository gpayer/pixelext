package nodes

import "image/color"

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
}
