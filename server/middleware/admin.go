package middleware

import (
	"log"
	"net/http"
	"verve-hrms/internal/common"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// AdminMiddleware is a middleware that checks jwt claims to see if the user is an admin
func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*jwt.Token) //echo handles missing/malformed token response
		if !ok {
			log.Printf("middleware.admin: error asserting token")
		}

		claims, ok := user.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("middleware.admin: error asserting claims: %v", user.Claims)
			return c.JSON(http.StatusBadRequest, common.APIResponse{
				Message: "invalid claims data",
				Data:    nil,
			})
		}

		isAdmin, ok := claims["admin"].(bool)
		if !ok {
			log.Printf("middleware.admin: error asserting admin status: %v", claims["admin"])
			return c.JSON(http.StatusBadRequest, common.APIResponse{
				Message: "admin status not found",
				Data:    nil,
			})
		}

		if !isAdmin {
			log.Printf("middleware.admin: user is not an admin")
			return c.JSON(http.StatusUnauthorized, common.APIResponse{
				Message: "user is not an admin",
				Data:    nil,
			})
		}

		return next(c)
	}
}
