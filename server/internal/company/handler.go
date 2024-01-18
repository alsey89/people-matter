package company

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"verve-hrms/internal/common"
	"verve-hrms/internal/schema"
)

type CompanyHandler struct {
	companyService *CompanyService
}

func NewCompanyHandler(companyService *CompanyService) *CompanyHandler {
	return &CompanyHandler{companyService: companyService}
}

//! Company ------------------------------------------------------------

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

	log.Printf("company.h.create_company: newCompany: %v", newCompany)

	companyData, err := ch.companyService.CreateNewCompanyAndReturnList(newCompany)
	if err != nil {
		log.Printf("company.h.create_company: %v", err)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.JSON(http.StatusNotFound, common.APIResponse{
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

func (ch *CompanyHandler) GetCompanyDataExpandDefault(c echo.Context) error {

	CompanyData, err := ch.companyService.GetCompanyListAndExpandDefault()
	if err != nil {
		log.Printf("company.h.get_company_data_expand_default: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, common.APIResponse{
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, common.APIResponse{
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

func (ch *CompanyHandler) UpdateCompany(c echo.Context) error {
	stringCompanyID := c.Param("company_id")

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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, common.APIResponse{
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, common.APIResponse{
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
			return c.JSON(http.StatusNotFound, common.APIResponse{
				Message: "department already exists",
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, common.APIResponse{
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, common.APIResponse{
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

//! Title ------------------------------------------------------------

func (ch *CompanyHandler) CreateTitle(c echo.Context) error {
	stringCompanyID := c.Param("company_id")
	if stringCompanyID == "" {
		log.Printf("company.h.create_title: empty company id parameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no company id",
			Data:    nil,
		})
	}

	companyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
	if err != nil {
		log.Printf("company.h.create_title: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing company id",
			Data:    nil,
		})
	}

	newTitle := new(schema.Title)

	err = c.Bind(newTitle)
	if err != nil {
		log.Printf("company.h.create_title: error binding title data: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	companyData, err := ch.companyService.CreateNewTitleAndReturnCompanyListAndExpandID(companyID, newTitle)
	if err != nil {
		log.Printf("company.h.create_title: %v", err)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.JSON(http.StatusConflict, common.APIResponse{
				Message: "title already exists",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error creating title data",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "title data has been created",
		Data:    companyData,
	})
}

func (ch *CompanyHandler) UpdateTitle(c echo.Context) error {
	stringCompanyID := c.Param("company_id")
	if stringCompanyID == "" {
		log.Printf("company.h.update_title: empty company id parameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no company id",
			Data:    nil,
		})
	}

	companyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
	if err != nil {
		log.Printf("company.h.update_title: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing company id",
			Data:    nil,
		})
	}

	stringTitleID := c.Param("title_id")
	if stringTitleID == "" {
		log.Printf("company.h.update_title: empty title id parameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no title id",
			Data:    nil,
		})
	}

	titleID, err := common.ConvertStringOfNumbersToUint(stringTitleID)
	if err != nil {
		log.Printf("company.h.update_title: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing title id",
			Data:    nil,
		})
	}

	dataToUpdate := new(schema.Title)
	err = c.Bind(dataToUpdate)
	if err != nil {
		log.Printf("company.h.update_title: error binding title data: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	companyData, err := ch.companyService.UpdateTitleAndReturnCompanyListAndExpandID(companyID, titleID, dataToUpdate)
	if err != nil {
		log.Printf("company.h.update_title: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, common.APIResponse{
				Message: "no title data",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error updating title data",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "title data has been updated",
		Data:    companyData,
	})
}

func (ch *CompanyHandler) DeleteTitle(c echo.Context) error {
	stringCompanyID := c.Param("company_id")
	if stringCompanyID == "" {
		log.Printf("company.h.delete_title: empty company id parameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no company id",
			Data:    nil,
		})
	}

	companyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
	if err != nil {
		log.Printf("company.h.delete_title: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing company id",
			Data:    nil,
		})
	}

	stringTitleID := c.Param("title_id")
	if stringTitleID == "" {
		log.Printf("company.h.delete_title: empty title id parameter")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no title id",
			Data:    nil,
		})
	}

	titleID, err := common.ConvertStringOfNumbersToUint(stringTitleID)
	if err != nil {
		log.Printf("company.h.update_title: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing title id",
			Data:    nil,
		})
	}

	companyData, err := ch.companyService.DeleteTitleAndReturnCompanyListAndExpandID(companyID, titleID)
	if err != nil {
		log.Printf("company.h.delete_title: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, common.APIResponse{
				Message: "no title data",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error deleting title",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "title has been deleted",
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

	companyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
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

	companyData, err := ch.companyService.CreateNewLocationAndReturnCompanyListAndExpandID(companyID, newLocation)
	if err != nil {
		log.Printf("company.h.create_location: %v", err)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.JSON(http.StatusConflict, common.APIResponse{
				Message: "location already exists",
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

	companyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
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

	locationID, err := common.ConvertStringOfNumbersToUint(stringLocationID)
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

	companyData, err := ch.companyService.UpdateLocationAndReturnCompanyListAndExpandID(companyID, locationID, dataToUpdate)
	if err != nil {
		log.Printf("company.h.update_location: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, common.APIResponse{
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

	companyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
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

	locationID, err := common.ConvertStringOfNumbersToUint(stringLocationID)
	if err != nil {
		log.Printf("company.h.delete_location: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error parsing location id",
			Data:    nil,
		})
	}

	companyData, err := ch.companyService.DeleteLocationAndReturnCompanyListAndExpandID(companyID, locationID)
	if err != nil {
		log.Printf("company.h.delete_location: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, common.APIResponse{
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
