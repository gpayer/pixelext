package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"pixelext/nodes"
	"pixelext/services"
	"pixelext/ui"
	"runtime"
	"runtime/pprof"

	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type demo struct {
	nodes.BaseNode
	text1, text2, text3 *ui.Text
	sprite              *nodes.Sprite
	rotslider           *nodes.BaseNode
	sprite2             *nodes.Sprite
	totalT              float64
}

func newDemo() *demo {
	d := &demo{BaseNode: *nodes.NewBaseNode("demo")}
	d.Self = d
	return d
}

type backgroundGrid struct {
	nodes.Canvas
}

func newBackgroundGrid(w, h float64) *backgroundGrid {
	b := &backgroundGrid{
		Canvas: *nodes.NewCanvas("backgroundGrid", w, h),
	}
	b.Self = b
	return b
}

func (b *backgroundGrid) Init() {
	size := b.GetCanvas().Bounds().Size()
	im := imdraw.New(nil)
	im.SetColorMask(colornames.Gray)
	var i float64
	for i = 100; i < size.X; i += 100 {
		im.Push(pixel.V(i, 0), pixel.V(i, size.Y))
		im.Line(1)
	}
	for i = 100; i < size.Y; i += 100 {
		im.Push(pixel.V(0, i), pixel.V(size.X, i))
		im.Line(1)
	}
	im.Draw(b.GetCanvas())
}

func makeText(x, y float64, name string, al nodes.Alignment) *ui.Text {
	text := ui.NewText(name, "basic")
	text.Printf("ABCDEF\nqpfjXX")
	text.SetPos(pixel.V(x, y))
	text.SetAlignment(al)
	return text
}

