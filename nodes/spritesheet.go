package nodes

import (
	"github.com/faiface/pixel"
)

type SpriteSheet struct {
	pic   pixel.Picture
	rects []pixel.Rect
}

func NewSpriteSheet(spritesheet pixel.Picture, dx int, dy int) *SpriteSheet {
	o := &SpriteSheet{
		pic: spritesheet,
	}
	fx, fy := float64(dx), float64(dy)
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += fx {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += fy {
			o.rects = append(o.rects, pixel.R(x, y, x+fx, y+fy))
		}
	}
	return o
}

func (spr *SpriteSheet) NewSprite(idx int) *Sprite {
	sprite := NewSprite("", spr.pic)
	sprite.SetBounds(spr.rects[idx])
	return sprite
}

func (spr *SpriteSheet) Length() int {
	return len(spr.rects)
}
