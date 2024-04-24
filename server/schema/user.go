package schema

import (
	"time"
)

// User model ---------------------------------------------------------------
type User struct {
	//* account information
	BaseModel      // This includes fields ID, CreatedAt, UpdatedAt, DeletedAt
	CompanyID uint `json:"companyId" gorm:"onUpdate:CASCADE;onDelete:CASCADE"`

	Email     string     `json:"email"      gorm:"uniqueIndex"`
	Password  string     `json:"-"          gorm:"type:varchar(100)"` //* Password is not returned in JSON
	AvatarURL string     `json:"avatarUrl"  gorm:"type:text"`
	Role      string     `json:"role"       gorm:"enum:admin,manager,user;default:user"`
	IsActive  bool       `json:"isActive"   gorm:"default:true"`
	LastLogin *time.Time `json:"lastLogin"  gorm:"default:null"`

	//* personal information
	FirstName        string            `json:"firstName"`
	MiddleName       *string           `json:"middleName"`
	LastName         string            `json:"lastName"`
	Nickname         string            `json:"nickname"`
	ContactInfo      *ContactInfo      `json:"contactInfo"       gorm:"foreignKey:UserID"`
	EmergencyContact *EmergencyContact `json:"emergencyContact"  gorm:"foreignKey:UserID"`

	//*job information
	UserPosition *UserPosition `json:"assignedJob" gorm:"foreignKey:UserID"`

	//* salary information
	SalaryID *uint     `json:"salaryId"  gorm:"foreignKey:UserID"`
	Payments []Payment `json:"Payments" gorm:"foreignKey:UserID"`

	//* leave & attendance
	Leave      []Leave      `json:"leave"`
	Attendance []Attendance `json:"attendance"`
}

// enum constants for User.Role
const (
	RoleAdmin   = "admin"
	RoleManager = "manager"
	RoleUser    = "user"
)

type ContactInfo struct {
	BaseModel
	UserID     uint   `json:"userId"`
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
	Mobile     string `json:"mobile"`
	Email      string `json:"email"`
}

type EmergencyContact struct {
	BaseModel
	UserID     uint    `json:"userId"`
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName"`
	LastName   string  `json:"lastName"`
	Relation   string  `json:"relation"`
	Mobile     string  `json:"mobile"`
	Email      string  `json:"email"`
}
