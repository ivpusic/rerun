package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type config struct {
	Port   int
	Watch  []string
	Ignore []string
	Cmd    string
}

func parseConf(path string) (*config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	conf := &config{}
	err = json.Unmarshal(content, &conf)
	if err != nil {
		return nil, errors.New("Error while running unmarshal! " + err.Error())
	}

	return conf, nil
}
