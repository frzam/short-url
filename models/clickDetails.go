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

// ClickDetails contains the complete details one a client clicks on the shorturl.
type ClickDetails struct {
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	CurrentTime time.Time `json:"current_time"`
	IPInfo      IPInfo    `json:"ip_info"`
}

// IPInfo contains the ip related information. It is the reponse from ipstack api.
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

// Languages stuct is used denote one language.
type Languages struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Native string `json:"native"`
}

// Location defines the location of a client in geometrical context.
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

// GetIPInfo is used to call ipstack api and it returns the IPInfo instance.
// This contains complete information about one ip.
func GetIPInfo(ip string) IPInfo {
	apiKey := os.Getenv("ipstack_api_key")
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

// InsertClickDetails is used to insert clickDetails object inside click_details collection.
// It returns error if any.
func (cd *ClickDetails) InsertClickDetails() error {
	collection := GetMongoClient().Database("shorturl").Collection("click_details")

	res, err := collection.InsertOne(ctx, cd)
	if err != nil {
		log.Println("Error while InsertClickDetails() : ", err)
		return err
	}
	log.Println("ID : ", res.InsertedID)
	return nil
}

// GetTotalClicksDetails is used to return the slice of *ClickDetails.
// It calls GetNdayClicksDetails with days as zero. Thus the called method returns all
// the *ClickDetails from 00-00-0000.
func (cd *ClickDetails) GetTotalClicksDetails(skip, limit int64) ([]*ClickDetails, error) {
	return cd.GetNdayClicksDetails(0, skip, limit)
}

// GetNdayClicksDetails returns the slice of ClickDetails and error if any for past n days.
// It takes skip and limit params to filter out the result (pagination).
func (cd *ClickDetails) GetNdayClicksDetails(days int, skip, limit int64) ([]*ClickDetails, error) {
	collection := GetMongoClient().Database("shorturl").Collection("click_details")
	var from time.Time
	if days != 0 {
		from = time.Now().AddDate(0, 0, -1*days)
	}

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
	return clickDetails, nil
}

// DeleteClickDetails is used to delete the complete click details data basis shorturl.
// It returns error if any.
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
	return int(count), nil
}

// GetClicksDetailsByCountry is used to retrieve the click details by country.
// It is passed with country, skip and limit.
// It returns the slice of *ClickDetails and error if any.
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
	return clickDetails, err
}

// GetClicksDetailsByCity is used to get the click details by city.
// It takes skip and limit and uses it in filter.
// It returns slice of *ClickDetails and error if any.
func (cd *ClickDetails) GetClicksDetailsByCity(city string, skip, limit int64) ([]*ClickDetails, error) {
	collection := GetMongoClient().Database("shorturl").Collection("click_details")
	filter := bson.M{
		"shorturl":    cd.ShortURL,
		"ipinfo.city": city,
	}
	opts := &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	}
	res, err := collection.Find(ctx, filter, opts)
	if err != nil {
		log.Println("Error in GetClicksDetailsByCity() : ", err)
		return nil, err
	}
	var clickDetails []*ClickDetails
	for res.Next(ctx) {
		var cd *ClickDetails
		err := res.Decode(&cd)
		if err != nil {
			log.Println("Error while Decoding in GetClicksDetailsBydCity() : ", err)
		} else {
			clickDetails = append(clickDetails, cd)
		}
	}
	return clickDetails, nil
}

// GetClicksCountByIP is returns the total count for one ip.
// It returns clickCount and error if any.
func (cd *ClickDetails) GetClicksCountByIP(ip string) (int, error) {
	collection := GetMongoClient().Database("shorturl").Collection("click_details")
	filter := bson.M{
		"shorturl":  cd.ShortURL,
		"ipinfo.ip": ip,
	}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Println("Error in GetClicksCountByIP : ", err)
		return -1, err
	}
	return int(count), nil
}

// GetClicksDetailsByIP is used to get the slice of ClickDetails for one particular client IP.
// It takes skip and limit param to filter the result.
func (cd *ClickDetails) GetClicksDetailsByIP(ip string, skip, limit int64) ([]*ClickDetails, error) {
	collection := GetMongoClient().Database("shorturl").Collection("click_details")
	filter := bson.M{
		"shorturl":  cd.ShortURL,
		"ipinfo.ip": ip,
	}
	opts := &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	}
	res, err := collection.Find(ctx, filter, opts)
	if err != nil {
		log.Println("Error in GetClicksDetailsByIP : ", err)
		return nil, err
	}
	var clickDetails []*ClickDetails
	for res.Next(ctx) {
		var cd *ClickDetails
		err := res.Decode(&cd)
		if err != nil {
			log.Println("Error while Decode of GetClicksDetailsByIP : ", err)
		} else {
			clickDetails = append(clickDetails, cd)
		}
	}
	return clickDetails, nil
}
