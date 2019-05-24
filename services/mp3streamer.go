package services

import (
	"bytes"
	"io"

	"github.com/gpayer/go-audio-service/snd"
	"github.com/hajimehoshi/go-mp3"
)

type Mp3Streamer struct {
	path      string
	bufreader *bytes.Reader
	readUntil uint32
	decoder   *mp3.Decoder
	samples   *snd.Samples
	noteEnded bool
	allRead   bool
}

func NewMp3Streamer(buf []byte, path string) *Mp3Streamer {
	m := &Mp3Streamer{
		readUntil: 0,
		allRead:   false,
		path:      path,
	}
	m.bufreader = bytes.NewReader(buf)
	_, _ = m.bufreader.Seek(0, io.SeekStart)
	var err error
	m.decoder, err = mp3.NewDecoder(m.bufreader)
	if err != nil {
		panic(err)
	}

	m.samples = snd.NewSamples(uint32(m.decoder.SampleRate()), int(m.decoder.Length()/4))

	return m
}

func (m *Mp3Streamer) Read(samples *snd.Samples) {
	m.ReadStateless(samples, 440, snd.EmptyNoteState)
}

func (m *Mp3Streamer) ReadStateless(samples *snd.Samples, freq float32, state *snd.NoteState) {
	var copyFrom, dstOffset, readEnd uint32
	length := uint32(len(samples.Frames))
	dstOffset = 0
	if state.Timecode+length > uint32(len(m.samples.Frames)) {
		readEnd = uint32(len(m.samples.Frames))
		m.noteEnded = true
	} else {
		readEnd = state.Timecode + length
		m.noteEnded = false
	}
	if m.readUntil > 0 {
		if state.Timecode < m.readUntil {
			copyFrom = state.Timecode
			if state.Timecode+length >= m.readUntil {
				dstOffset = m.readUntil - state.Timecode
			} else {
				dstOffset = length
			}
			for i := uint32(0); i < dstOffset; i++ {
				samples.Frames[i] = snd.Sample{
					L: m.samples.Frames[copyFrom].L * state.Volume,
					R: m.samples.Frames[copyFrom].R * state.Volume,
				}
				copyFrom++
			}
		}
	}

	if m.allRead {
		return
	}

	intbytes := make([]byte, 4)
	for i := m.readUntil; i < readEnd; i++ {
		var err error
		if !m.allRead {
			_, err = m.decoder.Read(intbytes)
		}
		if err == io.EOF || m.allRead {
			m.samples.Frames[i] = snd.Sample{L: 0, R: 0}
			m.allRead = true
			m.decoder = nil
			m.bufreader = nil
			ResourceManager().samples[m.path] = m.samples
		} else if err != nil && err != io.EOF {
			panic(err)
		} else {
			intvalL := int16(intbytes[0]) | int16(intbytes[1])<<8
			intvalR := int16(intbytes[2]) | int16(intbytes[3])<<8

			m.samples.Frames[i] = snd.Sample{
				L: float32(intvalL) / 32768.0,
				R: float32(intvalR) / 32768.0,
			}
		}
		if i >= state.Timecode {
			samples.Frames[dstOffset] = snd.Sample{
				L: m.samples.Frames[i].L * state.Volume,
				R: m.samples.Frames[i].R * state.Volume,
			}
			dstOffset++
		}
		m.readUntil++
	}
}

func (m *Mp3Streamer) NoteEnded() bool {
	return m.noteEnded
}
