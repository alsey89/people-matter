package identity

import (
	"net/http"

	"github.com/alsey89/people-matter/internal/common/constant"
	"github.com/alsey89/people-matter/internal/common/errmgr"
	"github.com/alsey89/people-matter/internal/common/extractor"
	"github.com/alsey89/people-matter/internal/common/response"
	"github.com/alsey89/people-matter/internal/common/util"
	"github.com/alsey89/people-matter/internal/schema"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Returns a middleware that checks if the JWT token is valid for the auth scope.
// Scope and scope-related settings are defined in config.
func (d *Domain) MustBeAuthenticated() echo.MiddlewareFunc {
	return d.params.TokenManager.GetJWTMiddleware("jwt_auth")
}

// Returns a middleware that checks if the JWT token is valid for the intended scope.
// Scope and scope-related settings are defined in config.
func (d *Domain) MustHaveValidConfirmationToken(tokenScope string) echo.MiddlewareFunc {
	return d.params.TokenManager.GetJWTMiddleware(tokenScope)
}

// Returns a middleware that resolves TenantID from the context TenantIdentifier.
// If the Tenant is not found in the database, it returns a 403.
func (d *Domain) MustResolveFSPID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var err error
			traceID := uuid.NewString()

			db := d.params.DB.GetDB()

			tenantIdentifier, ok := c.Get(constant.ContextTenantIdentifier).(string)
			if !ok || tenantIdentifier == "" {
				d.logger.Error("MustResolveFSPID:", zap.String("tenantIdentifier", tenantIdentifier), zap.String("traceID", traceID))
				return c.JSON(http.StatusForbidden, response.APIResponse{
					Message: "tenant error",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeTenant,
						TraceID: traceID,
					},
				})
			}

			var fsp schema.Tenant
			err = db.Model(&schema.Tenant{}).
				Where("tenant_identifier = ?", tenantIdentifier).
				First(&fsp).
				Error
			if err != nil {
				d.logger.Error("MustResolveFSPID:", zap.Error(err), zap.String("traceID", traceID))
				return c.JSON(http.StatusForbidden, response.APIResponse{
					Message: "tenant error",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeTenant,
						TraceID: traceID,
					},
				})
			}

			c.Set(constant.ContextTenantID, fsp.ID)

			return next(c)
		}
	}
}

// Returns a middleware that checks if the authenticated user's TenantID matches the resolved TenantID.
// If the FSPIDs do not match, it returns a 403.
func (d *Domain) MustMatchTenantIdentifierAndToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var err error
			traceID := uuid.NewString()

			contextFSPID, err := extractor.ExtractFSPIDFromContext(c)
			if err != nil {
				d.logger.Error("MustMatchTenantIdentifierAndToken:", zap.Error(err), zap.String("traceID", traceID))
				return c.JSON(http.StatusForbidden, response.APIResponse{
					Message: "error in tenant identification chain",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeTenant,
						TraceID: traceID,
					},
				})
			}

			_, claims, err := extractor.ExtractTokenAndClaimsFromContext(c)
			if err != nil {
				d.logger.Error("MustMatchTenantIdentifierAndToken:", zap.Error(err), zap.String("traceID", traceID))
				return c.JSON(http.StatusBadRequest, response.APIResponse{
					Message: "error in token chain",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeToken,
						TraceID: traceID,
					},
				})
			}
			tokenFSPID, ok := claims["fspId"].(float64) //*deserialized claims default to float64
			if !ok {
				d.logger.Error("MustMatchTenantIdentifierAndToken: error extracting fspId from token", zap.Any("claims", claims), zap.String("traceID", traceID))
				return c.JSON(http.StatusBadRequest, response.APIResponse{
					Message: "error in token chain",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeToken,
						TraceID: traceID,
					},
				})
			}

			if *contextFSPID != uint(tokenFSPID) {
				d.logger.Warn("MustMatchTenantIdentifierAndToken: tenant mismatch detected!", zap.Uint("contextFSPID", *contextFSPID), zap.Float64("tokenFSPID", tokenFSPID), zap.String("traceID", traceID))
				return c.JSON(http.StatusForbidden, response.APIResponse{
					Message: "error in tenant identification chain",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeInvalidFSPRole,
						TraceID: traceID,
					},
				})
			}

			return next(c)
		}
	}
}

