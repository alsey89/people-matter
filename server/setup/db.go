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

	// Open the connection
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true, //* This is needed to translate postgres errors to gorm errors
	})
	if err != nil {
		panic(err)
	}

	//set up automigrate
	err = db.AutoMigrate(&user.User{}, &user.ContactInfo{}, &user.EmergencyContact{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to local MongoDB!")

	return client
}
