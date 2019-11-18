package updater

import (
	"bytes"
	"strconv"

	"github.com/colt3k/utils/mathut"
)

type AppConfig struct {
	OS                string `json:"os"`
	Arch              string `json:"arch"`
	Name              string `json:"name"`
	Timestamp         int64  `json:"timestamp"`
	Hash              string `json:"hash"`
	Version           string `json:"version"`
	Changelog         string `json:"changelog"`
	BaseURL           string
	URL               string
	ArchiveName       string
	User              []byte
	Pass              []byte
	DisableVerifyCert bool
}

func (a *AppConfig) String() string {
	var byt bytes.Buffer
	byt.WriteString("{  OS: ")
	byt.WriteString(a.OS)
	byt.WriteString(",  Arch: ")
	byt.WriteString(a.Arch)
	byt.WriteString(",  Name: ")
	byt.WriteString(a.Name)
	byt.WriteString(",  TS: ")
	byt.WriteString(mathut.FmtInt(int(a.Timestamp)))
	byt.WriteString(",  Hash: ")
	byt.WriteString(a.Hash)
	byt.WriteString(",  Version: ")
	byt.WriteString(a.Version)
	byt.WriteString(",  Changelog: ")
	byt.WriteString(a.Changelog)
	byt.WriteString(",  BaseURL: ")
	byt.WriteString(a.BaseURL)
	byt.WriteString(",  URL: ")
	byt.WriteString(a.URL)
	byt.WriteString(",  ArchiveName: ")
	byt.WriteString(a.ArchiveName)
	byt.WriteString(",  DisableVerifyCert: ")
	byt.WriteString(strconv.FormatBool(a.DisableVerifyCert))
	byt.WriteString("  }")
	return byt.String()
}

func NewUser(user, passOrToken, urlPrefix, repository string) *Connection {
	t := new(Connection)
	t.User = user
	t.PassOrToken = passOrToken
	t.URLPrefix = urlPrefix
	t.Repository = repository
	return t
}

type Connection struct {
	Name                string
	User                string
	PassOrToken         string
	URLPrefix           string
	Repository          string
	Path                string
	OnAvailable         string
	available           bool
	OnHostNamePrefix    string
	hostNamePfx         bool
	OnHostNameSuffix    string
	hostNameSuffix      bool
	OnAvailableViaHTTP  bool
	DisableValidateCert bool
}

func (c *Connection) Available() bool {
	return c.available
}
func (c *Connection) HostPfx() bool {
	return c.hostNamePfx
}
func (c *Connection) HostSuffix() bool {
	return c.hostNameSuffix
}
func (c *Connection) SetAvailable(val bool) {
	c.available = val
}
func (c *Connection) SetHostPfx(val bool) {
	c.hostNamePfx = val
}
func (c *Connection) SetHostSfx(val bool) {
	c.hostNameSuffix = val
}

type Version struct {
	Version   string
	BuildDate string
}
