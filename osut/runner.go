package osut

import (
	"os/exec"

	"os"

	log "github.com/colt3k/nglog/ng"
)

//RunAppWithPath run the application with the specific path
func RunAppWithPath(app, path string) ([]byte, error) {
	log.Logf(log.INFO, "Executing...%s %s\n", app, path)
	fnd, err := exec.LookPath(app)
	if err != nil {
		return nil, err
	}

	cmd, err := exec.Command(fnd, path).CombinedOutput()
	log.Logln(log.INFO, "Completed Execution.")
	return cmd, err
}

//RunAppNoPath run the application wherever you can find it
func RunAppNoPath(app string) error {
	log.Logln(log.INFO, "Starting Open...")
	fnd, err := exec.LookPath(app)
	if err != nil {
		return err
	}

	cmd := exec.Command(fnd)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Start()
}

//RunAppMac run application on mac
func RunAppMac(app string) error {
	log.Logln(log.INFO, "Starting RunAppMac...")
	fnd, err := exec.LookPath("/usr/bin/open")
	if err != nil {
		return err
	}
	log.Printf("Available at %s\n", fnd)

	cmd := exec.Command(fnd, "-a", app)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Start()
}
