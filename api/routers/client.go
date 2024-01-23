package routers

import (
    "website/api/handlers"

	"net/http"

	"github.com/gorilla/mux"
)

// ConfigureClientRoutes sets up client-related routes on a provided Gorilla Mux router
func ConfigureClientRoutes(router *mux.Router) {
    // Create a subrouter for client-related routes under the "/client" path
    clientRouter := router.PathPrefix("/client").Subrouter()

    // Define routes for client-related endpoints
	clientRouter.HandleFunc("/{server_id}", handlers.ClientGet).Methods(http.MethodGet)
}
