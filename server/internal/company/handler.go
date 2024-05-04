package company

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/alsey89/people-matter/internal/common"
	"github.com/alsey89/people-matter/schema"
)

// ! Company ------------------------------------------------------------

// fetches company data *with preloaded department, location, position data*
func (d *Domain) GetCompanyHandler(c echo.Context) error {
	companyID, err := common.GetCompanyIDFromToken(c)
	if err != nil {
		d.logger.Error("[GetCompanyHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting company id from token",
			Data:    nil,
		})
	}

	companyData, err := d.GetCompanyWithDetails(companyID)

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "detailed company data retrieved",
		Data:    companyData,
	})
}

func (d *Domain) CreateCompanyHandler(c echo.Context) error {
	newCompany := new(NewCompany)

	err := c.Bind(newCompany)
	if err != nil {
		d.logger.Error("[CreateCompanyHandler] error binding newCompany data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	err = d.CreateNewCompanyAndAdminUser(newCompany)
	if err != nil {
		d.logger.Error("[createCompanyHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error creating company",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "company created",
		Data:    nil,
	})
}

func (d *Domain) UpdateCompanyHandler(c echo.Context) error {
	companyID, err := common.GetCompanyIDFromToken(c)
	if err != nil {
		d.logger.Error("[UpdateCompanyHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting company id from token",
			Data:    nil,
		})
	}

	dataToUpdate := new(schema.Company)

	err = c.Bind(dataToUpdate)
	if err != nil {
		d.logger.Error("[UpdateCompanyHandler] error binding dataToUpdate", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	err = d.UpdateCompany(companyID, dataToUpdate)

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "company data updated",
		Data:    nil,
	})
}

func (d *Domain) DeleteCompanyHandler(c echo.Context) error {
	companyID, err := common.GetCompanyIDFromToken(c)
	if err != nil {
		d.logger.Error("[DeleteCompanyHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting company id from token",
			Data:    nil,
		})
	}

	err = d.DeleteCompany(companyID)
	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "company deleted",
		Data:    nil,
	})
}

//! Department ------------------------------------------------------------

func (d *Domain) CreateDepartmentHandler(c echo.Context) error {
	companyID, err := common.GetCompanyIDFromToken(c)
	if err != nil {
		d.logger.Error("[CreateDepartmentHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting company id from token",
			Data:    nil,
		})
	}
	if companyID == nil {
		d.logger.Error("[CreateDepartmentHandler] companyID is nil")
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "company id is nil",
			Data:    nil,
		})
	}

	newDepartment := new(schema.Department)

	newDepartment.CompanyID = *companyID

	err = c.Bind(newDepartment)
	if err != nil {
		d.logger.Error("[CreateDepartmentHandler] error binding newDepartment data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	err = d.CreateDepartment(newDepartment)

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "department created",
		Data:    nil,
	})
}

func (d *Domain) UpdateDepartmentHandler(c echo.Context) error {
	companyID, err := common.GetCompanyIDFromToken(c)
	if err != nil {
		d.logger.Error("[UpdateDepartmentHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting company id from token",
			Data:    nil,
		})
	}

	departmentID, err := common.GetIDFromParam("department_id", c)
	if err != nil {
		d.logger.Error("[UpdateDepartmentHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting department id from param",
			Data:    nil,
		})
	}

	dataToUpdate := new(schema.Department)

	err = c.Bind(dataToUpdate)
	if err != nil {
		d.logger.Error("[UpdateDepartmentHandler] error binding dataToUpdate", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	err = d.UpdateDepartment(companyID, departmentID, dataToUpdate)
	if err != nil {
		d.logger.Error("[UpdateDepartmentHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error updating department",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "department updated",
		Data:    nil,
	})
}

func (d *Domain) DeleteDepartmentHandler(c echo.Context) error {
	companyID, err := common.GetCompanyIDFromToken(c)
	if err != nil {
		d.logger.Error("[DeleteDepartmentHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting company id from token",
			Data:    nil,
		})
	}

	departmentID, err := common.GetIDFromParam("department_id", c)
	if err != nil {
		d.logger.Error("[DeleteDepartmentHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting department id from param",
			Data:    nil,
		})
	}

	err = d.DeleteDepartment(companyID, departmentID)
	if err != nil {
		d.logger.Error("[DeleteDepartmentHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error deleting department",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "department has been deleted",
		Data:    nil,
	})
}

//! Location ------------------------------------------------------------

