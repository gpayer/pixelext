package ui

import (
	"strconv"

	"github.com/gpayer/pixelext/nodes"
)

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

func (i *IntInputBox) Update(dt float64) {
	i.InputBox.Update(dt)
	ev := nodes.Events()
	if !ev.IsMouseScrollHandled() && ev.MouseHovering(i) {
		scrolly := int64(ev.MouseScroll().Y)
		if scrolly != 0 {
			ival, err := i.IntValue()
			if err == nil {
				i.SetIntValue(ival + scrolly)
				i.onchange(i.Value())
			}
		}
	}
}
