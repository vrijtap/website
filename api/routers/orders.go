package routers

import (
    "website/api/handlers"

    "net/http"

    "github.com/gorilla/mux"
)

// ConfigureOrderRoutes sets up order-related routes on a provided Gorilla Mux router
func ConfigureOrderRoutes(router *mux.Router) {
    // Create a subrouter for order-related routes under the "/order" path
    orderRouter := router.PathPrefix("/order").Subrouter()

    // Define routes for order-related endpoints
    orderRouter.HandleFunc("", handlers.OrderPost).Methods(http.MethodPost)
    orderRouter.HandleFunc("/{order_id}", handlers.OrderUpdateStatus).Methods(http.MethodPost)
}
