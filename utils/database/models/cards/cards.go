package cards

import (
	"context"
	"time"
	"website/utils/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Card represents data for an RFID card
type Card struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ServerID	 uint64				`bson:"server_id"`
	Beers        uint               `bson:"beers"`
	LastPurchase time.Time          `bson:"last_purchase"`
}

// findHighestQR finds the highest QR value in the "cards" collection in MongoDB
func findHighestServerID(ctx context.Context) (uint64, error) {
	collection := database.GetCollection("cards")

	// Set sorting options to sort by "server_id" field in descending order (value: -1)
	findOptions := options.FindOne().SetSort(bson.D{{Key: "server_id", Value: -1}})

	// Find the document with the highest "server_id" value and decode it into the 'card' variable
	var card Card
	err := collection.FindOne(ctx, bson.M{}, findOptions).Decode(&card)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err
	}
	return card.ServerID, nil
}

// New creates a new Card instance with default values and a unique QR value
func New(ctx context.Context) (Card, error) {
	// Find the highest QR value in the database
	highest, err := findHighestServerID(ctx)
	if err != nil {
		return Card{}, err
	}

	// Increment the highest QR value to generate a unique QR for the new card
	return Card{
		ServerID:     highest + 1,
		Beers:        0,
		LastPurchase: time.Time{},
	}, nil
}

// GetByServerID retrieves a card document from MongoDB by its Server ID
func GetByServerID(ctx context.Context, ServerID uint64) (*Card, error) {
	collection := database.GetCollection("cards")
	filter := bson.M{"server_id": ServerID}
	var card Card
	err := collection.FindOne(ctx, filter).Decode(&card)
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// GetByID retrieves a card document from MongoDB by its ObjectID
func GetByID(ctx context.Context, cardID primitive.ObjectID) (*Card, error) {
	collection := database.GetCollection("cards")
	filter := bson.M{"_id": cardID}
	var card Card
	err := collection.FindOne(ctx, filter).Decode(&card)
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// Insert adds a new card document to the "cards" collection in MongoDB for testing
func Insert(ctx context.Context, card *Card) error {
	collection := database.GetCollection("cards")
	_, err := collection.InsertOne(ctx, card)
	return err
}

// UpdateByID updates an existing card document in the "cards" collection in MongoDB by its ID
func UpdateByID(ctx context.Context, cardID primitive.ObjectID, updates bson.M) error {
	collection := database.GetCollection("cards")
	filter := bson.M{"_id": cardID}
	update := bson.M{"$set": updates}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}
