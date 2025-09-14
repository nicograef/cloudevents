package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nicograef/cloudevents/database/api"
	"github.com/nicograef/cloudevents/database/config"
	"github.com/nicograef/cloudevents/database/database"
)

func main() {
	cfg := config.Load()

	appDatabase, err := database.LoadFromJSONFile(cfg.DataDir)
	if err != nil {
		fmt.Println("No existing database found, creating a new one.")
		appDatabase = database.New()
	} else {
		fmt.Println("Loaded existing database from file.")
	}

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
	}

	http.HandleFunc("/add", api.NewAddEventHandler(*appDatabase))

	// Set up signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("Signal received, persisting database...")
		if err := appDatabase.PersistToJsonFile(cfg.DataDir); err != nil {
			fmt.Println("Error persisting database:", err)
		}
		os.Exit(0)
	}()

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("Server error:", err)
		if persistErr := appDatabase.PersistToJsonFile(cfg.DataDir); persistErr != nil {
			fmt.Println("Error persisting database:", persistErr)
		}
		panic(err)
	}
}
