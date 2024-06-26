package middleware

import (
	"errors"
	"net/http"

	"github.com/alsey89/people-matter/internal/common"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// middleware to check if user is admin
func MustBeAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get user role from token
		role, err := getRoleFromToken(c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.APIResponse{
				Message: "error getting role from token",
				Data:    nil,
			})
		}

		if role != "admin" && role != "root" {
			return c.JSON(http.StatusUnauthorized, common.APIResponse{
				Message: "unauthorized",
				Data:    nil,
			})
		}

		return next(c)
	}
}

// middleware to check if user is manager or above
func MustBeManager(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get user role from token
		role, err := getRoleFromToken(c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.APIResponse{
				Message: "error getting role from token",
				Data:    nil,
			})
		}

		if role != "manager" && role != "admin" && role != "root" {
			return c.JSON(http.StatusUnauthorized, common.APIResponse{
				Message: "unauthorized",
				Data:    nil,
			})
		}

		return next(c)
	}
}

// ------------------------------

func getRoleFromToken(c echo.Context) (string, error) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return "", errors.New("[common.GetroleFromToken] error getting user from token")
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("[common.GetroleFromToken] error getting claims from token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", errors.New("[common.GetroleFromToken] error getting role from token")
	}

	return role, nil
}
