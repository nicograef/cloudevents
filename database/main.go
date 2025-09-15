package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nicograef/cloudevents/database/app"
	"github.com/nicograef/cloudevents/database/config"
)

func main() {
	cfg := config.Load()

	app, err := app.NewApp(cfg)
	if err != nil {
		fmt.Printf("Failed to create app: %v\n", err)
		os.Exit(1)
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
		fmt.Printf("Application error: %v\n", err)
		os.Exit(1)
	}
}
