package ioexif

import (
	"os"
	"strconv"
	"strings"

	"github.com/rwcarlsen/goexif/exif"

	log "github.com/colt3k/nglog/ng"
	ers "github.com/colt3k/nglog/ers/bserr"
)

//FileExif data store for exif data
type FileExif struct {
	Data map[string]string
}

//NewExif create a new FilExif store
func New() *FileExif {
	tmp := &FileExif{}
	return tmp
}

//ReadLatLongData find lat/long data and return
func (x *FileExif) ReadLatLongData(f *os.File) {
	xif, err := exif.Decode(f)
	ers.NotErr(err, "exif: read error")
	lat, long, err := xif.LatLong()
	ers.NotErr(err, "exif: read latlong data")

	x.Data["lat"] = strconv.FormatFloat(lat, 'f', -1, 64)
	x.Data["long"] = strconv.FormatFloat(long, 'f', -1, 64)

}

//ReadALLDataAsJSON read data from exif metadata and set to our datastore
func (x *FileExif) ReadALLDataAsJSON(f *os.File) {
	xif, err := exif.Decode(f)
	ers.NotErr(err, "exif: read data as json")
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
	xif, err := exif.Decode(f)

	if err != nil {
		log.Logf(log.WARN, "file_exif ReadDataByKey, no exif data %s\n%+v", f.Name(),err)
		return nil
	}

	for _, d := range key {
		val, err := xif.Get(d)
		//if ers.NotErr(err, f.Name()) {
		if err == nil {
			switch d {
			case exif.ThumbJPEGInterchangeFormat:
				data, err := val.Int64(0)
				ers.NotErrSkipTrace(err)
				x.Data["jpgformat"] = strconv.FormatInt(data, 10)
			case exif.ExifVersion:
				str := string(val.Val)
				x.Data["exifversion"] = str
			case exif.PixelYDimension:
				if ers.NotErrSkipTrace(err) {
					data, err := val.Int64(0)
					ers.NotErrSkipTrace(err)
					x.Data["ydimension_height"] = strconv.FormatInt(data, 10)
				}
			case exif.PixelXDimension:
				if ers.NotErrSkipTrace(err) {
					data, err := val.Int64(0)
					ers.NotErrSkipTrace(err)
					x.Data["xdimension_width"] = strconv.FormatInt(data, 10)
				}
			case exif.FocalLength:
				numer, denom, _ := val.Rat2(0)
				x.Data["focalnumerator"] = strconv.FormatInt(numer, 10)
				x.Data["focaldenominator"] = strconv.FormatInt(denom, 10)
			case exif.ThumbJPEGInterchangeFormatLength:
				data, err := val.Int64(0)
				ers.NotErrSkipTrace(err)
				x.Data["jpgintercahngeformatlength"] = strconv.FormatInt(data, 10)
			case exif.Make:
				str, err := val.StringVal()
				ers.NotErrSkipTrace(err)
				x.Data["make"] = str
			case exif.Model:
				str, err := val.StringVal()
				ers.NotErrSkipTrace(err)
				x.Data["model"] = str
			case exif.Flash:
				data, err := val.Int64(0)
				ers.NotErrSkipTrace(err)
				x.Data["flash"] = strconv.FormatInt(data, 10)
			case exif.ExposureTime:
				str, err := val.Rat(0)
				ers.NotErrSkipTrace(err)
				x.Data["exposuretimenumerator"] = strconv.FormatInt(str.Num().Int64(), 10)
				x.Data["exposuretimedenominator"] = strconv.FormatInt(str.Denom().Int64(), 10)
			case exif.DigitalZoomRatio:
				one, two, err := val.Rat2(0)
				ers.NotErrSkipTrace(err)
				x.Data["digitalzoomratio"] = strconv.FormatInt(one, 10) + "/" + strconv.FormatInt(two, 10)
			case exif.DateTimeOriginal:
				str, err := val.StringVal()
				ers.NotErrSkipTrace(err)
				x.Data["datecreated"] = str
			case exif.DateTime:
				str, err := val.StringVal()
				ers.NotErrSkipTrace(err)
				x.Data["datetime"] = str
			case exif.DateTimeDigitized:
				str, err := val.StringVal()
				ers.NotErrSkipTrace(err)
				x.Data["datedigitized"] = str
			}
		}
	}

	return nil
}

//Exif determine all exif meta data attached to file i.e. jpg type
func (x *FileExif) Exif(fileName string) error {

	file, err := os.Open(fileName)
	if err != nil {
		log.Logf(log.FATAL,"issue opening\n%+v", err)
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
		err := x.ReadDataByKey(file, keys)
		if err != nil {
			return err
		}
	}
	return nil
}
