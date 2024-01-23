package handlers

import (
	"website/utils/database/models/cards"
	"website/web/templates"
	
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// ClientGet handles GET requests to the client page
func ClientGet(w http.ResponseWriter, r *http.Request) {
	// Extract card_id variable from the URL path parameters.
	vars := mux.Vars(r)
	serverID := vars["server_id"]

	// Parse the serverID
	id, err := strconv.ParseUint(serverID, 10, 64)
	if err != nil {
		http.Error(w, "Supplied wrong id", http.StatusBadRequest)
		return
	}

	// Retrieve card information from the database
	card, err := cards.GetByServerID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to retrieve card", http.StatusInternalServerError)
		return
	}
	
	// Convert the price to float for formatting
	price, err := strconv.ParseFloat(os.Getenv("PRICE"), 64)

	// Setup the client page variables
	data := struct {
		Name    string
		Price   string
		Beers   uint
		ID		uint
	}{
		Name:  	os.Getenv("NAME"),
		Price: 	fmt.Sprintf("%.2f", price),
		Beers: 	card.Beers,
		ID:		uint(card.ServerID),
	}

	// Set the Content-Type header to specify that the response is HTML
	w.Header().Set("Content-Type", "text/html")

	// Render the client page
	err = templates.RenderHTML(w, "client.html", data)
	if err != nil {
		errMsg := "Failed to render HTML template"
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
}
