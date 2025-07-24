package main

import (
	"log"
	"net/http"
)

type State struct {
}

func main() {
	log.Println("Jesus is Lord!")

	state := State{}
	router := state.Routes()

	err := http.ListenAndServe(":1000", router)
	if err != nil {
		log.Fatal(err)
	}
}
