package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"short-url/handlers"
	"time"

	_ "short-url/models"

	"github.com/gorilla/mux"
)

func main() {
	currentTime := time.Now().Format("20060102150405")
	fmt.Println("currentTime : ", "google.com"+currentTime)

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
