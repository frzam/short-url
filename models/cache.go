package models

import (
	"errors"
	"fmt"
	"log"
	"os"

	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

// rdb is the redis client for the redisDB.
var rdb *redis.Client

// init is used to initialize the redis db client.
func init() {
	_ = godotenv.Load()
	host := os.Getenv("cache_db_host")
	port := os.Getenv("cache_db_port")
	uri := fmt.Sprintf("%s:%s", host, port)

	rdb = redis.NewClient(&redis.Options{
		Addr:     uri,
		DB:       0,
		Password: "",
	})
	pong, _ := rdb.Ping(ctx).Result()
	log.Println("Redis Ping - ", pong)
}

// GetRedisClient func is used instead of directly using the package level redis.Client
func GetRedisClient() *redis.Client {
	return rdb
}

type Cache interface {
	Set() error
	Get() error
}

// Set is used to set the short-url with original url. It sets for 2 hours.
func (url *URL) Set() error {
	return GetRedisClient().Set(ctx, url.ShortURL, []byte(fmt.Sprintf("%v", url)), time.Hour*2).Err()
}

// Get is used to get the cache original url quickly.
func (url *URL) Get() error {
	res := GetRedisClient().Get(ctx, url.ShortURL)
	err := res.Scan(url)
	if err != nil {
		fmt.Println("Error in Get() : ", err)
	}
	fmt.Println("url inside Get() : ", url)
	return err
}

// Get is used to get the click details from cache.
func (cd *ClickDetails) Get() error {
	return GetRedisClient().Get(ctx, cd.IPInfo.IP).Scan(&cd)
}

// Set is used to set the ip with with the clickDetails for 2 hours.
func (cd *ClickDetails) Set() error {
	return GetRedisClient().Set(ctx, cd.IPInfo.IP, []byte(fmt.Sprintf("%v", cd)), time.Hour*2).Err()
}

// Get is used to get the the IPInfo from ip or it returns the error.
func (ipInfo *IPInfo) Get() error {
	res := GetRedisClient().Get(ctx, ipInfo.IP)
	if res == nil {
		return errors.New("Not Found")
	}
	return res.Scan(&ipInfo)

}

// Set is used to set the ip with ipInfo for two hours in cache.
func (ipInfo IPInfo) Set() error {
	return GetRedisClient().Set(ctx, ipInfo.IP, []byte(fmt.Sprintf("%v", ipInfo)), time.Hour*2).Err()
}
