package handlers

import (
	"website/internal/password"
	"website/internal/jwt"
	"website/web/templates"
	
	"encoding/json"
	"net/http"
	"os"
	"sync"
)

// Globals for changing the password
var (
	passwordMutex   sync.Mutex
	changePassword  = false
)

// ErrorResponse represents the structure for reporting password update errors
type ErrorResponse struct {
	Errors  []string `json:"errors"`
	Message string   `json:"message"`
}

// validatePassword validates if a provided password is up to standards
func validatePassword(w http.ResponseWriter, pwd string) []error {
	// Validate the provided password
	validationErrors := password.Validate(pwd)

	if validationErrors != nil {
		// Prepare the error response with populated error messages
		errorResponse := ErrorResponse{
			Errors:  make([]string, len(validationErrors)),
			Message: "Update the password according to these errors",
		}

		// Populate the error messages directly in the struct
		for i, err := range validationErrors {
			errorResponse.Errors[i] = err.Error()
		}

		// Marshal the struct to JSON and send the response
		jsonResponse, err := json.Marshal(errorResponse)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return validationErrors
		}

		// Write the response in a single step
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return validationErrors
	}

	return nil
}

// OwnerGet handles GET requests meant for viewing the statistics page
func OwnerGet(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	var data interface{}
	var page string

	// Check the authentication
	if authHeader != "" {
		passwordMutex.Lock()
		if changePassword == true {
			// Setup the login page variables
			data = struct {
				Method string
				Action string
			}{
				"PUT",
				"Change Password",
			}
			page = "login.html"
		} else {
			// Setup the owner page variables
			data = struct {
				Name              string
				RaspberryEndpoint string
			}{
				os.Getenv("NAME"),
				os.Getenv("RASPBERRY_ENDPOINT"),
			}
			page = "owner.html"
		}
		passwordMutex.Unlock()
	} else {
		// Setup the login page variables
		data = struct {
			Method string
			Action string
		}{
			"POST",
			"Enter Password",
		}
		page = "login.html"
	}

	// Set the Content-Type header to specify that the response is HTML
	w.Header().Set("Content-Type", "text/html")

	// Render the page
	err := templates.RenderHTML(w, page, data)
	if err != nil {
		errMsg := "Failed to render HTML template"
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
}

// OwnerLogin handles POST requests meant for gaining auth
func OwnerLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the passed password
	pwd := r.FormValue("password")
	if pwd == "" {
		http.Error(w, "Password is missing", http.StatusBadRequest)
		return
	}

	// Check if the password is correct
	if password.Check(pwd) {
		// Generate a JWT token
		tokenString, err := jwt.CreateToken()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// Set the JWT token as a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			SameSite: http.SameSiteStrictMode,
			Secure:   false,
		})

		// If the default password was used, request a change
		if pwd == os.Getenv("PASSWORD_DEFAULT") {
			passwordMutex.Lock()
			changePassword = true
			passwordMutex.Unlock()
		}

		// Write back the auth token
		w.Header().Set("Authorization", "Bearer "+tokenString)
		w.WriteHeader(http.StatusOK)
		return
	} else {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}
}

// OwnerPut handles PUT requests meant for changing the password
func OwnerPut(w http.ResponseWriter, r *http.Request) {
	// Check the authentication
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the password from the form
	pwd := r.FormValue("password")
	if err := validatePassword(w, pwd); err != nil {
		return
	}

	// Change the password
	err := password.ChangeTo(pwd)
	if err != nil {
		http.Error(w, "Failed to save password", http.StatusInternalServerError)
		return
	}

	// Indicate that the password was changed
	passwordMutex.Lock()
	changePassword = false
	passwordMutex.Unlock()
}
