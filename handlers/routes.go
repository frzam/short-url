package handlers

import (
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	DB     *mongo.Client
	Cache  *redis.Client
	Router *mux.Router
}

func NewServer() *Server {
	return &Server{
		Router: mux.NewRouter(),
	}
}

func (s *Server) Routes() {
	// Adding a middel middlewas.routere to Log each request url and IP.
	s.Router.Use(LoggingMiddleware)
	// Creating a file server object.
	fs := http.FileServer(http.Dir("assets"))
	s.Router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	// Generate shorturl handles.Router.
	s.Router.HandleFunc("/generate", GenerateHandler)
	// Serving IndexHandler
	s.Router.HandleFunc("/", IndexHandler)

	// API for Click Details:
	s.Router.HandleFunc("/api/v1/{shorturl}/{days}", TotalDetailsNdaysHandler).Methods(http.MethodGet)
	s.Router.HandleFunc("/api/v1/{shorturl}/country/{country}", TotalDetailsByCountryHandler).Methods(http.MethodGet)
	s.Router.HandleFunc("/api/v1/{shorturl}/city/{city}", TotalDetailsByCityHandler).Methods(http.MethodGet)
	s.Router.HandleFunc("/api/v1/{shorturl}/ip/{ip}", TotalDetailsByIP).Methods(http.MethodGet)
	s.Router.HandleFunc("/api/v1/{shorturl}/totalcount", TotalCountHandler).Methods(http.MethodGet)
	s.Router.HandleFunc("/api/v1/{shorturl}/totalcount/{days}", TotalCountNdaysHandler).Methods(http.MethodGet)
	s.Router.HandleFunc("/api/v1/{shorturl}/ip/{ip}/totalcount", ClickCountsByIP).Methods(http.MethodGet)
	s.Router.HandleFunc("/api/v1/{shorturl}", DeleteClickDetailsHandler).Methods(http.MethodDelete)
	s.Router.HandleFunc("/api/v1/{shorturl}", s.getClickDetailsHandler()).Methods(http.MethodGet)

	// Get the original url from shorturl
	s.Router.HandleFunc("/{[a-zA-Z0-9_.-]*}", Redirect)
}
