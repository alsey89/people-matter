package user

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/spf13/viper"
)

type Repository interface {
	Create(newUser User) (User, error)
	Read(objID *primitive.ObjectID) (User, error)
	ReadByEmail(email string) (User, error)
	Update(objID *primitive.ObjectID, updateData User) (User, error)
	Delete(objID *primitive.ObjectID) error
	CheckIfEmailInUse(email string) (int64, error)
}

type UserRepository struct {
	client *mongo.Client
}

func NewUserRepository(client *mongo.Client) *UserRepository {
	return &UserRepository{client: client}
}

// ! Basic CRUD operations ------------------------------------------------------

func (ur *UserRepository) Create(newUser User) (*User, error) {
	coll := ur.client.Database(viper.GetString("DB_NAME")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := coll.InsertOne(ctx, newUser)
	if err != nil {
		log.Printf("Error in UserRepository Create, Error Type: %s, Error Details: %+v", reflect.TypeOf(err), err)
		return nil, fmt.Errorf("r.create: %w", err)
	}

	newUser.ID = result.InsertedID.(primitive.ObjectID)
	return &newUser, nil
}

func (ur *UserRepository) Read(objID *primitive.ObjectID) (*User, error) {
	coll := ur.client.Database(viper.GetString("DB_NAME")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser User
	err := coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&existingUser)
	if err != nil {
		return nil, fmt.Errorf("r.read: %w", err)
	}

	return &existingUser, nil
}

func (ur *UserRepository) Update(objID *primitive.ObjectID, updateData User) (*User, error) {
	coll := ur.client.Database(viper.GetString("DB_NAME")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": updateData}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedUser User
	err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedUser)
	if err != nil {
		return nil, fmt.Errorf("r.update: %w", err)
	}

	return &updatedUser, nil
}

func (ur *UserRepository) Delete(objID *primitive.ObjectID) error {
	coll := ur.client.Database(viper.GetString("DB_NAME")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := coll.DeleteOne(ctx, bson.M{"_id": objID})
	if result.DeletedCount == 0 {
		return fmt.Errorf("r.delete: %w", ErrUserNotFound)
	}
	if err != nil {
		return fmt.Errorf("r.delete: %w", err)
	}

	return nil
}

//! Specific operations ------------------------------------------------------

func (ur *UserRepository) ReadByEmail(email string) (*User, error) {
	coll := ur.client.Database(viper.GetString("DB_NAME")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser User
	err := coll.FindOne(ctx, bson.M{"email": email}).Decode(&existingUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("r.read_by_email: %w", ErrUserNotFound)
		}
		return nil, fmt.Errorf("r.read_by_email: %w", err)
	}

	return &existingUser, nil
}

func (ur *UserRepository) CountUsersByEmail(email string) (int64, error) {
	coll := ur.client.Database(viper.GetString("DB_NAME")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := coll.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return 0, fmt.Errorf("r.count_users_by_email: %w", err)
	}

	return count, nil
}
