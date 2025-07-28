package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
)

var path_to_templates = "../../templates/"

type TemplateData struct {
	IP   string
	Data map[string]any
}

func (app *State) Home(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]any)

	if app.Session.Exists(r.Context(), "test") {
		msg := app.Session.GetString(r.Context(), "test")
		data["test"] = msg
		app.Render(w, r, "home.html", &TemplateData{Data: data})
		return
	}

	app.Session.Put(r.Context(), "test", "hit at "+time.Now().UTC().String())
	app.Render(w, r, "home.html", &TemplateData{Data: data})
}

func (app *State) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	form := NewForm(r.PostForm)
	form.Required("email", "password")

	if !form.Valid() {
		fmt.Fprint(w, "failed validation")
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	log.Println(email, password)

	fmt.Fprint(w, email)
}

func (app *State) Render(w http.ResponseWriter, r *http.Request, template_file string, data *TemplateData) error {
	t, err := template.ParseFiles(path.Join(path_to_templates, template_file), path.Join(path_to_templates, "base.html"))
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
