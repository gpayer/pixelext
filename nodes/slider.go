package nodes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Slider struct {
	BaseNode
	min, max, current float32
	dirty             bool
	canvas            *pixelgl.Canvas
	onchange          func(v float32)
}

func (s *Slider) Init() {
	s.canvas = pixelgl.NewCanvas(s.GetBounds())
	s.dirty = true
	s.SetExtraOffset(pixel.V(s.GetBounds().W()/2, s.GetBounds().H()/2))
}

func (s *Slider) Mounted() {
}

func (s *Slider) Unmounted() {
}

func (s *Slider) Update(dt float64) {
	ev := Events()
	bounds := s.GetBounds()
	if ev.Clicked(pixelgl.MouseButtonLeft, s) {
		pos := ev.LocalMousePosition(s)
		zblPos := pos.Sub(bounds.Min)
		//fmt.Printf("clicked: %v %v %v %v\n", s.extraoffset, bounds, pos, zblPos)
		s.current = s.min + (s.max-s.min)*float32(zblPos.X/bounds.W())
		s.dirty = true
	} else if ev.MouseScroll().Y != 0 {
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
		bounds = s.canvas.Bounds()
		s.onchange(s.current)
		s.canvas.Clear(colornames.Black)
		im := imdraw.New(nil)
		im.Color = colornames.White
		min := bounds.Min
		max := bounds.Max
		im.Push(min, pixel.V(min.X, max.Y), max, pixel.V(max.X, min.Y))
		im.Polygon(2)
		currentw := (bounds.W() - 2) * float64(s.current/(s.max-s.min))
		im.Color = colornames.Skyblue
		innerOrig := min.Add(pixel.V(1, 1))
		im.Push(innerOrig, innerOrig.Add(pixel.V(currentw, 0)), innerOrig.Add(pixel.V(currentw, bounds.H()-2)), innerOrig.Add(pixel.V(0, bounds.H()-2)))
		im.Polygon(0)
		im.Draw(s.canvas)
		s.dirty = false
	}
}

func (s *Slider) Draw(win *pixelgl.Window, mat pixel.Matrix) {
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
		BaseNode: *NewBaseNode(name),
		min:      min, max: max, current: current,
		onchange: func(_ float32) {},
	}
	sl.Self = sl
	return sl
}
