package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nicograef/cloudevents/database/database"
)

func main() {
	server := &http.Server{
		Addr: ":5000",
	}

	appDatabase, err := database.LoadDatabaseFromJsonFile()
	if err != nil {
		fmt.Println("No existing database found, creating a new one.")
		appDatabase = database.New()
	} else {
		fmt.Println("Loaded existing database from file.")
	}

	// Set up signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("Signal received, persisting database...")
		if err := appDatabase.PersistToJsonFile(); err != nil {
			fmt.Println("Error persisting database:", err)
		}
		os.Exit(0)
	}()

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("Server error:", err)
		if persistErr := appDatabase.PersistToJsonFile(); persistErr != nil {
			fmt.Println("Error persisting database:", persistErr)
		}
		panic(err)
	}
}
