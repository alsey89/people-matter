package memorial

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/alsey89/people-matter/internal/common/errmgr"
	"github.com/alsey89/people-matter/internal/schema"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CURATOR --------------------------------------------------------

// contributors

func (d *Domain) serviceGetContributors(TenantID uint, memorialID uint) ([]schema.UserMemorialRole, error) {
	var err error

	db := d.params.DB.GetDB()

	existingUserMemorialRoles := []schema.UserMemorialRole{}

	err = db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Joins("JOIN memorial_roles ON memorial_roles.id = user_memorial_roles.memorial_role_id").
		Where("memorial_roles.name IN ?", []schema.MemorialRoleConst{schema.RoleMemContributor}).
		Preload("User").
		Preload("MemorialRole").
		Find(&existingUserMemorialRoles).
		Error
	if err != nil {
		return nil, fmt.Errorf("getTeamService: %w", err)
	}

	return existingUserMemorialRoles, nil
}
func (d *Domain) serviceInviteContributor(TenantID uint, memorialID uint, adminID uint, contributorEmail string, relationship schema.RelationshipConst) error {
	var err error

	// Check if the memorial exists
	memorial, err := d.FindMemorialByID(memorialID)
	if err != nil {
		return fmt.Errorf("serviceInviteContributor: %w", err)
	}
	if memorial == nil {
		return fmt.Errorf("serviceInviteContributor: Memorial not found")
	}

	// Check if user account already exists for Tenant
	user, err := d.params.Identity.FindUserByEmail(
		nil,              // db 		*gorm.DB
		TenantID,         // TenantID 		uint
		contributorEmail, // email 		string
	)
	if err != nil {
		return fmt.Errorf("serviceInviteContributor: %w", err)
	}

	// Admin cannot invite themselves
	if user != nil && user.ID == adminID {
		return fmt.Errorf("serviceInviteContributor: %w", errmgr.ErrUserIsCurator)
	}

	// Check if user already has a role in the memorial
	if user != nil {
		userMemorialRoles, err := d.params.Identity.QueryRolesByLevel(
			TenantID, // TenantID 	uint
			user.ID,  // userID 	uint
		)
		if err != nil {
			return fmt.Errorf("serviceInviteContributor: %w", err)
		}
		if userMemorialRoles != nil {
			for _, role := range userMemorialRoles.Memorial {
				if role.MemorialID == memorialID {
					return fmt.Errorf("serviceInviteContributor: %w", errmgr.ErrUserHasMemorialRole)
				}
			}
		}
	}

	// Check if the user has an invitation to the memorial
	invitation, err := d.params.Identity.FindContributorInvitationByEmail(TenantID, memorialID, contributorEmail)
	if err != nil {
		return fmt.Errorf("serviceInviteContributor: %w", err)
	}
	if invitation != nil {
		return fmt.Errorf("serviceInviteContributor: %w", errmgr.ErrUserHasInvitation)
	}

	createdInvitation, err := d.params.Identity.CreateMemorialContributorInvitation(
		nil,              // db 					*gorm.DB
		TenantID,         // TenantID 					uint
		adminID,          // adminID 				uint
		contributorEmail, // email 					string
		relationship,     // relationship 			schema.RelationshipConst
		&memorialID,      // memorialID 			uint
	)
	if err != nil {
		return fmt.Errorf("serviceInviteContributor: %w", err)
	}

	urlPath := fmt.Sprintf("/memorial/%v/invitation/%v/confirm?token=%s", memorial.ID, createdInvitation.ID, createdInvitation.Token)
	variables := map[string]interface{}{
		"memorial_title": memorial.Title,
	}

	//TODO: ??notify the user that they have been invited to the memorial??
	//TODO: if user already has account send email with link to accept invitation
	//TODO: if user does not have account send email with link to create account and accept invitation
	go d.params.TransMail.SendMail(
		TenantID,         // TenantID 				uint
		contributorEmail, // recipientEmail 	string
		6500238,          // templateID 		int
		&urlPath,         // urlPath 			*string
		variables,        //variables 			map[string]interface{},
	)

	return nil
}
func (d *Domain) serviceReinviteContributor(TenantID uint, memorialID uint, invitationID uint) error {
	var err error

	// Check if the memorial exists
	memorial, err := d.FindMemorialByID(memorialID)
	if err != nil {
		return fmt.Errorf("serviceInviteContributor: %w", err)
	}
	if memorial == nil {
		return fmt.Errorf("serviceInviteContributor: Memorial not found")
	}

	// Check if the invitation exists
	existingInvitation, err := d.params.Identity.FindContributorInvitationByID(
		TenantID,     // TenantID 			uint
		invitationID, // invitationID 	uint
	)
	if err != nil {
		return fmt.Errorf("serviceInviteContributor: %w", err)
	}
	if existingInvitation == nil {
		return fmt.Errorf("serviceInviteContributor: %w", errmgr.ErrInvitationNotFound)
	}

	// refresh the invitation
	refreshedInvitation, err := d.params.Identity.RefreshContributorInvitationAndToken(
		nil,          // db 					*gorm.DB
		TenantID,     // TenantID 					uint
		memorialID,   // memorialID 			uint
		invitationID, // invitationID 			uint
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("serviceInviteContributor: %w", errmgr.ErrInvitationNotFound)
		}
		return fmt.Errorf("serviceInviteContributor: %w", err)
	}

	//reinvite the user
	urlPath := fmt.Sprintf("/memorial/%v/invitation/%v/confirm?token=%s", memorial.ID, refreshedInvitation.ID, refreshedInvitation.Token)
	variables := map[string]interface{}{
		"memorial_title": memorial.Title,
	}

	go d.params.TransMail.SendMail(
		TenantID,                         // TenantID 				uint
		refreshedInvitation.InviteeEmail, // recipientEmail 		string
		6500238,                          // templateID 			int
		&urlPath,                         // urlPath 			*string
		variables,                        //variables 			map[string]interface{},
	)

	return nil
}
func (d *Domain) serviceDeleteContributor(TenantID uint, memorialID uint, contributorMemorialRoleID uint, userIsNotified bool) error {
	var err error

	db := d.params.DB.GetDB()

	// Check if the memorial exists
	existingMemorial, err := d.FindMemorialByID(memorialID)
	if err != nil {
		return fmt.Errorf("serviceDeleteContributor: %w", err)
	}
	if existingMemorial == nil {
		return fmt.Errorf("serviceDeleteContributor: Memorial not found")
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// Remove the user from the memorial ROLE
		err = d.params.Identity.RemoveUserMemorialRoleByID(
			tx,                        // db 					*gorm.DB
			TenantID,                  // TenantID 				uint
			contributorMemorialRoleID, // userMemorialRoleID 	uint
		)
		if err != nil {
			return fmt.Errorf("serviceDeleteContributor: %w", err)
		}

		if userIsNotified {
			//TODO: ??notify the user that they have been removed from the memorial??
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("serviceDeleteContributor: %w", err)
	}

	return nil
}

// applications

func (d *Domain) serviceGetContributorApplications(TenantID uint, memorialID uint) ([]schema.Application, error) {
	var err error

	db := d.params.DB.GetDB()

	existingApplications := []schema.Application{}

	err = db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Preload("Applicant").
		Find(&existingApplications).
		Error
	if err != nil {
		return nil, fmt.Errorf("serviceGetContributorApplicants: %w", err)
	}

	return existingApplications, nil
}
func (d *Domain) serviceAcceptContributorApplication(TenantID uint, memorialID uint, applicationID uint) error {
	var err error

	// check if the memorial exists
	memorial, err := d.FindMemorialByID(memorialID)
	if err != nil {
		return fmt.Errorf("serviceAcceptContributorApplication: %w", err)
	}
	if memorial == nil {
		return fmt.Errorf("serviceAcceptContributorApplication: %w", errmgr.ErrMemorialNotFound)
	}

	// check if the application exists
	existingApplication, err := d.params.Identity.FindContributorApplicationByID(
		TenantID,      // TenantID 			uint
		applicationID, // applicationID 	uint
	)
	if err != nil {
		return fmt.Errorf("serviceAcceptContributorApplication: %w", err)
	}
	if existingApplication == nil {
		return fmt.Errorf("serviceAcceptContributorApplication: %w", errmgr.ErrApplicationNotFound)
	}

	// check if the user is already a contributor
	userMemorialRoles, err := d.params.Identity.QueryRolesByLevel(
		TenantID,                        // TenantID 	uint
		existingApplication.ApplicantID, // userID 	uint
	)
	if err != nil {
		return fmt.Errorf("serviceAcceptContributorApplication: %w", err)
	}
	if userMemorialRoles != nil {
		for _, role := range userMemorialRoles.Memorial {
			if role.MemorialID == memorialID {
				return fmt.Errorf("serviceAcceptContributorApplication: %w", errmgr.ErrUserHasMemorialRole)
			}
		}
	}

	// assign memorial contributor role
	_, err = d.params.Identity.AssignOrUpdateMemorialRole(
		nil,                               // db 					*gorm.DB
		TenantID,                          // TenantID 				uint
		existingApplication.ApplicantID,   // userID 				uint
		memorialID,                        // memorialID 			uint
		schema.RoleMemContributor,         // roleName 			schema.MemorialRoleConst
		&existingApplication.Relationship, // relationship 		*schema.RelationshipConst
	)
	if err != nil {
		return fmt.Errorf("serviceAcceptContributorApplication: %w", err)
	}

	// update the application status
	existingApplication.Status = schema.ApplicationStatusAccepted
	err = d.params.DB.GetDB().
		Updates(&existingApplication).
		Error
	if err != nil {
		return fmt.Errorf("serviceAcceptContributorApplication: %w", err)
	}

	// TODO: notify the user that their application has been accepted

	return nil
}
func (d *Domain) serviceRejectContributorApplication(TenantID uint, memorialID uint, applicationID uint) error {
	var err error

	db := d.params.DB.GetDB()

	// check if the memorial exists
	memorial, err := d.FindMemorialByID(memorialID)
	if err != nil {
		return fmt.Errorf("serviceRejectContributorApplication: %w", err)
	}
	if memorial == nil {
		return fmt.Errorf("serviceRejectContributorApplication: %w", errmgr.ErrMemorialNotFound)
	}

	// check if the application exists
	application, err := d.params.Identity.FindContributorApplicationByID(
		TenantID,      // TenantID 			uint
		applicationID, // applicationID 	uint
	)
	if err != nil {
		return fmt.Errorf("serviceRejectContributorApplication: %w", err)
	}
	if (application.Status == schema.ApplicationStatusDeclined) || (application.Status == schema.ApplicationStatusAccepted) {
		return fmt.Errorf("serviceRejectContributorApplication: %w", errmgr.ErrApplicationResponded)
	}

	// update the application status
	application.Status = schema.ApplicationStatusDeclined
	err = db.
		Updates(&application).
		Error
	if err != nil {
		return fmt.Errorf("serviceRejectContributorApplication: %w", err)
	}

	// TODO: notify the user that their application has been rejected

	return nil
}

// invitations

func (d *Domain) serviceGetContributorInvitations(TenantID uint, memorialID uint) ([]schema.Invitation, error) {
	var err error

	db := d.params.DB.GetDB()

	existingInvitations := []schema.Invitation{}

	err = db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Preload("Inviter").
		Find(&existingInvitations).
		Error
	if err != nil {
		return nil, fmt.Errorf("serviceGetContributorInvitations: %w", err)
	}

	return existingInvitations, nil
}
func (d *Domain) serviceDeleteContributorInvitation(TenantID uint, memorialID uint, invitationID uint) error {
	var err error

	db := d.params.DB.GetDB()

	err = db.Transaction(func(tx *gorm.DB) error {
		// Check if the invitation exists
		existingInvitation := schema.Invitation{}
		err = db.
			Where("fsp_id = ?", TenantID).
			Where("memorial_id = ?", memorialID).
			Where("id = ?", invitationID).
			First(&existingInvitation).
			Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("deleteContributorInvitation: Invitation not found")
			}
			return fmt.Errorf("deleteContributorInvitation: %w", err)
		}

		// Remove the invitation
		err = db.
			Unscoped().
			Delete(&existingInvitation).
			Error
		if err != nil {
			return fmt.Errorf("deleteContributorInvitation: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("deleteContributorInvitation: %w", err)
	}

	return nil
}

// contributions
func (d *Domain) serviceGetMemorialContributions(TenantID uint, memorialID uint) (getContribtionsResponse, error) {
	var err error

	db := d.params.DB.GetDB()
	allMemorialContributions := getContribtionsResponse{}

	// Get condolence elements
	condolenceElements := []schema.ContributionCondolenceElement{}
	err = db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Preload("Contributor").
		Preload("ContributorMemorialRole").
		Find(&condolenceElements).
		Error
	if err != nil {
		return allMemorialContributions, fmt.Errorf("serviceGetMemorialContributions: %w", err)
	}

	// Get and sign gallery elements
	galleryElements := []schema.ContributionGalleryElement{}
	err = db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Preload("Contributor").
		Preload("ContributorMemorialRole").
		Find(&galleryElements).
		Error
	if err != nil {
		return allMemorialContributions, fmt.Errorf("serviceGetMemorialContributions: %w", err)
	}
	for i := range galleryElements {
		galleryElements[i].ElementMediaURL, err = d.signS3GetURL(
			galleryElements[i].ElementMediaURL, // url 				string
			true,                               // signThumbnail 	bool
		)
		if err != nil {
			return allMemorialContributions, fmt.Errorf("serviceGetMemorialContributions: %w", err)
		}
	}

	// Get and sign story elements
	storyElements := []schema.ContributionStoryElement{}
	err = db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Preload("Contributor").
		Preload("ContributorMemorialRole").
		Find(&storyElements).
		Error
	if err != nil {
		return allMemorialContributions, fmt.Errorf("serviceGetMemorialContributions: %w", err)
	}
	for i := range storyElements {
		storyElements[i].ElementMediaURL, err = d.signS3GetURL(
			storyElements[i].ElementMediaURL, // url 				string
			true,                             // signThumbnail 	bool
		)
		if err != nil {
			return allMemorialContributions, fmt.Errorf("serviceGetMemorialContributions: %w", err)
		}
	}

	// Get and optionally sign timeline elements
	timelineElements := []schema.ContributionTimelineElement{}
	err = db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Preload("Contributor").
		Preload("ContributorMemorialRole").
		Find(&timelineElements).
		Error
	if err != nil {
		return allMemorialContributions, fmt.Errorf("serviceGetMemorialContributions: %w", err)
	}
	for i := range timelineElements {
		if timelineElements[i].ElementMediaURL != nil {
			signedURL, err := d.signS3GetURL(
				*timelineElements[i].ElementMediaURL, // url 				string
				true,                                 // signThumbnail 	bool
			)
			if err != nil {
				return allMemorialContributions, err
			}
			timelineElements[i].ElementMediaURL = &signedURL
		}
	}

	allMemorialContributions.CondolenceElements = condolenceElements
	allMemorialContributions.GalleryElements = galleryElements
	allMemorialContributions.StoryElements = storyElements
	allMemorialContributions.TimelineElements = timelineElements

	return allMemorialContributions, nil
}
func (d *Domain) serviceUpdateContributionCondolenceElementState(TenantID uint, memorialID uint, elementID uint, state schema.ContributionStateConst) error {
	db := d.params.DB.GetDB()

	// Check if the contribution condolence element exists
	var existingElement schema.ContributionCondolenceElement
	err := db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Where("id = ?", elementID).
		First(&existingElement).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("serviceUpdateContributionCondolenceState: %w", errmgr.ErrContributionElementNotFound)
		}
		return fmt.Errorf("serviceUpdateContributionCondolenceState: %w", err)
	}

	// Update the state
	existingElement.ContributionState = state
	err = db.
		Select("ContributionState").
		Updates(&existingElement).
		Error
	if err != nil {
		return fmt.Errorf("serviceUpdateContributionCondolenceState: %w", err)
	}

	return nil
}
func (d *Domain) serviceUpdateContributionGalleryElementState(TenantID uint, memorialID uint, elementID uint, state schema.ContributionStateConst) error {
	db := d.params.DB.GetDB()

	// Check if the contribution gallery element exists
	var existingElement schema.ContributionGalleryElement
	err := db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Where("id = ?", elementID).
		First(&existingElement).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("serviceUpdateContributionGalleryState: %w", errmgr.ErrContributionElementNotFound)
		}
		return fmt.Errorf("serviceUpdateContributionGalleryState: %w", err)
	}

	// Update the state
	existingElement.ContributionState = state
	err = db.
		Select("ContributionState").
		Updates(&existingElement).
		Error
	if err != nil {
		return fmt.Errorf("serviceUpdateContributionGalleryState: %w", err)
	}

	return nil
}
func (d *Domain) serviceUpdateContributionTimelineElementState(TenantID uint, memorialID uint, elementID uint, state schema.ContributionStateConst) error {
	db := d.params.DB.GetDB()

	// Check if the contribution timeline element exists
	var existingElement schema.ContributionTimelineElement
	err := db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Where("id = ?", elementID).
		First(&existingElement).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("serviceUpdateContributionTimelineState: %w", errmgr.ErrContributionElementNotFound)
		}
		return fmt.Errorf("serviceUpdateContributionTimelineState: %w", err)
	}

	// Update the state
	existingElement.ContributionState = state
	err = db.
		Select("ContributionState").
		Updates(&existingElement).
		Error
	if err != nil {
		return fmt.Errorf("serviceUpdateContributionTimelineState: %w", err)
	}

	return nil
}
func (d *Domain) serviceUpdateContributionStoryElementState(TenantID uint, memorialID uint, elementID uint, state schema.ContributionStateConst) error {
	db := d.params.DB.GetDB()

	// Check if the contribution story element exists
	var existingElement schema.ContributionStoryElement
	err := db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Where("id = ?", elementID).
		First(&existingElement).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("serviceUpdateContributionStoryState: %w", errmgr.ErrContributionElementNotFound)
		}
		return fmt.Errorf("serviceUpdateContributionStoryState: %w", err)
	}

	// Update the state
	existingElement.ContributionState = state
	err = db.
		Select("ContributionState").
		Updates(&existingElement).
		Error
	if err != nil {
		return fmt.Errorf("serviceUpdateContributionStoryState: %w", err)
	}

	return nil
}

