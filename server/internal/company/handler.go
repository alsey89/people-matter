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
		log.Printf("company.h.fetch_company_data_expand_default: %v", err)
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
		log.Printf("company.h.get_company_data: %v", err)
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
	log.Println("company.h.update_company: not implemented")
	return nil
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

	companyData, err := ch.companyService.DeleteCompanyAndReturnListAndExpandDefault(companyID)
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
