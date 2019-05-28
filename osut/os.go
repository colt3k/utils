package osut

import (
	"os"
	"runtime"
	"strings"

	log "github.com/colt3k/nglog/ng"
)

// Platform stores the os state
type Platform struct {
	OS           string
	OSAbbrv      string
	Arch         string
	Shell        *Shell
	VersionMajor int
	VersionMinor int
}
type Shell struct {
	shell        string
	shellArg     string
	escapedQuote string
}

func OS() *Platform {
	s := shellCmd()

	return &Platform{
		OS:           runtime.GOOS,
		OSAbbrv:      string([]rune(runtime.GOOS)[:2]),
		Arch:         runtime.GOARCH,
		Shell:        s,
		VersionMajor: OSVersionMaj(),
		VersionMinor: OSVersionMinor(),
	}
}

func shellCmd() *Shell {
	if Windows() {
		return &Shell{
			shell:        "cmd",
			shellArg:     "/c",
			escapedQuote: "\\\"",
		}
	} else if Linux() {
		return &Shell{
			shell:        "bash",
			shellArg:     "-c",
			escapedQuote: "\"",
		}
	} else if Mac() {
		return &Shell{
			shell:        "bash",
			shellArg:     "-c",
			escapedQuote: "\"",
		}
	}
	return &Shell{}
}

//Windows is this windows?
func Windows() bool {
	if runtime.GOOS == "windows" {
		return true
	}
	return false
}

//Linux is this Linux
func Linux() bool {
	if runtime.GOOS == "linux" {
		return true
	}
	return false
}

//Mac is this a mac
func Mac() bool {
	if runtime.GOOS == "darwin" {
		return true
	}
	return false
}

func Android() bool {
	if runtime.GOOS == "android" {
		return true
	}
	return false
}
func Hostname() string {
	host, err := os.Hostname()
	if err != nil {
		log.Logf(log.ERROR, "issue retrieving hostname\n%+v", err)
	}
	return host
}

func OSVersionMinor() int {
	_, min := OSVersion()
	return min
}
func OSVersionMaj() int {
	maj, _ := OSVersion()
	return maj
}

func OSDistro() string {
	return strings.ToLower(strings.TrimSpace(Distro()))
}

func FullVer() string {
	return strings.ToLower(strings.TrimSpace(FVersion()))
}