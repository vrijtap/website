package routers

import (
	"backend/api/handlers"
	"backend/internal/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

// ConfigureOwnerRoutes sets up owner-related routes on a provided Gorilla Mux router.
func ConfigureOwnerRoutes(router *mux.Router) {
	// Create a subrouter for owner-related routes under the "/owner" path.
	ownerRouter := router.PathPrefix("/owner").Subrouter()

	ownerRouter.Use(middleware.AuthenticationMiddleware)

	// Handle GET requests to the "/owner" endpoint with the OwnerGET handler function.
	ownerRouter.HandleFunc("", handlers.OwnerGet).Methods(http.MethodGet)

	// Handle POST requests to the "/owner" endpoint with the OwnerLogin handler function.
	ownerRouter.HandleFunc("", handlers.OwnerLogin).Methods(http.MethodPost)

	// Handle PUT requests to the "/owner" endpoint with the OwnerPUT handler function.
	ownerRouter.HandleFunc("", handlers.OwnerPut).Methods(http.MethodPut)
}
