package auth

import (
	"fmt"
	"verve-hrms/internal/schema"
	"verve-hrms/internal/user"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userService *user.UserService
}

func NewAuthService(userService *user.UserService) *AuthService {
	return &AuthService{userService: userService}
}

func (as *AuthService) Signup(email string, password string) (*schema.User, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("auth.s.signup: %w", err)
	}

	newUser := schema.User{
		Email:    email,
		Password: string(hashedPassword),
		IsAdmin:  false,
	}

	createdUser, err := as.userService.CreateNewAccount(&newUser) //* this adds the ID to newUser
	if err != nil {
		return nil, fmt.Errorf("auth.s.signup: %w", err)
	}

	return createdUser, nil
}

func (as *AuthService) Signin(email string, password string) (*schema.User, error) {

	user, err := as.userService.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("auth.s.signin: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, fmt.Errorf("auth.s.signin: %w", ErrInvalidCredentials)
		}
		return nil, fmt.Errorf("auth.s.signin: %w", err)
	}

	return user, nil
}
