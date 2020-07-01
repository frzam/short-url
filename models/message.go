package models

import (
	"encoding/json"
	"net/http"
)

// message is used to convert status and message into complex form that
// can be parsed into the response. If the response is correct then status is true else false.
// If the status is false then message contains details about the problem.
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// message is takes statusCode and data and parses into the ResponseWriter with Content-Type
// as "application/json". In short it takes the data and writes the response.
func Respond(w http.ResponseWriter, statusCode int, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
