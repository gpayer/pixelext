package services

import (
	"testing"

	"github.com/gpayer/go-audio-service/mix"
	"github.com/gpayer/go-audio-service/snd"
	"github.com/stretchr/testify/assert"
)

func newMockAudioManager() *AudioManagerStruct {
	audioManager = &AudioManagerStruct{
		mixer:         mix.NewMixer(44100),
		samplesconfig: make(map[string]*samplesConfig),
		samplesgain:   0.5,
		maxsamples:    20,
	}
	audioManager.mixer.SetGain(0.5)
	audioManager.musicchannel = audioManager.mixer.GetChannel()
	audioManager.musicchannel.SetGain(0.5)
	audioManager.musicfader = NewFader()
	audioManager.musicchannel.SetReadable(audioManager.musicfader)
	return audioManager
}

func TestMaxSamplesCleanup(t *testing.T) {
	assert := assert.New(t)
	samples := snd.NewSamples(100, 10)
	ResourceManager().samples["test1"] = samples
	ResourceManager().samples["test2"] = samples
	ResourceManager().samples["test3"] = samples
	a := newMockAudioManager()
	a.SetMaxSamples(2)

	a.PlaySample("test1", 1.0)
	out := snd.NewSamples(100, 100)
	a.mixer.Read(out)
	a.PlaySample("test2", 1.0)
	assert.Len(a.samplesconfig, 2)
	assert.Equal(false, a.samplesconfig["test1"].ch.Enabled())
	assert.Equal(true, a.samplesconfig["test2"].ch.Enabled())

	a.PlaySample("test3", 1.0)
	assert.Len(a.samplesconfig, 2)
	for name, _ := range a.samplesconfig {
		if name == "test1" {
			assert.Fail("Expected \"test1\" to NOT be in a.sampleconfig")
		}
	}

	a.PlaySample("test1", 1.0)
	assert.Len(a.samplesconfig, 2)
}
