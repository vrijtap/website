package main

import (
	"context"
	"log"
	"os"
	"backend/internal/environment"
	"backend/internal/password"
	"backend/internal/server"
	"backend/utils/database"
	"backend/utils/database/models/cards"
	"backend/web/templates"
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
	// Load configurations
	if err := environment.Init(`./.env`); err != nil {
		log.Fatalf("Error loading environment: %v", err)
	}

	// Load passwords
	if err := password.Init(`./password.env`); err != nil {
		log.Fatalf("Error loading password: %v", err)
	}

	// Load templates
	templates.LoadTemplates()

	// Initialize MongoDB
	if err := database.Connect(os.Getenv("MONGO_URI")); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Close database on dereference
	defer func() {
		if err := database.Close(); err != nil {
			log.Fatalf("Error closing MongoDB connection: %v", err)
		}
	}()

	// Create the admin card
	initAdminCard()

	// Start the API based on the environment
	if os.Getenv("ENVIRONMENT") == "production" {
		// Create HTTPS server
		httpsServer := server.NewHTTPSServer(
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
		httpServer := server.NewHTTPServer(
			os.Getenv("PORT_HTTP"),
			os.Getenv("PORT_HTTPS"),
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
