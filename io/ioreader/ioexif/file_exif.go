package ioexif

import (
	"github.com/colt3k/utils/debug"
	"os"
	"strconv"
	"strings"

	"github.com/rwcarlsen/goexif/exif"

	"log"
)

// FileExif data store for exif data
type FileExif struct {
	Data map[string]string
}

// New create a new FilExif store
func New() *FileExif {
	tmp := &FileExif{}
	return tmp
}

// ReadLatLongData find lat/long data and return
func (x *FileExif) ReadLatLongData(f *os.File) {
	xif, err := exif.Decode(f)
	if err != nil {
		log.Printf("ERROR: exif read %v\n", err)
		debug.PrintStack()
	}

	lat, long, err := xif.LatLong()
	if err != nil {
		log.Printf("ERROR: exif read latlong data %v\n", err)
		debug.PrintStack()
	}

	x.Data["lat"] = strconv.FormatFloat(lat, 'f', -1, 64)
	x.Data["long"] = strconv.FormatFloat(long, 'f', -1, 64)

}

// ReadALLDataAsJSON read data from exif metadata and set to our datastore
func (x *FileExif) ReadALLDataAsJSON(f *os.File) {
	xif, err := exif.Decode(f)
	if err != nil {
		log.Printf("ERROR: exif read data as json %v\n", err)
		debug.PrintStack()
	}
	json, _ := xif.MarshalJSON()
	x.Data["jsondata"] = string(json)
}

/*
ReadDataByKey find data by key
exif.FocalLength
exif.PixelXDimension
exif.PixelYDimension
exif.ExifVersion
exif.ThumbJPEGInterchangeFormat
exif.ThumbJPEGInterchangeFormatLength
exif.Model
exif.Make
exif.Flash
exif.ExposureTime
exif.DigitalZoomRatio
*/
func (x *FileExif) ReadDataByKey(f *os.File, key []exif.FieldName) error {
	xif, errDecode := exif.Decode(f)
	if errDecode != nil {
		log.Printf("WARN: file_exif ReadDataByKey, no exif data %s\n%+v\n", f.Name(), errDecode)
		return nil
	}

	for _, d := range key {
		val, errGetKey := xif.Get(d)
		//if ers.NotErr(err, f.Name()) {
		if errGetKey == nil {
			switch d {
			case exif.ThumbJPEGInterchangeFormat:
				data, errTiffVal := val.Int64(0)
				if errTiffVal != nil {
					log.Println(errTiffVal)
				}
				x.Data["jpgformat"] = strconv.FormatInt(data, 10)
			case exif.ExifVersion:
				str := string(val.Val)
				x.Data["exifversion"] = str
			case exif.PixelYDimension:
				if errGetKey != nil {
					log.Println(errGetKey)
				} else {
					data, errTiffVal := val.Int64(0)
					if errTiffVal != nil {
						log.Println(errTiffVal)
					}
					x.Data["ydimension_height"] = strconv.FormatInt(data, 10)
				}
			case exif.PixelXDimension:
				if errGetKey != nil {
					log.Println(errGetKey)
				} else {
					data, errTiffVal := val.Int64(0)
					if errTiffVal != nil {
						log.Println(errTiffVal)
					}
					x.Data["xdimension_width"] = strconv.FormatInt(data, 10)
				}
			case exif.FocalLength:
				numer, denom, _ := val.Rat2(0)
				x.Data["focalnumerator"] = strconv.FormatInt(numer, 10)
				x.Data["focaldenominator"] = strconv.FormatInt(denom, 10)
			case exif.ThumbJPEGInterchangeFormatLength:
				data, errTiffVal := val.Int64(0)
				if errTiffVal != nil {
					log.Println(errTiffVal)
				}
				x.Data["jpgintercahngeformatlength"] = strconv.FormatInt(data, 10)
			case exif.Make:
				str, errTiffVal := val.StringVal()
				if errTiffVal != nil {
					log.Println(errTiffVal)
				}
				x.Data["make"] = str
			case exif.Model:
				str, errTiffVal := val.StringVal()
				if errTiffVal != nil {
					log.Println(errTiffVal)
				}
				x.Data["model"] = str
			case exif.Flash:
				data, errTiffVal := val.Int64(0)
				if errTiffVal != nil {
					log.Println(errTiffVal)
				}
				x.Data["flash"] = strconv.FormatInt(data, 10)
			case exif.ExposureTime:
				str, errTiffVal := val.Rat(0)
				if errTiffVal != nil {
					log.Println(errTiffVal)
				}
				x.Data["exposuretimenumerator"] = strconv.FormatInt(str.Num().Int64(), 10)
				x.Data["exposuretimedenominator"] = strconv.FormatInt(str.Denom().Int64(), 10)
			case exif.DigitalZoomRatio:
				one, two, errTiffVal := val.Rat2(0)
				if errTiffVal != nil {
					log.Println(errTiffVal)
				}
				x.Data["digitalzoomratio"] = strconv.FormatInt(one, 10) + "/" + strconv.FormatInt(two, 10)
			case exif.DateTimeOriginal:
				str, errTiffVal := val.StringVal()
				if errTiffVal != nil {
					log.Println(errTiffVal)
				}
				x.Data["datecreated"] = str
			case exif.DateTime:
				str, errTiffVal := val.StringVal()
				if errTiffVal != nil {
					log.Println(errTiffVal)
				}
				x.Data["datetime"] = str
			case exif.DateTimeDigitized:
				str, errTiffVal := val.StringVal()
				if errTiffVal != nil {
					log.Println(errTiffVal)
				}
				x.Data["datedigitized"] = str
			}
		}
	}

	return nil
}

// Exif determine all exif meta data attached to file i.e. jpg type
func (x *FileExif) Exif(fileName string) error {

	file, errOpenFile := os.Open(fileName)
	if errOpenFile != nil {
		log.Fatalf("ERROR: issue opening\n%+v", errOpenFile)
	}
	//Tell the program to call the following function when the current function returns
	defer file.Close()

	x.Data = make(map[string]string)
	if strings.Index(strings.ToLower(file.Name()), "jpg") > -1 ||
		strings.Index(strings.ToLower(file.Name()), "jpeg") > -1 {
		keys := make([]exif.FieldName, 0)
		keys = append(keys, exif.PixelXDimension)
		keys = append(keys, exif.PixelYDimension)
		keys = append(keys, exif.Model)
		keys = append(keys, exif.Make)
		keys = append(keys, exif.DateTimeOriginal)
		keys = append(keys, exif.DateTime)
		keys = append(keys, exif.DateTimeDigitized)
		errReadDataByKey := x.ReadDataByKey(file, keys)
		if errReadDataByKey != nil {
			return errReadDataByKey
		}
	}
	return nil
}
