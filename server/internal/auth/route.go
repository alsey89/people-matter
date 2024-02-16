package auth

import (
	"github.com/labstack/echo/v4"
)

// SetupAuthRoutes configures the routes for authentication.
func SetupAuthRoutes(g *echo.Group, authHandler *AuthHandler) {
	g.POST("/signin", authHandler.Signin)
	g.POST("/signup", authHandler.Signup)
	g.POST("/signout", authHandler.Signout)
	g.GET("/check", authHandler.CheckAuth)
	g.GET("/csrf", authHandler.GetCSRFToken)
}
