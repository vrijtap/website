package routers

import (
	"net/http"
	"github.com/gorilla/mux"
	"backend/api/handlers"
)

// ConfigureOwnerRoutes sets up owner-related routes on a provided Gorilla Mux router.
func ConfigureOwnerRoutes(r *mux.Router) {
	// Create a subrouter for owner-related routes under the "/owner" path.
	ownerRouter := r.PathPrefix("/owner").Subrouter()

	// Handle GET requests to the "/owner" endpoint with the OwnerGET handler function.
	ownerRouter.HandleFunc("", handlers.OwnerGet).Methods(http.MethodGet)

	// Handle POST requests to the "/owner" endpoint with the OwnerPOST handler function.
	ownerRouter.HandleFunc("", handlers.OwnerPost).Methods(http.MethodPost)

	// Handle PUT requests to the "/owner" endpoint with the OwnerPUT handler function.
	ownerRouter.HandleFunc("", handlers.OwnerPut).Methods(http.MethodPut)
}
