package company

import (
	"github.com/labstack/echo/v4"
)

// SetupCompanyRoutes configures the routes for company management.
func SetupCompanyRoutes(g *echo.Group, companyHandler *CompanyHandler) {
	g.GET("/default", companyHandler.GetCompanyDataExpandDefault)
	g.GET("/:company_id", companyHandler.GetCompanyDataExpandID)

	g.POST("", companyHandler.CreateCompany)
	g.PUT("/:company_id", companyHandler.UpdateCompany)
	g.DELETE("/:company_id", companyHandler.DeleteCompany)

	g.POST("/:company_id/department", companyHandler.CreateDepartment)
	g.PUT("/:company_id/department/:department_id", companyHandler.UpdateDepartment)
	g.DELETE("/:company_id/department/:department_id", companyHandler.DeleteDepartment)

	g.POST("/:company_id/location", companyHandler.CreateLocation)
	g.PUT("/:company_id/location/:location_id", companyHandler.UpdateLocation)
	g.DELETE("/:company_id/location/:location_id", companyHandler.DeleteLocation)

	g.POST("/:company_id/job", companyHandler.CreateJob)
	g.PUT("/:company_id/job/:job_id", companyHandler.UpdateJob)
	g.DELETE("/:company_id/job/:job_id", companyHandler.DeleteJob)
}