// publish/export
func (d *Domain) serviceExportMemorial(TenantID uint, memorialID uint) (*uint, error) {
	var err error

	db := d.params.DB.GetDB()

	// Check if the memorial exists
	memorial, err := d.FindMemorialByID(memorialID)
	if err != nil {
		return nil, fmt.Errorf("serviceExportMemorial: %w", err)
	}
	if memorial == nil {
		return nil, fmt.Errorf("serviceExportMemorial: %w", errmgr.ErrMemorialNotFound)
	}

	var createdExportID *uint

	err = db.Transaction(func(tx *gorm.DB) error {

		// Check if there is are incomplete exports
		incompleteExports, err := d.FindAndHandleIncompleteExports(tx, TenantID, memorialID)
		if err != nil {
			return fmt.Errorf("serviceExportMemorial: %w", err)
		}
		if incompleteExports != nil {
			return fmt.Errorf("serviceExportMemorial: %w", errmgr.ErrExportInProgress)
		}

		// add entry to the export table
		export := schema.Export{
			TenantID:   TenantID,
			MemorialID: memorialID,
			State:      schema.ExportStateRequested,
		}
		err = tx.
			Create(&export).
			Error
		if err != nil {
			return fmt.Errorf("serviceExportMemorial: %w", err)
		}

		createdExportID = &export.ID

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("serviceExportMemorial: %w", err)
	}

	//todo: generate a shortlived token

	//pass the URL and token to the expot pipeline
	variables := map[string]string{
		"NUXT_PUBLIC_API_URL":     "https://sci.curate.memorial/api/v1",
		"NUXT_PUBLIC_MEMORIAL_ID": "2",
		"NUXT_PUBLIC_FSP_ID":      "1",
		"AUTH_TOKEN":              "my-secret-token-from-backend",
		"EXPORT_ID":               "1",
	}
	err = d.triggerExportPipeline(
		variables, // variables 	map[string]interface{}
	)
	if err != nil {
		return nil, fmt.Errorf("serviceExportMemorial: %w", err)
	}

	return createdExportID, nil
}
func (d *Domain) serviceGetAllExports(TenantID uint, memorialID uint) ([]schema.Export, error) {
	var err error

	db := d.params.DB.GetDB()

	exports := []schema.Export{}

	err = db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Order("id DESC"). // Move Order before Find
		Find(&exports).
		Error
	if err != nil {
		return nil, fmt.Errorf("serviceGetAllExports: %w", err)
	}

	return exports, nil
}
func (d *Domain) serviceGetExportState(TenantID uint, memorialID uint, exportID uint) (*schema.ExportStateConst, error) {
	var err error

	db := d.params.DB.GetDB()

	export := schema.Export{}

	err = db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Where("id = ?", exportID).
		First(&export).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("serviceGetExportState: %w", errmgr.ErrExportNotFound)
		}
		return nil, fmt.Errorf("serviceGetExportState: %w", err)
	}

	return &export.State, nil
}
func (d *Domain) serviceUpdateExportState(TenantID uint, memorialID uint, exportID uint, state schema.ExportStateConst) error {
	var err error

	db := d.params.DB.GetDB()

	// Check if the export exists
	export := schema.Export{}
	err = db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Where("id = ?", exportID).
		First(&export).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("serviceUpdateExportState: %w", errmgr.ErrExportNotFound)
		}
		return fmt.Errorf("serviceUpdateExportState: %w", err)
	}

	// Update the state
	export.State = state
	err = db.
		Select("State").
		Updates(&export).
		Error
	if err != nil {
		return fmt.Errorf("serviceUpdateExportState: %w", err)
	}

	return nil
}

