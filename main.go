package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/oschwald/geoip2-golang"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if os.Getenv("ENVIRONMENT") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	router := gin.Default()
	router.Use(cors.New(config))
	router.GET("/", handler)
	fmt.Println(os.Getenv("ERSA"))
	router.Run("0.0.0.0:" + os.Getenv("PORT"))
}

var ctx = context.Background()

func redisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func handler(c *gin.Context) {
	log.Println("Hi from handler")
	clientIp := c.ClientIP()
	redisClient := redisClient()

	countryCode, success := getCountryCodeFromRedis(clientIp, redisClient)

	if !success {
		countryCode, success = getCountryCodeFromDb(clientIp)
	}

	if success {
		setCachedValue(clientIp, countryCode, redisClient)
	}

	response := codeResponse{Ip: clientIp,
		CountryCode: countryCode,
		Success:     success,
	}
	c.IndentedJSON(http.StatusOK, response)

}

type codeResponse struct {
	Ip          string `json:"ip"`
	CountryCode string `json:"countryCode"`
	Success     bool   `json:"success"`
}

func getCountryCodeFromDb(ip string) (string, bool) {
	db, err := geoip2.Open("db/country.mmdb")
	if err != nil {
		return "", false
	}

	record, err := db.Country(net.ParseIP(ip))
	if err != nil {
		return "", false
	}
	log.Println("Found in db")
	return record.Country.IsoCode, true
}

func getCountryCodeFromRedis(ip string, client *redis.Client) (string, bool) {
	log.Printf("Hi from redis ip: %s\n", ip)
	val, err := client.Get(ctx, ip).Result()
	if err != nil {
		log.Println(err)
		log.Println("Redis error")
		return "", false
	}
	log.Printf("Found in redis for ip %s value %s\n", ip, val)
	return val, true
}

func setCachedValue(ip string, countryCode string, client *redis.Client) bool {
	err := client.Set(ctx, ip, countryCode, 0)

	if err != nil {
		return false
	}

	return true
}
