package company

import (
	"fmt"

	"github.com/alsey89/people-matter/schema"
	"golang.org/x/crypto/bcrypt"
)

//! Company ------------------------------------------------------------

// Create new company and admin user, send confirmation email
func (d *Domain) CreateNewCompanyAndAdminUser(newCompanyForm *NewCompany) (*schema.User, error) {
	db := d.params.Database.GetDB()

	// terminate if user already exists
	var existingUser schema.User
	if db.Where("email = ?", newCompanyForm.AdminEmail).First(&existingUser).Error == nil {
		return nil, ErrUserExists
	}

	newCompany := schema.Company{
		Name:        newCompanyForm.CompanyName,
		CompanySize: newCompanyForm.CompanySize,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newCompanyForm.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("[CreateNewCompanyAndAdminUser] Error hashing password: %w", err)
	}

	newAdminUser := schema.User{
		Email:    newCompanyForm.AdminEmail,
		Role:     "admin",
		Password: string(hashedPassword),
	}

	// *----- Transaction Start -----
	tx := db.Begin()

	// Create company
	if err := tx.Create(&newCompany).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("[CreateNewCompanyAndAdminUser] Error creating company: %w", err)
	}

	// Assign company ID to the new admin user
	newAdminUser.CompanyID = newCompany.ID

	// Create user
	if err := tx.Create(&newAdminUser).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("[CreateNewCompanyAndAdminUser] Error creating user: %w", err)
	}

	tx.Commit()
	// *----- Transaction End -----

	return &newAdminUser, nil
}

// Get company data without preloading
func (d *Domain) GetCompany(companyID *uint) (*schema.Company, error) {
	db := d.params.Database.GetDB()

	var existingCompany schema.Company

	result := db.First(&existingCompany, companyID)
	if result.Error != nil {
		return nil, fmt.Errorf("[GetCompany] %w", result.Error)
	}

	return &existingCompany, nil
}

// Get company data with preloading
func (d *Domain) GetCompanyWithDetails(companyID *uint) (*schema.Company, error) {
	db := d.params.Database.GetDB()

	var existingCompany schema.Company

	result := db.
		Preload("Departments").
		Preload("Locations").
		Preload("Positions").
		First(&existingCompany, companyID)
	if result.Error != nil {
		return nil, fmt.Errorf("[GetCompanyAndExpand] %w", result.Error)
	}

	return &existingCompany, nil
}

// Update company data
func (d *Domain) UpdateCompany(companyID *uint, newData *schema.Company) error {
	db := d.params.Database.GetDB()

	//using a map instead of struct allows for null values
	dataToUpdate := map[string]interface{}{
		"Name":    newData.Name,
		"LogoURL": newData.LogoURL,
		"Website": newData.Website,
		"Email":   newData.Email,

		"Phone":      newData.Phone,
		"Address":    newData.Address,
		"City":       newData.City,
		"State":      newData.State,
		"Country":    newData.Country,
		"PostalCode": newData.PostalCode,
	}

	result := db.Model(&schema.Company{}).Where("id = ?", companyID).Updates(dataToUpdate)
	if result.Error != nil {
		return fmt.Errorf("[UpdateCompany] %w", result.Error)
	}

	return nil
}

// Delete company data
func (d *Domain) DeleteCompany(companyID *uint) error {
	db := d.params.Database.GetDB()

	result := db.Delete(&schema.Company{}, companyID)
	if result.Error != nil {
		return fmt.Errorf("[DeleteCompany] %w", result.Error)
	}

	return nil
}

// //! Department ------------------------------------------------------------

func (d *Domain) CreateDepartment(newDepartment *schema.Department) error {
	db := d.params.Database.GetDB()

	result := db.Create(newDepartment)
	if result.Error != nil {
		return fmt.Errorf("[CreateDepartment] %w", result.Error)
	}

	return nil
}

