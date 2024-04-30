package schema

type Company struct {
	BaseModel
	// ------------------------------------------------------------------------------------------------
	Name    string `json:"name"`
	LogoURL string `json:"logoUrl"`
	Website string `json:"website"`
	Email   string `json:"email"`
	// ------------------------------------------------------------------------------------------------
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postalCode"`
	// ------------------------------------------------------------------------------------------------
	Departments []Department `json:"departments" gorm:"foreignKey:CompanyID"`
	Locations   []Location   `json:"locations" gorm:"foreignKey:CompanyID"`
	Positions   []Position   `json:"jobs" gorm:"foreignKey:CompanyID"`
	// ------------------------------------------------------------------------------------------------
}

type Department struct {
	BaseModel
	CompanyID uint `json:"company_id" gorm:"onUpdate:CASCADE onDelete:CASCADE;not null"`
	// ------------------------------------------------------------------------------------------------
	Name        string `json:"name"`
	Description string `json:"description"`
	// ------------------------------------------------------------------------------------------------
}

type Location struct {
	BaseModel
	CompanyID uint `json:"company_id" gorm:"onUpdate:CASCADE onDelete:CASCADE;not null"`
	// ------------------------------------------------------------------------------------------------
	Name         string `json:"name"`
	IsHeadOffice bool   `json:"isHeadOffice"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	PostalCode   string `json:"postalCode"`
	// ------------------------------------------------------------------------------------------------
}