// CONTRIBUTOR -----------------------------------------------------

func (d *Domain) serviceGetContributions(TenantID uint, memorialID uint, contributorID uint) (getContribtionsResponse, error) {
	var err error
	db := d.params.DB.GetDB()
	response := getContribtionsResponse{}

	// Helper inline function to sign URLs if they exist
	signURL := func(url string, signThumbnail bool) (string, error) {
		if url == "" {
			return url, nil
		}

		if signThumbnail {
			convertedURL, err := d.params.S3.ConvertObjectURLOGToThumb(url)
			if err != nil {
				return "", fmt.Errorf("serviceGetContributions: %w", err)
			}
			if convertedURL == nil {
				return "", fmt.Errorf("serviceGetContributions: %s", errmgr.ErrNilCheckFailed)
			}
			url = *convertedURL
		}

		signedURL, err := d.params.S3.GetPresignedReadURL(
			context.Background(), // ctx	context.Context
			url,                  // s3URL	string
			24*time.Hour,         // exp	time.Duration
		)
		if err != nil {
			return "", fmt.Errorf("serviceGetContributions: signURL: %w", err)
		}
		if signedURL == nil {
			return "", fmt.Errorf("serviceGetContributions: signURL: %s", errmgr.ErrNilCheckFailed)
		}

		return *signedURL, nil
	}

	// Get condolence elements
	condolenceElements := []schema.ContributionCondolenceElement{}
	err = db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Where("contributor_id = ?", contributorID).
		Find(&condolenceElements).
		Error
	if err != nil {
		return response, fmt.Errorf("serviceGetContributions: %w", err)
	}

	// Get and sign gallery elements
	galleryElements := []schema.ContributionGalleryElement{}
	err = db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Where("contributor_id = ?", contributorID).
		Find(&galleryElements).
		Error
	if err != nil {
		return response, fmt.Errorf("serviceGetContributions: %w", err)
	}
	for i := range galleryElements {
		galleryElements[i].ElementMediaURL, err = signURL(
			galleryElements[i].ElementMediaURL, // url 				string
			true,                               // signThumbnail 	bool
		)
		if err != nil {
			return response, err
		}
	}

	// Get and sign story elements
	storyElements := []schema.ContributionStoryElement{}
	err = db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Where("contributor_id = ?", contributorID).
		Find(&storyElements).
		Error
	if err != nil {
		return response, fmt.Errorf("serviceGetContributions: %w", err)
	}
	for i := range storyElements {
		storyElements[i].ElementMediaURL, err = signURL(
			storyElements[i].ElementMediaURL, // url 				string
			true,                             // signThumbnail 	bool
		)
		if err != nil {
			return response, err
		}
	}

	// Get and optionally sign sign timeline elements
	timelineElements := []schema.ContributionTimelineElement{}
	err = db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Where("contributor_id = ?", contributorID).
		Find(&timelineElements).
		Error
	if err != nil {
		return response, fmt.Errorf("serviceGetContributions: %w", err)
	}
	for i := range timelineElements {
		if timelineElements[i].ElementMediaURL != nil {
			signedURL, err := signURL(
				*timelineElements[i].ElementMediaURL, // url 				string
				true,                                 // signThumbnail 	bool
			)
			if err != nil {
				return response, err
			}
			timelineElements[i].ElementMediaURL = &signedURL
		}
	}

	response.CondolenceElements = condolenceElements
	response.GalleryElements = galleryElements
	response.StoryElements = storyElements
	response.TimelineElements = timelineElements

	return response, nil
}
func (d *Domain) serviceGetPresignedUploadURL(TenantID uint, memorialID uint, contributorID uint) (*string, *map[string]string, error) {
	context := context.Background()

	data := S3KeyData{
		TenantID:      int(TenantID),
		MemorialID:    int(memorialID),
		ContributorID: int(contributorID),
		Date:          time.Now().Format("2006-01-02"),
		UUID:          uuid.NewString(),
	}

	var keyBuffer bytes.Buffer
	if err := d.s3ContribKeyTmpl.Execute(&keyBuffer, data); err != nil {
		return nil, nil, fmt.Errorf("serviceGetPresignedUploadURL: %w", err)
	}

	keyBufferString := keyBuffer.String()

	metadata := map[string]string{
		"fsp-id":         fmt.Sprint(TenantID),
		"memorial-id":    fmt.Sprint(memorialID),
		"contributor-id": fmt.Sprint(contributorID),
		"date":           time.Now().Format("2006-01-02"),
	}

	presignedUploadURL, err := d.params.S3.GetPresignedUploadURL(
		context,                       // ctx 				context.Context
		keyBufferString,               // key 				string
		d.config.S3PresignedURLExpiry, // exp 				time.Duration
		metadata,                      // metadata 			map[string]string
	)
	if err != nil {
		return nil, nil, fmt.Errorf("serviceGetPresignedUploadURL: %w", err)
	}
	if presignedUploadURL == nil {
		return nil, nil, fmt.Errorf("serviceGetPresignedUploadURL: %s", errmgr.ErrNilCheckFailed)
	}

	return presignedUploadURL, &metadata, nil
}

