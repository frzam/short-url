package handlers

import (
	"log"
	"net/http"
)

// LoggingMiddleware is used to log the request URL and IP of the client.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipAddress := r.Header.Get("X-Real-IP")
		if ipAddress == "" {
			ipAddress = r.Header.Get("X-Forwarded-For")
		}
		if ipAddress == "" {
			ipAddress = r.RemoteAddr
		}
		log.Println("URL : ", r.URL.Path)
		log.Println("IP : ", ipAddress)
		next.ServeHTTP(w, r)
	})
}
