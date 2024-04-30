package auth

import (
	"fmt"

	"github.com/alsey89/people-matter/schema"
	"golang.org/x/crypto/bcrypt"
)

// Hash password, create new user, and return user *with new ID*
func (a *Domain) SignupService(email string, password string, companyID uint) (*schema.User, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("SignupService: %w", err)
	}

	newUser := schema.User{
		CompanyID: companyID,
		FirstName: "New",
		LastName:  "User",
		Email:     email,
		Password:  string(hashedPassword),
		Role:      "user",
	}

	createdUser, err := a.CreateNewUser(&newUser) //* this adds the ID to newUser
	if err != nil {
		return nil, fmt.Errorf("SignupService: %w", err)
	}

	return createdUser, nil
}

// Search for user by email, compare password, and return user if successful
func (a *Domain) SigninService(email string, password string) (*schema.User, error) {

	user, err := a.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("SigninService: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, fmt.Errorf("SigninService: %w", ErrInvalidCredentials)
		}
		return nil, fmt.Errorf("SigninService: %w", err)
	}

	return user, nil
}

// ---------------------------------------------------------------------------

// Create new user and return user *with new ID*
func (a *Domain) CreateNewUser(newUser *schema.User) (*schema.User, error) {
	db := a.params.Database.GetDB()

	result := db.Create(newUser)
	if result.Error != nil {
		return nil, fmt.Errorf("CreateNewUser: %w", result.Error)
	}

	//add returned ID from Create to newUser
	createdUserID := uint(result.RowsAffected)

	newUser.ID = createdUserID

	return newUser, nil
}

// Search for user by email and return user
func (a *Domain) GetUserByEmail(email string) (*schema.User, error) {
	db := a.params.Database.GetDB()

	var user *schema.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("GetUserByEmail: %w", result.Error)
	}

	return user, nil
}
