package ioimage

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"log"
)

type ImageMeta struct {
	width   int
	height  int
	imgtype string
}

func NewImageMeta() *ImageMeta {
	return new(ImageMeta)
}
func (i *ImageMeta) Dimensions(imagePath string) (int, int, string) {

	if i.height > 0 || i.width > 0 {
		return i.width, i.height, i.imgtype
	}

	file, errOpenFile := os.Open(imagePath)
	defer file.Close()
	if errOpenFile != nil {
		log.Printf("ERROR: \n%+v\n", errOpenFile)
	}

	cfg, imgtype, errDecodeConfig := image.DecodeConfig(file)
	if errDecodeConfig != nil {
		log.Printf("ERROR: %s:\n%+v\n", imagePath, errDecodeConfig)
	}

	return cfg.Width, cfg.Height, imgtype
}
