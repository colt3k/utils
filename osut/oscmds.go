package osut

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

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

// CallCmd execute local apps
func CallCmd(command string) (string, error) {
	log.Logf(log.DBGL2, "sh -c %v", command)
	cmd := exec.Command("sh", "-c", command)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("issue processing call (%v)\n%v", command, err)
	}
	output := string(stdoutStderr)
	return output, nil
}

// CallCmdNoWait execute app if not already running
func CallCmdNoWait(command, app string) (int, error) {
	log.Logf(log.DBGL2, "sh -c %v", command)
	cmd := exec.Command("sh", "-c", command)
	err := cmd.Start()
	if err != nil {
		return 0, err
	}

	pid := cmd.Process.Pid
	// sleep give process time to spin up before app exits
	time.Sleep(5 * time.Second)

	id, found := FindProcess(app)
	log.Logf(log.INFO, "id %v, found %v", id, found)
	if found {
		pid = id
	}

	return pid, nil
}

// FindProcess find by application name return pid and if running true/false
func FindProcess(appName string) (int, bool) {
	pCurPid := os.Getppid()
	curPid := os.Getpid()
	log.Logf(log.DBGL2, "> Parent Pid: %v, Cur Pid: %v", pCurPid, curPid)
	cmd := exec.Command("sh", "-c", "ps aux | grep "+appName+" | grep -v grep")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		var exitErr *exec.ExitError
		switch {
		case errors.As(err, &exitErr):
			exitCode := exitErr.ExitCode()
			if exitCode == 1 {
				return curPid, false
			}
		default:
			log.Logf(log.ERROR, "!!! %v", err)
		}
	}
	output := string(stdoutStderr)
	trim := strings.TrimSpace(output)
	log.Logf(log.DBGL2, "findProcess data : %v", trim)

	if len(trim) > 0 {
		pids := make(map[int]int)
		lines := strings.Split(trim, "\n")
		for _, j := range lines {
			fields := strings.Fields(j)
			// Ensure the process cmd starts with this command and isn't included in another
			// search each field for appName, could be a diff field based on OS or full path with name
			for l, m := range fields {
				if l >= 10 {
					if strings.Contains(m, appName) {
						fP, _ := strconv.Atoi(fields[1])
						pids[fP] = fP
					}
				}
			}
		}
		log.Logf(log.DBGL2, "> Cur Pid: %v, PS found: %v", curPid, pids)
		// If found in map and len of map is more than one we found it running don't delete lock
		pid := 0
		for k := range pids {
			if pids[k] != curPid && pids[k] != pCurPid {
				pid = k
			}
		}
		if pid != 0 {
			return pid, true
		}
		return curPid, false
	}
	return curPid, false
}
