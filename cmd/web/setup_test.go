package main

import (
	"os"
	"testing"

	"github.com/blakehulett7/mctestface-webapp/pkg/repository/dbrepo"
)

var app State

func TestMain(m *testing.M) {
	app.Session = GetSession()
	app.DB = &dbrepo.TestDBRepo{}

	os.Exit(m.Run())
}
