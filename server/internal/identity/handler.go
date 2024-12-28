package identity

import (
	"errors"
	"net/http"
	"time"

	"github.com/alsey89/people-matter/internal/common/errmgr"
	"github.com/alsey89/people-matter/internal/common/extractor"
	"github.com/alsey89/people-matter/internal/common/response"
	"github.com/alsey89/people-matter/internal/schema"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Forms ----------------------------------------------------------
func (d *Domain) handlerGetRelationships(c echo.Context) error {
	relationships := []schema.RelationshipConst{
		schema.RelationshipSelf,
		schema.RelationshipFriend,
		schema.RelationshipMother,
		schema.RelationshipFather,
		schema.RelationshipBrother,
		schema.RelationshipSister,
		schema.RelationshipSon,
		schema.RelationshipDaughter,
		schema.RelationshipAunt,
		schema.RelationshipUncle,
		schema.RelationshipGrandmother,
		schema.RelationshipGrandfather,
		schema.RelationshipGreatGrandmother,
		schema.RelationshipGreatGrandfather,
		schema.RelationshipNiece,
		schema.RelationshipNephew,
		schema.RelationshipCousin,
		schema.RelationshipSecondCousin,
		schema.RelationshipStepMother,
		schema.RelationshipStepFather,
		schema.RelationshipStepSon,
		schema.RelationshipStepdaughter,
	}

	return c.JSON(http.StatusOK, response.APIResponse{Message: "relationships", Data: relationships})
}

// Auth -----------------------------------------------------------
func (d *Domain) csrfHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, response.APIResponse{Message: "csrf token", Data: nil})
}

func (d *Domain) signinHandler(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("signinHandler:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if FSPID == nil {
		d.logger.Error("signinHandler:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	signinReq := new(signinRequest)

	err = c.Bind(signinReq)
	if err != nil {
		d.logger.Error("signinHandler:", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "payload binding/validation error",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInput,
				TraceID: traceID,
			},
		})
	}
	err = c.Validate(signinReq)
	if err != nil {
		d.logger.Error("signinHandler:", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "payload binding/validation error",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInput,
				TraceID: traceID,
			},
		})
	}

	user, rolesByLevel, err := d.signinService(*FSPID, signinReq.Email, signinReq.Password)
	switch {
	case err != nil:
		d.logger.Error("signinHandler:", zap.Error(err), zap.String("traceID", traceID))
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return c.JSON(http.StatusNotFound, response.APIResponse{
				Message: "user not found",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeUserNotFound,
					TraceID: traceID,
				},
			})
		case errors.Is(err, ErrEmailNotConfirmed):
			return c.JSON(http.StatusForbidden, response.APIResponse{
				Message: "user not confirmed",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeEmailNotConfirmed,
					TraceID: traceID,
				},
			})
		case errors.Is(err, ErrInvalidCredentials):
			return c.JSON(http.StatusUnauthorized, response.APIResponse{
				Message: "invalid credentials",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInvalidCredentials,
					TraceID: traceID,
				},
			})
		default:
			return c.JSON(http.StatusInternalServerError, response.APIResponse{
				Message: "something went wrong",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
		}
	}

	additionalClaims := map[string]interface{}{
		"userId": user.ID,
		"email":  user.Email,
	}

	additionalClaims["fspId"] = rolesByLevel.FSP.FSPID
	additionalClaims["fspRole"] = rolesByLevel.FSP.FSPRole.Name

	//todo: currently, the first memorial role is considered the active memorial upon sign in
	//todo: feat: pass active memorial by query param or other means
	var activeMemorialRole *schema.UserMemorialRole

	if len(rolesByLevel.Memorial) > 0 {
		activeMemorialRole = &rolesByLevel.Memorial[0]
		additionalClaims["activeMemID"] = activeMemorialRole.MemorialID
		additionalClaims["activeMemRole"] = activeMemorialRole.MemorialRole.Name
	}

	cookie, err := d.params.TokenManager.GenerateTokenAndHTTPonlyCookie(d.config.JWTAuthScope, additionalClaims)
	if err != nil {
		d.logger.Error("signinHandler:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(http.StatusInternalServerError, response.APIResponse{
			Message: "something went wrong",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInternal,
				TraceID: traceID,
			},
		})
	}
	if cookie == nil {
		d.logger.Error("signinHandler:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
		return c.JSON(http.StatusInternalServerError, response.APIResponse{
			Message: "something went wrong",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInternal,
				TraceID: traceID,
			},
		})
	}

	c.SetCookie(cookie)

	signingResData := signinResponseData{
		User:               user,
		RolesByLevel:       rolesByLevel,
		ActiveMemorialRole: activeMemorialRole,
	}

	return c.JSON(http.StatusOK, response.APIResponse{Message: "user signed in", Data: signingResData})
}

