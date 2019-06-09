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
	theme       *Theme
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
	root._updateFromTheme(s.theme)
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
	Events().reset()
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

func (s *SceneManagerStruct) SetTheme(theme *Theme) {
	s.theme = theme
	if s.root != nil {
		s.root._updateFromTheme(theme)
	}
	s.redraw = true
}

func (s *SceneManagerStruct) Theme() *Theme {
	return s.theme
}

func (s *SceneManagerStruct) SetClearColor(color color.Color) {
	s.clearColor = color
}

func (s *SceneManagerStruct) Pause() {
	if s.root != nil {
		s.root.Pause()
	}
}

func (s *SceneManagerStruct) Unpause() {
	if s.root != nil {
		s.root.Unpause()
	}
}

func init() {
	sceneManager = &SceneManagerStruct{first: true, clearColor: color.RGBA{0, 0, 0, 255}, theme: DefaultTheme()}
}
