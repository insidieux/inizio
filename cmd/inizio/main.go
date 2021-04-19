package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/insidieux/inizio/cmd/inizio/internal"
	"github.com/insidieux/inizio/internal/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := internal.NewCommand().ExecuteContext(ctx); err != nil {
		logger.GetLogger().Fatalf(`Failed to execute command: %s`, err.Error())
	}
}
