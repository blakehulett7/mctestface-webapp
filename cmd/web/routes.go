package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *State) Routes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(app.AddIpToContext)
	router.Use(app.Session.LoadAndSave)

	router.Get("/", app.Home)
	router.Post("/login", app.Login)

	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../static"))))

	return router
}
