package company

type NewCompany struct {
	CompanyName string `json:"companyName"`
	CompanySize string `json:"companySize"`

	AdminEmail      string `json:"adminEmail"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}
