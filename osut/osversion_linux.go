package osut

import (
	"fmt"
	"os/exec"
	"strings"
	"io/ioutil"

	"github.com/colt3k/utils/mathut"
)

// OSVersion
func OSVersion() (int, int){

	cmd := exec.Command("/bin/sh", "-c", "grep -oE '[0-9]+\\.[0-9]+' /etc/system-release")
	r, err := cmd.StdoutPipe()
	err = cmd.Start()
	out, err := ioutil.ReadAll(r)
	err = cmd.Wait()
	if err != nil {
		fmt.Println("getInfo:",err)
	}
	parts := strings.Split(string(out),".")

	return int(mathut.ParseInt(parts[0])),int(mathut.ParseInt(parts[1]))

}

func Distro() string {
	cmd := exec.Command("/bin/sh", "-c", "grep -oE '^ID=.*' /etc/os-release")
	r, err := cmd.StdoutPipe()
	err = cmd.Start()
	out, err := ioutil.ReadAll(r)
	err = cmd.Wait()
	if err != nil {
		fmt.Println("err",err)
	}
	parts := strings.Split(string(out), "=")
	part := strings.TrimSpace(parts[1])
	if part[0] == '"' {
		part = part[1:(len(part)-1)]
	}
	return part
}

func FVersion() string {
	cmd := exec.Command("/bin/sh", "-c", "grep -oE '^VERSION_ID=.*' /etc/os-release")
	r, err := cmd.StdoutPipe()
	err = cmd.Start()
	out, err := ioutil.ReadAll(r)
	err = cmd.Wait()
	if err != nil {
		fmt.Println("err",err)
	}
	parts := strings.Split(string(out), "=")
	part := strings.TrimSpace(parts[1])
	if part[0] == '"' {
		part = part[1:(len(part)-1)]
	}
	return part
}
