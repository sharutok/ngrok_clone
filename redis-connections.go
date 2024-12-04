package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var ctx = context.Background()

func ConnectToRedis() *redis.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	i, err := strconv.Atoi(os.Getenv("REDIS_DB"))

	if err != nil {
		log.Println(err, "error in db redis")
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       i,
	})

	return rdb
}
