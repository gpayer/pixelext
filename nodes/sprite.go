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
	s.sprite = pixel.NewSprite(s.pic, s.pic.Bounds())
	//s.SetExtraOffset(pixel.V(s.pic.Bounds().W()/2, s.pic.Bounds().H()/2))
}

func (s *Sprite) Draw(win *pixelgl.Window, mat pixel.Matrix) {
	s.sprite.Draw(win, mat)
}

func (s *Sprite) Contains(point pixel.Vec) bool {
	bounds := s.pic.Bounds()
	bounds = bounds.Moved(bounds.Size().Scaled(-.5))
	return bounds.Contains(point)
}

func NewSprite(name string, pic pixel.Picture) *Sprite {
	s := &Sprite{
		BaseNode: *NewBaseNode(name),
		pic:      pic,
	}
	s.Self = s
	return s
}