func (d *Domain) UpdateDepartment(companyID *uint, departmentID *uint, newData *schema.Department) error {
	db := d.params.Database.GetDB()

	dataToUpdate := map[string]interface{}{
		"Name":        newData.Name,
		"Description": newData.Description,
	}

	result := db.
		Model(&schema.Department{}).
		Where("id = ? AND company_id = ?", departmentID, companyID).
		Updates(dataToUpdate)
	if result.Error != nil {
		return fmt.Errorf("[UpdateDepartment] %w", result.Error)
	}

	return nil
}

func (d *Domain) DeleteDepartment(companyID *uint, departmentID *uint) error {
	db := d.params.Database.GetDB()

	result := db.
		Model(&schema.Department{}).
		Where("company_id = ? AND id = ?", companyID, departmentID).
		Delete(&schema.Department{})
	if result.Error != nil {
		return fmt.Errorf("[DeleteDepartment] %w", result.Error)
	}

	return nil
}

//! Location ------------------------------------------------------------

func (d *Domain) CreateLocation(newLocation *schema.Location) error {
	db := d.params.Database.GetDB()

	result := db.Create(newLocation)
	if result.Error != nil {
		return fmt.Errorf("[CreateLocation] %w", result.Error)
	}

	return nil
}

func (d *Domain) UpdateLocation(companyID *uint, locationID *uint, newData *schema.Location) error {

	db := d.params.Database.GetDB()

	dataToUpdate := map[string]interface{}{
		"Name":         newData.Name,
		"IsHeadOffice": newData.IsHeadOffice,
		"Address":      newData.Address,
		"City":         newData.City,
		"State":        newData.State,
		"Country":      newData.Country,
		"PostalCode":   newData.PostalCode,
	}

	result := db.
		Model(&schema.Location{}).
		Where("company_id = ? AND id = ?", companyID, locationID).
		Updates(dataToUpdate)
	if result.Error != nil {
		return fmt.Errorf("[UpdateLocation] %w", result.Error)
	}

	return nil
}

func (d *Domain) DeleteLocation(companyID *uint, locationID *uint) error {
	db := d.params.Database.GetDB()

	result := db.
		Model(&schema.Location{}).
		Where("company_id = ? AND id = ?", companyID, locationID).
		Delete(&schema.Location{})
	if result.Error != nil {
		return fmt.Errorf("[DeleteLocation] %w", result.Error)
	}

	return nil
}

//! Position ------------------------------------------------------------

func (d *Domain) CreatePosition(newPosition *schema.Position) error {
	db := d.params.Database.GetDB()

	result := db.Create(newPosition)
	if result.Error != nil {
		return fmt.Errorf("[CreatePosition] %w", result.Error)
	}

	return nil
}

func (d *Domain) UpdatePosition(companyID *uint, positionID *uint, newData *schema.Position) error {
	db := d.params.Database.GetDB()

	dataToUpdate := map[string]interface{}{
		"Title":          newData.Title,
		"Description":    newData.Description,
		"Duties":         newData.Duties,
		"Qualifications": newData.Qualifications,
		"Experience":     newData.Experience,
		"MinSalary":      newData.MinSalary,
		"MaxSalary":      newData.MaxSalary,
		"DepartmentID":   newData.DepartmentID,
		"ManagerID":      newData.ManagerID,
	}

	result := db.
		Model(&schema.Position{}).
		Where("company_id = ? AND id = ?", companyID, positionID).
		Updates(dataToUpdate)
	if result.Error != nil {
		return fmt.Errorf("[UpdatePosition] %w", result.Error)
	}

	return nil
}

func (d *Domain) DeletePosition(companyID *uint, positionID *uint) error {
	db := d.params.Database.GetDB()

	result := db.
		Model(&schema.Position{}).
		Where("company_id = ? AND id = ?", companyID, positionID).
		Delete(&schema.Position{})
	if result.Error != nil {
		return fmt.Errorf("[DeletePosition] %w", result.Error)
	}

	return nil
}
