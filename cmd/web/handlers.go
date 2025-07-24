package main

import (
	"html/template"
	"log"
	"net/http"
)

type TemplateData struct {
	IP   string
	Data map[string]any
}

func (app *State) Home(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "home.html", &TemplateData{})
}

func (app *State) Render(w http.ResponseWriter, r *http.Request, template_file string, data *TemplateData) error {
	log.Println("./templates/" + template_file)
	t, err := template.ParseFiles("../../templates/" + template_file)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	err = t.Execute(w, data.Data)
	if err != nil {
		return err
	}

	return nil
}
