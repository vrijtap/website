package fakepay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Create a struct with the same structure as TransactionInput in the fakepay api
type FakepayTransactionInput struct {
    Amount      float64 `json:"amount"`
    WebhookURL  string  `json:"webhook_url"`
    WebhookKey  string  `json:"webhook_key"`
    RedirectURL string  `json:"redirect_url"`
}

// FakepayProcessor prepares and executes a transaction,
// then modifies the URL and returns it as the redirection URL.
func FakepayProcessor(input FakepayTransactionInput) (string, error) {
	// Convert transaction input to JSON
	transactionData, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	// Create a new POST request to CreateTransaction endpoint on the transaction server
	transactionURL := os.Getenv("PAYMENT_GATE_URL")
	req, err := http.NewRequest("POST", transactionURL, bytes.NewBuffer(transactionData))
	if err != nil {
		return "", err
	}

	// Set the Content-Type header for the request
	req.Header.Set("Content-Type", "application/json")

	// Set the Authorization header for the request
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("PAYMENT_GATE_KEY")))

	// Make the POST request with the prepared request object
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("Failed to create transaction: %d", response.StatusCode)
	}

	// Read the response body to get the URL
	var responseMap map[string]string
	if err := json.NewDecoder(response.Body).Decode(&responseMap); err != nil {
		return "", err
	}

	// Get the redirection URL
	transactionURL, exists := responseMap["url"]
	if !exists {
		return "", fmt.Errorf("Invalid transaction response")
	}

	return transactionURL, nil
}
