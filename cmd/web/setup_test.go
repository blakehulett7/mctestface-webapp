package main

import (
	"log"
	"os"
	"testing"

	"github.com/blakehulett7/mctestface-webapp/pkg/db"
)

var app State

func TestMain(m *testing.M) {
	app.Session = GetSession()
	app.DSN = "host=localhost port=5433 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5"

	conn, err := app.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	app.DB = db.PostgresConn{DB: conn}

	os.Exit(m.Run())
}
