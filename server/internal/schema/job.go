package schema

import (
	"time"

	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	UserID       uint  `json:"userId"`
	DepartmentID uint  `json:"departmentId"`
	TitleID      uint  `json:"titleId"`
	LocationID   uint  `json:"locationId"`
	ManagerID    *uint `json:"managerId"`

	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`

	DirectReports []Job `gorm:"foreignKey:ManagerID"` // Jobs where this job is the manager
}

// callback function to sync title and department
func (j *Job) AfterSave(tx *gorm.DB) (err error) {
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
