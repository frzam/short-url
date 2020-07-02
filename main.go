package main

import (
	"log"
	"net/http"
	"os"
	"short-url/handlers"

	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
}

func main() {
	// Get environment variables relating port, env i.e DEV or PROD, fullchain and privKey for
	// SSL Certificate.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	env := os.Getenv("env")
	fullchain := os.Getenv("fullchain")
	privkey := os.Getenv("privkey")

	// Open a file "info.log", all the application logs are saved inside it,
	// It closes when the main() exits.
	file, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error while opening info.log : ", err)
	}
	defer file.Close()
	log.SetOutput(file)

	// Create a new router instance
	s := newServer()
	s.routes()

	// Starting Server.
	log.Println("Starting Server at : ", port)

	// If the environment is PROD then we are serving all the content on 443 with SSL
	// but for our DEV environment we are still using 8080 without any ssl.
	// In PROD if the Request comes on Port 80 when we are redirecting onto 443.
	if env == "PROD" {
		go http.ListenAndServe(":80", http.HandlerFunc(redirectTLS))

		log.Fatal(http.ListenAndServeTLS(":"+port, fullchain, privkey, s.router))

	} else {
		log.Fatal(http.ListenAndServe(":"+port, s.router))
	}

	// Closing th Mongo and Redis Client when the main() exits.
	//defer models.GetMongoClient().Disconnect(context.TODO())
	//defer models.GetRedisClient().Close()

}

func newServer() *server {
	return &server{
		router: mux.NewRouter(),
	}
}

func (s *server) routes() {
	// Adding a middel middlewas.routere to Log each request url and IP.
	s.router.Use(handlers.LoggingMiddleware)
	// Creating a file server object.
	fs := http.FileServer(http.Dir("assets"))
	s.router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	// Generate shorturl handles.router.
	s.router.HandleFunc("/generate", handlers.GenerateHandler)
	// Serving IndexHandler
	s.router.HandleFunc("/", handlers.IndexHandler)

	// API for Click Details:
	s.router.HandleFunc("/api/v1/{shorturl}/{days}", handlers.TotalDetailsNdaysHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/api/v1/{shorturl}/country/{country}", handlers.TotalDetailsByCountryHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/api/v1/{shorturl}/city/{city}", handlers.TotalDetailsByCityHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/api/v1/{shorturl}/ip/{ip}", handlers.TotalDetailsByIP).Methods(http.MethodGet)
	s.router.HandleFunc("/api/v1/{shorturl}/totalcount", handlers.TotalCountHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/api/v1/{shorturl}/totalcount/{days}", handlers.TotalCountNdaysHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/api/v1/{shorturl}/ip/{ip}/totalcount", handlers.ClickCountsByIP).Methods(http.MethodGet)
	s.router.HandleFunc("/api/v1/{shorturl}", handlers.DeleteClickDetailsHandler).Methods(http.MethodDelete)
	s.router.HandleFunc("/api/v1/{shorturl}", handlers.GetClickDetailsHandler()).Methods(http.MethodGet)

	// Get the original url from shorturl
	s.router.HandleFunc("/{[a-zA-Z0-9_.-]*}", handlers.Redirect)
}

// redirectTLS redirects HTTP Request to HTTPS. It is called when the port
// from the client is 80, then redirectTLS serves the HTTPS version with the same URL
// and query params. In short it takes "http://..." and converts it to "https://.."
func redirectTLS(w http.ResponseWriter, r *http.Request) {
	target := "https://" + r.Host + r.URL.Path
	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}
	http.Redirect(w, r, target, http.StatusTemporaryRedirect)
}
