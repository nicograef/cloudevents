package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/nicograef/cloudevents/queue/api"
	"github.com/nicograef/cloudevents/queue/config"
	"github.com/nicograef/cloudevents/queue/queue"
)

type App struct {
	Queue  queue.Queue
	Server *http.Server
	Config config.Config
	router *http.ServeMux
	wg     sync.WaitGroup
}

// NewApp creates a new application instance
func NewApp(cfg config.Config) (*App, error) {
	appQueue := queue.NewQueue(cfg.Capacity)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	router := http.NewServeMux()

	return &App{
		Queue:  appQueue,
		Server: server,
		Config: cfg,
		router: router,
	}, nil
}

// SetupRoutes configures HTTP routes
func (app *App) SetupRoutes() {
	app.router.HandleFunc("POST /enqueue", api.NewEnqueueHandler(app.Queue))
	app.router.HandleFunc("GET /health", api.NewHealthHandler())
}

// Run starts the application with graceful shutdown
func (app *App) Run(ctx context.Context) error {
	app.SetupRoutes()

	// Start queue consumer in goroutine
	app.startQueueConsumer()

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		fmt.Printf("Starting server on port %d\n", app.Config.Port)
		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
		fmt.Println("Shutdown signal received, gracefully stopping...")
		return app.Shutdown()
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	}
}

// startQueueConsumer starts the queue consumer goroutine
func (app *App) startQueueConsumer() {
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		for item := range app.Queue.Queue {
			app.Queue.HandleQueueItem(item, app.Config.ConsumerURL, queue.SendToWebhook)
		}
	}()
}

// Shutdown gracefully stops the application
func (app *App) Shutdown() error {
	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := app.Server.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	// Close queue channel to stop consumer
	close(app.Queue.Queue)

	// Wait for consumer goroutine to finish
	app.wg.Wait()

	fmt.Println("Shutdown complete")
	return nil
}
