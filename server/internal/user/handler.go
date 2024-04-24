package user

import (
	"errors"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/alsey89/hrms/internal/common"
	"github.com/alsey89/hrms/internal/schema"
)

type UserHandler struct {
	userService *UserService
}

func NewUserHandler(userService *UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// ! Auth ------------------------------------------------------------
func (uh *UserHandler) GetCurrentUser(c echo.Context) error {
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

	userData, err := uh.userService.GetSingleUserByIDAndExpand(uintID)
	if err != nil {
		log.Printf("user.h.get_user: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting user data",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "user data has been retrieved",
		Data:    userData,
	})
}

// ! All Users List -----------------------------------------------------------
func (uh *UserHandler) GetUsersList(c echo.Context) error {
	userList, err := uh.userService.GetAllUsersWithRole()
	if err != nil {
		log.Printf("user.h.get_all_users: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting all user data",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "users data has been retrieved",
		Data:    userList,
	})
}

func (uh *UserHandler) CreateListUser(c echo.Context) error {
	var newUser schema.User
	err := c.Bind(&newUser)
	if err != nil {
		log.Printf("user.h.create_user: error binding request %v", err)
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	userList, err := uh.userService.CreateListUserAndGetAllUsersWithRole(&newUser)
	if err != nil {
		log.Printf("user.h.create_user: %v", err)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.JSON(http.StatusConflict, common.APIResponse{
				Message: "user already exists",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong while creating user",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, common.APIResponse{
		Message: "user has been created",
		Data:    userList,
	})
}

func (uh *UserHandler) UpdateListUser(c echo.Context) error {
	stringUserID := c.Param("user_id")
	if stringUserID == "" {
		log.Printf("user.h.update_user: user_id is empty")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "user_id is empty",
			Data:    nil,
		})
	}

	uintUserID, err := common.ConvertStringOfNumbersToUint(stringUserID)
	if err != nil {
		log.Printf("user.h.update_user: %v", err)
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "user_id is invalid format",
			Data:    nil,
		})
	}

	var updateData schema.User
	err = c.Bind(&updateData)
	if err != nil {
		log.Printf("user.h.update_user: error binding request %v", err)
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "error binding request",
			Data:    nil,
		})
	}

	userList, err := uh.userService.UpdateListUserAndGetAllUsersWithRole(uintUserID, updateData)
	if err != nil {
		log.Printf("user.h.update_user: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error updating user",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "user data has been updated",
		Data:    userList,
	})
}

func (uh *UserHandler) DeleteListUser(c echo.Context) error {
	stringUserID := c.Param("user_id")
	if stringUserID == "" {
		log.Printf("user.h.delete_user: user_id is empty")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "user_id is empty",
			Data:    nil,
		})
	}

	uintUserID, err := common.ConvertStringOfNumbersToUint(stringUserID)
	if err != nil {
		log.Printf("user.h.delete_user: %v", err)
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "user_id is invalid format",
			Data:    nil,
		})
	}

	userList, err := uh.userService.DeleteListUserAndGetAllUsersWithRole(uintUserID)
	if err != nil {
		log.Printf("user.h.delete_user: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, common.APIResponse{
				Message: "user not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error deleting user",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "user has been deleted",
		Data:    userList,
	})
}

// ! Single User Details -----------------------------------------------------------
func (uh *UserHandler) GetUserDetails(c echo.Context) error {
	stringUserID := c.Param("user_id")
	if stringUserID == "" {
		log.Printf("user.h.get_user_by_id: user_id is empty")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "user_id is empty",
			Data:    nil,
		})
	}

	uintUserID, err := common.ConvertStringOfNumbersToUint(stringUserID)
	if err != nil {
		log.Printf("user.h.get_user_by_id: %v", err)
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "user_id is invalid format",
			Data:    nil,
		})
	}

	userData, err := uh.userService.GetSingleUserByIDAndExpand(uintUserID)
	if err != nil {
		log.Printf("user.h.get_user_by_id: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error getting user by Id",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "user data has been retrieved",
		Data:    userData,
	})
}
