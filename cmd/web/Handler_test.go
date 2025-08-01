package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

func Test_app_Home(t *testing.T) {
	tests := []struct {
		Name         string
		PutInSession string
		ExpectedHTML string
	}{
		{"first visit", "", "<small>From session:"},
		{"second visit", "kyrie eleison", "<small>From session: kyrie eleison"},
	}

	for _, test := range tests {
		req, _ := http.NewRequest("GET", "/", nil)
		req = AddContextAndSessionToRequest(req, app)
		app.Session.Destroy(req.Context())

		if test.PutInSession != "" {
			app.Session.Put(req.Context(), "test", test.PutInSession)
		}

		w := httptest.NewRecorder()

		handler := http.HandlerFunc(app.Home)
		handler.ServeHTTP(w, req)

		if w.Result().StatusCode != http.StatusOK {
			t.Error("Expected status 200 but did not get it")
		}

		body, _ := io.ReadAll(w.Body)
		if !strings.Contains(string(body), test.ExpectedHTML) {
			t.Errorf("%s failed: not getting back expected %s", test.Name, test.ExpectedHTML)
		}
	}
}

func Test_app_RenderBadTemplate(t *testing.T) {
	path_to_templates = "./testdata/"

	req, _ := http.NewRequest("GET", "/", nil)
	req = AddContextAndSessionToRequest(req, app)

	w := httptest.NewRecorder()

	err := app.Render(w, req, "bad.html", &TemplateData{})
	if err == nil {
		t.Error("Expected bad template error, but did not get it")
	}

	path_to_templates = "../../templates/"
}

func GetCtx(r *http.Request) context.Context {
	return context.WithValue(r.Context(), ContextUserKey, "unknown")
}

func AddContextAndSessionToRequest(r *http.Request, app State) *http.Request {
	r = r.WithContext(GetCtx(r))
	ctx, _ := app.Session.Load(r.Context(), r.Header.Get("X-Session"))
	return r.WithContext(ctx)
}
