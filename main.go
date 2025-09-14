package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nicograef/qugo/qugo"
)

func main() {
	server := &http.Server{
		Addr: ":3000",
	}

	// Get queue size and consumer URL from env
	queueSize := 1000
	if v := os.Getenv("QUEUE_SIZE"); v != "" {
		fmt.Sscanf(v, "%d", &queueSize)
	}
	consumerURL := os.Getenv("CONSUMER_URL")
	if consumerURL == "" {
		consumerURL = "http://localhost:4000"
	}

	appQueue := make(chan qugo.Message, queueSize)

	http.HandleFunc("/", qugo.NewEnqueueHandler(appQueue))

	// Consumer goroutine: reads from channel and calls webhook
	go func() {
		for msg := range appQueue {
			// Send message to webhook
			resp, err := qugo.SendToWebhook(consumerURL, msg)
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
		fmt.Printf("Signal %s received, shutting down...\n", <-sigs)
		os.Exit(0)
	}()

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
