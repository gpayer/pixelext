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