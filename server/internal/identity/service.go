package identity

import (
	"errors"
	"fmt"
	"time"

	"github.com/alsey89/people-matter/internal/common/errmgr"
	"github.com/alsey89/people-matter/internal/schema"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ! Applications and Invitations ------------------------------------

func (d *Domain) acceptInvitationSignupService(FSPID uint, memorialID uint, token string, email string, firstName string, lastName string, password string) error {
	var err error

	// Check if the token is valid
	existingInvitation, err := d.FindInvitationByToken(
		FSPID,      // FSPID 		uint
		memorialID, // memorialID 	uint
		token,      // token 		string
	)
	if err != nil {
		return fmt.Errorf("acceptInvitationSignupService: %w", err)
	}
	if existingInvitation == nil {
		return fmt.Errorf("acceptInvitationSignupService: %w", errmgr.ErrInvitationNotFound)
	}
	if existingInvitation.Status != schema.InvitationStatusPending {
		return fmt.Errorf("acceptInvitationSignupService: %w", errmgr.ErrInvitationResponded)
	}
	if existingInvitation.ExpiresOn.Before(time.Now()) {
		return fmt.Errorf("acceptInvitationSignupService: %w", errmgr.ErrInvitationExpired)
	}

	// Check if the user already exists
	existingUser, err := d.FindUserByEmail(nil, FSPID, email)
	if err != nil {
		return fmt.Errorf("acceptInvitationSignupService: %w", err)
	}
	if existingUser != nil {
		return fmt.Errorf("acceptInvitationSignupService: %w", errmgr.ErrEmailInUse)
	}

	// Create the new user
	createdUser, err := d.signupService(
		FSPID,      //FSPID 				uint
		&firstName, //firstName 			*string
		&lastName,  //lastName 				*string
		email,      //email 				string
		password,   //password 				string
	)
	if err != nil {
		return fmt.Errorf("acceptInvitationSignupService: %w", err)
	}

	err = d.acceptInvitationService(
		FSPID,          // FSPID 		uint
		memorialID,     // memorialID 	uint
		createdUser.ID, // userID 		uint
		token,          // token 		string
	)
	if err != nil {
		return fmt.Errorf("acceptInvitationSignupService: %w", err)
	}

	// assign contributor role to the user
	_, err = d.AssignOrUpdateMemorialRole(
		nil,                              // db				*gorm.DB
		FSPID,                            // FSPID			uint
		createdUser.ID,                   // userID			uint
		memorialID,                       // memorialID		uint
		schema.RoleMemContributor,        // roleToAssign	schema.RoleConst
		&existingInvitation.Relationship, // relationship	*schema.RelationshipConst
	)
	if err != nil {
		return fmt.Errorf("acceptInvitationSignupService: %w", err)
	}

	//todo: send email to the curator

	return nil
}

func (d *Domain) acceptInvitationService(FSPID uint, memorialID uint, userID uint, token string) error {
	var err error
	db := d.params.DB.GetDB()

	// Check if the token is valid
	existingInvitation, err := d.FindInvitationByToken(
		FSPID,      // FSPID 		uint
		memorialID, // memorialID 	uint
		token,      // token 		string
	)
	if err != nil {
		return fmt.Errorf("acceptInvitationService: %w", err)
	}
	if existingInvitation == nil {
		return fmt.Errorf("acceptInvitationService: %w", errmgr.ErrInvitationNotFound)
	}
	if existingInvitation.Status != schema.InvitationStatusPending {
		return fmt.Errorf("acceptInvitationService: %w", errmgr.ErrInvitationResponded)
	}
	if existingInvitation.ExpiresOn.Before(time.Now()) {
		return fmt.Errorf("acceptInvitationService: %w", errmgr.ErrInvitationExpired)
	}

	// Check if the user exists
	existingUser, err := d.FindUserByID(nil, FSPID, userID)
	if err != nil {
		return fmt.Errorf("acceptInvitationService: %w", err)
	}
	if existingUser == nil {
		return fmt.Errorf("acceptInvitationService: %w", errmgr.ErrUserNotFound)
	}
	if existingUser.Email != existingInvitation.InviteeEmail {
		return fmt.Errorf("acceptInvitationService: %w", errmgr.ErrInvitationNotForUser)
	}

	toUpdate := schema.Invitation{
		Status: schema.InvitationStatusAccepted,
		IsUsed: true,
	}

	err = db.Model(&schema.Invitation{}).
		Where("fsp_id = ?", FSPID).
		Where("memorial_id = ?", memorialID).
		Where("token = ?", token).
		Updates(&toUpdate).Error
	if err != nil {
		return fmt.Errorf("acceptInvitationService: %w", err)
	}

	// assign contributor role to the user
	_, err = d.AssignOrUpdateMemorialRole(
		nil,                              // db				*gorm.DB
		FSPID,                            // FSPID			uint
		existingUser.ID,                  // userID			uint
		memorialID,                       // memorialID		uint
		schema.RoleMemContributor,        // roleToAssign	schema.RoleConst
		&existingInvitation.Relationship, // relationship	*schema.RelationshipConst
	)
	if err != nil {
		return fmt.Errorf("acceptInvitationSignupService: %w", err)
	}

	// todo: send email to the curator

	return nil
}

func (d *Domain) applicationSignupService(FSPID uint, memorialID uint, firstName string, lastName string, relationship schema.RelationshipConst, email string, password string) (*schema.Application, error) {
	var user *schema.User
	var err error

	// Check if the memorial exists
	existingMemorial, err := d.FindMemorialByID(
		FSPID,      // FSPID 		uint
		memorialID, // memorialID 	uint
	)
	if err != nil {
		return nil, fmt.Errorf("applicationService: %w", err)
	}
	if existingMemorial == nil {
		return nil, fmt.Errorf("applicationService: %w", ErrMemorialNotFound)
	}

	// Create the new user
	user, err = d.signupService(
		FSPID,      //FSPID 				uint
		&firstName, //firstName 			*string
		&lastName,  //lastName 			*string
		email,      //email 				string
		password,   //password 			string
	)
	if err != nil {
		return nil, fmt.Errorf("applicationSignupService: %w", err)
	}

	// Create a new application for the contributor role
	newApplication, err := d.CreateMemorialContributorApplication(
		nil,                               // db 				*gorm.DB
		FSPID,                             // FSPID 			uint
		user.ID,                           // applicantID 		uint
		memorialID,                        // memorialID 		uint
		relationship,                      // relationship 		schema.RelationshipConst
		schema.ApplicationTypeContributor, // applicationType 	schema.ApplicationTypeConst
	)
	if err != nil {
		return nil, fmt.Errorf("applicationSignupService: %w", err)
	}

	// TODO: Send application received email to the new user
	// go d.params.TransMail.SendMemorialApplicationReceivedEmail("Application Received", user.Email)
	// TODO: Notify curator for approval
	// go d.params.TransMail.SendApplicationReceivedEmail("Application Pending", user.Email)

	return newApplication, nil
}

func (d *Domain) applicationService(FSPID uint, memorialID uint, applicantID uint, relationship schema.RelationshipConst) (*schema.Application, error) {
	db := d.params.DB.GetDB()

	// Check if the memorial exists
	existingMemorial, err := d.FindMemorialByID(
		FSPID,      // FSPID 		uint
		memorialID, // memorialID 	uint
	)
	if err != nil {
		return nil, fmt.Errorf("applicationService: %w", err)
	}
	if existingMemorial == nil {
		return nil, fmt.Errorf("applicationService: %w", ErrMemorialNotFound)
	}

	// Check if the user already has a role in the memorial
	existingUserMemorialRole, err := d.FindUserMemorialRoleByUserID(
		nil,         // db 				*gorm.DB
		FSPID,       // FSPID 			uint
		applicantID, // userID 			uint
		memorialID,  // memorialID 		uint
	)
	if err != nil {
		return nil, fmt.Errorf("applicationService: %w", err)
	}
	if existingUserMemorialRole != nil {
		return nil, fmt.Errorf("applicationService: %w", ErrUserHasMemorialRole)
	}

	// Check if the user already has an application for the memorial
	existingApplication, err := d.FindContributorApplicationByApplicantID(
		FSPID,       // FSPID 			uint
		memorialID,  // memorialID 		uint
		applicantID, // email 			string
	)
	if err != nil {
		return nil, fmt.Errorf("applicationService: %w", err)
	}
	if existingApplication != nil {
		return nil, fmt.Errorf("applicationService: %w", ErrUserHasApplication)
	}

	// Create a new application for the contributor role
	newApplication, err := d.CreateMemorialContributorApplication(
		db,                                // db 				*gorm.DB
		FSPID,                             // FSPID 			uint
		applicantID,                       // applicantID 		uint
		memorialID,                        // memorialID 		uint
		relationship,                      // relationship 		schema.RelationshipConst
		schema.ApplicationTypeContributor, // applicationType 	schema.ApplicationTypeConst
	)
	if err != nil {
		return nil, fmt.Errorf("applicationService: %w", err)
	}

	//todo: send email to applicant that the application has been received
	//todo: send email to the curator

	return newApplication, nil
}

func (d *Domain) CreateFSPInvitation(FSPID uint, inviterID uint, inviteeEmail string, FSPRole schema.FSPRoleConst) (*schema.Invitation, error) {
	db := d.params.DB.GetDB()
	token := uuid.New().String()
	expirationDate := time.Now().Add(time.Hour * 24 * 7) // 1 week

	newInvitation := schema.Invitation{
		InvitationType: schema.InvitationTypeFSP,
		FSPID:          FSPID,
		InviterID:      inviterID,
		InviteeEmail:   inviteeEmail,
		Status:         schema.InvitationStatusPending,
		InvitedOn:      time.Now(),
		ExpiresOn:      expirationDate,
		Token:          token,
		FSPRole:        FSPRole,
	}

	err := db.Create(&newInvitation).Error
	if err != nil {
		return nil, fmt.Errorf("createFSPInvitation: %w", err)
	}

	return &newInvitation, nil
}

func (d *Domain) CreateMemorialContributorInvitation(db *gorm.DB, FSPID uint, inviterID uint, inviteeEmail string, relationship schema.RelationshipConst, memorialID *uint) (*schema.Invitation, error) {
	// If db is nil, get the regular DB instance
	// this allows the function to be used in transactions
	if db == nil {
		db = d.params.DB.GetDB()
	}

	token := uuid.New().String()
	expirationDate := time.Now().Add(time.Hour * 24 * 7) // 1 week

	newInvitation := schema.Invitation{
		InvitationType: schema.InvitationTypeMemorial,
		FSPID:          FSPID,
		InviterID:      inviterID,
		InviteeEmail:   inviteeEmail,
		Relationship:   relationship,
		Status:         schema.InvitationStatusPending,
		InvitedOn:      time.Now(),
		ExpiresOn:      expirationDate,
		Token:          token,
		MemorialID:     memorialID,
		MemorialRole:   schema.RoleMemContributor,
	}

	err := db.Create(&newInvitation).Error
	if err != nil {
		return nil, fmt.Errorf("CreateMemorialContributorInvitation: %w", err)
	}

	return &newInvitation, nil
}

func (d *Domain) RefreshContributorInvitationAndToken(db *gorm.DB, FSPID uint, memorialID uint, invitationID uint) (*schema.Invitation, error) {
	// If db is nil, get the regular DB instance
	if db == nil {
		db = d.params.DB.GetDB()
	}

	newToken := uuid.New().String()
	expirationDate := time.Now().Add(time.Hour * 24 * 7) // 1 week

	var updatedInvitation schema.Invitation

	// Update the invitation and return the updated row
	err := db.Model(&schema.Invitation{}).
		Where("fsp_id = ?", FSPID).
		Where("memorial_id = ?", memorialID).
		Where("id = ?", invitationID).
		Clauses(clause.Returning{}).
		Updates(schema.Invitation{
			Token:     newToken,
			ExpiresOn: expirationDate,
		}).
		Scan(&updatedInvitation).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("RefreshContributorInvitationAndToken: %w", errmgr.ErrInvitationNotFound)
		}
		return nil, fmt.Errorf("RefreshContributorInvitationAndToken: %w", err)
	}

	return &updatedInvitation, nil
}