func (d *Domain) handlerSwitchActiveMemorial(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerSwitchActiveMemorial:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if FSPID == nil {
		d.logger.Error("handlerSwitchActiveMemorial:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	targetMemorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerSwitchActiveMemorial:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "invalid memorial ID",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInput,
				TraceID: traceID,
			},
		})
	}
	if targetMemorialID == nil {
		d.logger.Error("handlerSwitchActiveMemorial:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
		return c.JSON(http.StatusInternalServerError, response.APIResponse{
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
		d.logger.Error("handlerSwitchActiveMemorial:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "error in token chain",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeToken,
				TraceID: traceID,
			},
		})
	}
	userID, ok := claims["userId"].(float64)
	if !ok {
		d.logger.Error("handlerSwitchActiveMemorial: %w", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "error in token chain",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeToken,
				TraceID: traceID,
			},
		})
	}
	userEmail, ok := claims["email"].(string)
	if !ok {
		d.logger.Error("handlerSwitchActiveMemorial: %w", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "error in token chain",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeToken,
				TraceID: traceID,
			},
		})
	}

	rolesByLevel, err := d.serviceSwitchActiveMemorial(*FSPID, uint(userID), *targetMemorialID)
	if err != nil {
		d.logger.Error("handlerSwitchActiveMemorial:", zap.Error(err), zap.String("traceID", traceID))
		switch {
		case errors.Is(err, errmgr.ErrMemorialRoleNotFound):
			return c.JSON(http.StatusNotFound, response.APIResponse{
				Message: "memorial role not found",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeMemorialRoleNotFound,
					TraceID: traceID,
				},
			})
		default:
			return c.JSON(http.StatusInternalServerError, response.APIResponse{
				Message: "something went wrong",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
		}
	}

	var activeMemorialRole schema.UserMemorialRole

	//issue new token with the targetMemorialID as the activeMemorialID
	additionalClaims := map[string]interface{}{
		"userId":  userID,
		"email":   userEmail,
		"fspId":   rolesByLevel.FSP.FSPID,
		"fspRole": rolesByLevel.FSP.FSPRole.Name,
	}

	if len(rolesByLevel.Memorial) > 0 {
		for _, role := range rolesByLevel.Memorial {
			if role.MemorialID == *targetMemorialID {
				additionalClaims["activeMemID"] = role.MemorialID
				additionalClaims["activeMemRole"] = role.MemorialRole.Name
				activeMemorialRole = role
				break
			}
		}
	}

	cookie, err := d.params.TokenManager.GenerateTokenAndHTTPonlyCookie(d.config.JWTAuthScope, additionalClaims)
	if err != nil {
		d.logger.Error("handlerSwitchActiveMemorial:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(http.StatusInternalServerError, response.APIResponse{
			Message: "something went wrong",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInternal,
				TraceID: traceID,
			},
		})
	}
	if cookie == nil {
		d.logger.Error("handlerSwitchActiveMemorial:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
		return c.JSON(http.StatusInternalServerError, response.APIResponse{
			Message: "something went wrong",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInternal,
				TraceID: traceID,
			},
		})
	}

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, response.APIResponse{
		Message: "active memorial switched",
		Data: signinResponseData{
			User:               nil,
			RolesByLevel:       rolesByLevel,
			ActiveMemorialRole: &activeMemorialRole,
		},
	})
}

