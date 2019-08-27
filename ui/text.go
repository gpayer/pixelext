package ui

import (
	"fmt"
	"image/color"
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

func NewText(name string) *Text {
	t := &Text{
		UIBase: *NewUIBase(name),
	}
	t.Self = t
	t.UISelf = t
	styles := t.GetStyles()
	styles.Text.Atlas = "basic"
	t.txt = text.New(pixel.ZV, nodes.FontService.Get(styles.Text.Atlas))
	return t
}

func NewTextCustom(name, atlasname string, textcolor color.Color) *Text {
	t := NewText(name)
	styles := t.GetStyles()
	styles.Text.Color = textcolor
	styles.Text.Atlas = atlasname
	t.OverrideStyles(styles)
	t.txt = text.New(pixel.ZV, nodes.FontService.Get(styles.Text.Atlas))
	t.txt.Color = textcolor
	return t
}

func (t *Text) Text() *text.Text {
	return t.txt
}

func (t *Text) innerPrintf(format string, a ...interface{}) {
	fmt.Fprintf(&t.content, format, a...)
	fmt.Fprintf(t.txt, format, a...)
	t.SetExtraOffset(pixel.V(-t.txt.Bounds().W()/2, -t.txt.Bounds().H()/2-t.txt.Dot.Y+t.txt.Atlas().Descent()))
	t.SetSize(t.txt.Bounds().Size())
	nodes.SceneManager().Redraw()
}

func (t *Text) Printf(format string, a ...interface{}) {
	t.innerPrintf(format, a...)
}

func (t *Text) GetContent() string {
	return t.content.String()
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
	if oldstyles.Text.Atlas != styles.Text.Atlas {
		t.txt = text.New(pixel.ZV, nodes.FontService.Get(styles.Text.Atlas))
		redraw = true
	}
	t.UIBase.SetStyles(styles)
	if redraw {
		t.txt.Color = styles.Text.Color
		t.txt.Orig = pixel.ZV
		t.txt.Clear()
		fmt.Fprint(t.txt, t.content.String())
		t.SetExtraOffset(pixel.V(-t.txt.Bounds().W()/2, -t.txt.Bounds().H()/2-t.txt.Dot.Y+t.txt.Atlas().Descent()))
		t.SetSize(t.txt.Bounds().Size())
		nodes.SceneManager().Redraw()
	}
}

func (t *Text) UpdateFromTheme(theme *nodes.Theme) {
	if t.overrideStyles {
		return
	}
	t.SetStyles(theme.Text)
}
