package sessions

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// Store is the session store.
var Store *sessions.CookieStore

// Initialize initializes the session store with a server secret.
func Initialize(secret string) {
	Store = sessions.NewCookieStore([]byte(secret))
}

// Get retrieves the session for the current request, and updates cookie attributes.
func Get(w http.ResponseWriter, r *http.Request, name string) (*sessions.Session, error) {
	session, err := Store.Get(r, name)

	// Check if the session cookie already exists
	if err == nil {
		// Set SameSite, Secure, and HttpOnly attributes if the cookie exists
		session.Options.MaxAge = 3600
		session.Options.SameSite = http.SameSiteStrictMode
		session.Options.Secure = true
		session.Options.HttpOnly = true
	}

	return session, err
}
