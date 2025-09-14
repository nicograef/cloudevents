package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"context"
	"sync"

	"github.com/nicograef/qugo/api"
	"github.com/nicograef/qugo/config"
	"github.com/nicograef/qugo/core"
)

func main() {
	cfg := config.Load()

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
	}

	queue := make(chan core.Message, cfg.QueueSize)

	http.HandleFunc("/", api.NewEnqueueHandler(queue))

	// WaitGroup for consumer goroutine
	var wg sync.WaitGroup
	wg.Add(1)

	// Consumer goroutine: reads from channel and calls webhook
	go func() {
		defer wg.Done()
		for msg := range queue {
			resp, err := core.SendToWebhook(cfg.ConsumerUrl, msg)
			if err != nil {
				fmt.Printf("Error sending to webhook: %v\n", err)
			} else {
				fmt.Printf("Webhook response: %s\n", resp)
			}
		}
	}()

	// Set up signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("Signal received, shutting down...")
		// Gracefully shutdown HTTP server
		ctx, cancel := context.WithTimeout(context.Background(), 5_000_000_000) // 5s
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			fmt.Printf("HTTP server Shutdown: %v\n", err)
		}
		// Close queue channel to stop consumer
		close(queue)
		// Wait for consumer goroutine to finish
		wg.Wait()
		fmt.Println("Shutdown complete.")
		os.Exit(0)
	}()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
