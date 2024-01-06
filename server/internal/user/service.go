package user

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	userRepository *UserRepository
}

func NewUserService(userRepository *UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (us *UserService) CreateNewAccount(newUser User) (*User, error) {
	createdUser, err := us.userRepository.Create(newUser)
	if err != nil {
		return nil, fmt.Errorf("s.create_new_account: %w", err)
	}

	return createdUser, nil
}

func (us *UserService) GetUserByEmail(email string) (*User, error) {
	user, err := us.userRepository.ReadByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("s.get_user_by_email: %w", err)
	}

	return user, nil
}

func (us *UserService) GetUserByID(objID *primitive.ObjectID) (*User, error) {
	user, err := us.userRepository.Read(objID)
	if err != nil {
		return nil, fmt.Errorf("s.get_user_by_id: %w", err)
	}

	return user, nil
}

func (us *UserService) UpdateUser(objID *primitive.ObjectID, updateData User) (*User, error) {
	updatedUser, err := us.userRepository.Update(objID, updateData)
	if err != nil {
		return nil, fmt.Errorf("s.update_user: %w", err)
	}

	return updatedUser, nil
}

func (us *UserService) DeleteUser(objID *primitive.ObjectID) error {
	err := us.userRepository.Delete(objID)
	if err != nil {
		return fmt.Errorf("s.delete_user: %w", err)
	}

	return nil
}
