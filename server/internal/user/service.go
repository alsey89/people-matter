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

//! All Users List -----------------------------------------------------------

func (us *UserService) GetAllUsersWithRole() ([]*schema.User, error) {

	existingUsers, err := us.userRepository.ReadAllAndPreloadJobTitleDepartment()
	if err != nil {
		return nil, fmt.Errorf("user.s.get_all_users: %w", err)
	}

	return existingUsers, nil // Return the filtered users
}

func (us *UserService) CreateListUserAndGetAllUsersWithRole(newUser *schema.User) ([]*schema.User, error) {
	_, err := us.userRepository.Create(newUser)
	if err != nil {
		return nil, fmt.Errorf("user.s.create_new_user: %w", err)
	}

	expandedUserList, err := us.GetAllUsersWithRole()
	if err != nil {
		return nil, fmt.Errorf("user.s.create_new_user: %w", err)
	}

	return expandedUserList, nil
}

func (us *UserService) UpdateListUserAndGetAllUsersWithRole(userID uint, updateData schema.User) ([]*schema.User, error) {
	_, err := us.userRepository.Update(userID, updateData)
	if err != nil {
		return nil, fmt.Errorf("user.s.update_user: %w", err)
	}

	expandedUserList, err := us.GetAllUsersWithRole()
	if err != nil {
		return nil, fmt.Errorf("user.s.create_new_user: %w", err)
	}

	return expandedUserList, nil
}

func (us *UserService) DeleteListUserAndGetAllUsersWithRole(userID uint) ([]*schema.User, error) {
	err := us.userRepository.Delete(userID)
	if err != nil {
		return nil, fmt.Errorf("s.delete_user: %w", err)
	}

	expandedUserList, err := us.GetAllUsersWithRole()
	if err != nil {
		return nil, fmt.Errorf("user.s.create_new_user: %w", err)
	}

	return expandedUserList, nil
}

// ! Single User Details -----------------------------------------------------------
func (us *UserService) GetSingleUserByIDAndExpand(ID uint) (*schema.User, error) {
	var existingUser *schema.User
	var err error

	existingUser, err = us.userRepository.ReadByIDAndPreloadJobTitleDepartment(ID)
	if err != nil {
		return nil, fmt.Errorf("user.s.get_user_by_id: %w", err)
	}

	return existingUser, nil
}
