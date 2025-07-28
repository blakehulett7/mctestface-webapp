package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

var app State

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

func Test_application_IpFromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), ContextUserKey, "test_value")

	got := app.IpFromContext(ctx)

	if got != "test_value" {
		t.Error("did not get the correct ip out of the context")
	}
}
