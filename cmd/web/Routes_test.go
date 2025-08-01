package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_app_Routes(t *testing.T) {
	registered := []struct {
		Route  string
		Method string
	}{
		{"/", "GET"},
		{"/login", "POST"},
		{"/user/profile", "GET"},
		{"/static/", "GET"},
	}

	router := app.Routes()

	chi_routes := router.(chi.Routes)

	for _, route := range registered {
		if !routeExists(route.Route, route.Method, chi_routes) {
			t.Errorf("route %s is not registered", route.Route)
		}
	}
}

func routeExists(test_route, test_method string, chi_routes chi.Routes) bool {
	found := false

	chi.Walk(chi_routes, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.EqualFold(method, test_method) && strings.EqualFold(route, test_route) {
			found = true
		}
		return nil
	})

	return found
}
