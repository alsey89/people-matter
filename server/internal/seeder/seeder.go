package seeder

import (
	"errors"
	"fmt"

	"github.com/alsey89/people-matter/internal/schema"

	"gorm.io/gorm"
)

func (d *Domain) seedTenants() error {
	db := d.params.DB.GetDB()

	fsps := []schema.Company{
		{
			TenantIdentifier:  "sunming",
			Name:              "Sunming Opthamology",
			LogoURL:           "https://sunming-eye.com.tw/wp-content/uploads/2024/06/logo1.png",
			Email:             "contact@sunming-eye.com.tw",
			Phone:             "(123) 456-7890",
			Website:           "https://sunming-eye.com.tw",
			ContactAddress:    "123 Main Street",
			ContactCity:       "Taipei",
			ContactCountry:    "Taiwan",
			ContactPostalCode: "100",
			BillingAddress:    "123 Main Street",
			BillingCity:       "Taipei",
			BillingCountry:    "Taiwan",
			BillingPostalCode: "100",
			BranchQuota:       3,
			BranchQuotaUsed:   1,
			EmployeeQuota:     50,
			EmployeeQuotaUsed: 25,
		},
		//more clinics
		{
			TenantIdentifier:  "revere",
			Name:              "Revere",
			LogoURL:           "https://revere.com/wp-content/uploads/2024/06/logo1.png",
			Email:             "contact@reverehere.com",
			Phone:             "(123) 456-7890",
			Website:           "https://reverehere.com",
			ContactAddress:    "123 Main Street",
			ContactCity:       "New York",
			ContactCountry:    "USA",
			ContactPostalCode: "10001",
			BillingAddress:    "123 Main Street",
			BillingCity:       "New York",
			BillingCountry:    "USA",
			BillingPostalCode: "10001",
			BranchQuota:       1,
			BranchQuotaUsed:   1,
			EmployeeQuota:     25,
			EmployeeQuotaUsed: 10,
		},
		{
			TenantIdentifier:  "staging",
			Name:              "Staging",
			LogoURL:           "https://staging.com/wp-content/uploads/2024/06/logo1.png",
			Email:             "contact@staging.com",
			Phone:             "(123) 456-7890",
			Website:           "https://staging.com",
			ContactAddress:    "123 Main Street",
			ContactCity:       "New York",
			ContactCountry:    "USA",
			ContactPostalCode: "10001",
			BillingAddress:    "123 Main Street",
			BillingCity:       "New York",
			BillingCountry:    "USA",
			BillingPostalCode: "10001",
			BranchQuota:       1,
			BranchQuotaUsed:   1,
			EmployeeQuota:     25,
			EmployeeQuotaUsed: 5,
		},
	}

	for _, fsp := range fsps {
		var existingFSP schema.Company
		if err := db.Where("tenant_identifier = ?", fsp.TenantIdentifier).First(&existingFSP).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&fsp).Error; err != nil {
					return fmt.Errorf("failed to create Tenant %s: %w", fsp.Name, err)
				}
			} else {
				return fmt.Errorf("failed to check Tenant %s: %w", fsp.Name, err)
			}
		}
	}

	return nil
}
