package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.TODO()

	streams, err := rdb.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"tickets", "$"},
	}).Result()

	if err != nil {
		panic(err)
	}

	for _, stream := range streams {
		fmt.Println("Stream:", stream.Stream)
		for _, message := range stream.Messages {
			fmt.Println("ID:", message.ID)
			for k, v := range message.Values {
				fmt.Println(k, v)
			}
		}
	}

}
