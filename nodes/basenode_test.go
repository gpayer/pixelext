package nodes

import (
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
	n.self = n
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
	m.self = m
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

	m.AddChild(b)
	m.AddChild(c)
	m.AddChild(a)

	r1 := m.children[0].(*BaseNode)
	r2 := m.children[1].(*BaseNode)
	r3 := m.children[2].(*BaseNode)
	assert.Equal(r1.Name, "a")
	assert.Equal(r2.Name, "b")
	assert.Equal(r3.Name, "c")
}
