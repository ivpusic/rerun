package main

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type processManager struct {
	port  int
	cmd   string
	oscmd []*exec.Cmd
}

func (pm *processManager) getPid() (int, error) {
	addr := ":" + strconv.Itoa(pm.port)
	var command *exec.Cmd

	if runtime.GOOS == "windows" {
		scripts := os.Getenv("GOPATH") + "/src/github.com/ivpusic/go-hotreload/hr/scripts"
		command = exec.Command(scripts+"/getpid_win.bat", addr)
	} else {
		command = exec.Command("lsof", "-t", "-i", addr, "-s", "TCP:LISTEN")
	}

	out, err := command.Output()
	if err != nil {
		return 0, errors.New("Error while executing command! " + err.Error())
	}

	// pid not found
	if pidStr := strings.TrimSpace(string(out[:])); len(pidStr) == 0 {
		return 0, nil
	} else {
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			return 0, errors.New("Error while converting pid to integer! " + err.Error())
		} else {
			return pid, nil
		}
	}

}

// will try to find PID of process which is runing on defined port
// if it succeed it will kill it
func (pm *processManager) killOnPort(showerr bool) {
	pid, err := pm.getPid()
	if err != nil {
		if showerr {
			logger.Error("Error while finding process to kill! " + err.Error())
		}
		return
	}

	if pid == 0 {
		logger.Debug("PID not found!")
		return
	}

	pidProc, err := os.FindProcess(pid)
	if err != nil {
		if showerr {
			logger.Error("Error while finding process to kill! " + err.Error())
		}
		return
	}

	pidProc.Kill()
}

func removeEmpty(cmd []string) []string {
	newInd := 0
	newCmd := make([]string, len(cmd))
	for _, str := range cmd {
		if len(strings.TrimSpace(str)) != 0 {
			newCmd[newInd] = str
			newInd++
		}
	}
	return newCmd
}

// will run defined command
func (pm *processManager) run() {
	logger.Debug("starting process")

	cmds := strings.Split(pm.cmd, "&")
	pm.oscmd = make([]*exec.Cmd, len(cmds))

	for ind, command := range cmds {
		split := removeEmpty(strings.Split(command, " "))
		pm.oscmd[ind] = exec.Command(split[0], split[1:]...)
		pm.oscmd[ind].Stdout = os.Stdout
		pm.oscmd[ind].Stdin = os.Stdin
		pm.oscmd[ind].Stderr = os.Stderr
		go func(ind int) {
			err := pm.oscmd[ind].Start()
			if err != nil {
				logger.Error(err.Error())
			}
		}(ind)
	}
}

func (pm *processManager) stop() {
	logger.Debug("stopping process")

	if pm.oscmd == nil {
		return
	}

	for _, cmd := range pm.oscmd {
		cmd.Process.Kill()
	}
	pm.killOnPort(true)
}
