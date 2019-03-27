package nodes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Node interface {
	_init()
	_mount()
	_unmount()
	_update(dt float64)
	_draw(win *pixelgl.Window, mat pixel.Matrix)
	_getMat() pixel.Matrix
	_getLastMat() pixel.Matrix
	_getZindex() int
	GetContainerBounds() pixel.Rect
	GetBounds() pixel.Rect
	GetOrigin() pixel.Vec
	GetPos() pixel.Vec
	SetPos(pos pixel.Vec)
	GetExtraOffset() pixel.Vec
	Children() []Node
}

type Initializable interface {
	Init()
}

type Mountable interface {
	Mount()
	Unmount()
}

type Updateable interface {
	Update(dt float64)
}

type Drawable interface {
	Draw(win *pixelgl.Window, mat pixel.Matrix)
}
