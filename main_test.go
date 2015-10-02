package main

import (
	"flag"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	TEST_MODE = true
	flag.Parse()
	os.Exit(m.Run())
}
