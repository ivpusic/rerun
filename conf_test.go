package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConf(t *testing.T) {
	*ignore = ""
	*args = ""
	*suffixes = ""

	cnf, err := parseConf("./test/config.json")
	assert.Nil(t, err, "Error should not happen if file is ok")

	AssertArraysEq(t, []string{"some/path/to/ignore", "some/path/to/ignore"}, cnf.Ignore)
	AssertArraysEq(t, []string{"dev", "test"}, cnf.Args)
	AssertArraysEq(t, []string{".go", ".html", ".tpl"}, cnf.Suffixes)
	assert.Equal(t, true, cnf.Attrib)
}

func TestParseConfWithCliArgs(t *testing.T) {
	*ignore = "path1,path2"
	*args = "arg1,arg2"
	*suffixes = ".go,.html"
	*attrib = false

	cnf, err := loadConfiguration()
	assert.Nil(t, err, "Error should not happen if file is ok")

	assert.True(t, len(cnf.Ignore) >= 2)
	pathPrefix := os.Getenv("GOPATH") + "/src/github.com/ivpusic/rerun"

	assert.True(t, contains(cnf.Ignore, pathPrefix + "/path1"))
	assert.True(t, contains(cnf.Ignore, pathPrefix + "/path2"))
	AssertArraysEq(t, []string{"arg1", "arg2"}, cnf.Args)
	AssertArraysEq(t, []string{".go", ".html"}, cnf.Suffixes)

	assert.False(t, cnf.Attrib)
}
