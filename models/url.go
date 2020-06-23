package models

// TO DO :
// 1. Connect to MongoDB and create a new client.	--> Done.
// 2. Create a new database and a collection to store the details of url.	--> Done.
// 3. Make func to add a document in collection.	--> Done.
// 4. Make func to get a document from collection.	--> Done.
// 5. Make func to delete a document from collection.	--> Done.

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

var client *mongo.Client
var ctx = context.Background()

type URL struct {
	ShortURL       string    `bson:"short_url"`
	OriginalURL    string    `bson:"original_url"`
	CreationDate   time.Time `bson:"creation_date"`
	ExpirationDate time.Time `bson:"expiration_date"`
	UserID         int       `bson:"user_id"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error while loading the .env file")
	}
	name := os.Getenv("primaryDB_name")
	host := os.Getenv("primaryDB_host")
	port := os.Getenv("primaryDB_port")

	uri := fmt.Sprintf("%s://%s:%s", name, host, port)
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Error while Ping to mongoDB : ", err)
	}

	// 	_, _ = cd.GetClickDetails()
	// 	//_ = cd.DeteteClickDetails()
	// 	_ = cd.SetCacheClickDetails()
	// 	_ = cd.GetCacheClickDetails()
	// 	fmt.Println("cd : ", cd.IPInfo)

	cd := &ClickDetails{
		ShortURL:    "s-url",
		CurrentTime: time.Now(),
	}
	//_ = cd.InsertClickDetails()
	_, _ = cd.GetNdayClicksDetails(2, 20, 10)

	_, _ = cd.GetTotalClicksCount()
	_, _ = cd.GetNdayClicksCount(2)
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
func (url *URL) DeleteURL() error {
	log.Println("Inside DeleteURL()")
	collection := GetMongoClient().Database("shorturl").Collection("url")
	filter := bson.M{
		"_id": url.ShortURL,
	}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Println("Error while DeleteURL : ", err)
		return err
	}
	return nil
}

func (url *URL) prepareURL() {
	url.ShortURL = generateHash(url.OriginalURL)
	url.CreationDate = time.Now()
	url.ExpirationDate = time.Now().AddDate(0, 1, 0)
}

// GetURL is used to retrive original url basis the short url and current date.
// If the current date is after or equal to expiry date then we are not returning anything.
// It returns the original url and error.
func (url *URL) GetURL() (string, error) {
	fmt.Println("GetURL() Called!")
	collection := GetMongoClient().Database("shorturl").Collection("url")
	filter := bson.M{
		"_id": url.ShortURL,
		"expiration_date": bson.M{
			"$gt": time.Now(),
		},
	}
	err := collection.FindOne(ctx, filter).Decode(&url) //
	if err != nil {
		log.Println("Error while getting the document :", err)
		return "", nil
	}
	_ = url.SetCacheURL()
	return url.OriginalURL, nil
}

func (url *URL) AddClickDetails(ip string) error {
	// TO DO : think about calling the Cache for same ip.
	cd := &ClickDetails{
		OriginalURL: url.OriginalURL,
		ShortURL:    url.ShortURL,
		CurrentTime: time.Now(),
		IPInfo:      GetIPInfo(ip),
	}
	err := cd.InsertClickDetails()
	if err != nil {
		log.Println("Error while Calling InsertClickDetails() : ", err)
		return err
	}
	// If in cache don't again insert in cache. Otherwise
	err = cd.SetCacheClickDetails()
	if err != nil {
		log.Println("Error while Calling SetCacheClickDetails() ")
	}
	return err
}

// GetMongoClient  gives the mongoDB client.
func GetMongoClient() *mongo.Client {
	return client
}
