package services

import (
	"github.com/gpayer/go-audio-service/generators"
	"github.com/gpayer/go-audio-service/mix"
	"github.com/gpayer/go-audio-service/notes"
	"github.com/gpayer/go-audio-service/snd"
)

type AudioManagerStruct struct {
	music               snd.Readable
	musicfader          *Fader
	musicchannel        *mix.Channel
	sampleschannels     map[string]*mix.Channel
	samplesmultiplexers map[string]*notes.NoteMultiplexer
	samplesgain         float32
	mixer               *mix.Mixer
	output              *snd.Output
}

var audioManager *AudioManagerStruct

func AudioManager() *AudioManagerStruct {
	if audioManager == nil {
		audioManager = &AudioManagerStruct{
			mixer:               mix.NewMixer(44100),
			sampleschannels:     make(map[string]*mix.Channel),
			samplesmultiplexers: make(map[string]*notes.NoteMultiplexer),
			samplesgain:         0.5,
		}
		output, err := snd.NewOutput(44100, 2048)
		if err != nil {
			panic(err)
		}
		audioManager.output = output
		audioManager.output.SetReadable(audioManager.mixer)
		audioManager.mixer.SetGain(0.5)
		audioManager.musicchannel = audioManager.mixer.GetChannel()
		audioManager.musicchannel.SetGain(0.5)
		audioManager.musicfader = NewFader()
		audioManager.musicchannel.SetReadable(audioManager.musicfader)
	}
	return audioManager
}

func (a *AudioManagerStruct) SetMusicGain(gain float32) {
	a.musicchannel.SetGain(gain)
}

func (a *AudioManagerStruct) MusicGain() float32 {
	return a.musicchannel.Gain()
}

func (a *AudioManagerStruct) SetMasterGain(gain float32) {
	a.mixer.SetGain(gain)
}

func (a *AudioManagerStruct) MasterGain() float32 {
	return a.mixer.Gain()
}

func (a *AudioManagerStruct) SetSamplesGain(gain float32) {
	a.samplesgain = gain
	for _, ch := range a.sampleschannels {
		ch.SetGain(gain)
	}
}

func (a *AudioManagerStruct) SamplesGain() float32 {
	return a.samplesgain
}

func (a AudioManagerStruct) PlaySample(name string, gain float32) {
	samples, err := ResourceManager().LoadSample(name)
	if err != nil {
		panic(err)
	}
	ch, ok := a.sampleschannels[name]
	if !ok {
		ch = a.mixer.GetChannel()
		ch.SetGain(a.samplesgain)
		a.sampleschannels[name] = ch
		gen := generators.NewSample(samples)
		gen.SetPlayFull(true)
		mx := notes.NewNoteMultiplexer()
		mx.SetReadable(gen)
		ch.SetReadable(mx)
		a.samplesmultiplexers[name] = mx
	}
	a.samplesmultiplexers[name].SendNoteEvent(notes.NewNoteEvent(notes.Pressed, 440, gain))
}

func (a *AudioManagerStruct) PlayMusicFromPath(path string, fadein uint32, loop bool) {
	if ResourceManager().HasSample(path) {
		s, _ := ResourceManager().LoadSample(path)
		a.PlayMusic(s, fadein, loop)
	} else {
		m, err := ResourceManager().CreateMp3Streamer(path)
		if err != nil {
			panic(err)
		}
		audioManager.musicfader.SetOn(false)
		audioManager.musicfader.SetLoop(loop)
		audioManager.music = m
		audioManager.musicfader.SetReadable(audioManager.music)
		if fadein == 0 {
			fadein = 1
		}
		audioManager.musicfader.FadeIn(fadein)
		_ = audioManager.output.Start()
	}
}

func (a *AudioManagerStruct) PlayMusic(samples *snd.Samples, fadein uint32, loop bool) {
	audioManager.musicfader.SetOn(false)
	audioManager.musicfader.SetLoop(loop)
	g := generators.NewSample(samples)
	g.SetPlayFull(true)
	audioManager.music = g
	audioManager.musicfader.SetReadable(audioManager.music)
	if fadein == 0 {
		fadein = 1
	}
	audioManager.musicfader.FadeIn(fadein)
	_ = audioManager.output.Start()
}

func (a *AudioManagerStruct) FadeOut(fadeout uint32) {
	if fadeout == 0 {
		fadeout = 1
	}
	audioManager.musicfader.FadeOut(fadeout)
}
