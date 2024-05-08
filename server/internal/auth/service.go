package auth

import (
	"fmt"

	"github.com/alsey89/people-matter/schema"
	"golang.org/x/crypto/bcrypt"
)

// Search for user by email, compare password, and return user if successful.
// If user is a manager, fetch location information.
func (d *Domain) SignIn(email string, password string) (*schema.User, error) {
	db := d.params.Database.GetDB()

	var user schema.User

	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("[SignIn] %w", result.Error)
	}

	if !user.IsActive {
		return nil, fmt.Errorf("[SignIn]: %w", ErrUserNotConfirmed)
	}

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

func (d *Domain) ConfirmEmail(userID uint, companyID uint) error {
	db := d.params.Database.GetDB()

	//query user and mark as confirmed
	result := db.Model(&schema.User{}).
		Where("id = ? AND company_id = ?", userID, companyID).
		Update("is_active", true)
	if result.Error != nil {
		return fmt.Errorf("[Confirmation]: %w", result.Error)
	}

	return nil
}