// Returns a middleware that checks if the authenticated user's role is an Tenant superadmin.
// If the user's role is not Tenant SuperAdmin level, it returns a 403.
func (d *Domain) MustBeFSPSuperAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var err error
			traceID := uuid.NewString()

			_, claims, err := extractor.ExtractTokenAndClaimsFromContext(c)
			if err != nil {
				d.logger.Error("MustBeFSPSuperAdmin: %w", zap.Error(err), zap.String("traceID", traceID))
				return c.JSON(http.StatusBadRequest,
					response.APIResponse{
						Message: "error in token chain",
						Data:    nil,
						Error: errmgr.ErrResponse{
							Code:    errmgr.ErrCodeToken,
							TraceID: traceID,
						},
					})
			}
			fspRole, ok := claims["fspRole"].(string)
			if !ok {
				d.logger.Error("MustBeFSPSuperAdmin: error extracting fspRole from token", zap.String("traceID", traceID))
				return c.JSON(http.StatusForbidden,
					response.APIResponse{
						Message: "error in token chain",
						Data:    nil,
						Error: errmgr.ErrResponse{
							Code:    errmgr.ErrCodeToken,
							TraceID: traceID,
						},
					})
			}

			if fspRole != string(schema.RoleFSPSuperAdmin) {
				d.logger.Warn("MustBeFSPSuperAdmin:", zap.String("fspRole", fspRole), zap.String("traceID", traceID))
				return c.JSON(http.StatusForbidden,
					response.APIResponse{
						Message: "invalid fsp role",
						Data:    nil,
						Error: errmgr.ErrResponse{
							Code:    errmgr.ErrCodeInvalidFSPRole,
							TraceID: traceID,
						},
					})
			}

			return next(c)
		}
	}
}

// Returns a middleware that checks if the authenticated user's role is an Tenant admin or superadmin.
// If the user's role is not Tenant Admin or Tenant Super Admin, it returns a 403.
func (d *Domain) MustBeFSPAdminOrHigher() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var err error
			traceID := uuid.NewString()

			_, claims, err := extractor.ExtractTokenAndClaimsFromContext(c)
			if err != nil {
				d.logger.Error("MustBeFSPAdmin:", zap.Error(err), zap.String("traceID", traceID))
				return c.JSON(
					http.StatusBadRequest,
					response.APIResponse{
						Message: "error in token chain",
						Data:    nil,
						Error: errmgr.ErrResponse{
							Code:    errmgr.ErrCodeToken,
							TraceID: traceID,
						},
					})
			}
			fspRole, ok := claims["fspRole"].(string)
			if !ok {
				d.logger.Error("MustBeFSPAdmin: error extracting fspRole from token", zap.String("traceID", traceID))
				return c.JSON(http.StatusForbidden,
					response.APIResponse{
						Message: "error in token chain",
						Data:    nil,
						Error: errmgr.ErrResponse{
							Code:    errmgr.ErrCodeToken,
							TraceID: traceID,
						},
					})
			}

			if fspRole != string(schema.RoleFSPAdmin) && fspRole != string(schema.RoleFSPSuperAdmin) {
				d.logger.Warn("MustBeFSPAdmin:", zap.String("fspRole", fspRole))
				return c.JSON(http.StatusForbidden,
					response.APIResponse{
						Message: "invalid fsp role",
						Data:    nil,
						Error: errmgr.ErrResponse{
							Code:    errmgr.ErrCodeInvalidFSPRole,
							TraceID: traceID,
						},
					})
			}

			return next(c)
		}
	}
}

