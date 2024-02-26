package app

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func (a *App) processMessageB(_ context.Context, message redis.XMessage) error {
	// Обработка сообщения
	fmt.Printf("Message: %v\n", message)

	return nil
}
