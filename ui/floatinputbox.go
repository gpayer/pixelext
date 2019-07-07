package ui

import (
	"strconv"

	"github.com/gpayer/pixelext/nodes"
)

type FloatInputBox struct {
	InputBox
}

func NewFloatInputBox(name, atlasname string, w, h float64) *FloatInputBox {
	f := &FloatInputBox{
		InputBox: *NewInputBox(name, atlasname, w, h),
	}
	f.Self = f
	f.UISelf = f
	return f
}

func (f *FloatInputBox) FloatValue() (float64, error) {
	return strconv.ParseFloat(f.Value(), 64)
}

func (f *FloatInputBox) SetFloatValue(v float64) {
	f.SetValue(strconv.FormatFloat(v, 'g', -1, 64))
}

func (f *FloatInputBox) Update(dt float64) {
	f.InputBox.Update(dt)
	ev := nodes.Events()
	if !ev.IsMouseScrollHandled() && ev.MouseHovering(f) {
		scrolly := ev.MouseScroll().Y
		if scrolly != 0 {
			v, err := f.FloatValue()
			if err == nil {
				f.SetFloatValue(v + scrolly)
				f.onchange(f.Value())
			}
		}
	}
}
