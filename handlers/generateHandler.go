package handlers

import (
	"log"
	"net/http"
	"short-url/models"
)

func GenerateHandler(w http.ResponseWriter, r *http.Request) {
	originalURL := r.FormValue("url")
	if originalURL == "" {
		IndexHandler(w, r)
		return
	}
	url := &models.URL{
		OriginalURL: originalURL,
		UserID:      1,
	}
	err := url.InsertURL()
	if err != nil {
		log.Println("Error while Calling InsertURL() : ", err)
		return
	}
	tpl.Execute(w, url)
}
