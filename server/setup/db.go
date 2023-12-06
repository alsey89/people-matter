package setup

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func GetMongoClient() *mongo.Client {
	if client != nil {
		return client
	}

	// Local MongoDB URI
	localURI := "mongodb://mongodb:27017"

	opts := options.Client().ApplyURI(localURI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	err = client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to local MongoDB!")

	return client
}
