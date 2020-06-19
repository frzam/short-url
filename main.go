package main

import (
	"log"
	"net/http"
	"os"
	"short-url/handlers"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Create a new router instance
	r := mux.NewRouter()

	// Creating a file server object.
	fs := http.FileServer(http.Dir("assets"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	// Serving IndexHandler
	r.HandleFunc("/", handlers.IndexHandler)

	// Starting Server.
	log.Println("Starting Server at : ", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
