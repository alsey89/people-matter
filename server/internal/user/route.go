package user

import (
	"github.com/labstack/echo/v4"
)

// SetupCompanyRoutes configures the routes for company management.
func SetupUserRoutes(g *echo.Group, userHandler *UserHandler) {
	g.GET("/current", userHandler.GetCurrentUser)

	// ! List
	g.GET("/list", userHandler.GetUsersList)
	g.POST("/list", userHandler.CreateListUser)
	g.PUT("/list/:user_id", userHandler.UpdateListUser)
	g.DELETE("/list/:user_id", userHandler.DeleteListUser)
	// ! Single User Details
	g.GET("/:user_id", userHandler.GetUserDetails)
}
