package ui

import (
	"github.com/gpayer/pixelext/nodes"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type Slider struct {
	UIBase
	min, max, current float32
	dirty             bool
	canvas            *pixelgl.Canvas
	onchange          func(v float32)
}

func (s *Slider) Init() {
	size := s.Size()
	s.canvas = pixelgl.NewCanvas(pixel.R(0, 0, size.X, size.Y))
	s.dirty = true
}

func (s *Slider) Mounted() {
}

func (s *Slider) Unmounted() {
}

func (s *Slider) Update(dt float64) {
	ev := nodes.Events()
	size := s.Size()
	if ev.Clicked(pixelgl.MouseButtonLeft, s) {
		pos := ev.LocalMousePosition(s)
		xPos := pos.X + size.X/2
		//fmt.Printf("clicked: %v %v %v %v\n", s.GetExtraOffset(), size, pos, zblPos)
		s.current = s.min + (s.max-s.min)*float32(xPos/size.X)
		s.dirty = true
	} else if ev.MouseScroll().Y != 0 && ev.MouseHovering(s) {
		s.current += float32(0.1 * ev.MouseScroll().Y)
		if s.current > s.max {
			s.current = s.max
		}
		if s.current < s.min {
			s.current = s.min
		}
		s.dirty = true
	}
	if s.dirty {
		styles := s.GetStyles()
		bounds := s.canvas.Bounds()
		s.onchange(s.current)
		s.canvas.Clear(styles.Background.Color)
		im := imdraw.New(nil)
		im.Color = styles.Border.Color
		min := bounds.Min
		max := bounds.Max
		bw := styles.Border.Width
		if bw > 0 {
			im.Push(min, pixel.V(min.X, max.Y), max, pixel.V(max.X, min.Y))
			im.Polygon(bw)
		}
		currentw := (bounds.W() - 2) * float64(s.current/(s.max-s.min))
		im.Color = styles.Foreground.Color
		innerOrig := min.Add(pixel.V(bw/2, bw/2))
		im.Push(innerOrig, innerOrig.Add(pixel.V(currentw, 0)), innerOrig.Add(pixel.V(currentw, bounds.H()-bw)), innerOrig.Add(pixel.V(0, bounds.H()-bw)))
		im.Polygon(0)
		im.Draw(s.canvas)
		s.dirty = false
		nodes.SceneManager().Redraw()
	}
}

func (s *Slider) SetStyles(styles *nodes.Styles) {
	s.BaseNode.SetStyles(styles)
	s.dirty = true
}

func (s *Slider) Draw(win pixel.Target, mat pixel.Matrix) {
	s.canvas.Draw(win, mat)
}

func (s *Slider) SetValue(v float32) {
	if v >= s.min && v <= s.max {
		s.current = v
		s.dirty = true
	}
}

func (s *Slider) Value() float32 {
	return s.current
}

func (s *Slider) OnChange(fn func(v float32)) {
	s.onchange = fn
}

func NewSlider(name string, min, max, current float32) *Slider {
	sl := &Slider{
		UIBase: *NewUIBase(name),
		min:    min, max: max, current: current,
		onchange: func(_ float32) {},
	}
	sl.Self = sl
	sl.UISelf = sl
	return sl
}
