package schema

import (
	"time"
)

type Memorial struct {
	BaseModel
	TenantID uint `json:"fspId"`
	// Memorial Information ---------------------
	Identifier      string `json:"identifier"` //should be firstname_lastname_dob (time.Time)
	IdentifierIsSet bool   `json:"identifierIsSet" gorm:"default:false"`
	HeaderImageURL  string `json:"headerImageUrl"`

	Title       string     `json:"title"` //should be first name + last name + 's + memorial
	FirstName   string     `json:"firstName"`
	LastName    string     `json:"lastName"`
	Description string     `json:"description"`
	Pronoun     string     `json:"pronoun"`
	DOB         *time.Time `json:"dob"`
	DOD         *time.Time `json:"dod"`
	Hometown    string     `json:"hometown"`
	Obituary    string     `json:"obituary"`
	Quote       string     `json:"quote"`

	// Inheritor Information ---------------------
	InheritorName  *string `json:"inheritorName"`
	InheritorEmail *string `json:"inheritorEmail"`

	// Associations ---------------------
	UserMemorialRoles []UserMemorialRole `json:"userMemorialRoles"`
}