func (d *Domain) CreateMemorialContributorApplication(db *gorm.DB, FSPID uint, applicantID uint, memorialID uint, relationship schema.RelationshipConst, applicationType schema.ApplicationTypeConst) (*schema.Application, error) {
	// If db is nil, get the regular DB instance
	// this allows the function to be used in transactions
	if db == nil {
		db = d.params.DB.GetDB()
	}

	var err error

	// Create the new application
	newApplication := schema.Application{
		FSPID:           FSPID,                           // FSPID 				uint
		ApplicantID:     applicantID,                     // ApplicantID 		uint
		MemorialID:      memorialID,                      // MemorialID 		uint
		Relationship:    relationship,                    // Relationship 		schema.RelationshipConst
		ApplicationType: applicationType,                 // ApplicationType 	schema.ApplicationTypeConst
		Status:          schema.ApplicationStatusPending, // Status 			schema.ApplicationStatusConst
		AppliedOn:       time.Now(),                      // AppliedOn 		time.Time
	}

	// Save the new application
	err = db.Create(&newApplication).Error
	if err != nil {
		return nil, fmt.Errorf("CreateMemorialContributorApplication: %w", err)
	}

	return &newApplication, nil
}

// ! Auth Service ----------------------------------------------------

func (d *Domain) signinService(FSPID uint, email string, password string) (*schema.User, *rolesByLevel, error) {
	db := d.params.DB.GetDB()

	var existingUser schema.User

	err := db.Model(&schema.User{}).
		Where("fsp_id = ?", FSPID).
		Where("email = ?", email).
		First(&existingUser).
		Error
	if err != nil {
		return nil, nil, fmt.Errorf("signinService: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(password))
	if err != nil {
		return nil, nil, fmt.Errorf("signinService: %w", ErrInvalidCredentials)
	}

	if !existingUser.EmailConfirmed {
		return nil, nil, fmt.Errorf("signinService: %w", ErrEmailNotConfirmed)
	}

	rolesByLevel, err := d.QueryRolesByLevel(FSPID, existingUser.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("signinService: %w", err)
	}

	return &existingUser, rolesByLevel, nil
}

func (d *Domain) serviceSwitchActiveMemorial(FSPID uint, userID uint, targetMemorialID uint) (*rolesByLevel, error) {
	db := d.params.DB.GetDB()

	// Check if the user exists
	existingUser, err := d.FindUserByID(db, FSPID, userID)
	if err != nil {
		return nil, fmt.Errorf("switchActiveMemorialService: %w", err)
	}
	if existingUser == nil {
		return nil, fmt.Errorf("switchActiveMemorialService: %w", errmgr.ErrUserNotFound)
	}

	rolesByLevel, err := d.QueryRolesByLevel(FSPID, userID)
	if err != nil {
		return nil, fmt.Errorf("switchActiveMemorialService: %w", err)
	}

	userHasRoleInTargetMemorial := func(roles []schema.UserMemorialRole, targetMemorialID uint) bool {
		for _, role := range roles {
			if role.MemorialID == targetMemorialID {
				return true
			}
		}
		return false
	}
	if !userHasRoleInTargetMemorial(rolesByLevel.Memorial, targetMemorialID) {
		return nil, fmt.Errorf("switchActiveMemorialService: %w", errmgr.ErrMemorialRoleNotFound)
	}

	return rolesByLevel, nil
}

func (d *Domain) signupService(FSPID uint, firstName *string, lastName *string, email string, password string) (*schema.User, error) {
	// Set default first and last name if not provided
	if firstName == nil {
		defaultFirstName := "New"
		firstName = &defaultFirstName
	}
	if lastName == nil {
		defaultLastName := "User"
		lastName = &defaultLastName
	}

	// Check if the user already exists
	existingUser, err := d.FindUserByEmail(nil, FSPID, email)
	if err != nil {
		return nil, fmt.Errorf("signupService: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("signupService: %w", ErrEmailAlreadyInUse)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("signupService: %w", err)
	}

	hashedPasswordString := string(hashedPassword)

	// Create a new user with a default FSP user role
	newUser, err := d.CreateUserAndFSPRole(
		nil,                   // db
		FSPID,                 // FSPID
		*firstName,            // firstName
		*lastName,             // lastName
		email,                 // email
		&hashedPasswordString, // passwordHash as string
		schema.RoleFSPUser,    // startingRole
	)
	if err != nil {
		return nil, fmt.Errorf("signupService: %w", err)
	}

	// Generate email confirmation token
	additionalClaims := map[string]interface{}{
		"email": newUser.Email,
	}
	confirmationToken, err := d.params.TokenManager.GenerateToken(d.config.JWTEmailConfirmationScope, additionalClaims)
	if err != nil {
		return nil, fmt.Errorf("signupService: %w", err)
	}
	if confirmationToken == nil {
		return nil, fmt.Errorf("signupService: %s", "confirmationToken is nil")
	}

	confirmationPath := fmt.Sprintf("/auth/signup/confirm?token=%s", *confirmationToken)

	variables := map[string]interface{}{
		"role": "User",
	}

	d.params.TransMail.SendMail(
		FSPID,             // FSPID 			uint
		email,             // recipientEmail 	string
		6405058,           // templateID 		int
		&confirmationPath, // urlPath 			*string
		variables,         // variables 		map[string]interface{}
	)

	return newUser, nil
}

func (d *Domain) confirmEmailService(FSPID uint, email string) error {
	db := d.params.DB.GetDB()

	var existingUser schema.User

	err := db.Model(&schema.User{}).
		Where("fsp_id = ?", FSPID).
		Where("email = ?", email).
		First(&existingUser).Error
	if err != nil {
		return fmt.Errorf("confirmEmailService: %w", err)
	}

	if existingUser.EmailConfirmed {
		return fmt.Errorf("confirmEmailService: %w", ErrEmailAlreadyConfirmed)
	}

	err = db.Model(&schema.User{}).Where("email = ?", email).Update("email_confirmed", true).Error
	if err != nil {
		return fmt.Errorf("confirmEmailService: %w", err)
	}

	return nil
}

// Request a password reset email.
// Checks if the user exists. Send email with reset link if yes, otherwise return an error.
func (d *Domain) requestResetPasswordService(FSPID uint, email string) error {

	// Check if the user exists
	existingUser, err := d.FindUserByEmail(nil, FSPID, email)
	if err != nil {
		return fmt.Errorf("RequestResetPasswordService: %w", err)
	}
	if existingUser == nil {
		return fmt.Errorf("RequestResetPasswordService: %w", errmgr.ErrUserNotFound)
	}

	// Check if the user exists
	_, err = d.FindUserByEmail(nil, FSPID, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("GeneratePasswordResetLink: %w", errmgr.ErrUserNotFound)
		}
		return fmt.Errorf("GeneratePasswordResetLink: %w", err)
	}

	// Generate a password reset token
	additionalClaims := map[string]interface{}{
		"email": email,
	}
	resetToken, err := d.params.TokenManager.GenerateToken(d.config.JWTPasswordResetScope, additionalClaims)
	if err != nil {
		return fmt.Errorf("GeneratePasswordResetLink: %w", err)
	}
	if resetToken == nil {
		return fmt.Errorf("GeneratePasswordResetLink: %w", errmgr.ErrNilCheckFailed)
	}

	resetUrlPath := fmt.Sprintf("/auth/password/reset/confirm?token=%s", *resetToken)

	// Send the reset email
	go d.params.TransMail.SendMail(
		FSPID,         // FSPID 			uint
		email,         // recipientEmail 	string
		6458354,       // templateID 		int
		&resetUrlPath, //urlPath 			*string
		nil,           //variables 			map[string]interface{}
	)

	return nil
}

// Wrapper function around updatePasswordServiceWithConfirmation that sets alsoConfirmAccount to false.
func (d *Domain) updatePasswordServiceWithoutConfirmation(FSPID uint, email string, password string) error {
	return d.updatePasswordServiceWithConfirmation(FSPID, email, password, false)
}

// Base function also confirms the email if the alsoConfirmAccount flag is set to true.
// This is useful when user accounts are created by an admin and the user is required to set their password. Which in itself can be considered as confirming the email.
func (d *Domain) updatePasswordServiceWithConfirmation(FSPID uint, email string, password string, alsoConfirmAccount bool) error {
	db := d.params.DB.GetDB()

	// Fetch existing user
	var existingUser schema.User
	err := db.Model(&schema.User{}).
		Where("fsp_id = ?", FSPID).
		Where("email = ?", email).
		First(&existingUser).
		Error
	if err != nil {
		return fmt.Errorf("updatePasswordService: %w", err)
	}

	// Confirm email if not already confirmed and the flag is true
	if !existingUser.EmailConfirmed && alsoConfirmAccount {
		err := d.confirmEmailService(FSPID, email)
		if err != nil {
			return fmt.Errorf("updatePasswordService: %w", err)
		}
	}

	// Check if the new password matches the old one
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(password))
	if err == nil {
		return fmt.Errorf("updatePasswordService: %w", ErrNewPasswordIsOldPassword)
	}

	// Hash and update the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("updatePasswordService: password hashing failed: %w", err)
	}

	err = db.Model(&schema.User{}).
		Where("email = ?", email).
		Update("password_hash", string(hashedPassword)).
		Error
	if err != nil {
		return fmt.Errorf("updatePasswordService: password update failed: %w", err)
	}

	return nil
}

