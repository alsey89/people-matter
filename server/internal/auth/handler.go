package auth

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/alsey89/people-matter/internal/common"
)

func (a *Domain) CreateAccountHandler(c echo.Context) error {

	creds := new(SignupCredentials)
	err := c.Bind(creds)
	if err != nil {
		a.logger.Error("[signupHandler] error binding credentials", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "credentials error",
			Data:    nil,
		})
	}

	// validate email
	email := creds.Email
	if !common.EmailValidator(email) {
		a.logger.Error("[signupHandler] email validation failed")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "invalid email",
			Data:    nil,
		})
	}

	// validate password
	password := creds.Password
	confirmPassword := creds.ConfirmPassword
	if password != confirmPassword {
		a.logger.Error("[signupHandler] password confirmation check failed")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "passwords do not match",
			Data:    nil,
		})
	}

	// validate company
	// companyID := c.Get("companyID").(uint)

	newUser, err := a.SignupService(email, password)
	if err != nil {
		a.logger.Error("[signupHandler] error signing up user", zap.Error(err))
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.JSON(http.StatusConflict, common.APIResponse{
				Message: "email not available",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	claims := Claims{
		ID:        newUser.ID, // Store the ObjectId
		CompanyID: newUser.CompanyID,
		Role:      newUser.Role,
		Email:     newUser.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		a.logger.Error("[signupHandler] error signing jwt with claims", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "token error",
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
		Message: "user has been signed up and signed in",
		Data:    newUser,
	})
}

func (a *Domain) SigninHandler(c echo.Context) error {
	creds := new(SigninCredentials)
	err := c.Bind(creds)
	if err != nil {
		log.Printf("auth.h.signin: error binding credentials: %v", err)
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "invalid form data",
			Data:    nil,
		})
	}

	email := creds.Email
	password := creds.Password

	existingUser, err := a.SigninService(email, password)
	if err != nil {
		log.Printf("auth.h.signin: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, common.APIResponse{
				Message: "user not found",
				Data:    nil,
			})
		} else if errors.Is(err, ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, common.APIResponse{
				Message: "invalid credentials",
				Data:    nil,
			})
		} else {
			return c.JSON(http.StatusInternalServerError, common.APIResponse{
				Message: "something went wrong",
				Data:    nil,
			})
		}
	}

	claims := Claims{
		ID:        existingUser.ID, // Store the ObjectId
		CompanyID: existingUser.CompanyID,
		Role:      existingUser.Role,
		Email:     existingUser.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		log.Printf("auth.h.signin: error signing jwt with claims: %v", err)
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

func (a *Domain) SignoutHandler(c echo.Context) error {
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

func (a *Domain) CheckAuth(c echo.Context) error {
	user, ok := c.Get("user").(*jwt.Token) //echo handles missing/malformed token response
	if !ok {
		log.Printf("auth.check_auth: error asserting token")
		return c.JSON(http.StatusUnauthorized, common.APIResponse{
			Message: "token error",
			Data:    nil,
		})
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("auth.check_auth: error asserting claims: %v", user.Claims)
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "invalid claims data",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "success",
		Data: echo.Map{
			"authenticated": true,
			"role":          claims["role"],
		},
	})
}

func (a *Domain) GetCSRFToken(c echo.Context) error {
	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "csrf token has been set",
		Data:    nil,
	})
}
