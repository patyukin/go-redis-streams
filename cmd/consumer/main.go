package main

import (
	"context"
	"github.com/patyukin/go-redis-streams/internal/app"
	"github.com/patyukin/go-redis-streams/internal/config"
	"github.com/patyukin/go-redis-streams/internal/processor"
	"github.com/patyukin/go-redis-streams/internal/storage"
	"github.com/patyukin/go-redis-streams/internal/streamer"
	"github.com/patyukin/go-redis-streams/pkg/logger"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.LoadYamlConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	l, err := logger.InitLogger(cfg.Logger.Path, "dev")
	if err != nil {
		l.Fatal("Failed to init logger", zap.String("error", err.Error()))
	}

	ctx = context.WithValue(ctx, "logger", l)

	strg, err := storage.NewStorage(cfg.MySQL.DSN)
	if err != nil {
		l.Fatal("failed connect MySQL")
	}

	logger.Debug(ctx, "connected to mysql")

	// storage
	reader, err := streamer.NewRedisStreamer(ctx, cfg)
	processors := []app.Processor{
		processor.NewBookProcessor(reader, strg),
		processor.NewJournalProcessor(reader, strg),
		processor.NewBookProcessor(reader, strg),
		processor.NewJournalProcessor(reader, strg),
		processor.NewBookProcessor(reader, strg),
		processor.NewJournalProcessor(reader, strg),
	}
	if err != nil {
		l.Fatal("Failed to init redis", zap.String("error", err.Error()))
	}

	logger.Debug(ctx, "connected to redis")

	a, err := app.NewApp(cfg, processors)
	if err != nil {
		l.Fatal("Failed to init app", zap.String("error", err.Error()))
	}

	errCh := make(chan error)
	go func() {
		err = a.RunConsume(ctx)
		if err != nil {
			errCh <- err
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case err = <-errCh:
		l.Error("Failed to run", zap.String("error", err.Error()))
	case res := <-sigChan:
		if res == syscall.SIGINT || res == syscall.SIGTERM {
			l.Info("Signal received", zap.String("signal", res.String()))
		} else if res == syscall.SIGHUP {
			l.Info("Signal received", zap.String("signal", res.String()))
		}
	}

	cancel()

	// gf shutdown
	err = reader.Close()
	if err != nil {
		l.Warn("Failed to close consumer", zap.String("error", err.Error()))
	}
}
