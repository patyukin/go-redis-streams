package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"math/rand"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Println("Publisher started")

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", "127.0.0.1", "6379"),
	})
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Unbale to connect to Redis", err)
	}

	log.Println("Connected to Redis server")

	for i := 0; i < 3000; i++ {
		err = publishTicketReceivedEvent(ctx, redisClient)
		time.Sleep(time.Second * 5)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func publishTicketReceivedEvent(ctx context.Context, client *redis.Client) error {
	log.Println("Publishing event to Redis")

	err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: "tickets",
		MaxLen: 0,
		ID:     "",
		Values: map[string]interface{}{
			"whatHappened": string("ticket received"),
			"ticketID":     int(rand.Intn(100000000)),
			"ticketData":   string("some ticket data"),
		},
	}).Err()

	return err
}
