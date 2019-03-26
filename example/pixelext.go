package main

import (
	"fmt"
	"math"
	"pixelext/nodes"
	"pixelext/services"

	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type demo struct {
	nodes.BaseNode
	text1, text2, text3 *nodes.Text
	sprite              *nodes.Sprite
	rotslider           *nodes.BaseNode
}

func newDemo() *demo {
	d := &demo{BaseNode: *nodes.NewBaseNode("demo")}
	d.Self = d
	return d
}

func makeText(x, y float64, name string, al nodes.Alignment) *nodes.Text {
	text := nodes.NewText(name, "basic")
	text.Printf("ABCDEF")
	text.SetPos(pixel.V(x, y))
	text.SetZeroAlignment(al)
	return text
}

func (d *demo) Init() {
	text := makeText(100, 100, "text1", nodes.AlignmentBottomLeft)
	d.AddChild(text)
	d.text1 = text

	text = makeText(100, 200, "text2", nodes.AlignmentBottomCenter)
	d.AddChild(text)

	text = makeText(100, 300, "text2", nodes.AlignmentBottomRight)
	d.AddChild(text)

	text = makeText(100, 400, "text2", nodes.AlignmentCenterLeft)
	d.AddChild(text)

	text = makeText(100, 500, "text2", nodes.AlignmentCenter)
	d.AddChild(text)
	d.text2 = text

	text = makeText(200, 100, "text2", nodes.AlignmentCenterRight)
	d.AddChild(text)

	text = makeText(200, 200, "text2", nodes.AlignmentTopLeft)
	d.AddChild(text)

	text = makeText(200, 300, "text2", nodes.AlignmentTopCenter)
	d.AddChild(text)

	text = makeText(200, 400, "text2", nodes.AlignmentTopRight)
	d.AddChild(text)
	d.text3 = text

	sltext := nodes.NewBaseNode("sltext")
	sltext.SetBounds(pixel.R(0, 0, 100, 30))
	sltext.SetPos(pixel.V(400, 500))
	sltext.SetRot(0.0)
	d.AddChild(sltext)

	slidertext := makeText(50, 15, "slidertext", nodes.AlignmentCenter)
	slidertext.SetZIndex(10)
	sltext.AddChild(slidertext)
	d.rotslider = sltext

	slider := nodes.NewSlider("slider1", 0, 1, 0.5)
	slider.SetBounds(pixel.R(0, 0, 100, 30))
	slider.SetPos(pixel.V(0, 0))
	slider.OnChange(func(v float32) {
		slidertext.Text().Clear()
		slidertext.Printf("%.2f", v)
		fmt.Printf("new value: %.2f\n", v)
	})
	sltext.AddChild(slider)

	slider = nodes.NewSlider("slider2", 0, 1, 0.5)
	slider.SetBounds(pixel.R(0, 0, 100, 30))
	slider.SetPos(pixel.V(400, 400))
	slider.SetZeroAlignment(nodes.AlignmentCenter)
	d.AddChild(slider)

	slider = nodes.NewSlider("slider3", 0, 1, 0.5)
	slider.SetBounds(pixel.R(0, 0, 50, 20))
	slider.SetPos(pixel.V(500, 300))
	slider.SetZeroAlignment(nodes.AlignmentTopRight)
	d.AddChild(slider)

	pic, err := services.ResourceManager().LoadPicture("gopher.png")
	if err != nil {
		panic(err)
	}
	sprite := nodes.NewSprite("sprite1", pic)
	sprite.SetZeroAlignment(nodes.AlignmentCenter)
	sprite.SetPos(pixel.V(600, 300))
	d.AddChild(sprite)
	d.sprite = sprite
}

func (d *demo) Update(dt float64) {
	dphi := math.Pi * dt
	d.text1.SetRot(d.text1.GetRot() + dphi)
	d.text2.SetRot(d.text2.GetRot() + dphi)
	d.text3.SetRot(d.text3.GetRot() + dphi)
	d.sprite.SetRot(d.sprite.GetRot() + dphi)

	dphi = math.Pi / 5.0 * dt
	if nodes.Events().Pressed(pixelgl.KeyA) {
		d.rotslider.SetRot(d.rotslider.GetRot() + dphi)
	} else if nodes.Events().Pressed(pixelgl.KeyD) {
		d.rotslider.SetRot(d.rotslider.GetRot() - dphi)
	}
}

func Run() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "pixelext example",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
	})
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)
	nodes.Events().SetWin(win)

	root := newDemo()

	nodes.SceneManager().SetRoot(root)

	im := imdraw.New(nil)
	im.SetColorMask(colornames.Gray)
	var i float64
	for i = 100; i < 800; i += 100 {
		im.Push(pixel.V(i, 0), pixel.V(i, 600))
		im.Line(1)
	}
	for i = 100; i < 600; i += 100 {
		im.Push(pixel.V(0, i), pixel.V(800, i))
		im.Line(1)
	}

	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyEscape) {
			break
		}
		win.Clear(colornames.Black)
		im.Draw(win)
		nodes.SceneManager().Run(pixel.IM)
		win.Update()
	}
}

func main() {
	pixelgl.Run(Run)
}
