package services

import (
	"github.com/gpayer/go-audio-service/generators"
	"github.com/gpayer/go-audio-service/mix"
	"github.com/gpayer/go-audio-service/notes"
	"github.com/gpayer/go-audio-service/snd"
)

type samplesConfig struct {
	ch           *mix.Channel
	mx           *notes.NoteMultiplexer
	secondChance bool
}

type AudioManagerStruct struct {
	maxsamples    int
	music         snd.Readable
	musicfader    *Fader
	musicchannel  *mix.Channel
	samplesconfig map[string]*samplesConfig
	samplesgain   float32
	mixer         *mix.Mixer
	output        *snd.Output
}

var audioManager *AudioManagerStruct

func AudioManager() *AudioManagerStruct {
	if audioManager == nil {
		audioManager = &AudioManagerStruct{
			mixer:         mix.NewMixer(44100),
			samplesconfig: make(map[string]*samplesConfig),
			samplesgain:   0.5,
			maxsamples:    20,
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

func (a *AudioManagerStruct) newSamplesConfig(samples *snd.Samples) *samplesConfig {
	c := &samplesConfig{
		ch:           a.mixer.GetChannel(),
		mx:           notes.NewNoteMultiplexer(),
		secondChance: false,
	}
	c.ch.SetReadable(c.mx)
	gen := generators.NewSample(samples)
	gen.SetPlayFull(true)
	c.mx.SetReadable(gen)
	return c
}

func (a *AudioManagerStruct) SetMaxSamples(max int) {
	a.maxsamples = max
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
	for _, config := range a.samplesconfig {
		config.ch.SetGain(gain)
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
	_, ok := a.samplesconfig[name]
	if !ok {
		checkAll := false
		for len(a.samplesconfig) >= a.maxsamples {
			removed := false
			for name, config := range a.samplesconfig {
				if checkAll || config.mx.ActiveNotes() == 0 {
					if config.secondChance {
						removed = true
						delete(a.samplesconfig, name)
						break
					} else {
						config.secondChance = true
					}
				}
			}
			if !removed {
				checkAll = true
			}
		}
		a.samplesconfig[name] = a.newSamplesConfig(samples)
	}
	a.samplesconfig[name].ch.SetEnabled(true)
	a.samplesconfig[name].secondChance = false
	a.samplesconfig[name].mx.SendNoteEvent(notes.NewNoteEvent(notes.Pressed, 440, gain))

	for _, config := range a.samplesconfig {
		if config.mx.ActiveNotes() == 0 {
			config.ch.SetEnabled(false)
		}
	}
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
