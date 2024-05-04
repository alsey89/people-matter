package schema

import (
	"time"
)

// User Schema ---------------------------------------------------------------

// defines the role of the user
type RoleEnum string

const (
	RoleAdmin   RoleEnum = "admin"
	RoleManager RoleEnum = "manager"
	RoleUser    RoleEnum = "user"
)

type User struct {
	BaseModel
	CompanyID uint `json:"company_id" gorm:"onUpdate:CASCADE onDelete:CASCADE;not null"`
	// ------------------------------------------------------------------------------------------------
	Email      string      `json:"email"      gorm:"uniqueIndex"`
	Password   string      `json:"-"          gorm:"type:varchar(100)"` //* Password is not returned in JSON
	AvatarURL  string      `json:"avatarUrl"  gorm:"type:text"`
	Role       RoleEnum    `json:"role"       gorm:"type:enum('admin','manager','user');default:'user'"`
	LastLogin  *time.Time  `json:"lastLogin"  gorm:"default:null"`
	IsArchived bool        `json:"isArchived" gorm:"default:false"`
	Documents  []*Document `json:"documents"  gorm:"foreignKey:UserID"`
	// ------------------------------------------------------------------------------------------------
	FirstName        string            `json:"firstName"`
	MiddleName       *string           `json:"middleName"`
	LastName         string            `json:"lastName"`
	Nickname         string            `json:"nickname"`
	ContactInfo      *ContactInfo      `json:"contactInfo"       gorm:"foreignKey:UserID"`
	EmergencyContact *EmergencyContact `json:"emergencyContact"  gorm:"foreignKey:UserID"`
	// ------------------------------------------------------------------------------------------------
	UserPosition *UserPosition `json:"assignedJob" gorm:"foreignKey:UserID"`
	// ------------------------------------------------------------------------------------------------
	SalaryID *uint      `json:"salaryId"  gorm:"foreignKey:UserID"`
	Payments []*Payment `json:"payments" gorm:"foreignKey:UserID"`
	// ------------------------------------------------------------------------------------------------
	Leave      []*Leave      `json:"leave" gorm:"foreignKey:UserID"`
	Attendance []*Attendance `json:"attendance" gorm:"foreignKey:UserID"`
	// ------------------------------------------------------------------------------------------------
}

type ContactInfo struct {
	BaseModel
	CompanyID uint `json:"companyId" gorm:"onUpdate:CASCADE;onDelete:CASCADE"`
	// ------------------------------------------------------------------------------------------------
	UserID     uint   `json:"userId" gorm:"onUpdate:CASCADE;onDelete:CASCADE"`
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
	Mobile     string `json:"mobile"`
	Email      string `json:"email"`
	// ------------------------------------------------------------------------------------------------
}

type EmergencyContact struct {
	BaseModel
	CompanyID uint `json:"companyId" gorm:"onUpdate:CASCADE;onDelete:CASCADE"`
	// ------------------------------------------------------------------------------------------------
	UserID     uint    `json:"userId" gorm:"onUpdate:CASCADE;onDelete:CASCADE"`
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName"`
	LastName   string  `json:"lastName"`
	Relation   string  `json:"relation"`
	Mobile     string  `json:"mobile"`
	Email      string  `json:"email"`
	// ------------------------------------------------------------------------------------------------
}
