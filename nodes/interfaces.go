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
	GetName() string
	GetContainerBounds() pixel.Rect
	GetBounds() pixel.Rect
	SetBounds(r pixel.Rect)
	SetBoundsInternal(r pixel.Rect)
	SetOrigin(origin pixel.Vec)
	GetOrigin() pixel.Vec
	GetPos() pixel.Vec
	SetPos(pos pixel.Vec)
	SetRot(rot float64)
	GetRot() float64
	SetScale(scale pixel.Vec)
	GetScale() pixel.Vec
	SetRotPoint(rotpoint pixel.Vec)
	GetRotPoint() pixel.Vec
	SetZeroAlignment(a Alignment)
	SetZIndex(z int)
	GetExtraOffset() pixel.Vec
	Show()
	Hide()
	AddChild(child Node)
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
