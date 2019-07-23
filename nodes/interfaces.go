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
	calcLastMat()
	GetLastMat() pixel.Matrix
	SetLastMat(mat pixel.Matrix)
	Initialized() bool
	Parent() Node
	SetParent(p Node)
	SetName(name string)
	GetName() string
	SetOrigin(origin pixel.Vec)
	GetOrigin() pixel.Vec
	GetPos() pixel.Vec
	LocalToGlobalPos(local pixel.Vec) pixel.Vec
	GlobalToLocalPos(global pixel.Vec) pixel.Vec
	SetPos(pos pixel.Vec)
	SetRot(rot float64)
	GetRot() float64
	SetScale(scale pixel.Vec)
	GetScale() pixel.Vec
	SetRotPoint(rotpoint pixel.Vec)
	GetRotPoint() pixel.Vec
	ZIndex() int
	SetZIndex(z int)
	GetExtraOffset() pixel.Vec
	Show()
	Hide()
	AddChild(child Node)
	Children() []Node
	SortChildren()
	RemoveChild(child Node)
	RemoveChildren()
	ChildChanged()
	SetStyles(styles *Styles)
	GetStyles() *Styles
	SetSize(size pixel.Vec)
	Contains(point pixel.Vec) bool
	_updateFromTheme(theme *Theme)
	UpdateFromTheme(theme *Theme)
	SetPausable(pausable bool)
	Pause()
	Unpause()
	SetLocked(locked bool)
	Locked() bool
	Iterate(fn func(n Node))
	CopyFrom(from Node)
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

type Sizer interface {
	Size() pixel.Vec
}
