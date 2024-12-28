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
			Description: "Super Admin at Tenant Level",
		},
		{
			Name:        schema.RoleFSPAdmin,
			Description: "Admin at Tenant Level",
		},
		{
			Name:        schema.RoleFSPUser,
			Description: "User at Tenant Level",
		},
	}

	for _, role := range predefinedFSPRoles {
		var existingRole schema.FSPRole
		if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&role).Error; err != nil {
					return fmt.Errorf("failed to create Tenant role %s: %w", role.Name, err)
				}
			} else {
				return fmt.Errorf("failed to check Tenant role %s: %w", role.Name, err)
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
