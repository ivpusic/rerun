package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	verbose  = kingpin.Flag("verbose", "Verbose mode. It will show rerun internal messages. Default: false").Short('v').Bool()
	ignore   = kingpin.Flag("ignore", "List of ignored files and directories.").Default("").Short('i').String()
	args     = kingpin.Flag("args", "Application arguments.").Default("").Short('a').String()
	suffixes = kingpin.Flag("suffixes", "File suffixes to watch.").Short('s').String()
	confPath = kingpin.Flag("config", "JSON configuration location").Short('c').String()
	attrib   = kingpin.Flag("attrib", "Also watch attribute changes").Short('b').Bool()
)

type config struct {
	Ignore   []string
	Args     []string
	Suffixes []string
	build    string
}

func parseConf(path string) (*config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Did not find specified configuration file %s", path)
	}

	conf := &config{}
	err = json.Unmarshal(content, &conf)
	if err != nil {
		return nil, fmt.Errorf("Error while unmarshaling %s", err.Error())
	}

	return conf, nil
}

func loadConfiguration() (*config, error) {
	if !TEST_MODE {
		kingpin.Version("1.0.0")
		kingpin.Parse()
	}

	conf := &config{}
	var err error

	if len(*confPath) > 0 {
		conf, err = parseConf(*confPath)

		if err != nil {
			return nil, err
		}
	}

	if len(*ignore) > 0 {
		conf.Ignore = append(conf.Ignore, strings.Split(*ignore, ",")...)
	}

	if len(*args) > 0 {
		conf.Args = append(conf.Args, strings.Split(*args, ",")...)
	}

	if len(*suffixes) > 0 {
		conf.Suffixes = append(conf.Suffixes, strings.Split(*suffixes, ",")...)
	}

	if len(conf.Suffixes) == 0 {
		conf.Suffixes = append(conf.Suffixes, ".go")
	}

	buildName := "application-build-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	if runtime.GOOS == "windows" {
		buildName += ".exe"
	}

	conf.build, err = convertAbsolute(buildName)
	if err != nil {
		return nil, err
	}

	// ignore build files
	conf.Ignore = append(conf.Ignore, conf.build)

	// make absolute paths out of ignored files
	conf.Ignore = parseGlobs(conf.Ignore)
	conf.Ignore = convertAbsolutes(conf.Ignore)

	return conf, nil
}
