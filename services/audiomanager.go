package services

import (
	"github.com/gpayer/go-audio-service/generators"
	"github.com/gpayer/go-audio-service/mix"
	"github.com/gpayer/go-audio-service/snd"
)

type AudioManagerStruct struct {
	music        *generators.Sample
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
		audioManager.musicchannel = audioManager.mixer.GetChannel()
	}
	return audioManager
}

func (a *AudioManagerStruct) PlayMusic(samples *snd.Samples) {

}
