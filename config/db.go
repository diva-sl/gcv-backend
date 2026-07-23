package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var MongoClient *mongo.Client

// ConnectDB initializes the connection to MongoDB
func ConnectDB() {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "GCV"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("MongoDB Connection Error: %v", err)
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB Ping Failed: %v", err)
	}

	MongoClient = client
	DB = client.Database(dbName)
	log.Printf("Connected successfully to MongoDB Database: %s", dbName)
}

// GetCollection returns a handler to a specific collection
func GetCollection(collectionName string) *mongo.Collection {
	return DB.Collection(collectionName)
}