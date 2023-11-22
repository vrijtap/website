package server

import (
	"backend/api/routers"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

// NewHTTPServer creates and configures an HTTP server.
func NewHTTPServer(portHTTP, portHTTPS string) *http.Server {
	// Return the server
	return &http.Server{
		Addr:    fmt.Sprintf(`:%s`, portHTTP),
		Handler: routers.HTTPRouter(portHTTPS),
	}
}

// NewHTTPSServer creates and configures an HTTPS server.
func NewHTTPSServer(portHTTPS, certFilePath, keyFilePath string) *http.Server {
	// Load in the TLS certificates
	tlsCert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil {
		log.Fatalf("Error loading certificate and key file: %v", err)
	}

	// Create the TLS configuration
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}

	// Return the server
	return &http.Server{
		Addr:      fmt.Sprintf(`:%s`, portHTTPS),
		Handler:   routers.HTTPSRouter(),
		TLSConfig: tlsConfig,
	}
}
