package main

import (
	"log"
	"net/http"
	"os"
	"short-url/handlers"
)

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
	s := handlers.NewServer()
	s.Routes()

	// Starting Server.
	log.Println("Starting Server at : ", port)

	// If the environment is PROD then we are serving all the content on 443 with SSL
	// but for our DEV environment we are still using 8080 without any ssl.
	// In PROD if the Request comes on Port 80 when we are redirecting onto 443.
	if env == "PROD" {
		go http.ListenAndServe(":80", http.HandlerFunc(redirectTLS))

		log.Fatal(http.ListenAndServeTLS(":"+port, fullchain, privkey, s.Router))

	} else {
		log.Fatal(http.ListenAndServe(":"+port, s.Router))
	}

	// Closing th Mongo and Redis Client when the main() exits.
	//defer models.GetMongoClient().Disconnect(context.TODO())
	//defer models.GetRedisClient().Close()

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
