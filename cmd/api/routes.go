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

	router.Post("/authenticate", app.Authenticate)
	router.Post("/refresh", app.Refresh)

	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		payload := struct {
			Message string `json:"message"`
		}{"Gratia Plena"}

	})

	router.Route("/users", func(mux chi.Router) {
		mux.Get("/", app.GetAllUsers)
		mux.Post("/", app.CreateUser)

		mux.Get("/{userId}", app.GetUserDetails)
		mux.Put("/{userId}", app.UpdateUser)
		mux.Delete("/{userId}", app.DeleteUser)
	})

	return router
}
