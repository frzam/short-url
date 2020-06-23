package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func GetIPInfo(ip string) IPInfo {
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
	if err != nil {
		log.Println("Error while GetIPInfo Decode : ", err)
	}
	return ipInfo
}

func (cd *ClickDetails) InsertClickDetails() error {
	collection := GetMongoClient().Database("shorturl").Collection("click_details")

	res, err := collection.InsertOne(ctx, cd)
	if err != nil {
		log.Println("Error while InsertClickDetails() : ", err)
		return err
	}
	fmt.Println("Id : ", res.InsertedID)
	return nil
}

func (cd *ClickDetails) GetTotalClicksDetails(skip, limit int64) ([]*ClickDetails, error) {
	return cd.GetNdayClicksDetails(0, skip, limit)
}

func (cd *ClickDetails) GetNdayClicksDetails(days int, skip, limit int64) ([]*ClickDetails, error) {
	collection := GetMongoClient().Database("shorturl").Collection("click_details")
	var from time.Time
	if days != 0 {
		from = time.Now().AddDate(0, 0, -1*days)
	}
	fmt.Println("from : ", from)

	opts := &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	}
	filter := bson.M{
		"shorturl": cd.ShortURL,
		"currenttime": bson.M{
			"$gt": from,
		},
	}
	res, err := collection.Find(ctx, filter, opts)
	if err != nil {
		log.Println("Error while GetClickDetails() : ", err)
		return nil, err
	}
	var clickDetails []*ClickDetails
	for res.Next(ctx) {
		var c ClickDetails
		err = res.Decode(&c)
		if err != nil {
			log.Println("Error in Decode of GetClickDetails() : ", err)
		}
		clickDetails = append(clickDetails, &c)
	}
	fmt.Println(clickDetails)
	return clickDetails, nil
}

func (cd *ClickDetails) DeleteClickDetails() error {
	collection := GetMongoClient().Database("shorturl").Collection("click_details")
	filter := bson.M{
		"shorturl": cd.ShortURL,
	}
	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Println("Error in Delete of DeleteClickDetails : ", err)
		return err
	}
	log.Println("Deleted Properly ", cd.ShortURL)
	return nil
}

// GetTotalClicksCount will return the total clicks count for a particular url.
func (cd *ClickDetails) GetTotalClicksCount() (int, error) {
	collection := GetMongoClient().Database("shorturl").Collection("click_details")
	filter := bson.M{
		"shorturl": cd.ShortURL,
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Println("Error in GetTotalCount : ", err)
		return -1, err
	}
	fmt.Println("Count : ", count)
	return int(count), nil
}

// GetNdayClicksCount will return the click counts for past n days.
func (cd *ClickDetails) GetNdayClicksCount(days int) (int, error) {
	collection := GetMongoClient().Database("shorturl").Collection("click_details")
	filter := bson.M{
		"shorturl": cd.ShortURL,
		"currenttime": bson.M{
			"$gt": time.Now().AddDate(0, 0, -1*days),
		},
	}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Println("Error in GetTodayClicksCount : ", err)
		return -1, err
	}
	fmt.Println("Today Count : ", count)
	return int(count), nil
}

func (cd *ClickDetails) GetClicksDetailsByCountry(country string, skip, limit int64) ([]*ClickDetails, error) {
	collection := GetMongoClient().Database("shorturl").Collection("click_details")
	filter := &bson.M{
		"shorturl":           cd.ShortURL,
		"ipinfo.countryname": country,
	}
	opts := &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	}
	res, err := collection.Find(ctx, filter, opts)
	if err != nil {
		log.Println("Error in GetClicksDetailsByCountry : ", err)
		return nil, err
	}
	var clickDetails []*ClickDetails
	for res.Next(ctx) {
		var cd ClickDetails
		err = res.Decode(&cd)
		if err != nil {
			log.Println("Error while Decode in GetClicksDetailsByCountry : ", err)
		} else {
			clickDetails = append(clickDetails, &cd)
		}
	}
	fmt.Println("clickDetalis : ", clickDetails)
	return clickDetails, err
}
