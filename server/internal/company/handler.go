package company

import (
	"net/http"

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

func (d *Domain) handlerGetCompanyAccount(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	CompanyID, err := extractor.ExtractTenantIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerGetCompanyAccount:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the Company identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if CompanyID == nil {
		d.logger.Error("handlerGetCompanyAccount:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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
		*CompanyID, //CompanyID 			uint
		true,       //preloadDetails	bool
	)
	if err != nil {
		d.logger.Error("handlerGetCompanyAccount:", zap.Error(err), zap.String("traceID", traceID))
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
func (d *Domain) handlerUpdateCompanyAccount(c echo.Context) error {
	var err error
	traceID := uuid.NewString()

	CompanyID, err := extractor.ExtractTenantIDFromContext(c)
	if err != nil {
		d.logger.Error("handlerUpdateCompanyAccount:", zap.Error(err), zap.String("traceID", traceID))
		return c.JSON(
			http.StatusBadRequest,
			response.APIResponse{
				Message: "error in the Company identification chain",
				Data:    nil,
				Error: errmgr.ErrResponse{
					Code:    errmgr.ErrCodeTenant,
					TraceID: traceID,
				},
			})
	}
	if CompanyID == nil {
		d.logger.Error("handlerUpdateCompanyAccount:", zap.Error(errmgr.ErrNilCheckFailed), zap.String("traceID", traceID))
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

	putMemorialReq := schema.Company{}

	err = c.Bind(&putMemorialReq)
	if err != nil {
		d.logger.Error("handlerUpdateCompanyAccount:", zap.Error(err), zap.String("traceID", traceID))
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

	err = d.updateAccountService(*CompanyID, putMemorialReq)
	if err != nil {
		d.logger.Error("handlerUpdateCompanyAccount:", zap.Error(err), zap.String("traceID", traceID))
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

// Positions --------------------------------------------

func (d *Domain) handlerGetPositions(c echo.Context) error {
	return c.JSON(http.StatusOK, response.APIResponse{Message: "positions", Data: nil})
}
func (d *Domain) handlerCreatePosition(c echo.Context) error {
	return c.JSON(http.StatusOK, response.APIResponse{Message: "position created", Data: nil})
}
func (d *Domain) handlerUpdatePosition(c echo.Context) error {
	return c.JSON(http.StatusOK, response.APIResponse{Message: "position updated", Data: nil})
}
func (d *Domain) handlerDeletePosition(c echo.Context) error {
	return c.JSON(http.StatusOK, response.APIResponse{Message: "position deleted", Data: nil})
}
