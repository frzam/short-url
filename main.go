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
	log.Fatal(http.ListenAndServe(":"+port, r))
}
