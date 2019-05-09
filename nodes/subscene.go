package nodes

import (
	"github.com/faiface/pixel"
)

type SubScene struct {
	Canvas
	root Node
}

func NewSubScene(name string, w, h float64) *SubScene {
	s := &SubScene{
		Canvas: *NewCanvas(name, w, h),
	}
	s.Self = s
	s.canvas.SetBounds(pixel.R(-w/2, -h/2, w/2, h/2))
	return s
}

func (s *SubScene) SetRoot(root Node) {
	if s.root != nil {
		root._unmount()
	}
	s.root = root
	s.root._init()
	s.root._mount()
}

func (s *SubScene) SetSize(size pixel.Vec) {
	s.canvas.SetBounds(pixel.R(-size.X/2, -size.Y/2, size.X/2, size.Y/2))
	SceneManager().Redraw()
}

func (s *SubScene) Update(dt float64) {
	if s.root != nil {
		s.root._update(dt)
	}
}

func (s *SubScene) Draw(win pixel.Target, mat pixel.Matrix) {
	s.Clear(s.GetStyles().Background.Color)
	if s.root != nil {
		s.root._draw(s.GetCanvas(), pixel.IM)
	}
	s.Canvas.Draw(win, mat)
	if s.root != nil {
		s.root._setLastMat(mat)
	}
}
