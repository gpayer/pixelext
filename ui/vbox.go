package ui

import (
	"fmt"
	"math"

	"github.com/gpayer/pixelext/nodes"

	"github.com/faiface/pixel"
)

type VBox struct {
	UIBase
	background *nodes.BorderBox
}

func NewVBox(name string) *VBox {
	v := &VBox{
		UIBase:     *NewUIBase(name),
		background: nodes.NewBorderBox("__bg", 1, 1),
	}
	v.Self = v
	v.UISelf = v
	return v
}

func (v *VBox) Init() {
	v.background.SetPos(pixel.ZV)
	v.background.SetStyles(v.GetStyles())
	v.background.SetZIndex(-1)
	v.AddChild(v.background)
}

func (v *VBox) recalcPositions() {
	padding := v.GetStyles().Padding
	borderWidth := v.GetStyles().Border.Width
	ypos := padding
	maxx := 0.0
	fmt.Printf("num children: %d\n", len(v.Children()))
	for _, child := range v.Children() {
		uichild, ok := child.(UINode)
		if ok && child.GetName() != "__bg" && !child.IsRemove() {
			childbounds := uichild.Size()
			ypos += childbounds.Y + 2*padding
			if childbounds.X > maxx {
				maxx = childbounds.X
			}
		}
	}
	size := pixel.V(math.Round(maxx+2*padding+borderWidth), math.Round(ypos-padding+borderWidth))
	fmt.Printf("vbox size: %v\n", size)
	v.SetSize(size)
	v.background.SetSize(size)

	ypos = math.Round(size.Y/2 - padding - borderWidth/2)
	maxxHalf := math.Round(maxx / 2)
	for _, child := range v.Children() {
		uichild, ok := child.(UINode)
		if ok && child.GetName() != "__bg" && !child.IsRemove() {
			switch v.halignment {
			case nodes.HAlignmentLeft:
				uichild.SetAlignment(nodes.AlignmentTopLeft)
				uichild.SetPos(pixel.V(-maxxHalf, ypos))
			case nodes.HAlignmentCenter:
				uichild.SetAlignment(nodes.AlignmentTopCenter)
				uichild.SetPos(pixel.V(0, ypos))
			case nodes.HAlignmentRight:
				uichild.SetAlignment(nodes.AlignmentTopRight)
				uichild.SetPos(pixel.V(maxxHalf, ypos))
			}
			childbounds := uichild.Size()
			ypos -= childbounds.Y + 2*padding
		}
	}
	nodes.SceneManager().Redraw()
}

func (v *VBox) AddChild(child nodes.Node) {
	v.UIBase.AddChild(child)
	v.recalcPositions()
}

func (v *VBox) SetStyles(styles *nodes.Styles) {
	v.UIBase.SetStyles(styles)
	v.background.SetStyles(styles)
	v.recalcPositions()
}

func (v *VBox) RemoveChild(child nodes.Node) {
	v.UIBase.RemoveChild(child)
	v.recalcPositions()
}

func (v *VBox) RemoveChildren() {
	for _, ch := range v.Children() {
		if ch.GetName() != "__bg" {
			ch.SetRemove(true)
		}
	}
	v.recalcPositions()
}
