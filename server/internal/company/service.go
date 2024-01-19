package company

import (
	"fmt"
	"log"
	"verve-hrms/internal/schema"
)

type CompanyService struct {
	CompanyRepository *CompanyRepository
}

func NewCompanyService(CompanyRepository *CompanyRepository) *CompanyService {
	return &CompanyService{CompanyRepository: CompanyRepository}
}

//! Company ------------------------------------------------------------

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

	expandedCompany, err := cs.CompanyRepository.CompanyReadAndExpand(companyID)
	if err != nil {
		return nil, fmt.Errorf("company.s.fetch_company_list_and_expand_by_id: %w", err)
	}

	var existingCompanies []*schema.Company
	existingCompanies, err = cs.CompanyRepository.CompanyReadAll()
	if err != nil {
		return nil, fmt.Errorf("company.s.fetch_company_list_and_expand_by_id: %w", err)
	}

	CompanyInterfaceData := &CompanyInterfaceData{
		CompanyList:     existingCompanies,
		ExpandedCompany: expandedCompany,
	}

	return CompanyInterfaceData, nil
}

func (cs *CompanyService) UpdateCompanyAndReturnCompanyListAndExpandID(companyId uint, newData *schema.Company) (*CompanyInterfaceData, error) {

	updatedCompany, err := cs.CompanyRepository.CompanyUpdate(companyId, newData)
	if err != nil {
		return nil, fmt.Errorf("company.s.update_company: %w", err)
	}

	companyID := updatedCompany.ID

	companyData, err := cs.GetCompanyListAndExpandByID(companyID)
	if err != nil {
		return nil, fmt.Errorf("company.s.create_title_and_return_company_list_and_expand_id: %w", err)
	}

	return companyData, nil
}

func (cs *CompanyService) DeleteCompanyAndReturnCompanyListAndExpandDefault(companyID uint) (*CompanyInterfaceData, error) {

	err := cs.CompanyRepository.CompanyDelete(companyID)
	if err != nil {
		return nil, fmt.Errorf("company.s.delete_company_by_id: %w", err)
	}

	companyData, err := cs.GetCompanyListAndExpandDefault()
	if err != nil {
		return nil, fmt.Errorf("company.s.delete_company_by_id: %w", err)
	}

	return companyData, nil
}

//! Department ------------------------------------------------------------

func (cs *CompanyService) CreateNewDepartmentAndReturnCompanyListAndExpandID(companyID uint, newDepartment *schema.Department) (*CompanyInterfaceData, error) {
	newDepartment.CompanyID = companyID

	_, err := cs.CompanyRepository.DepartmentCreate(newDepartment)
	if err != nil {
		return nil, fmt.Errorf("company.s.create_department_and_return_company_list_and_expand_id: %w", err)
	}

	companyData, err := cs.GetCompanyListAndExpandByID(companyID)
	if err != nil {
		return nil, fmt.Errorf("company.s.create_title_and_return_company_list_and_expand_id: %w", err)
	}

	return companyData, nil
}

func (cs *CompanyService) UpdateDepartmentAndReturnCompanyListAndExpandID(companyID uint, departmentID uint, dataToUpdate *schema.Department) (*CompanyInterfaceData, error) {
	dataToUpdate.CompanyID = companyID

	_, err := cs.CompanyRepository.DepartmentUpdate(departmentID, dataToUpdate)
	if err != nil {
		return nil, fmt.Errorf("company.s.update_department_and_return_company_list_and_expand_id: %w", err)
	}

	companyData, err := cs.GetCompanyListAndExpandByID(companyID)
	if err != nil {
		return nil, fmt.Errorf("company.s.create_title_and_return_company_list_and_expand_id: %w", err)
	}

	return companyData, nil
}

func (cs *CompanyService) DeleteDepartmentAndReturnCompanyListAndExpandID(companyID uint, departmentID uint) (*CompanyInterfaceData, error) {
	err := cs.CompanyRepository.DepartmentDelete(departmentID)
	if err != nil {
		return nil, fmt.Errorf("company.s.delete_department_and_return_company_list_and_expand_id: %w", err)
	}

	companyData, err := cs.GetCompanyListAndExpandByID(companyID)
	if err != nil {
		return nil, fmt.Errorf("company.s.create_title_and_return_company_list_and_expand_id: %w", err)
	}

	return companyData, nil
}

