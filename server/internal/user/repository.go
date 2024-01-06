package user

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

type Repository interface {
	Create(newUser User) (User, error)
	Read(objID *primitive.ObjectID) (User, error)
	ReadByEmail(email string) (User, error)
	Update(objID *primitive.ObjectID, updateData User) (User, error)
	Delete(objID *primitive.ObjectID) error
}

type UserRepository struct {
	client *gorm.DB
}

func NewUserRepository(client *gorm.DB) *UserRepository {
	return &UserRepository{client: client}
}

// ! Basic CRUD operations ------------------------------------------------------
func (ur UserRepository) Create(newUser User) (*User, error) {
	result := ur.client.Create(&newUser)
	if result.Error != nil {
		return nil, fmt.Errorf("ur.create: %w", result.Error)
	}

	return &newUser, nil
}

func (ur UserRepository) Read(objID *primitive.ObjectID) (*User, error) {
	var user User
	result := ur.client.First(&user, objID)
	if result.Error != nil {
		return nil, fmt.Errorf("ur.read: %w", result.Error)
	}

	return &user, nil
}

func (ur UserRepository) ReadByEmail(email string) (*User, error) {
	var user User
	result := ur.client.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("ur.read_by_email: %w", result.Error)
	}

	return &user, nil
}

func (ur UserRepository) Update(objID *primitive.ObjectID, updateData User) (*User, error) {
	var user User
	result := ur.client.First(&user, objID)
	if result.Error != nil {
		return nil, fmt.Errorf("ur.update: %w", result.Error)
	}

	result = ur.client.Model(&user).Updates(updateData)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (ur UserRepository) Delete(objID *primitive.ObjectID) error {
	var user User
	result := ur.client.Delete(&user, objID)
	if result.Error != nil {
		return fmt.Errorf("ur.delete: %w", result.Error)
	}

	return nil
}
