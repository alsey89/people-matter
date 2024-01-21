package schema

import (
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	Name       string  `json:"name"`
	LogoURL    *string `json:"logoUrl"`
	Website    *string `json:"website"`
	Email      *string `json:"email"`
	Phone      *string `json:"phone"`
	Address    *string `json:"address"`
	City       *string `json:"city"`
	State      *string `json:"state"`
	Country    *string `json:"country"`
	PostalCode *string `json:"postalCode"`

	// Relationships
	Departments []Department `json:"departments" gorm:"foreignKey:CompanyID"`
	Titles      []Title      `json:"titles" gorm:"foreignKey:CompanyID"`
	Locations   []Location   `json:"locations" gorm:"foreignKey:CompanyID"`
}

type Title struct {
	gorm.Model
	CompanyID   uint   `json:"companyId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Department struct {
	gorm.Model
	CompanyID   uint   `json:"companyId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Location struct {
	gorm.Model
	CompanyID    uint   `json:"companyId"`
	Name         string `json:"name"`
	IsHeadOffice bool   `json:"isHeadOffice"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	PostalCode   string `json:"postalCode"`
}
