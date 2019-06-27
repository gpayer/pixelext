package nodes

type Alignment int

const (
	AlignmentFixed Alignment = iota
	AlignmentBottomLeft
	AlignmentCenterLeft
	AlignmentTopLeft
	AlignmentBottomCenter
	AlignmentCenter
	AlignmentTopCenter
	AlignmentBottomRight
	AlignmentCenterRight
	AlignmentTopRight
)

type HorizontalAlignment int

const (
	HAlignmentLeft HorizontalAlignment = iota
	HAlignmentCenter
	HAlignmentRight
)

type VerticalAlignment int

const (
	VAlignmentTop VerticalAlignment = iota
	VAlignmentCenter
	VAlignmentBottom
)
