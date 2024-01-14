package company

import (
	"fmt"
	"verve-hrms/internal/schema"

	"gorm.io/gorm"
)

type Repository interface {
	Create(newUser *schema.User) (*schema.User, error)
	Read(UserID uint) (*schema.User, error)
	ReadByEmail(email string) (*schema.User, error)
	Update(UserID uint, updateData *schema.User) (*schema.User, error)
	Delete(UserID uint) error
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
		return nil, fmt.Errorf("company.r.create: %w", result.Error)
	}

	newCompany.ID = uint(result.RowsAffected)

	return newCompany, nil
}

func (cr CompanyRepository) CompanyRead(CompanyID uint) (*schema.Company, error) {
	var company schema.Company
	result := cr.client.First(&company, CompanyID)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.read: %w", result.Error)
	}

	return &company, nil
}

func (cr CompanyRepository) CompanyReadAndExpand(CompanyID uint) (*schema.Company, error) {
	var company schema.Company
	result := cr.client.Preload("Titles").Preload("Departments").Preload("Locations").First(&company, CompanyID)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.read_and_expand: %w", result.Error)
	}

	return &company, nil
}

func (cr CompanyRepository) CompanyReadAll() ([]*schema.Company, error) {
	var companies []*schema.Company
	result := cr.client.Find(&companies)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.read_all: %w", result.Error)
	}
	if len(companies) == 0 {
		return nil, fmt.Errorf("company.r.read_all: %w", ErrEmptyTable)
	}

	return companies, nil
}

func (cr CompanyRepository) CompanyReadAllAndExpand() ([]*schema.Company, error) {
	var companies []*schema.Company
	result := cr.client.Preload("Titles").Preload("Departments").Preload("Locations").Find(&companies)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.read_all_and_expand: %w", result.Error)
	}

	return companies, nil
}

func (cr CompanyRepository) CompanyUpdate(CompanyID uint, updateData *schema.Company) (*schema.Company, error) {
	var company schema.Company
	result := cr.client.First(&company, CompanyID)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.update: %w", result.Error)
	}

	result = cr.client.Model(&company).Updates(updateData)
	if result.Error != nil {
		return nil, fmt.Errorf("company.r.update: %w", result.Error)
	}

	return &company, nil
}

func (cr CompanyRepository) CompanyDelete(CompanyID uint) error {
	var company schema.Company
	result := cr.client.Delete(&company, CompanyID)
	if result.Error != nil {
		return fmt.Errorf("company.r.delete: %w", result.Error)
	}

	return nil
}

//! Departments ------------------------------------------------------

//! Titles      ------------------------------------------------------

//! Locations   ------------------------------------------------------
