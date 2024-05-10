package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/alsey89/people-matter/internal/common"
)

func (d *Domain) SigninHandler(c echo.Context) error {
	creds := new(SigninCredentials)
	err := c.Bind(creds)
	if err != nil {
		d.logger.Error("[SigninHandler] error binding credentials", zap.Error(err))
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "invalid form data",
			Data:    nil,
		})
	}

	email := creds.Email
	password := creds.Password

	existingUser, err := d.SignIn(email, password)
	switch {
	case err != nil:
		d.logger.Error("[SigninHandler]", zap.Error(err))
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return c.JSON(http.StatusNotFound, common.APIResponse{
				Message: "user not found",
				Data:    nil,
			})
		case errors.Is(err, ErrUserNotConfirmed):
			return c.JSON(http.StatusForbidden, common.APIResponse{
				Message: "user not confirmed",
				Data:    nil,
			})
		case errors.Is(err, ErrInvalidCredentials):
			return c.JSON(http.StatusUnauthorized, common.APIResponse{
				Message: "invalid credentials",
				Data:    nil,
			})
		default:
			return c.JSON(http.StatusInternalServerError, common.APIResponse{
				Message: "something went wrong",
				Data:    nil,
			})
		}
	}

	claims := Claims{
		ID:        existingUser.ID,
		CompanyID: existingUser.CompanyID,
		Role:      existingUser.Role,
		Email:     existingUser.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// include location id for manager role
	if existingUser.Role == "manager" {
		claims.LocationID = &existingUser.UserPosition.Location.ID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		d.logger.Error("[SigninHandler] error signing jwt with claims", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = t
	cookie.HttpOnly = true
	cookie.Secure = viper.GetBool("IS_PRODUCTION")
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(time.Hour * 72)

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "user has been signed in",
		Data:    existingUser,
	})
}

func (d *Domain) SignoutHandler(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = ""
	cookie.HttpOnly = true
	cookie.Secure = viper.GetBool("IS_PRODUCTION")
	cookie.Path = "/"
	cookie.Expires = time.Unix(0, 0) //* set the cookie to expire immediately

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "user has been signed out",
		Data:    nil,
	})
}

func (d *Domain) ConfirmationHandler(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		d.logger.Error("[ConfirmationHandler] token is empty")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "token is empty",
			Data:    nil,
		})
	}

	claims := &Claims{}
	parsedClaims, err := d.ParseToken(token, claims)
	if err != nil {
		d.logger.Error("[ConfirmationHandler] error parsing token", zap.Error(err))
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "invalid token",
			Data:    nil,
		})
	}

	err = d.ConfirmEmail(parsedClaims.ID, parsedClaims.CompanyID)
	if err != nil {
		d.logger.Error("[ConfirmationHandler]", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "user has been confirmed",
		Data:    nil,
	})
}

func (d *Domain) CheckAuth(c echo.Context) error {
	_, ok := c.Get("user").(*jwt.Token) //echo jwt middleware handles missing/malformed token response
	if !ok {
		d.logger.Error("[CheckAuth] error getting user token")
		return c.JSON(http.StatusUnauthorized, common.APIResponse{
			Message: "token error",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "authenticated",
		Data:    nil,
	})
}

func (d *Domain) GetCSRFToken(c echo.Context) error {
	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "csrf token has been set",
		Data:    nil,
	})
}
