package main

import (
	"log"
	"net/url"
	"os"

	"github.com/i-tinerary/cotton/server"
)

func main() {
	port := os.Getenv("PORT")
	redisRawURL := os.Getenv("REDIS_URL")

	redisURL, err := url.Parse(redisRawURL)
	if err != nil {
		log.Fatal("Parsing redisurl: ", err)
	}

	log.Fatal("Serving: ", server.Serve(port, redisURL))
}
