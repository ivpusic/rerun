package main

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var envSchemes = map[string]*regexp.Regexp{
	"windows": regexp.MustCompile("%([a-zA-Z0-9]+)%"),
	"linux":   regexp.MustCompile("$([a-zA-Z0-9]+)"),
	"darwin":  regexp.MustCompile("$([a-zA-Z0-9]+)"),
}

func convertAbsolute(path string) (string, error) {

	// if the path contains an environment variable at its beginning
	// then replace that with its value
	reg, ok := envSchemes[runtime.GOOS]
	if ok && reg.MatchString(path) {
		if matches := reg.FindAllStringSubmatch(path, -1); len(matches) > 0 {
			// replace environment variable with its value
			path = strings.Replace(path, matches[0][0], os.Getenv(matches[0][1]), 1)
		}
	}

	abs, err := filepath.Abs(path)
	if err == nil {
		return abs, nil
	} else {
		return path, err
	}
}

// takes an array of (maybe) relative paths and convert them to their absolute representatives
func convertAbsolutes(paths []string) []string {
	for ind, path := range paths {
		if newPath, err := convertAbsolute(path); err == nil {
			paths[ind] = newPath
		} else {
			logger.Errorf("Error while attempting to translate file path %q to absolute path: %q", path, err.Error())
		}
	}
	return paths
}
