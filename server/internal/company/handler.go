package company

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/alsey89/hrms/internal/common"
	"github.com/alsey89/hrms/internal/schema"
)

type CompanyHandler struct {
	companyService *CompanyService
}

func NewCompanyHandler(companyService *CompanyService) *CompanyHandler {
	return &CompanyHandler{companyService: companyService}
}

// ! Company ------------------------------------------------------------
// todo: switch to single company data
func (ch *CompanyHandler) GetCompanyDataExpandDefault(c echo.Context) error {

	CompanyData, err := ch.companyService.GetCompanyListAndExpandDefault()
	if err != nil {
		log.Printf("company.h.get_company_data_expand_default: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrEmptyTable) {
			return c.JSON(http.StatusNoContent, common.APIResponse{
				Message: "no company data",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error retrieving company data",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "company data has been retrieved",
		Data:    CompanyData,
	})
}

func (ch *CompanyHandler) GetCompanyDataExpandID(c echo.Context) error {
	stringCompanyID := c.Param("company_id")
	if stringCompanyID == "" {
		log.Printf("company.h.get_company_data_expand_id: empty company id parameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no company id",
			Data:    nil,
		})
	}

	companyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
	if err != nil {
		log.Printf("company.h.get_company_data_expand_default: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing company id",
			Data:    nil,
		})
	}

	CompanyData, err := ch.companyService.GetCompanyListAndExpandByID(companyID)
	if err != nil {
		log.Printf("company.h.fetch_company_data_expand_id: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrEmptyTable) {
			return c.JSON(http.StatusNoContent, common.APIResponse{
				Message: "no company data",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error retrieving company data",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "company data has been retrieved",
		Data:    CompanyData,
	})
}

func (ch *CompanyHandler) CreateCompany(c echo.Context) error {
	newCompany := new(schema.Company)

	err := c.Bind(newCompany)
	if err != nil {
		log.Printf("company.h.create_company: error binding company data: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	companyData, err := ch.companyService.CreateNewCompanyAndReturnList(newCompany)
	if err != nil {
		log.Printf("company.h.create_company: %v", err)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.JSON(http.StatusConflict, common.APIResponse{
				Message: "company already exists",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error creating company data",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "company data has been created",
		Data:    companyData,
	})
}

func (ch *CompanyHandler) UpdateCompany(c echo.Context) error {
	stringCompanyID := c.Param("company_id")
	if stringCompanyID == "" {
		log.Printf("company.h.update_company: empty company id parameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no company id",
			Data:    nil,
		})
	}

	companyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
	if err != nil {
		log.Printf("company.h.update_company: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing company id",
			Data:    nil,
		})
	}

	dataToUpdate := new(schema.Company)

	err = c.Bind(dataToUpdate)
	if err != nil {
		log.Printf("company.h.update_company: error binding company data: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	companyData, err := ch.companyService.UpdateCompanyAndReturnCompanyListAndExpandID(companyID, dataToUpdate)
	if err != nil {
		log.Printf("company.h.update_company: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrEmptyTable) {
			return c.JSON(http.StatusNoContent, common.APIResponse{
				Message: "no company data",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error updating company data",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "company data has been updated",
		Data:    companyData,
	})
}

func (ch *CompanyHandler) DeleteCompany(c echo.Context) error {
	stringCompanyID := c.Param("company_id")
	if stringCompanyID == "" {
		log.Printf("company.h.delete_company: empty company id parameterparameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no company id",
			Data:    nil,
		})
	}

	companyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
	if err != nil {
		log.Printf("company.h.delete_company: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing company id",
			Data:    nil,
		})
	}

	companyData, err := ch.companyService.DeleteCompanyAndReturnCompanyListAndExpandDefault(companyID)
	if err != nil {
		log.Printf("company.h.delete_company: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrEmptyTable) {
			return c.JSON(http.StatusNoContent, common.APIResponse{
				Message: "no company data",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error deleting company",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "company has been deleted",
		Data:    companyData,
	})
}

//! Department ------------------------------------------------------------

func (ch *CompanyHandler) CreateDepartment(c echo.Context) error {
	stringCompanyID := c.Param("company_id")
	if stringCompanyID == "" {
		log.Printf("company.h.create_department: empty company id parameterparameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no company id",
			Data:    nil,
		})
	}

	companyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
	if err != nil {
		log.Printf("company.h.create_department: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing company id",
			Data:    nil,
		})
	}

	newDepartment := new(schema.Department)

	err = c.Bind(newDepartment)
	if err != nil {
		log.Printf("company.h.create_department: error binding department data: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	companyData, err := ch.companyService.CreateNewDepartmentAndReturnCompanyListAndExpandID(companyID, newDepartment)
	if err != nil {
		log.Printf("company.h.create_department: %v", err)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.JSON(http.StatusConflict, common.APIResponse{
				Message: "department already exists",
				Data:    nil,
			})
		}
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrEmptyTable) {
			return c.JSON(http.StatusNoContent, common.APIResponse{
				Message: "no department data",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error creating department data",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "department data has been created",
		Data:    companyData,
	})
}

func (ch *CompanyHandler) UpdateDepartment(c echo.Context) error {
	stringCompanyID := c.Param("company_id")
	if stringCompanyID == "" {
		log.Printf("company.h.update_department: empty company id parameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no company id",
			Data:    nil,
		})
	}

	companyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
	if err != nil {
		log.Printf("company.h.update_department: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing company id",
			Data:    nil,
		})
	}

	stringDepartmentID := c.Param("department_id")
	if stringDepartmentID == "" {
		log.Printf("company.h.update_department: empty department id parameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no department id",
			Data:    nil,
		})
	}

	departmentID, err := common.ConvertStringOfNumbersToUint(stringDepartmentID)
	if err != nil {
		log.Printf("company.h.update_department: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing department id",
			Data:    nil,
		})
	}

	dataToUpdate := new(schema.Department)
	err = c.Bind(dataToUpdate)
	if err != nil {
		log.Printf("company.h.update_department: error binding department data: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	companyData, err := ch.companyService.UpdateDepartmentAndReturnCompanyListAndExpandID(companyID, departmentID, dataToUpdate)
	if err != nil {
		log.Printf("company.h.update_department: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrEmptyTable) {
			return c.JSON(http.StatusNoContent, common.APIResponse{
				Message: "no department data",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error updating department data",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "department data has been updated",
		Data:    companyData,
	})
}

func (ch *CompanyHandler) DeleteDepartment(c echo.Context) error {
	stringCompanyID := c.Param("company_id")
	if stringCompanyID == "" {
		log.Printf("company.h.delete_department: empty company id parameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no company id",
			Data:    nil,
		})
	}

	companyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
	if err != nil {
		log.Printf("company.h.delete_department: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing company id",
			Data:    nil,
		})
	}

	stringDepartmentID := c.Param("department_id")
	if stringDepartmentID == "" {
		log.Printf("company.h.delete_department: empty department id parameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no department id",
			Data:    nil,
		})
	}

	departmentID, err := common.ConvertStringOfNumbersToUint(stringDepartmentID)
	if err != nil {
		log.Printf("company.h.delete_department: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing department id",
			Data:    nil,
		})
	}

	companyData, err := ch.companyService.DeleteDepartmentAndReturnCompanyListAndExpandID(companyID, departmentID)
	if err != nil {
		log.Printf("company.h.delete_department: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrEmptyTable) {
			return c.JSON(http.StatusNoContent, common.APIResponse{
				Message: "no department data",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error deleting department",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "department has been deleted",
		Data:    companyData,
	})
}

//! Location ------------------------------------------------------------

func (ch *CompanyHandler) CreateLocation(c echo.Context) error {
	stringCompanyID := c.Param("company_id")
	if stringCompanyID == "" {
		log.Printf("company.h.create_location: empty company id parameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no company id",
			Data:    nil,
		})
	}

	uintCompanyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
	if err != nil {
		log.Printf("company.h.create_location: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing company id",
			Data:    nil,
		})
	}

	newLocation := new(schema.Location)

	err = c.Bind(newLocation)
	if err != nil {
		log.Printf("company.h.create_location: error binding location data: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	companyData, err := ch.companyService.CreateNewLocationAndReturnCompanyListAndExpandID(uintCompanyID, newLocation)
	if err != nil {
		log.Printf("company.h.create_location: %v", err)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.JSON(http.StatusConflict, common.APIResponse{
				Message: "location already exists",
				Data:    nil,
			})
		}
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrEmptyTable) {
			return c.JSON(http.StatusNoContent, common.APIResponse{
				Message: "no location data",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error creating location data",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "location data has been created",
		Data:    companyData,
	})
}

func (ch *CompanyHandler) UpdateLocation(c echo.Context) error {
	stringCompanyID := c.Param("company_id")
	if stringCompanyID == "" {
		log.Printf("company.h.update_location: empty company id parameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no company id",
			Data:    nil,
		})
	}

	uintCompanyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
	if err != nil {
		log.Printf("company.h.update_location: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing company id",
			Data:    nil,
		})
	}

	stringLocationID := c.Param("location_id")
	if stringLocationID == "" {
		log.Printf("company.h.update_location: empty location id parameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no location id",
			Data:    nil,
		})
	}

	uinLocationID, err := common.ConvertStringOfNumbersToUint(stringLocationID)
	if err != nil {
		log.Printf("company.h.update_location: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing location id",
			Data:    nil,
		})
	}

	dataToUpdate := new(schema.Location)
	err = c.Bind(dataToUpdate)
	if err != nil {
		log.Printf("company.h.update_location: error binding location data: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	companyData, err := ch.companyService.UpdateLocationAndReturnCompanyListAndExpandID(uintCompanyID, uinLocationID, dataToUpdate)
	if err != nil {
		log.Printf("company.h.update_location: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrEmptyTable) {
			return c.JSON(http.StatusNoContent, common.APIResponse{
				Message: "no location data",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error updating location data",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "location data has been updated",
		Data:    companyData,
	})
}

func (ch *CompanyHandler) DeleteLocation(c echo.Context) error {
	stringCompanyID := c.Param("company_id")
	if stringCompanyID == "" {
		log.Printf("company.h.delete_location: empty company id parameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no company id",
			Data:    nil,
		})
	}

	uintCompanyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
	if err != nil {
		log.Printf("company.h.delete_location: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing company id",
			Data:    nil,
		})
	}

	stringLocationID := c.Param("location_id")
	if stringLocationID == "" {
		log.Printf("company.h.delete_location: empty location id parameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no location id",
			Data:    nil,
		})
	}

	uintLocationID, err := common.ConvertStringOfNumbersToUint(stringLocationID)
	if err != nil {
		log.Printf("company.h.delete_location: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing location id",
			Data:    nil,
		})
	}

	companyData, err := ch.companyService.DeleteLocationAndReturnCompanyListAndExpandID(uintCompanyID, uintLocationID)
	if err != nil {
		log.Printf("company.h.delete_location: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrEmptyTable) {
			return c.JSON(http.StatusNoContent, common.APIResponse{
				Message: "no location data",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error deleting location",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "location has been deleted",
		Data:    companyData,
	})
}

//! Jobs ------------------------------------------------------------

func (ch *CompanyHandler) CreateJob(c echo.Context) error {
	stringCompanyID := c.Param("company_id")
	if stringCompanyID == "" {
		log.Printf("company.h.create_job: empty company id parameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no company id",
			Data:    nil,
		})
	}

	uintCompanyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
	if err != nil {
		log.Printf("company.h.create_job: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing company id",
			Data:    nil,
		})
	}

	newJob := new(schema.Job)
	err = c.Bind(newJob)
	if err != nil {
		log.Printf("company.h.create_job: error binding job data: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	companyData, err := ch.companyService.CreateNewJobAndReturnCompanyListAndExpandID(uintCompanyID, newJob)
	if err != nil {
		log.Printf("company.h.create_job: %v", err)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.JSON(http.StatusConflict, common.APIResponse{
				Message: "job already exists",
				Data:    nil,
			})
		}
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrEmptyTable) {
			return c.JSON(http.StatusNoContent, common.APIResponse{
				Message: "no job data",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error creating job data",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "job data has been created",
		Data:    companyData,
	})
}

func (ch *CompanyHandler) UpdateJob(c echo.Context) error {

	stringCompanyID := c.Param("company_id")
	if stringCompanyID == "" {
		log.Printf("company.h.update_job: empty company id parameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no company id",
			Data:    nil,
		})
	}

	uintCompanyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
	if err != nil {
		log.Printf("company.h.update_job: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing company id",
			Data:    nil,
		})
	}

	stringJobID := c.Param("job_id")
	if stringJobID == "" {
		log.Printf("company.h.update_job: empty job id parameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no job id",
			Data:    nil,
		})
	}

	uintJobID, err := common.ConvertStringOfNumbersToUint(stringJobID)
	if err != nil {
		log.Printf("company.h.update_job: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing job id",
			Data:    nil,
		})
	}

	dataToUpdate := new(schema.Job)

	err = c.Bind(dataToUpdate)
	if err != nil {
		log.Printf("company.h.update_job: error binding job data: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	companyData, err := ch.companyService.UpdateJobAndReturnCompanyListAndExpandID(uintCompanyID, uintJobID, dataToUpdate)
	if err != nil {
		log.Printf("company.h.update_job: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrEmptyTable) {
			return c.JSON(http.StatusNoContent, common.APIResponse{
				Message: "no job data",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error updating job data",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "job data has been created",
		Data:    companyData,
	})
}

func (ch *CompanyHandler) DeleteJob(c echo.Context) error {
	stringCompanyID := c.Param("company_id")
	if stringCompanyID == "" {
		log.Printf("company.h.delete_job: empty company id parameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no company id",
			Data:    nil,
		})
	}

	uintCompanyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
	if err != nil {
		log.Printf("company.h.delete_job: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing company id",
			Data:    nil,
		})
	}

	stringJobID := c.Param("job_id")
	if stringJobID == "" {
		log.Printf("company.h.delete_job: empty job id parameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no job id",
			Data:    nil,
		})
	}

	uintJobID, err := common.ConvertStringOfNumbersToUint(stringJobID)
	if err != nil {
		log.Printf("company.h.delete_job: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing job id",
			Data:    nil,
		})
	}

	companyData, err := ch.companyService.DeleteJobAndReturnCompanyListAndExpandID(uintCompanyID, uintJobID)
	if err != nil {
		log.Printf("company.h.delete_job: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrEmptyTable) {
			return c.JSON(http.StatusNoContent, common.APIResponse{
				Message: "no job data",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error deleting job",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "job has been deleted",
		Data:    companyData,
	})

}
