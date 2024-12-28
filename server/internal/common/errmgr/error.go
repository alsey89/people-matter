package errmgr

import (
	"errors"
)

// Coded errors for frontend to handle

type ErrCode string

const (

	// GENERAL ERRORS ----------------------------------------------------------

	ErrCodeInternal     ErrCode = "ERR_INTERNAL"      // 500: ALL server errors are obfuscated
	ErrCodeInvalidInput ErrCode = "ERR_INVALID_INPUT" // 400: non-specific bad request
	ErrCodeUnauthorized ErrCode = "ERR_UNAUTHORIZED"  // 401: non-specific unauthorized
	ErrCodeForbidden    ErrCode = "ERR_FORBIDDEN"     // 403: non-specific forbidden

	/*
		An error in the tenant identification process:
		1. missing tenant identifier
		2. tenant identifier cannot be resolved into a fspID
		3. resolved fspID does not match fspID in token
		4. missing fields in any part of the chain
		Note: this is considered 400 error because tenant identification chain starts from the client
	*/
	ErrCodeTenant ErrCode = "ERR_TENANT"

	/*
		An error in the token extraction and resolution process (*from context*):
		1. missing "user" (which respresents the token) in context
		2. user/token, claims, etc. is not valid type
		3. missing one more more claims fields
		Note: this is considered 400 error because it's technically client-side, even though the client has minimal control
	*/
	ErrCodeToken ErrCode = "ERR_TOKEN"

	/*
		An error in the input binding/validation process:
		1. binding failed
		2. validation failed
		3. missing or invalid path and query parameters
	*/
	ErrCodeInput ErrCode = "ERR_INPUT"

	// DOMAIN ERRORS -----------------------------------------------------------

	// identity.auth

	ErrCodeUserNotFound             ErrCode = "ERR_USER_NOT_FOUND"               // 404: user not found in database
	ErrCodeInvalidCredentials       ErrCode = "ERR_INVALID_CREDENTIALS"          // 401: credential match failed
	ErrCodeEmailNotConfirmed        ErrCode = "ERR_EMAIL_NOT_CONFIRMED"          // 403: user exists but email not confirmed
	ErrCodeEmailConfirmed           ErrCode = "ERR_EMAIL_CONFIRMED"              // 409: email is already confirmed
	ErrCodeEmailInUse               ErrCode = "ERR_EMAIL_IN_USE"                 // 409: email already in use
	ErrCodeNewPasswordIsOldPassword ErrCode = "ERR_NEW_PASSWORD_IS_OLD_PASSWORD" // 409: new password is not allowed to be the same as the old password

	// identity.role

	ErrCodeInvalidMemRole       ErrCode = "ERR_INVALID_MEMORIAL_ROLE"    // 403: invalid memorial role
	ErrCodeInvalidFSPRole       ErrCode = "ERR_INVALID_FSP_ROLE"         // 403: invalid fsp role
	ErrCodeMemorialRoleNotFound ErrCode = "ERR_MEMORIAL_ROLE_NOT_FOUND"  // 404: memorial role not found
	ErrCodeUserHasMemorialRole  ErrCode = "ERR_USER_HAS_MEMORIAL_ROLE"   // 409: user already has a role in this memorial
	ErrCodeFSPRoleNotFound      ErrCode = "ERR_FSP_ROLE_NOT_FOUND"       // 404: fsp role not found
	ErrCodeUserHasFSPRole       ErrCode = "ERR_USER_HAS_FSP_ROLE"        // 409: user already has a role in this fsp
	ErrCodeUserIsLastSuperAdmin ErrCode = "ERR_USER_IS_LAST_SUPER_ADMIN" // 409: user is the last superadmin in this fsp
	ErrCodeApplicationNotFound  ErrCode = "ERR_APPLICATION_NOT_FOUND"    // 404: application not found in db
	ErrCodeUserHasApplication   ErrCode = "ERR_USER_HAS_APPLICATION"     // 409: user already has an application in this memorial
	ErrCodeInvitationNotFound   ErrCode = "ERR_INVITATION_NOT_FOUND"     // 404: invitation not found in db
	ErrCodeUserHasInvitation    ErrCode = "ERR_USER_HAS_INVITATION"      // 409: user already has an invitation
	ErrCodeInvitationResponded  ErrCode = "ERR_INVITATION_RESPONDED"     // 409: invitation has already been responded to
	ErrCodeInvitationExpired    ErrCode = "ERR_INVITATION_EXPIRED"       // 403: invitation has expired
	ErrCodeInvitationNotForUser ErrCode = "ERR_INVITATION_NOT_FOR_USER"  // 403: invitation is not for this user
	ErrCodeUserIsCurator        ErrCode = "ERR_USER_IS_CURATOR"          // 409: curator cannot invite themselves

	// memorial

	ErrCodeMemorialNotFound              ErrCode = "ERR_MEMORIAL_NOT_FOUND"               // 404: memorial not found
	ErrCodeContributionElementNotFound   ErrCode = "ERR_CONTRIBUTION_ELEMENT_NOT_FOUND"   // 404: element not found, this error is common for all types of contribution elements
	ErrCodeContributionElementImmutable  ErrCode = "ERR_CONTRIBUTION_ELEMENT_IMMUTABLE"   // 403: contributor has flagged this element as immutable, only they can edit it
	ErrCodeContributionElementNotPending ErrCode = "ERR_CONTRIBUTION_ELEMENT_NOT_PENDING" // 403: contribution is not in the pending state, only pending contributions can be edited or approved

	ErrCodeExportNotFound   ErrCode = "ERR_EXPORT_NOT_FOUND"   // 404: export not found
	ErrCodeExportInProgress ErrCode = "ERR_EXPORT_IN_PROGRESS" // 409: an export is in progress
)

