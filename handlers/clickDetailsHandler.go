package handlers

import (
	"fmt"
	"log"
	"net/http"
	"short-url/models"
	"strconv"

	"github.com/gorilla/mux"
)

// GetClickDetailsHandler is used to get the click details for one particular shorturl.
// It takes two optional params skip and limit.
// Path: GET /api/v1/{shorturl}?skip=0&limit=100
func (s *Server) getClickDetailsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Inside the getClickDetailsHandler")
		shortURL := mux.Vars(r)["shorturl"]
		skip, limit := getSkipAndLimit(r)
		cd := &models.ClickDetails{
			ShortURL: shortURL,
		}
		res, err := cd.GetTotalClicksDetails(skip, limit)
		if err != nil {
			log.Println("Error while calling GetClickDetails() : ", err)
			return
		}
		if res == nil {
			models.Respond(w, http.StatusBadRequest, models.Message(false, "No Data is found."))
			return
		}
		resp := models.Message(true, "Success")
		resp["data"] = res
		models.Respond(w, http.StatusOK, resp)

	}
}

// DeleteClickDetailsHandler is used to delete all the details of a shorturl.
// Path : DELETE /api/v1/{shorturl}
func DeleteClickDetailsHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := mux.Vars(r)["shorturl"]
	if shortURL == "" {
		models.Respond(w, http.StatusBadRequest, models.Message(false, "Empty shorturl."))
		return
	}
	cd := &models.ClickDetails{
		ShortURL: shortURL,
	}
	err := cd.DeleteClickDetails()
	if err != nil {
		log.Println("Error while Callin DeleteClickDetails() : ", err)
		models.Respond(w, http.StatusInternalServerError, models.Message(false, "Internal Server error."))
		return
	}
	models.Respond(w, http.StatusAccepted, models.Message(true, "Deleted!"))
}

// TotalCountHandler returns total count of hits for one particular shorturl.
// It taskes the shorturl from the url param.
// Path : GET /api/v1/{shorturl}/totalcount
func TotalCountHandler(w http.ResponseWriter, r *http.Request) {
	shorturl := mux.Vars(r)["shorturl"]
	if shorturl == "" {
		models.Respond(w, http.StatusBadRequest, models.Message(false, "Invalid shorturl."))
		return
	}
	cd := &models.ClickDetails{
		ShortURL: shorturl,
	}
	count, err := cd.GetTotalClicksCount()
	if err != nil {
		models.Respond(w, http.StatusInternalServerError, models.Message(false, "Internal Server Error"))
		return
	}
	resp := models.Message(true, "Success")
	resp["total_count"] = count
	models.Respond(w, http.StatusOK, resp)
}

// TotalCountNdaysHandler is used to get the total hit count for past n days.
// It takes shorturl from the url param and days(days >= 1).
// Path : GET /api/v1/{shorturl}/totalcount/{days}
func TotalCountNdaysHandler(w http.ResponseWriter, r *http.Request) {
	shorturl := mux.Vars(r)["shorturl"]
	days := mux.Vars(r)["days"]
	if shorturl == "" || days == "" {
		models.Respond(w, http.StatusBadRequest, models.Message(false, "Invalid Path params"))
	}
	d, err := strconv.Atoi(days)
	if err != nil {
		log.Println("Error in TotalCountNdaysHandler : ", err)
		d = 1
	}
	cd := &models.ClickDetails{
		ShortURL: shorturl,
	}
	count, err := cd.GetNdayClicksCount(d)
	if err != nil {
		models.Respond(w, http.StatusInternalServerError, models.Message(false, "Internal Server Error."))
		return
	}
	resp := models.Message(true, "Success")
	resp["count"] = count
	models.Respond(w, http.StatusOK, resp)
}

// TotalDetailsNdaysHandler is used to return the total data response for past n days.
// Days should always be more than 0.
// It takes the skip and limit query param.. Max(limit) <= 100
// Path : GET /api/v1/{shorturl}/{days}
func TotalDetailsNdaysHandler(w http.ResponseWriter, r *http.Request) {
	shorturl := mux.Vars(r)["shorturl"]
	days := mux.Vars(r)["days"]
	if days == "" || shorturl == "" {
		models.Respond(w, http.StatusBadRequest, models.Message(false, "Invalid Path params."))
		return
	}
	skip, limit := getSkipAndLimit(r)
	cd := &models.ClickDetails{
		ShortURL: shorturl,
	}
	d, err := strconv.Atoi(days)
	if err != nil {
		log.Println("Error in TotalDetailsNdaysHandler : ", err)
		d = 1
	}
	data, err := cd.GetNdayClicksDetails(d, skip, limit)
	if err != nil {
		models.Respond(w, http.StatusInternalServerError, models.Message(false, "Internal Server Error"))
		return
	}
	resp := models.Message(true, "Success")
	resp["data"] = data
	models.Respond(w, http.StatusOK, resp)

}

