package company

import (
	"fmt"
	"verve-hrms/internal/schema"

	"gorm.io/gorm"
)

type Repository interface {
	CompanyCreate(newCompany *schema.Company) (*schema.Company, error)
	CompanyRead(CompanyID uint) (*schema.Company, error)
	CompanyReadAndExpand(CompanyID uint) (*schema.Company, error)
	CompanyReadAll() ([]*schema.Company, error)
	CompanyReadAndExpandAll() ([]*schema.Company, error)
	CompanyUpdate(CompanyID uint, updateData *schema.Company) (*schema.Company, error)
	CompanyDelete(CompanyID uint) error

	DepartmentCreate(newDepartment *schema.Department) (*schema.Department, error)
	DepartmentRead(DepartmentID uint) (*schema.Department, error)
	DepartmentReadAll() ([]*schema.Department, error)
	DepartmentUpdate(DepartmentID uint, updateData *schema.Department) (*schema.Department, error)
	DepartmentDelete(DepartmentID uint) error

	LocationCreate(newLocation *schema.Location) (*schema.Location, error)
	LocationRead(LocationID uint) (*schema.Location, error)
	LocationReadAll() ([]*schema.Location, error)
	LocationUpdate(LocationID uint, updateData *schema.Location) (*schema.Location, error)
	LocationDelete(LocationID uint) error
}

type CompanyRepository struct {
	client *gorm.DB
}

func NewCompanyRepository(client *gorm.DB) *CompanyRepository {
	return &CompanyRepository{client: client}
}

//! Company     ------------------------------------------------------

func (cr CompanyRepository) CompanyCreate(newCompany *schema.Company) (*schema.Company, error) {
	result := cr.client.Create(newCompany)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.company_create: %w", result.Error)
	}

	return newCompany, nil
}

func (cr CompanyRepository) CompanyRead(CompanyID uint) (*schema.Company, error) {
	var company schema.Company
	result := cr.client.First(&company, CompanyID)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.company_read: %w", result.Error)
	}

	return &company, nil
}

func (cr CompanyRepository) CompanyReadAndExpand(CompanyID uint) (*schema.Company, error) {
	var company schema.Company
	result := cr.client.Preload("Departments").Preload("Locations").Preload("Jobs").First(&company, CompanyID)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.company_read_and_expand: %w", result.Error)
	}

	return &company, nil
}

func (cr CompanyRepository) CompanyReadAll() ([]*schema.Company, error) {
	var companies []*schema.Company
	result := cr.client.Find(&companies)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.company_read_all: %w", result.Error)
	}
	if len(companies) == 0 {
		return nil, fmt.Errorf("company.r.company_read_all: %w", ErrEmptyTable)
	}

	return companies, nil
}

func (cr CompanyRepository) CompanyReadAndExpandAll() ([]*schema.Company, error) {
	var companies []*schema.Company
	result := cr.client.Preload("Departments").Preload("Locations").Preload("Jobs").Find(&companies)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.company_read_all_and_expand: %w", result.Error)
	}
	if len(companies) == 0 {
		return nil, fmt.Errorf("company.r.company_read_all_and_expand: %w", ErrEmptyTable)
	}

	return companies, nil
}

func (cr CompanyRepository) CompanyUpdate(CompanyID uint, updateData *schema.Company) (*schema.Company, error) {
	var company schema.Company

	updateMap := map[string]interface{}{
		"Name":       updateData.Name,
		"LogoURL":    updateData.LogoURL,
		"Website":    updateData.Website,
		"Email":      updateData.Email,
		"Phone":      updateData.Phone,
		"Address":    updateData.Address,
		"City":       updateData.City,
		"State":      updateData.State,
		"Country":    updateData.Country,
		"PostalCode": updateData.PostalCode,
	}

	result := cr.client.Model(&company).Where("ID = ?", CompanyID).Updates(updateMap)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.company_update: %w", result.Error)
	}

	// Check if any row was affected, if not, the company does not exist.
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("company.r.company_update: %w", gorm.ErrRecordNotFound)
	}

	return &company, nil
}

