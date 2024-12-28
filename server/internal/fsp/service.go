package fsp

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alsey89/people-matter/internal/common/errmgr"
	"github.com/alsey89/people-matter/internal/identity"
	"github.com/alsey89/people-matter/internal/schema"

	"gorm.io/gorm"
)

// Form Data ------------------------------------------------------

// Account ---------------------------------------------------------

func (d *Domain) GetAccountService(TenantID uint, preloadDetails bool) (*schema.Tenant, error) {
	db := d.params.DB.GetDB()

	existingFSP := schema.Tenant{}

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

func (d *Domain) updateAccountService(TenantID uint, updatedFSP schema.Tenant) error {
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

// Team -----------------------------------------------------------

func (d *Domain) getTeamService(TenantID uint) ([]schema.UserFSPRole, error) {
	db := d.params.DB.GetDB()

	existingTeam := []schema.UserFSPRole{}

	err := db.
		Where("fsp_id = ?", TenantID).
		Joins("JOIN fsp_roles ON fsp_roles.id = user_fsp_roles.fsp_role_id").
		Where("fsp_roles.name IN ?", []schema.FSPRoleConst{schema.RoleFSPAdmin, schema.RoleFSPSuperAdmin}).
		Preload("User").
		Preload("FSPRole").
		Find(&existingTeam).
		Error
	if err != nil {
		return nil, fmt.Errorf("getTeamService: %w", err)
	}

	return existingTeam, nil
}

func (d *Domain) postTeamService(TenantIdentifier string, TenantID uint, email string, startingRole schema.FSPRoleConst) error {
	var err error

	// Check if user with email already exists
	existingUser, err := d.params.Identity.FindUserByEmail(nil, TenantID, email)
	if err != nil {
		return fmt.Errorf("postTeamService: %w", err)
	}

	// 1. user exists
	if existingUser != nil {
		// Fetch user's current roles & check if they already have the starting role
		rolesByLevel, err := d.params.Identity.QueryRolesByLevel(TenantID, existingUser.ID)
		if err != nil {
			return fmt.Errorf("postTeamService: %w", err)
		}
		if rolesByLevel != nil && rolesByLevel.Tenant.FSPRole.Name == startingRole {
			return fmt.Errorf("postTeamService: %w", ErrTeamMemberHasRole)
		}

		// If user does not have the starting role, assign it
		_, err = d.params.Identity.AssignOrUpdateFSPRole(nil, TenantID, existingUser.ID, startingRole)
		if err != nil {
			return fmt.Errorf("postTeamService: %w", err)
		}

		urlPath := fmt.Sprintf("/%s", startingRole)

		variables := map[string]interface{}{
			"role": startingRole,
		}

		go d.params.TransMail.SendMail(
			TenantID,  // TenantID 			uint
			email,     // recipientEmail 	string
			6405058,   // templateID 		int
			&urlPath,  // urlPath 			*string
			variables, // variables 		map[string]interface{}
		)

		return nil
	}

	// 2. user does not exist, create the user and send email with link to reset password
	createdUser, err := d.params.Identity.CreateUserAndFSPRole(
		nil,      // db 			*gorm.DB
		TenantID, // TenantID 		uint
		"New",    // firstName 	string
		"User",   // lastName 	string
		email,    // email 		string
		nil,      // passwordHash *string
		startingRole,
	)
	if err != nil {
		return fmt.Errorf("postTeamService: %w", err)
	}

	fullName := fmt.Sprintf("%s %s", createdUser.FirstName, createdUser.LastName)

	variables := map[string]interface{}{
		"fullName": fullName,
		"role":     startingRole,
	}

	resetUrlPath, err := d.params.Identity.GeneratePasswordResetTokenAndPath(
		TenantID, // TenantID 			uint
		email,    // email 			string
	)
	if err != nil {
		return fmt.Errorf("postTeamService: %w", err)
	}
	if resetUrlPath == nil {
		return fmt.Errorf("postTeamService: %w", errmgr.ErrNilCheckFailed)
	}

	go d.params.TransMail.SendMail(
		TenantID,     // TenantID 				uint
		email,        // recipientEmail 	string
		6458393,      // templateID 		int
		resetUrlPath, // urlPath 			*string
		variables,    // variables 			map[string]interface{}
	)

	return nil
}

func (d *Domain) putTeamService(TenantID uint, teamMemberID uint, updatedRole schema.FSPRoleConst) error {
	var err error

	_, err = d.params.Identity.AssignOrUpdateFSPRole(nil, TenantID, teamMemberID, updatedRole)
	if err != nil {
		return fmt.Errorf("putTeamService: %w", err)
	}

	return nil
}

func (d *Domain) deleteTeamService(TenantID uint, teamMemberID uint, notifyUser bool) error {
	var err error

	// Check if the user is a member of the Tenant admin team, throw error if not
	// Check if the user is the only Super Admin, throw error if so
	userFSPRoles, err := d.getTeamService(TenantID)
	if err != nil {
		return fmt.Errorf("deleteTeamService: %w", err)
	}

	var teamMember *schema.User
	var teamMemberRole *schema.FSPRole
	superAdminCount := 0

	for _, userFSPRole := range userFSPRoles {
		if userFSPRole.User.ID == teamMemberID {
			teamMember = userFSPRole.User
			teamMemberRole = userFSPRole.FSPRole
		}
		if userFSPRole.FSPRole.Name == schema.RoleFSPSuperAdmin {
			superAdminCount++
		}
	}

	if teamMember == nil {
		return fmt.Errorf("deleteTeamService: %s", "user is not a member of the Tenant admin team")
	}
	if superAdminCount <= 1 && teamMemberRole.Name == schema.RoleFSPSuperAdmin {
		return fmt.Errorf("deleteTeamService: %w", errmgr.ErrUserIsLastSuperAdmin)
	}

	// Change user role to FSPUser
	_, err = d.params.Identity.AssignOrUpdateFSPRole(nil, TenantID, teamMemberID, schema.RoleFSPUser)
	if err != nil {
		return fmt.Errorf("deleteTeamService: %w", err)
	}

	if notifyUser {
		go d.params.TransMail.SendMail(
			TenantID,         // TenantID 				uint
			teamMember.Email, // recipientEmail 	string
			6458368,          // templateID 		int
			nil,              // urlPath 			*string
			nil,              // variables 			map[string]interface{}
		)
	}

	return nil
}

// Memorial -------------------------------------------------------

func (d *Domain) getAllMemorials(TenantID uint) ([]schema.Memorial, error) {
	db := d.params.DB.GetDB()

	existingMemorials := []schema.Memorial{}

	err := db.
		Where("fsp_id = ?", TenantID).
		Preload("UserMemorialRoles.MemorialRole").
		Preload("UserMemorialRoles.User").
		Find(&existingMemorials).
		Error
	if err != nil {
		return nil, fmt.Errorf("getMemorialService: %w", err)
	}

	return existingMemorials, nil
}

// Creates a new memorial, a new user, and assigns the user the curator role for the memorial.
// Checks if the memorial already exists, if it does, it uses the existing memorial.
// Checks if the user already exists, if it does, it uses the existing user.
// Checks if the user already has a role for the memorial, if it does, it updates the role.
// Intended for use when Tenant admin creates a new memorial.
func (d *Domain) createOrUpdateMemorialWithUserAndCuratorRole(TenantID uint, firstName string, lastName string, DOB *time.Time, DOD *time.Time, emailOfTheCurator string, relationship schema.RelationshipConst) error {
	var err error

	db := d.params.DB.GetDB()

	formattedName := strings.Trim(
		fmt.Sprintf(
			"%s_%s",
			strings.ToLower(firstName),
			strings.ToLower(lastName),
		),
		" ",
	)

	identifier := fmt.Sprintf("%s_%s", formattedName, DOB.Format("2006-01-02"))

	newMemorial := schema.Memorial{
		TenantID:        TenantID,
		Identifier:      identifier, // identifier is a combination of first name, last name, and dob
		IdentifierIsSet: false,
		FirstName:       firstName,
		LastName:        lastName,
		Title:           fmt.Sprintf("%s %s's Memorial", firstName, lastName),
		DOB:             DOB,
		DOD:             DOD,
	}

	// Check if memorial with identifier already exists
	existingMemorial, err := d.FindMemorialByIdentifier(identifier)
	if err != nil {
		return fmt.Errorf("createOrUpdateMemorialCuratorAndUserRole: %w", err)
	}

	// Check if the curator already has an account
	existingUser, err := d.params.Identity.FindUserByEmail(nil, TenantID, emailOfTheCurator)
	if err != nil {
		return fmt.Errorf("createOrUpdateMemorialCuratorAndUserRole: %w", err)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		var memorialID uint

		// if memorial does not exist, create it
		// if memorial exists, use the existing memorial
		if existingMemorial == nil {
			err = tx.Create(&newMemorial).Error
			if err != nil {
				return fmt.Errorf("createOrUpdateMemorialCuratorAndUserRole: %w", err)
			}
			memorialID = newMemorial.ID
		} else {
			memorialID = existingMemorial.ID
		}

		var user *schema.User

		// if user does not exist, create it
		// if user exists, use the existing user
		if existingUser == nil {
			user, err = d.params.Identity.CreateUserAndFSPRole(
				tx,                 // db 			*gorm.DB
				TenantID,           // TenantID 		uint
				"New",              // firstName 	string
				"User",             // lastName 	string
				emailOfTheCurator,  // email 		string
				nil,                // passwordHash  *string
				schema.RoleFSPUser, // startingRole schema.RoleConst
			)
			if err != nil && !errors.Is(err, identity.ErrEmailAlreadyInUse) {
				return fmt.Errorf("createOrUpdateMemorialCuratorAndUserRole: %w", err)
			}
		} else {
			user = existingUser
		}

		// assign or update memorial to the user
		_, err = d.params.Identity.AssignOrUpdateMemorialRole(
			tx,                    // db				*gorm.DB
			TenantID,              // TenantID 			uint
			user.ID,               // userID 	    	uint
			memorialID,            // memorialID 		uint
			schema.RoleMemCurator, // roleIDToAssign	uint
			&relationship,         // relationship		*schema.RelationshipConst
		)
		if err != nil {
			return fmt.Errorf("createOrUpdateMemorialCuratorAndUserRole: %w", err)
		}

		// send email to set password and confirm their account
		setPWUrlPath, err := d.params.Identity.GenerateSetPasswordTokenAndPath(
			TenantID,          // TenantID 			uint
			emailOfTheCurator, // email 			string
		)
		if err != nil {
			return fmt.Errorf("createOrUpdateMemorialCuratorAndUserRole: %w", err)
		}
		if setPWUrlPath == nil {
			return fmt.Errorf("createOrUpdateMemorialCuratorAndUserRole: %w", errmgr.ErrNilCheckFailed)
		}

		variables := map[string]interface{}{
			"memorial_title": newMemorial.Title,
		}
		go d.params.TransMail.SendMail(
			TenantID,          // TenantID 			uint
			emailOfTheCurator, // recipientEmail 	string
			6399852,           // templateID 		int
			setPWUrlPath,      // urlPath 			*string
			variables,         // variables 		map[string]interface{}
		)

		return nil
	})
	if err != nil {
		return fmt.Errorf("createOrUpdateMemorialCuratorAndUserRole: %w", err)
	}

	return nil
}

// Creates a new user and assigns the user the given memorial role for the memorial.
// Checks if the memorial already exists, returns an error if it does not exist.
// Checks if the user already exists, if it does, it uses the existing user.
// Checks if the user already has a role for the memorial, if it does, it updates the role.
// Intended for use when Tenant admin assigns a role to a user for a memorial.
func (d *Domain) createOrUpdateUserWithMemorialRole(TenantID uint, memorialID uint, roleToAssign schema.MemorialRoleConst, userEmail string) error {
	var err error

	db := d.params.DB.GetDB()

	// Check if memorial ID Exists
	existingMemorial, err := d.FindMemorialByID(memorialID)
	if err != nil {
		return fmt.Errorf("createOrUpdateMemorialCuratorAndUserRole: %w", err)
	}
	if existingMemorial == nil {
		return fmt.Errorf("createOrUpdateMemorialCuratorAndUserRole: Memorial does not exist")
	}

	// Check if the curator already has an account
	existingUser, err := d.params.Identity.FindUserByEmail(nil, TenantID, userEmail)
	if err != nil {
		return fmt.Errorf("createOrUpdateMemorialCuratorAndUserRole: %w", err)
	}

	err = db.Transaction(func(tx *gorm.DB) error {

		var user *schema.User

		// if user does not exist, create it
		// if user exists, use the existing user
		if existingUser == nil {
			user, err = d.params.Identity.CreateUserAndFSPRole(
				tx,                 // db 			*gorm.DB
				TenantID,           // TenantID 		uint
				"New",              // firstName 	string
				"User",             // lastName 	string
				userEmail,          // email 		string
				nil,                // passwordHash *string
				schema.RoleFSPUser, // startingRole schema.RoleConst
			)
			if err != nil && !errors.Is(err, identity.ErrEmailAlreadyInUse) {
				return fmt.Errorf("createOrUpdateMemorialCuratorAndUserRole: %w", err)
			}
		} else {
			user = existingUser
		}

		// assign or update memorial to the user
		_, err = d.params.Identity.AssignOrUpdateMemorialRole(
			tx,           // db				*gorm.DB
			TenantID,     // TenantID 			uint
			user.ID,      // userID 	    uint
			memorialID,   // memorialID 	uint
			roleToAssign, // roleToAssign	schema.MemorialRoleConst
			nil,          // relationship	*schema.RelationshipConst
		)
		if err != nil {
			return fmt.Errorf("createOrUpdateMemorialCuratorAndUserRole: %w", err)
		}
		return nil

	})
	if err != nil {
		return fmt.Errorf("createOrUpdateMemorialCuratorAndUserRole: %w", err)
	}

	return nil
}

// Creates a new memorial. Checks if a memorial exists by identifier and creates it if not.
func (d *Domain) CreateMemorial(db *gorm.DB, TenantID uint, firstName string, lastName string, identifier string, DOB *time.Time, DOD *time.Time) (*schema.Memorial, error) {
	// if db is not set, use the default db
	// this allows the function to be used in a transaction
	if db == nil {
		db = d.params.DB.GetDB()
	}

	//if identifier is not set, create the default identifier
	if identifier == "" {
		formattedName := strings.Trim(
			fmt.Sprintf(
				"%s_%s",
				strings.ToLower(firstName),
				strings.ToLower(lastName),
			),
			" ",
		)
		identifier = fmt.Sprintf("%s_%s", formattedName, DOB.Format("2006-01-02"))
	}

	// Check if memorial with identifier already exists
	existingMemorial, err := d.FindMemorialByIdentifier(identifier)
	if err != nil {
		return nil, fmt.Errorf("createMemorial: %w", err)
	}

	// if memorial exists, return the existing memorial
	if existingMemorial != nil {
		return existingMemorial, nil
	}

	// if memorial does not exist, create it
	newMemorial := schema.Memorial{
		TenantID:        TenantID,
		Identifier:      identifier,
		IdentifierIsSet: false,
		FirstName:       firstName,
		LastName:        lastName,
		Title:           fmt.Sprintf("%s %s's Memorial", firstName, lastName),
		DOB:             DOB,
		DOD:             DOD,
	}

	err = db.Create(&newMemorial).Error
	if err != nil {
		return nil, fmt.Errorf("createMemorial: %w", err)
	}

	return &newMemorial, nil
}

func (d *Domain) updateMemorial(TenantID uint, memorialID uint, updatedMemorial schema.Memorial) error {
	var err error

	db := d.params.DB.GetDB()

	err = db.
		Where("id = ?", memorialID).
		Where("fsp_id = ?", TenantID).
		Updates(&updatedMemorial).
		Error
	if err != nil {
		return fmt.Errorf("updateMemorial: %w", err)
	}

	return nil
}

func (d *Domain) deleteMemorial(TenantID uint, memorialID uint) error {
	var err error

	db := d.params.DB.GetDB()

	err = db.
		Where("id = ?", memorialID).
		Where("fsp_id = ?", TenantID).
		Delete(&schema.Memorial{}).
		Error
	if err != nil {
		return fmt.Errorf("deleteMemorial: %w", err)
	}

	return nil
}

// Helper ---------------------------------------------------------

// Finds a Tenant by its identifier, returns the Tenant if it exists
// **Returns nil without an error** if the Tenant does not exist
func (d *Domain) FindFSPByID(TenantID uint) (*schema.Tenant, error) {
	db := d.params.DB.GetDB()

	existingFSP := schema.Tenant{}

	err := db.
		Where("id = ?", TenantID).
		First(&existingFSP).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("FindFSPByID: %w", err)
	}

	return &existingFSP, nil
}

// Finds a memorial by its identifier, returns the memorial if it exists
// **Returns nil without an error** if the memorial does not exist
func (d *Domain) FindMemorialByIdentifier(identifier string) (*schema.Memorial, error) {
	db := d.params.DB.GetDB()

	existingMemorial := schema.Memorial{}

	err := db.
		Where("identifier = ?", identifier).
		First(&existingMemorial).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("FindMemorialByIdentifier: %w", err)
	}

	return &existingMemorial, nil
}

// Finds a memorial by its ID, returns the memorial if it exists
// **Returns nil without an error** if the memorial does not exist
func (d *Domain) FindMemorialByID(memorialID uint) (*schema.Memorial, error) {
	db := d.params.DB.GetDB()

	existingMemorial := schema.Memorial{}

	err := db.
		Where("id = ?", memorialID).
		First(&existingMemorial).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("FindMemorialByID: %w", err)
	}

	return &existingMemorial, nil
}