func (d *Domain) serviceCreateContributionGalleryElement(TenantID uint, memorialID uint, contributorID uint, isImmutable bool, elementTitle string, elementDescription string, elementDate time.Time, elementMediaType string, elementMediaURL string, elementLocation *string, ElementGooglePlaceID *string) error {
	var err error

	db := d.params.DB.GetDB()

	// Check if the memorial exists
	memorial, err := d.FindMemorialByID(memorialID)
	if err != nil {
		return fmt.Errorf("serviceCreateContributionGalleryElement: %w", err)
	}
	if memorial == nil {
		return fmt.Errorf("serviceCreateContributionGalleryElement: %w", errmgr.ErrMemorialNotFound)
	}

	// Check if the contributor exists
	contributor, err := d.params.Identity.FindUserByID(
		nil,           // db 		*gorm.DB
		TenantID,      // TenantID 	uint
		contributorID, // userID 	uint
	)
	if err != nil {
		return fmt.Errorf("serviceCreateContributionGalleryElement: %w", err)
	}
	if contributor == nil {
		return fmt.Errorf("serviceCreateContributionGalleryElement: %w", errmgr.ErrUserNotFound)
	}

	// Create the contribution gallery element
	newContributionGalleryElement := schema.ContributionGalleryElement{
		TenantID:          TenantID,
		MemorialID:        memorialID,
		ContributorID:     contributorID,
		ContributionState: schema.ContributionStatePending,
		IsImmutable:       isImmutable,

		ElementTitle:         elementTitle,
		ElementDescription:   elementDescription,
		ElementDate:          elementDate,
		ElementMediaType:     elementMediaType,
		ElementMediaURL:      elementMediaURL,
		ElementLocation:      elementLocation,
		ElementGooglePlaceID: ElementGooglePlaceID,
		HasEXIF:              false,
		UseEXIF:              false,
	}

	err = db.
		Create(&newContributionGalleryElement).
		Error
	if err != nil {
		return fmt.Errorf("serviceCreateContributionGalleryElement: %w", err)
	}

	return nil
}
func (d *Domain) serviceCreateContributionTimelineElement(TenantID uint, memorialID uint, contributorID uint, isImmutable bool, elementTitle string, elementDescription string, elementDate time.Time, elementEventType schema.EventTypeConst, elementMediaURL *string, elementLocation *string, ElementGooglePlaceID *string) error {
	var err error

	db := d.params.DB.GetDB()

	// Check if the memorial exists
	memorial, err := d.FindMemorialByID(memorialID)
	if err != nil {
		return fmt.Errorf("serviceCreateContributionTimelineElement: %w", err)
	}
	if memorial == nil {
		return fmt.Errorf("serviceCreateContributionTimelineElement: %w", errmgr.ErrMemorialNotFound)
	}

	// Check if the contributor exists
	contributor, err := d.params.Identity.FindUserByID(
		nil,           // db 		*gorm.DB
		TenantID,      // TenantID 	uint
		contributorID, // userID 	uint
	)
	if err != nil {
		return fmt.Errorf("serviceCreateContributionTimelineElement: %w", err)
	}
	if contributor == nil {
		return fmt.Errorf("serviceCreateContributionTimelineElement: %w", errmgr.ErrUserNotFound)
	}

	// Create the contribution timeline element
	newContributionTimelineElement := schema.ContributionTimelineElement{
		TenantID:          TenantID,
		MemorialID:        memorialID,
		ContributorID:     contributorID,
		ContributionState: schema.ContributionStatePending,
		IsImmutable:       isImmutable,

		ElementTitle:         elementTitle,
		ElementDescription:   elementDescription,
		ElementDate:          elementDate,
		ElementEventType:     elementEventType,
		ElementMediaURL:      elementMediaURL,
		ElementLocation:      elementLocation,
		ElementGooglePlaceID: ElementGooglePlaceID,
	}

	err = db.
		Create(&newContributionTimelineElement).
		Error
	if err != nil {
		return fmt.Errorf("serviceCreateContributionTimelineElement: %w", err)
	}

	return nil
}
func (d *Domain) serviceCreateContributionStoryElement(TenantID uint, memorialID uint, contributorID uint, isImmutable bool, elementTitle string, elementDescription string, elementMediaURL *string, elementAuthor string) error {
	var err error

	db := d.params.DB.GetDB()

	// Check if the memorial exists
	memorial, err := d.FindMemorialByID(memorialID)
	if err != nil {
		return fmt.Errorf("serviceCreateContributionStoryElement: %w", err)
	}
	if memorial == nil {
		return fmt.Errorf("serviceCreateContributionStoryElement: %w", errmgr.ErrMemorialNotFound)
	}

	// Check if the contributor exists
	contributor, err := d.params.Identity.FindUserByID(
		nil,           // db 		*gorm.DB
		TenantID,      // TenantID 	uint
		contributorID, // userID 	uint
	)
	if err != nil {
		return fmt.Errorf("serviceCreateContributionStoryElement: %w", err)
	}
	if contributor == nil {
		return fmt.Errorf("serviceCreateContributionStoryElement: %w", errmgr.ErrUserNotFound)
	}

	// Set the author to the contributor's name if not provided
	if elementAuthor == "" {
		elementAuthor = fmt.Sprintf("%s %s", contributor.FirstName, contributor.LastName)
	}

	// Create the contribution story element
	newContributionStoryElement := schema.ContributionStoryElement{
		TenantID:          TenantID,
		MemorialID:        memorialID,
		ContributorID:     contributorID,
		ContributionState: schema.ContributionStatePending,
		IsImmutable:       isImmutable,

		ElementTitle:       elementTitle,
		ElementDescription: elementDescription,
		ElementAuthor:      elementAuthor,
	}
	if elementMediaURL != nil {
		newContributionStoryElement.ElementMediaURL = *elementMediaURL
	}

	err = db.
		Create(&newContributionStoryElement).
		Error
	if err != nil {
		return fmt.Errorf("serviceCreateContributionStoryElement: %w", err)
	}

	return nil
}
func (d *Domain) serviceCreateContributionCondolenceElement(TenantID uint, memorialID uint, contributorID uint, isImmutable bool, elementTitle string, elementDescription string, elementAuthor string, designElementID string) error {
	var err error

	db := d.params.DB.GetDB()

	// Check if the memorial exists
	memorial, err := d.FindMemorialByID(memorialID)
	if err != nil {
		return fmt.Errorf("serviceCreateContributionCondolenceElement: %w", err)
	}
	if memorial == nil {
		return fmt.Errorf("serviceCreateContributionCondolenceElement: %w", errmgr.ErrMemorialNotFound)
	}

	// Check if the contributor exists
	contributor, err := d.params.Identity.FindUserByID(
		nil,           // db 		*gorm.DB
		TenantID,      // TenantID 	uint
		contributorID, // userID 	uint
	)
	if err != nil {
		return fmt.Errorf("serviceCreateContributionCondolenceElement: %w", err)
	}
	if contributor == nil {
		return fmt.Errorf("serviceCreateContributionCondolenceElement: %w", errmgr.ErrUserNotFound)
	}

	// Create the contribution story element
	newContributionCondolenceElement := schema.ContributionCondolenceElement{
		TenantID:          TenantID,
		MemorialID:        memorialID,
		ContributorID:     contributorID,
		ContributionState: schema.ContributionStatePending,
		IsImmutable:       isImmutable,

		ElementTitle:       elementTitle,
		ElementDescription: elementDescription,
		ElementAuthor:      elementAuthor,
		DesignElementID:    designElementID,
	}

	err = db.
		Create(&newContributionCondolenceElement).
		Error
	if err != nil {
		return fmt.Errorf("serviceCreateContributionCondolenceElement: %w", err)
	}

	return nil
}

