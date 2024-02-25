package main

import (
	"context"
	"github.com/patyukin/go-redis-streams/internal/streamer/redis"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Println("Publisher started")

	redisClient := redis.NewRedisStreamer(ctx)

	for i := 0; i < 3000; i++ {
		err := publishTicketReceivedEvent(ctx, redisClient)
		time.Sleep(time.Second * 5)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func publishTicketReceivedEvent(ctx context.Context, redisClient *redis.Streamer) error {
	log.Println("Publishing event to Redis")
	err := redisClient.Publish(ctx)
	return err
}
