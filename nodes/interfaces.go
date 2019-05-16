package nodes

import (
	"github.com/faiface/pixel"
)

type Node interface {
	_init()
	_mount()
	_unmount()
	_update(dt float64)
	_draw(win pixel.Target, mat pixel.Matrix)
	_getMat() pixel.Matrix
	_getLastMat() pixel.Matrix
	_setLastMat(mat pixel.Matrix)
	_getZindex() int
	Initialized() bool
	GetName() string
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
	SetZIndex(z int)
	GetExtraOffset() pixel.Vec
	Show()
	Hide()
	AddChild(child Node)
	Children() []Node
	RemoveChild(child Node)
	RemoveChildren()
	SetStyles(styles *Styles)
	GetStyles() *Styles
	SetSize(size pixel.Vec)
	Contains(point pixel.Vec) bool
	_updateFromTheme(theme *Theme)
	UpdateFromTheme(theme *Theme)
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
	Draw(win pixel.Target, mat pixel.Matrix)
}
