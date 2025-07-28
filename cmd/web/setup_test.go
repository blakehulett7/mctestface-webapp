package main

import (
	"os"
	"testing"
)

var app State

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
