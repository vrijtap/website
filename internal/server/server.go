package server

import (
	"website/api/routers"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

// CreateHTTPServer creates and configures an HTTP server.
func CreateHTTPServer(portHTTP string) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", portHTTP),
		Handler: routers.CreateApiRouter(),
	}
}

// CreateHTTPSServer creates and configures an HTTPS server with TLS certificates.
func CreateHTTPSServer(portHTTPS, certFilePath, keyFilePath string) *http.Server {
	// Load TLS certificates
	tlsCert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil {
		log.Fatalf("Error loading certificate and key file: %v", err)
	}

	// Configure TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}

	return &http.Server{
		Addr:      fmt.Sprintf(":%s", portHTTPS),
		Handler:   routers.CreateApiRouter(),
		TLSConfig: tlsConfig,
	}
}
