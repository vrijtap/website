package app

import (
	"website/web/templates"

	"fmt"
)

// Initialize initializes the application
func Initialize(relativeRootFolder string) error {
    // Load templates from the templates folder
    if err := templates.Load(fmt.Sprintf("%sweb/templates/", relativeRootFolder)); err != nil {
        return fmt.Errorf("failed to load .html templates: %v", err)
    }

    return nil
}
