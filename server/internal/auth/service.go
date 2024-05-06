package auth

import (
	"fmt"
	"log"

	"github.com/alsey89/people-matter/schema"
	"golang.org/x/crypto/bcrypt"
)

// // Hash password, create new user, and return user *with new ID*
// func (a *Domain) SignupService(email string, password string, companyID uint) (*schema.User, error) {

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return nil, fmt.Errorf("SignupService: %w", err)
// 	}

// 	newUser := schema.User{
// 		CompanyID: companyID,
// 		FirstName: "New",
// 		LastName:  "User",
// 		Email:     email,
// 		Password:  string(hashedPassword),
// 		Role:      "user",
// 	}

// 	createdUser, err := a.CreateNewUser(&newUser) //* this adds the ID to newUser
// 	if err != nil {
// 		return nil, fmt.Errorf("SignupService: %w", err)
// 	}

// 	return createdUser, nil
// }

// Search for user by email, compare password, and return user if successful.
// If user is a manager, fetch location information.
func (a *Domain) SignIn(email string, password string) (*schema.User, error) {
	db := a.params.Database.GetDB()

	var user schema.User

	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("[SignIn] %w", result.Error)
	}

	log.Printf("User: %+v", user)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, fmt.Errorf("[SignIn]: %w", ErrInvalidCredentials)
		}
		return nil, fmt.Errorf("[SignIn]: %w", err)
	}

	// If user is a manager, fetch location information
	if user.Role == "manager" {
		result = db.
			Joins("UserPosition").
			Joins("Location").
			Where("id = ?", user.ID).
			Find(&user)
		if result.Error != nil {
			return nil, fmt.Errorf("[SignIn]: %w", result.Error)
		}
	}

	return &user, nil
}
