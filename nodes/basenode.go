package nodes

import (
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
	scale                     float64
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
	for _, child := range b.children {
		child._update(dt)
	}
	updateable, ok := b.self.(Updateable)
	if ok {
		updateable.Update(dt)
	}
}

func (b *BaseNode) _draw(win *pixelgl.Window, mat pixel.Matrix) {
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
		scale:    1,
		rot:      0,
		origin:   pixel.ZV,
		zindex:   0,
	}
	b.self = b
	b.calcMat()
	return b
}

func (b *BaseNode) calcMat() {
	b.mat = pixel.IM.Rotated(b.origin, b.rot).Scaled(b.origin, b.scale).Moved(b.pos)
}
