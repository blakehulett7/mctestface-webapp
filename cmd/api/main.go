package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/blakehulett7/mctestface-webapp/pkg/repository"
	"github.com/blakehulett7/mctestface-webapp/pkg/repository/dbrepo"
)

const port = 8090

type Bridge struct {
	DSN       string
	DB        repository.DatabaseRepo
	Domain    string
	JWTSecret string
}

func main() {
	var app Bridge

	// This flag package is very useful for command line tools
	flag.StringVar(&app.Domain, "domain", "example.com", "Domain for app")
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5433 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "Sancta Maria", "signing secret")
	flag.Parse()

	conn, err := app.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	app.DB = &dbrepo.PostgresDBRepo{DB: conn}

	log.Println("In nomine Patris, et Filii, et Spiritus Sancti")

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.Routes())
	if err != nil {
		log.Println(err)
	}
}
