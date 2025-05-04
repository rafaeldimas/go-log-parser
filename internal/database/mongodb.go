package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(options.Credential{Username: "root", Password: "password", AuthSource: "admin"})

	return mongo.Connect(context.Background(), clientOptions)
}

func NewDatabase(client *mongo.Client) *mongo.Database {
	return client.Database("logs")
}

func NewCollection(database *mongo.Database) *mongo.Collection {
	return database.Collection("logs")
}

func Close(client *mongo.Client) {
	client.Disconnect(context.Background())
}
