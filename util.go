package main

import (
	"fmt"
	"testing"
)

func contains(arr []string, path string) bool {
	for _, abs := range arr {
		if path == abs {
			return true
		}
	}

	return false
}

func AssertArraysEq(t *testing.T, a, b []string) {
	if a == nil && b == nil {
		return
	}

	failMsg := fmt.Sprintf("%v and %v not equal", a, b)

	if a == nil || b == nil {
		t.Error(failMsg)
	}

	if len(a) != len(b) {
		t.Error(failMsg)
	}

	for i := range a {
		if a[i] != b[i] {
			t.Error(failMsg)
		}
	}

	return
}
