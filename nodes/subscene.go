package nodes

import (
	"github.com/faiface/pixel"
)

type SubScene struct {
	Canvas
	root Node
	size pixel.Vec
}

func NewSubScene(name string, w, h float64) *SubScene {
	s := &SubScene{
		Canvas: *NewCanvas(name, w, h),
		size:   pixel.V(w, h),
	}
	s.Self = s
	s.canvas.SetBounds(pixel.R(-w/2, -h/2, w/2, h/2))
	return s
}

func (s *SubScene) SetRoot(root Node) {
	if s.root != nil {
		root._unmount()
	}
	root._init()
	s.root = root
	s.root._mount()
	s.root._updateFromTheme(SceneManager().Theme())
}

func (s *SubScene) SetSize(size pixel.Vec) {
	s.Canvas.SetSize(size)
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
		s.root.SetLastMat(mat)
	}
}

/*func (s *SubScene) Mount() {
	s.Canvas.Mount()
	w, h := s.size.XY()
	s.Canvas.SetSize(s.size)
	s.canvas.SetBounds(pixel.R(-w/2, -h/2, w/2, h/2))
}

func (s *SubScene) Unmount() {
	s.Canvas.Unmount()
}*/