// Auth Helpers -----------------------------------------------------

// Generates a password reset token and path for the given email.
// !WARNING: It does *NOT* check if the user exists.
func (d *Domain) GeneratePasswordResetTokenAndPath(FSPID uint, email string) (*string, error) {

	// Generate a password reset token
	additionalClaims := map[string]interface{}{
		"email": email,
	}
	resetToken, err := d.params.TokenManager.GenerateToken(d.config.JWTPasswordResetScope, additionalClaims)
	if err != nil {
		return nil, fmt.Errorf("GeneratePasswordResetLink: %w", err)
	}
	if resetToken == nil {
		return nil, fmt.Errorf("GeneratePasswordResetLink: %s", "resetToken is nil")
	}

	// construct the reset URL path
	resetUrlPath := fmt.Sprintf("/auth/password/reset/confirm?token=%s", *resetToken)

	return &resetUrlPath, nil
}

// Generates a set password token and path for the given email.
// !WARNING: It does *NOT* check if the user exists.
func (d *Domain) GenerateSetPasswordTokenAndPath(FSPID uint, email string) (*string, error) {

	// Generate a set password token
	additionalClaims := map[string]interface{}{
		"email": email,
	}
	setPasswordToken, err := d.params.TokenManager.GenerateToken(d.config.JWTPasswordResetScope, additionalClaims)
	if err != nil {
		return nil, fmt.Errorf("GenerateSetPasswordTokenAndPath: %w", err)
	}
	if setPasswordToken == nil {
		return nil, fmt.Errorf("GenerateSetPasswordTokenAndPath: %s", "setPasswordToken is nil")
	}

	// construct the set password URL path
	setPasswordUrlPath := fmt.Sprintf("/auth/password/set/confirm?token=%s", *setPasswordToken)

	return &setPasswordUrlPath, nil
}