func (d *Domain) serviceUpdateContributionGalleryElement(TenantID uint, memorialID uint, elementID uint, updaterID uint, isImmutable bool, elementTitle string, elementDescription string, elementDate time.Time, elementLocation *string, ElementGooglePlaceID *string) error {
	db := d.params.DB.GetDB()

	// Check if the contribution gallery element exists
	var existingElement schema.ContributionGalleryElement
	err := db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Where("id = ?", elementID).
		First(&existingElement).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("serviceUpdateContributionGalleryElement: %w", errmgr.ErrContributionElementNotFound)
		}
		return fmt.Errorf("serviceUpdateContributionGalleryElement: %w", err)
	}

	// if the element is not in pending state throw an error
	if existingElement.ContributionState != schema.ContributionStatePending {
		return fmt.Errorf("serviceUpdateContributionGalleryElement: %w", errmgr.ErrContributionElementNotPending)
	}

	// if the element is flagged as immutable the updater must be the contributor themself
	if existingElement.IsImmutable && existingElement.ContributorID != updaterID {
		return fmt.Errorf("serviceUpdateContributionStoryElement: %w", errmgr.ErrContributionElementImmutable)
	}

	// Update fields
	existingElement.IsImmutable = isImmutable
	existingElement.ElementTitle = elementTitle
	existingElement.ElementDescription = elementDescription
	existingElement.ElementDate = elementDate
	existingElement.ElementLocation = elementLocation
	existingElement.ElementGooglePlaceID = ElementGooglePlaceID

	err = db.
		Select(
			"IsImmutable",
			"ElementTitle",
			"ElementDescription",
			"ElementDate",
			"ElementLocation",
			"ElementGooglePlaceID",
		).
		Updates(&existingElement).
		Error
	if err != nil {
		return fmt.Errorf("serviceUpdateContributionGalleryElement: %w", err)
	}

	return nil
}
func (d *Domain) serviceUpdateContributionTimelineElement(TenantID uint, memorialID uint, elementID uint, updaterID uint, isImmutable bool, elementTitle string, elementDescription string, elementDate time.Time, elementEventType schema.EventTypeConst, elementLocation *string, ElementGooglePlaceID *string) error {
	db := d.params.DB.GetDB()

	// Check if the contribution timeline element exists
	var existingElement schema.ContributionTimelineElement
	err := db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Where("id = ?", elementID).
		First(&existingElement).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("serviceUpdateContributionTimelineElement: %w", errmgr.ErrContributionElementNotFound)
		}
		return fmt.Errorf("serviceUpdateContributionTimelineElement: %w", err)
	}

	// if the element is not in pending state throw an error
	if existingElement.ContributionState != schema.ContributionStatePending {
		return fmt.Errorf("serviceUpdateContributionGalleryElement: %w", errmgr.ErrContributionElementNotPending)
	}

	// if the element is flagged as immutable the updater must be the contributor themself
	if existingElement.IsImmutable && existingElement.ContributorID != updaterID {
		return fmt.Errorf("serviceUpdateContributionStoryElement: %w", errmgr.ErrContributionElementImmutable)
	}

	// Update fields
	existingElement.IsImmutable = isImmutable
	existingElement.ElementTitle = elementTitle
	existingElement.ElementDescription = elementDescription
	existingElement.ElementDate = elementDate
	existingElement.ElementEventType = elementEventType
	existingElement.ElementLocation = elementLocation
	existingElement.ElementGooglePlaceID = ElementGooglePlaceID

	err = db.
		Select(
			"IsImmutable",
			"ElementTitle",
			"ElementDescription",
			"ElementDate",
			"ElementEventType",
			"ElementLocation",
			"ElementGooglePlaceID",
		).
		Updates(&existingElement).
		Error
	if err != nil {
		return fmt.Errorf("serviceUpdateContributionTimelineElement: %w", err)
	}

	return nil
}
func (d *Domain) serviceUpdateContributionStoryElement(TenantID uint, memorialID uint, elementID uint, updaterID uint, isImmutable bool, elementTitle string, elementDescription string, elementAuthor string) error {
	db := d.params.DB.GetDB()

	// Check if the contribution story element exists
	var existingElement schema.ContributionStoryElement
	err := db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Where("id = ?", elementID).
		First(&existingElement).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("serviceUpdateContributionStoryElement: %w", errmgr.ErrContributionElementNotFound)
		}
		return fmt.Errorf("serviceUpdateContributionStoryElement: %w", err)
	}

	// if the element is not in pending state throw an error
	if existingElement.ContributionState != schema.ContributionStatePending {
		return fmt.Errorf("serviceUpdateContributionGalleryElement: %w", errmgr.ErrContributionElementNotPending)
	}

	// if the element is flagged as immutable the updater must be the contributor themself
	if existingElement.IsImmutable && existingElement.ContributorID != updaterID {
		return fmt.Errorf("serviceUpdateContributionStoryElement: %w", errmgr.ErrContributionElementImmutable)
	}

	// Update fields
	existingElement.IsImmutable = isImmutable
	existingElement.ElementTitle = elementTitle
	existingElement.ElementDescription = elementDescription
	existingElement.ElementAuthor = elementAuthor

	err = db.
		Select(
			"IsImmutable",
			"ElementTitle",
			"ElementDescription",
			"ElementAuthor",
		).
		Updates(&existingElement).
		Error
	if err != nil {
		return fmt.Errorf("serviceUpdateContributionStoryElement: %w", err)
	}

	return nil
}
func (d *Domain) serviceUpdateContributionCondolenceElement(TenantID uint, memorialID uint, elementID uint, updaterID uint, isImmutable bool, elementTitle string, elementDescription string, elementAuthor string, designElementID string) error {
	db := d.params.DB.GetDB()

	// Check if the contribution condolence element exists
	var existingElement schema.ContributionCondolenceElement
	err := db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Where("id = ?", elementID).
		First(&existingElement).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("serviceUpdateContributionCondolenceElement: %w", errmgr.ErrContributionElementNotFound)
		}
		return fmt.Errorf("serviceUpdateContributionCondolenceElement: %w", err)
	}

	// if the element is not in pending state throw an error
	if existingElement.ContributionState != schema.ContributionStatePending {
		return fmt.Errorf("serviceUpdateContributionGalleryElement: %w", errmgr.ErrContributionElementNotPending)
	}

	// if the element is flagged as immutable the updater must be the contributor themself
	if existingElement.IsImmutable && existingElement.ContributorID != updaterID {
		return fmt.Errorf("serviceUpdateContributionStoryElement: %w", errmgr.ErrContributionElementImmutable)
	}

	// Update fields
	existingElement.IsImmutable = isImmutable
	existingElement.ElementTitle = elementTitle
	existingElement.ElementDescription = elementDescription
	existingElement.ElementAuthor = elementAuthor
	existingElement.DesignElementID = designElementID

	err = db.
		Select(
			"IsImmutable",
			"ElementTitle",
			"ElementDescription",
			"ElementAuthor",
			"DesignElementID",
		).
		Updates(&existingElement).
		Error
	if err != nil {
		return fmt.Errorf("serviceUpdateContributionCondolenceElement: %w", err)
	}

	return nil
}

