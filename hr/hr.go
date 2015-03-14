package main

import (
	"fmt"
	"github.com/ivpusic/golog"
	"gopkg.in/alecthomas/kingpin.v1"
	"gopkg.in/fsnotify.v1"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	_verbose = kingpin.Flag("verbose", "Verbose mode. It will show internal messages of go-hotreload. Default: false").Short('v').Bool()
	_cmd     = kingpin.Flag("cmd", "Command to execute on each reload. Default: 'go run main.go'").Short('c').String()
	_watch   = kingpin.Flag("watch", "Comma separated list of directories to watch. Default: ['.']").Short('w').String()
	_ignore  = kingpin.Flag("ignore", "Comma separated list of directories to ignore. Default: []").Short('i').String()
	_port    = kingpin.Flag("port", "Port on which app is running. Default: 3000").Short('p').Int()
	_conf    = kingpin.Flag("conf", "Path to json config. Default: ''").String()
	verbose  = false
	logger   = golog.GetLogger("github.com/ivpusic/go-hotreload/hr")
)

func contains(arr []string, path string) bool {
	for _, abs := range arr {
		if strings.Index(path, abs) == 0 {
			return true
		}
	}

	return false
}

func main() {
	kingpin.Version("0.0.3")
	kingpin.Parse()

	cmd := []string{"go", "run", "main.go"}
	watch := []string{"."}
	ignore := []string{}
	port := 3000

	confPath := *_conf
	if len(confPath) > 0 {
		conf, err := parseConf(confPath)

		if err != nil {
			logger.Error(err.Error())
			if conf == nil {
				logger.Error("Terminating due to missing configuration file")
				return
			}
		}

		cmd = strings.Split(conf.Cmd, " ")
		watch = conf.Watch
		ignore = conf.Ignore
		port = conf.Port
		verbose = conf.Verbose
	}

	if len(*_cmd) > 0 {
		cmd = strings.Split(*_cmd, " ")
	}

	if len(*_watch) > 0 {
		watch = strings.Split(*_watch, ",")
	}

	if len(*_ignore) > 0 {
		ignore = strings.Split(*_ignore, ",")
	}

	if *_port > 0 {
		port = *_port
	}

	if *_verbose {
		verbose = *_verbose
	}

	if verbose {
		logger.Level = golog.DEBUG
	} else {
		logger.Level = golog.INFO
	}

	watch = convertAbsolutes(watch)
	ignore = convertAbsolutes(ignore)

	pm := processManager{
		port: port,
		cmd:  cmd[0],
		args: cmd[1:],
	}

	// if there is old process listening on specified port, kill it
	pm.killOnPort(false)

	logger.Debug("will run with" + strings.Join(cmd, " ") + " command")
	logger.Debug("will listen on port " + strconv.Itoa(port))
	for _, val := range watch {
		logger.Debug("watching: " + val)
	}

	for _, val := range ignore {
		logger.Debug("ignoring: " + val)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Panic(err.Error())
	}

	done := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				file := ev.Name
				if strings.HasSuffix(file, ".go") {
					abs, err := filepath.Abs(file)
					if err == nil && contains(watch, abs) && !contains(ignore, abs) {
						logger.Debug("reloading...")
						pm.stop()
						pm.run()
					} else {
						logger.Debug("ignoring change on file: " + file)
					}
				}
			case err := <-watcher.Errors:
				logger.Error("error: " + err.Error())
				done <- true
			}
		}
	}()

	for _, val := range watch {
		err = watcher.Add(val)
		if err != nil {
			logger.Error("error: " + err.Error())
		}
	}

	pm.run()
	logger.Info("go-hotreload started.")
	logger.Info("Watching for changes.")

	if !verbose {
		logger.Info("Use -v flag for verbose mode")
	}

	<-done
	watcher.Close()
}
