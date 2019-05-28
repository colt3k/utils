package osut

import "github.com/gonutz/w32"

func OSVersion() (int, int) {

	v := w32.GetVersion()
	major, minor := v&0xFF, v&0xFF00>>8
	return int(major),int(minor)
}

func Distro() string {
	return "windows"
}

func FVersion() string {
	return "NA"
}