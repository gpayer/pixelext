package services

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"

	"github.com/faiface/pixel"
)

type ResourceManagerStruct struct {
	pics     map[string]pixel.Picture
	basepath string
}

func (r *ResourceManagerStruct) LoadPicture(path string) (pixel.Picture, error) {
	if len(r.basepath) > 0 {
		var sb strings.Builder
		sb.WriteString(r.basepath)
		sb.WriteRune('/')
		sb.WriteString(path)
		path = sb.String()
	}
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

var resourceManager *ResourceManagerStruct

func ResourceManager() *ResourceManagerStruct {
	return resourceManager
}

func init() {
	resourceManager = &ResourceManagerStruct{
		pics:     make(map[string]pixel.Picture, 0),
		basepath: "",
	}
}
