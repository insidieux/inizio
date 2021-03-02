package main

import (
	"github.com/insidieux/inizio/cmd/inizio/internal"
	"github.com/insidieux/inizio/internal/logger"
	"github.com/sethvargo/go-signalcontext"
)

func main() {
	ctx, cancel := signalcontext.OnInterrupt()
	if cancel != nil {
		defer cancel()
	}
	if err := internal.NewCommand().ExecuteContext(ctx); err != nil {
		logger.GetLogger().Fatalf(`Failed to execute command: %s`, err.Error())
	}
}
