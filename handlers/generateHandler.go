package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"short-url/models"
)

type captchaReq struct {
	secret string
	token  string
}

type captchaResp struct {
	success      bool
	challenge_ts string
	hostname     string
}

func GenerateHandler(w http.ResponseWriter, r *http.Request) {
	originalURL := r.FormValue("url")
	if originalURL == "" {
		IndexHandler(w, r)
		return
	}
	_ = verifyCaptcha(r)
	host := os.Getenv("host")
	url := &models.URL{
		OriginalURL: originalURL,
		UserID:      1,
	}
	err := url.InsertURL()
	if err != nil {
		log.Println("Error while Calling InsertURL() : ", err)
		return
	}
	err = url.SetCacheURL()
	if err != nil {
		log.Println("Error while Calling SetCacheURL() : ", err)
		return
	}
	url.ShortURL = host + url.ShortURL
	tpl.Execute(w, url)
}

func verifyCaptcha(r *http.Request) bool {
	endPoint := "https://www.google.com/recaptcha/api/siteverify"
	capReq := captchaReq{
		secret: os.Getenv("privateToken"),
		token:  r.FormValue("g-recaptcha-response"),
	}
	req, err := json.Marshal(capReq)
	if err != nil {
		log.Println("Error in verifyCaptcha while marshalling : ", err)
	}
	resp, err := http.Post(endPoint, "application/json", bytes.NewBuffer(req))
	if err != nil {
		log.Println("Error while Calling Captcha Service : ", err)
		return false
	}
	var capResp captchaResp
	err = json.NewDecoder(resp.Body).Decode(&capResp)
	if err != nil {
		log.Println("Error while Decode Captcha Response : ", err)
	}
	defer resp.Body.Close()
	return true
}
