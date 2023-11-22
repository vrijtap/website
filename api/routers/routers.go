package routers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Middleware for implementing security headers
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        next.ServeHTTP(w, r)
    })
}

// HTTPRouter configures the HTTP router for redirecting to HTTPS.
func HTTPRouter(portHTTPS string) http.Handler {
	r := mux.NewRouter()

	// Add security middleware
	r.Use(SecurityHeadersMiddleware)

	// Serve static files from the "web/static" directory.
	fs := http.FileServer(http.Dir("web/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Configure client, owner, order and payment routes.
	ConfigureClientRoutes(r)
	ConfigureOwnerRoutes(r)
	ConfigureOrderRoutes(r)

	return r
}

// HTTPSRouter configures the HTTPS router for serving static files and application routes.
func HTTPSRouter() *mux.Router {
	r := mux.NewRouter()

	// Add security middleware
	r.Use(SecurityHeadersMiddleware)

	// Serve static files from the "web/static" directory.
	fs := http.FileServer(http.Dir("web/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Configure client, owner, order and payment routes.
	ConfigureClientRoutes(r)
	ConfigureOwnerRoutes(r)
	ConfigureOrderRoutes(r)

	return r
}
