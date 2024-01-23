package main

import (
	"context"
	"log"
	"os"
	"website/internal/app"
	"website/internal/server"
	"website/utils/database"
	"website/utils/database/models/cards"
)

// initAdminCard initializes an admin (testing) card for the backend if it doesn't already exist.
func initAdminCard() {
	card, err := cards.GetByServerID(context.TODO(), 0)
	if card == nil && err != nil {
		card, err := cards.New(context.TODO())
		if err != nil {
			log.Fatalf("Error creating test card: %v", err)
		}
		card.ServerID = 0
		cards.Insert(context.TODO(), &card)
	}
}

func main() {
	// Initialize the application
    if err := app.Initialize("./"); err != nil {
        log.Fatalf("[Error] %v", err)
    }

	// Close database on dereference
	defer func() {
		if err := database.Disconnect(); err != nil {
			log.Fatalf("Error closing MongoDB connection: %v", err)
		}
	}()

	// Create the admin card
	initAdminCard()

	// Start the API based on the environment
	if os.Getenv("ENVIRONMENT") == "production" {
		// Create HTTPS server
		httpsServer := server.CreateHTTPSServer(
			os.Getenv("PORT_HTTPS"),
			os.Getenv("PATH_CERT_FILE"),
			os.Getenv("PATH_KEY_FILE"),
		)

		// Close server on dereference
		defer func() {
			if err := httpsServer.Close(); err != nil {
				log.Fatalf("Error closing HTTPS server: %v", err)
			}
		}()

		// Start HTTPS server
		log.Printf("Listening to port %s for HTTPS requests...\n", os.Getenv("PORT_HTTPS"))
		if err := httpsServer.ListenAndServeTLS("", ""); err != nil {
			log.Printf("%v", err)
		}
	} else {
		// Create HTTP server
		httpServer := server.CreateHTTPServer(
			os.Getenv("PORT_HTTP"),
		)

		// Close server on dereference
		defer func() {
			if err := httpServer.Close(); err != nil {
				log.Fatalf("Error closing HTTP server: %v", err)
			}
		}()

		// Start HTTP server
		log.Printf("Listening to port %s for HTTP requests...\n", os.Getenv("PORT_HTTP"))
		if err := httpServer.ListenAndServe(); err != nil {
			log.Printf("%v", err)
		}
	}
}
