package schema

import (
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	Name       string  `json:"name"`
	Website    *string `json:"website"`
	Email      *string `json:"email"`
	Phone      *string `json:"phone"`
	Address    *string `json:"address"`
	City       *string `json:"city"`
	State      *string `json:"state"`
	Country    *string `json:"country"`
	PostalCode *string `json:"postalCode"`

	// Relationships
	Departments []Department `gorm:"foreignKey:CompanyID"`
	Titles      []Title      `gorm:"foreignKey:CompanyID"`
	Locations   []Location   `gorm:"foreignKey:CompanyID"`
}

type Title struct {
	gorm.Model
	CompanyID   uint   `json:"companyId"`
	Name        string `json:"name"`
	Description string `json:"description"`

	// Relationships
	Job []Job `json:"job"`
}

type Department struct {
	gorm.Model
	CompanyID   uint   `json:"companyId"`
	Name        string `json:"name"`
	Description string `json:"description"`

	// Relationships
	Job []Job `json:"job"`
}

type Location struct {
	gorm.Model
	CompanyID  uint   `json:"companyId"`
	Name       string `json:"name"`
	HeadOffice bool   `json:"headOffice"`
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postalCode"`

	// Relationships
	Job []Job `json:"job"`
}