func (d *Domain) CreateLocationHandler(c echo.Context) error {
	companyID, err := common.GetCompanyIDFromToken(c)
	if err != nil {
		d.logger.Error("[CreateLocationHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting company id from token",
			Data:    nil,
		})
	}
	// to avoid nil pointer error
	if companyID == nil {
		d.logger.Error("[CreateLocationHandler] companyID is nil")
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "company id is nil",
			Data:    nil,
		})
	}

	newLocation := new(schema.Location)

	newLocation.CompanyID = *companyID

	err = c.Bind(newLocation)
	if err != nil {
		d.logger.Error("[CreateLocationHandler] error binding newLocation data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	err = d.CreateLocation(newLocation)

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "location created",
		Data:    nil,
	})
}

func (d *Domain) UpdateLocationHandler(c echo.Context) error {
	companyID, err := common.GetCompanyIDFromToken(c)
	if err != nil {
		d.logger.Error("[UpdateLocationHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting company id from token",
			Data:    nil,
		})
	}

	locationID, err := common.GetIDFromParam("location_id", c)
	if err != nil {
		d.logger.Error("[UpdateLocationHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting location id from param",
			Data:    nil,
		})
	}

	dataToUpdate := new(schema.Location)

	err = c.Bind(dataToUpdate)
	if err != nil {
		d.logger.Error("[UpdateLocationHandler] error binding dataToUpdate", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	err = d.UpdateLocation(companyID, locationID, dataToUpdate)
	if err != nil {
		d.logger.Error("[UpdateLocationHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error updating location",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "location data updated",
		Data:    nil,
	})
}

func (d *Domain) DeleteLocationHandler(c echo.Context) error {
	companyID, err := common.GetCompanyIDFromToken(c)
	if err != nil {
		d.logger.Error("[DeleteLocationHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting company id from token",
			Data:    nil,
		})
	}

	locationID, err := common.GetIDFromParam("location_id", c)
	if err != nil {
		d.logger.Error("[DeleteLocationHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting location id from param",
			Data:    nil,
		})
	}

	err = d.DeleteLocation(companyID, locationID)
	if err != nil {
		d.logger.Error("[DeleteLocationHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error deleting location",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "location deleted",
		Data:    nil,
	})
}

//! Position ------------------------------------------------------------

func (d *Domain) CreatePositionHandler(c echo.Context) error {
	companyID, err := common.GetCompanyIDFromToken(c)
	if err != nil {
		d.logger.Error("[CreatePositionHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting company id from token",
			Data:    nil,
		})
	}
	// to avoid nil pointer error
	if companyID == nil {
		d.logger.Error("[CreatePositionHandler] companyID is nil")
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "company id is nil",
			Data:    nil,
		})
	}

	newPosition := new(schema.Position)

	newPosition.CompanyID = *companyID

	err = c.Bind(newPosition)
	if err != nil {
		d.logger.Error("[CreatePositionHandler] error binding newPosition data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	err = d.CreatePosition(newPosition)
	if err != nil {
		d.logger.Error("[CreatePositionHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error creating position",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "position data has been created",
		Data:    nil,
	})
}

func (d *Domain) UpdatePositionHandler(c echo.Context) error {
	companyID, err := common.GetCompanyIDFromToken(c)
	if err != nil {
		d.logger.Error("[UpdatePositionHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting company id from token",
			Data:    nil,
		})
	}

	positionID, err := common.GetIDFromParam("position_id", c)
	if err != nil {
		d.logger.Error("[UpdatePositionHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting position id from param",
			Data:    nil,
		})
	}

	dataToUpdate := new(schema.Position)

	err = c.Bind(dataToUpdate)
	if err != nil {
		d.logger.Error("[UpdatePositionHandler] error binding dataToUpdate", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	err = d.UpdatePosition(companyID, positionID, dataToUpdate)
	if err != nil {
		d.logger.Error("[UpdatePositionHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error updating position",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "position updated",
		Data:    nil,
	})
}

func (d *Domain) DeletePositionHandler(c echo.Context) error {
	companyID, err := common.GetCompanyIDFromToken(c)
	if err != nil {
		d.logger.Error("[DeletePositionHandlers]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting company id from token",
			Data:    nil,
		})
	}

	positionID, err := common.GetIDFromParam("position_id", c)
	if err != nil {
		d.logger.Error("[DeletePositionHandlers]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting position id from param",
			Data:    nil,
		})
	}

	err = d.DeletePosition(companyID, positionID)
	if err != nil {
		d.logger.Error("[DeletePositionHandlers]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error deleting position",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "position deleted",
		Data:    nil,
	})

}
