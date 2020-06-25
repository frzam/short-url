package main

import (
	"log"
	"net/http"
	"os"
	"short-url/handlers"

	_ "short-url/models"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
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

	// API:
	r.HandleFunc("/api/v1/{shorturl}", handlers.GetClickDetailsHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/{shorturl}", handlers.DeleteClickDetailsHandler).Methods(http.MethodDelete)
	r.HandleFunc("/api/v1/{shorturl}/totalcount", handlers.TotalCountHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/{shorturl}/totalcount/{days}", handlers.TotalCountNdaysHandler).Methods(http.MethodGet)
	// Get the original url from shorturl
	r.HandleFunc("/{[a-zA-Z0-9_.-]*}", handlers.Redirect)
	// Starting Server.
	log.Println("Starting Server at : ", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
