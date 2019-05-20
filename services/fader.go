package services

import "github.com/gpayer/go-audio-service/snd"

type Fader struct {
	readable      snd.Readable
	on            bool
	fadeRemaining uint32
	dV            float32
	notestate     snd.NoteState
}

func (f *Fader) SetReadable(r snd.Readable) {
	f.readable = r
}

func (f *Fader) Read(samples *snd.Samples) {
	length := len(samples.Frames)
	if !f.on {
		for i := 0; i < length; i++ {
			samples.Frames[i].L = 0
			samples.Frames[i].R = 0
		}
	}
}

func (f *Fader) ReadStateless(samples *snd.Samples, freq float32, state *snd.NoteState) {
	f.Read(samples)
}
