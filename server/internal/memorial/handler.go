package memorial

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alsey89/people-matter/internal/common/errmgr"
	"github.com/alsey89/people-matter/internal/common/extractor"
	"github.com/alsey89/people-matter/internal/common/response"
	"github.com/alsey89/people-matter/internal/schema"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Form Data ------------------------------------------------------
func (d *Domain) handlerGetMemorialFormRole(c echo.Context) error {
	memorialFormRoles := []schema.MemorialRoleConst{
		schema.RoleMemSelf,
		schema.RoleMemCurator,
		schema.RoleMemContributor,
		schema.RoleMemContributorApplicant,
		schema.RoleMemInsitu,
	}

	return c.JSON(http.StatusOK, response.APIResponse{Message: "form role", Data: memorialFormRoles})
}
func (d *Domain) handlerGetTimelineElementType(c echo.Context) error {
	eventTypes := []schema.EventTypeConst{
		schema.EventTypeBirth,
		schema.EventTypeDeath,
		schema.EventTypeMarriage,
		schema.EventTypeDivorce,
		schema.EventTypeEducation,
		schema.EventTypeGraduation,
		schema.EventTypeEmployment,
		schema.EventTypeRetirement,
		schema.EventTypePromotion,
		schema.EventTypeAward,
		schema.EventTypeAchievement,
		schema.EventTypeMilestone,
		schema.EventTypeOther,
	}

	return c.JSON(http.StatusOK, response.APIResponse{Message: "timeline element type", Data: eventTypes})
}

// CURATOR --------------------------------------------------------

