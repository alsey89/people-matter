package errmgr

import (
	"errors"
	"net/http"
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
)

type StatusCodeMap map[ErrCode]int

var statusCodeMap = StatusCodeMap{
	ErrCodeInternal:     http.StatusInternalServerError,
	ErrCodeInvalidInput: http.StatusBadRequest,
	ErrCodeUnauthorized: http.StatusUnauthorized,
	ErrCodeForbidden:    http.StatusForbidden,
	ErrCodeTenant:       http.StatusBadRequest,
	ErrCodeToken:        http.StatusBadRequest,
	ErrCodeInput:        http.StatusBadRequest,

	ErrCodeUserNotFound:             http.StatusNotFound,
	ErrCodeInvalidCredentials:       http.StatusUnauthorized,
	ErrCodeEmailNotConfirmed:        http.StatusForbidden,
	ErrCodeEmailConfirmed:           http.StatusConflict,
	ErrCodeEmailInUse:               http.StatusConflict,
	ErrCodeNewPasswordIsOldPassword: http.StatusConflict,
}

// Type error messages for backend to match with errors.Is

var (
	// GENERAL ERRORS ----------------------------------------------------------

	ErrNilCheckFailed = errors.New("nil check failed")

	// DOMAIN ERRORS -----------------------------------------------------------

	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailNotConfirmed  = errors.New("email not confirmed")
	ErrEmailInUse         = errors.New("email already in use")
)

// Custom error type to return in api responses
type ErrResponse struct {
	Code    ErrCode `json:"code"`
	TraceID string  `json:"traceId"`
}

// GetStatusCode returns the http status code for a given ErrCode
// Falls back to 500 if the code is not found
func GetStatusCode(code ErrCode) int {
	status, ok := statusCodeMap[code]
	if !ok {
		return http.StatusInternalServerError
	}

	return status
}
