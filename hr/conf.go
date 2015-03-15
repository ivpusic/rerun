package main

import (
	"encoding/json"
	"fmt"
	"github.com/ivpusic/golog"
	"io/ioutil"
	"strings"
)

type config struct {
	Port    int
	Watch   []string
	Ignore  []string
	Cmd     string
	Verbose bool
}

func parseConf(path string) (*config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Did not find specified configuration file %q", path)
	}

	conf := &config{}
	err = json.Unmarshal(content, &conf)
	if err != nil {
		return nil, fmt.Errorf("Error while unmarshaling %q", err.Error())
	}

	return conf, nil
}

func loadConfiguration() {
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

		cmd = conf.Cmd
		watch = conf.Watch
		ignore = conf.Ignore
		port = conf.Port
		verbose = conf.Verbose
	}

	if len(*_cmd) > 0 {
		cmd = *_cmd
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
}
