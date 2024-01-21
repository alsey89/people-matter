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

	TitleCreate(newTitle *schema.Title) (*schema.Title, error)
	TitleRead(TitleID uint) (*schema.Title, error)
	TitleReadAll() ([]*schema.Title, error)
	TitleUpdate(TitleID uint, updateData *schema.Title) (*schema.Title, error)
	TitleDelete(TitleID uint) error

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
	result := cr.client.Preload("Titles").Preload("Departments").Preload("Locations").First(&company, CompanyID)
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
	result := cr.client.Preload("Titles").Preload("Departments").Preload("Locations").Find(&companies)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.company_read_all_and_expand: %w", result.Error)
	}

	return companies, nil
}

func (cr CompanyRepository) CompanyUpdate(CompanyID uint, updateData *schema.Company) (*schema.Company, error) {
	var company schema.Company

	result := cr.client.Model(&company).Where("ID = ?", CompanyID).Updates(updateData)
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
	result := cr.client.Delete(&company, CompanyID)
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

	result := cr.client.Model(&department).Where("ID = ?", DepartmentID).Updates(updateData)
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
	result := cr.client.Delete(&department, DepartmentID)
	if result.Error != nil {
		return fmt.Errorf("company.r.department_delete: %w", result.Error)
	}

	return nil
}

//! Titles      ------------------------------------------------------

func (cr CompanyRepository) TitleCreate(newTitle *schema.Title) (*schema.Title, error) {
	result := cr.client.Create(newTitle)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.title_create: %w", result.Error)
	}

	return newTitle, nil
}

func (cr CompanyRepository) TitleRead(TitleID uint) (*schema.Title, error) {
	var title schema.Title
	result := cr.client.First(&title, TitleID)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.title_read: %w", result.Error)
	}
	return &title, nil
}

func (cr CompanyRepository) TitleReadAll() ([]*schema.Title, error) {
	var titles []*schema.Title
	result := cr.client.Find(&titles)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.title_read_all: %w", result.Error)
	}
	if len(titles) == 0 {
		return nil, fmt.Errorf("company.r.title_read_all: %w", ErrEmptyTable)
	}
	return titles, nil
}

func (cr CompanyRepository) TitleUpdate(TitleID uint, updateData *schema.Title) (*schema.Title, error) {
	var title schema.Title

	result := cr.client.Model(&title).Where("ID = ?", TitleID).Updates(updateData)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.title_update: %w", result.Error)
	}

	// Check if any row was affected, if not, the title does not exist.
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("company.r.title_update: %w", gorm.ErrRecordNotFound)
	}

	return &title, nil
}

func (cr CompanyRepository) TitleDelete(TitleID uint) error {
	var title schema.Title
	result := cr.client.Delete(&title, TitleID)
	if result.Error != nil {
		return fmt.Errorf("company.r.title_delete: %w", result.Error)
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
	return locations, nil
}

func (cr CompanyRepository) LocationUpdate(locationID uint, updateData *schema.Location) (*schema.Location, error) {
	var location schema.Location

	result := cr.client.Model(&location).Where("ID = ?", locationID).Select("IsHeadOffice").Updates(updateData)
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
	result := cr.client.Delete(&location, LocationID)
	if result.Error != nil {
		return fmt.Errorf("company.r.location_delete: %w", result.Error)
	}

	return nil
}
