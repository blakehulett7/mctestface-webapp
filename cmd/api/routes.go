package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Bridge) Routes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	//router.Use(app.EnableCORS)

	return router
}