func (d *Domain) handlerApplicationSignup(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerApplicationSignup:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if FSPID == nil {
		d.logger.Error("handlerApplicationSignup:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	signupReq := new(applicationSignupRequest)
	err = c.Bind(signupReq)
	if err != nil {
		d.logger.Error("handlerApplicationSignup:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "payload binding/validation error",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInput,
				TraceID: traceID,
			},
		})
	}
	err = c.Validate(signupReq)
	if err != nil {
		d.logger.Error("handlerApplicationSignup:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "payload binding/validation error",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInput,
				TraceID: traceID,
			},
		})
	}

	_, err = d.applicationSignupService(
		*FSPID,                     // FSPID 			uint
		uint(signupReq.MemorialID), // memorialID 		uint
		signupReq.FirstName,        // firstName 		string
		signupReq.LastName,         // lastName 		string
		signupReq.Relationship,     // relationship 	schema.RelationshipConst
		signupReq.Email,            // email 			string
		signupReq.Password,         // password 		string
	)
	if err != nil {
		d.logger.Error("handlerApplicationSignup:", zap.Error(err), zap.String("traceID", traceID))
		switch {
		case errors.Is(err, ErrMemorialNotFound):
			return c.JSON(http.StatusNotFound, response.APIResponse{
				Message: "memorial not found",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeMemorialNotFound,
					TraceID: traceID,
				},
			})
		case errors.Is(err, ErrEmailAlreadyInUse):
			return c.JSON(http.StatusConflict, response.APIResponse{
				Message: "email already in use",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeEmailInUse,
					TraceID: traceID,
				},
			})
		case errors.Is(err, ErrUserHasApplication):
			return c.JSON(http.StatusConflict, response.APIResponse{
				Message: "user already has an application for this memorial",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeUserHasApplication,
					TraceID: traceID,
				},
			})
		default:
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
	}

	return c.JSON(http.StatusOK, response.APIResponse{
		Message: "application submitted",
		Data:    nil,
	})
}

func (d *Domain) handlerApplication(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerApplication:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if FSPID == nil {
		d.logger.Error("handlerApplication:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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
		d.logger.Error("handlerApplication:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "error in token chain",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeToken,
				TraceID: traceID,
			},
		})
	}
	userID, ok := claims["userId"].(float64)
	if !ok {
		d.logger.Error("handlerApplication: %w", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "error in token chain",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeToken,
				TraceID: traceID,
			},
		})
	}

	applicationReq := new(applicationRequest)
	err = c.Bind(applicationReq)
	if err != nil {
		d.logger.Error("handlerApplication:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "invalid request body",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInput,
				TraceID: traceID,
			},
		})
	}
	err = c.Validate(applicationReq)
	if err != nil {
		d.logger.Error("handlerApplication:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "missing or invalid fields",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInput,
				TraceID: traceID,
			},
		})
	}

	_, err = d.applicationService(
		*FSPID,                          // FSPID 			uint
		uint(applicationReq.MemorialID), // memorialID 		uint
		uint(userID),                    // applicantID 		uint
		applicationReq.Relationship,     // relationship 	schema.RelationshipConst
	)
	if err != nil {
		d.logger.Error("handlerApplication:", zap.Error(err), zap.String("traceID", traceID))
		if errors.Is(err, ErrMemorialNotFound) {
			return c.JSON(http.StatusNotFound, response.APIResponse{
				Message: "memorial not found",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeMemorialNotFound,
					TraceID: traceID,
				},
			})
		}
		if errors.Is(err, ErrUserHasMemorialRole) {
			return c.JSON(http.StatusConflict, response.APIResponse{
				Message: "already has role in memorial",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeUserHasMemorialRole,
					TraceID: traceID,
				},
			})
		}
		if errors.Is(err, ErrUserHasApplication) {
			return c.JSON(http.StatusConflict, response.APIResponse{
				Message: "already applied",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeUserHasApplication,
					TraceID: traceID,
				},
			})
		}
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "error signing up",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInternal,
				TraceID: traceID,
			},
		})
	}

	return c.JSON(http.StatusOK, response.APIResponse{
		Message: "application submitted",
		Data:    nil,
	})
}

