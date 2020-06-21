package models

import (
	"fmt"
	"log"
	"os"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var rdb *redis.Client

func init() {
	_ = godotenv.Load()
	host := os.Getenv("cacheDB_host")
	port := os.Getenv("cacheDB_port")
	uri := fmt.Sprintf("%s:%s", host, port)

	rdb = redis.NewClient(&redis.Options{
		Addr:     uri,
		DB:       0,
		Password: "",
	})
	pong, _ := rdb.Ping(ctx).Result()
	fmt.Println("Ping - ", pong)
}

func GetRedisClient() *redis.Client {
	return rdb
}

func (url *URL) SetCacheURL() error {
	fmt.Println("SetCacheURL() Called !")
	return GetRedisClient().Set(ctx, url.ShortURL, url.OriginalURL, time.Hour*2).Err()
}

func (url *URL) GetCacheURL() (string, error) {
	return GetRedisClient().Get(ctx, url.ShortURL).Result()
}

func (cd *ClickDetails) GetCacheClickDetails() error {
	return GetRedisClient().Get(ctx, cd.IPInfo.IP).Scan(&cd)
}

func (cd *ClickDetails) SetCacheClickDetails() error {
	log.Println("IP : ", cd.IPInfo.IP)
	log.Println("cd : ", *cd)
	return GetRedisClient().HMSet(ctx, cd.IPInfo.IP, cd, time.Minute*2).Err()
}
