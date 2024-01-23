package app

import (
	"website/internal/password"
	"website/web/templates"
	"website/utils/database"
	"website/utils/database/models/cards"

	"fmt"
	"os"
	"context"
	"net/http"
	"log"
	"time"

	"github.com/joho/godotenv"
)

// initAdminCard initializes an admin (testing) card for the backend if it doesn't already exist.
func initAdminCard() error {
	// Check if there is an admin card already
	card, err := cards.GetByServerID(context.TODO(), 0)
	if card == nil && err != nil {
		// Initialize the admin card
		card, err := cards.New(context.TODO())
		if err != nil {
			return fmt.Errorf("failed to create a new admin card: %v", err)
		}

		// Insert the admin card into the database
		card.ServerID = 0
		err = cards.Insert(context.TODO(), &card)
		if err != nil {
			return fmt.Errorf("failed to insert the admin card into the database: %v", err)
		}
	}

	return nil
}

// Initialize initializes the application
func Initialize(relativeRootFolder string) error {
	// Load configurations from .env file
    if err := godotenv.Load(fmt.Sprintf("%s.env", relativeRootFolder)); err != nil {
        return fmt.Errorf("failed to load environment configurations from .env file: %v", err)
    }

	// Load passwords
	if err := password.Init(`./password.env`); err != nil {
		return fmt.Errorf("failed to initialize the owner password: %v", err)
	}

    // Load templates from the templates folder
    if err := templates.Load(fmt.Sprintf("%sweb/templates/", relativeRootFolder)); err != nil {
        return fmt.Errorf("failed to load .html templates: %v", err)
    }

	// Initialize the database connection
    if err := database.Connect(os.Getenv("MONGO_URI"), "backend"); err != nil {
        return fmt.Errorf("unable to establish connection to the database: %v", err)
    }

	// Setup the admin card if needed
	if err := initAdminCard(); err != nil {
		return err
	}

    return nil
}

// Clean is a function that performs cleanup operations, closing the server and disconnecting from the database.
func Clean(server *http.Server) error {
    log.Println("Shutting down gracefully...")
    var errs []error

    // Attempt to close the HTTP server
    ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()
    if err := server.Shutdown(ctx); err != nil {
        errs = append(errs, fmt.Errorf("unable to shutdown the server: %v", err))
    }

    // Attempt to disconnect from the database
    if err := database.Disconnect(); err != nil {
        errs = append(errs, fmt.Errorf("unable to disconnect the database: %v", err))
    }

	// No errors, return nil
    if len(errs) == 0 {
        return nil
    }

    // If there are errors, return them as a multi-error
    return errs[0]
}
