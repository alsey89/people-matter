package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/alsey89/people-matter/internal/common"
)

// ! User ------------------------------------------------------------

func (d *Domain) GetCurrentUserHandler(c echo.Context) error {
	companyID, err := common.GetCompanyIDFromToken(c)
	if err != nil {
		d.logger.Error("[getCurrentUser]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting company id from token",
			Data:    nil,
		})
	}

	userID, err := common.GetUserIDFromToken(c)
	if err != nil {
		d.logger.Error("[getCurrentUser]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting user id from token",
			Data:    nil,
		})
	}

	user, err := d.GetUser(companyID, userID)
	if err != nil {
		d.logger.Error("[getCurrentUser]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting user",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "success",
		Data:    user,
	})
}

func (d *Domain) GetAllUsersHandler(c echo.Context) error {
	companyID, err := common.GetCompanyIDFromToken(c)
	if err != nil {
		d.logger.Error("[GetAllUsersHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting company id from token",
			Data:    nil,
		})
	}

	users, err := d.GetAllUsers(companyID)
	if err != nil {
		d.logger.Error("[GetAllUsersHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting users",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "success",
		Data:    users,
	})
}

func (d *Domain) GetAllLocationUsersHandler(c echo.Context) error {
	companyID, err := common.GetCompanyIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting company id from token",
			Data:    nil,
		})
	}

	locationID, err := common.GetIDFromParam("location_id", c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting location id from param",
			Data:    nil,
		})
	}

	users, err := d.GetUsersByLocation(companyID, locationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting location users",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "success",
		Data:    users,
	})
}
