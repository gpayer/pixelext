package services

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/gpayer/go-audio-service/snd"
	"golang.org/x/image/font"

	"github.com/faiface/pixel"
	"github.com/go-audio/wav"
	"github.com/hajimehoshi/go-mp3"
)

type fontFaces struct {
	font  *truetype.Font
	faces map[float64]font.Face
}

func newFontFaces(fnt *truetype.Font) *fontFaces {
	return &fontFaces{
		font:  fnt,
		faces: make(map[float64]font.Face, 0),
	}
}

type ResourceManagerStruct struct {
	pics     map[string]pixel.Picture
	samples  map[string]*snd.Samples
	fonts    map[string]*fontFaces
	basepath string
}

func (r *ResourceManagerStruct) fixPath(path string) string {
	if len(r.basepath) > 0 {
		var sb strings.Builder
		sb.WriteString(r.basepath)
		sb.WriteRune('/')
		sb.WriteString(path)
		path = sb.String()
	}
	return path
}

func (r *ResourceManagerStruct) LoadPicture(path string) (pixel.Picture, error) {
	path = r.fixPath(path)
	p, ok := r.pics[path]
	if ok {
		return p, nil
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	p = pixel.PictureDataFromImage(img)
	r.pics[path] = p
	return p, nil
}

func (r *ResourceManagerStruct) HasSample(path string) bool {
	path = r.fixPath(path)
	_, ok := r.samples[path]
	return ok
}

// LoadAsync controls whether sound samples are loaded asynchronously in the background.
// After the first 44000 frames a goroutine is spawned to load the rest.
var LoadAsync bool = true

func (r *ResourceManagerStruct) LoadSample(path string) (*snd.Samples, error) {
	path = r.fixPath(path)
	s, ok := r.samples[path]
	if ok {
		return s, nil
	}

	if strings.HasSuffix(path, ".wav") {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		decoder := wav.NewDecoder(f)
		intbuf, err := decoder.FullPCMBuffer()
		if err != nil {
			return nil, err
		}
		buf := intbuf.AsFloat32Buffer()

		samples := snd.NewSamples(uint32(buf.Format.SampleRate), buf.NumFrames())
		for i := 0; i < buf.NumFrames(); i++ {
			samples.Frames[i].L = buf.Data[i*buf.Format.NumChannels]
			if buf.Format.NumChannels > 1 {
				samples.Frames[i].R = buf.Data[i*buf.Format.NumChannels+1]
			}
		}
		r.samples[path] = samples

		return samples, nil
	} else if strings.HasSuffix(path, ".mp3") {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		d, err := mp3.NewDecoder(f)
		if err != nil {
			f.Close()
			return nil, err
		}

		samples := snd.NewSamples(uint32(d.SampleRate()), int(d.Length()/4))

		errCh := make(chan error, 0)

		go func() {
			defer f.Close()
			intbytes := make([]byte, 4)
			sentOk := false
			for i := 0; i < len(samples.Frames); i++ {
				_, err := d.Read(intbytes)
				if err == io.EOF {
					if !sentOk {
						errCh <- io.EOF
					}
					return
				} else if err != nil {
					if !sentOk {
						errCh <- err
					}
					return
				}
				intvalL := int16(intbytes[0]) | int16(intbytes[1])<<8
				intvalR := int16(intbytes[2]) | int16(intbytes[3])<<8

				samples.Frames[i] = snd.Sample{
					L: float32(intvalL) / 32768.0,
					R: float32(intvalR) / 32768.0,
				}
				if i > 144000 && !sentOk {
					sentOk = true
					errCh <- io.EOF
				}
			}
			if !sentOk {
				errCh <- io.EOF
			}
		}()
		err = <-errCh
		if err != nil && err != io.EOF {
			return nil, err
		}
		r.samples[path] = samples

		return samples, nil
	}

	return nil, fmt.Errorf("unsupported sound file format")
}

func (r *ResourceManagerStruct) CreateMp3Streamer(path string) (*Mp3Streamer, error) {
	path = r.fixPath(path)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	var buf []byte = make([]byte, info.Size())
	_, err = io.ReadFull(f, buf)
	if err != nil {
		return nil, err
	}

	return NewMp3Streamer(buf, path), nil
}

func (r *ResourceManagerStruct) LoadTTF(path string, size float64) (font.Face, error) {
	path = r.fixPath(path)

	font, ok := r.fonts[path]
	if ok {
		face, ok := font.faces[size]
		if ok {
			return face, nil
		}
	} else {

		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		font, err := truetype.Parse(bytes)
		if err != nil {
			return nil, err
		}
		r.fonts[path] = newFontFaces(font)
	}

	face := truetype.NewFace(r.fonts[path].font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	})

	r.fonts[path].faces[size] = face

	return face, nil
}

func (r *ResourceManagerStruct) SetBasePath(basepath string) {
	r.basepath = basepath
}

var resourceManager *ResourceManagerStruct

func ResourceManager() *ResourceManagerStruct {
	return resourceManager
}

func init() {
	resourceManager = &ResourceManagerStruct{
		pics:     make(map[string]pixel.Picture, 0),
		samples:  make(map[string]*snd.Samples, 0),
		fonts:    make(map[string]*fontFaces),
		basepath: "",
	}
}
