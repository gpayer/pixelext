package services

import (
	"github.com/gpayer/go-audio-service/generators"
	"github.com/gpayer/go-audio-service/mix"
	"github.com/gpayer/go-audio-service/snd"
)

type AudioManagerStruct struct {
	music        *generators.Sample
	musicfader   *Fader
	musicchannel *mix.Channel
	mixer        *mix.Mixer
	output       *snd.Output
}

var audioManager *AudioManagerStruct

func AudioManager() *AudioManagerStruct {
	if audioManager == nil {
		audioManager = &AudioManagerStruct{
			mixer: mix.NewMixer(44100),
		}
		output, err := snd.NewOutput(44100, 2048)
		if err != nil {
			panic(err)
		}
		audioManager.output = output
		audioManager.output.SetReadable(audioManager.mixer)
		audioManager.mixer.SetGain(0.8)
		audioManager.musicchannel = audioManager.mixer.GetChannel()
		audioManager.musicchannel.SetGain(0.5)
		audioManager.musicfader = NewFader()
		audioManager.musicchannel.SetReadable(audioManager.musicfader)
	}
	return audioManager
}

func (a *AudioManagerStruct) PlayMusic(samples *snd.Samples, fadein uint32, loop bool) {
	audioManager.musicfader.SetOn(false)
	audioManager.musicfader.SetLoop(loop)
	audioManager.music = generators.NewSample(samples)
	audioManager.music.SetPlayFull(true)
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
