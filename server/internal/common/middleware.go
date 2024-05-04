package common

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// middleware to check if user is admin
func MustBeAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get user role from token
		role, err := GetRoleFromToken(c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, APIResponse{
				Message: "error getting role from token",
				Data:    nil,
			})
		}

		if role != "admin" {
			return c.JSON(http.StatusUnauthorized, APIResponse{
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
		role, err := GetRoleFromToken(c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, APIResponse{
				Message: "error getting role from token",
				Data:    nil,
			})
		}

		if role != "manager" && role != "admin" {
			return c.JSON(http.StatusUnauthorized, APIResponse{
				Message: "unauthorized",
				Data:    nil,
			})
		}

		return next(c)
	}
}

func GetRoleFromToken(c echo.Context) (string, error) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return "", errors.New("[common.GetRoleFromToken] error getting user from token")
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("[common.GetRoleFromToken] error getting claims from token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", errors.New("[common.GetRoleFromToken] error getting role from token")
	}

	return role, nil
}
