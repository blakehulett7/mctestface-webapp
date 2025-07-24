package main

import (
	"log"
	"net/http"
)

type State struct {
}

func main() {
	log.Println("Dominus Iesus Christus")

	app := State{}
	router := app.Routes()

	err := http.ListenAndServe(":1000", router)
	if err != nil {
		log.Fatal(err)
	}
}
