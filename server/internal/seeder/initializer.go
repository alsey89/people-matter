package seeder

import (
	"errors"
	"fmt"

	"github.com/alsey89/people-matter/internal/schema"

	"gorm.io/gorm"
)

func (d *Domain) initializeFSPRoles() error {
	db := d.params.DB.GetDB()

	predefinedFSPRoles := []schema.FSPRole{
		{
			Name:        schema.RoleFSPSuperAdmin,
			Description: "Super Admin at FSP Level",
		},
		{
			Name:        schema.RoleFSPAdmin,
			Description: "Admin at FSP Level",
		},
		{
			Name:        schema.RoleFSPUser,
			Description: "User at FSP Level",
		},
	}

	for _, role := range predefinedFSPRoles {
		var existingRole schema.FSPRole
		if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&role).Error; err != nil {
					return fmt.Errorf("failed to create FSP role %s: %w", role.Name, err)
				}
			} else {
				return fmt.Errorf("failed to check FSP role %s: %w", role.Name, err)
			}
		}
	}

	return nil
}

func (d *Domain) initializeMemorialRoles() error {
	db := d.params.DB.GetDB()

	predefinedMemorialRoles := []schema.MemorialRole{
		{
			Name:        schema.RoleMemSelf,
			Description: "Curator is self",
		},
		{
			Name:        schema.RoleMemCurator,
			Description: "Curator at Memorial Level",
		},
		{
			Name:        schema.RoleMemContributor,
			Description: "Contributor at Memorial Level",
		},
		{
			Name:        schema.RoleMemInsitu,
			Description: "In Situ at Memorial Level",
		},
	}

	for _, role := range predefinedMemorialRoles {
		var existingRole schema.MemorialRole
		if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&role).Error; err != nil {
					return fmt.Errorf("failed to create Memorial role %s: %w", role.Name, err)
				}
			} else {
				return fmt.Errorf("failed to check Memorial role %s: %w", role.Name, err)
			}
		}
	}

	return nil
}

func (d *Domain) intializeCountries() error {
	db := d.params.DB.GetDB()

	countries := []schema.Country{
		{
			// BaseModel: schema.BaseModel{ID: 1},
			Name: "United States",
			Code: "US"},
		{
			// BaseModel: schema.BaseModel{ID: 2},
			Name: "China",
			Code: "CH"},
		{
			// BaseModel: schema.BaseModel{ID: 3},
			Name: "Mexico",
			Code: "MX"},
	}

	for _, country := range countries {
		var existingCountry schema.Country
		if err := db.Where("name = ?", country.Name).First(&existingCountry).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&country).Error; err != nil {
					return fmt.Errorf("failed to create country %s: %w", country.Name, err)
				}
			} else {
				return fmt.Errorf("failed to check country %s: %w", country.Name, err)
			}
		}
	}
	return nil
}

func (d *Domain) initializeStateProvinces() error {
	db := d.params.DB.GetDB()

	stateProvinces := []schema.StateProvince{
		{
			// BaseModel: schema.BaseModel{ID: 1},
			Name:      "California",
			Code:      "CA",
			CountryID: 1},
		{
			// BaseModel: schema.BaseModel{ID: 2},
			Name:      "Texas",
			Code:      "TX",
			CountryID: 1},
		{
			// BaseModel: schema.BaseModel{ID: 3},
			Name:      "Illinois",
			Code:      "IL",
			CountryID: 1},
		{
			// BaseModel: schema.BaseModel{ID: 4},
			Name:      "Beijing",
			Code:      "BJ",
			CountryID: 2},
		{
			// BaseModel: schema.BaseModel{ID: 5},
			Name:      "Shanghai",
			Code:      "SH",
			CountryID: 2},
		{
			// BaseModel: schema.BaseModel{ID: 6},
			Name:      "Mexico City",
			Code:      "MC",
			CountryID: 3},
		{
			// BaseModel: schema.BaseModel{ID: 7},
			Name:      "Guadalajara",
			Code:      "GJ",
			CountryID: 3},
	}

	for _, stateProvince := range stateProvinces {
		var existingStateProvince schema.StateProvince
		if err := db.Where("name = ?", stateProvince.Name).First(&existingStateProvince).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&stateProvince).Error; err != nil {
					return fmt.Errorf("failed to create state province %s: %w", stateProvince.Name, err)
				}
			} else {
				return fmt.Errorf("failed to check state province %s: %w", stateProvince.Name, err)
			}
		}
	}

	return nil
}