func (d *Domain) serviceDeleteContributionGalleryElement(TenantID uint, memorialID uint, elementID uint) error {
	var err error

	db := d.params.DB.GetDB()

	// Check if the memorial exists
	memorial, err := d.FindMemorialByID(memorialID)
	if err != nil {
		return fmt.Errorf("serviceCreateContributionGalleryElement: %w", err)
	}
	if memorial == nil {
		return fmt.Errorf("serviceCreateContributionGalleryElement: %w", errmgr.ErrMemorialNotFound)
	}

	// find and delete the gallery element
	err = db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Where("id = ?", elementID).
		Delete(&schema.ContributionGalleryElement{}).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("serviceCreateContributionGalleryElement: %w", errmgr.ErrContributionElementNotFound)
		}
		return fmt.Errorf("serviceCreateContributionGalleryElement: %w", err)
	}

	return nil
}
func (d *Domain) serviceDeleteContributionTimelineElement(TenantID uint, memorialID uint, elementID uint) error {
	var err error

	db := d.params.DB.GetDB()

	// Check if the memorial exists
	memorial, err := d.FindMemorialByID(memorialID)
	if err != nil {
		return fmt.Errorf("serviceCreateContributionTimelineElement: %w", err)
	}
	if memorial == nil {
		return fmt.Errorf("serviceCreateContributionTimelineElement: %w", errmgr.ErrMemorialNotFound)
	}

	// find and delete the timeline element
	err = db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Where("id = ?", elementID).
		Delete(&schema.ContributionTimelineElement{}).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("serviceCreateContributionTimelineElement: %w", errmgr.ErrContributionElementNotFound)
		}
		return fmt.Errorf("serviceCreateContributionTimelineElement: %w", err)
	}

	return nil
}
func (d *Domain) serviceDeleteContributionStoryElement(TenantID uint, memorialID uint, elementID uint) error {
	var err error

	db := d.params.DB.GetDB()

	// Check if the memorial exists
	memorial, err := d.FindMemorialByID(memorialID)
	if err != nil {
		return fmt.Errorf("serviceCreateContributionStoryElement: %w", err)
	}
	if memorial == nil {
		return fmt.Errorf("serviceCreateContributionStoryElement: %w", errmgr.ErrMemorialNotFound)
	}

	// find and delete the story element
	err = db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Where("id = ?", elementID).
		Delete(&schema.ContributionStoryElement{}).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("serviceCreateContributionStoryElement: %w", errmgr.ErrContributionElementNotFound)
		}
		return fmt.Errorf("serviceCreateContributionStoryElement: %w", err)
	}

	return nil
}
func (d *Domain) serviceDeleteContributionCondolenceElement(TenantID uint, memorialID uint, elementID uint) error {
	var err error

	db := d.params.DB.GetDB()

	// Check if the memorial exists
	memorial, err := d.FindMemorialByID(memorialID)
	if err != nil {
		return fmt.Errorf("serviceCreateContributionCondolenceElement: %w", err)
	}
	if memorial == nil {
		return fmt.Errorf("serviceCreateContributionCondolenceElement: %w", errmgr.ErrMemorialNotFound)
	}

	// find and delete the condolence element
	err = db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Where("id = ?", elementID).
		Delete(&schema.ContributionCondolenceElement{}).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("serviceCreateContributionCondolenceElement: %w", errmgr.ErrContributionElementNotFound)
		}
		return fmt.Errorf("serviceCreateContributionCondolenceElement: %w", err)
	}

	return nil
}

