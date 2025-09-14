package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"context"
	"sync"

	"github.com/nicograef/qugo/api"
	"github.com/nicograef/qugo/config"
	"github.com/nicograef/qugo/queue"
)

func main() {
	cfg := config.Load()
	appQueue := queue.NewQueue(cfg.Capacity)

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
	}
	http.HandleFunc("/", api.NewEnqueueHandler(appQueue))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for item := range appQueue.Queue {
			appQueue.HandleQueueItem(item, cfg.ConsumerURL, queue.SendToWebhook)
		}
	}()

	// Set up signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Println("Signal received, shutting down...")
		// Gracefully shutdown HTTP server
		ctx, cancel := context.WithTimeout(context.Background(), 5_000_000_000) // 5s
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		// Close queue channel to stop consumer
		close(appQueue.Queue)
		// Wait for consumer goroutine to finish
		wg.Wait()
		log.Println("Shutdown complete.")
		os.Exit(0)
	}()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
