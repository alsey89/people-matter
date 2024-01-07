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

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", pgHost, pgUser, pgPassword, pgDB, pgPort)

	var err error
	client, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true, // ! this is needed to translate postgres errors to gorm errors
	})
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
