package main

import (
	"html/template"
	"net/http"
	"path"
)

var path_to_templates = "../../templates/"

type TemplateData struct {
	IP   string
	Data map[string]any
}

func (app *State) Home(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "home.html", &TemplateData{})
}

func (app *State) Render(w http.ResponseWriter, r *http.Request, template_file string, data *TemplateData) error {
	t, err := template.ParseFiles(path.Join(path_to_templates, template_file))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	data.IP = app.IpFromContext(r.Context())

	err = t.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}
