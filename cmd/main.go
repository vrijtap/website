package main

import (
	"website/internal/app"
	"website/internal/server"

	"log"
	"os"
	"os/signal"
	"syscall"
	"sync"
	"net/http"
)

func main() {
	// Initialize the application
    if err := app.Initialize("./"); err != nil {
        log.Fatalf("[Error] %v", err)
    }

	// Start the API based on the environment
	var apiServer *http.Server

	// Create a channel to receive interrupt signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Create a WaitGroup to synchronize goroutines
	var wg sync.WaitGroup

	// Goroutine to handle interrupt signal
	go func() {
		// Wait for the interrupt signal
		<-interrupt

		// Increment WaitGroup counter to indicate the start of this goroutine
		wg.Add(1)

		// Clean up resources and gracefully exit
		if err := app.Clean(apiServer); err != nil {
			log.Printf("[Warning] %v", err)
		}

		// Exit the program
		os.Exit(0)
	}()

	// Start the server
	if os.Getenv("ENVIRONMENT") == "production" {
		// Create HTTPS server
		apiServer = server.CreateHTTPSServer(
			os.Getenv("PORT_HTTPS"),
			os.Getenv("PATH_CERT_FILE"),
			os.Getenv("PATH_KEY_FILE"),
		)

		// Start HTTPS server
		log.Printf("Listening to port %s for HTTPS requests...\n", os.Getenv("PORT_HTTPS"))
		if err := apiServer.ListenAndServeTLS("", ""); err != nil {
			log.Printf("%v", err)
		}
	} else {
		// Create HTTP server
		apiServer = server.CreateHTTPServer(
			os.Getenv("PORT_HTTP"),
		)

		// Start HTTP server
		log.Printf("Listening to port %s for HTTP requests...\n", os.Getenv("PORT_HTTP"))
		if err := apiServer.ListenAndServe(); err != nil {
			log.Printf("%v", err)
		}
	}

	// Wait for priority and exit
	wg.Wait()
	if err := app.Clean(apiServer); err != nil {
		log.Printf("[Warning] %v", err)
	}
	os.Exit(1)
}
