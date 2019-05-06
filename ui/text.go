package ui

import (
	"fmt"
	"github.com/gpayer/pixelext/nodes"
	"strings"

	"github.com/faiface/pixel"
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
		txt:    text.New(pixel.V(0, 2), nodes.FontService.Get(atlasname)), // TODO: correct this cheap workaround for real
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
	t.SetExtraOffset(pixel.V(-t.txt.Bounds().W()/2, t.txt.Bounds().H()/2-t.txt.Atlas().LineHeight()))
	t.SetSize(t.txt.Bounds().Size())
	nodes.SceneManager().Redraw()
}

func (t *Text) Draw(win pixel.Target, mat pixel.Matrix) {
	t.txt.Draw(win, mat)
}

func (t *Text) Clear() {
	t.content.Reset()
	t.txt.Clear()
	nodes.SceneManager().Redraw()
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
		nodes.SceneManager().Redraw()
	}
}
