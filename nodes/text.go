package nodes

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

type Text struct {
	BaseNode
	txt *text.Text
}

func NewText(name, atlasname string) *Text {
	t := &Text{
		BaseNode: *NewBaseNode(name),
		txt:      text.New(pixel.ZV, FontService.Get(atlasname)),
	}
	t.Self = t
	return t
}

func (t *Text) Text() *text.Text {
	return t.txt
}

func (t *Text) Printf(format string, a ...interface{}) {
	fmt.Fprintf(t.txt, format, a...)
	t.SetBounds(t.txt.Bounds())
}

func (t *Text) Draw(win *pixelgl.Window, mat pixel.Matrix) {
	t.txt.Draw(win, mat)
}

func (t *Text) Clear() {
	t.txt.Clear()
}
