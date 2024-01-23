package orders

import (
    "website/utils/database"

    "context"
    "time"
    "math"
    "errors"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Order represents an order for beer on a card
type Order struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    CardID      primitive.ObjectID `bson:"card_id"`
    OrderDate   time.Time          `bson:"order_date"`
    Status      string             `bson:"status"`
    Quantity    uint               `bson:"quantity"`
    TotalAmount float64            `bson:"total_amount"`
}

// New creates a new Order instance with default values
func New(cardID primitive.ObjectID, quantity uint, price float64) Order {
    return Order{
        CardID:      cardID,
        OrderDate:   time.Now(),
        Status:      "Pending",
        Quantity:    quantity,
        TotalAmount: math.Round(float64(quantity) * price*100)/100,
    }
}

// GetByID retrieves an order document from MongoDB by its ID
func GetByID(ctx context.Context, orderID primitive.ObjectID) (*Order, error) {
    // Setup the database request
	collection := database.GetCollection("orders")
    filter := bson.M{"_id": orderID}

    // Get the order from the collection "orders"
    var order Order
    err := collection.FindOne(ctx, filter).Decode(&order)
    if err != nil {
        return nil, err
    }

    // If no error was received, return the order
    return &order, nil
}

// Insert adds a new order document to the "orders" collection in MongoDB
func Insert(ctx context.Context, order *Order) (*Order, error) {
    // Setup the database request
    collection := database.GetCollection("orders")

    // Insert the order into the collection "orders"
    insertOneResult , err := collection.InsertOne(ctx, order)
    if err != nil {
        return nil, err
    }

    // Assert the InsertedID as a primitive.ObjectID
	id, ok := insertOneResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("Failed to assert InsertedID as primitive.ObjectID")
	}

    // If the assertion succeeds, return the inserted order
    order.ID = id
    return order, nil
}

// UpdateStatus updates the status of an existing order document in the "orders" collection in MongoDB
func UpdateStatus(ctx context.Context, orderID primitive.ObjectID, newStatus string) error {
    // Setup the database request
    collection := database.GetCollection("orders")
    filter := bson.M{"_id": orderID}
    update := bson.M{"$set": bson.M{"status": newStatus}}

    // Update the order with the new status
    _, err := collection.UpdateOne(ctx, filter, update)
    return err
}
