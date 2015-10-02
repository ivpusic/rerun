package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	paths := []string{"/path/a", "/path/b", "/path/c"}
	assert.True(t, contains(paths, "/path/a"))
	assert.False(t, contains(paths, "/path/a/b"))
	assert.False(t, contains(paths, "/path/d"))
}

func TestArraysEq(t *testing.T) {
	a1 := []string{"a", "b"}
	a2 := []string{"a", "b", "c"}
	a3 := []string{"a", "b"}

	mock := &testing.T{}
	AssertArraysEq(mock, a1, a2)
	assert.True(t, mock.Failed(), "Should fail when arrays are not the same")

	mock = &testing.T{}
	AssertArraysEq(mock, a1, a3)
	assert.False(t, mock.Failed(), "Should not fail when arrays are the same")
}
