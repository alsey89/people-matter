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

func (d *Domain) seedFSPs() error {
	db := d.params.DB.GetDB()

	fsps := []schema.FSP{
		{
			TenantIdentifier:  "sci",
			Name:              "SCI",
			LogoURL:           "https://www.sci-corp.com/dfsmedia/042808e1630c49a48950d5077d6556eb/33440-50075",
			FSPType:           schema.FSPTypeCemetery,
			BusinessTypeID:    1,
			BusinessType:      schema.BusinessTypeLLC,
			CRN:               "LLC1234567",
			EIN:               "12-3456789",
			BillingAddress:    "123 Main Street, Springfield",
			Address:           "123 Main Street, Springfield",
			PostalCode:        "62701",
			StateProvinceID:   3,
			CountryID:         1,
			Phone:             "(123) 456-7890",
			Website:           "http://sci.com",
			Established:       "1984",
			EmployeeCount:     "500-2500",
			ParentCompany:     "SCI Global",
			MemorialQuota:     20000,
			MemorialQuotaUsed: 10000,
			StorageQuota:      2,
			StorageQuotaUsed:  1,
		},
		{
			TenantIdentifier:  "monumentalists",
			Name:              "Monumentalists",
			LogoURL:           "resources/brand2.png",
			FSPType:           schema.FSPTypeMonument,
			BusinessTypeID:    2,
			BusinessType:      schema.BusinessTypeCorporation,
			CRN:               "C1234567",
			EIN:               "98-7654321",
			BillingAddress:    "456 Elm Avenue, Dallas",
			Address:           "456 Elm Avenue, Dallas",
			PostalCode:        "75201",
			StateProvinceID:   2,
			CountryID:         1,
			Phone:             "(345) 678-9012",
			Website:           "http://monumentalists.com",
			Established:       "2019",
			EmployeeCount:     "2-9",
			MemorialQuota:     2000,
			MemorialQuotaUsed: 1000,
			StorageQuota:      0.5,
			StorageQuotaUsed:  0.25,
		},
		{
			TenantIdentifier:  "momandpop",
			Name:              "Mom and Pop",
			LogoURL:           "resources/brand3.png",
			FSPType:           schema.FSPTypeFuneralHome,
			BusinessTypeID:    3,
			BusinessType:      schema.BusinessTypeSoleProp,
			CRN:               "1234567",
			EIN:               "56-7890123",
			BillingAddress:    "789 Oak Drive, San Francisco",
			Address:           "789 Oak Drive, San Francisco",
			PostalCode:        "94102",
			StateProvinceID:   1,
			CountryID:         1,
			Phone:             "(234) 567-8901",
			Website:           "http://mompop.com",
			Established:       "1935",
			EmployeeCount:     "10-19",
			MemorialQuota:     1000,
			MemorialQuotaUsed: 500,
			StorageQuota:      0.5,
			StorageQuotaUsed:  0.25,
		},
		{
			TenantIdentifier:  "staging",
			Name:              "Staging",
			LogoURL:           "resources/brand_staging.png",
			FSPType:           schema.FSPTypeCemetery,
			BusinessTypeID:    4,
			BusinessType:      schema.BusinessTypeCorporation,
			CRN:               "STG1234567",
			EIN:               "98-7654322",
			BillingAddress:    "101 Staging Lane, Testville",
			Address:           "101 Staging Lane, Testville",
			PostalCode:        "99999",
			StateProvinceID:   4,
			CountryID:         1,
			Phone:             "(123) 555-7890",
			Website:           "http://staging.com",
			Established:       "2023",
			EmployeeCount:     "10-50",
			MemorialQuota:     5000,
			MemorialQuotaUsed: 100,
			StorageQuota:      1,
			StorageQuotaUsed:  0.1,
		},
	}

	for _, fsp := range fsps {
		var existingFSP schema.FSP
		if err := db.Where("name = ?", fsp.Name).First(&existingFSP).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&fsp).Error; err != nil {
					return fmt.Errorf("failed to create FSP %s: %w", fsp.Name, err)
				}
			} else {
				return fmt.Errorf("failed to check FSP %s: %w", fsp.Name, err)
			}
		}
	}

	return nil
}

