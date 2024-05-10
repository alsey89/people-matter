package company

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
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
	if err != nil {
		d.logger.Error("[GetCompanyHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting company data",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "detailed company data retrieved",
		Data:    companyData,
	})
}

func (d *Domain) CreateCompanyHandler(c echo.Context) error {

	form := new(NewCompany)
	err := c.Bind(form)
	if err != nil {
		d.logger.Error("[signupHandler] error binding credentials", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "form error",
			Data:    nil,
		})
	}

	// validate email
	email := form.AdminEmail
	if !common.EmailValidator(email) {
		d.logger.Error("[signupHandler] email validation failed")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "invalid email",
			Data:    nil,
		})
	}

	// validate password
	password := form.Password
	confirmPassword := form.ConfirmPassword
	if password != confirmPassword {
		d.logger.Error("[signupHandler] password confirmation check failed")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "passwords do not match",
			Data:    nil,
		})
	}

	// validate company name
	companyName := form.CompanyName
	if companyName == "" {
		d.logger.Error("[signupHandler] company name is empty")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "company name required",
			Data:    nil,
		})
	}

	createdAdminUser, err := d.CreateNewCompanyAndAdminUser(form)
	if err != nil {
		d.logger.Error("[createCompanyHandler]", zap.Error(err))
		if errors.Is(err, ErrUserExists) {
			return c.JSON(http.StatusConflict, common.APIResponse{
				Message: "user already exists",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error creating company and admin user",
			Data:    nil,
		})
	}

	//generate jwt token
	additionalClaims := jwt.MapClaims{
		"Id":        createdAdminUser.ID,
		"companyId": createdAdminUser.CompanyID,
	}

	token, err := d.params.Auth.GenerateToken(additionalClaims)
	if err != nil {
		d.logger.Error("[CreateNewCompanyAndAdminUser]", zap.Error(err))
	}

	//send confirmation email
	// todo define send mail function in gogetter
	m := d.params.Mailer.NewMessage()
	m.SetHeader("From", "hello@peoplematter.app")
	m.SetHeader("To", form.AdminEmail)
	m.SetHeader("Subject", "Welcome to People Matter")
	m.SetBody("text/html",
		"<p>Welcome to People Matter</p><p>Your account has been created. Please click the link below to confirm your email address.</p><a href=\"http://localhost:3000/onboarding/confirmation?token="+*token+"\">Confirm Email</a>",
	)
	err = d.params.Mailer.Send(m)
	if err != nil {
		return fmt.Errorf("[CreateNewCompanyAndAdminUser] Error sending confirmation email %w", err)
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "company and admin user created",
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
	if err != nil {
		d.logger.Error("[UpdateCompanyHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error updating company",
			Data:    nil,
		})
	}

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
	if err != nil {
		d.logger.Error("[DeleteCompanyHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error deleting company",
			Data:    nil,
		})
	}

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
	if err != nil {
		d.logger.Error("[CreateDepartmentHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error creating department",
			Data:    nil,
		})
	}

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