// Returns a middleware that checks if the authenticated user is curator of current memorial.
// If the user's role is not curator of the current memorial, it returns a 403.
func (d *Domain) MustBeCuratorOfCurrentMemorial() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var err error
			traceID := uuid.NewString()

			memorialId, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
			if err != nil {
				d.logger.Error("MustBeCuratorOfCurrentMemorial:", zap.Error(err), zap.String("traceID", traceID))
				return c.JSON(
					http.StatusBadRequest,
					response.APIResponse{
						Message: "missing or invalid path/query parameters",
						Data:    nil,
						Error: errmgr.ErrResponse{
							Code:    errmgr.ErrCodeInput,
							TraceID: traceID,
						},
					})
			}
			if memorialId == nil {
				d.logger.Error("MustBeCuratorOfCurrentMemorial:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
				return c.JSON(
					http.StatusInternalServerError,
					response.APIResponse{
						Message: "something went wrong",
						Data:    nil,
						Error: errmgr.ErrResponse{
							Code:    errmgr.ErrCodeInternal,
							TraceID: traceID,
						},
					})
			}

			_, claims, err := extractor.ExtractTokenAndClaimsFromContext(c)
			if err != nil {
				d.logger.Error("MustBeCuratorOfCurrentMemorial:", zap.Error(err), zap.String("traceID", traceID))
				return c.JSON(
					http.StatusBadRequest,
					response.APIResponse{
						Message: "error in token chain",
						Data:    nil,
						Error: errmgr.ErrResponse{
							Code:    errmgr.ErrCodeToken,
							TraceID: traceID,
						},
					})
			}

			activeMemorialID, ok := claims["activeMemID"].(float64) //*deserialized claims default to float64
			if !ok {
				d.logger.Error("MustBeCuratorOfCurrentMemorial: error extracting activeMemID from claims", zap.Any("claims", claims), zap.String("traceID", traceID))
				return c.JSON(
					http.StatusBadRequest,
					response.APIResponse{
						Message: "error in token chain",
						Data:    nil,
						Error: errmgr.ErrResponse{
							Code:    errmgr.ErrCodeToken,
							TraceID: traceID,
						},
					})
			}

			UINTActiveMemorialID, err := util.Float64ToUINT(activeMemorialID)
			if err != nil {
				d.logger.Error("MustBeCuratorOfCurrentMemorial:", zap.Error(err), zap.String("traceID", traceID))
				return c.JSON(
					http.StatusInternalServerError,
					response.APIResponse{
						Message: "something went wrong",
						Data:    nil,
						Error: errmgr.ErrResponse{
							Code:    errmgr.ErrCodeInternal,
							TraceID: traceID,
						},
					})
			}
			if UINTActiveMemorialID == nil {
				d.logger.Error("MustBeCuratorOfCurrentMemorial:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
				return c.JSON(
					http.StatusInternalServerError,
					response.APIResponse{
						Message: "something went wrong",
						Data:    nil,
						Error: errmgr.ErrResponse{
							Code:    errmgr.ErrCodeInternal,
							TraceID: traceID,
						},
					})
			}

			// extract userID from token
			userID, ok := claims["userId"].(float64) //*deserialized claims default to float64
			if !ok {
				d.logger.Error("MustBeCuratorOfCurrentMemorial: error extracting userID from claims", zap.Any("claims", claims), zap.String("traceID", traceID))
				return c.JSON(http.StatusForbidden, response.APIResponse{
					Message: "MemorialRole error",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeToken,
						TraceID: traceID,
					},
				})
			}

			//! reject if the memorialID in the path param does not match the activeMemID in the token
			if *memorialId != *UINTActiveMemorialID {
				d.logger.Warn("MustBeCuratorOfCurrentMemorial: Unauthorized Memorial Access",
					zap.Uint("memorialID", *memorialId),
					zap.Uint("expectedMemorialID", *UINTActiveMemorialID),
					zap.Float64("userID", userID),
					zap.String("traceID", traceID),
				)
				return c.JSON(http.StatusForbidden, response.APIResponse{
					Message: "invalid memorial role",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeInvalidMemRole,
						TraceID: traceID,
					},
				})
			}

			activeMemorialRole, ok := claims["activeMemRole"].(string)
			if !ok {
				d.logger.Error("MustBeCuratorOfCurrentMemorial: error extracting activeMemRole from claims", zap.Any("claims", claims), zap.String("traceID", traceID))
				return c.JSON(
					http.StatusBadRequest,
					response.APIResponse{
						Message: "error in token chain",
						Data:    nil,
						Error: errmgr.ErrResponse{
							Code:    errmgr.ErrCodeToken,
							TraceID: traceID,
						},
					})
			}

			//! reject if the activeMemRole in the token is not "curator"
			if activeMemorialRole != string(schema.RoleMemCurator) {
				d.logger.Warn("MustBeCuratorOfCurrentMemorial: Unauthorized Role Access",
					zap.String("memorialRole", activeMemorialRole),
					zap.String("expectedMemorialRole", string(schema.RoleMemCurator)),
					zap.Uint("memorialID", *memorialId),
					zap.Float64("userID", userID),
					zap.String("traceID", traceID),
				)
				return c.JSON(http.StatusForbidden, response.APIResponse{
					Message: "invalid memorial role",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeInvalidMemRole,
						TraceID: traceID,
					},
				})
			}

			return next(c)
		}
	}
}

