package ioexif

import (
	"os"

	"github.com/rwcarlsen/goexif/exif"

	log "github.com/colt3k/nglog/ng"
	ers "github.com/colt3k/nglog/ers/bserr"
	"github.com/colt3k/utils/io/ioreader/ioimage"
)

func ExampleFileExif_ReadALLDataAsJSON() {
	var file = "/Users/gcollins/Desktop/tester/pics2/12-23-05_0543.jpg"

	f, err := os.Open(file)
	ers.NotErr(err)
	defer f.Close()

	d := &FileExif{}

	d.Data = make(map[string]string)
	//d.ReadALLDataAsJSON(f)		Use to find fields

	keys := make([]exif.FieldName, 0)
	keys = append(keys, exif.FocalLength)
	keys = append(keys, exif.PixelXDimension)
	keys = append(keys, exif.PixelYDimension)
	keys = append(keys, exif.ExifVersion)
	keys = append(keys, exif.ThumbJPEGInterchangeFormat)
	keys = append(keys, exif.ThumbJPEGInterchangeFormatLength)
	keys = append(keys, exif.Model)
	keys = append(keys, exif.Make)
	keys = append(keys, exif.Flash)
	keys = append(keys, exif.ExposureTime)
	keys = append(keys, exif.DigitalZoomRatio)
	keys = append(keys, exif.DateTimeOriginal)
	keys = append(keys, exif.DateTime)
	keys = append(keys, exif.DateTimeDigitized)
	err = d.ReadDataByKey(f, keys)
	if err != nil {
		log.Logf(log.ERROR,"issue reading by key\n%+v", err)
	}
	//w x h
	t := ioimage.NewImageMeta()
	w, h, imgtype := t.Dimensions(file)
	log.Logf(log.INFO, "W: %d, H: %d Type: %s", w, h, imgtype)
	log.Logln(log.INFO, "print data found")
	for k, v := range d.Data {
		log.Println(k, ":", v)
	}

	//log.Println(d.JSONData)

	//log.Println("")

	/*
		Output:
		hi
	*/
}
