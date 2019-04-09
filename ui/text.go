package ui

import (
	"fmt"
	"pixelext/nodes"
	"strings"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

type Text struct {
	UIBase
	txt     *text.Text
	content strings.Builder
}

func NewText(name, atlasname string) *Text {
	t := &Text{
		UIBase: *NewUIBase(name),
		txt:    text.New(pixel.ZV, nodes.FontService.Get(atlasname)),
	}
	t.Self = t
	t.UISelf = t
	return t
}

func (t *Text) Text() *text.Text {
	return t.txt
}

func (t *Text) Printf(format string, a ...interface{}) {
	fmt.Fprintf(&t.content, format, a...)
	fmt.Fprintf(t.txt, format, a...)
	t.SetExtraOffset(pixel.V(0, t.txt.Bounds().H()-t.txt.Atlas().LineHeight()))
	//t.SetBounds(t.txt.Bounds())
}

func (t *Text) Draw(win *pixelgl.Window, mat pixel.Matrix) {
	t.txt.Draw(win, mat)
}

func (t *Text) Clear() {
	t.content.Reset()
	t.txt.Clear()
}

func (t *Text) SetStyles(styles *nodes.Styles) {
	redraw := false
	if t.GetStyles().Text.Color != styles.Text.Color {
		redraw = true
	}
	t.BaseNode.SetStyles(styles)
	if redraw {
		t.txt.Color = styles.Text.Color
		t.txt.Clear()
		fmt.Fprint(t.txt, t.content.String())
	}
}
