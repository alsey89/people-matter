package common

import (
	"fmt"
	"net/mail"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// validate email address
func EmailValidator(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// get current date and time in string format
func GetCurrentDateTimeString() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

// format time to string
func ConvertTimeToString(t time.Time) string {
	return t.Format("2006/01/02 15:04:05")
}

// take string of numbers and convert to uint
func ConvertStringOfNumbersToUint(str string) (uint, error) {
	if str == "" {
		return 0, fmt.Errorf("[common.ConvertStringOfNumbersToUint] string is empty")
	}
	num, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("[common.ConvertStringOfNumbersToUint] error parsing string to uint: %w", err)
	}
	uintNum := uint(num)
	return uintNum, nil
}

// get uint userID from JWT token
func GetUserIDFromToken(c echo.Context) (*uint, error) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return nil, fmt.Errorf("[common.GetCompanyIDFromToken] error asserting token. It seems to be of type: %T", c.Get("user"))
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("[common.GetCompanyIDFromToken] error asserting claims: %v. It seems to be of type: %T", user.Claims, user.Claims)
	}

	ID, ok := claims["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("[common.GetUserIDFromToken] error asserting ID: %v. It seems to be of type: %T", claims["id"], claims["id"])
	}

	uintID := uint(ID)

	return &uintID, nil
}

// get uint companyID from JWT token
func GetCompanyIDFromToken(c echo.Context) (*uint, error) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return nil, fmt.Errorf("[common.GetCompanyIDFromToken] error asserting token. It seems to be of type: %T", c.Get("user"))
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("[common.GetCompanyIDFromToken] error asserting claims: %v. It seems to be of type: %T", user.Claims, user.Claims)
	}

	ID, ok := claims["companyId"].(float64)
	if !ok {
		return nil, fmt.Errorf("[common.GetCompanyIDFromToken] error asserting ID: %v. It seems to be of type: %T", claims["companyId"], claims["companyId"])
	}

	uintID := uint(ID)

	return &uintID, nil
}

// get uint ID from URL parameter
func GetIDFromParam(key string, c echo.Context) (*uint, error) {
	value := c.Param(key)
	if value == "" {
		return nil, fmt.Errorf("[common.GetValueFromParam] %s is empty", key)
	}

	uintValue, err := ConvertStringOfNumbersToUint(value)
	if err != nil {
		return nil, fmt.Errorf("[common.GetValueFromParam] %w", err)
	}

	return &uintValue, nil
}