func (d *Domain) handlerGetDashboard(c echo.Context) error {
	return c.JSON(http.StatusOK, response.APIResponse{Message: "dashboard", Data: nil})
}
func (d *Domain) handlerGetContributors(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerGetContributors:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerGetContributors:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusInternalServerError,
			response.APIResponse{
				Message: "error getting contributors",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
	}

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(
		c,            // c				echo.Context
		"memorialID", // key			string
	)
	if err != nil {
		d.logger.Error("handlerGetContributors:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusInternalServerError,
			response.APIResponse{
				Message: "missing or invalid path/query parameters",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
	}
	if memorialID == nil {
		d.logger.Error("handlerGetContributors:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error getting contributors",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}

	contributors, err := d.serviceGetContributors(*FSPID, *memorialID)
	if err != nil {
		d.logger.Error("handlerGetContributors: failed to get contributors from service", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusInternalServerError,
			response.APIResponse{
				Message: "error getting contributors",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
	}

	return c.JSON(http.StatusOK, response.APIResponse{
		Message: "contributors",
		Data:    contributors,
	})
}
func (d *Domain) handlerGetContributorApplications(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerGetContributorApplications:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerGetContributorApplications:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusInternalServerError,
			response.APIResponse{
				Message: "error getting contributor applicants",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
	}

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerGetContributorApplications:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerGetContributorApplications:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	applications, err := d.serviceGetContributorApplications(*FSPID, *memorialID)
	if err != nil {
		d.logger.Error("handlerGetContributorApplications: failed to get contributor applications from service", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusInternalServerError,
			response.APIResponse{
				Message: "error getting contributors",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
	}

	return c.JSON(http.StatusOK, response.APIResponse{Message: "contributor applications", Data: applications})
}
func (d *Domain) handlerAcceptContributorApplication(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerAcceptContributorApplication:", zap.Error(err), zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerAcceptContributorApplication:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "missing or invalid path/query parameters",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
	}

	applicationID, err := extractor.ExtractIDFromPathParamAsUINT(c, "applicationID")
	if err != nil {
		d.logger.Error("handlerAcceptContributorApplication:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "missing or invalid path/query parameters",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
	}

	err = d.serviceAcceptContributorApplication(
		*FSPID,         // FSPID			uint
		*memorialID,    // memorialID		uint
		*applicationID, // applicationID	uint
	)
	if err != nil {
		d.logger.Error("handlerAcceptContributorApplication: failed to accept contributor application", zap.Error(err), zap.String("traceID", traceID))
		switch {
		case errors.Is(err, errmgr.ErrMemorialRoleNotFound):
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "memorial not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeMemorialNotFound,
						TraceID: traceID,
					},
				})
		case errors.Is(err, errmgr.ErrApplicationNotFound):
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "application not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeApplicationNotFound,
						TraceID: traceID,
					},
				})
		case errors.Is(err, errmgr.ErrUserHasMemorialRole):
			return c.JSON(
				http.StatusConflict,
				response.APIResponse{
					Message: "user already has role in memorial",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeUserHasMemorialRole,
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "contributor accepted", Data: nil})
}
func (d *Domain) handlerRejectContributorApplication(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerRejectContributorApplication: failed to extract FSPID", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "erorr in the tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if FSPID == nil {
		d.logger.Error("handlerRejectContributorApplication:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusInternalServerError,
			response.APIResponse{
				Message: "error rejecting contributor",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
	}

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerRejectContributorApplication:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusInternalServerError,
			response.APIResponse{
				Message: "missing or invalid path/query parameters",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
	}

	applicationID, err := extractor.ExtractIDFromPathParamAsUINT(
		c,               // c				echo.Context
		"applicationID", // key				string
	)
	if err != nil {
		d.logger.Error("handlerRejectContributorApplication: failed to extract applicationID", zap.Error(err), zap.String("traceID", traceID))
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
	if applicationID == nil {
		d.logger.Error("handlerRejectContributorApplication:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	err = d.serviceRejectContributorApplication(
		*FSPID,         // FSPID			uint
		*memorialID,    // memorialID		uint
		*applicationID, // applicationID	uint
	)
	if err != nil {
		d.logger.Error("handlerRejectContributorApplication: failed to reject contributor application", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusInternalServerError,
			response.APIResponse{
				Message: "error rejecting contributor",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
	}

	return c.JSON(http.StatusOK, response.APIResponse{Message: "contributor rejected", Data: nil})
}
func (d *Domain) handlerGetContributorInvitations(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerGetContributorInvitations: failed to extract FSPID", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerGetContributorInvitations: FSPID is nil", zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerGetContributorInvitations: failed to extract memorialID", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerGetContributorInvitations: memorialID is nil", zap.String("traceID", traceID))
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

	invitations, err := d.serviceGetContributorInvitations(*FSPID, *memorialID)
	if err != nil {
		d.logger.Error("handlerGetContributorInvitations: failed to get contributor invitations from service", zap.Error(err), zap.String("traceID", traceID))
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "contributor invitees", Data: invitations})
}
func (d *Domain) handlerInviteContributor(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerInviteContributor:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerInviteContributor:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusInternalServerError,
			response.APIResponse{
				Message: "something went wrong",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerInviteContributor:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerInviteContributor:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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
		d.logger.Error("handlerInviteContributor:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}
	curatorID, ok := claims["userId"].(float64)
	if !ok {
		d.logger.Error("handlerInviteContributor: userId not found in claims", zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}
	UINTCuratorID := uint(curatorID)

	var postContributorReq postContributorRequest

	err = c.Bind(&postContributorReq)
	if err != nil {
		d.logger.Error("handlerInviteContributor:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}
	err = c.Validate(&postContributorReq)
	if err != nil {
		d.logger.Error("handlerInviteContributor:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}

	err = d.serviceInviteContributor(
		*FSPID,                          // FSPID			uint
		*memorialID,                     // memorialID		uint
		UINTCuratorID,                   // adminID		uint
		postContributorReq.Email,        // email			string
		postContributorReq.Relationship, // relationship	schema.RelationshipConst
	)
	if err != nil {
		d.logger.Error("handlerInviteContributor:", zap.Error(err), zap.String("traceID", traceID))
		switch {
		case errors.Is(err, errmgr.ErrUserHasMemorialRole):
			return c.JSON(
				http.StatusConflict, response.APIResponse{
					Message: "user already has role in memorial.",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeUserHasMemorialRole,
						TraceID: traceID,
					},
				})
		case errors.Is(err, errmgr.ErrUserIsCurator):
			return c.JSON(
				http.StatusConflict,
				response.APIResponse{
					Message: "curator cannot invite themselves",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeUserIsCurator,
						TraceID: traceID,
					},
				})
		case errors.Is(err, errmgr.ErrUserHasInvitation):
			return c.JSON(
				http.StatusConflict,
				response.APIResponse{
					Message: "user already has an invitation",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeUserHasInvitation,
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "invited contributor", Data: nil})
}
func (d *Domain) handlerDeleteContributorInvitation(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerDeleteContributorInvitation:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerDeleteContributorInvitation:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusInternalServerError,
			response.APIResponse{
				Message: "something went wrong",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerDeleteContributorInvitation:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerDeleteContributorInvitation:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusInternalServerError,
			response.APIResponse{
				Message: "something went wrong",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			},
		)
	}

	invitationID, err := extractor.ExtractIDFromPathParamAsUINT(c, "invitationID")
	if err != nil {
		d.logger.Error("handlerDeleteContributorInvitation: failed to extract invitationID", zap.Error(err), zap.String("traceID", traceID))
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
	if invitationID == nil {
		d.logger.Error("handlerDeleteContributorInvitation:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusInternalServerError,
			response.APIResponse{
				Message: "something went wrong",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			},
		)
	}

	err = d.serviceDeleteContributorInvitation(
		*FSPID,        // FSPID			uint
		*memorialID,   // memorialID		uint
		*invitationID, // invitationID	uint
	)
	if err != nil {
		d.logger.Error("handlerDeleteContributorInvitation: failed to delete contributor invitation", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusInternalServerError,
			response.APIResponse{
				Message: "error deleting contributor invitation",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
	}

	return c.JSON(http.StatusOK, response.APIResponse{Message: "deleted contributor invitation", Data: nil})
}
func (d *Domain) handlerReinviteContributor(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerReinviteContributor:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerReinviteContributor:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusInternalServerError,
			response.APIResponse{
				Message: "something went wrong",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerReinviteContributor:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerReinviteContributor:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	invitationID, err := extractor.ExtractIDFromPathParamAsUINT(c, "invitationID")
	if err != nil {
		d.logger.Error("handlerReinviteContributor:", zap.Error(err), zap.String("traceID", traceID))
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
	if invitationID == nil {
		d.logger.Error("handlerReinviteContributor:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	err = d.serviceReinviteContributor(
		*FSPID,        // FSPID			uint
		*memorialID,   // memorialID		uint
		*invitationID, // invitationID	uint
	)
	if err != nil {
		d.logger.Error("handlerReinviteContributor:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusInternalServerError,
			response.APIResponse{
				Message: "error reinviting contributor",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
	}

	return c.JSON(http.StatusOK, response.APIResponse{Message: "reinvite contributor", Data: nil})
}
func (d *Domain) handlerUpdateContributor(c echo.Context) error {
	return c.JSON(http.StatusOK, response.APIResponse{Message: "update contributor", Data: nil})
}
func (d *Domain) handlerDeleteContributor(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerDeleteContributor:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusInternalServerError,
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
		d.logger.Error("handlerDeleteContributor: FSPID is nil", zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerDeleteContributor: failed to extract memorialID", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerDeleteContributor:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	contributorMemorialRoleID, err := extractor.ExtractIDFromPathParamAsUINT(c, "contributorMemorialRoleID")
	if err != nil {
		d.logger.Error("handlerDeleteContributor:", zap.Error(err), zap.String("traceID", traceID))
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
	if contributorMemorialRoleID == nil {
		d.logger.Error("handlerDeleteContributor:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	userIsNotified, err := extractor.ExtractBoolFromQueryParam(
		c,                // c				echo.Context
		"userIsNotified", // key			string
		false,            // defaultValue	bool
	)
	if err != nil {
		d.logger.Error("handlerDeleteContributor: failed to extract userIsNotified", zap.Error(err), zap.String("traceID", traceID))
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
	if userIsNotified == nil {
		d.logger.Error("handlerDeleteContributor:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	err = d.serviceDeleteContributor(
		*FSPID,                     // FSPID			uint
		*memorialID,                // memorialID		uint
		*contributorMemorialRoleID, // contributorID	uint
		*userIsNotified,            // userIsNotified	bool
	)
	if err != nil {
		d.logger.Error("handlerDeleteContributor:", zap.Error(err), zap.String("traceID", traceID))
		switch {
		case errors.Is(err, errmgr.ErrUserNotFound):
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "user not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeUserNotFound,
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "deleted contributor", Data: nil})
}

func (d *Domain) handlerGetMemorialContributions(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerGetMemorialContributions:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerGetMemorialContributions: FSPID is nil", zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerGetMemorialContributions:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerGetMemorialContributions:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	contributions, err := d.serviceGetMemorialContributions(*FSPID, *memorialID)
	if err != nil {
		d.logger.Error("handlerGetMemorialContributions:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusInternalServerError,
			response.APIResponse{
				Message: "error getting memorial contributions",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
	}

	return c.JSON(http.StatusOK, response.APIResponse{Message: "memorial contributions", Data: contributions})
}
func (d *Domain) handlerUpdateContributionCondolenceElementState(c echo.Context) error {
	var err error

	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerUpdateContributionCondolenceElementState:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerUpdateContributionCondolenceElementState:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerUpdateContributionCondolenceElementState:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerUpdateContributionCondolenceElementState:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	elementID, err := extractor.ExtractIDFromPathParamAsUINT(c, "contributionCondolenceElementID")
	if err != nil {
		d.logger.Error("handlerUpdateContributionCondolenceElementState:", zap.Error(err), zap.String("traceID", traceID))
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
	if elementID == nil {
		d.logger.Error("handlerUpdateContributionCondolenceElementState:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	var putCondolenceElementStateReq putCondolenceElementStateRequest

	err = c.Bind(&putCondolenceElementStateReq)
	if err != nil {
		d.logger.Error("handlerUpdateContributionCondolenceElementState:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}
	err = c.Validate(&putCondolenceElementStateReq)
	if err != nil {
		d.logger.Error("handlerUpdateContributionCondolenceElementState:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}

	var contributionState schema.ContributionStateConst
	if putCondolenceElementStateReq.IsApproved {
		contributionState = schema.ContributionStatePrivate
	} else {
		contributionState = schema.ContributionStatePending
	}

	err = d.serviceUpdateContributionCondolenceElementState(
		*FSPID,            // FSPID					uint
		*memorialID,       // memorialID				uint
		*elementID,        // elementID				uint
		contributionState, // state					schema.ContributionState
	)
	if err != nil {
		d.logger.Error("handlerUpdateContributionCondolenceElementState:", zap.Error(err), zap.String("traceID", traceID))
		if errors.Is(err, errmgr.ErrContributionElementNotFound) {
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "element not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeContributionElementNotFound,
						TraceID: traceID,
					},
				})
		}
		return c.JSON(
			http.StatusInternalServerError,
			response.APIResponse{
				Message: "error updating contribution element state",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
	}

	return c.JSON(http.StatusOK, response.APIResponse{Message: "updated contribution element state", Data: nil})
}
func (d *Domain) handlerUpdateContributionGalleryElementState(c echo.Context) error {
	var err error

	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerApproveOrUnapproveContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerApproveOrUnapproveContributionGalleryElement:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerApproveOrUnapproveContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerApproveOrUnapproveContributionGalleryElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	elementID, err := extractor.ExtractIDFromPathParamAsUINT(c, "contributionGalleryElementID")
	if err != nil {
		d.logger.Error("handlerApproveOrUnapproveContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if elementID == nil {
		d.logger.Error("handlerApproveOrUnapproveContributionGalleryElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	var putGalleryElementStateReq putGalleryElementStateRequest

	err = c.Bind(&putGalleryElementStateReq)
	if err != nil {
		d.logger.Error("handlerApproveOrUnapproveContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}
	err = c.Validate(&putGalleryElementStateReq)
	if err != nil {
		d.logger.Error("handlerApproveOrUnapproveContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}

	var contributionState schema.ContributionStateConst
	if putGalleryElementStateReq.IsApproved {
		contributionState = schema.ContributionStatePrivate
	} else {
		contributionState = schema.ContributionStatePending
	}

	err = d.serviceUpdateContributionGalleryElementState(
		*FSPID,            // FSPID					uint
		*memorialID,       // memorialID				uint
		*elementID,        // elementID				uint
		contributionState, // state					schema.ContributionState
	)
	if err != nil {
		d.logger.Error("handlerApproveOrUnapproveContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
		if errors.Is(err, errmgr.ErrContributionElementNotFound) {
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "element not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeContributionElementNotFound,
						TraceID: traceID,
					},
				})
		}
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "updated contribution gallery element state", Data: nil})
}
func (d *Domain) handlerUpdateContributionStoryElementState(c echo.Context) error {
	var err error

	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerUpdateContributionStoryElementState:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerUpdateContributionStoryElementState: FSPID is nil", zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerUpdateContributionStoryElementState:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerUpdateContributionStoryElementState: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	elementID, err := extractor.ExtractIDFromPathParamAsUINT(c, "contributionStoryElementID")
	if err != nil {
		d.logger.Error("handlerUpdateContributionStoryElementState:", zap.Error(err), zap.String("traceID", traceID))
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
	if elementID == nil {
		d.logger.Error("handlerUpdateContributionStoryElementState: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	var putStoryElementStateReq putStoryElementStateRequest

	err = c.Bind(&putStoryElementStateReq)
	if err != nil {
		d.logger.Error("handlerUpdateContributionStoryElementState:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}
	err = c.Validate(&putStoryElementStateReq)
	if err != nil {
		d.logger.Error("handlerUpdateContributionStoryElementState:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}

	var contributionState schema.ContributionStateConst
	if putStoryElementStateReq.IsApproved {
		contributionState = schema.ContributionStatePrivate
	} else {
		contributionState = schema.ContributionStatePending
	}

	err = d.serviceUpdateContributionStoryElementState(
		*FSPID,            // FSPID					uint
		*memorialID,       // memorialID				uint
		*elementID,        // elementID				uint
		contributionState, // state					schema.ContributionState
	)
	if err != nil {
		d.logger.Error("handlerUpdateContributionStoryElementState:", zap.Error(err), zap.String("traceID", traceID))
		if errors.Is(err, errmgr.ErrContributionElementNotFound) {
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "element not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeContributionElementNotFound,
						TraceID: traceID,
					},
				})
		}
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "updated contribution story element state", Data: nil})
}
func (d *Domain) handlerUpdateContributionTimelineElementState(c echo.Context) error {
	var err error

	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerUpdateContributionTimelineElementState:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerUpdateContributionTimelineElementState: FSPID is nil", zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerUpdateContributionTimelineElementState:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerUpdateContributionTimelineElementState: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	elementID, err := extractor.ExtractIDFromPathParamAsUINT(c, "contributionTimelineElementID")
	if err != nil {
		d.logger.Error("handlerUpdateContributionTimelineElementState:", zap.Error(err), zap.String("traceID", traceID))
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
	if elementID == nil {
		d.logger.Error("handlerUpdateContributionTimelineElementState: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	var putTimelineElementStateReq putTimelineElementStateRequest

	err = c.Bind(&putTimelineElementStateReq)
	if err != nil {
		d.logger.Error("handlerUpdateContributionTimelineElementState:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}
	err = c.Validate(&putTimelineElementStateReq)
	if err != nil {
		d.logger.Error("handlerUpdateContributionTimelineElementState:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}

	var contributionState schema.ContributionStateConst
	if putTimelineElementStateReq.IsApproved {
		contributionState = schema.ContributionStatePrivate
	} else {
		contributionState = schema.ContributionStatePending
	}

	err = d.serviceUpdateContributionTimelineElementState(
		*FSPID,            // FSPID					uint
		*memorialID,       // memorialID				uint
		*elementID,        // elementID				uint
		contributionState, // state					schema.ContributionState
	)
	if err != nil {
		d.logger.Error("handlerUpdateContributionTimelineElementState:", zap.Error(err), zap.String("traceID", traceID))
		if errors.Is(err, errmgr.ErrContributionElementNotFound) {
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "element not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeContributionElementNotFound,
						TraceID: traceID,
					},
				})
		}
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "updated contribution timeline element state", Data: nil})
}

func (d *Domain) handlerExportMemorial(c echo.Context) error {
	var err error

	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerPublishMemorial:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerPublishMemorial: FSPID is nil", zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerPublishMemorial:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerPublishMemorial: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	createdExportID, err := d.serviceExportMemorial(
		*FSPID,      // FSPID			uint
		*memorialID, // memorialID		uint
	)
	if err != nil {
		d.logger.Error("handlerPublishMemorial:", zap.Error(err), zap.String("traceID", traceID))
		if errors.Is(err, errmgr.ErrMemorialNotFound) {
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "memorial not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeMemorialNotFound,
						TraceID: traceID,
					},
				})
		}
		if errors.Is(err, errmgr.ErrExportInProgress) {
			return c.JSON(
				http.StatusConflict,
				response.APIResponse{
					Message: "an export is already in progress",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeExportInProgress,
						TraceID: traceID,
					},
				})
		}
		return c.JSON(
			http.StatusInternalServerError,
			response.APIResponse{
				Message: "error publishing memorial",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
	}
	if createdExportID == nil {
		d.logger.Error("handlerPublishMemorial:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	returnPayload := map[string]uint{"exportId": *createdExportID}

	return c.JSON(http.StatusOK, response.APIResponse{Message: "published memorial", Data: returnPayload})
}
func (d *Domain) handlerGetExportState(c echo.Context) error {
	var err error

	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerGetExportState:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerGetExportState: FSPID is nil", zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerGetExportState:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerGetExportState: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	exportID, err := extractor.ExtractIDFromPathParamAsUINT(c, "exportID")
	if err != nil {
		d.logger.Error("handlerGetExportState:", zap.Error(err), zap.String("traceID", traceID))
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
	if exportID == nil {
		d.logger.Error("handlerGetExportState: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	exportState, err := d.serviceGetExportState(
		*FSPID,      // FSPID			uint
		*memorialID, // memorialID		uint
		*exportID,   // exportID		uint
	)
	if err != nil {
		d.logger.Error("handlerGetExportState:", zap.Error(err), zap.String("traceID", traceID))
		if errors.Is(err, errmgr.ErrExportNotFound) {
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "export not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeExportNotFound,
						TraceID: traceID,
					},
				})
		}
		return c.JSON(
			http.StatusInternalServerError,
			response.APIResponse{
				Message: "error getting export state",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
	}

	return c.JSON(http.StatusOK, response.APIResponse{Message: "export state", Data: exportState})
}
func (d *Domain) handlerGetAllExports(c echo.Context) error {
	var err error

	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerGetAllExports:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerGetAllExports: FSPID is nil", zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerGetAllExports:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerGetAllExports: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	exports, err := d.serviceGetAllExports(
		*FSPID,      // FSPID			uint
		*memorialID, // memorialID		uint
	)
	if err != nil {
		d.logger.Error("handlerGetAllExports:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusInternalServerError,
			response.APIResponse{
				Message: "error getting all exports",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
	}

	return c.JSON(http.StatusOK, response.APIResponse{Message: "all exports", Data: exports})
}
func (d *Domain) handlerUpdateExportState(c echo.Context) error {
	var err error

	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerUpdateExportState:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerUpdateExportState: FSPID is nil", zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerUpdateExportState:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerUpdateExportState: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	exportID, err := extractor.ExtractIDFromPathParamAsUINT(c, "exportID")
	if err != nil {
		d.logger.Error("handlerUpdateExportState:", zap.Error(err), zap.String("traceID", traceID))
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
	if exportID == nil {
		d.logger.Error("handlerUpdateExportState: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	var patchExportStateReq patchExportStateRequest

	err = c.Bind(&patchExportStateReq)
	if err != nil {
		d.logger.Error("handlerUpdateExportState:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}
	err = c.Validate(&patchExportStateReq)
	if err != nil {
		d.logger.Error("handlerUpdateExportState:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}

	//call service
	err = d.serviceUpdateExportState(
		*FSPID,                          // FSPID				uint
		*memorialID,                     // memorialID			uint
		*exportID,                       // exportID			uint
		patchExportStateReq.ExportState, // state				schema.ExportStateConst
	)
	if err != nil {
		d.logger.Error("handlerUpdateExportState:", zap.Error(err), zap.String("traceID", traceID))
		if errors.Is(err, errmgr.ErrExportNotFound) {
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "export not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeExportNotFound,
						TraceID: traceID,
					},
				})
		}
		return c.JSON(
			http.StatusInternalServerError,
			response.APIResponse{
				Message: "error updating export state",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
	}

	return c.JSON(http.StatusOK, response.APIResponse{Message: "update export state", Data: nil})
}

// CONTRIBUTOR ----------------------------------------------------

func (d *Domain) handlerGetContributions(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerGetContributions:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerGetContributions: FSPID is nil")
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerGetContributions:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "missing or invalid path/query parameters",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
	}
	if memorialID == nil {
		d.logger.Error("handlerGetContributions:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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
		d.logger.Error("handlerGetContributions:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}
	contributorID, ok := claims["userId"].(float64)
	if !ok {
		d.logger.Error("handlerGetContributions: userId not found in claims", zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}

	contributions, err := d.serviceGetContributions(
		*FSPID,              // FSPID			uint
		*memorialID,         // memorialID		uint
		uint(contributorID), // contributorID	uint
	)
	if err != nil {
		d.logger.Error("handlerGetContributions:", zap.Error(err), zap.String("traceID", traceID))
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "contributions", Data: contributions})
}
func (d *Domain) handlerGetPresignedUploadURL(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerGetPresignedUploadURL:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerGetPresignedUploadURL:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerInviteContributor:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "missing or invalid path/query parameters",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInternal,
					TraceID: traceID,
				},
			})
	}
	if memorialID == nil {
		d.logger.Error("handlerGetPresignedUploadURL:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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
		d.logger.Error("handlerGetPresignedUploadURL:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}
	contributorID, ok := claims["userId"].(float64)
	if !ok {
		d.logger.Error("handlerGetPresignedUploadURL: userId not found in claims", zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}

	presignedUploadURL, metadata, error := d.serviceGetPresignedUploadURL(
		*FSPID,              // FSPID			uint
		*memorialID,         // memorialID		uint
		uint(contributorID), // contributorID	uint
	)
	if error != nil {
		d.logger.Error("handlerGetPresignedUploadURL:", zap.Error(error), zap.String("traceID", traceID))
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
	if presignedUploadURL == nil || metadata == nil {
		d.logger.Error("handlerGetPresignedUploadURL:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	objectURL := strings.Split(*presignedUploadURL, "?")[0]

	s3PSUploadURLWithMetadata := map[string]interface{}{
		"presignedUploadUrl": presignedUploadURL,
		"objectUrl":          objectURL,
		"metadata":           metadata,
	}

	return c.JSON(http.StatusOK, response.APIResponse{Message: "presigned upload url", Data: s3PSUploadURLWithMetadata})
}

func (d *Domain) handlerCreateContributionGalleryElement(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerCreateContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerCreateContributionGalleryElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerCreateContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerCreateContributionGalleryElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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
		d.logger.Error("handlerCreateContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}
	userID, ok := claims["userId"].(float64)
	if !ok {
		d.logger.Error("handlerCreateContributionGalleryElement: userId not found in claims", zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}

	var postGalleryElementReq postGalleryElementRequest

	err = c.Bind(&postGalleryElementReq)
	if err != nil {
		d.logger.Error("handlerCreateContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))

		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}
	err = c.Validate(postGalleryElementReq)
	if err != nil {
		d.logger.Error("handlerCreateContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))

		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}

	//parse string date to time.Time
	FormattedElementDate, err := time.Parse("2006-01", postGalleryElementReq.ElementDate)
	if err != nil {
		d.logger.Error("handlerCreateContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "invalid date/time input",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}

	err = d.serviceCreateContributionGalleryElement(
		*FSPID,                            // FSPID					uint
		*memorialID,                       // memorialID				uint
		uint(userID),                      // contributorID			uint
		postGalleryElementReq.IsImmutable, // isImmutable			bool

		postGalleryElementReq.ElementTitle,         // elementTitle			string
		postGalleryElementReq.ElementDescription,   // elementDescription	string
		FormattedElementDate,                       // elementDate			string
		postGalleryElementReq.ElementMediaType,     // elementMediaType		string
		postGalleryElementReq.ElementMediaURL,      // elementURL			string
		postGalleryElementReq.ElementLocation,      // elementLocation		*string
		postGalleryElementReq.ElementGooglePlaceID, // elementGooglePlaceID	*string
	)
	if err != nil {
		d.logger.Error("handlerCreateContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "create contribution gallery element", Data: nil})
}
func (d *Domain) handlerCreateContributionTimelineElement(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerCreateContributionTimelineElement:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerCreateContributionTimelineElement:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerCreateContributionTimelineElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerCreateContributionTimelineElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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
		d.logger.Error("handlerCreateContributionTimelineElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}
	userID, ok := claims["userId"].(float64)
	if !ok {
		d.logger.Error("handlerCreateContributionTimelineElement: userId not found in claims", zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}

	var postTimelineElementReq postTimelineElementRequest

	err = c.Bind(&postTimelineElementReq)
	if err != nil {
		d.logger.Error("handlerCreateContributionTimelineElement:", zap.Error(err), zap.String("traceID", traceID))

		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}
	err = c.Validate(postTimelineElementReq)
	if err != nil {
		d.logger.Error("handlerCreateContributionTimelineElement:", zap.Error(err), zap.String("traceID", traceID))

		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}

	//parse string date to time.Time
	formattedElementDate, err := time.Parse("2006-01", postTimelineElementReq.ElementDate)
	if err != nil {
		d.logger.Error("handlerCreateContributionTimelineElement: failed to parse date", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "invalid date time input",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}

	err = d.serviceCreateContributionTimelineElement(
		*FSPID,                             // FSPID					uint
		*memorialID,                        // memorialID				uint
		uint(userID),                       // contributorID			uint
		postTimelineElementReq.IsImmutable, // isImmutable			bool

		postTimelineElementReq.ElementTitle,         // elementTitle				string
		postTimelineElementReq.ElementDescription,   // elementDescription		string
		formattedElementDate,                        // elementDate				string
		postTimelineElementReq.ElementEventType,     // elementEventType		schema.EventTypeConst
		postTimelineElementReq.ElementMediaURL,      // elementURL				*string
		postTimelineElementReq.ElementLocation,      // elementLocation			*string
		postTimelineElementReq.ElementGooglePlaceID, // elementGooglePlaceID	*string
	)
	if err != nil {
		d.logger.Error("handlerCreateContributionTimelineElement:", zap.Error(err), zap.String("traceID", traceID))
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "create contribution timeline element", Data: nil})
}
func (d *Domain) handlerCreateContributionStoryElement(c echo.Context) error {
	var err error

	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerCreateContributionStoryElement:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerCreateContributionStoryElement:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerCreateContributionStoryElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerCreateContributionStoryElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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
		d.logger.Error("handlerCreateContributionStoryElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}

	userID, ok := claims["userId"].(float64)
	if !ok {
		d.logger.Error("handlerCreateContributionTimelineElement: userId not found in claims", zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}

	var postStoryElementReq postStoryElementRequest

	err = c.Bind(&postStoryElementReq)
	if err != nil {
		d.logger.Error("handlerCreateContributionStoryElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}
	err = c.Validate(postStoryElementReq)
	if err != nil {
		d.logger.Error("handlerCreateContributionStoryElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}

	err = d.serviceCreateContributionStoryElement(
		*FSPID,                          // FSPID					uint
		*memorialID,                     // memorialID				uint
		uint(userID),                    // contributorID			uint
		postStoryElementReq.IsImmutable, // isImmutable			bool

		postStoryElementReq.ElementTitle,       // elementTitle				string
		postStoryElementReq.ElementDescription, // elementDescription		string
		postStoryElementReq.ElementMediaURL,    // elementURL				*string
		postStoryElementReq.ElementAuthor,      // elementAuthor			string
	)
	if err != nil {
		d.logger.Error("handlerCreateContributionStoryElement:", zap.Error(err), zap.String("traceID", traceID))
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "create contribution story element", Data: nil})
}
func (d *Domain) handlerCreateContributionCondolenceElement(c echo.Context) error {
	var err error

	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerCreateContributionCondolenceElement:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerCreateContributionCondolenceElement:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerCreateContributionCondolenceElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerCreateContributionCondolenceElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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
		d.logger.Error("handlerCreateContributionCondolenceElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}
	userID, ok := claims["userId"].(float64)
	if !ok {
		d.logger.Error("handlerCreateContributionCondolenceElement: userId not found in claims", zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}

	var postCondolenceElementReq postCondolenceElementRequest

	err = c.Bind(&postCondolenceElementReq)
	if err != nil {
		d.logger.Error("handlerCreateContributionCondolenceElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}
	err = c.Validate(postCondolenceElementReq)
	if err != nil {
		d.logger.Error("handlerCreateContributionCondolenceElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}

	err = d.serviceCreateContributionCondolenceElement(
		*FSPID,                                // FSPID								uint
		*memorialID,                           // memorialID						uint
		uint(userID),                          // contributorID				uint
		postCondolenceElementReq.IsImmutable,  // isImmutable					bool
		postCondolenceElementReq.ElementTitle, // elementTitle					string
		postCondolenceElementReq.ElementDescription, // elementDescription		string
		postCondolenceElementReq.ElementAuthor,      // elementAuthor			string
		postCondolenceElementReq.DesignElementID,    // designElementID			*uint
	)
	if err != nil {
		d.logger.Error("handlerCreateContributionCondolenceElement:", zap.Error(err), zap.String("traceID", traceID))
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "create contribution condolence element", Data: nil})
}

func (d *Domain) handlerUpdateContributionGalleryElement(c echo.Context) error {
	var err error

	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerUpdateContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerUpdateContributionGalleryElement:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerUpdateContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerUpdateContributionGalleryElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	elementID, err := extractor.ExtractIDFromPathParamAsUINT(c, "contributionGalleryElementID")
	if err != nil {
		d.logger.Error("handlerUpdateContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if elementID == nil {
		d.logger.Error("handlerUpdateContributionGalleryElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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
		d.logger.Error("handlerUpdateContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}

	userID, ok := claims["userId"].(float64)
	if !ok {
		d.logger.Error("handlerUpdateContributionGalleryElement: userId not found in claims", zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}

	var putGalleryElementReq putGalleryElementRequest

	err = c.Bind(&putGalleryElementReq)
	if err != nil {
		d.logger.Error("handlerUpdateContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}
	err = c.Validate(putGalleryElementReq)
	if err != nil {
		d.logger.Error("handlerUpdateContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}

	//parse string date to time.Time
	FormattedElementDate, err := time.Parse("2006-01", putGalleryElementReq.ElementDate)
	if err != nil {
		d.logger.Error("handlerUpdateContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "invalid date/time input",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}

	err = d.serviceUpdateContributionGalleryElement(
		*FSPID,                           // FSPID					uint
		*memorialID,                      // memorialID				uint
		*elementID,                       // elementID				uint
		uint(userID),                     // updatedID				uint
		putGalleryElementReq.IsImmutable, // isImmutable			bool

		putGalleryElementReq.ElementTitle,         // elementTitle				string
		putGalleryElementReq.ElementDescription,   // elementDescription		string
		FormattedElementDate,                      // elementDate				string
		putGalleryElementReq.ElementLocation,      // elementLocation			*string
		putGalleryElementReq.ElementGooglePlaceID, // elementGooglePlaceID		*string
	)
	if err != nil {
		d.logger.Error("handlerUpdateContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
		if errors.Is(err, errmgr.ErrContributionElementNotFound) {
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "element not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeContributionElementNotFound,
						TraceID: traceID,
					},
				})
		}
		if errors.Is(err, errmgr.ErrContributionElementImmutable) {
			return c.JSON(
				http.StatusForbidden,
				response.APIResponse{
					Message: "element is immutable",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeContributionElementImmutable,
						TraceID: traceID,
					},
				})
		}
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "updated contribution gallery element", Data: nil})
}
func (d *Domain) handlerUpdateContributionTimelineElement(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerUpdateContributionTimelineElement:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerUpdateContributionTimelineElement:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerUpdateContributionTimelineElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerUpdateContributionTimelineElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	elementID, err := extractor.ExtractIDFromPathParamAsUINT(c, "contributionTimelineElementID")
	if err != nil {
		d.logger.Error("handlerUpdateContributionTimelineElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if elementID == nil {
		d.logger.Error("handlerUpdateContributionTimelineElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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
		d.logger.Error("handlerUpdateContributionTimelineElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}

	userID, ok := claims["userId"].(float64)
	if !ok {
		d.logger.Error("handlerUpdateContributionTimelineElement: userId not found in claims", zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}

	var putTimelineElementReq putTimelineElementRequest

	err = c.Bind(&putTimelineElementReq)
	if err != nil {
		d.logger.Error("handlerUpdateContributionTimelineElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}
	err = c.Validate(putTimelineElementReq)
	if err != nil {
		d.logger.Error("handlerUpdateContributionTimelineElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}

	//parse string date to time.Time
	formattedElementDate, err := time.Parse("2006-01", putTimelineElementReq.ElementDate)
	if err != nil {
		d.logger.Error("handlerUpdateContributionTimelineElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "invalid date/time input",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}

	err = d.serviceUpdateContributionTimelineElement(
		*FSPID,                            // FSPID						uint
		*memorialID,                       // memorialID				uint
		*elementID,                        // elementID					uint
		uint(userID),                      // updaterID					uint
		putTimelineElementReq.IsImmutable, // isImmutable				bool

		putTimelineElementReq.ElementTitle,         // elementTitle				string
		putTimelineElementReq.ElementDescription,   // elementDescription		string
		formattedElementDate,                       // elementDate				string
		putTimelineElementReq.ElementEventType,     // elementEventType			schema.EventTypeConst
		putTimelineElementReq.ElementLocation,      // elementLocation			*string
		putTimelineElementReq.ElementGooglePlaceID, // elementGooglePlaceID		*string
	)
	if err != nil {
		d.logger.Error("handlerUpdateContributionTimelineElement:", zap.Error(err), zap.String("traceID", traceID))
		if errors.Is(err, errmgr.ErrContributionElementNotFound) {
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "element not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeContributionElementNotFound,
						TraceID: traceID,
					},
				})
		}
		if errors.Is(err, errmgr.ErrContributionElementImmutable) {
			return c.JSON(
				http.StatusForbidden,
				response.APIResponse{
					Message: "element is immutable",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeContributionElementImmutable,
						TraceID: traceID,
					},
				})
		}
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "updated contribution timeline element", Data: nil})
}
func (d *Domain) handlerUpdateContributionStoryElement(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerUpdateContributionStoryElement:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerUpdateContributionStoryElement:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerUpdateContributionStoryElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerUpdateContributionStoryElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	elementID, err := extractor.ExtractIDFromPathParamAsUINT(c, "contributionStoryElementID")
	if err != nil {
		d.logger.Error("handlerUpdateContributionStoryElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if elementID == nil {
		d.logger.Error("handlerUpdateContributionStoryElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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
		d.logger.Error("handlerUpdateContributionStoryElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}

	userID, ok := claims["userId"].(float64)
	if !ok {
		d.logger.Error("handlerUpdateContributionStoryElement: userId not found in claims", zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}

	var putStoryElementReq putStoryElementRequest

	err = c.Bind(&putStoryElementReq)
	if err != nil {
		d.logger.Error("handlerUpdateContributionStoryElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}
	err = c.Validate(putStoryElementReq)
	if err != nil {
		d.logger.Error("handlerUpdateContributionStoryElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}

	err = d.serviceUpdateContributionStoryElement(
		*FSPID,                         // FSPID					uint
		*memorialID,                    // memorialID				uint
		*elementID,                     // elementID				uint
		uint(userID),                   // updaterID				uint
		putStoryElementReq.IsImmutable, // isImmutable				bool

		putStoryElementReq.ElementTitle,       // elementTitle			string
		putStoryElementReq.ElementDescription, // elementDescription	string
		putStoryElementReq.ElementAuthor,      // elementAuthor			string
	)
	if err != nil {
		d.logger.Error("handlerUpdateContributionStoryElement:", zap.Error(err), zap.String("traceID", traceID))
		if errors.Is(err, errmgr.ErrContributionElementNotFound) {
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "element not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeContributionElementNotFound,
						TraceID: traceID,
					},
				})
		}
		if errors.Is(err, errmgr.ErrContributionElementImmutable) {
			return c.JSON(
				http.StatusForbidden,
				response.APIResponse{
					Message: "element is immutable",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeContributionElementImmutable,
						TraceID: traceID,
					},
				})
		}
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "updated contribution story element", Data: nil})
}
func (d *Domain) handlerUpdateContributionCondolenceElement(c echo.Context) error {
	var err error

	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerUpdateContributionCondolenceElement:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerUpdateContributionCondolenceElement:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerUpdateContributionCondolenceElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerUpdateContributionCondolenceElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	elementID, err := extractor.ExtractIDFromPathParamAsUINT(c, "contributionCondolenceElementID")
	if err != nil {
		d.logger.Error("handlerUpdateContributionCondolenceElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if elementID == nil {
		d.logger.Error("handlerUpdateContributionCondolenceElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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
		d.logger.Error("handlerUpdateContributionCondolenceElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}
	userID, ok := claims["userId"].(float64)
	if !ok {
		d.logger.Error("handlerUpdateContributionCondolenceElement: userId not found in claims", zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the token chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeToken,
					TraceID: traceID,
				},
			})
	}

	var putCondolenceElementReq putCondolenceElementRequest

	err = c.Bind(&putCondolenceElementReq)
	if err != nil {
		d.logger.Error("handlerUpdateContributionCondolenceElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}
	err = c.Validate(putCondolenceElementReq)
	if err != nil {
		d.logger.Error("handlerUpdateContributionCondolenceElement:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "payload binding/validation error",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}

	err = d.serviceUpdateContributionCondolenceElement(
		*FSPID,                               // FSPID								uint
		*memorialID,                          // memorialID						uint
		*elementID,                           // elementID						uint
		uint(userID),                         // updaterID						uint
		putCondolenceElementReq.IsImmutable,  // isImmutable					bool
		putCondolenceElementReq.ElementTitle, // elementTitle					string
		putCondolenceElementReq.ElementDescription, // elementDescription		string
		putCondolenceElementReq.ElementAuthor,      // elementAuthor			string
		putCondolenceElementReq.DesignElementID,    // designElementID			*uint
	)
	if err != nil {
		d.logger.Error("handlerUpdateContributionCondolenceElement:", zap.Error(err), zap.String("traceID", traceID))
		if errors.Is(err, errmgr.ErrContributionElementNotFound) {
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "element not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeContributionElementNotFound,
						TraceID: traceID,
					},
				})
		}
		if errors.Is(err, errmgr.ErrContributionElementImmutable) {
			return c.JSON(
				http.StatusForbidden,
				response.APIResponse{
					Message: "element is immutable",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeContributionElementImmutable,
						TraceID: traceID,
					},
				})
		}
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "updated contribution condolence element", Data: nil})
}

func (d *Domain) handlerDeleteContributionGalleryElement(c echo.Context) error {
	var err error

	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerDeleteContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerDeleteContributionGalleryElement:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerDeleteContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerDeleteContributionGalleryElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	galleryElementID, err := extractor.ExtractIDFromPathParamAsUINT(c, "contributionGalleryElementID")
	if err != nil {
		d.logger.Error("handlerDeleteContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if galleryElementID == nil {
		d.logger.Error("handlerDeleteContributionGalleryElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	err = d.serviceDeleteContributionGalleryElement(
		*FSPID,            // FSPID					uint
		*memorialID,       // memorialID				uint
		*galleryElementID, // contributionGalleryElementID	uint
	)
	if err != nil {
		d.logger.Error("handlerDeleteContributionGalleryElement:", zap.Error(err), zap.String("traceID", traceID))
		switch {
		case errors.Is(err, errmgr.ErrMemorialNotFound):
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "memorial not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeMemorialNotFound,
						TraceID: traceID,
					},
				})
		case errors.Is(err, errmgr.ErrContributionElementNotFound):
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "element not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeContributionElementNotFound,
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "deleted contribution gallery element", Data: nil})
}
func (d *Domain) handlerDeleteContributionTimelineElement(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerDeleteContributionTimelineElement:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerDeleteContributionTimelineElement:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerDeleteContributionTimelineElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerDeleteContributionTimelineElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	timelineElementID, err := extractor.ExtractIDFromPathParamAsUINT(c, "contributionTimelineElementID")
	if err != nil {
		d.logger.Error("handlerDeleteContributionTimelineElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if timelineElementID == nil {
		d.logger.Error("handlerDeleteContributionTimelineElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	err = d.serviceDeleteContributionTimelineElement(
		*FSPID,             // FSPID					uint
		*memorialID,        // memorialID				uint
		*timelineElementID, // contributionTimelineElementID	uint
	)
	if err != nil {
		d.logger.Error("handlerDeleteContributionTimelineElement:", zap.Error(err), zap.String("traceID", traceID))
		switch {
		case errors.Is(err, errmgr.ErrMemorialNotFound):
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "memorial not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeMemorialNotFound,
						TraceID: traceID,
					},
				})
		case errors.Is(err, errmgr.ErrContributionElementNotFound):
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "element not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeContributionElementNotFound,
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "deleted contribution timeline element", Data: nil})
}
func (d *Domain) handlerDeleteContributionStoryElement(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerDeleteContributionStoryElement:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerDeleteContributionStoryElement:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerDeleteContributionStoryElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerDeleteContributionStoryElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	storyElementID, err := extractor.ExtractIDFromPathParamAsUINT(c, "contributionStoryElementID")
	if err != nil {
		d.logger.Error("handlerDeleteContributionStoryElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if storyElementID == nil {
		d.logger.Error("handlerDeleteContributionStoryElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	err = d.serviceDeleteContributionStoryElement(
		*FSPID,          // FSPID					uint
		*memorialID,     // memorialID				uint
		*storyElementID, // contributionStoryElementID	uint
	)
	if err != nil {
		d.logger.Error("handlerDeleteContributionStoryElement:", zap.Error(err), zap.String("traceID", traceID))
		switch {
		case errors.Is(err, errmgr.ErrMemorialNotFound):
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "memorial not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeMemorialNotFound,
						TraceID: traceID,
					},
				})
		case errors.Is(err, errmgr.ErrContributionElementNotFound):
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "element not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeContributionElementNotFound,
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "deleted contribution story element", Data: nil})
}
func (d *Domain) handlerDeleteContributionCondolenceElement(c echo.Context) error {
	var err error

	traceID := uuid.NewString()

	FSPID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerDeleteContributionCondolenceElement:", zap.Error(err), zap.String("traceID", traceID))
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
		d.logger.Error("handlerDeleteContributionCondolenceElement:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	memorialID, err := extractor.ExtractIDFromPathParamAsUINT(c, "memorialID")
	if err != nil {
		d.logger.Error("handlerDeleteContributionCondolenceElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if memorialID == nil {
		d.logger.Error("handlerDeleteContributionCondolenceElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	condolenceElementID, err := extractor.ExtractIDFromPathParamAsUINT(c, "contributionCondolenceElementID")
	if err != nil {
		d.logger.Error("handlerDeleteContributionCondolenceElement:", zap.Error(err), zap.String("traceID", traceID))
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
	if condolenceElementID == nil {
		d.logger.Error("handlerDeleteContributionCondolenceElement: ", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	err = d.serviceDeleteContributionCondolenceElement(
		*FSPID,               // FSPID					uint
		*memorialID,          // memorialID				uint
		*condolenceElementID, // contributionCondolenceElementID	uint
	)
	if err != nil {
		d.logger.Error("handlerDeleteContributionCondolenceElement:", zap.Error(err), zap.String("traceID", traceID))
		switch {
		case errors.Is(err, errmgr.ErrMemorialNotFound):
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "memorial not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeMemorialNotFound,
						TraceID: traceID,
					},
				})
		case errors.Is(err, errmgr.ErrContributionElementNotFound):
			return c.JSON(
				http.StatusNotFound,
				response.APIResponse{
					Message: "element not found",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeContributionElementNotFound,
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "deleted contribution condolence element", Data: nil})
}

// USER ----------------------------------------------------------
func (d *Domain) handlerApplyToBeContributor(c echo.Context) error {

	return c.JSON(http.StatusOK, response.APIResponse{Message: "apply to be contributor", Data: nil})
}
func (d *Domain) handlerAcceptContributorInvitation(c echo.Context) error {

	return c.JSON(http.StatusOK, response.APIResponse{Message: "accept contributor invitation", Data: nil})
}
