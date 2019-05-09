package nodes

import (
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var sceneManager *SceneManagerStruct

func SceneManager() *SceneManagerStruct {
	return sceneManager
}

type SceneManagerStruct struct {
	root        Node
	last        time.Time
	first       bool
	redraw      bool
	clearColor  color.Color
	redrawCount int
	win         *pixelgl.Window
}

func (s *SceneManagerStruct) SetWin(win *pixelgl.Window) {
	s.win = win
	Events().setWin(win)
}

func (s *SceneManagerStruct) Win() *pixelgl.Window {
	return s.win
}

func (s *SceneManagerStruct) Redraw() {
	s.redraw = true
}

func (s *SceneManagerStruct) NeedsRedraw() bool {
	return s.redraw
}

func (s *SceneManagerStruct) SetRoot(root Node) {
	if s.root != nil {
		s.root._unmount()
	}
	root._init()
	s.root = root
	root._mount()
	s.Redraw()
}

func (s *SceneManagerStruct) Run(mat pixel.Matrix) {
	if s.first {
		s.last = time.Now()
		s.first = false
		s.redraw = true
	}
	dt := time.Since(s.last).Seconds()
	s.last = time.Now()
	s.root._update(dt)
	if s.redrawCount < 60 { // dirty hack
		s.redraw = true
		s.redrawCount++
	}
	if s.redraw {
		s.win.Clear(s.clearColor)
		s.root._draw(Events().win, mat)
		s.win.Update()
		s.redraw = false
	} else {
		s.win.UpdateInput()
		sleepRemaining := time.Until(s.last.Add(17 * time.Millisecond))
		if sleepRemaining > 0 {
			time.Sleep(sleepRemaining)
		}
	}
}

func init() {
	sceneManager = &SceneManagerStruct{first: true, clearColor: color.RGBA{0, 0, 0, 1}}
}
