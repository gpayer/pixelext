package ui

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/gpayer/pixelext/nodes"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
)

const maxCacheText = 1000

var cacheTextFifo chan *text.Text

type Text struct {
	UIBase
	txt     *text.Text
	content strings.Builder
}

func NewText(name, atlasname string) *Text {
	t := &Text{
		UIBase: *NewUIBase(name),
	}
	t.Self = t
	t.UISelf = t
	t.overrideStyles = true // TODO: make this method without atlas parameter
	styles := t.GetStyles()
	styles.Text.Atlas = atlasname
	t.attachText()
	return t
}

func NewTextCustom(name, atlasname string, textcolor color.Color) *Text {
	t := NewText(name, atlasname)
	t.overrideStyles = true
	t.txt.Color = textcolor
	styles := t.GetStyles()
	styles.Text.Color = textcolor
	styles.Text.Atlas = atlasname
	return t
}

func (t *Text) attachText() {
	if t.txt == nil {
		styles := t.GetStyles()
		select {
		case cachedTxt := <-cacheTextFifo:
			t.txt = cachedTxt
			if t.txt == nil {
				panic(fmt.Errorf("t.txt is nil"))
			}
			content := t.content.String()
			content = (content + " ")[:len(content)]
			t.content.Reset()
			t.txt.Clear()
			t.txt.Color = t.GetStyles().Text.Color
			t.innerPrintf(content)
			t.SetPos(t.GetOrigPos())
		default:
			t.txt = text.New(pixel.ZV, nodes.FontService.Get(styles.Text.Atlas))
		}
	}
}

func (t *Text) detachText() {
	if t.txt == nil {
		return
	}
	select {
	case cacheTextFifo <- t.txt:
	default:
	}
	t.txt = nil
}

func (t *Text) Text() *text.Text {
	t.attachText()
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
	t.attachText()
	t.innerPrintf(format, a...)
}

func (t *Text) GetContent() string {
	return t.content.String()
}

func (t *Text) Mount() {
	t.attachText()
}

func (t *Text) Unmount() {
	t.detachText()
}

func (t *Text) Draw(win pixel.Target, mat pixel.Matrix) {
	t.txt.Draw(win, mat)
}

func (t *Text) Clear() {
	t.attachText()
	t.content.Reset()
	t.txt.Clear()
	nodes.SceneManager().Redraw()
}

func (t *Text) Size() pixel.Vec {
	t.attachText()
	return t.UIBase.Size()
}

func (t *Text) SetStyles(styles *nodes.Styles) {
	t.attachText()
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
		nodes.SceneManager().Redraw()
	}
}

func (t *Text) UpdateFromTheme(theme *nodes.Theme) {
	if t.overrideStyles {
		return
	}
	t.SetStyles(theme.Text)
}

func init() {
	cacheTextFifo = make(chan *text.Text, maxCacheText)
}
