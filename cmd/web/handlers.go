package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
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

type UploadedFile struct {
	OriginalFileName string
	FileSize         int64
}

func (app *State) UploadFiles(r *http.Request, uploadDir string) ([]*UploadedFile, error) {
	var uploadedFiles []*UploadedFile

	err := r.ParseMultipartForm(int64(1024 * 1024 * 5))
	if err != nil {
		return nil, fmt.Errorf("File must be less than 5MB")
	}

	for _, fileHeaders := range r.MultipartForm.File {
		for _, header := range fileHeaders {
			uploadedFiles, err = func(uploadedFiles []*UploadedFile) ([]*UploadedFile, error) {
				var uploadedFile UploadedFile
				infile, err := header.Open()
				if err != nil {
					return nil, err
				}
				defer infile.Close()

				uploadedFile.OriginalFileName = header.Filename

				var outfile *os.File
				defer outfile.Close()

				outfile, err = os.Create(filepath.Join(uploadDir, uploadedFile.OriginalFileName))
				if err != nil {
					return nil, err
				}

				fileSize, err := io.Copy(outfile, infile)
				if err != nil {
					return nil, err
				}

				uploadedFile.FileSize = fileSize
				uploadedFiles = append(uploadedFiles, &uploadedFile)

				return uploadedFiles, nil

			}(uploadedFiles)

			if err != nil {
				return uploadedFiles, err
			}
		}
	}

	return uploadedFiles, nil
}

func (app *State) UploadProfilePicture(w http.ResponseWriter, r *http.Request) {
	files, err := app.UploadFiles(r, "../../static/img/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := app.Session.Get(r.Context(), "user").(data.User)

	image := data.UserImage{
		UserID:   user.ID,
		FileName: files[0].OriginalFileName,
	}

	_, err = app.DB.InsertUserImage(image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedUser, err := app.DB.GetUser(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app.Session.Put(r.Context(), "user", updatedUser)
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}
