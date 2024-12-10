package handlers

import (
	"errors"
	"html/template"
	"log"
	"net/http"
)

var (
	ErrExecTemplate  = errors.New("error executing template")
	ErrParseTemplate = errors.New("error parsing templates")
	templates        *template.Template
)

func ParseTemplates() error {
	parsed, err := template.ParseGlob("./templ/*")
	if err != nil {
		return errors.Join(ErrParseTemplate, err)
	}
	templates = parsed
	return nil
}

func Execute(w http.ResponseWriter, r *http.Request, name string, data any) {
	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		err = errors.Join(ErrExecTemplate, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(newErr(r, err))
	}
}
