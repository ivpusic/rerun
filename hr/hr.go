package main

import (
	"github.com/howeyc/fsnotify"
	"gopkg.in/alecthomas/kingpin.v1"
	"log"
	"path/filepath"
	"strings"
)

var (
	_cmd    = kingpin.Flag("cmd", "Command to execute on each reload. Default: 'go run main.go'").String()
	_watch  = kingpin.Flag("watch", "Comma separated list of directories to watch. Default: ['.']").String()
	_ignore = kingpin.Flag("ignore", "Comma separated list of directories to ignore. Default: []").String()
	_port   = kingpin.Flag("port", "Port on which app is running. Default: 3000").Int()
	_conf   = kingpin.Flag("conf", "Path to json config. Default: ''").String()
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

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	cmd := []string{"go", "run", "main.go"}
	watch := []string{"."}
	ignore := []string{}
	port := 3000

	confPath := *_conf
	if len(confPath) > 0 {
		conf, err := parseConf(confPath)
		if err != nil {
			log.Fatal(err)
		}

		cmd = strings.Split(conf.Cmd, " ")
		watch = conf.Watch
		ignore = conf.Ignore
		port = conf.Port

	}

	watch = makeAbsPaths(watch)
	ignore = makeAbsPaths(ignore)

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

	pm := processManager{
		port: port,
		cmd:  cmd[0],
		args: cmd[1:],
	}

	// if there is old process listening on specified port, kill it
	pm.killOnPort(false)

	log.Printf("will run with '%s' command", cmd)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
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
						log.Println("reloading...")
						pm.stop()
						pm.run()
					}
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
				done <- true
			}
		}
	}()

	for _, val := range watch {
		err = watcher.Watch(val)
		if err != nil {
			log.Panic(err)
		}
	}

	pm.run()
	log.Print("watching for changes...")

	<-done
	watcher.Close()
}
