package nodes

import (
	"fmt"
	"image/color"
	"github.com/gpayer/pixelext/services"

	"github.com/faiface/pixel"
	"github.com/gpayer/go-tiled"
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

func TileMapsFromTmx(tmx *tiled.Map) ([]*TileMapLayer, error) {
	if tmx.Orientation != "orthogonal" {
		return nil, fmt.Errorf("only orthogonal tilemaps supported")
	}
	var tilemaplayers []*TileMapLayer
	for z, layer := range tmx.Layers {
		// TODO: only tilemap layers
		var tileset *tiled.Tileset
		for _, tile := range layer.Tiles {
			if tile.ID > 0 {
				tileset = tile.Tileset
			}
		}
		pic, err := services.ResourceManager().LoadPicture(tileset.Image.Source)
		if err != nil {
			return nil, err
		}
		t := &TileMapLayer{
			BaseNode:    *NewBaseNode(layer.Name),
			batch:       pixel.NewBatch(&pixel.TrianglesData{}, pic),
			spritesheet: NewSpriteSheet(pic, tmx.TileWidth, tmx.TileHeight),
			dx:          tmx.TileWidth,
			dy:          tmx.TileHeight,
			dirty:       false,
		}
		t.Self = t
		t.SetZIndex(z)
		t.batch.SetColorMask(color.RGBA{255, 255, 255, uint8(255.0 * layer.Opacity)})

		w := float64(tmx.Width * tmx.TileWidth)
		h := float64(tmx.Height * tmx.TileHeight)
		tw := float64(tmx.TileWidth)
		th := float64(tmx.TileHeight)
		tw2 := float64(tmx.TileWidth) / 2
		th2 := float64(tmx.TileHeight) / 2
		var x, y, dx, dy float64
		switch tmx.RenderOrder {
		case "right-down":
			x = -w/2 + tw2
			y = h/2 - th2
			dx = tw
			dy = -th
		case "right-up":
			x = -w/2 + tw2
			y = -h/2 + th2
			dx = tw
			dy = th
		case "left-down":
			x = w/2 - tw2
			y = h/2 - th2
			dx = -tw
			dy = -th
		case "left-up":
			x = w/2 - tw2
			y = -h/2 + th2
			dx = -tw
			dy = th
		}

		inity := y
		for tx := 0; tx < tmx.Width; tx++ {
			y = inity
			for ty := 0; ty < tmx.Height; ty++ {
				tileid := layer.Tiles[ty*tmx.Width+tx].ID
				t.AddTile(x, y, int(tileid))
				y += dy
			}
			x += dx
		}

		tilemaplayers = append(tilemaplayers, t)
	}
	return tilemaplayers, nil
}
