package main

import (
	"context"
	"github.com/patyukin/go-redis-streams/internal/app"
	"github.com/patyukin/go-redis-streams/internal/streamer/redis"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := redis.NewRedisStreamer(ctx)
	a := app.NewApp(client)

	done := make(chan struct{})

	go func() {
		err := a.Run(ctx)
		if err != nil {
			done <- struct{}{}
			log.Fatal(err)
		}
	}()

	go func() {
		err := a.Run(ctx)
		if err != nil {
			done <- struct{}{}
			log.Fatal(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	select {
	case <-done:
	case <-ch:
	}

	gracefullShotdown()
}

func publishTicketReceivedEvent(ctx context.Context, redisClient *redis.Streamer) error {
	log.Println("Publishing event to Redis")
	err := redisClient.Publish(ctx)
	return err
}
