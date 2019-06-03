package nodes

import (
	"math"
	"sort"

	"github.com/faiface/pixel"
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
	pausable, paused          bool
	remove                    bool
}

func (b *BaseNode) _getMat() pixel.Matrix {
	return b.mat
}

func (b *BaseNode) _getLastMat() pixel.Matrix {
	return b.lastmat
}

func (b *BaseNode) _setLastMat(mat pixel.Matrix) {
	b.lastmat = mat
	for _, child := range b.children {
		child._setLastMat(child._getMat().Chained(mat))
	}
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
	SceneManager().Redraw()
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
	if b.paused {
		return
	}
	updateable, ok := b.Self.(Updateable)
	if ok {
		updateable.Update(dt)
	}

	sortchildren := false
restart:
	for i, child := range b.children {
		if child.IsRemove() {
			b.children[i] = b.children[len(b.children)-1]
			b.children[len(b.children)-1] = nil
			b.children = b.children[:len(b.children)-1]
			newChildren := make([]Node, len(b.children)-1)
			copy(newChildren, b.children[:len(b.children)-1])
			b.children = newChildren
			child._unmount()
			sortchildren = true
			goto restart
		}
	}

	if sortchildren {
		b.sortChildren()
	}
}

func (b *BaseNode) _draw(win pixel.Target, mat pixel.Matrix) {
	b.lastmat = mat
	if !b.active || !b.show || b.remove {
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

func (b *BaseNode) IsRemove() bool {
	return b.remove
}

func (b *BaseNode) SetRemove(remove bool) {
	b.remove = remove
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
		pausable:      true,
		paused:        false,
		remove:        false,
	}
	b.Self = b
	b.calcMat()
	return b
}

func (b *BaseNode) calcMat() {
	b.mat = pixel.IM.Moved(b.extraoffset).ScaledXY(b.origin, b.scale).Rotated(b.rotpoint, b.rot).Moved(b.pos)
	SceneManager().Redraw()
}

func (b *BaseNode) Initialized() bool {
	return b.initialized
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
	SceneManager().Redraw()
}

func (b *BaseNode) SetExtraOffset(extra pixel.Vec) {
	b.extraoffset = pixel.V(math.Round(extra.X), math.Round(extra.Y))
	b.calcMat()
}

func (b *BaseNode) GetExtraOffset() pixel.Vec {
	return b.extraoffset
}

func (b *BaseNode) Show() {
	if !b.show {
		b.show = true
		SceneManager().Redraw()
	}
	b.active = true
}

func (b *BaseNode) Hide() {
	if b.show {
		b.show = false
		SceneManager().Redraw()
	}
	b.active = false
}

func (b *BaseNode) SetActive(active bool) {
	b.active = active
	SceneManager().Redraw()
}

func (b *BaseNode) sortChildren() {
	sort.SliceStable(b.children, func(i, j int) bool {
		less := (b.children[i]._getZindex() < b.children[j]._getZindex())
		return less
	})
}

func (b *BaseNode) AddChild(child Node) {
	child.SetRemove(false)
	b.children = append(b.children, child)
	b.sortChildren()
	child._init()
	SceneManager().Redraw()
}

func (b *BaseNode) RemoveChild(child Node) {
	for _, ch := range b.children {
		if child == ch {
			child.SetRemove(true)
			SceneManager().Redraw()
			break
		}
	}
}

func (b *BaseNode) RemoveChildren() {
	for _, ch := range b.children {
		ch.SetRemove(true)
	}
	SceneManager().Redraw()
}

func (b *BaseNode) Children() []Node {
	return b.children
}

func (b *BaseNode) SetStyles(styles *Styles) {
	b.styles = styles
	SceneManager().Redraw()
}

func (b *BaseNode) GetStyles() *Styles {
	return b.styles
}

func (b *BaseNode) SetSize(size pixel.Vec) {}

func (b *BaseNode) Contains(point pixel.Vec) bool {
	return false
}

func (b *BaseNode) _updateFromTheme(theme *Theme) {
	if theme == nil {
		return
	}
	b.Self.UpdateFromTheme(theme)
	for _, ch := range b.children {
		ch._updateFromTheme(theme)
	}
}

func (b *BaseNode) UpdateFromTheme(theme *Theme) {}

func (b *BaseNode) SetPausable(pausable bool) {
	b.pausable = pausable
}

func (b *BaseNode) Pause() {
	if b.pausable {
		b.paused = true
		for _, c := range b.children {
			c.Pause()
		}
	}
}

func (b *BaseNode) Unpause() {
	if b.pausable {
		b.paused = false
		for _, c := range b.children {
			c.Unpause()
		}
	}
}
