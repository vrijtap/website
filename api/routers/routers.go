package routers

import (
	"net/http"
	"github.com/gorilla/mux"
	"website/internal/middleware"
)

// CreateAPIRouter creates a router that routes to API endpoints
func CreateAPIRouter() http.Handler {
	router := mux.NewRouter()

	// Add middleware functionalities
	router.Use(middleware.SecurityHeadersMiddleware)

	// Serve static files from the "web/static" directory.
	fileServer := http.FileServer(http.Dir("web/static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	// Configure client, owner, order, and payment routes.
	ConfigureClientRoutes(router)
	ConfigureOwnerRoutes(router)
	ConfigureOrderRoutes(router)

	return router
}