func (d *Domain) handlerAcceptInvitationSignup(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerAcceptInvitationSignup:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if FSPID == nil {
		d.logger.Error("handlerAcceptInvitationSignup:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	signupReq := new(invitationSignupRequest)
	err = c.Bind(signupReq)
	if err != nil {
		d.logger.Error("handlerAcceptInvitationSignup:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "payload binding/validation error",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInput,
				TraceID: traceID,
			},
		})
	}
	err = c.Validate(signupReq)
	if err != nil {
		d.logger.Error("handlerAcceptInvitationSignup:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "payload binding/validation error",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInput,
				TraceID: traceID,
			},
		})
	}

	err = d.acceptInvitationSignupService(
		*FSPID,                     // FSPID 			uint
		uint(signupReq.MemorialID), // memorialID 		uint
		signupReq.Token,            // token 			string
		signupReq.Email,            // email 			string
		signupReq.FirstName,        // firstName 		string
		signupReq.LastName,         // lastName 		string
		signupReq.Password,         // password 		string
	)
	if err != nil {
		d.logger.Error("handlerAcceptInvitationSignup:", zap.Error(err), zap.String("traceID", traceID))
		switch {
		case errors.Is(err, errmgr.ErrInvitationNotFound):
			return c.JSON(http.StatusNotFound, response.APIResponse{
				Message: "invitation not found",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInvitationNotFound,
					TraceID: traceID,
				},
			})
		case errors.Is(err, errmgr.ErrInvitationExpired):
			return c.JSON(http.StatusUnauthorized, response.APIResponse{
				Message: "invitation expired",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInvitationNotFound,
					TraceID: traceID,
				},
			})
		case errors.Is(err, errmgr.ErrInvitationResponded):
			return c.JSON(http.StatusConflict, response.APIResponse{
				Message: "invitation already responded to",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInvitationResponded,
					TraceID: traceID,
				},
			})
		case errors.Is(err, errmgr.ErrEmailInUse):
			return c.JSON(http.StatusConflict, response.APIResponse{
				Message: "email already in use",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeEmailInUse,
					TraceID: traceID,
				},
			})
		case errors.Is(err, ErrUserHasMemorialRole):
			return c.JSON(http.StatusConflict, response.APIResponse{
				Message: "already has role in memorial",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeUserHasMemorialRole,
					TraceID: traceID,
				},
			})
		default:
			return c.JSON(http.StatusInternalServerError, response.APIResponse{
				Message: "something went wrong",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
		}
	}

	return c.JSON(http.StatusOK, response.APIResponse{
		Message: "user signed up",
		Data:    nil,
		Error: errmgr.ErrResponse{
			Code:    errmgr.ErrCodeInternal,
			TraceID: traceID,
		},
	})
}

func (d *Domain) handlerAcceptInvitation(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerAcceptInvitation:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if FSPID == nil {
		d.logger.Error("handlerAcceptInvitation:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	acceptInvitationReq := new(acceptInvitationRequest)

	err = c.Bind(acceptInvitationReq)
	if err != nil {
		d.logger.Error("handlerAcceptInvitation:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "payload binding/validation error",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInput,
				TraceID: traceID,
			},
		})
	}
	err = c.Validate(acceptInvitationReq)
	if err != nil {
		d.logger.Error("handlerAcceptInvitation:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "payload binding/validation error",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInput,
				TraceID: traceID,
			},
		})
	}

	_, claims, err := extractor.ExtractTokenAndClaimsFromContext(c)
	if err != nil {
		d.logger.Error("handlerAcceptInvitation:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "error in token chain",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeToken,
				TraceID: traceID,
			},
		})
	}
	userID, ok := claims["userId"].(float64)
	if !ok {
		d.logger.Error("handlerAcceptInvitation: %w", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "error in token chain",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeToken,
				TraceID: traceID,
			},
		})
	}

	err = d.acceptInvitationService(
		*FSPID,                               // FSPID 				uint
		uint(acceptInvitationReq.MemorialID), // memorialID 		uint
		uint(userID),                         // userID 			uint
		acceptInvitationReq.Token,            // token 				string
	)
	if err != nil {
		d.logger.Error("handlerAcceptInvitation:", zap.Error(err), zap.String("traceID", traceID))
		switch {
		//! UNAUTHORIZED ACCESS
		case errors.Is(err, errmgr.ErrInvitationNotForUser):
			d.logger.Warn("handlerAcceptInvitation: UNAUTHORIZED ACCESS DETECTED!",
				zap.String("method", c.Request().Method),
				zap.String("url", c.Request().URL.String()),
				zap.Any("headers", c.Request().Header),
				zap.Uint("userID", uint(userID)),
				zap.Uint("FSPID", *FSPID),
				zap.Uint("MemorialID", uint(acceptInvitationReq.MemorialID)),
				zap.Error(err),
				zap.String("traceID", traceID),
			)
			return c.JSON(http.StatusForbidden, response.APIResponse{
				Message: "invitation not for this user",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInvitationNotForUser,
					TraceID: traceID,
				},
			})
		case errors.Is(err, errmgr.ErrInvitationNotFound):
			return c.JSON(http.StatusNotFound, response.APIResponse{
				Message: "invitation not found",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInvitationNotFound,
					TraceID: traceID,
				},
			})
		case errors.Is(err, errmgr.ErrInvitationExpired):
			return c.JSON(http.StatusUnauthorized, response.APIResponse{
				Message: "invitation expired",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInvitationNotFound,
					TraceID: traceID,
				},
			})
		case errors.Is(err, errmgr.ErrInvitationResponded):
			return c.JSON(http.StatusConflict, response.APIResponse{
				Message: "invitation already responded to",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInvitationResponded,
					TraceID: traceID,
				},
			})
		case errors.Is(err, ErrUserHasMemorialRole):
			return c.JSON(http.StatusConflict, response.APIResponse{
				Message: "already has role in memorial",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeUserHasMemorialRole,
					TraceID: traceID,
				},
			})
		default:
			return c.JSON(http.StatusInternalServerError, response.APIResponse{
				Message: "something went wrong",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
		}
	}

	return c.JSON(http.StatusOK, response.APIResponse{
		Message: "invitation accepted",
		Data:    nil,
	})
}

