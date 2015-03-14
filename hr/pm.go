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
	args  []string
	oscmd *exec.Cmd
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

	pidStr := strings.TrimSpace(string(out[:]))
	// pid not found
	if len(pidStr) == 0 {
		return 0, nil
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, errors.New("Error while converting pid to integer! " + err.Error())
	}

	return pid, nil
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

// will run defined command
func (pm *processManager) run() {
	logger.Debug("starting process")

	pm.oscmd = exec.Command(pm.cmd, pm.args...)
	pm.oscmd.Stdout = os.Stdout
	pm.oscmd.Stdin = os.Stdin
	pm.oscmd.Stderr = os.Stderr

	err := pm.oscmd.Start()
	if err != nil {
		logger.Error(err.Error())
	}
}

func (pm *processManager) stop() {
	logger.Debug("stopping process")

	if pm.oscmd == nil {
		return
	}

	pm.oscmd.Process.Kill()
	pm.killOnPort(true)
}
