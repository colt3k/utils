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

	file, err := os.Open(imagePath)
	defer file.Close()
	if err != nil {
		log.Printf("ERROR: \n%+v\n", err)
	}

	cfg, imgtype, err := image.DecodeConfig(file)
	if err != nil {
		log.Printf("ERROR: %s:\n%+v\n", imagePath, err)
	}

	return cfg.Width, cfg.Height, imgtype
}
