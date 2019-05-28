package store

import (
	"path/filepath"
	"time"

	"strconv"

	log "github.com/colt3k/nglog/ng"
	"github.com/colt3k/utils/file"
	"github.com/colt3k/utils/io/ioreader/ioimage"
)

//FileStore datastore for File data
type FileStore struct {
	File     string `json:"file"`
	file     *file.File
	Name     string `json:"name"`
	Dir      bool   `json:"directory"`
	Symlink  bool   `json:"symlink"`
	Ext      string `json:"file_ext"`
	Hash     string `json:"hash"`
	Path     string
	SameName bool              `json:"samename"`
	Status   string            `json:"status"`
	LastMod  int64             `json:"lastmod"`
	Size     int64             `json:"size"`
	SizeHR   string            `json:"sizehr"`
	Time     time.Time         `json:"time"`
	Exif     map[string]string `json:"exif"`
}

//Init initialize the data store
func (fs *FileStore) Init() {

}

func (fs *FileStore) GetExt() *string {
	if len(fs.Ext) > 0 {
		return &fs.Ext
	}
	return &fs.Ext
}

//GetLastMod get the last modification data from file
func (fs *FileStore) GetLastMod() *int64 {
	if fs.LastMod > 0 {
		return &fs.LastMod
	}
	return &fs.LastMod
}

//META printout meta data
func (fs *FileStore) META() {

	log.Println("StoreObject [file=", fs.File, ", hash=", fs.Hash, ", status=", fs.Status, "]")
}

//GetEXIF retrive exif data on file if jpg or jpeg
func (fs *FileStore) GetEXIF() map[string]string {

	//log.Println("Type:",filepath.Ext(fs.File))
	if len(fs.Exif) <= 0 && (filepath.Ext(fs.File) == ".jpg" || filepath.Ext(fs.File) == ".jpeg") {
		//exf := ioexif.NewExif()
		//exf.Exif(fs.File)
		//fs.Exif = exf.Data
		if len(fs.Exif["xdimension"]) <= 0 {
			t := ioimage.NewImageMeta()
			w, h, _ := t.Dimensions(fs.File)
			fs.Exif["xdimension"] = strconv.Itoa(w)
			fs.Exif["ydimension"] = strconv.Itoa(h)
		}
		return fs.Exif
	}
	return nil
}
