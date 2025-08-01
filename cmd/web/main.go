package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/blakehulett7/mctestface-webapp/pkg/db"
)

type State struct {
	DB      db.PostgresConn
	DSN     string
	Session *scs.SessionManager
}

func main() {
	log.Println("Dominus Iesus Christus")

	app := State{}

	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5433 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection")
	flag.Parse()
	conn, err := app.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}

	app.DB = db.PostgresConn{conn}
	app.Session = GetSession()

	router := app.Routes()

	err = http.ListenAndServe(":1000", router)
	if err != nil {
		log.Fatal(err)
	}
}
