package nodes

import (
	"sort"
	"testing"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/stretchr/testify/assert"
)

type nakedNode struct {
	BaseNode
}

func newNakedNode() *nakedNode {
	n := &nakedNode{
		BaseNode: *NewBaseNode("test"),
	}
	n.Self = n
	return n
}

type mockNode struct {
	BaseNode
	success, mounted, unmounted, updated, drawn bool
}

func newMockNode() *mockNode {
	m := &mockNode{
		BaseNode:  *NewBaseNode("test"),
		success:   false,
		mounted:   false,
		unmounted: false,
		updated:   false,
		drawn:     false,
	}
	m.Self = m
	return m
}

func (m *mockNode) Init() {
	m.success = true
}

func (m *mockNode) Mount() {
	m.mounted = true
}

func (m *mockNode) Unmount() {
	m.unmounted = true
}

func (m *mockNode) Update(dt float64) {
	m.updated = true
}

func (m *mockNode) Draw(win *pixelgl.Window, mat pixel.Matrix) {
	panic("not implemented")
}

func TestInit(t *testing.T) {
	assert := assert.New(t)
	m := newMockNode()
	m._init()
	n := newNakedNode()
	n._init()

	assert.True(m.success)
}

func TestMountable(t *testing.T) {
	assert := assert.New(t)
	m := newMockNode()
	m._mount()
	n := newNakedNode()
	n._mount()
	m._unmount()
	n._unmount()

	assert.True(m.mounted)
	assert.True(m.unmounted)
}

func TestUpdate(t *testing.T) {
	assert := assert.New(t)
	m := newMockNode()
	m._update(1)
	n := newNakedNode()
	n._update(1)

	assert.True(m.updated)

	m.active = false
	m.updated = false
	m._update(1)
	assert.False(m.updated)
}

func TestChildren(t *testing.T) {
	assert := assert.New(t)
	m := newMockNode()
	a := NewBaseNode("a")
	a.zindex = 0
	b := NewBaseNode("b")
	b.zindex = 1
	c := NewBaseNode("c")
	c.zindex = 2

	m.AddChild(c)
	m.AddChild(a)
	m.AddChild(b)

	r1 := m.children[0].(*BaseNode)
	r2 := m.children[1].(*BaseNode)
	r3 := m.children[2].(*BaseNode)
	assert.Equal("a", r1.Name)
	assert.Equal("b", r2.Name)
	assert.Equal("c", r3.Name)
}

func TestSortSlice(t *testing.T) {
	assert := assert.New(t)
	sl := []int{0, 3, 6, -1, 23, 2}
	sort.Slice(sl, func(i, j int) bool {
		return sl[i] < sl[j]
	})
	assert.Equal(sl, []int{-1, 0, 2, 3, 6, 23})
}

func TestContainerBounds(t *testing.T) {
	assert := assert.New(t)
	c := NewBaseNode("container")
	n1 := NewBaseNode("n1")
	n1.SetBounds(pixel.R(0, 0, 100, 100))
	n1.SetPos(pixel.V(0, 0))
	c.AddChild(n1)

	cr := c.GetContainerBounds()
	assert.Equal(pixel.R(0, 0, 100, 100), cr)

	n2 := NewBaseNode("n2")
	n2.SetBounds(pixel.R(0, 0, 20, 20))
	n2.SetOrigin(pixel.V(10, 10))
	n2.SetPos(pixel.V(100, 100))
	c.AddChild(n2)

	cr = c.GetContainerBounds()
	assert.Equal(pixel.R(0, 0, 110, 110), cr)
}

func TestContainerBounds2(t *testing.T) {
	assert := assert.New(t)
	c := NewBaseNode("container")
	n1 := NewBaseNode("n1")
	n1.SetBounds(pixel.R(0, 0, 100, 100))
	n1.SetPos(pixel.V(0, 0))
	n1.SetZeroAlignment(AlignmentCenter)
	c.AddChild(n1)

	cr := c.GetContainerBounds()
	assert.Equal(pixel.R(-50, -50, 50, 50), cr)

	n2 := NewBaseNode("n2")
	n2.SetBounds(pixel.R(0, 0, 20, 20))
	n2.SetPos(pixel.V(100, 100))
	n2.SetZeroAlignment(AlignmentTopRight)
	c.AddChild(n2)

	cr = c.GetContainerBounds()
	assert.Equal(pixel.R(-50, -50, 100, 100), cr)
}