//! Title ------------------------------------------------------------

func (cs *CompanyService) CreateNewTitleAndReturnCompanyListAndExpandID(companyID uint, newTitle *schema.Title) (*CompanyInterfaceData, error) {
	newTitle.CompanyID = companyID

	_, err := cs.CompanyRepository.TitleCreate(newTitle)
	if err != nil {
		return nil, fmt.Errorf("company.s.create_title_and_return_company_list_and_expand_id: %w", err)
	}

	companyData, err := cs.GetCompanyListAndExpandByID(companyID)
	if err != nil {
		return nil, fmt.Errorf("company.s.create_title_and_return_company_list_and_expand_id: %w", err)
	}

	return companyData, nil
}

func (cs *CompanyService) UpdateTitleAndReturnCompanyListAndExpandID(companyID uint, titleID uint, dataToUpdate *schema.Title) (*CompanyInterfaceData, error) {
	dataToUpdate.ID = titleID

	_, err := cs.CompanyRepository.TitleUpdate(titleID, dataToUpdate)
	if err != nil {
		return nil, fmt.Errorf("company.s.update_title_and_return_company_list_and_expand_id: %w", err)
	}

	companyData, err := cs.GetCompanyListAndExpandByID(companyID)
	if err != nil {
		return nil, fmt.Errorf("company.s.create_title_and_return_company_list_and_expand_id: %w", err)
	}

	return companyData, nil
}

func (cs *CompanyService) DeleteTitleAndReturnCompanyListAndExpandID(companyID uint, titleID uint) (*CompanyInterfaceData, error) {
	err := cs.CompanyRepository.TitleDelete(titleID)
	if err != nil {
		return nil, fmt.Errorf("company.s.delete_title_and_return_company_list_and_expand_id: %w", err)
	}

	companyData, err := cs.GetCompanyListAndExpandByID(companyID)
	if err != nil {
		return nil, fmt.Errorf("company.s.create_title_and_return_company_list_and_expand_id: %w", err)
	}

	return companyData, nil
}

//! Location ------------------------------------------------------------

func (cs *CompanyService) CreateNewLocationAndReturnCompanyListAndExpandID(companyID uint, newLocation *schema.Location) (*CompanyInterfaceData, error) {
	newLocation.CompanyID = companyID

	log.Printf("newLocation: %v", newLocation)

	_, err := cs.CompanyRepository.LocationCreate(newLocation)
	if err != nil {
		return nil, fmt.Errorf("company.s.create_location_and_return_company_list_and_expand_id: %w", err)
	}

	companyData, err := cs.GetCompanyListAndExpandByID(companyID)
	if err != nil {
		return nil, fmt.Errorf("company.s.create_location_and_return_company_list_and_expand_id: %w", err)
	}

	return companyData, nil
}

func (cs *CompanyService) UpdateLocationAndReturnCompanyListAndExpandID(companyID uint, locationID uint, dataToUpdate *schema.Location) (*CompanyInterfaceData, error) {
	dataToUpdate.ID = locationID

	_, err := cs.CompanyRepository.LocationUpdate(locationID, dataToUpdate)
	if err != nil {
		return nil, fmt.Errorf("company.s.update_location_and_return_company_list_and_expand_id: %w", err)
	}

	companyData, err := cs.GetCompanyListAndExpandByID(companyID)
	if err != nil {
		return nil, fmt.Errorf("company.s.update_location_and_return_company_list_and_expand_id: %w", err)
	}

	return companyData, nil
}

func (cs *CompanyService) DeleteLocationAndReturnCompanyListAndExpandID(companyID uint, locationID uint) (*CompanyInterfaceData, error) {
	err := cs.CompanyRepository.LocationDelete(locationID)
	if err != nil {
		return nil, fmt.Errorf("company.s.delete_location_and_return_company_list_and_expand_id: %w", err)
	}

	companyData, err := cs.GetCompanyListAndExpandByID(companyID)
	if err != nil {
		return nil, fmt.Errorf("company.s.delete_location_and_return_company_list_and_expand_id: %w", err)
	}

	return companyData, nil
}
