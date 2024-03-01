package main

import (
	"context"
	"github.com/patyukin/go-redis-streams/internal/config"
	"github.com/patyukin/go-redis-streams/internal/streamer/redis"
	"github.com/patyukin/go-redis-streams/pkg/logger"
	"go.uber.org/zap"
	"log"
	"math/rand"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadYamlConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	l, err := logger.InitLogger(cfg.Logger.Path, "dev")
	if err != nil {
		l.Fatal("Failed to init logger", zap.String("error", err.Error()))
	}

	wrapCtx := context.WithValue(ctx, "logger", l)

	client := redis.NewRedisStreamer(wrapCtx, cfg)

	for i := 0; i < 100000000000; i++ {
		r := rand.Intn(2)
		if r == 0 {
			err = client.Publish(ctx, "books")
		} else {
			err = client.Publish(ctx, "journals")
		}
		if err != nil {
			l.Fatal("Failed to publish", zap.String("error", err.Error()))
		}
	}
}
