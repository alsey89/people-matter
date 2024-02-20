package user

import (
	"fmt"
	"verve-hrms/internal/schema"

	"gorm.io/gorm"
)

type Repository interface {
	Create(newUser *schema.User) (*schema.User, error)
	Read(UserID uint) (*schema.User, error)
	ReadByEmail(email string) (*schema.User, error)
	Update(UserID uint, updateData *schema.User) (*schema.User, error)
	Delete(UserID uint) error
}

type UserRepository struct {
	client *gorm.DB
}

func NewUserRepository(client *gorm.DB) *UserRepository {
	return &UserRepository{client: client}
}

// Basic CRUD operations ------------------------------------------------------

func (ur UserRepository) Create(newUser *schema.User) (*schema.User, error) {
	result := ur.client.Create(newUser)
	if result.Error != nil {
		return nil, fmt.Errorf("user.r.create: %w", result.Error)
	}

	return newUser, nil
}

func (ur UserRepository) ReadByEmail(email string) (*schema.User, error) {
	var user schema.User
	result := ur.client.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("user.r.read_by_email: %w", result.Error)
	}

	return &user, nil
}

func (ur UserRepository) Update(UserID uint, updateData schema.User) (*schema.User, error) {
	var user schema.User

	updateDataMap := map[string]interface{}{
		"FirstName":  updateData.FirstName,
		"MiddleName": updateData.MiddleName,
		"LastName":   updateData.LastName,
		"Nickname":   updateData.Nickname,

		"Email":     updateData.Email,
		"AvatarURL": updateData.AvatarURL,
		"Role":      updateData.Role,
		"IsActive":  updateData.IsActive,
	}

	result := ur.client.Model(&user).Where("id = ?", UserID).Updates(updateDataMap)
	if result.Error != nil {
		return nil, fmt.Errorf("user.r.update: %w", result.Error)
	}

	return &user, nil
}

func (ur UserRepository) Delete(UserID uint) error {
	result := ur.client.Unscoped().Delete(&schema.User{}, "id = ?", UserID)
	if result.Error != nil {
		return fmt.Errorf("user.r.delete: %w", result.Error)
	}

	return nil
}

// Expanded operations --------------------------------------------------------

// note: Preloads roles (assignedJob > job > title and department)
func (ur UserRepository) ReadAllAndPreloadJobTitleDepartment() ([]*schema.User, error) {
	var users []*schema.User
	result := ur.client.
		Preload("AssignedJob.Job").
		Preload("AssignedJob.Job.Department").
		Preload("AssignedJob.Job.Location").
		Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("user.r.read_all_and_preload_title_department: %w", result.Error)
	}

	return users, nil
}

// note: Preloads all association for user (contactInfo, emergencyContact, assignedJob, payments, leave, attendance)
func (ur UserRepository) ReadByIDAndPreloadJobTitleDepartment(UserID uint) (*schema.User, error) {
	var user schema.User
	result := ur.client.
		Preload("ContactInfo").
		Preload("EmergencyContact").
		Preload("AssignedJob.Job").
		Preload("Payments").
		Preload("Leave").
		Preload("Attendance").
		First(&user, "id = ?", UserID)
	if result.Error != nil {
		return nil, fmt.Errorf("user.r.read: %w", result.Error)
	}

	return &user, nil
}
