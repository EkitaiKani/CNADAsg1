package handlers

import (
	"net/http"
	"html/template"

	"CNADASG1/templates"
)

type HomeHandler struct {
	Templates *template.Template
}

func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {

	err := templates.Templates.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *HomeHandler) Login(w http.ResponseWriter, r *http.Request) {

	err := templates.Templates.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *HomeHandler) Register(w http.ResponseWriter, r *http.Request) {

	err := templates.Templates.ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

