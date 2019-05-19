package nodes

type animData struct {
	startV, endV, startT, duration float64
	dV                             float64
	animFn                         func(v float64)
}

type Animation struct {
	BaseNode
	anims   []*animData
	totalT  float64
	endT    float64
	onEndFn func()
}

func NewAnimation(name string) *Animation {
	a := &Animation{
		BaseNode: *NewBaseNode(name),
		totalT:   0.0,
		onEndFn:  func() {},
	}
	a.Self = a
	a.SetActive(false)
	return a
}

func (a *Animation) Add(startValue, endValue, startTime, duration float64, animFn func(v float64)) {
	anim := &animData{startValue, endValue, startTime, duration, 0, animFn}
	anim.dV = (endValue - startValue) / duration
	a.anims = append(a.anims, anim)
	if startTime+duration > a.endT {
		a.endT = startTime + duration
	}
}

func (a *Animation) OnEnd(fn func()) {
	a.onEndFn = fn
}

func (a *Animation) Start() {
	a.totalT = 0
	a.SetActive(true)
}

func (a *Animation) Update(dt float64) {
	a.totalT += dt
	if a.totalT > a.endT {
		a.SetActive(false)
		a.onEndFn()
	}
	for _, anim := range a.anims {
		if a.totalT >= anim.startT && a.totalT < anim.startT+anim.duration {
			anim.animFn(anim.startV + (a.totalT-anim.startT)*anim.dV)
		} else if a.totalT > anim.startT+anim.duration {
			anim.animFn(anim.endV)
		}
	}
}
