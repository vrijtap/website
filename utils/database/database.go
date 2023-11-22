package database

import (
    "context"
    "fmt"
    "sync"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB client and synchronization mutex
var (
    client     *mongo.Client
    clientLock sync.Mutex
)

// Connect initializes the MongoDB client and establishes a connection.
func Connect(uri string) error {
    // Use the SetServerAPIOptions() method to set the Stable API version to 1
    apiOpts := options.ServerAPI(options.ServerAPIVersion1)
    clientOpts := options.Client().ApplyURI(uri).SetServerAPIOptions(apiOpts)

    // Connect to MongoDB
    c, err := mongo.Connect(context.Background(), clientOpts)
    if err != nil {
        return err
    }

    // Ping MongoDB to verify the connection
    err = c.Ping(context.Background(), nil)
    if err != nil {
        return err
    }

    // Lock the mutex to safely set the client
    clientLock.Lock()
    defer clientLock.Unlock()

    // Store the initialized client
    client = c

    // Log the initialization
    fmt.Println("Connected to MongoDB")

    return nil
}

// GetClient returns the MongoDB client.
func GetClient() *mongo.Client {
    // Lock the mutex to safely retrieve the client
    clientLock.Lock()
    defer clientLock.Unlock()

    return client
}

// Close disconnects from MongoDB and closes the client.
func Close() error {
    // Lock the mutex to safely close the client
    clientLock.Lock()
    defer clientLock.Unlock()

    // Check if the client is already closed
    if client == nil {
        return nil
    }

    // Disconnect from MongoDB
    err := client.Disconnect(context.Background())
    if err != nil {
        return err
    }

    fmt.Println("MongoDB connection closed.")
    client = nil

    return nil
}
