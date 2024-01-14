package company

import "verve-hrms/internal/schema"

type CompanyInterfaceData struct {
	companyList     []*schema.Company
	expandedCompany *schema.Company
}
