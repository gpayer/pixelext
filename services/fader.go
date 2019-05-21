package services

import (
	"sync"

	"github.com/gpayer/go-audio-service/notes"
	"github.com/gpayer/go-audio-service/snd"
)

type Fader struct {
	mtx           sync.Mutex
	readable      snd.Readable
	on            bool
	fadeRemaining uint32
	dV            float32
	notestate     *snd.NoteState
	loop          bool
}

func NewFader() *Fader {
	return &Fader{
		on: false,
		notestate: &snd.NoteState{
			Volume: 1.0,
		},
		loop: true,
	}
}

func (f *Fader) SetReadable(r snd.Readable) {
	f.readable = r
}

func (f *Fader) Read(samples *snd.Samples) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	length := uint32(len(samples.Frames))
	if !f.on || f.readable == nil {
		for i := uint32(0); i < length; i++ {
			samples.Frames[i].L = 0
			samples.Frames[i].R = 0
		}
	}

	noteaware, isNoteAware := f.readable.(notes.NoteAware)

	f.readable.ReadStateless(samples, 440, f.notestate)
	f.notestate.Timecode += length
	if f.fadeRemaining > 0 {
		if f.fadeRemaining < length {
			f.fadeRemaining = 0
			if f.dV > 0 {
				f.notestate.Volume = 1
				f.notestate.ReleaseTimecode = f.notestate.Timecode
				f.notestate.On = false
			} else {
				f.notestate.Volume = 0
				f.on = false
			}
		} else {
			f.fadeRemaining -= length
			f.notestate.Volume += f.dV * float32(length)
		}
	}

	if isNoteAware {
		if noteaware.NoteEnded() {
			if f.loop {
				f.notestate.Timecode = 0
				f.notestate.ReleaseTimecode = 0
				f.notestate.On = false
			} else {
				f.on = false
			}
		}
	}
}

func (f *Fader) ReadStateless(samples *snd.Samples, freq float32, state *snd.NoteState) {
	f.Read(samples)
}

func (f *Fader) FadeIn(samples uint32) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	f.dV = 1.0 / float32(samples)
	f.fadeRemaining = samples
	f.on = true
	f.notestate.Volume = 0
	f.notestate.ReleaseTimecode = 0
	f.notestate.On = true
	f.notestate.Timecode = 0
}

func (f *Fader) FadeOut(samples uint32) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	f.dV = -1.0 / float32(samples)
	f.fadeRemaining = samples
	f.on = true
	f.notestate.Volume = 1
	if f.notestate.ReleaseTimecode == 0 {
		f.notestate.ReleaseTimecode = f.notestate.Timecode
	}
	f.notestate.On = false
}

func (f *Fader) SetOn(on bool) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	f.on = on
}

func (f *Fader) SetLoop(loop bool) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	f.loop = loop
}
