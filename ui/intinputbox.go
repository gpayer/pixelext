package ui

import "strconv"

type IntInputBox struct {
	InputBox
}

func NewIntInputBox(name, atlasname string, w, h float64) *IntInputBox {
	i := &IntInputBox{
		InputBox: *NewInputBox(name, atlasname, w, h),
	}
	i.Self = i
	i.UISelf = i
	return i
}

func (i *IntInputBox) IntValue() (int64, error) {
	return strconv.ParseInt(i.Value(), 10, 64)
}

func (i *IntInputBox) SetIntValue(v int64) {
	i.SetValue(strconv.FormatInt(v, 10))
}
