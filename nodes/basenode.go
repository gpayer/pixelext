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
	parent                    Node
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
	locked                    bool
	drawable                  Drawable
	updateable                Updateable
	mounted                   bool
}

func (b *BaseNode) _getMat() pixel.Matrix {
	return b.mat
}

func (b *BaseNode) GetLastMat() pixel.Matrix {
	return b.lastmat
}

func (b *BaseNode) _setLastMat(mat pixel.Matrix) {
	b.lastmat = mat
	for _, child := range b.children {
		if child == nil {
			continue
		}
		child._setLastMat(child._getMat().Chained(mat))
	}
}

func (b *BaseNode) _init() {
	if b.initialized {
		return
	}
	for _, child := range b.children {
		if child == nil {
			continue
		}
		child._init()
	}
	init, ok := b.Self.(Initializable)
	if ok {
		init.Init()
	}
	b.initialized = true

	updateable, ok := b.Self.(Updateable)
	if ok {
		b.updateable = updateable
	}
	drawable, ok := b.Self.(Drawable)
	if ok {
		b.drawable = drawable
	}
}

func (b *BaseNode) _mount() {
	for _, child := range b.children {
		if child == nil {
			continue
		}
		child._mount()
	}
	if b.mounted {
		return
	}
	b.Self.UpdateFromTheme(SceneManager().Theme())
	mountable, ok := b.Self.(Mountable)
	if ok {
		mountable.Mount()
	}
	b.mounted = true
	SceneManager().Redraw()
}

func (b *BaseNode) _unmount() {
	for _, child := range b.children {
		if child == nil {
			continue
		}
		child._unmount()
	}
	if !b.mounted {
		return
	}
	mountable, ok := b.Self.(Mountable)
	if ok {
		mountable.Unmount()
	}
	b.mounted = false
}

func (b *BaseNode) _update(dt float64) {
	if !b.active {
		return
	}
	for i := len(b.children) - 1; i >= 0; i-- {
		child := b.children[i]
		if child == nil {
			continue
		}
		child._update(dt)
	}
	if b.paused {
		return
	}

	if b.updateable != nil {
		b.updateable.Update(dt)
	}

	sortchildren := false

	newsize := 0
	for _, child := range b.children {
		if child == nil {
			sortchildren = true
		} else {
			newsize++
		}
	}

	if sortchildren {
		newchildren := make([]Node, newsize)
		cur := 0
		for _, child := range b.children {
			if child != nil {
				newchildren[cur] = child
				cur++
			}
		}
		b.children = newchildren
		b.SortChildren()
	}
}

func (b *BaseNode) _draw(win pixel.Target, mat pixel.Matrix) {
	b.lastmat = mat
	if !b.active || !b.show {
		return
	}

	if b.drawable != nil {
		b.drawable.Draw(win, mat)
	}
	for _, child := range b.children {
		if child == nil {
			continue
		}
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
		pausable:      true,
		paused:        false,
		locked:        false,
		mounted:       false,
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

func (b *BaseNode) Parent() Node {
	return b.parent
}

func (b *BaseNode) SetParent(p Node) {
	b.parent = p
}

func (b *BaseNode) SetName(name string) {
	b.Name = name
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

func (b *BaseNode) LocalToGlobalPos(local pixel.Vec) pixel.Vec {
	return b.GetLastMat().Project(local.Sub(b.Self.GetExtraOffset()))
}

func (b *BaseNode) GlobalToLocalPos(global pixel.Vec) pixel.Vec {
	if b.parent == nil {
		return global
	}
	return b.parent.GetLastMat().Unproject(global).Add(b.Self.GetExtraOffset())
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

func (b *BaseNode) ZIndex() int {
	return b.zindex
}

func (b *BaseNode) SetZIndex(z int) {
	b.zindex = z
	if b.parent != nil {
		b.parent.SortChildren()
	}
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

func (b *BaseNode) SortChildren() {
	sort.SliceStable(b.children, func(i, j int) bool {
		if b.children[i] == nil {
			return true
		} else if b.children[j] == nil {
			return false
		}
		less := (b.children[i].ZIndex() < b.children[j].ZIndex())
		return less
	})
}

func (b *BaseNode) AddChild(child Node) {
	b.children = append(b.children, child)
	child.SetParent(b.Self)
	b.SortChildren()
	child._init()
	child._mount()
	SceneManager().Redraw()
}

func (b *BaseNode) RemoveChild(child Node) {
	for i, ch := range b.children {
		if child == ch {
			ch._unmount()
			b.children[i] = nil
			SceneManager().Redraw()
			break
		}
	}
}

func (b *BaseNode) RemoveChildren() {
	for i, ch := range b.children {
		if ch == nil || ch.Locked() {
			continue
		}
		ch._unmount()
		b.children[i] = nil
	}
	SceneManager().Redraw()
}

func (b *BaseNode) Children() []Node {
	return b.children
}

func (b *BaseNode) ChildChanged() {
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
		if ch == nil {
			continue
		}
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
			if c == nil {
				continue
			}
			c.Pause()
		}
	}
}

func (b *BaseNode) Unpause() {
	if b.pausable {
		b.paused = false
		for _, c := range b.children {
			if c == nil {
				continue
			}
			c.Unpause()
		}
	}
}

func (b *BaseNode) SetLocked(locked bool) {
	b.locked = locked
}

func (b *BaseNode) Locked() bool {
	return b.locked
}

func (b *BaseNode) Iterate(fn func(n Node)) {
	if b.Locked() {
		return
	}
	fn(b)
	for _, c := range b.Self.Children() {
		c.Iterate(fn)
	}
}
