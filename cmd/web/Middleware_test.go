package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blakehulett7/mctestface-webapp/pkg/data"
)

func Test_application_AddIPToContext(t *testing.T) {
	tests := []struct {
		HeaderName  string
		HeaderValue string
		Addr        string
		EmptyAddr   bool
	}{
		{"", "", "", false},
		{"", "", "", true},
		{"X-Forwarded-For", "192.3.2.1", "", false},
		{"", "", "hello:world", false},
	}

	dummy_handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value(ContextUserKey)
		if val == nil {
			t.Error(ContextUserKey, "not present")
		}

		ip, ok := val.(string)
		if !ok {
			t.Error("not string")
		}
		t.Log(ip)
	})

	for _, test := range tests {
		HandlerToTest := app.AddIpToContext(dummy_handler)

		req := httptest.NewRequest("GET", "http://testing", nil)

		if test.EmptyAddr {
			req.RemoteAddr = ""
		}

		if len(test.HeaderName) > 0 {
			req.Header.Add(test.HeaderName, test.HeaderValue)
		}

		if len(test.Addr) > 0 {
			req.RemoteAddr = test.Addr
		}

		HandlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	}
}

func Test_middleware_AuthMiddleware(t *testing.T) {
	input_handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})

	tests := []struct {
		Name     string
		IsAuthed bool
	}{
		{"logged in", true},
		{"not logged in", false},
	}

	for _, test := range tests {
		handler_to_test := app.AuthMiddleware(input_handler)

		req := httptest.NewRequest("GET", "http://testing", nil)
		req = AddContextAndSessionToRequest(req, app)

		if test.IsAuthed {
			app.Session.Put(req.Context(), "user", data.User{ID: 1})
		}

		w := httptest.NewRecorder()
		handler_to_test.ServeHTTP(w, req)

		if test.IsAuthed && w.Code != http.StatusOK {
			t.Errorf("%s failed: expected status %d but got %d\n", test.Name, http.StatusOK, w.Code)
		}

		if !test.IsAuthed && w.Code != http.StatusTemporaryRedirect {
			t.Errorf("%s failed: expected status %d but got %d\n", test.Name, http.StatusTemporaryRedirect, w.Code)
		}
	}
}

func Test_application_IpFromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), ContextUserKey, "test_value")

	got := app.IpFromContext(ctx)

	if got != "test_value" {
		t.Error("did not get the correct ip out of the context")
	}
}
