package fsp

import (
	"errors"
	"net/http"
	"time"

	"github.com/alsey89/people-matter/internal/common/errmgr"
	"github.com/alsey89/people-matter/internal/common/extractor"
	"github.com/alsey89/people-matter/internal/common/response"
	"github.com/alsey89/people-matter/internal/schema"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Data for forms --------------------------------------------

// Dashboard

func (d *Domain) handlerGetDashboard(c echo.Context) error {
	return c.JSON(http.StatusOK, response.APIResponse{Message: "dashboard data", Data: nil})
}

// Account --------------------------------------------

func (d *Domain) handlerGetFSPAccount(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	TenantID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerGetFSPAccount:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if TenantID == nil {
		d.logger.Error("handlerGetFSPAccount:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	exstingFSP, err := d.GetAccountService(
		*TenantID, //TenantID 			uint
		true,      //preloadDetails	bool
	)
	if err != nil {
		d.logger.Error("handlerGetFSPAccount:", zap.Error(err), zap.String("traceID", traceID))
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "account", Data: exstingFSP})
}
func (d *Domain) handlerUpdateFSPAccount(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	TenantID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerUpdateFSPAccount:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if TenantID == nil {
		d.logger.Error("handlerUpdateFSPAccount:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	putMemorialReq := schema.Tenant{}

	err = c.Bind(&putMemorialReq)
	if err != nil {
		d.logger.Error("handlerUpdateFSPAccount:", zap.Error(err), zap.String("traceID", traceID))
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

	err = d.updateAccountService(*TenantID, putMemorialReq)
	if err != nil {
		d.logger.Error("handlerUpdateFSPAccount:", zap.Error(err), zap.String("traceID", traceID))
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "account updated", Data: nil})
}

// Team --------------------------------------------

func (d *Domain) handlerGetFSPTeam(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	TenantID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerGetFSPTeam:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if TenantID == nil {
		d.logger.Error("handlerGetFSPTeam:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	existingTeam, err := d.getTeamService(*TenantID)
	if err != nil {
		d.logger.Error("teamHandler:", zap.Error(err), zap.String("traceID", traceID))
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "team", Data: existingTeam})
}
func (d *Domain) handlerAddFSPTeam(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	TenantID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerAddFSPTeam:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if TenantID == nil {
		d.logger.Error("handlerAddFSPTeam:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	tenantIdentifier, err := extractor.ExtractTenantIdentifierFromContext(c)
	if err != nil {
		d.logger.Error("handlerAddFSPTeam:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if tenantIdentifier == nil {
		d.logger.Error("handlerAddFSPTeam:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	var postTeamReq postTeamRequest

	err = c.Bind(&postTeamReq)
	if err != nil {
		d.logger.Error("handlerAddFSPTeam:", zap.Error(err), zap.String("traceID", traceID))
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
	err = c.Validate(postTeamReq)
	if err != nil {
		d.logger.Error("handlerAddFSPTeam:", zap.Error(err), zap.String("traceID", traceID))
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

	err = d.postTeamService(
		*tenantIdentifier, //tenantIdentifier string
		*TenantID,         //TenantID 			uint
		postTeamReq.Email, //email 			string
		postTeamReq.Role,  //startingRole 	schema.RoleConst
	)
	if err != nil {
		if errors.Is(err, ErrTeamMemberHasRole) {
			return c.JSON(
				http.StatusConflict,
				response.APIResponse{
					Message: "user already has a role",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeUserHasFSPRole,
						TraceID: traceID,
					},
				})
		}
		d.logger.Error("handlerAddFSPTeam:", zap.Error(err))
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "team member added", Data: nil})
}
func (d *Domain) handlerUpdateFSPTeam(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	TenantID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerUpdateFSPTeam:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if TenantID == nil {
		d.logger.Error("handlerUpdateFSPTeam:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	teamMemberID, err := extractor.ExtractIDFromPathParamAsUINT(c, "teamMemberID")
	if err != nil {
		d.logger.Error("handlerUpdateFSPTeam:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "invalid or missing path/query parameters",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}
	if teamMemberID == nil {
		d.logger.Error("handlerUpdateFSPTeam:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	var putTeamReq putTeamRequest

	err = c.Bind(&putTeamReq)
	if err != nil {
		d.logger.Error("handlerUpdateFSPTeam:", zap.Error(err), zap.String("traceID", traceID))
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
	err = c.Validate(putTeamReq)
	if err != nil {
		d.logger.Error("handlerUpdateFSPTeam:", zap.Error(err), zap.String("traceID", traceID))
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

	err = d.putTeamService(*TenantID, *teamMemberID, putTeamReq.Role)
	if err != nil {
		d.logger.Error("handlerUpdateFSPTeam:", zap.Error(err), zap.String("traceID", traceID))
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "team member updated", Data: nil})
}
func (d *Domain) handlerDeleteFSPTeam(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	TenantID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerDeleteFSPTeam:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if TenantID == nil {
		d.logger.Error("handlerDeleteFSPTeam:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	tenantIdentifier, err := extractor.ExtractTenantIdentifierFromContext(c)
	if err != nil {
		d.logger.Error("handlerDeleteFSPTeam:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if tenantIdentifier == nil {
		d.logger.Error("handlerDeleteFSPTeam:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	teamMemberID, err := extractor.ExtractIDFromPathParamAsUINT(c, "teamMemberID")
	if err != nil {
		d.logger.Error("handlerDeleteFSPTeam:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "invalid or missing path/query parameters",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}
	if teamMemberID == nil {
		d.logger.Error("handlerDeleteFSPTeam:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	var deleteTeamReq deleteTeamRequest

	err = c.Bind(&deleteTeamReq)
	if err != nil {
		d.logger.Error("handlerDeleteFSPTeam:", zap.Error(err), zap.String("traceID", traceID))
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
	err = c.Validate(deleteTeamReq)
	if err != nil {
		d.logger.Error("handlerDeleteFSPTeam:", zap.Error(err), zap.String("traceID", traceID))
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

	err = d.deleteTeamService(
		*TenantID,                //TenantID 				uint
		*teamMemberID,            //teamMemberID 		uint
		deleteTeamReq.NotifyUser, //notifyUser 			bool
	)
	if err != nil {
		d.logger.Error("handlerDeleteFSPTeam:", zap.Error(err), zap.String("traceID", traceID))
		if errors.Is(err, ErrUserIsLastSuperAdmin) {
			return c.JSON(
				http.StatusConflict,
				response.APIResponse{
					Message: "cannot delete last super admin",
					Data:    nil,
					Error: errmgr.ErrResponse{
						Code:    errmgr.ErrCodeUserIsLastSuperAdmin,
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "team member deleted", Data: nil})
}

// Memorial --------------------------------------------

func (d *Domain) handlerGetMemorials(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	TenantID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerGetMemorials:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if TenantID == nil {
		d.logger.Error("handlerGetMemorials:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	existingMemorials, err := d.getAllMemorials(*TenantID)
	if err != nil {
		d.logger.Error("handlerGetMemorials:", zap.Error(err), zap.String("traceID", traceID))
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "memorials", Data: existingMemorials})
}
func (d *Domain) handlerAddMemorial(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	TenantID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerAddMemorial:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if TenantID == nil {
		d.logger.Error("handlerAddMemorial:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	var postMemorialReq postMemorialRequest

	err = c.Bind(&postMemorialReq)
	if err != nil {
		d.logger.Error("handlerAddMemorial:", zap.Error(err), zap.String("traceID", traceID))
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
	err = c.Validate(postMemorialReq)
	if err != nil {
		d.logger.Error("handlerAddMemorial:", zap.Error(err), zap.String("traceID", traceID))
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

	//parse the date strings
	var DOB *time.Time
	var DOD *time.Time
	if postMemorialReq.DOBString != nil {
		parsedDOB, err := time.Parse("2006-01-02", *postMemorialReq.DOBString)
		if err != nil {
			d.logger.Error("handlerAddMemorial:", zap.Error(err), zap.String("traceID", traceID))
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
		DOB = &parsedDOB
	}
	if postMemorialReq.DODString != nil {
		parsedDOD, err := time.Parse("2006-01-02", *postMemorialReq.DODString)
		if err != nil {
			d.logger.Error("handlerAddMemorial:", zap.Error(err), zap.String("traceID", traceID))
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
		DOD = &parsedDOD
	}

	err = d.createOrUpdateMemorialWithUserAndCuratorRole(
		*TenantID,                           //TenantID 			uint
		postMemorialReq.FirstName,           //FirstName 		string
		postMemorialReq.LastName,            //LastName 		string
		DOB,                                 //DOB 			*time.Time
		DOD,                                 //DOD 			*time.Time
		postMemorialReq.CuratorEmail,        //CuratorEmail 	string
		postMemorialReq.CuratorRelationship, //relationship 	string
	)
	if err != nil {
		d.logger.Error("handlerAddMemorial:", zap.Error(err), zap.String("traceID", traceID))
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "memorial added", Data: nil})
}
func (d *Domain) handlerUpdateMemorial(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	TenantID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerUpdateMemorial:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if TenantID == nil {
		d.logger.Error("handlerUpdateMemorial:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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
		d.logger.Error("handlerUpdateMemorial:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "invalid or missing path/query parameters",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}
	if memorialID == nil {
		d.logger.Error("handlerUpdateMemorial:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	putMemorialReq := putMemorialRequest{}

	err = c.Bind(&putMemorialReq)
	if err != nil {
		d.logger.Error("handlerUpdateMemorial:", zap.Error(err), zap.String("traceID", traceID))
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
	err = c.Validate(putMemorialReq)
	if err != nil {
		d.logger.Error("handlerUpdateMemorial:", zap.Error(err), zap.String("traceID", traceID))
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

	//parse the date strings
	var DOB *time.Time
	var DOD *time.Time
	if putMemorialReq.DOBString != nil {
		parsedDOB, err := time.Parse("2006-01-02", *putMemorialReq.DOBString)
		if err != nil {
			d.logger.Error("handlerUpdateMemorial:", zap.Error(err), zap.String("traceID", traceID))
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
		DOB = &parsedDOB
	}
	if putMemorialReq.DODString != nil {
		parsedDOD, err := time.Parse("2006-01-02", *putMemorialReq.DODString)
		if err != nil {
			d.logger.Error("handlerUpdateMemorial:", zap.Error(err), zap.String("traceID", traceID))
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
		DOD = &parsedDOD
	}

	updatedMemorial := schema.Memorial{}
	if putMemorialReq.Title != nil {
		updatedMemorial.Title = *putMemorialReq.Title
	}
	if putMemorialReq.FirstName != nil {
		updatedMemorial.FirstName = *putMemorialReq.FirstName
	}
	if putMemorialReq.LastName != nil {
		updatedMemorial.LastName = *putMemorialReq.LastName
	}
	updatedMemorial.DOB = DOB
	updatedMemorial.DOD = DOD

	err = d.updateMemorial(
		*TenantID,       //TenantID 			uint
		*memorialID,     //memorialID 		uint
		updatedMemorial, //updatedMemorial 	schema.Memorial
	)
	if err != nil {
		d.logger.Error("handlerUpdateMemorial:", zap.Error(err), zap.String("traceID", traceID))
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "memorial updated", Data: nil})
}
func (d *Domain) handlerDeleteMemorial(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	TenantID, err := extractor.ExtractFSPIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerDeleteMemorial:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the tenant identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if TenantID == nil {
		d.logger.Error("handlerDeleteMemorial:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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
		d.logger.Error("handlerDeleteMemorial:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "invalid or missing path/query parameters",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeInput,
					TraceID: traceID,
				},
			})
	}
	if memorialID == nil {
		d.logger.Error("handlerDeleteMemorial:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	err = d.deleteMemorial(*TenantID, *memorialID)
	if err != nil {
		d.logger.Error("handlerDeleteMemorial:", zap.Error(err), zap.String("traceID", traceID))
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

	return c.JSON(http.StatusOK, response.APIResponse{Message: "memorial deleted", Data: nil})
}

// Partners --------------------------------------------

func (d *Domain) handlerGetPartners(c echo.Context) error {
	return c.JSON(http.StatusOK, response.APIResponse{Message: "partners", Data: nil})
}
