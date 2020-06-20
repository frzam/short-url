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
	originalURL, err := url.GetURL()
	if err != nil {
		log.Println("err : ", err)
	}
	fmt.Println("originalURL : ", originalURL)
	http.Redirect(w, r, originalURL, http.StatusFound)
}
