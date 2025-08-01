package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/blakehulett7/mctestface-webapp/pkg/data"
)

var path_to_templates = "../../templates/"

type TemplateData struct {
	Data  map[string]any
	Error string
	Flash string
	IP    string
	User  data.User
}

func (app *State) Authenticate(r *http.Request, user *data.User, password string) bool {
	valid, err := user.PasswordMatches(password)
	if !valid || err != nil {
		return false
	}

	app.Session.Put(r.Context(), "user", user)
	return true
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
		app.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := app.DB.GetUserByEmail(email)
	if err != nil {
		log.Println(err)
		app.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if !app.Authenticate(r, user, password) {
		app.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	app.Session.RenewToken(r.Context())

	app.Session.Put(r.Context(), "flash", "Log in successful!")
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

func (app *State) Profile(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "profile.html", &TemplateData{})
}

func (app *State) Render(w http.ResponseWriter, r *http.Request, template_file string, td *TemplateData) error {
	t, err := template.ParseFiles(path.Join(path_to_templates, template_file), path.Join(path_to_templates, "base.html"))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	td.IP = app.IpFromContext(r.Context())

	td.Error = app.Session.PopString(r.Context(), "error")
	td.Flash = app.Session.PopString(r.Context(), "flash")

	return t.Execute(w, td)
}
