package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