func (d *demo) Init() {
	styles := nodes.DefaultStyles()
	styles.Text.Color = colornames.Chartreuse
	styles.Border.Width = 5
	styles.Border.Color = colornames.Fuchsia
	styles.Element.EnabledColor = colornames.Darkgreen
	styles.Background.Color = color.RGBA{R: 30, G: 30, B: 30, A: 128}

	text := makeText(100, 100, "text1", nodes.AlignmentBottomLeft)
	d.AddChild(text)
	d.text1 = text

	text = makeText(200, 100, "text2", nodes.AlignmentBottomCenter)
	d.AddChild(text)

	text = makeText(300, 100, "text3", nodes.AlignmentBottomRight)
	d.AddChild(text)

	text = makeText(100, 200, "text4", nodes.AlignmentCenterLeft)
	d.AddChild(text)

	text = makeText(200, 200, "text5", nodes.AlignmentCenter)
	d.AddChild(text)
	d.text2 = text

	text = makeText(300, 200, "text6", nodes.AlignmentCenterRight)
	text.Clear()
	text.Printf("ABC\nqpg\nXYZ")
	d.AddChild(text)

	text = makeText(100, 300, "text7", nodes.AlignmentTopLeft)
	d.AddChild(text)

	text = makeText(200, 300, "text8", nodes.AlignmentTopCenter)
	d.AddChild(text)

	text = makeText(300, 300, "text9", nodes.AlignmentTopRight)
	text.SetStyles(styles)
	d.AddChild(text)
	d.text3 = text

	sltext := nodes.NewBaseNode("sltext")
	//sltext.SetSize(pixel.V(100, 30))
	sltext.SetPos(pixel.V(450, 515))
	sltext.SetRot(0.0)
	d.AddChild(sltext)

	slidertext := makeText(0, 0, "slidertext", nodes.AlignmentCenter)
	slidertext.SetZIndex(10)
	sltext.AddChild(slidertext)
	sltext.SetRotPoint(pixel.V(-50, -15))
	d.rotslider = sltext

	slider := ui.NewSlider("slider1", 0, 1, 0.5)
	slider.SetSize(pixel.V(100, 30))
	slider.SetAlignment(nodes.AlignmentCenter)
	slider.SetPos(pixel.V(0, 0))
	slider.OnChange(func(v float32) {
		slidertext.Text().Clear()
		slidertext.Printf("%.2f", v)
		fmt.Printf("new value: %.2f\n", v)
	})
	slider.SetStyles(styles)
	sltext.AddChild(slider)

	slider = ui.NewSlider("slider2", 0, 1, 0.5)
	slider.SetSize(pixel.V(100, 30))
	slider.SetPos(pixel.V(400, 400))
	slider.SetAlignment(nodes.AlignmentCenter)
	slider.OnChange(func(v float32) {
		fmt.Printf("slider2: new value: %.2f\n", v)
	})
	d.AddChild(slider)

	slider = ui.NewSlider("slider3", 0, 1, 0.5)
	slider.SetSize(pixel.V(50, 20))
	slider.SetPos(pixel.V(500, 300))
	slider.SetAlignment(nodes.AlignmentTopRight)
	slider.OnChange(func(v float32) {
		fmt.Printf("slider3: new value: %.2f\n", v)
	})
	d.AddChild(slider)

	pic, err := services.ResourceManager().LoadPicture("gopher.png")
	if err != nil {
		panic(err)
	}
	sprite := nodes.NewSprite("sprite1", pic)
	sprite.SetPos(pixel.V(600, 300))
	d.AddChild(sprite)
	d.sprite = sprite

	hbox := ui.NewHBox("hbox1")
	hbox.SetAlignment(nodes.AlignmentBottomLeft)
	hbox.SetPos(pixel.V(100, 700))
	hbox.SetStyles(styles)
	d.AddChild(hbox)

	text = ui.NewText("f1", "basic")
	text.Printf("Field1")
	hbox.AddChild(text)
	text = ui.NewText("f2", "basic")
	text.Printf("Field2\nLine2")
	txtstyle := text.GetStyles()
	txtstyle.Text.Color = colornames.Gold
	text.SetStyles(txtstyle)
	hbox.AddChild(text)
	text = ui.NewText("f3", "basic")
	text.Printf("Field3")
	hbox.AddChild(text)

	bbox := nodes.NewBorderBox("bbox", 50, 50)
	bbox.SetStyles(styles)
	bbox.SetPos(pixel.V(500, 700))
	bbox.SetZIndex(-1)
	d.AddChild(bbox)

	canvas := nodes.NewCanvas("canvas", 100, 100)
	canvas.Clear(colornames.Royalblue)
	im := imdraw.New(nil)
	im.Push(pixel.V(0, 2), pixel.V(100, 2))
	im.Line(4)
	im.Push(pixel.V(98, 0), pixel.V(98, 100))
	im.Line(4)
	im.Push(pixel.V(100, 98), pixel.V(0, 98))
	im.Line(4)
	im.Push(pixel.V(2, 0), pixel.V(2, 100))
	im.Line(4)
	im.Draw(canvas.Canvas())
	canvas.SetPos(pixel.V(200, 200))
	canvas.SetZIndex(-1)
	d.AddChild(canvas)

	button := ui.NewButton("btn1", 75, 40, "Click me!")
	button.SetAlignment(nodes.AlignmentBottomLeft)
	button.SetPos(pixel.V(710, 710))
	button.OnClick(func() {
		fmt.Println("Button clicked!")
	})
	d.AddChild(button)

	btngroup := ui.NewButtonGroup("btngroup", 0)
	btngroup.SetPos(pixel.V(710, 640))
	btngroup.AddButton("First Choice", "1", 0)
	btngroup.AddButton("Nr. 2", "2", 0)
	btngroup.AddButton("Three", "3", 0)
	btngroup.OnChange(func(v string) {
		fmt.Println(v)
	})
	d.AddChild(btngroup)

	grid := ui.NewGrid("grid", 3)
	grid.SetAlignment(nodes.AlignmentTopLeft)
	grid.SetPos(pixel.V(910, 590))
	d.AddChild(grid)

	text = ui.NewText("", "basic")
	text.Printf("One Field")
	grid.AddChild(text)

	button = ui.NewButton("btn1", 0, 0, "Btn!")
	grid.AddChild(button)

	text = ui.NewText("", "basic")
	text.Printf("Last")
	grid.AddChild(text)

	text = ui.NewText("", "basic")
	text.Printf("Next line")
	grid.AddChild(text)

	subscene := nodes.NewSubScene("subscene1", 100, 100)
	subscene.SetPos(pixel.V(1000, 200))
	subscene.SetRot(3.14 / 4)
	d.AddChild(subscene)

	subroot := nodes.NewBaseNode("subroot")
	d.sprite2 = nodes.NewSprite("sprite2", pic)
	d.sprite2.SetPos(pixel.V(0, 0))
	subroot.AddChild(d.sprite2)
	subscene.SetRoot(subroot)

	subscene = nodes.NewSubScene("subscene2", 100, 100)
	subscene.SetPos(pixel.V(1100, 300))
	subscene.SetRot(-3.14 / 4)
	d.AddChild(subscene)

	subroot = nodes.NewBaseNode("subroot2")

	slider = ui.NewSlider("subslider", 0, 1, 0.5)
	slider.SetSize(pixel.V(100, 30))
	slider.SetPos(pixel.V(0, -20))
	slider.SetAlignment(nodes.AlignmentCenter)
	slider.OnChange(func(v float32) {
		fmt.Printf("subslider: new value: %.2f\n", v)
	})
	subroot.AddChild(slider)
	subscene.SetRoot(subroot)

	vbox := ui.NewVBox("vbox1")
	vbox.SetPos(pixel.V(1100, 400))

	text = ui.NewText("", "basic")
	text.Printf("Line 1")
	vbox.AddChild(text)
	text = ui.NewText("", "basic")
	text.Printf("Line 2 is pretty long")
	vbox.AddChild(text)

	d.AddChild(vbox)

	vscroll := ui.NewVScroll("vscroll", 100, 50)
	vscroll.SetPos(pixel.V(150, 575))
	d.AddChild(vscroll)
	vbox = ui.NewVBox("innervbox")
	text = ui.NewText("", "basic")
	text.Printf("Line 1")
	vbox.AddChild(text)
	text = ui.NewText("", "basic")
	text.Printf("Line 2 is pretty long")
	vbox.AddChild(text)
	text = ui.NewText("", "basic")
	text.Printf("Line 3 is not long")
	vbox.AddChild(text)
	text = ui.NewText("", "basic")
	text.Printf("And another line")
	vbox.AddChild(text)

	vscroll.SetInner(vbox)

	dropdown := ui.NewDropDown("dropdown", "basic", 100, 20, 60)
	dropdown.SetPos(pixel.V(150, 490))
	dropdown.OnChange(func(v string) {
		fmt.Printf("dropdown: %s\n", v)
	})
	dropdown.AddValue("Choice 1", "c1")
	dropdown.AddValue("Another one", "c2")
	dropdown.AddValue("One really long choice!!", "c3")
	d.AddChild(dropdown)
}

