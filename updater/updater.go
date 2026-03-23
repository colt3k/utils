// Package updater defines the shared metadata used by the update backends.
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
	Version           string `json:"version"`
	Changelog         string `json:"changelog"`
	BaseURL           string
	URL               string
	ArchiveName       string
	User              []byte
	Pass              []byte
	Bearer            bool
	DisableVerifyCert bool
	Issue             string
}

// String returns a compact printable view of the remote app metadata.
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

// NewUser builds a basic authenticated update connection definition.
func NewUser(user, passOrToken, urlPrefix, repository string) *Connection {
	t := new(Connection)
	t.User = user
	t.PassOrToken = passOrToken
	t.URLPrefix = urlPrefix
	t.Repository = repository
	return t
}

// Connection describes a candidate update source and its gating rules.
type Connection struct {
	Name                string
	HostName            string
	User                string
	PassOrToken         string
	Bearer              bool
	URLPrefix           string
	Repository          string
	Path                string
	OnAvailable         string
	OnAvailableTimeout  int
	available           bool
	OnHostNamePrefix    string
	hostNamePfx         bool
	OnHostNameSuffix    string
	hostNameSuffix      bool
	OnAvailableViaHTTP  bool
	DisableValidateCert bool
	AQLSupport          bool
}

// Available reports whether the connection passed availability checks.
func (c *Connection) Available() bool {
	return c.available
}

// HostPfx reports whether hostname prefix gating matched.
func (c *Connection) HostPfx() bool {
	return c.hostNamePfx
}

// HostSuffix reports whether hostname suffix gating matched.
func (c *Connection) HostSuffix() bool {
	return c.hostNameSuffix
}

// SetAvailable stores the current availability check result.
func (c *Connection) SetAvailable(val bool) {
	c.available = val
}

// SetHostPfx stores the current hostname prefix match state.
func (c *Connection) SetHostPfx(val bool) {
	c.hostNamePfx = val
}

// SetHostSfx stores the current hostname suffix match state.
func (c *Connection) SetHostSfx(val bool) {
	c.hostNameSuffix = val
}

// Version holds the local version metadata used for comparisons.
type Version struct {
	Version   string
	BuildDate string
}
