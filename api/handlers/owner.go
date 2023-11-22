package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"backend/internal/password"
	"backend/utils/sessions"
	"backend/web/templates"
)

type ErrorResponse struct {
    Errors  []string `json:"errors"`
    Message string   `json:"message"`
}

// OwnerGet handles GET requests meant for viewing the statistics page
func OwnerGet(w http.ResponseWriter, r *http.Request) {
	// Get the owner session
	session, err := sessions.Get(w, r, "owner-auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the owner is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !auth || !ok {
		// Render the login page if the user is not authenticated
		templates.RenderTemplate(w, "login.html", struct {
			Method string
			Action string
		}{"POST", "Enter Password"})
	} else if check, ok := session.Values["changePassword"].(bool); check && ok {
		// Render the update password page if the flag is set
		templates.RenderTemplate(w, "login.html", struct {
			Method string
			Action string
		}{"PUT", "Change Password"})
	} else {
		// Render the owner page if everything is correct
		templates.RenderTemplate(w, "owner.html", struct {
			Name string
		}{"WIP"})
	}
}

// OwnerPost handles POST requests meant for gaining auth
func OwnerPost(w http.ResponseWriter, r *http.Request) {
	// Check if the password is present
	pwd := r.FormValue("password")
	if pwd == "" {
		http.Error(w, "Password is missing", http.StatusBadRequest)
		return
	}

	// Get the owner session
	session, err := sessions.Get(w, r, "owner-auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the password is correct
	if password.Check(pwd) {
		// Set the session variables
		session.Values["authenticated"] = true
		if pwd == os.Getenv("PASSWORD_DEFAULT") {
			session.Values["changePassword"] = true
		}

		// Save the session
		err := session.Save(r, w)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}
}

// validatePassword validates if a provided password is up to standards
func validatePassword(w http.ResponseWriter, pwd string) []error {
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

func OwnerPut(w http.ResponseWriter, r *http.Request) {
	// Get the owner session
	session, err := sessions.Get(w, r, "owner-auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the owner is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !auth || !ok {
		http.Error(w, "Session is not valid or authenticated", http.StatusUnauthorized)
		return
	}

	// Get the password from the form
	pwd := r.FormValue("password")

	// Validate the input
	if err := validatePassword(w, pwd); err != nil {
		return
	}

	// Change the password
	err = password.ChangeTo(pwd)
	if err != nil {
		http.Error(w, "Failed to save password", http.StatusInternalServerError)
		return
	}

	// Update the user session
	session.Values["changePassword"] = false
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