//! User Service ---------------------------------------------------

// Creates a new user and assigns the user a role at the FSP level.
// It checks if the user with the email already exists.
// If so, it returns an error.
// If not, it creates a new user and assigns the user the given role.
// If the passwordHash is nil, the user will be created without a password.
// Optional db parameter allows passing a transaction instance.
func (d *Domain) CreateUserAndFSPRole(db *gorm.DB, FSPID uint, firstName string, lastName string, email string, passwordHash *string, startingRole schema.FSPRoleConst) (*schema.User, error) {

	var err error

	// If db is nil, get the regular DB instance
	// this allows the function to be used in transactions
	if db == nil {
		db = d.params.DB.GetDB()
	}

	var newUser *schema.User

	err = db.Transaction(func(tx *gorm.DB) error {
		existingUser, err := d.FindUserByEmail(tx, FSPID, email)
		if err != nil {
			return fmt.Errorf("CreateUserAndFSPRole: %w", err)
		}
		if existingUser != nil {
			return fmt.Errorf("CreateUserAndFSPRole: %w", ErrEmailAlreadyInUse)
		}

		newUser = &schema.User{
			FSPID:     FSPID,
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
		}

		// We create the user without a password if passwordHash is nil
		// This is useful for creating users with social login
		// And in the case of admin-created users, where the password will be set later
		if passwordHash != nil {
			newUser.PasswordHash = *passwordHash
		}

		err = tx.Create(newUser).Error
		if err != nil {
			return fmt.Errorf("CreateUserAndFSPRole: %w", err)
		}

		_, err = d.AssignOrUpdateFSPRole(
			tx,           //db 				*gorm.DB
			FSPID,        //FSPID 			uint
			newUser.ID,   //userID 			uint
			startingRole) //roleToAssign 	schema.RoleConst
		if err != nil {
			return fmt.Errorf("CreateUserAndFSPRole: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

//! Role Service ---------------------------------------------------

// AssignOrUpdateFSPRole assigns or updates a user's role at the FSP level.
// It checks if the user already has a UserFSPRole.
// If no record is found, it creates a new UserFSPRole with the given role.
// If a record is found, it updates the existing UserRole with the given role.
func (d *Domain) AssignOrUpdateFSPRole(db *gorm.DB, FSPID uint, userID uint, roleToAssign schema.FSPRoleConst) (*schema.UserFSPRole, error) {
	var err error

	// If db is nil, get the regular DB instance
	// This allows the function to be used in transactions
	if db == nil {
		db = d.params.DB.GetDB()
	}

	// Ensure user exists
	existingUser, err := d.FindUserByID(db, FSPID, userID)
	if err != nil {
		return nil, fmt.Errorf("AssignOrUpdateFSPRole: %w", err)
	}
	if existingUser == nil {
		return nil, fmt.Errorf("AssignOrUpdateFSPRole: %w", errmgr.ErrUserNotFound)
	}

	// Resolve the role ID from the role name
	roleIDToAssign, err := d.ResolveFSPRoleID(roleToAssign)
	if err != nil {
		return nil, fmt.Errorf("AssignOrUpdateFSPRole: %w", err)
	}
	if roleIDToAssign == nil {
		return nil, fmt.Errorf("AssignOrUpdateFSPRole: %s", "roleIDToAssign is nil")
	}

	// Check if user already has an FSP-level role
	existingUserFSPRole, err := d.FindUserFSPRole(
		db,     //db 		*gorm.DB
		FSPID,  //FSPID 	uint
		userID, //userID 	uint
	)
	if err != nil {
		return nil, fmt.Errorf("AssignOrUpdateFSPRole: %w", err)
	}

	// If No FSP-level UserRole exists, create a new one with the given role
	if existingUserFSPRole == nil {
		newUserRole := schema.UserFSPRole{
			FSPID:     FSPID,
			UserID:    userID,
			FSPRoleID: *roleIDToAssign,
		}
		err = db.Create(&newUserRole).Error
		if err != nil {
			return nil, fmt.Errorf("AssignOrUpdateFSPRole: %w", err)
		}
		return &newUserRole, nil
	}

	// If FSP-level role exists, update the role
	existingUserFSPRole.FSPRoleID = *roleIDToAssign
	err = db.
		Model(&schema.UserFSPRole{}).
		Where("id = ?", existingUserFSPRole.ID).
		Updates(existingUserFSPRole).
		Error
	if err != nil {
		return nil, fmt.Errorf("AssignOrUpdateFSPRole: %w", err)
	}

	return existingUserFSPRole, nil
}

// AssignOrUpdateMemorialRole assigns or updates a user's role for the given Memorial.
// It checks if the user already has a UserMemorialRole for the given Memorial.
// If no record is found, it creates a new UserMemorialRole with the given role.
// If a record is found, it updates the existing UserMemorialRole with the given role.
func (d *Domain) AssignOrUpdateMemorialRole(db *gorm.DB, FSPID uint, userID uint, memorialID uint, roleToAssign schema.MemorialRoleConst, relationship *schema.RelationshipConst) (*schema.UserMemorialRole, error) {
	// If db is nil, get the regular DB instance
	// this allows the function to be used in transactions
	if db == nil {
		db = d.params.DB.GetDB()
	}

	// Ensure user exists
	existingUser, err := d.FindUserByID(db, FSPID, userID)
	if err != nil {
		return nil, fmt.Errorf("AssignOrUpdateMemorialRole: %w", err)
	}
	if existingUser == nil {
		return nil, fmt.Errorf("AssignOrUpdateMemorialRole: %w", errmgr.ErrUserNotFound)
	}

	// Resolve the role ID from the role name
	roleIDToAssign, err := d.ResolveMemorialRoleID(roleToAssign)
	if err != nil {
		return nil, fmt.Errorf("AssignOrUpdateMemorialRole: %w", err)
	}
	if roleIDToAssign == nil {
		return nil, fmt.Errorf("AssignOrUpdateMemorialRole: %s", "roleIDToAssign is nil")
	}

	// Check if user already has a role at the given Memorial
	existingMemorialRole, err := d.FindUserMemorialRoleByUserID(
		db,         //db 			*gorm.DB
		FSPID,      //FSPID 		uint
		userID,     //userID 		uint
		memorialID, //memorialID 	uint
	)
	if err != nil {
		return nil, fmt.Errorf("AssignOrUpdateMemorialRole: %w", err)
	}

	// If user has no role at the given Memorial, create a new one with the given role
	// return the new userRole
	if existingMemorialRole == nil {
		newUserRole := schema.UserMemorialRole{
			FSPID:          FSPID,
			UserID:         userID,
			MemorialID:     memorialID,
			MemorialRoleID: *roleIDToAssign,
		}
		if relationship != nil {
			newUserRole.Relationship = *relationship
		}

		err = db.Create(&newUserRole).Error
		if err != nil {
			return nil, fmt.Errorf("AssignOrUpdateMemorialRole: %w", err)
		}

		return &newUserRole, nil
	}

	// If user has a role at the given Memorial, update the role to the given role
	updatedMemorialRole := *existingMemorialRole
	updatedMemorialRole.MemorialRoleID = *roleIDToAssign
	if relationship != nil {
		updatedMemorialRole.Relationship = *relationship
	}

	err = db.
		Model(&schema.UserMemorialRole{}).
		Where("id = ?", existingMemorialRole.ID).
		Updates(existingMemorialRole).
		Error
	if err != nil {
		return nil, fmt.Errorf("AssignOrUpdateMemorialRole: %w", err)
	}

	return &updatedMemorialRole, nil
}

// Removes a user's role at the given Memorial.
// It checks if the user exists and if they already have a UserMemorialRole for the given Memorial.
// If not, it returns an error.
func (d *Domain) RemoveUserMemorialRoleByUserID(db *gorm.DB, FSPID uint, userID uint, memorialID uint) error {
	// If db is nil, get the regular DB instance
	// this allows the function to be used in transactions
	if db == nil {
		db = d.params.DB.GetDB()
	}

	// Ensure user exists
	existingUser, err := d.FindUserByID(db, FSPID, userID)
	if err != nil {
		return fmt.Errorf("RemoveMemorialRole: %w", err)
	}
	if existingUser == nil {
		return fmt.Errorf("RemoveMemorialRole: %w", errmgr.ErrUserNotFound)
	}

	// Check if user already has a role at the given Memorial
	existingMemorialRole, err := d.FindUserMemorialRoleByUserID(
		db,         //db 			*gorm.DB
		FSPID,      //FSPID 		uint
		userID,     //userID 		uint
		memorialID, //memorialID 	uint
	)
	if err != nil {
		return fmt.Errorf("RemoveMemorialRole: %w", err)
	}
	if existingMemorialRole == nil {
		return fmt.Errorf("RemoveMemorialRole: %w", ErrRoleNotFound)
	}

	err = db.Delete(&existingMemorialRole).Error
	if err != nil {
		return fmt.Errorf("RemoveMemorialRole: %w", err)
	}

	return nil
}

// Removes a given UserMemorialRole by the UserMemorialRole ID.
// Returns an error if no UserMemorialRole was found.
func (d *Domain) RemoveUserMemorialRoleByID(db *gorm.DB, FSPID uint, userMemorialRoleID uint) error {
	// If db is nil, get the regular DB instance
	// this allows the function to be used in transactions
	if db == nil {
		db = d.params.DB.GetDB()
	}

	result := db.
		Unscoped(). //explicitly allow hard delete
		Where("fsp_id = ?", FSPID).
		Where("id = ?", userMemorialRoleID).
		Delete(&schema.UserMemorialRole{})
	if result.Error != nil {
		return fmt.Errorf("RemoveMemorialRoleByID: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("RemoveMemorialRoleByID: %w", ErrRoleNotFound)
	}

	return nil
}

func (d *Domain) QueryRolesByLevel(FSPID uint, userID uint) (*rolesByLevel, error) {
	db := d.params.DB.GetDB()

	var userFSPRole schema.UserFSPRole
	var userMemorialRoles []schema.UserMemorialRole

	// Fetch FSP role
	err := db.Model(&schema.UserFSPRole{}).
		Where("fsp_id = ?", FSPID).
		Where("user_id = ?", userID).
		Preload("FSPRole").
		Preload("FSP").
		First(&userFSPRole).
		Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("QueryRolesByLevel: %w", err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newFSPRole, err := d.AssignOrUpdateFSPRole(
			db,                 // db 				*gorm.DB
			FSPID,              // FSPID 			uint
			userID,             // userID 			uint
			schema.RoleFSPUser, // roleToAssign 	schema.FSPRoleConst
		)
		if err != nil {
			return nil, fmt.Errorf("QueryRolesByLevel: %w", err)
		}
		if newFSPRole == nil {
			return nil, fmt.Errorf("QueryRolesByLevel: %s", "newFSPRole is nil")
		}
		userFSPRole = *newFSPRole
	}

	// Fetch Memorial roles
	err = db.Model(&schema.UserMemorialRole{}).
		Where("fsp_id = ?", FSPID).
		Where("user_id = ?", userID).
		Preload("MemorialRole").
		Preload("Memorial").
		Find(&userMemorialRoles).
		Error
	if err != nil {
		return nil, fmt.Errorf("QueryRolesByLevel: %w", err)
	}

	rolesByLevel := &rolesByLevel{
		FSP:      userFSPRole,
		Memorial: userMemorialRoles,
	}

	return rolesByLevel, nil
}

func (d *Domain) ResolveFSPRoleID(roleName schema.FSPRoleConst) (*uint, error) {
	db := d.params.DB.GetDB()

	var role schema.FSPRole
	err := db.Model(&schema.FSPRole{}).
		Where("name = ?", roleName).
		First(&role).
		Error
	if err != nil {
		return nil, fmt.Errorf("ResolveFSPRoleID: %w", err)
	}

	return &role.ID, nil
}

func (d *Domain) ResolveMemorialRoleID(roleName schema.MemorialRoleConst) (*uint, error) {
	db := d.params.DB.GetDB()

	var role schema.MemorialRole
	err := db.Model(&schema.MemorialRole{}).
		Where("name = ?", roleName).
		First(&role).
		Error
	if err != nil {
		return nil, fmt.Errorf("ResolveMemorialRoleID: %w", err)
	}

	return &role.ID, nil
}

// Existence Checks ------------------------------------------------

// Find a user's UserFSPRole. **If the UserFSPRole is not found, it returns nil without an error**.
func (d *Domain) FindUserFSPRole(db *gorm.DB, FSPID uint, userID uint) (*schema.UserFSPRole, error) {
	var err error

	// If db is nil, get the regular DB instance
	// this allows the function to be used in transactions
	if db == nil {
		db = d.params.DB.GetDB()
	}

	var userRole schema.UserFSPRole
	err = db.
		Model(&schema.UserFSPRole{}).
		Where("fsp_id = ?", FSPID).
		Where("user_id = ?", userID).
		First(&userRole).
		Error

	if err != nil {
		// return nil without an error if the record is not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("findFSPLevelUserRole: %w", err)
	}

	return &userRole, nil
}

// Find a user's UserMemorialRole for the given memorial. **If the UserMemorialRole is not found, it returns nil without an error**.
func (d *Domain) FindUserMemorialRoleByUserID(db *gorm.DB, FSPID uint, userID uint, memorialID uint) (*schema.UserMemorialRole, error) {
	var err error

	// If db is nil, get the regular DB instance
	// this allows the function to be used in transactions
	if db == nil {
		db = d.params.DB.GetDB()
	}

	var userRole schema.UserMemorialRole
	err = db.
		Model(&schema.UserMemorialRole{}).
		Where("fsp_id = ?", FSPID).
		Where("user_id = ?", userID).
		Where("memorial_id = ?", memorialID).
		First(&userRole).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("findMemorialLevelUserRoles: %w", err)
	}

	return &userRole, nil
}

// Find a UserMemorialRole by the UserMemorialRole ID. **If the UserMemorialRole is not found, it returns nil without an error**.
func (d *Domain) FindUserMemorialRoleByID(db *gorm.DB, FSPID uint, userMemorialRoleID uint) (*schema.UserMemorialRole, error) {
	var err error

	// If db is nil, get the regular DB instance
	// this allows the function to be used in transactions
	if db == nil {
		db = d.params.DB.GetDB()
	}

	var userRole schema.UserMemorialRole
	err = db.
		Model(&schema.UserMemorialRole{}).
		Where("fsp_id = ?", FSPID).
		Where("id = ?", userMemorialRoleID).
		First(&userRole).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("findMemorialLevelUserRoles: %w", err)
	}

	return &userRole, nil
}

// Find a user by email. **If the user is not found, it returns nil without an error**.
func (d *Domain) FindUserByEmail(db *gorm.DB, FSPID uint, email string) (*schema.User, error) {
	var err error

	// If db is nil, get the regular DB instance
	// this allows the function to be used in transactions
	if db == nil {
		db = d.params.DB.GetDB()
	}

	var user schema.User
	err = db.
		Model(&schema.User{}).
		Where("fsp_id = ?", FSPID).
		Where("email = ?", email).
		First(&user).
		Error
	if err != nil {
		// return nil without an error if the record is not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("findUserByEmail: %w", err)
	}

	return &user, nil
}

// Find a user by ID. **If the user is not found, it returns nil without an error**.
func (d *Domain) FindUserByID(db *gorm.DB, FSPID uint, userID uint) (*schema.User, error) {
	var err error

	// If db is nil, get the regular DB instance
	// this allows the function to be used in transactions
	if db == nil {
		db = d.params.DB.GetDB()
	}

	var user schema.User
	err = db.
		Model(&schema.User{}).
		Where("fsp_id = ?", FSPID).
		Where("id = ?", userID).
		First(&user).
		Error

	if err != nil {
		// return nil without an error if the record is not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("findUserByID: %w", err)
	}

	return &user, nil
}

// Find a user's invitation by the invitation token. **If the invitation is not found, it returns nil without an error**.
func (d *Domain) FindInvitationByToken(FSPID uint, memorialID uint, token string) (*schema.Invitation, error) {
	db := d.params.DB.GetDB()

	var invitation schema.Invitation
	err := db.
		Model(&schema.Invitation{}).
		Where("fsp_id = ?", FSPID).
		Where("memorial_id = ?", memorialID).
		Where("token = ?", token).
		First(&invitation).
		Error
	if err != nil {
		// return nil without an error if the record is not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("findInvitationByToken: %w", err)
	}

	return &invitation, nil
}

// Find a user's invitation by the invitee's email. **If the invitation is not found, it returns nil without an error**.
func (d *Domain) FindInvitationByEmail(FSPID uint, memorialID uint, email string) (*schema.Invitation, error) {
	db := d.params.DB.GetDB()

	var invitation schema.Invitation
	err := db.
		Model(&schema.Invitation{}).
		Where("fsp_id = ?", FSPID).
		Where("invitee_email = ?", email).
		First(&invitation).
		Error
	if err != nil {
		// return nil without an error if the record is not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("findInvitationByEmail: %w", err)
	}

	return &invitation, nil
}

// Find a contributor invitation by the invitee's email. **If the invitation is not found, it returns nil without an error**.
func (d *Domain) FindContributorInvitationByEmail(FSPID uint, memorialID uint, email string) (*schema.Invitation, error) {
	var err error

	db := d.params.DB.GetDB()

	existingInvitation := schema.Invitation{}

	err = db.
		Where("fsp_id = ?", FSPID).
		Where("memorial_id = ?", memorialID).
		Where("invitee_email = ?", email).
		First(&existingInvitation).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("findContributorInvitationByEmail: %w", err)
	}

	return &existingInvitation, nil
}

// Find a contributor invitation by its ID. **If the invitation is not found, it returns nil without an error**.
func (d *Domain) FindContributorInvitationByID(FSPID uint, invitationID uint) (*schema.Invitation, error) {
	var err error

	db := d.params.DB.GetDB()

	existingInvitation := schema.Invitation{}

	err = db.
		Where("fsp_id = ?", FSPID).
		Where("id = ?", invitationID).
		First(&existingInvitation).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("findContributorInvitationByID: %w", err)
	}

	return &existingInvitation, nil
}

// Find a contributor application by its ID. **If the application is not found, it returns nil without an error**.
func (d *Domain) FindContributorApplicationByID(FSPID uint, applicationID uint) (*schema.Application, error) {
	var err error

	db := d.params.DB.GetDB()

	existingApplication := schema.Application{}

	err = db.
		Where("fsp_id = ?", FSPID).
		Where("id = ?", applicationID).
		First(&existingApplication).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("findContributorApplicationByID: %w", err)
	}

	return &existingApplication, nil
}

// Find a contributor application by the applicant's email. **If the application is not found, it returns nil without an error**.
func (d *Domain) FindContributorApplicationByEmail(FSPID uint, memorialID uint, email string) (*schema.Application, error) {
	var err error

	db := d.params.DB.GetDB()

	var existingApplication schema.Application

	// Joining the User table and filtering by User.Email
	err = db.
		Joins("JOIN users ON users.id = applications.applicant_id").
		Where("applications.fsp_id = ?", FSPID).
		Where("applications.memorial_id = ?", memorialID).
		Where("users.email = ?", email).
		First(&existingApplication).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("findContributorApplicationByEmail: %w", err)
	}

	return &existingApplication, nil
}

