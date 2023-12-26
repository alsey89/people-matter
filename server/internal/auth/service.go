package auth

import (
	"errors"
	"log"
	"verve-hrms/internal/user"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userService *user.UserService
}

func NewAuthService(userService *user.UserService) *AuthService {
	return &AuthService{userService: userService}
}

func (as *AuthService) Signup(email string, password string, username string) (*user.User, error) {
	emailIsAvailable, err := as.userService.IsEmailAvailable(email) //* using availability over existence because of return type (bool)
	if !emailIsAvailable {
		err = errors.New("email is already in use")
		return nil, err
	}
	if err != nil {
		log.Printf("Error checking email availability: %v", err)
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, err
	}

	newUser := user.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		IsAdmin:  false,
	}

	createdUser, err := as.userService.CreateNewAccount(newUser) //* this adds the ID to newUser
	if err != nil {
		log.Printf("Error creating new user: %v", err)
		return nil, err
	}

	return createdUser, nil
}

func (as *AuthService) Signin(email string, password string) (*user.User, error) {

	user, err := as.userService.GetUserByEmail(email)
	if err == mongo.ErrNoDocuments {
		err = errors.New("user not found")
		return nil, err
	}
	if err != nil {
		log.Printf("Error getting user by email: %v", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		err = errors.New("incorrect password")
		return nil, err
	}
	if err != nil {
		log.Printf("Error comparing password hashes: %v", err)
		return nil, err
	}

	return user, nil
}