// TotalDetailsByCountryHandler is used to get the click details per country.
// It accepts skip and limit as query params. limit <= 100.
// Path : GET /api/v1/{shorturl}/country/{country}
func TotalDetailsByCountryHandler(w http.ResponseWriter, r *http.Request) {
	shorturl := mux.Vars(r)["shorturl"]
	country := mux.Vars(r)["country"]
	if shorturl == "" || country == "" {
		models.Respond(w, http.StatusBadRequest, models.Message(false, "Invalid Path params"))
		return
	}
	cd := &models.ClickDetails{
		ShortURL: shorturl,
	}
	skip, limit := getSkipAndLimit(r)

	data, err := cd.GetClicksDetailsByCountry(country, skip, limit)
	if err != nil {
		models.Respond(w, http.StatusInternalServerError, models.Message(false, "Internal Server Error"))
		return
	}
	resp := models.Message(true, "Success")
	resp["data"] = data
	models.Respond(w, http.StatusOK, resp)
}

// TotalDetailsByCityHandler is used to get the click details per city.
// It takes two query param skip and limit where limit <= 100.
// Path : GET /api/v1/{shorturl}/city/{city}
func TotalDetailsByCityHandler(w http.ResponseWriter, r *http.Request) {
	shorturl := mux.Vars(r)["shorturl"]
	city := mux.Vars(r)["city"]
	if shorturl == "" || city == "" {
		models.Respond(w, http.StatusBadRequest, models.Message(false, "Invalid Path Params."))
		return
	}
	skip, limit := getSkipAndLimit(r)
	cd := &models.ClickDetails{
		ShortURL: shorturl,
	}

	data, err := cd.GetClicksDetailsByCity(city, skip, limit)
	if err != nil {
		models.Respond(w, http.StatusInternalServerError, models.Message(false, "Internal Server Error"))
		return
	}
	resp := models.Message(true, "Success")
	resp["data"] = data
	models.Respond(w, http.StatusOK, resp)
}

// TotalDetailsByIP returns the click details by particular IP.
// It uses two query params skip and limit where limit <= 100.
// Path : GET /api/v1/{shorturl}/ip/{ip}
func TotalDetailsByIP(w http.ResponseWriter, r *http.Request) {
	shorturl := mux.Vars(r)["shorturl"]
	ip := mux.Vars(r)["ip"]
	if shorturl == "" || ip == "" {
		models.Respond(w, http.StatusBadRequest, models.Message(false, "Invalid Path params"))
		return
	}
	skip, limit := getSkipAndLimit(r)
	cd := &models.ClickDetails{
		ShortURL: shorturl,
	}
	data, err := cd.GetClicksDetailsByIP(ip, skip, limit)
	if err != nil {
		models.Respond(w, http.StatusInternalServerError, models.Message(false, "Internal Server Error"))
		return
	}
	resp := models.Message(true, "Sucess")
	resp["data"] = data
	models.Respond(w, http.StatusOK, resp)
}

// ClickCountsByIP returns the total click count from one IP address.
// Path : GET /api/v1/{shorturl}/ip/{ip}/totalcount
func ClickCountsByIP(w http.ResponseWriter, r *http.Request) {
	shorturl := mux.Vars(r)["shorturl"]
	ip := mux.Vars(r)["ip"]
	if shorturl == "" || ip == "" {
		models.Respond(w, http.StatusBadRequest, models.Message(false, "Invalid Path Param"))
		return
	}
	cd := &models.ClickDetails{
		ShortURL: shorturl,
	}
	count, err := cd.GetClicksCountByIP(ip)
	if err != nil {
		models.Respond(w, http.StatusInternalServerError, models.Message(false, "Internal Server Error."))
		return
	}
	resp := models.Message(true, "Success")
	resp["total_count"] = count
	models.Respond(w, http.StatusOK, resp)
}

// getSkipAndLimit grabs the query params skip and limit,
// and check whether they are correct or not.
// If skip is not present in the param then skip = 0.
// If limit is not present or the limit is more than 100 then limit = 100.
func getSkipAndLimit(r *http.Request) (int64, int64) {
	skip := r.URL.Query().Get("skip")
	s, err := strconv.Atoi(skip)
	if err != nil {
		s = 0
	}
	limit := r.URL.Query().Get("limit")
	l, err := strconv.Atoi(limit)
	if err != nil {
		l = 100
	}
	if l > 100 {
		l = 100
	}
	return int64(s), int64(l)
}
