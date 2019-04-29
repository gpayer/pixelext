package nodes

import (
	"github.com/faiface/pixel"
)

type SpriteSheetInfo struct {
	MaxX, MaxY, MaxIdx int
}

func (i *SpriteSheetInfo) Idx(x, y int) int {
	return x*i.MaxY + y
}

type SpriteSheet struct {
	pic   pixel.Picture
	rects []pixel.Rect
	info  SpriteSheetInfo
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
	o.info.MaxX = int(spritesheet.Bounds().W() / fx)
	o.info.MaxY = int(spritesheet.Bounds().H() / fy)
	o.info.MaxIdx = len(o.rects)
	return o
}

func (spr *SpriteSheet) NewSprite(idx int) *Sprite {
	sprite := NewSprite("", spr.pic)
	sprite.SetBounds(spr.rects[idx])
	return sprite
}

func (spr *SpriteSheet) Length() int {
	return spr.info.MaxIdx
}

func (spr *SpriteSheet) Info() SpriteSheetInfo {
	return spr.info
}