func (cr CompanyRepository) CompanyDelete(CompanyID uint) error {
	var company schema.Company
	result := cr.client.Unscoped().Delete(&company, CompanyID)
	if result.Error != nil {
		return fmt.Errorf("company.r.company_delete: %w", result.Error)
	}

	return nil
}

//! Departments ------------------------------------------------------

func (cr CompanyRepository) DepartmentCreate(newDepartment *schema.Department) (*schema.Department, error) {
	result := cr.client.Create(newDepartment)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.create: %w", result.Error)
	}

	return newDepartment, nil
}

func (cr CompanyRepository) DepartmentRead(DepartmentID uint) (*schema.Department, error) {
	var department schema.Department
	result := cr.client.First(&department, DepartmentID)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.department_read: %w", result.Error)
	}
	return &department, nil
}

func (cr CompanyRepository) DepartmentReadAll() ([]*schema.Department, error) {
	var departments []*schema.Department
	result := cr.client.Find(&departments)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.department_read_all: %w", result.Error)
	}
	if len(departments) == 0 {
		return nil, fmt.Errorf("company.r.department_read_all: %w", ErrEmptyTable)
	}
	return departments, nil
}

func (cr CompanyRepository) DepartmentUpdate(DepartmentID uint, updateData *schema.Department) (*schema.Department, error) {
	var department schema.Department

	updateMap := map[string]interface{}{
		// "CompanyID": updateData.CompanyID, //* not part of updateData, will always initialize to 0 since it's not a pointer
		"Name":        updateData.Name,
		"Description": updateData.Description,
	}

	result := cr.client.Model(&department).Where("ID = ?", DepartmentID).Updates(updateMap)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.department_update: %w", result.Error)
	}

	// Check if any row was affected, if not, the department does not exist.
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("company.r.department_update: %w", gorm.ErrRecordNotFound)
	}

	return &department, nil
}

func (cr CompanyRepository) DepartmentDelete(DepartmentID uint) error {
	var department schema.Department
	result := cr.client.Unscoped().Delete(&department, DepartmentID)
	if result.Error != nil {
		return fmt.Errorf("company.r.department_delete: %w", result.Error)
	}

	return nil
}

//! Locations   ------------------------------------------------------

func (cr CompanyRepository) LocationCreate(newLocation *schema.Location) (*schema.Location, error) {
	result := cr.client.Create(newLocation)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.location_create: %w", result.Error)
	}

	return newLocation, nil
}

func (cr CompanyRepository) LocationRead(LocationID uint) (*schema.Location, error) {
	var location schema.Location
	result := cr.client.First(&location, LocationID)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.location_read: %w", result.Error)
	}
	return &location, nil
}

func (cr CompanyRepository) LocationReadAll() ([]*schema.Location, error) {
	var locations []*schema.Location
	result := cr.client.Find(&locations)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.location_read_all: %w", result.Error)
	}
	if len(locations) == 0 {
		return nil, fmt.Errorf("company.r.location_read_all: %w", ErrEmptyTable)
	}

	return locations, nil
}

func (cr CompanyRepository) LocationUpdate(locationID uint, updateData *schema.Location) (*schema.Location, error) {
	var location schema.Location

	updateMap := map[string]interface{}{
		// "CompanyID":    updateData.CompanyID, //* not part of updateData, will always initialize to 0 since it's not a pointer
		"Name":         updateData.Name,
		"IsHeadOffice": updateData.IsHeadOffice,
		"Phone":        updateData.Phone,
		"Address":      updateData.Address,
		"City":         updateData.City,
		"State":        updateData.State,
		"Country":      updateData.Country,
		"PostalCode":   updateData.PostalCode,
	}

	result := cr.client.Model(&location).Where("ID = ?", locationID).Updates(updateMap)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.location_update: %w", result.Error)
	}

	// Check if any row was affected, if not, the location does not exist.
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("company.r.location_update: %w", gorm.ErrRecordNotFound)
	}

	return &location, nil
}

func (cr CompanyRepository) LocationDelete(LocationID uint) error {
	var location schema.Location
	result := cr.client.Unscoped().Delete(&location, LocationID)
	if result.Error != nil {
		return fmt.Errorf("company.r.location_delete: %w", result.Error)
	}

	return nil
}
