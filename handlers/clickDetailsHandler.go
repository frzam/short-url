package handlers

import (
	"log"
	"net/http"
	"short-url/models"
	"short-url/utils"
	"strconv"

	"github.com/gorilla/mux"
)

// Path: Get /api/v1/{shorturl}?skip=0&limit=100
func GetClickDetailsHandler(w http.ResponseWriter, r *http.Request) {
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
		utils.Respond(w, http.StatusBadRequest, utils.Message(false, "No Data is found."))
		return
	}
	resp := utils.Message(true, "Success")
	resp["data"] = res
	utils.Respond(w, http.StatusOK, resp)
}

// Path : Delete /api/v1/{shorturl}
func DeleteClickDetailsHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := mux.Vars(r)["shorturl"]
	if shortURL == "" {
		utils.Respond(w, http.StatusBadRequest, utils.Message(false, "Empty shorturl."))
		return
	}
	cd := &models.ClickDetails{
		ShortURL: shortURL,
	}
	err := cd.DeleteClickDetails()
	if err != nil {
		log.Println("Error while Callin DeleteClickDetails() : ", err)
		utils.Respond(w, http.StatusInternalServerError, utils.Message(false, "Internal Server error."))
		return
	}
	utils.Respond(w, http.StatusAccepted, utils.Message(true, "Deleted!"))
}

// Path : GET /api/v1/{shorturl}/totalcount
func TotalCountHandler(w http.ResponseWriter, r *http.Request) {
	shorturl := mux.Vars(r)["shorturl"]
	if shorturl == "" {
		utils.Respond(w, http.StatusBadRequest, utils.Message(false, "Invalid shorturl."))
		return
	}
	cd := &models.ClickDetails{
		ShortURL: shorturl,
	}
	count, err := cd.GetTotalClicksCount()
	if err != nil {
		utils.Respond(w, http.StatusInternalServerError, utils.Message(false, "Internal Server Error"))
		return
	}
	resp := utils.Message(true, "Success")
	resp["total_count"] = count
	utils.Respond(w, http.StatusOK, resp)
}

func getSkipAndLimit(r *http.Request) (int64, int64) {
	skip := r.URL.Query().Get("skip")
	s, err := strconv.Atoi(skip)
	if err != nil {
		log.Println("Error in getSkipAndLimit skip : ", skip)
		s = 0
	}
	limit := r.URL.Query().Get("limit")
	l, err := strconv.Atoi(limit)
	if err != nil {
		log.Println("Error in getSkipAndLimit limit : ", limit)
		l = 100
	}
	if l > 100 {
		l = 100
	}

	return int64(s), int64(l)
}
