package user

import (
	"context"
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

func (ur *UserRepository) Create(newUser User) (User, error) {
	coll := ur.client.Database(viper.GetString("DB_NAME")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := coll.InsertOne(ctx, newUser)
	if err != nil {
		return User{}, err
	}

	newUser.ID = result.InsertedID.(primitive.ObjectID)
	return newUser, nil
}

func (ur *UserRepository) Read(objID *primitive.ObjectID) (User, error) {
	coll := ur.client.Database(viper.GetString("DB_NAME")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser User
	err := coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&existingUser)
	if err != nil {
		return User{}, err
	}

	return existingUser, nil
}

func (ur *UserRepository) ReadByEmail(email string) (User, error) {
	coll := ur.client.Database(viper.GetString("DB_NAME")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser User
	err := coll.FindOne(ctx, bson.M{"email": email}).Decode(&existingUser)
	if err != nil {
		return User{}, err
	}

	return existingUser, nil
}

func (ur *UserRepository) Update(objID *primitive.ObjectID, updateData User) (User, error) {
	coll := ur.client.Database(viper.GetString("DB_NAME")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": updateData}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedUser User
	err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedUser)
	if err != nil {
		return User{}, err
	}

	return updatedUser, nil
}

func (ur *UserRepository) Delete(objID *primitive.ObjectID) error {
	coll := ur.client.Database(viper.GetString("DB_NAME")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := coll.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (ur *UserRepository) CheckEmailAvailability(email string) (bool, error) {
	coll := ur.client.Database(viper.GetString("DB_NAME")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := coll.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}

	emailIsAvailable := count == 0

	return emailIsAvailable, nil
}
