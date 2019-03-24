package nodes

import (
	"sort"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type BaseNode struct {
	self                      nodeInternal
	Name                      string
	children                  []nodeInternal
	show, active, initialized bool
	mat                       pixel.Matrix
	pos                       pixel.Vec
	bounds                    pixel.Rect
	scale                     pixel.Vec
	rot                       float64
	origin                    pixel.Vec
	zindex                    int
}

func (b *BaseNode) _getMat() pixel.Matrix {
	return b.mat
}

func (b *BaseNode) _init() {
	if b.initialized {
		return
	}
	for _, child := range b.children {
		child._init()
	}
	init, ok := b.self.(Initializable)
	if ok {
		init.Init()
	}
	b.initialized = true
}

func (b *BaseNode) _mount() {
	for _, child := range b.children {
		child._mount()
	}
	mountable, ok := b.self.(Mountable)
	if ok {
		mountable.Mount()
	}
}

func (b *BaseNode) _unmount() {
	for _, child := range b.children {
		child._unmount()
	}
	mountable, ok := b.self.(Mountable)
	if ok {
		mountable.Unmount()
	}
}

func (b *BaseNode) _update(dt float64) {
	if !b.active {
		return
	}
	for _, child := range b.children {
		child._update(dt)
	}
	updateable, ok := b.self.(Updateable)
	if ok {
		updateable.Update(dt)
	}
}

func (b *BaseNode) _draw(win *pixelgl.Window, mat pixel.Matrix) {
	if !b.active || !b.show {
		return
	}
	for _, child := range b.children {
		child._draw(win, mat.Chained(child._getMat()))
	}
	drawable, ok := b.self.(Drawable)
	if ok {
		drawable.Draw(win, mat)
	}
}

func NewBaseNode(name string) *BaseNode {
	b := &BaseNode{
		Name:     name,
		children: make([]nodeInternal, 0),
		show:     true,
		active:   true,
		mat:      pixel.IM,
		pos:      pixel.ZV,
		bounds:   pixel.R(0, 0, 1, 1),
		scale:    pixel.V(1, 1),
		rot:      0,
		origin:   pixel.ZV,
		zindex:   0,
	}
	b.self = b
	b.calcMat()
	return b
}

func (b *BaseNode) calcMat() {
	b.mat = pixel.IM.ScaledXY(b.origin, b.scale).Rotated(b.origin, b.rot).Moved(b.pos)
}

func (b *BaseNode) SetPos(p pixel.Vec) {
	b.pos = p
	b.calcMat()
}

func (b *BaseNode) GetPos() pixel.Vec {
	return b.pos
}

func (b *BaseNode) SetRot(rot float64) {
	b.rot = rot
	b.calcMat()
}

func (b *BaseNode) GetRot() float64 {
	return b.rot
}

func (b *BaseNode) SetScale(scale pixel.Vec) {
	b.scale = scale
	b.calcMat()
}

func (b *BaseNode) GetScale() pixel.Vec {
	return b.scale
}

func (b *BaseNode) SetOrigin(origin pixel.Vec) {
	b.origin = origin
	b.calcMat()
}

func (b *BaseNode) GetOrigin() pixel.Vec {
	return b.origin
}

func (b *BaseNode) SetBounds(r pixel.Rect) {
	b.bounds = r
}

func (b *BaseNode) GetBounds() pixel.Rect {
	return b.bounds
}

func (b *BaseNode) Show() {
	b.show = true
}

func (b *BaseNode) Hide() {
	b.show = false
}

func (b *BaseNode) SetActive(active bool) {
	b.active = active
}

func (b *BaseNode) AddChild(child nodeInternal) {
	b.children = append(b.children, child)
	sort.Slice(b.children, func(i, j int) bool {
		childA := b.children[i].(*BaseNode)
		childB := b.children[j].(*BaseNode)
		return childA.zindex < childB.zindex
	})
}
