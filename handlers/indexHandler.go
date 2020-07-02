package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"short-url/models"
	"strings"
)

// Creating a template instance so that we can execute our data into it.
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("assets/html/index.html"))
}

// IndexHandler is used to handle "/" path (HOME).
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

// Redirect will the shorturl part from the url then create a url instance and check the RedisCache
// If data is present then it will update the ClickDetails. After the updation it will redirect to
// Original url.
func Redirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := r.URL.Path[1:]
		url := &models.URL{
			ShortURL: shortURL,
		}
		var err error
		// Calling the Redis to get the cache value
		_ = url.Get()
		fmt.Println("originalURL from Cache : ", url.OriginalURL)
		// If Not found then only call the MongoDB.
		if url.OriginalURL == "" {
			url.OriginalURL, err = url.GetURL()
			if err != nil {
				log.Println("err : ", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		// Get IP Address of the client.
		ip := getIPAddress(r)
		// Call AddClickDetails to Save the click details data in mongoDB.
		_ = url.AddClickDetails(ip)
		fmt.Println("url.OriginalURL :", url.OriginalURL)
		http.Redirect(w, r, url.OriginalURL, http.StatusFound)
	}

}

// getIPAddress from the request.
func getIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-for")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	if strings.Contains(ip, ":") {
		ip = ip[:strings.Index(ip, ":")]
	}
	return ip
}
