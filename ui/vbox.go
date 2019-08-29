package ui

import (
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
	v.background.SetLocked(true)
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
	notremove := 0
	notok := 0
	for _, child := range v.Children() {
		if child == nil {
			continue
		}
		uichild, ok := child.(UINode)
		if ok && child.GetName() != "__bg" {
			notremove++
			childbounds := uichild.Size()
			ypos += childbounds.Y + 2*padding
			if childbounds.X > maxx {
				maxx = childbounds.X
			}
		} else if !ok {
			notok++
		}
	}
	size := pixel.V(math.Round(maxx+2*padding+2*borderWidth), math.Round(ypos-padding+2*borderWidth))
	v.UISelf.SetSize(size)
	v.background.SetSize(size)

	ypos = math.Round(size.Y/2 - padding - borderWidth)
	maxxHalf := math.Round(maxx / 2)
	for _, child := range v.Children() {
		if child == nil {
			continue
		}
		uichild, ok := child.(UINode)
		if ok && child.GetName() != "__bg" {
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
	child.SetPosLocked(true)
}

func (v *VBox) SetStyles(styles *nodes.Styles) {
	v.UIBase.SetStyles(styles)
	v.background.SetStyles(styles)
	v.recalcPositions()
}

func (v *VBox) RemoveChild(child nodes.Node) {
	v.UIBase.RemoveChild(child)
	v.recalcPositions()
	child.SetPosLocked(false)
}

func (v *VBox) RemoveChildren() {
	for _, c := range v.Children() {
		if c != nil {
			c.SetPosLocked(false)
		}
	}
	v.UIBase.RemoveChildren()
	v.AddChild(v.background)
	v.recalcPositions()
}

func (v *VBox) ChildChanged() {
	v.recalcPositions()
}

func (v *VBox) SetHAlignment(val nodes.HorizontalAlignment) {
	v.UIBase.SetHAlignment(val)
	v.recalcPositions()
}