// Find a contributor application by the applicant's ID. **If the application is not found, it returns nil without an error**.
func (d *Domain) FindContributorApplicationByApplicantID(FSPID uint, memorialID uint, applicantID uint) (*schema.Application, error) {
	var err error

	db := d.params.DB.GetDB()

	var existingApplication schema.Application

	err = db.
		Where("fsp_id = ?", FSPID).
		Where("memorial_id = ?", memorialID).
		Where("applicant_id = ?", applicantID).
		First(&existingApplication).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("findContributorApplicationByApplicantID: %w", err)
	}

	return &existingApplication, nil
}

// Find a memorial by its ID. **If the memorial is not found, it returns nil without an error**.
func (d *Domain) FindMemorialByID(FSPID uint, memorialID uint) (*schema.Memorial, error) {
	var err error

	db := d.params.DB.GetDB()

	var existingMemorial schema.Memorial

	err = db.
		Where("fsp_id = ?", FSPID).
		Where("id = ?", memorialID).
		First(&existingMemorial).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("findMemorialByID: %w", err)
	}

	return &existingMemorial, nil
}

// Checkers -------------------------------------------------------
func (d *Domain) UserHasMemorialRole(FSPID uint, userID uint, memorialID uint, role schema.MemorialRoleConst) (bool, error) {
	db := d.params.DB.GetDB()

	existingUserMemorialRole := schema.UserMemorialRole{}

	err := db.
		Where("fsp_id = ?", FSPID).
		Where("user_id = ?", userID).
		Where("memorial_id = ?", memorialID).
		Joins("JOIN memorial_roles ON memorial_roles.id = user_memorial_roles.memorial_role_id").
		Where("memorial_roles.name = ?", role).
		First(&existingUserMemorialRole).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("UserHasMemorialRole: %w", err)
	}

	return true, nil
}

func (d *Domain) UserHasFSPRole(FSPID uint, userID uint, role schema.FSPRoleConst) (bool, error) {
	db := d.params.DB.GetDB()

	existingUserFSPRole := schema.UserFSPRole{}

	err := db.
		Where("fsp_id = ?", FSPID).
		Where("user_id = ?", userID).
		Preload("FSPRole").
		First(&existingUserFSPRole).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("UserHasFSPRole: %w", err)
	}

	if existingUserFSPRole.FSPRole.Name == role {
		return true, nil
	}

	return false, nil
}
