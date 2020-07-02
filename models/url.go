package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// client is *mongo.Client that is only used with GetMongoClient().
// In this package or outside this package.
var client *mongo.Client

// ctx is running in the background.
var ctx = context.Background()

// URL struct contains complete details about a url.
// It contains the OrginalURL provided by the client, it
// contains the UserID. URL collection in mongoDB contains data
// in the format of url struct.
type URL struct {
	ShortURL       string    `bson:"short_url"`
	OriginalURL    string    `bson:"original_url"`
	CreationDate   time.Time `bson:"creation_date"`
	ExpirationDate time.Time `bson:"expiration_date"`
	UserID         int       `bson:"user_id"`
}

// init is used to load the environment variables and establish mongoDB
// connection. It sets *mongo.Client and then this client is used by
// by all the methods and function using GetMongoClient().
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error while loading the .env file")
	}
	name := os.Getenv("primary_db_name")
	host := os.Getenv("primary_db_host")
	port := os.Getenv("primary_db_port")

	uri := fmt.Sprintf("%s://%s:%s", name, host, port)
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Error while Ping to mongoDB : ", err)
	}
}

// InsertURL is used to insert a new url into the collection.
// It calls the prepareURL method that is used to generate the url, currentDate
// expiry date. Then InsertURL() inserts the document with _id = shortURL.
// It returns the error.
func (url *URL) InsertURL() error {
	url.prepareURL()
	collection := GetMongoClient().Database("shorturl").Collection("url")
	data := bson.M{
		"_id":             url.ShortURL,
		"short_url":       url.ShortURL,
		"original_url":    url.OriginalURL,
		"creation_date":   url.CreationDate,
		"expiration_date": url.ExpirationDate,
		"user_id":         url.UserID,
	}
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		log.Println("Error while inserting into DB : ", err)
		return err
	}
	fmt.Println("id : ", res.InsertedID)
	return nil
}

// DeleteURL is used to delete a short url from the url collection.
// It doesn't delete the click details.
func (url *URL) DeleteURL() error {
	log.Println("URL Deleted : ", url.ShortURL)
	collection := GetMongoClient().Database("shorturl").Collection("url")
	filter := bson.M{
		"_id": url.ShortURL,
	}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Println("Error in DeleteURL() : ", err)
		return err
	}
	return nil
}

// prepareURL calls generateHash() to generateShortURL and it set the expiration
// time to one month from the current time.
func (url *URL) prepareURL() {
	url.ShortURL = generateHash(url.OriginalURL)
	url.CreationDate = time.Now()
	url.ExpirationDate = time.Now().AddDate(0, 1, 0)
}

// GetURL is used to retrive original url basis the short url and current date.
// If the current date is after or equal to expiry date then we are not returning anything.
// It returns the original url and error. Before returning it sets the into the cache.
func (url *URL) GetURL() (string, error) {
	collection := GetMongoClient().Database("shorturl").Collection("url")
	filter := bson.M{
		"_id": url.ShortURL,
		"expiration_date": bson.M{
			"$gt": time.Now(),
		},
	}
	err := collection.FindOne(ctx, filter).Decode(&url)
	if err != nil {
		log.Println("Error while getting the document :", err)
		return "", nil
	}
	_ = url.Set()
	return url.OriginalURL, nil
}

// AddClickDetails is used to add click details for a shorturl.
// It gets the originalURL from GetCacheURL, then it gets the GetCacheIPInfo.
// If it is not present then then it calls GetIPInfo to call the api to get the ip detals.
// Then it inserts the click details, and sets the caches the click details.
func (url *URL) AddClickDetails(ip string) error {
	_, _ = url.Get()
	cd := &ClickDetails{
		OriginalURL: url.OriginalURL,
		ShortURL:    url.ShortURL,
		IPInfo:      GetIPInfo(ip),
		CurrentTime: time.Now(),
	}
	err := cd.InsertClickDetails()
	if err != nil {
		log.Println("Error while Calling InsertClickDetails() : ", err)
		return err
	}
	// If in cache don't again insert in cache. Otherwise
	err = cd.Set()
	if err != nil {
		log.Println("Error while Calling SetCacheClickDetails() ")
	}
	return err
}

// GetMongoClient  gives the mongoDB client.
func GetMongoClient() *mongo.Client {
	return client
}
