package handlers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-redis/redis/v8"
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
}
