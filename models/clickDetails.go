package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type ClickDetails struct {
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	CurrentTime time.Time `json:"current_time"`
	IPInfo      IPInfo    `json:"ip_info"`
}

type IPInfo struct {
	IP            string   `json:"ip"`
	Type          string   `json:"type"`
	ContinentCode string   `json:"continent_code"`
	ContinentName string   `json:"continent_name"`
	CountryCode   string   `json:"country_code"`
	CountryName   string   `json:"country_name"`
	RegionCode    string   `json:"region_code"`
	RegionName    string   `json:"region_name"`
	City          string   `json:"city"`
	Zip           string   `json:"zip"`
	Latitude      float64  `json:"latitude"`
	Longitude     float64  `json:"longitude"`
	Location      Location `json:"location"`
}
type Languages struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Native string `json:"native"`
}
type Location struct {
	GeonameID               int         `json:"geoname_id"`
	Capital                 string      `json:"capital"`
	Languages               []Languages `json:"languages"`
	CountryFlag             string      `json:"country_flag"`
	CountryFlagEmoji        string      `json:"country_flag_emoji"`
	CountryFlagEmojiUnicode string      `json:"country_flag_emoji_unicode"`
	CallingCode             string      `json:"calling_code"`
	IsEu                    bool        `json:"is_eu"`
}

func GetIPInfo(ip string) (IPInfo, error) {
	apiKey := os.Getenv("ipstack_apiKey")
	if apiKey == "" {
		log.Println("apiKey is empty")
	}
	endPoint := fmt.Sprintf("http://api.ipstack.com/%s?access_key=%s", ip, apiKey)
	resp, err := http.Get(endPoint)
	if err != nil {
		log.Println("Error while Calling : ", endPoint, err)
	}
	defer resp.Body.Close()
	var ipInfo IPInfo
	err = json.NewDecoder(resp.Body).Decode(&ipInfo)
	return ipInfo, err
}

func (cd *ClickDetails) InsertClickDetails() error {
	collection := GetMongoClient().Database("shorturl").Collection("click_details")
	data := bson.M{
		"ip_info":      cd.IPInfo,
		"short_url":    cd.ShortURL,
		"original_url": cd.OriginalURL,
		"current_time": time.Now(),
	}
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		log.Println("Error while InsertClickDetails() : ", err)
		return err
	}
	fmt.Println("Id : ", res.InsertedID)
	return nil
}

func (cd *ClickDetails) GetClickDetails() (string, error) {
	collection := GetMongoClient().Database("shorturl").Collection("click_data")
	filter := bson.M{"short_url": cd.ShortURL}
	res, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println("Error while GetClickDetails() : ", err)
		return "", err
	}
	var clickDetails ClickDetails
	i := 0
	fmt.Println("res : ", res.Decode(&clickDetails))
	for res.Next(ctx) {
		i++
		log.Println("times : ", i)
		var c ClickDetails
		err = res.Decode(&c)
		if err != nil {
			log.Println("Error in Decode of GetClickDetails() : ", err)
		}
		//clickDetails = append(clickDetails, c)
	}
	log.Println("clickDetails")
	return "clickDetails", nil
}

func (cd *ClickDetails) DeteteClickDetails() error {
	collection := GetMongoClient().Database("shorturl").Collection("click_details")
	filter := bson.M{
		"short_url": cd.ShortURL,
	}
	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Println("Error in Delete of DeleteClickDetails : ", err)
		return err
	}
	log.Println("Deleted Properly ", cd.ShortURL)
	return nil
}
