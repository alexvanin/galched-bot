package grace

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func NewGracefulContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-ch:
			log.Print("ctx: caught interrupt signal")
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx
}
