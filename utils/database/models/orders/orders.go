package orders

import (
    "context"
    "time"
    "math"
    "backend/utils/database"

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
    var order Order
    collection := database.GetClient().Database("backend").Collection("orders")
    filter := bson.M{"_id": orderID}
    err := collection.FindOne(ctx, filter).Decode(&order)
    if err != nil {
        return nil, err
    }
    return &order, nil
}

// Insert adds a new order document to the "orders" collection in MongoDB
func Insert(ctx context.Context, order *Order) (*Order, error) {
    collection := database.GetClient().Database("backend").Collection("orders")
    result, err := collection.InsertOne(ctx, order)
    if err != nil {
        return nil, err
    }
    order.ID = result.InsertedID.(primitive.ObjectID)
    return order, nil
}

// UpdateStatus updates the status of an existing order document in the "orders" collection in MongoDB
func UpdateStatus(ctx context.Context, orderID primitive.ObjectID, newStatus string) error {
    collection := database.GetClient().Database("backend").Collection("orders")

    // Define a filter to specify which order to update based on its ID
    filter := bson.M{"_id": orderID}

    // Define an update operation using the $set operator to modify the "status" field
    update := bson.M{"$set": bson.M{"status": newStatus}}

    // Perform the update operation
    _, err := collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }

    return nil
}
