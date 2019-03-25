package nodes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type nodeInternal interface {
	_init()
	_mount()
	_unmount()
	_update(dt float64)
	_draw(win *pixelgl.Window, mat pixel.Matrix)
	_getMat() pixel.Matrix
	_getZindex() int
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
