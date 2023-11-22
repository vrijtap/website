package routers

import (
    "net/http"
    "github.com/gorilla/mux"
    "backend/api/handlers"
)

// ConfigureOrderRoutes configures the routes related to orders.
func ConfigureOrderRoutes(r *mux.Router) {
    orderRouter := r.PathPrefix("/order").Subrouter()

    // Handle POST requests for creating orders.
    orderRouter.HandleFunc("", handlers.OrderPost).Methods(http.MethodPost)

    // Handle PUT requests for updating order statuses.
    orderRouter.HandleFunc("/{order_id}", handlers.OrderUpdateStatus).Methods(http.MethodPost)
}