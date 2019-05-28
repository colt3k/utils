package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/colt3k/nglog/ng"
)

func main() {

	var serviceName string
	var serviceRunAs string
	var removeService bool
	flag.StringVar(&serviceName, "service", "", "the service to install")
	flag.StringVar(&serviceRunAs, "runas", "nobody", "user to run service as")
	flag.BoolVar(&removeService, "remove", false, "remove service binary and service spec file")
	flag.Parse()

	if serviceName == "" {
		log.Logf(log.FATAL, "service name missing")
	}

	if removeService {
		log.Println("removing service", serviceName)
		_,err := SudoNoFail("systemctl", "stop", serviceName)
		if err != nil {
			log.Logf(log.ERROR, "sudo issue %+v", err)
		}
		_,err = SudoNoFail("systemctl", "disable", serviceName)
		if err != nil {
			log.Logf(log.ERROR, "sudo issue %+v", err)
		}
		_,err = SudoNoFail("rm", serviceBinary(serviceName))
		if err != nil {
			log.Logf(log.ERROR, "sudo issue %+v", err)
		}
		_,err = SudoNoFail("rm", serviceDefinition(serviceName))
		if err != nil {
			log.Logf(log.ERROR, "sudo issue %+v", err)
		}
		return
	}

	log.Println("installing service", serviceName)

	_,err := SudoNoFail("systemctl", "stop", serviceName)
	if err != nil {
		log.Logf(log.ERROR, "sudo issue %+v", err)
	}
	_,err = SudoNoFail("systemctl", "disable", serviceName)
	if err != nil {
		log.Logf(log.ERROR, "sudo issue %+v", err)
	}

	source := fmt.Sprintf("%s/%s.go", serviceName, serviceName)
	_,err = Command("go", "build", "-o", serviceBinary(serviceName), source)
	if err != nil {
		log.Logf(log.ERROR, "sudo issue %+v", err)
	}

	templatePath := writeServiceFile(serviceName, serviceRunAs)

	_,err = Sudo("systemctl", "enable", templatePath)
	if err != nil {
		log.Logf(log.ERROR, "sudo issue %+v", err)
	}
	_,err = Sudo("systemctl", "start", serviceName)
	if err != nil {
		log.Logf(log.ERROR, "sudo issue %+v", err)
	}
}

func serviceBinary(serviceName string) string {
	return fmt.Sprintf("/opt/gosvc/%s", serviceName)
}

func serviceDefinition(serviceName string) string {
	return filepath.Join("/opt/gosvc", fmt.Sprintf("%s.service", serviceName))
}

func writeServiceFile(serviceName, serviceRunAs string) string {
	t := template.Must(template.New("tmpl").Parse(ServiceTemplate))
	templatePath := serviceDefinition(serviceName)
	file, err := os.OpenFile(templatePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	PanicOnError(err)
	err = t.Execute(file, map[string]string{"ServiceName": serviceName, "ServiceRunAs": serviceRunAs})
	PanicOnError(err)
	file.Close()
	return templatePath
}

var ServiceTemplate = `[Unit]
Description={{.ServiceName}} Service
After=network.target

[Service]
Type=simple
User={{.ServiceRunAs}}
WorkingDirectory=/opt/gosvc
ExecStart=/opt/gosvc/{{.ServiceName}}
Restart=on-failure

[Install]
WantedBy=multi-user.target
`

func SudoNoFail(cmd ...string) ([]byte, error) {
	command := exec.Command("sudo", cmd...)
	out, err := command.CombinedOutput()
	return out, err
}
func Sudo(cmd ...string) ([]byte, error) {
	command := exec.Command("sudo", cmd...)
	out, err := command.CombinedOutput()
	return out, err
}
func Command(app string, cmd ...string) ([]byte, error) {
	command := exec.Command(app, cmd...)
	out, err := command.CombinedOutput()
	return out, err
}
func PanicOnError(err error) {
	panic(err)
}
