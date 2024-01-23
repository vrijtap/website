package routers

import (
	"net/http"
	"github.com/gorilla/mux"
	"website/internal/middleware"
)

// HTTPRouter configures the HTTP router for redirecting to HTTPS.
func HTTPRouter(portHTTPS string) http.Handler {
	router := mux.NewRouter()

	// Add middleware functionalities
	router.Use(middleware.SecurityHeadersMiddleware)

	// Serve static files from the "web/static" directory.
	fileServer := http.FileServer(http.Dir("web/static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	// Configure client, owner, order and payment routes.
	ConfigureClientRoutes(router)
	ConfigureOwnerRoutes(router)
	ConfigureOrderRoutes(router)

	return router
}

// HTTPSRouter configures the HTTPS router for serving static files and application routes.
func HTTPSRouter() *mux.Router {
	router := mux.NewRouter()

	// Add middleware functionalities
	router.Use(middleware.SecurityHeadersMiddleware)
	router.Use(middleware.AuthenticationMiddleware)

	// Serve static files from the "web/static" directory.
	fileServer := http.FileServer(http.Dir("web/static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	// Configure client, owner, order and payment routes.
	ConfigureClientRoutes(router)
	ConfigureOwnerRoutes(router)
	ConfigureOrderRoutes(router)

	return router
}