func (d *Domain) confirmEmailHandler(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("confirmEmailHandler:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if FSPID == nil {
		d.logger.Error("confirmEmailHandler:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		d.logger.Error("confirmEmailHandler: user not found in context", zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest,
			response.APIResponse{
				Message: "token error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		d.logger.Error("confirmEmailHandler: claims not found in token", zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "token error",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeToken,
				TraceID: traceID,
			},
		})
	}

	email, ok := claims["email"].(string)
	if !ok {
		d.logger.Error("confirmEmailHandler: email not found in claims", zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest,
			response.APIResponse{
				Message: "token error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}

	err = d.confirmEmailService(*FSPID, email)
	switch {
	case err != nil:
		d.logger.Error("confirmEmailHandler:", zap.Error(err), zap.String("traceID", traceID))
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return c.JSON(http.StatusNotFound,
				response.APIResponse{
					Message: "user not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeUserNotFound,
						TraceID: traceID,
					},
				})
		case errors.Is(err, ErrEmailAlreadyConfirmed):
			return c.JSON(http.StatusConflict,
				response.APIResponse{
					Message: "already confirmed",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeEmailConfirmed,
						TraceID: traceID,
					},
				})
		default:
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
	}

	return c.JSON(http.StatusOK, response.APIResponse{
		Message: "email confirmed",
		Data:    nil,
	})
}

func (d *Domain) signoutHandler(c echo.Context) error {

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now(),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
		SameSite: http.SameSiteNoneMode,
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, response.APIResponse{
		Message: "user signed out",
		Data:    nil,
	})
}

func (d *Domain) requestResetPasswordHandler(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("requestResetPasswordHandler:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if FSPID == nil {
		d.logger.Error("requestResetPasswordHandler:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	resetPWReq := new(resetPasswordRequest)

	err = c.Bind(resetPWReq)
	if err != nil {
		d.logger.Error("requestResetPasswordHandler:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "payload binding/validation error",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInput,
				TraceID: traceID,
			},
		})
	}
	err = c.Validate(resetPWReq)
	if err != nil {
		d.logger.Error("requestResetPasswordHandler:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "payload binding/validation error",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInput,
				TraceID: traceID,
			},
		})
	}

	err = d.requestResetPasswordService(
		*FSPID,           // FSPID 			string
		resetPWReq.Email, // email 			string
	)
	if err != nil {
		d.logger.Error("requestResetPasswordHandler:", zap.Error(err), zap.String("traceID", traceID))
		switch {
		case errors.Is(err, ErrUserNotFound):
			return c.JSON(http.StatusNotFound, response.APIResponse{
				Message: "user not found",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeUserNotFound,
					TraceID: traceID,
				},
			})
		default:
			return c.JSON(http.StatusInternalServerError, response.APIResponse{
				Message: "something went wrong",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
		}
	}

	return c.JSON(http.StatusOK, response.APIResponse{
		Message: "password reset email sent",
		Data:    nil,
	})
}

