package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserService struct {
	userRepository *UserRepository
}

func NewUserService(userRepository *UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (us *UserService) CreateNewAccount(newUser User) (*User, error) {
	newUser, err := us.userRepository.Create(newUser)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (us *UserService) GetUserByEmail(email string) (*User, error) {
	user, err := us.userRepository.ReadByEmail(email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *UserService) GetUserByID(objID *primitive.ObjectID) (*User, error) {
	user, err := us.userRepository.Read(objID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *UserService) UpdateUser(objID *primitive.ObjectID, updateData User) (*User, error) {
	updatedUser, err := us.userRepository.Update(objID, updateData)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (us *UserService) DeleteUser(objID *primitive.ObjectID) error {
	err := us.userRepository.Delete(objID)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) IsEmailAvailable(email string) (bool, error) {
	userCount, err := us.userRepository.CountUsersByEmail(email)
	if err != nil {
		return false, err
	}

	isEmailAvailable := userCount == 0

	return isEmailAvailable, nil
}
