package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/nurtai325/kaspi-service/internal/auth"
)

const (
	loginTempl = "login.html"
	errorQuery = "error"
)

type loginTemplData struct {
	Error string
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		loginErr := r.URL.Query().Get(errorQuery)
		Execute(w, r, loginTempl, loginTemplData{Error: loginErr})
		return
	}
	r.ParseForm()
	phone := r.FormValue("phone")
	password := r.FormValue("password")
	err := auth.Login(w, r, phone, password)
	if err != nil {
		if !errors.Is(err, auth.ErrLogin) {
			err = fmt.Errorf("user %q password %q: %w", phone, password, err)
			log.Println(newErr(r, err))
		}
		return
	}
}
