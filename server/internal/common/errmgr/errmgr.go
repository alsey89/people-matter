package errmgr

import (
	"errors"
	"net/http"

	"go.uber.org/zap"
)

// Sentinel errors
var (
	ErrNilCheckFailed = errors.New("nil check failed")
	ErrTenant         = errors.New("tenant error")
	ErrPayload        = errors.New("payload binding or validation error")
	ErrPermission     = errors.New("invalid or insufficient permissions")

	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailUnverified    = errors.New("email not verified")
)

// Logs the error and returns an APIError that can be returned to the client.
func LogAndTranslateError(logger *zap.Logger, traceID string, err error) (string, int, APIError) {

	logger.Error("error occurred",
		zap.String("traceID", traceID),
		zap.Error(err),
	)

	switch {

	// ======================
	// GENERAL ERRORS
	// ======================

	case errors.Is(err, ErrTenant):
		return "Error in tenant chain",
			http.StatusNotFound,
			APIError{
				TraceID: traceID,
				Code:    "ERR_CODE_TENANT_CHAIN",
				Status:  http.StatusNotFound,
			}
	case errors.Is(err, ErrPayload):
		return "Payload binding or validation error",
			http.StatusBadRequest,
			APIError{
				TraceID: traceID,
				Code:    "ERR_CODE_PAYLOAD",
				Status:  http.StatusBadRequest,
			}
	case errors.Is(err, ErrPermission):
		return "Invalid or insufficient permissions",
			http.StatusForbidden,
			APIError{
				TraceID: traceID,
				Code:    "ERR_CODE_PERMISSION",
				Status:  http.StatusForbidden,
			}

	// ======================
	// IDENTITY DOMAIN ERRORS
	// ======================

	case errors.Is(err, ErrUserNotFound):
		return "User not found",
			http.StatusNotFound,
			APIError{
				TraceID: traceID,
				Code:    "ERR_CODE_USER_NOT_FOUND",
				Status:  http.StatusNotFound,
			}
	case errors.Is(err, ErrInvalidCredentials):
		return "Invalid credentials",
			http.StatusUnauthorized,
			APIError{
				TraceID: traceID,
				Code:    "ERR_CODE_INVALID_CREDENTIALS",
				Status:  http.StatusUnauthorized,
			}
	case errors.Is(err, ErrEmailUnverified):
		return "Email not verified",
			http.StatusUnauthorized,
			APIError{
				TraceID: traceID,
				Code:    "ERR_CODE_EMAIL_UNVERIFIED",
				Status:  http.StatusUnauthorized,
			}

	// ======================
	// DEFAULT FALLBACK
	// ======================

	default:
		return "Internal server error",
			http.StatusInternalServerError,
			APIError{
				TraceID: traceID,
				Code:    "ERR_CODE_INTERNAL",
				Status:  http.StatusInternalServerError,
			}

	}
}
