package schema

type Company struct {
	BaseModel
	TenantIdentifier string `json:"-" gorm:"not null;index"`
	// Company Information ---------------------
	Name    string `json:"name"`
	LogoURL string `json:"logoUrl"`
	// Contact Information ---------------------
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
	//Address Information ---------------------
	ContactAddress    string `json:"address"`
	ContactCity       string `json:"city"`
	ContactCountry    string `json:"country"`
	ContactPostalCode string `json:"postalCode"`

	BillingAddress    string `json:"billingAddress"`
	BillingCity       string `json:"billingCity"`
	BillingCountry    string `json:"billingCountry"`
	BillingPostalCode string `json:"billingPostalCode"`
	// Account Information ---------------------
	BranchQuota       int     `json:"branchQuota"`
	BranchQuotaUsed   int     `json:"branchQuotaUsed"`
	EmployeeQuota     float64 `json:"employeeQuota"`
	EmployeeQuotaUsed float64 `json:"employeeQuotaUsed"`
	// Associations
	Users *[]User `json:"users" gorm:"foreignKey:TenantID"`
}
