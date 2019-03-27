package nodes

import (
	"time"

	"github.com/faiface/pixel"
)

var sceneManager *SceneManagerStruct

func SceneManager() *SceneManagerStruct {
	return sceneManager
}

type SceneManagerStruct struct {
	root  Node
	last  time.Time
	first bool
}

func (s *SceneManagerStruct) SetRoot(root Node) {
	if s.root != nil {
		s.root._unmount()
	}
	root._init()
	s.root = root
	root._mount()
}

func (s *SceneManagerStruct) Run(mat pixel.Matrix) {
	if s.first {
		s.last = time.Now()
		s.first = false
	}
	dt := time.Since(s.last).Seconds()
	s.last = time.Now()
	s.root._update(dt)
	s.root._draw(Events().win, mat)
}

func init() {
	sceneManager = &SceneManagerStruct{first: true}
}
