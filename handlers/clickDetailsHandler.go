package handlers

import (
	"net/http"
	"short-url/models"
	"short-url/utils"
	"log"
	"github.com/gorilla/mux"
)

func GetClickDetailsHandler(w http.ResponseWriter, r *http.Request)  {
	shortURL := mux.Vars(r)["shorturl"]
	cd := &models.ClickDetails{
			ShortURL: shortURL,
	}
	res, err := cd.GetClickDetails()
	if err != nil{
		log.Println("Error while calling GetClickDetails() : ", err)
		return
	}
	if res == nil{
		utils.Respond(w, http.StatusBadRequest, utils.Message(false, "No Data is found."))
		return
	}
	resp := utils.Message(true, "Success")
	resp["data"] = res
	utils.Respond(w,http.StatusOK,resp)
}


func DeleteClickDetailsHandler(w http.ResponseWriter, r *http.Request)  {
	shortURL := mux.Vars(r)["shorturl"]
	if shortURL == ""{
		utils.Respond(w, http.StatusBadRequest, utils.Message(false,"Empty shorturl."))
		return
	}
	cd := &models.ClickDetails{
		ShortURL: shortURL,
	}
	err := cd.DeleteClickDetails()
	if err != nil{
		log.Println("Error while Callin DeleteClickDetails() : ",err)
		utils.Respond(w, http.StatusInternalServerError, utils.Message(false,"Internal Server error."))
		return
	}
	utils.Respond(w,http.StatusAccepted,utils.Message(true, "Deleted!"))
}