func (d *demo) Update(dt float64) {
	dphi := math.Pi * dt
	d.text1.SetRot(d.text1.GetRot() + dphi)
	d.text2.SetRot(d.text2.GetRot() + dphi)
	d.text3.SetRot(d.text3.GetRot() + dphi)

	d.sprite.SetRot(d.sprite.GetRot() + dphi)

	d.totalT += math.Pi * dt
	newscale := 1.0 + math.Sin(d.totalT)*.5
	d.sprite2.SetScale(pixel.V(newscale, newscale))
	if d.totalT >= 2*math.Pi {
		d.totalT -= 2 * math.Pi
	}

	dphi = math.Pi / 5.0 * dt
	if nodes.Events().Pressed(pixelgl.KeyA) {
		d.rotslider.SetRot(d.rotslider.GetRot() + dphi)
	} else if nodes.Events().Pressed(pixelgl.KeyD) {
		d.rotslider.SetRot(d.rotslider.GetRot() - dphi)
	}
}

var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

func Run() {
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		err = pprof.StartCPUProfile(f)
		if err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "pixelext example",
		Bounds: pixel.R(0, 0, 1200, 800),
		VSync:  true,
	})
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)
	nodes.Events().SetWin(win)

	root := newDemo()
	bgGrid := newBackgroundGrid(win.Bounds().W(), win.Bounds().H())
	bgGrid.SetZIndex(-2)
	bgGrid.SetPos(pixel.V(600, 400))
	root.AddChild(bgGrid)

	nodes.SceneManager().SetRoot(root)

	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyEscape) {
			break
		}
		nodes.SceneManager().Run(pixel.IM)
	}
	if *memprofile != "" {
		fmemprofile, err := os.Create(*memprofile)
		if err != nil {
			panic(err)
		}
		defer fmemprofile.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(fmemprofile); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}

func main() {
	flag.Parse()

	pixelgl.Run(Run)
}
