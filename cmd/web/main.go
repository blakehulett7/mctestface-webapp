package main

import (
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type State struct {
	Session *scs.SessionManager
}

func main() {
	log.Println("Dominus Iesus Christus")

	app := State{}
	app.Session = GetSession()

	router := app.Routes()

	err := http.ListenAndServe(":1000", router)
	if err != nil {
		log.Fatal(err)
	}
}
