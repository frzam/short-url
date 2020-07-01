package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"short-url/models"
)

// captchaReq is used for the captchaReq param for the google captcha api.
type captchaReq struct {
	secret string
	token  string
}

// captchaResp is used to get the response from the google captcha api.
type captchaResp struct {
	success      bool
	challenge_ts string
	hostname     string
}

// GenerateHandler is used to get the original url param from the request.
// It will verify captcha. It inserts the url and set the cache url.
// And execute the url object on the w.
func GenerateHandler(w http.ResponseWriter, r *http.Request) {
	originalURL := r.FormValue("url")
	if originalURL == "" {
		IndexHandler(w, r)
		return
	}
	_ = verifyCaptcha(r)
	url := &models.URL{
		OriginalURL: originalURL,
		UserID:      1,
	}
	err := url.InsertURL()
	if err != nil {
		log.Println("Error while Calling InsertURL() : ", err)
		return
	}
	// It sets the newly generated url in cache for fast url retrival.
	err = url.SetCacheURL()
	if err != nil {
		log.Println("Error while Calling SetCacheURL() : ", err)
		return
	}

	host := os.Getenv("host")
	url.ShortURL = host + url.ShortURL
	tpl.Execute(w, url)
}

// verifyCaptcha takes the token from the front end and calls the captcha api.
func verifyCaptcha(r *http.Request) bool {
	endPoint := "https://www.google.com/recaptcha/api/siteverify"
	capReq := captchaReq{
		secret: os.Getenv("private_token"),
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
