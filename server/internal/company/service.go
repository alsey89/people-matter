package company

import (
	"fmt"
	"verve-hrms/internal/schema"
)

type CompanyService struct {
	CompanyRepository *CompanyRepository
}

func NewCompanyService(CompanyRepository *CompanyRepository) *CompanyService {
	return &CompanyService{CompanyRepository: CompanyRepository}
}

func (cs *CompanyService) CreateNewCompanyAndReturnList(newCompany *schema.Company) (*CompanyInterfaceData, error) {
	createdCompany, err := cs.CompanyRepository.CompanyCreate(newCompany)
	if err != nil {
		return nil, fmt.Errorf("company.s.create_company: %w", err)
	}

	var existingCompanies []*schema.Company
	existingCompanies, err = cs.CompanyRepository.CompanyReadAll()
	if err != nil {
		return nil, fmt.Errorf("company.s.fetch_company_list_and_expand_default: %w", err)
	}

	CompanyInterfaceData := &CompanyInterfaceData{
		CompanyList:     existingCompanies,
		ExpandedCompany: createdCompany,
	}

	return CompanyInterfaceData, nil
}

func (cs *CompanyService) GetCompanyListAndExpandDefault() (*CompanyInterfaceData, error) {
	var existingCompanies []*schema.Company
	existingCompanies, err := cs.CompanyRepository.CompanyReadAll()
	if err != nil {
		return nil, fmt.Errorf("company.s.fetch_company_list_and_expand_default: %w", err)
	}

	// Define the Default Company to expand
	var targetCompanyID uint
	targetCompanyID = existingCompanies[0].ID

	// Fetch details for the targeted company
	expandedCompany, err := cs.CompanyRepository.CompanyReadAndExpand(targetCompanyID)
	if err != nil {
		return nil, fmt.Errorf("company.s.fetch_company_list_and_expand_default: %w", err)
	}

	CompanyInterfaceData := &CompanyInterfaceData{
		CompanyList:     existingCompanies,
		ExpandedCompany: expandedCompany,
	}

	return CompanyInterfaceData, nil
}

func (cs *CompanyService) GetCompanyListAndExpandByID(companyID uint) (*CompanyInterfaceData, error) {
	var existingCompanies []*schema.Company
	existingCompanies, err := cs.CompanyRepository.CompanyReadAll()
	if err != nil {
		return nil, fmt.Errorf("company.s.fetch_company_list_and_expand_by_id: %w", err)
	}

	expandedCompany, err := cs.CompanyRepository.CompanyReadAndExpand(companyID)
	if err != nil {
		return nil, fmt.Errorf("company.s.fetch_company_list_and_expand_by_id: %w", err)
	}

	CompanyInterfaceData := &CompanyInterfaceData{
		CompanyList:     existingCompanies,
		ExpandedCompany: expandedCompany,
	}

	return CompanyInterfaceData, nil
}

func (cs *CompanyService) DeleteCompanyAndReturnListAndExpandDefault(companyID uint) (*CompanyInterfaceData, error) {

	err := cs.CompanyRepository.CompanyDelete(companyID)
	if err != nil {
		return nil, fmt.Errorf("company.s.delete_company_by_id: %w", err)
	}

	CompanyInterfaceData, err := cs.GetCompanyListAndExpandDefault()
	if err != nil {
		return nil, fmt.Errorf("company.s.delete_company_by_id: %w", err)
	}

	return CompanyInterfaceData, nil
}
