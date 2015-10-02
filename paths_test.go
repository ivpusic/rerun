package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertAbsolutes(t *testing.T) {
	paths := []string{"test"}
	goPath := os.Getenv("GOPATH")
	expected := goPath + "/src/github.com/ivpusic/rerun/test"
	absolutes := convertAbsolutes(paths)

	assert.Equal(t, expected, absolutes[0])
}

func TestParseGlobs(t *testing.T) {
	globs := []string{"test/*"}
	expected := "test/config.json"
	parsedGlobs := parseGlobs(globs)

	assert.Equal(t, expected, parsedGlobs[0])
}
