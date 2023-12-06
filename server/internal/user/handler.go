package user

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"extesy-fullstack/internal/shared"
)

func GetCurrentUserHandler(c echo.Context, userRepository *UserRepository) error {
	user, ok := c.Get("user").(*jwt.Token) //echo handles missing/malformed token response
	if !ok {
		log.Printf("error asserting user")
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("error asserting claims: %v", user.Claims)
		return c.JSON(http.StatusBadRequest, shared.APIResponse{
			Message: "invalid claims data",
			Data:    nil,
		})
	}

	_id, ok := claims["id"].(string)
	if !ok {
		log.Printf("error asserting isAdmin: %v", claims["isAdmin"])
		return c.JSON(http.StatusBadRequest, shared.APIResponse{
			Message: "admin status not found",
			Data:    nil,
		})
	}

	objID, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		log.Printf("Error converting string to ObjectID: %v", err)
		return c.JSON(http.StatusInternalServerError, shared.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	existingUser, err := userRepository.Read(&objID)
	if err != nil {
		log.Printf("Error reading user data: %v", err)
		return c.JSON(http.StatusInternalServerError, shared.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, shared.APIResponse{
		Message: "user data has been retrieved",
		Data:    existingUser,
	})
}

func EditCurrentUserHandler(c echo.Context, userRepository *UserRepository) error {
	user, ok := c.Get("user").(*jwt.Token) //echo handles missing/malformed token response
	if !ok {
		log.Printf("error asserting user")
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("error asserting claims: %v", user.Claims)
		return c.JSON(http.StatusBadRequest, shared.APIResponse{
			Message: "invalid claims data",
			Data:    nil,
		})
	}

	_id, ok := claims["id"].(string)
	if !ok {
		log.Printf("error asserting id: %v", claims["id"])
		return c.JSON(http.StatusBadRequest, shared.APIResponse{
			Message: "admin status not found",
			Data:    nil,
		})
	}

	objID, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		log.Printf("Error converting string to ObjectID: %v", err)
		return c.JSON(http.StatusInternalServerError, shared.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	var updateData User
	err = c.Bind(&updateData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, shared.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	updatedUser, err := userRepository.Update(&objID, updateData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, shared.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, shared.APIResponse{
		Message: "user data has been updated",
		Data:    updatedUser,
	})
}
