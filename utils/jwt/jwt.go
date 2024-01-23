package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken() (string, error) {
	// Create a new JWT token expiring in an hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	
	// Sign the new JWT token
    tokenString, err := token.SignedString([]byte(os.Getenv("SERVER_SECRET")))
    if err != nil {
    	return "", err
    }

	// Return the token string
	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SERVER_SECRET")), nil
	})

	// Return parsing errors
	if err != nil {
		return err
	}

	// Check if the token is valid
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	// Access the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("Error extracting claims")
	}

	// Access the expiration claim
	expiration, ok := claims["exp"].(float64)
	if !ok {
		return fmt.Errorf("Error extracting expiration claim")
	}

	// Convert the expiration time to a Go time.Time object
	expirationTime := time.Unix(int64(expiration), 0)

	// Check if the token has expired
	if time.Now().After(expirationTime) {
		return fmt.Errorf("Token has expired")
	} else {
		return nil
	}
}
