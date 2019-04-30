package nodes

import (
	"github.com/faiface/pixel"
	"github.com/lafriks/go-tiled"
)

type TileInfo struct {
	x, y  float64
	idx   int
	drawn bool
}

type TileMapLayer struct {
	BaseNode
	batch       *pixel.Batch
	spritesheet *SpriteSheet
	dx, dy      int
	tiles       []TileInfo
	dirty       bool
}

func NewTileMapLayer(name string, pic pixel.Picture, dx, dy int) *TileMapLayer {
	t := &TileMapLayer{
		BaseNode:    *NewBaseNode(name),
		batch:       pixel.NewBatch(&pixel.TrianglesData{}, pic),
		spritesheet: NewSpriteSheet(pic, dx, dy),
		dx:          dx,
		dy:          dy,
		dirty:       false,
	}
	t.Self = t
	return t
}

func (t *TileMapLayer) AddTileGrid(x, y, idx int) {
	t.tiles = append(t.tiles, TileInfo{float64(x*t.dx + t.dx/2), float64(y*t.dy + t.dy/2), idx, false})
	t.dirty = true
	SceneManager().Redraw()
}

func (t *TileMapLayer) AddTile(x, y float64, idx int) {
	t.tiles = append(t.tiles, TileInfo{x, y, idx, false})
	t.dirty = true
	SceneManager().Redraw()
}

func (t *TileMapLayer) Draw(win pixel.Target, mat pixel.Matrix) {
	if t.dirty {
		for i, tile := range t.tiles {
			if !tile.drawn {
				spr := t.spritesheet.NewSprite(tile.idx)
				spr.Init()
				spr.Draw(t.batch, pixel.IM.Moved(pixel.V(tile.x, tile.y)).Chained(mat))
				t.tiles[i].drawn = true
			}
		}
		t.dirty = false
	}
	t.batch.SetMatrix(mat)
	t.batch.Draw(win)
}

func (t *TileMapLayer) SpriteSheet() *SpriteSheet {
	return t.spritesheet
}

func TileMapsFromTmx(tmx *tiled.Map) []*TileMapLayer {
	var tilemaplayers []*TileMapLayer
	for _, layer := range tmx.Layers {
		t := &TileMapLayer{
			BaseNode: *NewBaseNode(layer.Name),
			dx:       tmx.TileWidth,
			dy:       tmx.TileHeight,
			dirty:    false,
		}
		t.Self = t

	}
	return tilemaplayers
}
