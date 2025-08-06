package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"testing"
)

func Test_app_Handlers(t *testing.T) {
	tests := []struct {
		Name                    string
		Url                     string
		ExpectedStatusCode      int
		ExpectedUrl             string
		ExpectedFirstStatusCode int
	}{
		{"Home", "/", http.StatusOK, "/", http.StatusOK},
		{"NotFound", "/nirvana", http.StatusNotFound, "/nirvana", http.StatusNotFound},
		{"Profile", "/user/profile", http.StatusOK, "/", http.StatusTemporaryRedirect},
	}

	routes := app.Routes()

	test_server := httptest.NewTLSServer(routes)
	defer test_server.Close()

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	first_statuscode_client := &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for _, test := range tests {
		res, err := test_server.Client().Get(test_server.URL + test.Url)
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != test.ExpectedStatusCode {
			t.Errorf("%s failed: expected status %d but got status %d", test.Name, test.ExpectedStatusCode, res.StatusCode)
		}

		if res.Request.URL.Path != test.ExpectedUrl {
			t.Errorf("%s failed: expected final url %s but got %s\n", test.Name, test.ExpectedUrl, res.Request.URL.Path)
		}

		res_for_first_status_code, _ := first_statuscode_client.Get(test_server.URL + test.Url)
		if res_for_first_status_code.StatusCode != test.ExpectedFirstStatusCode {
			t.Errorf("%s failed: expected first status code %d but got %d", test.Name, test.ExpectedFirstStatusCode, res_for_first_status_code.StatusCode)
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

func Test_app_Login(t *testing.T) {
	tests := []struct {
		Name               string
		PostedData         url.Values
		ExpectedStatusCode int
		ExpectedLocation   string
	}{
		{
			Name: "valid login",
			PostedData: url.Values{
				"email":    {"admin@example.com"},
				"password": {"secret"},
			},
			ExpectedStatusCode: http.StatusSeeOther,
			ExpectedLocation:   "/user/profile",
		},
		{
			Name: "invalid form",
			PostedData: url.Values{
				"email":    {""},
				"password": {""},
			},
			ExpectedStatusCode: http.StatusSeeOther,
			ExpectedLocation:   "/",
		},
		{
			Name: "bad email",
			PostedData: url.Values{
				"email":    {"wrongemail@example.com"},
				"password": {"wrong password"},
			},
			ExpectedStatusCode: http.StatusSeeOther,
			ExpectedLocation:   "/",
		},
		{
			Name: "bad password",
			PostedData: url.Values{
				"email":    {"admin@example.com"},
				"password": {"wrong password"},
			},
			ExpectedStatusCode: http.StatusSeeOther,
			ExpectedLocation:   "/",
		},
	}

	for _, test := range tests {
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(test.PostedData.Encode()))
		req = AddContextAndSessionToRequest(req, app)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()

		handler := http.HandlerFunc(app.Login)
		handler.ServeHTTP(w, req)

		if w.Code != test.ExpectedStatusCode {
			t.Errorf("%s failed: Expected status code %d, but got %d\n", test.Name, test.ExpectedStatusCode, w.Code)
		}

		location, err := w.Result().Location()
		if err != nil {
			t.Errorf("%s failed: Error accessing the location returned from the handler %v\n", test.Name, err)
		}

		if location.String() != test.ExpectedLocation {
			t.Errorf("%s failed: Expected location %s, but got %s\n", test.Name, test.ExpectedLocation, location.String())
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

func Test_app_UploadFiles(t *testing.T) {
	pipeRead, pipeWrite := io.Pipe()
	writer := multipart.NewWriter(pipeWrite)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go simulatePNGUpload("./testdata/img.png", writer, t, wg)

	req := httptest.NewRequest("POST", "/", pipeRead)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	uploadedFiles, err := app.UploadFiles(req, "./testdata/uploads/")
	if err != nil {
		t.Error(err)
	}

	_, err = os.Stat(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].OriginalFileName))
	if os.IsNotExist(err) {
		t.Errorf("expected file to exist: %s", err.Error())
	}

	os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].OriginalFileName))
}

func simulatePNGUpload(fileToUpload string, writer *multipart.Writer, t *testing.T, wg *sync.WaitGroup) {
	defer writer.Close()
	defer wg.Done()

	part, err := writer.CreateFormFile("file", path.Base(fileToUpload))
	if err != nil {
		t.Error(err)
	}

	file, err := os.Open(fileToUpload)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		t.Error("error decoding image")
	}

	err = png.Encode(part, img)
	if err != nil {
		t.Error(err)
	}
}

func GetCtx(r *http.Request) context.Context {
	return context.WithValue(r.Context(), ContextUserKey, "unknown")
}

func AddContextAndSessionToRequest(r *http.Request, app State) *http.Request {
	r = r.WithContext(GetCtx(r))
	ctx, _ := app.Session.Load(r.Context(), r.Header.Get("X-Session"))
	return r.WithContext(ctx)
}
