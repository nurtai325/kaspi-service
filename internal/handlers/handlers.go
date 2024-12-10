package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/nurtai325/kaspi-service/internal/auth"
)

var (
	ErrMethodNotAllowed = errors.New("error method not allowed")
	ErrUnauthenticated  = errors.New("error user is not authenticated")
)

func Register() {
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/tailwind", HandleTailwind)
	http.HandleFunc("/favicon.ico", HandleFav)
	http.HandleFunc("/login", HandleLogin)

	// TODO: uncomment this
	// handlers.WithAuth("/", handlers.HandleClients)
	// handlers.WithAuth("/clients/create", handlers.HandleClientsCreate)
	http.HandleFunc("/", HandleClients)
	http.HandleFunc("/clients/create", HandleClientsCreate)
	http.HandleFunc("/clients/extend", HandleClientExtend)
	http.HandleFunc("/clients/cancel", HandleClientsCancel)
}

func withAuth(path string, next http.HandlerFunc) {
	http.HandleFunc(path, auth.Middleware(next))
}

func addCacheHeader(w http.ResponseWriter, seconds int) {
	w.Header().Add("Cache-Control", fmt.Sprintf("public;max-age=%d", seconds))
}

func newErr(r *http.Request, err error) error {
	return fmt.Errorf("%s %q: %w", r.Method, r.URL.RequestURI(), err)
}
