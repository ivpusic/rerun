package main

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type processManager struct {
	port  int
	cmd   string
	args  []string
	oscmd *exec.Cmd
}

// will try to find PID of process which is runing on defined port
// if it succeed it will kill it
func (pm *processManager) killOnPort(showerr bool) {
	addr := ":" + strconv.Itoa(pm.port)
	proc := exec.Command("lsof", "-t", "-i", addr, "-s", "TCP:LISTEN")

	out, err := proc.Output()
	if err != nil {
		if showerr {
			logErr("Error while executing fuser command! " + err.Error())
		}
		return
	}

	_pid := strings.TrimSpace(string(out[:]))
	pid, err := strconv.Atoi(_pid)
	if err != nil {
		if showerr {
			logErr("Error while converting pid to integer! " + err.Error())
		}
		return
	}

	pidProc, err := os.FindProcess(pid)
	if err != nil {
		if showerr {
			logErr("Error while finding process with pid " + _pid + "! " + err.Error())
		}
		return
	}

	pidProc.Kill()
}

// will run defined command
func (pm *processManager) run() {
	log("starting process")

	pm.oscmd = exec.Command(pm.cmd, pm.args...)
	pm.oscmd.Stdout = os.Stdout
	pm.oscmd.Stdin = os.Stdin
	pm.oscmd.Stderr = os.Stderr

	err := pm.oscmd.Start()
	if err != nil {
		logErr(err.Error())
	}
}

func (pm *processManager) stop() {
	log("stopping process")

	if pm.oscmd == nil {
		return
	}

	pm.oscmd.Process.Kill()
	pm.killOnPort(true)
}
