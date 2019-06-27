package ui

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/gpayer/pixelext/nodes"
)

func AlignUINode(uinode UINode, bounds pixel.Rect, valignment nodes.VerticalAlignment, halignment nodes.HorizontalAlignment) {
	var x, y float64
	var alignment nodes.Alignment
	doRound := uinode.GetStyles().RoundToPixels
	bounds = bounds.Norm()
	switch halignment {
	case nodes.HAlignmentLeft:
		x = bounds.Min.X
		switch valignment {
		case nodes.VAlignmentTop:
			alignment = nodes.AlignmentTopLeft
		case nodes.VAlignmentCenter:
			alignment = nodes.AlignmentCenterLeft
		case nodes.VAlignmentBottom:
			alignment = nodes.AlignmentBottomLeft
		}
	case nodes.HAlignmentCenter:
		x = bounds.Min.X + bounds.W()/2
		if doRound {
			x = math.Round(x)
		}
		switch valignment {
		case nodes.VAlignmentTop:
			alignment = nodes.AlignmentTopCenter
		case nodes.VAlignmentCenter:
			alignment = nodes.AlignmentCenter
		case nodes.VAlignmentBottom:
			alignment = nodes.AlignmentBottomCenter
		}
	case nodes.HAlignmentRight:
		x = bounds.Max.X
		switch valignment {
		case nodes.VAlignmentTop:
			alignment = nodes.AlignmentTopRight
		case nodes.VAlignmentCenter:
			alignment = nodes.AlignmentCenterRight
		case nodes.VAlignmentBottom:
			alignment = nodes.AlignmentBottomRight
		}
	}

	switch valignment {
	case nodes.VAlignmentTop:
		y = bounds.Max.Y
	case nodes.VAlignmentCenter:
		y = bounds.Min.Y + bounds.H()/2
		if doRound {
			y = math.Round(y)
		}
	case nodes.VAlignmentBottom:
		y = bounds.Min.Y
	}

	uinode.SetPos(pixel.V(x, y))
	uinode.SetAlignment(alignment)
}
