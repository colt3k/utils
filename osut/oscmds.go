package osut

import (
	"errors"
	"os/exec"

	"github.com/mgutz/str"

	log "github.com/colt3k/nglog/ng"
)

var (
	// ErrNoOpenCommand : When we don't know which command to use to open a file
	ErrNoOpenCommand = errors.New("unsure what command to use to open this file")
)

// OSCommand holds all the os commands
type OSCommand struct {
	Log      *log.StdLogger
	Platform *Platform
}

// NewOSCommand os command runner
func NewOSCommand(log *log.StdLogger) (*OSCommand, error) {
	osCommand := &OSCommand{
		Log:      log,
		Platform: OS(),
	}
	return osCommand, nil
}

// RunCommandWithOutput wrapper around commands returning their output and error
func (c *OSCommand) RunCommandWithOutput(command string) (string, error) {
	flds := make([]log.Fields, 0)
	flds = append(flds, log.Fields{"command": command})
	entry := log.WithFields(flds)
	entry.Info("RunCommand")

	splitCmd := str.ToArgv(command)
	log.Logln(log.INFO, splitCmd)
	cmdOut, err := exec.Command(splitCmd[0], splitCmd[1:]...).CombinedOutput()
	return sanitisedCommandOutput(cmdOut, err)
}

// RunCommand runs a command and just returns the error
func (c *OSCommand) RunCommand(command string) error {
	_, err := c.RunCommandWithOutput(command)
	return err
}

// GetOpenCommand get open command
func (c *OSCommand) GetOpenCommand() (string, string, error) {
	//NextStep open equivalents: xdg-open (linux), cygstart (cygwin), open (OSX)
	trailMap := map[string]string{
		"xdg-open": " &>/dev/null &",
		"cygstart": "",
		"open":     "",
	}
	for name, trail := range trailMap {
		if err := c.RunCommand("which " + name); err == nil {
			return name, trail, nil
		}
	}
	return "", "", ErrNoOpenCommand
}

// OpenFile opens a file with the given
func (c *OSCommand) OpenFile(filename string) (*exec.Cmd, error) {
	cmdName, cmdTrail, err := c.GetOpenCommand()
	if err != nil {
		return nil, err
	}
	err = c.RunCommand(cmdName + " " + filename + cmdTrail)
	return nil, err
}
func sanitisedCommandOutput(output []byte, err error) (string, error) {
	outputString := string(output)
	if err != nil {
		// errors like 'exit status 1' are not very useful so we'll create an error
		// from the combined output
		return outputString, errors.New(outputString)
	}
	return outputString, nil
}
