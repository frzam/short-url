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
	resp := utils.Message(true, "Success")
	resp["data"] = res
	utils.Respond(w,http.StatusOK,resp)
}