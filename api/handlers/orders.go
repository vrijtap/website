package handlers

import (
	"website/internal/payment/fakepay"
	"website/utils/database/models/cards"
	"website/utils/database/models/orders"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PaymentData represents the JSON data structure for payment requests.
type PaymentData struct {
	Quantity string `json:"quantity"`
	ID       string `json:"id"`
}

// Function to determine the correct webhook URL based on the request scheme (HTTP or HTTPS)
func getWebhookURL(r *http.Request, orderID string) string {
    scheme := "https"
    if r.TLS == nil {
        // Request is not over HTTPS, use HTTP instead
        scheme = "http"
    }
    return fmt.Sprintf("%s://%s/order/%s", scheme, r.Host, orderID)
}

// Function to determine the correct redirect URL based on the request scheme (HTTP or HTTPS)
func getRedirectURL(r *http.Request, clientID string) string {
    scheme := "https"
    if r.TLS == nil {
        // Request is not over HTTPS, use HTTP instead
        scheme = "http"
    }
    return fmt.Sprintf("%s://%s/client/%s", scheme, r.Host, clientID)
}

// OrderPost handles POST requests for creating orders
func OrderPost(w http.ResponseWriter, r *http.Request) {
	// Parse JSON data from the request body into paymentData struct
	var paymentData PaymentData
	if err := json.NewDecoder(r.Body).Decode(&paymentData); err != nil {
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	// Parse server ID from payment data
	id, err := strconv.ParseUint(paymentData.ID, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid server ID: %v", err), http.StatusBadRequest)
		return
	}

	// Fetch card details by server ID
	card, err := cards.GetByServerID(r.Context(), id)
	if err != nil {
		http.Error(w, "Could not fetch Card", http.StatusBadRequest)
		return
	}

	// Parse quantity from payment data
	quantity, err := strconv.ParseUint(paymentData.Quantity, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid quantity: %v", err), http.StatusBadRequest)
		return
	}

	// Parse price from environment variable
	price, err := strconv.ParseFloat(os.Getenv("PRICE"), 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid price: %v", err), http.StatusInternalServerError)
		return
	}

	// Create a new order with parsed data
	order := orders.New(card.ID, uint(quantity), price)
	_, err = orders.Insert(r.Context(), &order)
	if err != nil {
		http.Error(w, "Could not create an order", http.StatusInternalServerError)
		return
	}

	// Prepare redirection URL for the transaction
	var redirectURL string
	if os.Getenv("ENVIRONMENT") == "production" {
		redirectURL = "" // Actual payment gate space
	} else {
		// Prepare data for the transaction using fakepay processor
		transactionInput := fakepay.FakepayTransactionInput{
    		Amount:			order.TotalAmount,
    		WebhookURL:		getWebhookURL(r, order.ID.Hex()),
    		WebhookKey:		os.Getenv("WEBHOOK_KEY"),
    		RedirectURL:	getRedirectURL(r, paymentData.ID),
		}

		// Obtain the modified redirect URL from the fakepay processor
		redirectURL, err = fakepay.FakepayProcessor(transactionInput)
		if err != nil {
			http.Error(w, "Could not create a transaction", http.StatusInternalServerError)
			return
		}

		// Extract the host part without the port from r.Host
		hostParts := strings.Split(r.Host, ":")
		host := hostParts[0]

		// Modify the URL to replace "localhost" with the actual host
		redirectURL = strings.Replace(redirectURL, "localhost", host, -1)
	}

	// Return the URL in the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"url": redirectURL})
}

// OrderUpdateStatus handles POST requests for updating order statuses.
func OrderUpdateStatus(w http.ResponseWriter, r *http.Request) {
	// Check Authorization header for valid credentials
	if auth := r.Header.Get("Authorization"); auth != fmt.Sprintf("Bearer %s", os.Getenv("WEBHOOK_KEY")) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the order_id variable from the URL path parameters.
	vars := mux.Vars(r)
	orderID := vars["order_id"]

	// Convert hexadecimal string to primitive.ObjectID.
	objectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusInternalServerError)
		return
	}

	// Get the order details from the database using the order ID.
	order, err := orders.GetByID(r.Context(), objectID)
	if err != nil {
		http.Error(w, "Could not fetch order", http.StatusInternalServerError)
		return
	}

	// Check if the order was already handled.
	if order.Status != "Pending" {
		http.Error(w, "Order was already processed", http.StatusOK)
		return
	}

	// Parse the form data
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Access form fields by name
	status := r.FormValue("Status")

	// Update the order status.
	order.Status = status

	// Get the card details from the database using the card ID.
	card, err := cards.GetByID(r.Context(), order.CardID)
	if err != nil {
		http.Error(w, "Could not fetch card", http.StatusInternalServerError)
		return
	}

	// Define the card updates.
	updates := bson.M{
		"beers":         card.Beers + order.Quantity,
		"last_purchase": time.Now(),
	}

	// Update the card.
	if err := cards.UpdateByID(r.Context(), card.ID, updates); err != nil {
		http.Error(w, "Failed to update card", http.StatusInternalServerError)
		return
	}

	// Update the order status in the database.
	if err := orders.UpdateStatus(r.Context(), objectID, order.Status); err != nil {
		http.Error(w, "Failed to update status", http.StatusInternalServerError)
		return
	}
}
