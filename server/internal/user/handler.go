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

// auth
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

	userData, err := uh.userService.GetUserByIDAndExpand(uintID)
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

// users
func (uh *UserHandler) GetCompanyUsers(c echo.Context) error {
	stringCompanyID := c.Param("company_id")
	if stringCompanyID == "" {
		log.Printf("user.h.get_company_users: company_id is empty")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "company_id is empty",
			Data:    nil,
		})
	}

	uintCompanyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
	if err != nil {
		log.Printf("user.h.get_company_users: %v", err)
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	users, err := uh.userService.GetAllUsersAndExpand(uintCompanyID)
	if err != nil {
		log.Printf("user.h.get_all_users: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "users data has been retrieved",
		Data:    users,
	})
}

func (uh *UserHandler) CreateCompanyUser(c echo.Context) error {
	stringCompanyID := c.Param("company_id")
	if stringCompanyID == "" {
		log.Printf("user.h.create_company_user: company_id is empty")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "company_id is empty",
			Data:    nil,
		})
	}

	uintCompanyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
	if err != nil {
		log.Printf("user.h.create_company_user: %v", err)
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	var newUser schema.User
	err = c.Bind(&newUser)
	if err != nil {
		log.Printf("user.h.create_company_user: error binding request %v", err)
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	userList, err := uh.userService.CreateNewUserAndGetAllUsersAndExpand(uintCompanyID, &newUser)
	if err != nil {
		log.Printf("user.h.create_company_user: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, common.APIResponse{
		Message: "user has been created",
		Data:    userList,
	})
}

func (uh *UserHandler) UpdateCompanyUser(c echo.Context) error {
	stringCompanyID := c.Param("company_id")
	if stringCompanyID == "" {
		log.Printf("user.h.update_compnay_user: company_id is empty")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "company_id is empty",
			Data:    nil,
		})
	}

	uintCompanyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
	if err != nil {
		log.Printf("user.h.update_compnay_user: %v", err)
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	stringUserID := c.Param("id")
	if stringUserID == "" {
		log.Printf("user.h.update_compnay_user: user_id is empty")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "user_id is empty",
			Data:    nil,
		})
	}

	uintUserID, err := common.ConvertStringOfNumbersToUint(stringUserID)
	if err != nil {
		log.Printf("user.h.update_compnay_user: %v", err)
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	var updateData schema.User
	err = c.Bind(&updateData)
	if err != nil {
		log.Printf("user.h.update_compnay_user: error binding request %v", err)
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	userList, err := uh.userService.UpdateUserAndGetAllUsersAndExpand(uintCompanyID, uintUserID, updateData)
	if err != nil {
		log.Printf("user.h.update_compnay_user: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "user data has been updated",
		Data:    userList,
	})
}

func (uh *UserHandler) DeleteCompanyUser(c echo.Context) error {
	stringCompanyID := c.Param("company_id")
	if stringCompanyID == "" {
		log.Printf("user.h.delete_company_user: company_id is empty")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "company_id is empty",
			Data:    nil,
		})
	}

	uintCompanyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
	if err != nil {
		log.Printf("user.h.delete_company_user: %v", err)
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "company_id is invalid format",
			Data:    nil,
		})
	}

	stringUserID := c.Param("id")
	if stringUserID == "" {
		log.Printf("user.h.delete_company_user: user_id is empty")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "user_id is empty",
			Data:    nil,
		})
	}

	uintUserID, err := common.ConvertStringOfNumbersToUint(stringUserID)
	if err != nil {
		log.Printf("user.h.delete_company_user: %v", err)
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "user_id is invalid format",
			Data:    nil,
		})
	}

	userList, err := uh.userService.DeleteUserAndGetAllUsersAndExpand(uintCompanyID, uintUserID)
	if err != nil {
		log.Printf("user.h.delete_user: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "user has been deleted",
		Data:    userList,
	})
}
