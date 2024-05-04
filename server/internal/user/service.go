package user

import (
	"fmt"

	"github.com/alsey89/people-matter/schema"
)

//! User ------------------------------------------------------------

// Get all users in company
func (d *Domain) GetAllUsers(companyID *uint) ([]schema.User, error) {
	db := d.params.Database.GetDB()

	var users []schema.User

	result := db.Where("company_id = ?", companyID).Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("[GetAllUsers] %w", result.Error)
	}

	return users, nil
}

// Get all users by location
func (d *Domain) GetUsersByLocation(companyID *uint, locationID *uint) ([]schema.User, error) {
	db := d.params.Database.GetDB()

	var users []schema.User

	result := db.
		Model(&schema.User{}).
		Joins("UserPosition").
		Joins("Location").
		Where("users.company_id = ? AND locations.company_id =? AND locations.id = ?", companyID, companyID, locationID).
		Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("[GetUsersByLocation] %w", result.Error)
	}

	return users, nil
}

// Get all users by department
func (d *Domain) GetUsersByDepartment(companyID *uint, departmentID *uint) ([]schema.User, error) {
	db := d.params.Database.GetDB()
	var users []schema.User

	result := db.
		Model(&schema.User{}).
		Joins("JOIN user_positions ON user_positions.user_id = users.id").
		Joins("JOIN positions ON user_positions.position_id = positions.id").
		Joins("JOIN departments ON positions.department_id = departments.id").
		Where("users.company_id = ? AND departments.company_id =? AND departments.id = ?", companyID, companyID, departmentID).
		Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("[GetUsersByDepartment] %w", result.Error)
	}

	return users, nil
}

// Get a single user by ID
func (d *Domain) GetUser(companyID *uint, userID *uint) (*schema.User, error) {
	db := d.params.Database.GetDB()

	var existingUser schema.User

	result := db.Where("company_id = ? AND id = ?", companyID, userID).First(&existingUser)
	if result.Error != nil {
		return nil, fmt.Errorf("[GetUser] %w", result.Error)
	}

	return &existingUser, nil
}
