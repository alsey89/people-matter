package auth

import (
	"fmt"

	"github.com/alsey89/people-matter/schema"
	"golang.org/x/crypto/bcrypt"
)

func (a *Domain) SignupService(email string, password string) (*schema.User, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("auth.s.signup: %w", err)
	}

	newUser := schema.User{
		FirstName: "New",
		LastName:  "User",
		Email:     email,
		Password:  string(hashedPassword),
		Role:      "user",
	}

	createdUser, err := a.CreateNewUser(&newUser) //* this adds the ID to newUser
	if err != nil {
		return nil, fmt.Errorf("auth.s.signup: %w", err)
	}

	return createdUser, nil
}

func (a *Domain) SigninService(email string, password string) (*schema.User, error) {

	user, err := a.GetUserByEmail(email)
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

// ---------------------------------------------------------------------------

func (a *Domain) CreateNewUser(newUser *schema.User) (*schema.User, error) {
	db := a.params.Database.GetDB()

	result := db.Create(newUser)
	if result.Error != nil {
		return nil, fmt.Errorf("user.r.create: %w", result.Error)
	}

	//add returned ID from Create to newUser
	createdUserID := uint(result.RowsAffected)

	newUser.ID = createdUserID

	return newUser, nil
}

func (a *Domain) GetUserByEmail(email string) (*schema.User, error) {
	db := a.params.Database.GetDB()

	var user *schema.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("user.r.read_by_email: %w", result.Error)
	}

	return user, nil
}
