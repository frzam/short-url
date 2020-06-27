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
	host := os.Getenv("cacheDB_host")
	port := os.Getenv("cacheDB_port")
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

// SetCacheURL is used to set the short-url with original url. It sets for 2 hours.
func (url *URL) SetCacheURL() error {
	return GetRedisClient().Set(ctx, url.ShortURL, url.OriginalURL, time.Hour*2).Err()
}

// GetCacheURl is used to get the cache original url quickly.
func (url *URL) GetCacheURL() (string, error) {
	return GetRedisClient().Get(ctx, url.ShortURL).Result()
}

// GetCacheClickDetails is used to get the click details from cache.
func (cd *ClickDetails) GetCacheClickDetails() error {
	return GetRedisClient().Get(ctx, cd.IPInfo.IP).Scan(&cd)
}

// SetCacheClickDetails is used to set the ip with with the clickDetails for 2 hours.
func (cd *ClickDetails) SetCacheClickDetails() error {
	return GetRedisClient().Set(ctx, cd.IPInfo.IP, []byte(fmt.Sprintf("%v", cd)), time.Hour*2).Err()
}

// GetCacheIPInfo is used to get the the IPInfo from ip or it returns the error.
func GetCacheIPInfo(ip string) (IPInfo, error) {
	var ipInfo IPInfo
	res := GetRedisClient().Get(ctx, ip)
	if res == nil {
		return ipInfo, errors.New("Not Found")
	}
	res.Scan(&ipInfo)
	return ipInfo, nil
}

// SetCacheIPInfo is used to set the ip with ipInfo for two hours in cache.
func (ipInfo IPInfo) setCacheIPInfo() error {
	return GetRedisClient().Set(ctx, ipInfo.IP, []byte(fmt.Sprintf("%v", ipInfo)), time.Hour*2).Err()
}
