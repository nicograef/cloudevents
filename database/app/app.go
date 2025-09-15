package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/nicograef/cloudevents/database/api"
	"github.com/nicograef/cloudevents/database/config"
	"github.com/nicograef/cloudevents/database/database"
)

type App struct {
	Database *database.Database
	Server   *http.Server
	Config   config.Config
	router   *http.ServeMux
}

// NewApp creates a new application instance
func NewApp(cfg config.Config) (*App, error) {
	appDatabase, err := database.LoadFromJSONFile(cfg.DataDir)
	if err != nil {
		fmt.Println("No existing database found, creating a new one.")
		appDatabase = database.New()
	} else {
		fmt.Println("Loaded existing database from file.")
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	router := http.NewServeMux()

	return &App{
		Database: appDatabase,
		Server:   server,
		Config:   cfg,
		router:   router,
	}, nil
}

// SetupRoutes configures HTTP routes
func (app *App) SetupRoutes() {
	app.router.HandleFunc("POST /add", api.NewAddEventHandler(*app.Database))
	app.router.HandleFunc("GET /health", api.NewHealthHandler())
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
		fmt.Printf("Error shutting down server: %v\n", err)
	}

	// Persist database
	fmt.Println("Persisting database...")
	if err := app.Database.PersistToJsonFile(app.Config.DataDir); err != nil {
		return fmt.Errorf("error persisting database: %w", err)
	}

	fmt.Println("Shutdown complete")
	return nil
}
