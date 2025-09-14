package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nicograef/qugo/api"
	"github.com/nicograef/qugo/queue"
)

func main() {
	server := &http.Server{
		Addr: ":3000",
	}

	appQueue, err := queue.LoadQueueFromJsonFile()
	if err != nil {
		fmt.Println("No existing queue found, creating a new one.")
		appQueue = queue.New()
	} else {
		fmt.Println("Loaded existing queue from file.")
	}

	http.HandleFunc("/enqueue", api.NewEnqueueHandler(appQueue))

	// Set up signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Printf("Signal %s received, persisting queue...\n", <-sigs)
		if err := appQueue.PersistToJsonFile(); err != nil {
			fmt.Println("Error persisting queue:", err)
		}
		os.Exit(0)
	}()

	err = server.ListenAndServe()
	if err != nil {
		if err := appQueue.PersistToJsonFile(); err != nil {
			fmt.Println("Error persisting queue:", err)
		}
		panic(err)
	}

}
