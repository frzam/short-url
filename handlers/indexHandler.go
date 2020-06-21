package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"short-url/models"
)

// Creating a template instance so that we can execute our data into it.
var tpl = template.Must(template.ParseFiles("index.html"))

// IndexHandler is used to handle "/" path (HOME).
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:]
	fmt.Println("shortURL : ", shortURL)
	url := &models.URL{
		ShortURL: shortURL,
	}
	var originalURL string
	var err error
	// Calling the Redis to get the cache value
	originalURL, _ = url.GetCacheURL()
	// If Not found then only call the MongoDB.
	if originalURL == "" {
		originalURL, err = url.GetURL()
		if err != nil {
			log.Println("err : ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	fmt.Println("originalURL : ", originalURL)
	http.Redirect(w, r, originalURL, http.StatusFound)
}
