package main

import (
	"github.com/howeyc/fsnotify"
	"github.com/ivpusic/golog"
	"gopkg.in/alecthomas/kingpin.v1"
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

func makeAbsPaths(paths []string) []string {
	abspaths := make([]string, len(paths), (cap(paths)+1)*2)

	for ind, path := range paths {
		abs, err := filepath.Abs(path)

		if err == nil {
			path = abs
		}

		abspaths[ind] = abs
	}

	return abspaths
}

func contains(arr []string, path string) bool {
	for _, abs := range arr {
		if strings.Index(path, abs) == 0 {
			return true
		}
	}

	return false
}

func log(msg string) {
	if verbose {
		logger.Info(msg)
	}
}

func logErr(msg string) {
	if verbose {
		logger.Error(msg)
	}
}

func main() {
	kingpin.Version("0.0.2")
	kingpin.Parse()

	cmd := []string{"go", "run", "main.go"}
	watch := []string{"."}
	ignore := []string{}
	port := 3000

	confPath := *_conf
	if len(confPath) > 0 {
		conf, err := parseConf(confPath)
		if err != nil {
			logErr(err.Error())
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

	watch = makeAbsPaths(watch)
	ignore = makeAbsPaths(ignore)

	pm := processManager{
		port: port,
		cmd:  cmd[0],
		args: cmd[1:],
	}

	// if there is old process listening on specified port, kill it
	pm.killOnPort(false)

	log("will run with" + strings.Join(cmd, " ") + " command")
	log("will listen on port " + strconv.Itoa(port))
	for _, val := range watch {
		log("watching: " + val)
	}

	for _, val := range ignore {
		log("ignoring: " + val)
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
			case ev := <-watcher.Event:
				file := ev.Name

				if strings.Contains(file, ".go") {
					abs, err := filepath.Abs(file)
					if err == nil && contains(watch, abs) && !contains(ignore, abs) {
						log("reloading...")
						pm.stop()
						pm.run()
					} else {
						log("ignoring change on file: " + file)
					}
				}
			case err := <-watcher.Error:
				logErr("error: " + err.Error())
				done <- true
			}
		}
	}()

	for _, val := range watch {
		err = watcher.Watch(val)
		if err != nil {
			logErr("error: " + err.Error())
		}
	}

	pm.run()
	logger.Info("go-hotreload started.")
	logger.Info("Watching for changes.")
	logger.Info("Use -v flag for verbose mode")

	<-done
	watcher.Close()
}
