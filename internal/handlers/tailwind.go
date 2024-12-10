package handlers

import (
	"errors"
	"log"
	"net/http"
	"os"
)

var (
	ErrTailWindFile = errors.New("error serving tailwind file")
)

func HandleTailwind(w http.ResponseWriter, r *http.Request) {
	contents, err := os.ReadFile("./assets/tailwind.js")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
	addCacheHeader(w, 31536000)
	w.Write(contents)
}

func HandleFav(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/assets/favicon.ico", http.StatusMovedPermanently)
}
