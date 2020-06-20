package models

// TO DO :
// 1. Connect to MongoDB and create a new client.	--> Done.
// 2. Create a new database and a collection to store the details of url.
// 3. Make func to add a document in collection.
// 4. Make func to get a document from collection.
// 5. Make func to delete a document from collection.

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client
var ctx = context.Background()

type URL struct {
	ShortURL       string `json:"short_url"`
	OriginalURL    string `json:"original_url"`
	CreationDate   string `json:"creation_date"`
	ExpirationDate string `json:"expiration_date"`
	UserID         int    `json:"user_id"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading the .env file")
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
	fmt.Println("err : ", err)
	url := &URL{
		OriginalURL: "https://github.com/mongodb/mongo-go-driver#usage",
		UserID:      1,
	}
	url.InsertIntoDB()
}

func (url *URL) InsertIntoDB() {
	shortURL := generateHash(url.OriginalURL)
	url.ShortURL = shortURL
	url.CreationDate = time.Now().Format(time.RFC850)
	if url.ExpirationDate == "" {
		url.ExpirationDate = time.Now().AddDate(0, 1, 0).Format(time.RFC850)
	}
	collection := GetMongoClient().Database("shorturl").Collection("url")
	res, err := collection.InsertOne(ctx, url)
	if err != nil {
		log.Fatal("Error while inserting into DB : ", err)
	}
	fmt.Println("id : ", res.InsertedID)
}

// GetMongoClient  gives the mongoDB client.
func GetMongoClient() *mongo.Client {
	return client
}
