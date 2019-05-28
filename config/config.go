package config

import (
	"os"
	"path"
	"path/filepath"
	"runtime"

	goup "github.com/ufoscout/go-up"

	ers "github.com/colt3k/nglog/ers/bserr"
	log "github.com/colt3k/nglog/ng"
)

type Config struct {
	Util goup.GoUp
}

func NewConfig() *Config {
	tmp := &Config{}
	return tmp
}
func (c *Config) Load(file string) {

	if len(file) <= 0 {
		file = ".env"
	}
	ex, err := os.Executable()
	if err != nil {
		log.Logf(log.FATAL,"issue loading\n%+v", err)
	}
	exPath := filepath.Dir(ex)
	log.Logln(log.DEBUG, "Path to executable:", exPath)

	str, err := filepath.Abs(".")
	log.Logln(log.DEBUG, "Executed from:", str)

	_, currentFilePath, _, _ := runtime.Caller(1)
	appdir := path.Dir(currentFilePath)

	envdir := path.Join(exPath, file)
	ignoreFileNotFound := false
	c.Util, err = goup.NewGoUp().AddFile(envdir, ignoreFileNotFound).Build()
	if ers.NoPrintErr(err) {
		envdir = path.Join(appdir, file)
		c.Util, err = goup.NewGoUp().AddFile(envdir, ignoreFileNotFound).Build()
		if ers.NoPrintErr(err) {
			c.Util, err = goup.NewGoUp().AddFile(file, ignoreFileNotFound).Build()
			ers.StopErr(err)
		}
	}
}

func (c *Config) Save() {
}

func (c *Config) Delete() {
}