func (d *Domain) seedUsersAndFSPRoles() error {
	db := d.params.DB.GetDB()

	users := []schema.User{
		{
			FSPID:          1,
			Email:          "richard@reverehere.com",
			FirstName:      "Richard",
			LastName:       "Thompson",
			AvatarURL:      "resources/avatar1.png",
			PasswordHash:   "$2a$10$4wT/6yXlp7GnHcOOHQ5JS.4F6E3/K9PB7tC7SA74F0EVZAMdQguU2", // bcrypt hash
			EmailConfirmed: true,
		},
		{
			FSPID:          1,
			Email:          "michael@reverehere.com",
			FirstName:      "Michael",
			LastName:       "Chen",
			AvatarURL:      "resources/avatar2.png",
			PasswordHash:   "$2a$10$4wT/6yXlp7GnHcOOHQ5JS.4F6E3/K9PB7tC7SA74F0EVZAMdQguU2",
			EmailConfirmed: true,
		},
		{
			FSPID:          1,
			Email:          "joey@reverehere.com",
			FirstName:      "Joseph",
			LastName:       "Scully",
			AvatarURL:      "resources/avatar3.png",
			PasswordHash:   "$2a$10$4wT/6yXlp7GnHcOOHQ5JS.4F6E3/K9PB7tC7SA74F0EVZAMdQguU2",
			EmailConfirmed: true,
		},
		{
			FSPID:          4,
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

		existingUser, err = d.params.Identity.FindUserByEmail(nil, user.FSPID, user.Email)
		if err != nil {
			return fmt.Errorf("seedUsers: %w", err)
		}

		if existingUser == nil {
			// Create user and FSP role
			d.params.Identity.CreateUserAndFSPRole(
				nil,                      // db *gorm.DB
				user.FSPID,               // FSPID uint
				user.FirstName,           // firstName string
				user.LastName,            // lastName string
				user.Email,               // email string
				&user.PasswordHash,       // passwordHash *string
				schema.RoleFSPSuperAdmin, // startingRole schema.FSPRoleConst
			)
			// Update name, avatar, and email confirmed
			db.Model(&schema.User{}).Where("email = ? AND fsp_id = ?", user.Email, user.FSPID).Updates(schema.User{
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
			FSPID:      1,
			Title:      "Richard Thompson's Memorial",
			Identifier: "richard_thompson_1980-01-01",
			FirstName:  "Richard",
			LastName:   "Thompson",
			DOB:        func(t time.Time) *time.Time { return &t }(time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			FSPID:      1,
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
			memorial.FSPID,
			memorial.Identifier,
		)
		if err != nil {
			d.logger.Error("seedMemorials: %w", zap.Error(err))
		}

		if existingMemorial == nil {
			// Create memorial
			_, err = d.params.FSP.CreateMemorial(
				nil,
				memorial.FSPID,      // FSPID 		uint
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
		FSPID        uint
		Email        string
		MemorialID   uint
		Role         schema.MemorialRoleConst
		Relationship schema.RelationshipConst
	}

	// Define the memorial-curator pairs
	memorialCurators := []MemorialCurator{
		{
			FSPID:        1,
			Email:        "richard@reverehere.com",
			MemorialID:   1,
			Role:         schema.RoleMemCurator,
			Relationship: schema.RelationshipSelf,
		},
		{
			FSPID:        1,
			Email:        "michael@reverehere.com",
			MemorialID:   2,
			Role:         schema.RoleMemCurator,
			Relationship: schema.RelationshipSelf,
		},
		{
			FSPID:        1,
			Email:        "michael@reverhere.com",
			MemorialID:   1,
			Role:         schema.RoleMemContributor,
			Relationship: schema.RelationshipFriend,
		},
		{
			FSPID:        1,
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
			user, err := d.params.Identity.FindUserByEmail(tx, curator.FSPID, curator.Email)
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
				curator.FSPID,
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
				curator.FSPID,         // FSPID 		uint
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
			FSPID:           1,
			MemorialID:      2,
			ApplicantID:     3,
			ApplicationType: schema.ApplicationTypeContributor,
			Relationship:    schema.RelationshipBrother,
			Status:          schema.ApplicationStatusPending,
			AppliedOn:       time.Now(),
		},
		{
			FSPID:           1,
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
		if err := db.Where("fsp_id = ? AND memorial_id = ? AND applicant_id = ?", application.FSPID, application.MemorialID, application.ApplicantID).First(&existingApplication).Error; err != nil {
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
			FSPID:          1,
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
			FSPID:          1,
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
			invitation.FSPID,        // FSPID uint
			*invitation.MemorialID,  // memorialID *uint
			invitation.InviteeEmail, // email string
		)
		if err != nil {
			d.logger.Error("seedContributorInvitations: ", zap.Error(err))
		}

		if existingInvitation == nil {
			d.params.Identity.CreateMemorialContributorInvitation(
				nil,
				invitation.FSPID,        // FSPID 			uint
				invitation.InviterID,    // inviterID 		uint
				invitation.InviteeEmail, // inviteeEmail 	string
				invitation.Relationship, // relationship 	schema.RelationshipConst
				invitation.MemorialID,   // memorialID 		*uint
			)
		}

	}
	return nil
}
