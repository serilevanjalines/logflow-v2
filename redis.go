package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

func connectRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		log.Fatal("Unable to connect to Redis: ", err)
	}

	return client
}
