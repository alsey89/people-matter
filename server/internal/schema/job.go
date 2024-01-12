package schema

import (
	"time"

	"gorm.io/gorm"
)

type JobInfo struct {
	gorm.Model
	UserID       uint `json:"userId"`
	DepartmentID uint `json:"departmentId"` // Foreign key for Department
	TitleID      uint `json:"titleId"`      // Foreign key for Title
	LocationID   uint `json:"locationId"`   // Foreign key for Location

	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`

	// Hierarchy
	Superiors    []Superior    `gorm:"foreignKey:JobInfoID"`
	Subordinates []Subordinate `gorm:"foreignKey:JobInfoID"`

	Location Location `gorm:"foreignKey:LocationID"` // One-to-one relationship with Location
}

type Superior struct {
	gorm.Model
	JobInfoID uint `json:"jobInfoId"` // References JobInfo
	UserID    uint `json:"userId"`    // References the superior's User
}

type Subordinate struct {
	gorm.Model
	JobInfoID uint `json:"jobInfoId"` // References JobInfo
	UserID    uint `json:"userId"`    // References the subordinate's User
}

type Title struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"` // Optional: More details about the title/role

	// Relationships
	JobInfo []JobInfo `json:"jobInfo"` // One-to-many relationship with JobInfo
}

type Department struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"` // Optional: More details about the department

	// Relationships
	JobInfo []JobInfo `json:"jobInfo"` // One-to-many relationship with JobInfo
}

type Location struct {
	gorm.Model
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postalCode"`

	// Relationships
	JobInfo []JobInfo `json:"jobInfo"` // One-to-many relationship with JobInfo
}

// callback function to sync title and department
func (j *JobInfo) AfterSave(tx *gorm.DB) (err error) {
	// Fetch the associated user
	var user User
	if err := tx.First(&user, j.UserID).Error; err != nil {
		return err
	}

	// Fetch the new department name and title name
	var department Department
	var title Title
	if err := tx.First(&department, "id = ?", j.DepartmentID).Error; err != nil {
		return err
	}
	if err := tx.First(&title, "id = ?", j.TitleID).Error; err != nil {
		return err
	}

	// Update the user's department name and title name
	user.Department = department.Name
	user.Title = title.Name

	// Save the updated user back to the database
	return tx.Save(&user).Error
}
