package nodes

import (
	"github.com/faiface/pixel"
)

type Sprite struct {
	BaseNode
	pic    pixel.Picture
	sprite *pixel.Sprite
	bounds pixel.Rect
}

func (s *Sprite) Init() {
	s.sprite = pixel.NewSprite(s.pic, s.bounds)
}

func (s *Sprite) Draw(win pixel.Target, mat pixel.Matrix) {
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
	s.bounds = s.pic.Bounds()
	s.Self = s
	return s
}

func (s *Sprite) SetBounds(bounds pixel.Rect) {
	s.bounds = bounds
}
