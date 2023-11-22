package password

import (
	"errors"
	"os"
	"regexp"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type localConfig struct {
	filename string
	mu       sync.Mutex
}

var config = localConfig{
	filename: "",
	mu:       sync.Mutex{},
}

// hash creates a secure hash for a given password
func hash(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// save stores the hash in a file
func save(filename, hash string) error {
	if err := os.WriteFile(filename, []byte(hash), 0600); err != nil {
		return err
	}
	return nil
}

// read reads the hash from a file
func read(filename string) (string, error) {
	hash, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ChangeTo sets the password to its new value
func ChangeTo(password string) error {
	config.mu.Lock()
	defer config.mu.Unlock()

	// Hash the received password
	hash, err := hash(password)
	if err != nil {
		return err;
	}

	// Change all instances of the password to the new hash
	os.Setenv("PASSWORD", hash)
	if err := save(config.filename, hash); err != nil {
		return err;
	}

	return nil
}


// Init loads in the owner password
func Init(filename string) error {
	config.filename = filename

	// Try and read from an existing password file
	if hash, err := read(filename); err == nil {
		os.Setenv("PASSWORD", hash)
	} else {
		// Check if the standard password is defined
		stdPwd := os.Getenv("PASSWORD_DEFAULT")
		if stdPwd == "" {
			return errors.New("PASSWORD_DEFAULT environment variable is undeclared")
		}

		// Implement the standard password
		if err := ChangeTo(stdPwd); err != nil {
			return err
		}
	}

	return nil
}

// Check matches the stored password hash with the given password
func Check(password string) bool {
	config.mu.Lock()
	defer config.mu.Unlock()

	err := bcrypt.CompareHashAndPassword([]byte(os.Getenv("PASSWORD")), []byte(password))
	if err != nil {
		return false
	}
	return true
}

// Validate checks if the password matches requirements
func Validate(password string) []error {
	var errorsList []error

	// Define the password matching rules
	rules := []struct {
		pattern *regexp.Regexp
		message string
	}{
		{regexp.MustCompile(`[A-Z]`), "Password must contain at least one uppercase letter"},
		{regexp.MustCompile(`[a-z]`), "Password must contain at least one lowercase letter"},
		{regexp.MustCompile(`\d`), "Password must contain at least one digit"},
		{regexp.MustCompile(`[@#$%^&+=!(){}[\]\*\-_\\|;:'",<.>/?]`), "Password must contain at least one special character"},
	}

	// Check if the rules are followed
	for _, rule := range rules {
		if !rule.pattern.MatchString(password) {
			errorsList = append(errorsList, errors.New(rule.message))
		}
	}

	// Check if the password is long enough
	if len(password) < 8 {
		errorsList = append(errorsList, errors.New("Password must be at least 8 characters long"))
	}

	return errorsList
}
