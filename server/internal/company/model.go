package company

import "github.com/alsey89/hrms/internal/schema"

type CompanyInterfaceData struct {
	CompanyList     []*schema.Company `json:"companyList"`
	ExpandedCompany *schema.Company   `json:"expandedCompany"`
}