func (d *Domain) handlerResetPassword(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("confirmResetPasswordHandler:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if FSPID == nil {
		d.logger.Error("confirmResetPasswordHandler:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	confirmResetPWReq := new(confirmResetPasswordRequest)

	err = c.Bind(confirmResetPWReq)
	if err != nil {
		d.logger.Error("confirmResetPasswordHandler:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "payload binding/validation error",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInput,
				TraceID: traceID,
			},
		})
	}
	err = c.Validate(confirmResetPWReq)
	if err != nil {
		d.logger.Error("confirmResetPasswordHandler:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "payload binding/validation error",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInput,
				TraceID: traceID,
			},
		})
	}

	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		d.logger.Error("confirmResetPasswordHandler: user not found in context", zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "token error",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeToken,
				TraceID: traceID,
			},
		})
	}
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		d.logger.Error("confirmResetPasswordHandler: claims not found in token", zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "token error",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeToken,
				TraceID: traceID,
			},
		})
	}
	email, ok := claims["email"].(string)
	if !ok {
		d.logger.Error("confirmResetPasswordHandler: email not found in claims", zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "token error",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeToken,
				TraceID: traceID,
			},
		})
	}

	err = d.updatePasswordServiceWithoutConfirmation(*FSPID, email, confirmResetPWReq.Password)
	if err != nil {
		d.logger.Error("confirmResetPasswordHandler:", zap.Error(err), zap.String("traceID", traceID))
		switch {
		case errors.Is(err, ErrNewPasswordIsOldPassword):
			return c.JSON(http.StatusConflict, response.APIResponse{
				Message: "new password cannot be the same as the old password",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeNewPasswordIsOldPassword,
					TraceID: traceID,
				},
			})
		default:
			return c.JSON(http.StatusInternalServerError, response.APIResponse{
				Message: "something went wrong",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
		}
	}

	return c.JSON(http.StatusOK, response.APIResponse{
		Message: "password reset confirmed",
		Data:    nil,
	})
}

func (d *Domain) handlerSetPasswordAndConfirmEmail(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("confirmResetPasswordHandler:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if FSPID == nil {
		d.logger.Error("confirmResetPasswordHandler:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	confirmResetPWReq := new(confirmResetPasswordRequest)

	err = c.Bind(confirmResetPWReq)
	if err != nil {
		d.logger.Error("confirmResetPasswordHandler:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "payload binding/validation error",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInput,
				TraceID: traceID,
			},
		})
	}
	err = c.Validate(confirmResetPWReq)
	if err != nil {
		d.logger.Error("confirmResetPasswordHandler:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "payload binding/validation error",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeInput,
				TraceID: traceID,
			},
		})
	}

	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		d.logger.Error("confirmResetPasswordHandler: user not found in context", zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "token error",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeToken,
				TraceID: traceID,
			},
		})
	}
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		d.logger.Error("confirmResetPasswordHandler: claims not found in token", zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "token error",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeToken,
				TraceID: traceID,
			},
		})
	}
	email, ok := claims["email"].(string)
	if !ok {
		d.logger.Error("confirmResetPasswordHandler: email not found in claims", zap.String("traceID", traceID))
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "token error",
			Data:    nil,
			Error: errmgr.ErrResponse{
				Code:    errmgr.ErrCodeToken,
				TraceID: traceID,
			},
		})
	}

	err = d.updatePasswordServiceWithConfirmation(*FSPID, email, confirmResetPWReq.Password, true)
	if err != nil {
		d.logger.Error("confirmResetPasswordHandler:", zap.Error(err), zap.String("traceID", traceID))
		switch {
		case errors.Is(err, ErrNewPasswordIsOldPassword):
			return c.JSON(http.StatusConflict, response.APIResponse{
				Message: "new password cannot be the same as the old password",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeNewPasswordIsOldPassword,
					TraceID: traceID,
				},
			})
		default:
			return c.JSON(http.StatusInternalServerError, response.APIResponse{
				Message: "something went wrong",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
		}
	}

	return c.JSON(http.StatusOK, response.APIResponse{
		Message: "password reset confirmed",
		Data:    nil,
	})
}

func (d *Domain) authCheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, response.APIResponse{
		Message: "user is authenticated",
		Data:    nil,
	})
}
