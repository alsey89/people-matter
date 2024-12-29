package company

import (
	"errors"
	"fmt"

	"github.com/alsey89/people-matter/internal/schema"

	"gorm.io/gorm"
)

// Form Data ------------------------------------------------------

// Account ---------------------------------------------------------

func (d *Domain) GetAccountService(TenantID uint, preloadDetails bool) (*schema.Company, error) {
	db := d.params.DB.GetDB()

	existingFSP := schema.Company{}

	query := db.Where("id = ?", TenantID)
	if preloadDetails {
		query = query.Preload("Country").Preload("StateProvince")
	}

	err := query.First(&existingFSP).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("GetAccountService: no record found for TenantID %d", TenantID)
		}
		return nil, fmt.Errorf("GetAccountService: %w", err)
	}

	return &existingFSP, nil
}

func (d *Domain) updateAccountService(TenantID uint, updatedFSP schema.Company) error {
	db := d.params.DB.GetDB()

	err := db.
		Where("id = ?", TenantID).
		// Select(
		// 	"Name",
		// 	"LogoURL",
		// 	"FSPType",
		// 	"BusinessTypeID",
		// 	"EmployeeCount",
		// 	"ParentCompany",
		// 	"Subsidiaries",
		// 	"Email",
		// 	"Phone",
		// 	"Website",
		// 	"Address",
		// 	"PostalCode",
		// 	"CountryID",
		// 	"StateProvinceID",
		// 	"BillingAddress",
		// ).
		Updates(&updatedFSP).
		Error
	if err != nil {
		return fmt.Errorf("updateAccountService: %w", err)
	}

	return nil
}
