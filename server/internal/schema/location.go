package schema

type Location struct {
	BaseModel
	TenantIdentifier string `json:"-" gorm:"not null;index"`
	// Location Information ---------------------
	Name string `json:"name"`
	// Contact Information ---------------------
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
	//Address Information ---------------------
	ContactAddress    string `json:"address"`
	ContactCity       string `json:"city"`
	ContactCountry    string `json:"country"`
	ContactPostalCode string `json:"postalCode"`
	// Associations
	Positions []Position `json:"positions" gorm:"foreignKey:TenantID"`
}
