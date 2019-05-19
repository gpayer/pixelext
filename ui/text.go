package ui

import (
	"fmt"
	"strings"

	"github.com/gpayer/pixelext/nodes"

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
		txt:    text.New(pixel.ZV, nodes.FontService.Get(atlasname)), // TODO: correct this cheap workaround for real
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
	tH := t.txt.Atlas().Descent() + t.txt.Atlas().Ascent()
	t.SetExtraOffset(pixel.V(-t.txt.Bounds().W()/2, -tH/2+t.txt.Atlas().Descent()))
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
	oldstyles := t.GetStyles()
	if oldstyles.Text.Color != styles.Text.Color || oldstyles.Text.OffsetY != styles.Text.OffsetY {
		redraw = true
	}
	t.BaseNode.SetStyles(styles)
	if redraw {
		t.txt.Color = styles.Text.Color
		t.txt.Orig = pixel.V(0, styles.Text.OffsetY)
		t.txt.Clear()
		fmt.Fprint(t.txt, t.content.String())
		nodes.SceneManager().Redraw()
	}
}

func (t *Text) UpdateFromTheme(theme *nodes.Theme) {
	if t.overrideStyles {
		return
	}
	t.SetStyles(theme.Text)
}
