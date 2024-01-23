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

	// Define routes for owner-related endpoints
	ownerRouter.HandleFunc("", handlers.OwnerGet).Methods(http.MethodGet)
	ownerRouter.HandleFunc("", handlers.OwnerLogin).Methods(http.MethodPost)
	ownerRouter.HandleFunc("", handlers.OwnerPut).Methods(http.MethodPut)
}
