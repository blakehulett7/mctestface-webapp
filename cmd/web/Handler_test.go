package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_app_Handlers(t *testing.T) {
	tests := []struct {
		Name               string
		Url                string
		ExpectedStatusCode int
	}{
		{"Home", "/", http.StatusOK},
		{"NotFound", "/nirvana", http.StatusNotFound},
	}

	var app State
	routes := app.Routes()

	test_server := httptest.NewTLSServer(routes)
	defer test_server.Close()

	for _, test := range tests {
		res, err := test_server.Client().Get(test_server.URL + test.Url)
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != test.ExpectedStatusCode {
			t.Errorf("%s failed: expected status %d but got status %d", test.Name, test.ExpectedStatusCode, res.StatusCode)
		}
	}
}
