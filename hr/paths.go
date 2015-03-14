package main

import (
	"fmt"
	"runtime"
)

func convertAbsolute(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err == nil {
		return abs, nil
	} else {
		return path, err
	}
}

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
