package app

import (
	"context"
	"github.com/patyukin/go-redis-streams/internal/streamer/redis"
	goredis "github.com/redis/go-redis/v9"
	"log"
	"time"
)

type StreamerInterface interface {
	Publish(ctx context.Context) error
	Consume(ctx context.Context, stream string)
	LimitConsume(ctx context.Context, stream string, processMessage func(ctx context.Context, m goredis.XMessage) error)
}

type App struct {
	client StreamerInterface
}

func NewApp(client StreamerInterface) *App {
	return &App{
		client: client,
	}
}

func (a *App) RunTicket(ctx context.Context) error {
	client := redis.NewRedisStreamer(ctx)

	go client.LimitConsume(ctx, "tickets", a.processMessage)

	return nil
}

func (a *App) RunBook(ctx context.Context) error {
	client := redis.NewRedisStreamer(ctx)

	go client.LimitConsume(ctx, "books", a.processMessageB)

	return nil
}

func (a *App) RunPublisher(ctx context.Context) error {
	client := redis.NewRedisStreamer(ctx)

	for i := 0; i < 3000; i++ {
		err := publishTicketReceivedEvent(ctx, client)
		time.Sleep(time.Second * 5)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