// Type error messages for backend to match with errors.Is

var (
	// GENERAL ERRORS ----------------------------------------------------------

	ErrNilCheckFailed = errors.New("nil check failed")

	// DOMAIN ERRORS -----------------------------------------------------------

	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailNotConfirmed  = errors.New("email not confirmed")
	ErrEmailInUse         = errors.New("email already in use")

	ErrMemorialRoleNotFound = errors.New("memorial role not found")
	ErrUserHasMemorialRole  = errors.New("user already has a role in this memorial")
	ErrFSPRoleNotFound      = errors.New("fsp role not found")
	ErrUserHasFSPRole       = errors.New("user already has a role in this fsp")
	ErrUserIsLastSuperAdmin = errors.New("user is the last superadmin in this fsp")

	ErrMemorialNotFound     = errors.New("memorial not found")
	ErrApplicationNotFound  = errors.New("application not found")
	ErrUserHasApplication   = errors.New("user already has an application in this memorial")
	ErrApplicationResponded = errors.New("application has already been responded to")
	ErrInvitationNotFound   = errors.New("invitation not found")
	ErrUserHasInvitation    = errors.New("user already has an invitation")
	ErrInvitationResponded  = errors.New("invitation has already been responded to")
	ErrInvitationExpired    = errors.New("invitation has expired")
	ErrInvitationNotForUser = errors.New("invitation is not for this user")
	ErrUserIsCurator        = errors.New("user is already the curator of this memorial")

	ErrContributionElementNotFound   = errors.New("contribution element not found")
	ErrContributionElementImmutable  = errors.New("contribution element is immutable")
	ErrContributionElementNotPending = errors.New("contribution is not in the pending state")

	ErrExportNotFound   = errors.New("export not found")
	ErrExportInProgress = errors.New("an export is in progress")
)

// Custom error type to return in api responses
type ErrResponse struct {
	Code    ErrCode `json:"code"`
	TraceID string  `json:"traceId"`
}
