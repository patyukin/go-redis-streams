package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"math/rand"
	"sync"
)

type Streamer struct {
	c *redis.Client
}

func NewRedisStreamer(ctx context.Context) *Streamer {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", "127.0.0.1", "6379"),
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Unbale to connect to Redis", err)
	}

	return &Streamer{
		c: client,
	}
}

func (s *Streamer) Consume(ctx context.Context, stream string) {
	var wg *sync.WaitGroup

	cursor := "0"

	for {
		result, err := s.c.XRead(ctx, &redis.XReadArgs{
			Streams: []string{stream, cursor},
			Count:   10,
			Block:   0,
		}).Result()

		if err != nil {
			log.Fatalf("Failed to read from stream: %v\n", err)
			return
		}

		// Обработка сообщений из стрима
		for _, xStream := range result {
			for _, message := range xStream.Messages {
				wg.Add(1)
				go func(ctx context.Context, m redis.XMessage) {
					defer wg.Done()
					//processMessage(ctx, m)
				}(ctx, message)
			}
		}

		// Обновляем курсор
		if len(result) > 0 {
			cursor = result[0].Messages[len(result[0].Messages)-1].ID
		}
	}
}

func (s *Streamer) LimitConsume(ctx context.Context, stream string, processMessage func(ctx context.Context, m redis.XMessage) error) {

	var wg *sync.WaitGroup
	// Создаем пул горутин-работников (Worker Pool)
	workerPool := make(chan struct{}, 6)
	defer wg.Done()

	cursor := "0"

	for {
		result, err := s.c.XRead(ctx, &redis.XReadArgs{
			Streams: []string{stream, cursor},
			Count:   10,
			Block:   0,
		}).Result()

		if err != nil {
			log.Fatalf("Failed to read from stream: %v\n", err)
			return
		}

		for _, xStream := range result {
			for _, message := range xStream.Messages {

				workerPool <- struct{}{}
				wg.Add(1)
				go func(ctx context.Context, m redis.XMessage) {
					defer func() {
						<-workerPool
						wg.Done()
					}()
					err = processMessage(ctx, m)
					if err != nil {
						// TODO - обработка ошибки
						return
					}
				}(ctx, message)
			}
		}

		if len(result) > 0 {
			cursor = result[0].Messages[len(result[0].Messages)-1].ID
		}
	}
}

func (s *Streamer) Publish(ctx context.Context) error {
	err := s.c.XAdd(ctx, &redis.XAddArgs{
		Stream: "tickets",
		MaxLen: 0,
		ID:     "",
		Values: map[string]interface{}{
			"whatHappened": "ticket received",
			"ticketID":     rand.Intn(100000000),
			"ticketData":   "some ticket data",
		},
	}).Err()

	return err
}
