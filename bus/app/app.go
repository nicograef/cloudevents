package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nicograef/cloudevents/bus/api"
	"github.com/nicograef/cloudevents/bus/bus"
	"github.com/nicograef/cloudevents/bus/config"
)

type App struct {
	Server *http.Server
	Config config.Config
	router *http.ServeMux
}

// NewApp creates a new application instance
func NewApp(cfg config.Config) (*App, error) {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	router := http.NewServeMux()

	return &App{
		Server: server,
		Config: cfg,
		router: router,
	}, nil
}

// SetupRoutes configures HTTP routes
func (app *App) SetupRoutes() {
	app.router.HandleFunc("POST /publish", api.NewPublishHandler(bus.NewPublish(app.Config.Subscribers, bus.SendToWebhook)))
	app.router.HandleFunc("GET /health", api.NewHealthHandler())
	app.Server.Handler = app.router
}

// Run starts the application with graceful shutdown
func (app *App) Run(ctx context.Context) error {
	app.SetupRoutes()

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

// Shutdown gracefully stops the application
func (app *App) Shutdown() error {
	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := app.Server.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	fmt.Println("Shutdown complete")
	return nil
}
