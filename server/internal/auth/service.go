package auth

import (
	"fmt"
	"verve-hrms/internal/user"

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
		return nil, ErrEmailNotAvailable
	}
	if err != nil {
		return nil, fmt.Errorf("s.signup: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("s.signup: %w", err)
	}

	newUser := user.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		IsAdmin:  false,
	}

	createdUser, err := as.userService.CreateNewAccount(newUser) //* this adds the ID to newUser
	if err != nil {
		return nil, fmt.Errorf("s.signup: %w", err)
	}

	return createdUser, nil
}

func (as *AuthService) Signin(email string, password string) (*user.User, error) {

	user, err := as.userService.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("s.signin: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, fmt.Errorf("s.signin: %w", ErrInvalidCredentials)
		}
		return nil, fmt.Errorf("s.signin: %w", err)
	}

	return user, nil
}
