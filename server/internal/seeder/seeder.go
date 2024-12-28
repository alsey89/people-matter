package seeder

import (
	"errors"
	"fmt"
	"time"

	"github.com/alsey89/people-matter/internal/schema"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (d *Domain) seedTenants() error {
	db := d.params.DB.GetDB()

	fsps := []schema.Tenant{
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
		var existingFSP schema.Tenant
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

func (d *Domain) seedUsersAndFSPRoles() error {
	db := d.params.DB.GetDB()

	users := []schema.User{
		{
			TenantID:       1,
			Email:          "richard@reverehere.com",
			FirstName:      "Richard",
			LastName:       "Thompson",
			AvatarURL:      "resources/avatar1.png",
			PasswordHash:   "$2a$10$4wT/6yXlp7GnHcOOHQ5JS.4F6E3/K9PB7tC7SA74F0EVZAMdQguU2", // bcrypt hash
			EmailConfirmed: true,
		},
		{
			TenantID:       1,
			Email:          "michael@reverehere.com",
			FirstName:      "Michael",
			LastName:       "Chen",
			AvatarURL:      "resources/avatar2.png",
			PasswordHash:   "$2a$10$4wT/6yXlp7GnHcOOHQ5JS.4F6E3/K9PB7tC7SA74F0EVZAMdQguU2",
			EmailConfirmed: true,
		},
		{
			TenantID:       1,
			Email:          "joey@reverehere.com",
			FirstName:      "Joseph",
			LastName:       "Scully",
			AvatarURL:      "resources/avatar3.png",
			PasswordHash:   "$2a$10$4wT/6yXlp7GnHcOOHQ5JS.4F6E3/K9PB7tC7SA74F0EVZAMdQguU2",
			EmailConfirmed: true,
		},
		{
			TenantID:       4,
			Email:          "cjchen@reverehere.com",
			FirstName:      "CJ",
			LastName:       "Chen",
			AvatarURL:      "resources/avatar2.png",
			PasswordHash:   "$2a$10$4wT/6yXlp7GnHcOOHQ5JS.4F6E3/K9PB7tC7SA74F0EVZAMdQguU2",
			EmailConfirmed: true,
		},
	}

	for _, user := range users {
		var existingUser *schema.User
		var err error

		existingUser, err = d.params.Identity.FindUserByEmail(nil, user.TenantID, user.Email)
		if err != nil {
			return fmt.Errorf("seedUsers: %w", err)
		}

		if existingUser == nil {
			// Create user and Tenant role
			d.params.Identity.CreateUserAndFSPRole(
				nil,                      // db *gorm.DB
				user.TenantID,            // TenantID uint
				user.FirstName,           // firstName string
				user.LastName,            // lastName string
				user.Email,               // email string
				&user.PasswordHash,       // passwordHash *string
				schema.RoleFSPSuperAdmin, // startingRole schema.FSPRoleConst
			)
			// Update name, avatar, and email confirmed
			db.Model(&schema.User{}).Where("email = ? AND fsp_id = ?", user.Email, user.TenantID).Updates(schema.User{
				FirstName:      user.FirstName,
				LastName:       user.LastName,
				AvatarURL:      user.AvatarURL,
				EmailConfirmed: user.EmailConfirmed,
			})
		}
	}

	return nil
}

func (d *Domain) seedMemorials() error {
	// db := d.params.DB.GetDB()

	memorials := []schema.Memorial{
		{
			TenantID:   1,
			Title:      "Richard Thompson's Memorial",
			Identifier: "richard_thompson_1980-01-01",
			FirstName:  "Richard",
			LastName:   "Thompson",
			DOB:        func(t time.Time) *time.Time { return &t }(time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			TenantID:   1,
			Title:      "Michael Chen's Memorial",
			Identifier: "michael_chen_1985-01-01",
			FirstName:  "Michael",
			LastName:   "Chen",
			DOB:        func(t time.Time) *time.Time { return &t }(time.Date(1985, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
	}

	for _, memorial := range memorials {
		var existingMemorial *schema.Memorial
		var err error

		existingMemorial, err = d.params.Memorial.FindMemorialByIdentifier(
			memorial.TenantID,
			memorial.Identifier,
		)
		if err != nil {
			d.logger.Error("seedMemorials: %w", zap.Error(err))
		}

		if existingMemorial == nil {
			// Create memorial
			_, err = d.params.Tenant.CreateMemorial(
				nil,
				memorial.TenantID,   // TenantID 		uint
				memorial.FirstName,  // firstName 	string
				memorial.LastName,   // lastName 	string
				memorial.Identifier, // identifier 	string
				memorial.DOB,        // DOB 		*time.Time
				memorial.DOD,        // DOD 		*time.Time
			)
			if err != nil {
				d.logger.Error("seedMemorials: %w", zap.Error(err))
			}
		}
	}

	return nil
}

func (d *Domain) seedUserMemorialRoles() error {
	type MemorialCurator struct {
		TenantID     uint
		Email        string
		MemorialID   uint
		Role         schema.MemorialRoleConst
		Relationship schema.RelationshipConst
	}

	// Define the memorial-curator pairs
	memorialCurators := []MemorialCurator{
		{
			TenantID:     1,
			Email:        "richard@reverehere.com",
			MemorialID:   1,
			Role:         schema.RoleMemCurator,
			Relationship: schema.RelationshipSelf,
		},
		{
			TenantID:     1,
			Email:        "michael@reverehere.com",
			MemorialID:   2,
			Role:         schema.RoleMemCurator,
			Relationship: schema.RelationshipSelf,
		},
		{
			TenantID:     1,
			Email:        "michael@reverhere.com",
			MemorialID:   1,
			Role:         schema.RoleMemContributor,
			Relationship: schema.RelationshipFriend,
		},
		{
			TenantID:     1,
			Email:        "joey@reverehere.com",
			MemorialID:   2,
			Role:         schema.RoleMemContributor,
			Relationship: schema.RelationshipFriend,
		},
		// Add more memorial-curator pairs as needed
	}

	var seedErrors []error

	// Start a database transaction
	err := d.params.DB.GetDB().Transaction(func(tx *gorm.DB) error {
		for _, curator := range memorialCurators {
			// Fetch the user by email
			user, err := d.params.Identity.FindUserByEmail(tx, curator.TenantID, curator.Email)
			if err != nil {
				seedErrors = append(seedErrors, fmt.Errorf("failed to find user by email %s: %w", curator.Email, err))
				continue
			}

			if user == nil {
				seedErrors = append(seedErrors, fmt.Errorf("user with email %s not found", curator.Email))
				continue
			}

			// Check if the user already has the curator role for the memorial
			hasRole, err := d.params.Identity.UserHasMemorialRole(
				curator.TenantID,
				user.ID,
				curator.MemorialID,
				curator.Role,
			)
			if err != nil {
				seedErrors = append(seedErrors, fmt.Errorf("failed to check role for user %s on memorial %d: %w", curator.Email, curator.MemorialID, err))
				continue
			}

			if hasRole {
				continue
			}

			// Assign the curator role
			_, err = d.params.Identity.AssignOrUpdateMemorialRole(
				tx,                    // db 			*gorm.DB
				curator.TenantID,      // TenantID 		uint
				user.ID,               // userID 		uint
				curator.MemorialID,    // memorialID 	uint
				curator.Role,          // role 			schema.MemorialRoleConst
				&curator.Relationship, // relationship 	*schema.RelationshipConst
			)
			if err != nil {
				seedErrors = append(seedErrors, fmt.Errorf("failed to assign role to user %s for memorial %d: %w", curator.Email, curator.MemorialID, err))
				continue
			}
		}
		return nil
	})

	if err != nil {
		// Transaction-level error
		return fmt.Errorf("transaction failed: %w", err)
	}

	if len(seedErrors) > 0 {
		// Aggregate errors
		var combinedErrorMsg string
		for _, e := range seedErrors {
			combinedErrorMsg += e.Error() + "; "
		}
		return fmt.Errorf("seedUserMemorialRoles encountered errors: %s", combinedErrorMsg)
	}

	return nil
}

func (d *Domain) seedContributorApplications() error {
	db := d.params.DB.GetDB()

	applications := []schema.Application{
		{
			TenantID:        1,
			MemorialID:      2,
			ApplicantID:     3,
			ApplicationType: schema.ApplicationTypeContributor,
			Relationship:    schema.RelationshipBrother,
			Status:          schema.ApplicationStatusPending,
			AppliedOn:       time.Now(),
		},
		{
			TenantID:        1,
			MemorialID:      2,
			ApplicantID:     1,
			ApplicationType: schema.ApplicationTypeContributor,
			Relationship:    schema.RelationshipFather,
			Status:          schema.ApplicationStatusPending,
			AppliedOn:       time.Now(),
		},
	}

	for _, application := range applications {
		var existingApplication schema.Application
		if err := db.Where("fsp_id = ? AND memorial_id = ? AND applicant_id = ?", application.TenantID, application.MemorialID, application.ApplicantID).First(&existingApplication).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&application).Error; err != nil {
					return fmt.Errorf("failed to create application for applicant_id %d: %w", application.ApplicantID, err)
				}
			} else {
				return fmt.Errorf("failed to check application for applicant_id %d: %w", application.ApplicantID, err)
			}
		}
	}

	return nil
}

func (d *Domain) seedContributorInvitations() error {

	michaelChenMemorialID := uint(2)

	invitations := []schema.Invitation{
		{
			TenantID:       1,
			InvitationType: schema.InvitationTypeMemorial,
			InviterID:      2,
			InviteeEmail:   "johndoe1231@mailnesia.com",
			Relationship:   schema.RelationshipFriend,
			Status:         schema.InvitationStatusPending,
			InvitedOn:      time.Now(),
			ExpiresOn:      time.Now().AddDate(0, 0, 7),
			Token:          uuid.NewString(),
			MemorialID:     &michaelChenMemorialID,
			MemorialRole:   schema.RoleMemContributor,
		},
		{
			TenantID:       1,
			InvitationType: schema.InvitationTypeMemorial,
			InviterID:      2,
			InviteeEmail:   "johndoe1232@mailnesia.com",
			Relationship:   schema.RelationshipFriend,
			Status:         schema.InvitationStatusPending,
			InvitedOn:      time.Now(),
			ExpiresOn:      time.Now().AddDate(0, 0, 7),
			Token:          uuid.NewString(),
			MemorialID:     &michaelChenMemorialID,
			MemorialRole:   schema.RoleMemContributor,
		},
	}

	for _, invitation := range invitations {
		if invitation.MemorialID == nil {
			continue
		}

		existingInvitation, err := d.params.Identity.FindInvitationByEmail(
			invitation.TenantID,     // TenantID uint
			*invitation.MemorialID,  // memorialID *uint
			invitation.InviteeEmail, // email string
		)
		if err != nil {
			d.logger.Error("seedContributorInvitations: ", zap.Error(err))
		}

		if existingInvitation == nil {
			d.params.Identity.CreateMemorialContributorInvitation(
				nil,
				invitation.TenantID,     // TenantID 			uint
				invitation.InviterID,    // inviterID 		uint
				invitation.InviteeEmail, // inviteeEmail 	string
				invitation.Relationship, // relationship 	schema.RelationshipConst
				invitation.MemorialID,   // memorialID 		*uint
			)
		}

	}
	return nil
}
