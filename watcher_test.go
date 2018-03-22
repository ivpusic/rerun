package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/fsnotify/fsnotify.v1"
)

func TestIsFileImportant(t *testing.T) {
	assert.True(t, true)

	cnf := &config{build: "some-build", Suffixes: []string{".go"}, Ignore: []string{"some/file.c"}}
	pm := &processManager{conf: cnf}
	watcher := watcher{pm: pm}

	// fail if event is not supported
	event := fsnotify.Event{"some/file.go", fsnotify.Chmod}
	assert.False(t, watcher.isEventImportant(event))

	// fail if ext is not supported
	event = fsnotify.Event{"some/file.html", fsnotify.Write}
	assert.False(t, watcher.isEventImportant(event))

	// fail if file is ignored
	event = fsnotify.Event{"some/file.c", fsnotify.Write}
	assert.False(t, watcher.isEventImportant(event))

	// Test attrib event
	event = fsnotify.Event{"some/file.go", fsnotify.Chmod}
	assert.False(t, watcher.isEventImportant(event))

	// Test attrib event in case settings are different
	pm.conf.Attrib = true
	event = fsnotify.Event{"some/file.go", fsnotify.Chmod}
	assert.True(t, watcher.isEventImportant(event))

	// all ok
	event = fsnotify.Event{"some/file.go", fsnotify.Write}
	assert.True(t, watcher.isEventImportant(event))
}