// helper
func (d *Domain) triggerExportPipeline(variables map[string]string) error {
	var err error

	url := "https://gitlab.com/api/v4/projects/60498659/trigger/pipeline"
	branch := "static"
	token := d.config.GitLabTriggerToken

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	err = writer.WriteField("token", token)
	if err != nil {
		return fmt.Errorf("triggerGitLabPipeline: %w", err)
	}
	err = writer.WriteField("ref", branch)
	if err != nil {
		return fmt.Errorf("triggerGitLabPipeline: %w", err)
	}

	for key, val := range variables {
		fieldName := fmt.Sprintf("variables[%s]", key)
		err := writer.WriteField(fieldName, val)
		if err != nil {
			return fmt.Errorf("triggerGitLabPipeline: unable to write variable %s field: %w", key, err)
		}
	}
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("triggerGitLabPipeline: %w", err)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("triggerGitLabPipeline: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("triggerGitLabPipeline: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		responseData, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("triggerGitLabPipeline: non-2xx status code %d: %s", resp.StatusCode, string(responseData))
	}

	d.logger.Debug("triggerGitLabPipeline: pipeline triggered successfully")
	return nil
}

// ! Finders -------------------------------------------------------

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

// Finds a memorial by its identifier, returns the memorial if it exists
// **Returns nil without an error** if the memorial does not exist
func (d *Domain) FindMemorialByIdentifier(TenantID uint, identifier string) (*schema.Memorial, error) {
	db := d.params.DB.GetDB()

	existingMemorial := schema.Memorial{}

	err := db.
		Where("fsp_id = ?", TenantID).
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

// Finds incomplete exports (requested or running) for a memorial.
// If incomplete exports are found, they are failed if stale, and the fresh ones are returned.
// **Returns nil without an error** if no incomplete exports are found
func (d *Domain) FindAndHandleIncompleteExports(db *gorm.DB, TenantID uint, memorialID uint) ([]schema.Export, error) {
	// if db is not set, use the default db
	// this allows the function to be used in a transaction
	if db == nil {
		db = d.params.DB.GetDB()
	}

	allIncompleteExports := []schema.Export{}
	err := db.
		Where("fsp_id = ?", TenantID).
		Where("memorial_id = ?", memorialID).
		Where("state = ? OR state = ?", schema.ExportStateRequested, schema.ExportStateRunning).
		Find(&allIncompleteExports).
		Error
	if err != nil {
		return nil, fmt.Errorf("FindAndHandleIncompleteExports: %w", err)
	}
	if len(allIncompleteExports) == 0 {
		return nil, nil
	}

	// Check for stale exports & fail them
	staleIncompleteExports := []schema.Export{}
	freshIncompleteExports := []schema.Export{}
	now := time.Now()

	for _, export := range allIncompleteExports {
		if export.UpdatedAt.Before(now.Add(-1 * time.Hour)) {
			export.State = schema.ExportStateFailed
			staleIncompleteExports = append(staleIncompleteExports, export)
		} else {
			freshIncompleteExports = append(freshIncompleteExports, export)
		}
	}
	if len(staleIncompleteExports) > 0 {
		for _, staleExport := range staleIncompleteExports {
			deferencedStaleExport := staleExport

			err = db.
				Model(&schema.Export{}).
				Where("fsp_id = ?", TenantID).
				Where("memorial_id = ?", memorialID).
				Where("id = ?", deferencedStaleExport.ID).
				Update("state", schema.ExportStateFailed).
				Error
			if err != nil {
				return nil, fmt.Errorf("FindAndHandleIncompleteExports: %w", err)
			}
		}
	}

	if len(freshIncompleteExports) == 0 {
		return nil, nil
	}

	return freshIncompleteExports, nil
}

// ! Helpers -------------------------------------------------------

func (d *Domain) signS3GetURL(url string, signThumbnail bool) (string, error) {
	if url == "" {
		return url, nil
	}

	if signThumbnail {
		convertedURL, err := d.params.S3.ConvertObjectURLOGToThumb(url)
		if err != nil {
			return "", fmt.Errorf("signS3GetURL: %w", err)
		}
		if convertedURL == nil {
			return "", fmt.Errorf("signS3GetURL: %s", errmgr.ErrNilCheckFailed)
		}
		url = *convertedURL
	}

	signedURL, err := d.params.S3.GetPresignedReadURL(
		context.Background(), // ctx	context.Context
		url,                  // s3URL	string
		24*time.Hour,         // exp	time.Duration
	)
	if err != nil {
		return "", fmt.Errorf("signS3GetURL: %w", err)
	}
	if signedURL == nil {
		return "", fmt.Errorf("signS3GetURL: %s", errmgr.ErrNilCheckFailed)
	}

	return *signedURL, nil
}
