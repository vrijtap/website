package routers

import (
	"net/http"
	"github.com/gorilla/mux"
	"website/api/handlers"
)

// ConfigureClientRoutes sets up client-related routes
func ConfigureClientRoutes(r *mux.Router) {
    // Create a subrouter for client-related routes under the "/client" path.
    clientRouter := r.PathPrefix("/client").Subrouter()

    // Handle GET requests to the "/client/{card_id}" endpoint with the ClientGet handler function.
    clientRouter.HandleFunc("/{server_id}", handlers.ClientGet).Methods(http.MethodGet)
}