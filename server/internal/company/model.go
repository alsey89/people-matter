package company

import "verve-hrms/internal/schema"

type CompanyInterfaceData struct {
	CompanyList     []*schema.Company `json:"companyList"`
	ExpandedCompany *schema.Company   `json:"expandedCompany"`
}
