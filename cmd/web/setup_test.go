package main

import (
	"os"
	"testing"
)

var app State

func TestMain(m *testing.M) {
	app.Session = GetSession()

	os.Exit(m.Run())
}
