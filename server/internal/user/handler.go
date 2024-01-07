package user

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"verve-hrms/internal/common"
	"verve-hrms/internal/schema"
)

type UserHandler struct {
	userService *UserService
}

func NewUserHandler(userService *UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (uh *UserHandler) GetUser(c echo.Context) error {
	user, ok := c.Get("user").(*jwt.Token) //echo handles missing/malformed token response
	if !ok {
		log.Printf("user.h.get_user: error asserting token")
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("user.h.get_user: error asserting claims: %v", user.Claims)
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "invalid claims data",
			Data:    nil,
		})
	}

	ID, ok := claims["id"].(float64)
	if !ok {
		log.Printf("user.h.get_user: error asserting id: %v", claims["id"])
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "id not found",
			Data:    nil,
		})
	}

	uintID := uint(ID)

	userData, err := uh.userService.GetUserByID(uintID)
	if err != nil {
		log.Printf("user.h.get_user: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "user data has been retrieved",
		Data:    userData,
	})
}

func (uh *UserHandler) EditUser(c echo.Context) error {
	user, ok := c.Get("user").(*jwt.Token) //echo handles missing/malformed token response
	if !ok {
		log.Printf("user.h.edit_user: error asserting token")
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("user.h.edit_user: error asserting claims: %v", user.Claims)
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "invalid claims data",
			Data:    nil,
		})
	}

	ID, ok := claims["id"].(float64)
	if !ok {
		log.Printf("user.h.edit_user: error asserting id: %v", claims["id"])
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "admin status not found",
			Data:    nil,
		})
	}

	uintID := uint(ID)

	var updateData schema.User
	err := c.Bind(&updateData)
	if err != nil {
		log.Printf("user.h.edit_user: error binding request %v", err)
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	updatedUser, err := uh.userService.UpdateUser(uintID, updateData)
	if err != nil {
		log.Printf("user.h.edit_user: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "user data has been updated",
		Data:    updatedUser,
	})
}
