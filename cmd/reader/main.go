package main

import (
	"context"
	"github.com/patyukin/go-redis-streams/internal/app"
	"github.com/patyukin/go-redis-streams/internal/config"
	"log"
	"log/slog"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadYamlConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	a, err := app.NewApp(ctx, cfg, logger)
	if err != nil {
		log.Fatalf("Failed to init app: %v\n", err)
	}

	err = a.Run(ctx)
	if err != nil {
		log.Fatalf("Failed to run: %v\n", err)
	}
}
