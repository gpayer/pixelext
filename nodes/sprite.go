package nodes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Sprite struct {
	BaseNode
	pic    pixel.Picture
	sprite *pixel.Sprite
}

func (s *Sprite) Init() {
	if s.GetBounds().W()*s.GetBounds().H() == 0 {
		s.SetBounds(s.pic.Bounds())
	}
	s.sprite = pixel.NewSprite(s.pic, s.pic.Bounds())
	s.SetExtraOffset(pixel.V(s.pic.Bounds().W()/2, s.pic.Bounds().H()/2))
}

func (s *Sprite) Draw(win *pixelgl.Window, mat pixel.Matrix) {
	s.sprite.Draw(win, mat)
}

func NewSprite(name string, pic pixel.Picture) *Sprite {
	s := &Sprite{
		BaseNode: *NewBaseNode(name),
		pic:      pic,
	}
	s.Self = s
	return s
}
