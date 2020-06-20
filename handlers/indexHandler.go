package handlers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Creating a template instance so that we can execute our data into it.
var tpl = template.Must(template.ParseFiles("index.html"))

// IndexHandler is used to handle "/" path (HOME).
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

var ctx = context.Background()

func init() {
	rdb := redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
	pong, _ := rdb.Ping(ctx).Result()
	fmt.Println("Redis Ping : ", pong)

	// MongoDB.
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println("Error while Connecting to mongodb : ", err)
	}
	err = client.Ping(ctx, readpref.Primary())
	fmt.Println("MongoDb err : ", err)
	collection := client.Database("temp").Collection("users")
	res := collection.FindOne(ctx, nil)
	fmt.Println(res)
}
