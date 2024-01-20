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

//! Auth ------------------------------------------------------------

func (us *UserService) CreateNewUser(newUser *schema.User) (*schema.User, error) {
	var createdUser *schema.User
	var err error

	createdUser, err = us.userRepository.Create(newUser)
	if err != nil {
		return nil, fmt.Errorf("user.s.create_new_account: %w", err)
	}

	return createdUser, nil
}

func (us *UserService) GetUserByEmail(email string) (*schema.User, error) {
	var existingUser *schema.User
	var err error

	existingUser, err = us.userRepository.ReadByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("user.s.get_user_by_email: %w", err)
	}

	return existingUser, nil
}

//! User -----------------------------------------------------------

func (us *UserService) GetAllUsersAndExpandRoles() ([]*schema.User, error) {
	var existingUsers []*schema.User
	var err error

	existingUsers, err = us.userRepository.ReadAllAndExpandRoles()
	if err != nil {
		return nil, fmt.Errorf("user.s.get_all_users: %w", err)
	}

	return existingUsers, nil
}

func (us *UserService) CreateNewUserAndExpandRole(newUser *schema.User) (*schema.User, error) {
	var createdUser *schema.User
	var err error

	createdUser, err = us.userRepository.Create(newUser)
	if err != nil {
		return nil, fmt.Errorf("user.s.create_new_account: %w", err)
	}

	createdUser, err = us.userRepository.ReadAndExpandRole(createdUser.ID)
	if err != nil {
		return nil, fmt.Errorf("user.s.create_new_account: %w", err)
	}

	return createdUser, nil
}

func (us *UserService) GetUserByIDAndExpandRole(ID uint) (*schema.User, error) {
	var existingUser *schema.User
	var err error

	existingUser, err = us.userRepository.ReadAndExpandRole(ID)
	if err != nil {
		return nil, fmt.Errorf("user.s.get_user_by_id: %w", err)
	}

	return existingUser, nil
}

func (us *UserService) UpdateAndReturnUserAndExpandRole(ID uint, updateData schema.User) (*schema.User, error) {
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
