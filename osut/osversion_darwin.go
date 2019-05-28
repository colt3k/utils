package osut

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	log "github.com/colt3k/nglog/ng"
	"github.com/colt3k/utils/mathut"
)

// OSVersion Darwin 17.7.0 x86_64
func OSVersion() (int, int){

	cmd := exec.Command("uname","-r")
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Logln(log.ERROR, "os version:",err)
	}

	parts := strings.Split(strings.TrimSpace(out.String()),".")

	return int(mathut.ParseInt(parts[0])),int(mathut.ParseInt(parts[1]))
}

func Distro() string {
	cmd := exec.Command("/bin/sh", "-c", "uname -s")
	r, err := cmd.StdoutPipe()
	err = cmd.Start()
	out, err := ioutil.ReadAll(r)
	err = cmd.Wait()
	if err != nil {
		fmt.Println("err",err)
	}
	return string(out)
}

func FVersion() string {
	cmd := exec.Command("/bin/sh", "-c", "uname -r")
	r, err := cmd.StdoutPipe()
	err = cmd.Start()
	out, err := ioutil.ReadAll(r)
	err = cmd.Wait()
	if err != nil {
		fmt.Println("err",err)
	}
	return string(out)
}