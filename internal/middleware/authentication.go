package middleware

import (
	"backend/utils/jwt"
	"fmt"
    "strings"
	"net/http"
)

// AuthenticationMiddleware checks if the user has a valid JWT token
func AuthenticationMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try and get the token cookie or Authorization header
        cookie, err := r.Cookie("token")
        token := ""
        if err == nil {
            token = cookie.Value
	    } else {
            // Get the "Authorization" header from the request
	        authHeader := r.Header.Get("Authorization")

	        // Check if the header is missing or has an unexpected format
	        if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
                token = strings.TrimPrefix(authHeader, "Bearer ")
	        }
        }

        // Verify that the token was signed correctly
        err = jwt.VerifyToken(token)
        if err != nil {
            r.Header.Set("Authorization", "")
        } else {
            r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
        }

		// Pass on the Request
		next.ServeHTTP(w, r)
    })
}
