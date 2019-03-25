package main

import (
	"math"
	"pixelext/nodes"

	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type demo struct {
	nodes.BaseNode
	text1, text2, text3 *nodes.Text
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
}

func (d *demo) Update(dt float64) {
	dphi := math.Pi * dt
	d.text1.SetRot(d.text1.GetRot() + dphi)
	d.text2.SetRot(d.text2.GetRot() + dphi)
	d.text3.SetRot(d.text3.GetRot() + dphi)
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
