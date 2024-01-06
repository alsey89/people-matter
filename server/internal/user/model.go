package user

import (
	"time"

	"gorm.io/gorm"
)

// * User model ---------------------------------------------------------------
type User struct {
	//* account information
	gorm.Model            // This includes fields ID, CreatedAt, UpdatedAt, DeletedAt
	Email      string     `json:"email"      gorm:"uniqueIndex"`
	Password   string     `json:"-"          gorm:"type:varchar(100)"` //* Password is not returned in JSON
	IsAdmin    bool       `json:"isAdmin"    gorm:"default:false"`
	AvatarURL  *string    `json:"avatarUrl"  gorm:"type:text"`
	IsActive   bool       `json:"isActive"   gorm:"default:true"`
	IsVerified bool       `json:"isVerified" gorm:"default:false"`
	LastLogin  *time.Time `json:"lastLogin"  gorm:"default:null"`

	//* basic personal information
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName"`
	LastName   string  `json:"lastName"`
	Nickname   *string `json:"nickname"`
	//* detailed personal information
	ContactInfo      *ContactInfo      `json:"contactInfo"`
	EmergencyContact *EmergencyContact `json:"emergencyContact"`
}

type ContactInfo struct {
	gorm.Model
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
	gorm.Model
	UserID     uint    `json:"userId"`
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName"`
	LastName   string  `json:"lastName"`
	Relation   string  `json:"relation"`
	Mobile     string  `json:"mobile"`
	Email      string  `json:"email"`
}
