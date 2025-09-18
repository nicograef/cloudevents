package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nicograef/cloudevents/bus/app"
	"github.com/nicograef/cloudevents/bus/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("FATAL Configuration error: %v", err)
	}

	app, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("FATAL Failed to create app: %v", err)
	}

	// Set up signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Run application
	if err := app.Run(ctx); err != nil {
		log.Fatalf("FATAL Application error: %v", err)
	}
}
