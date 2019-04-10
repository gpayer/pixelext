package nodes

import (
	"math"
	"sort"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type BaseNode struct {
	Self                      Node
	Name                      string
	children                  []Node
	show, active, initialized bool
	mat, lastmat              pixel.Matrix
	pos                       pixel.Vec
	bounds                    pixel.Rect
	scale                     pixel.Vec
	rot                       float64
	rotpoint                  pixel.Vec
	origin                    pixel.Vec
	zindex                    int
	zeroalignment             Alignment
	extraoffset               pixel.Vec
	styles                    *Styles
}

func (b *BaseNode) _getMat() pixel.Matrix {
	return b.mat
}

func (b *BaseNode) _getLastMat() pixel.Matrix {
	return b.lastmat
}

func (b *BaseNode) _getZindex() int {
	return b.zindex
}

func (b *BaseNode) _init() {
	if b.initialized {
		return
	}
	for _, child := range b.children {
		child._init()
	}
	init, ok := b.Self.(Initializable)
	if ok {
		init.Init()
	}
	b.initialized = true
}

func (b *BaseNode) _mount() {
	for _, child := range b.children {
		child._mount()
	}
	mountable, ok := b.Self.(Mountable)
	if ok {
		mountable.Mount()
	}
}

func (b *BaseNode) _unmount() {
	for _, child := range b.children {
		child._unmount()
	}
	mountable, ok := b.Self.(Mountable)
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
	updateable, ok := b.Self.(Updateable)
	if ok {
		updateable.Update(dt)
	}
}

func (b *BaseNode) _draw(win *pixelgl.Window, mat pixel.Matrix) {
	b.lastmat = mat
	if !b.active || !b.show {
		return
	}
	drawable, ok := b.Self.(Drawable)
	if ok {
		drawable.Draw(win, mat)
	}
	for _, child := range b.children {
		child._draw(win, child._getMat().Chained(mat))
	}
}

func NewBaseNode(name string) *BaseNode {
	b := &BaseNode{
		Name:          name,
		children:      make([]Node, 0),
		show:          true,
		active:        true,
		mat:           pixel.IM,
		pos:           pixel.ZV,
		bounds:        pixel.R(0, 0, 0, 0),
		origin:        pixel.ZV,
		scale:         pixel.V(1, 1),
		rot:           0,
		rotpoint:      pixel.ZV,
		zindex:        0,
		zeroalignment: AlignmentBottomLeft,
		extraoffset:   pixel.ZV,
		styles:        DefaultStyles(),
	}
	b.Self = b
	b.calcMat()
	return b
}

func (b *BaseNode) calcMat() {
	b.mat = pixel.IM.Moved(b.extraoffset).ScaledXY(b.origin, b.scale).Rotated(b.rotpoint, b.rot).Moved(b.pos)
}

func (b *BaseNode) GetName() string {
	return b.Name
}

func (b *BaseNode) SetPos(p pixel.Vec) {
	if b.styles.RoundToPixels {
		b.pos = pixel.V(math.Round(p.X), math.Round(p.Y))
	} else {
		b.pos = p
	}
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

func (b *BaseNode) SetRotPoint(rotpoint pixel.Vec) {
	b.rotpoint = rotpoint
	b.calcMat()
}

func (b *BaseNode) GetRotPoint() pixel.Vec {
	return b.rotpoint
}

func (b *BaseNode) SetZIndex(z int) {
	b.zindex = z
}

func (b *BaseNode) SetExtraOffset(extra pixel.Vec) {
	b.extraoffset = pixel.V(math.Round(extra.X), math.Round(extra.Y))
	b.calcMat()
}

func (b *BaseNode) GetExtraOffset() pixel.Vec {
	return b.extraoffset
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

func (b *BaseNode) AddChild(child Node) {
	b.children = append(b.children, child)
	sort.SliceStable(b.children, func(i, j int) bool {
		less := (b.children[i]._getZindex() < b.children[j]._getZindex())
		return less
	})
	child._init()
}

func (b *BaseNode) Children() []Node {
	return b.children
}

func (b *BaseNode) SetStyles(styles *Styles) {
	b.styles = styles
}

func (b *BaseNode) GetStyles() *Styles {
	return b.styles
}

func (b *BaseNode) SetSize(size pixel.Vec) {}

func (b *BaseNode) Contains(point pixel.Vec) bool {
	return false
}
