package app

import (
	"website/internal/password"
	"website/web/templates"
	"website/utils/database"

	"fmt"
	"os"

	"github.com/joho/godotenv"
)

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

    return nil
}
