package user

import (
	"fmt"
	"verve-hrms/internal/schema"
)

type UserService struct {
	userRepository *UserRepository
}

func NewUserService(userRepository *UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (us *UserService) CreateNewAccount(newUser *schema.User) (*schema.User, error) {
	createdUser, err := us.userRepository.Create(newUser)
	if err != nil {
		return nil, fmt.Errorf("user.s.create_new_account: %w", err)
	}

	return createdUser, nil
}

func (us *UserService) GetAllUsers() ([]*schema.User, error) {
	users, err := us.userRepository.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("user.s.get_all_users: %w", err)
	}

	return users, nil
}

func (us *UserService) GetUserByEmail(email string) (*schema.User, error) {
	user, err := us.userRepository.ReadByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("user.s.get_user_by_email: %w", err)
	}

	return user, nil
}

func (us *UserService) GetUserByID(ID uint) (*schema.User, error) {
	user, err := us.userRepository.Read(ID)
	if err != nil {
		return nil, fmt.Errorf("user.s.get_user_by_id: %w", err)
	}

	return user, nil
}

func (us *UserService) UpdateUser(ID uint, updateData schema.User) (*schema.User, error) {
	updatedUser, err := us.userRepository.Update(ID, updateData)
	if err != nil {
		return nil, fmt.Errorf("user.s.update_user: %w", err)
	}

	return updatedUser, nil
}

func (us *UserService) DeleteUser(ID uint) error {
	err := us.userRepository.Delete(ID)
	if err != nil {
		return fmt.Errorf("s.delete_user: %w", err)
	}

	return nil
}
