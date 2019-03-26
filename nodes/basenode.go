package nodes

import (
	"sort"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type BaseNode struct {
	Self                      nodeInternal
	Name                      string
	children                  []nodeInternal
	show, active, initialized bool
	mat                       pixel.Matrix
	pos                       pixel.Vec
	bounds                    pixel.Rect
	scale                     pixel.Vec
	rot                       float64
	rotpoint                  pixel.Vec
	origin                    pixel.Vec
	zindex                    int
	zeroalignment             Alignment
	Extraoffset               pixel.Vec
}

func (b *BaseNode) _getMat() pixel.Matrix {
	return b.mat
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
	if !b.active || !b.show {
		return
	}
	for _, child := range b.children {
		child._draw(win, mat.Chained(child._getMat()))
	}
	drawable, ok := b.Self.(Drawable)
	if ok {
		drawable.Draw(win, mat)
	}
}

func NewBaseNode(name string) *BaseNode {
	b := &BaseNode{
		Name:          name,
		children:      make([]nodeInternal, 0),
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
		Extraoffset:   pixel.ZV,
	}
	b.Self = b
	b.calcMat()
	return b
}

func (b *BaseNode) calcMat() {
	relOrigin := b.origin.Sub(b.bounds.Min)
	relRotpoint := b.rotpoint.Sub(b.bounds.Min)
	relPos := b.pos.Sub(relOrigin)
	b.mat = pixel.IM.ScaledXY(relOrigin, b.scale).Rotated(relRotpoint, b.rot).Moved(relPos)
}

func (b *BaseNode) calcZero() {
	if b.zeroalignment != AlignmentFixed {
		whalf := b.bounds.W() / 2
		hhalf := b.bounds.H() / 2
		blAligned := b.bounds.Moved(b.bounds.Min.Scaled(-1)).Moved(b.Extraoffset)
		switch b.zeroalignment {
		case AlignmentBottomLeft:
			b.bounds = blAligned
		case AlignmentBottomCenter:
			b.bounds = blAligned.Moved(pixel.V(-whalf, 0))
		case AlignmentBottomRight:
			b.bounds = blAligned.Moved(pixel.V(-2*whalf, 0))
		case AlignmentCenterLeft:
			b.bounds = blAligned.Moved(pixel.V(0, -hhalf))
		case AlignmentCenter:
			b.bounds = blAligned.Moved(pixel.V(-whalf, -hhalf))
		case AlignmentCenterRight:
			b.bounds = blAligned.Moved(pixel.V(-2*whalf, -hhalf))
		case AlignmentTopLeft:
			b.bounds = blAligned.Moved(pixel.V(0, -2*hhalf))
		case AlignmentTopCenter:
			b.bounds = blAligned.Moved(pixel.V(-whalf, -2*hhalf))
		case AlignmentTopRight:
			b.bounds = blAligned.Moved(pixel.V(-2*whalf, -2*hhalf))
		default:
		}
		b.calcMat()
	}
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

func (b *BaseNode) SetRotPoint(rotpoint pixel.Vec) {
	b.rotpoint = rotpoint
	b.calcMat()
}

func (b *BaseNode) GetRotPoint() pixel.Vec {
	return b.rotpoint
}

func (b *BaseNode) SetBounds(r pixel.Rect) {
	b.bounds = r
	b.calcZero()
	b.calcMat()
}

func (b *BaseNode) GetBounds() pixel.Rect {
	return b.bounds
}

// SetZeroAlignment decides how to handle element boundaries. Unless the alignment is AlignementFixed
// the Min and Max values of the boundaries Rect are recalculated to have the zero point at the
// requested position.
func (b *BaseNode) SetZeroAlignment(a Alignment) {
	b.zeroalignment = a
	b.calcZero()
}

func (b *BaseNode) SetZIndex(z int) {
	b.zindex = z
}

func (b *BaseNode) SetExtraOffset(extra pixel.Vec) {
	b.Extraoffset = extra
	b.calcZero()
}

func (b *BaseNode) GetExtraOffset() pixel.Vec {
	return b.Extraoffset
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
	sort.SliceStable(b.children, func(i, j int) bool {
		less := (b.children[i]._getZindex() < b.children[j]._getZindex())
		return less
	})
	child._init()
}
