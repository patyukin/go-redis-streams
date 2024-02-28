package app

import (
	"context"
	"time"
)

func (a *App) streamSender(ctx context.Context) error {
	// Отправвка сообщения
	wrapCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := a.serviceProvider.streamer.Publish(wrapCtx)
	if err != nil {
		return err
	}

	return nil
}
