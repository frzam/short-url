package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"short-url/handlers"

	"short-url/models"

	"github.com/gorilla/mux"
)

// TO DO:
// Total Count Bug.		--> Done.
// SSL and Deploy
// Integrate catcha
// Write Comments and Deploy.

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	env := os.Getenv("env")
	fullchain := os.Getenv("fullchain")
	privkey := os.Getenv("privkey")
	defer models.GetMongoClient().Disconnect(context.TODO())
	defer models.GetRedisClient().Close()

	// Create a new router instance
	r := mux.NewRouter()
	r.Use(handlers.LoggingMiddleware)
	// Creating a file server object.
	fs := http.FileServer(http.Dir("assets"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	// Generate shorturl handler.
	r.HandleFunc("/generate", handlers.GenerateHandler)
	// Serving IndexHandler
	r.HandleFunc("/", handlers.IndexHandler)

	// API for Click Details:
	r.HandleFunc("/api/v1/{shorturl}/{days}", handlers.TotalDetailsNdaysHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/{shorturl}/country/{country}", handlers.TotalDetailsByCountryHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/{shorturl}/city/{city}", handlers.TotalDetailsByCityHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/{shorturl}/ip/{ip}", handlers.TotalDetailsByIP).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/{shorturl}/totalcount", handlers.TotalCountHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/{shorturl}/totalcount/{days}", handlers.TotalCountNdaysHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/{shorturl}/ip/{ip}/totalcount", handlers.ClickCountsByIP).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/{shorturl}", handlers.DeleteClickDetailsHandler).Methods(http.MethodDelete)
	r.HandleFunc("/api/v1/{shorturl}", handlers.GetClickDetailsHandler).Methods(http.MethodGet)

	// Get the original url from shorturl
	r.HandleFunc("/{[a-zA-Z0-9_.-]*}", handlers.Redirect)
	// Starting Server.
	log.Println("Starting Server at : ", port)

	if env == "PROD" {
		go http.ListenAndServe(":80", http.HandlerFunc(redirectTLS))

		log.Fatal(http.ListenAndServeTLS(":"+port, fullchain, privkey, r))

	} else {
		log.Fatal(http.ListenAndServe(":"+port, r))
	}
}

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	target := "https://" + r.Host + r.URL.Path
	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}
	http.Redirect(w, r, target, http.StatusTemporaryRedirect)
}