// Returns a middleware that checks if the authenticated user is a contributor or curator at the current memorial.
// If the user's role is not contributor or curator at the current memorial, it returns a 403.
func (d *Domain) MustBeContributorOrCuratorOfCurrentMemorial() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var err error
			traceID := uuid.NewString()

			//extract memorialID from path param
			memorialId, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
			if err != nil {
				d.logger.Error("MustBeContributorOrCuratorOfCurrentMemorial:", zap.Error(err), zap.String("traceID", traceID))
				return c.JSON(http.StatusBadRequest, response.APIResponse{
					Message: "missing or invalid path/query parameters",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeInput,
						TraceID: traceID,
					},
				})
			}
			if memorialId == nil {
				d.logger.Error("MustBeContributorOrCuratorOfCurrentMemorial:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
				return c.JSON(http.StatusInternalServerError, response.APIResponse{
					Message: "something went wrong",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeInternal,
						TraceID: traceID,
					},
				})
			}

			//extract claims from token
			_, claims, err := extractor.ExtractTokenAndClaimsFromContext(c)
			if err != nil {
				d.logger.Error("MustBeContributorOrCuratorOfCurrentMemorial:", zap.Error(err), zap.String("traceID", traceID))
				return c.JSON(http.StatusBadRequest, response.APIResponse{
					Message: "error in token chain",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeToken,
						TraceID: traceID,
					},
				})
			}

			activeMemorialID, ok := claims["activeMemID"].(float64) //*deserialized claims default to float64
			if !ok {
				d.logger.Error("MustBeContributorOrCuratorOfCurrentMemorial: error extracting activeMemID from claims", zap.Any("claims", claims), zap.String("traceID", traceID))
				return c.JSON(http.StatusBadRequest, response.APIResponse{
					Message: "error in token chain",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeToken,
						TraceID: traceID,
					},
				})
			}

			UINTActiveMemorialID, err := util.Float64ToUINT(activeMemorialID)
			if err != nil {
				d.logger.Error("MustBeContributorOrCuratorOfCurrentMemorial:", zap.Error(err), zap.String("traceID", traceID))
				return c.JSON(http.StatusInternalServerError, response.APIResponse{
					Message: "something went wrong",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeInternal,
						TraceID: traceID,
					},
				})
			}
			if UINTActiveMemorialID == nil {
				d.logger.Error("MustBeContributorOrCuratorOfCurrentMemorial:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
				return c.JSON(http.StatusInternalServerError, response.APIResponse{
					Message: "something went wrong",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeInternal,
						TraceID: traceID,
					},
				})
			}

			// extract userID from token
			userID, ok := claims["userId"].(float64) //*deserialized claims default to float64
			if !ok {
				d.logger.Error("MustBeContributorOrCuratorOfCurrentMemorial: error extracting userID from claims", zap.Any("claims", claims), zap.String("traceID", traceID))
				return c.JSON(
					http.StatusBadRequest,
					response.APIResponse{
						Message: "error in token chain",
						Data:    nil,
						Error: errmgr.ErrResponse{
							Code:    errmgr.ErrCodeToken,
							TraceID: traceID,
						},
					})
			}

			//! reject if the memorialID in the path param does not match the activeMemID in the token
			if *memorialId != *UINTActiveMemorialID {
				d.logger.Warn("MustBeContributorOrCuratorOfCurrentMemorial: Unauthorized Memorial Access",
					zap.Uint("memorialID", *memorialId),
					zap.Uint("expectedMemorialID", *UINTActiveMemorialID),
					zap.Float64("userID", userID),
					zap.String("traceID", traceID),
				)
				return c.JSON(http.StatusForbidden, response.APIResponse{
					Message: "invalid memorial role",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeInvalidMemRole,
						TraceID: traceID,
					},
				})
			}

			activeMemorialRole, ok := claims["activeMemRole"].(string)
			if !ok {
				d.logger.Error("MustBeContributorOrCuratorOfCurrentMemorial: error extracting activeMemRole from claims", zap.Any("claims", claims), zap.String("traceID", traceID))
				return c.JSON(http.StatusBadRequest, response.APIResponse{
					Message: "error in token chain",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeToken,
						TraceID: traceID,
					},
				})
			}

			//! reject if the activeMemRole in the token is not "contributor" or "curator"
			if activeMemorialRole != string(schema.RoleMemContributor) && activeMemorialRole != string(schema.RoleMemCurator) {
				d.logger.Warn("MustBeCuratorOfCurrentMemorial: Unauthorized Role Access",
					zap.String("memorialRole", activeMemorialRole),
					zap.String("expectedMemorialRole", string(schema.RoleMemContributor)),
					zap.Uint("memorialID", *memorialId),
					zap.Float64("userID", userID),
					zap.String("traceID", traceID),
				)
				return c.JSON(
					http.StatusForbidden,
					response.APIResponse{
						Message: "invalid memorial role",
						Data:    nil,
						Error: errmgr.ErrResponse{
							Code:    errmgr.ErrCodeInvalidMemRole,
							TraceID: traceID,
						},
					})
			}

			return next(c)
		}
	}
}
