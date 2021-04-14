package config

import (
	log "github.com/colt3k/nglog/ng"
	"os"
	"path"
	"path/filepath"
	"runtime"

	goup "github.com/ufoscout/go-up"
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
	if err != nil {
		envdir = path.Join(appdir, file)
		c.Util, err = goup.NewGoUp().AddFile(envdir, ignoreFileNotFound).Build()
		if err != nil {
			c.Util, err = goup.NewGoUp().AddFile(file, ignoreFileNotFound).Build()
			if err != nil {
				log.Logf(log.FATAL, "%v",err)
			}
		}
	}
}

func (c *Config) Save() {
}

func (c *Config) Delete() {
}
