package ui

import (
	"math"
	"pixelext/nodes"

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
	for _, child := range v.Children() {
		uichild, ok := child.(UINode)
		if ok && child.GetName() != "__bg" {
			childbounds := uichild.Size()
			ypos += childbounds.Y + 2*padding
			if childbounds.X > maxx {
				maxx = childbounds.X
			}
		}
	}
	size := pixel.V(math.Round(maxx+2*padding+borderWidth), math.Round(ypos-padding+borderWidth))
	v.SetSize(size)
	v.background.SetSize(size)

	ypos = size.Y/2 - padding - borderWidth/2
	for _, child := range v.Children() {
		uichild, ok := child.(UINode)
		if ok && child.GetName() != "__bg" {
			uichild.SetAlignment(nodes.AlignmentTopCenter)
			uichild.SetPos(pixel.V(0, ypos))
			childbounds := uichild.Size()
			ypos -= childbounds.Y + 2*padding
		}
	}
	nodes.SceneManager().Redraw()
}

func (v *VBox) AddChild(child nodes.Node) {
	v.BaseNode.AddChild(child)
	v.recalcPositions()
}

func (v *VBox) SetStyles(styles *nodes.Styles) {
	v.BaseNode.SetStyles(styles)
	v.background.SetStyles(styles)
	v.recalcPositions()
}
