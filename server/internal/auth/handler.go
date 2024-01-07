package auth

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"verve-hrms/internal/common"
)

type AuthHandler struct {
	authService *AuthService
}

func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Signup godoc
// @Summary User signup
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param creds body Credentials true "Signup Credentials"
// @Success 200 {object} common.APIResponse "user has been signed up and signed in"
// @Failure 400 {object} common.APIResponse "invalid email"
// @Failure 409 {object} common.APIResponse "email not available"
// @Failure 500 {object} common.APIResponse "something went wrong"
// @Router /auth/signup [post]
func (ah *AuthHandler) Signup(c echo.Context) error {
	creds := new(SignupCredentials)
	err := c.Bind(creds)
	if err != nil {
		log.Printf("error binding credentials: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	email := creds.Email
	if !common.EmailValidator(email) {
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "invalid email",
			Data:    nil,
		})
	}

	password := creds.Password
	confirmPassword := creds.ConfirmPassword

	if password != confirmPassword {
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "passwords do not match",
			Data:    nil,
		})
	}

	newUser, err := ah.authService.Signup(email, password)
	if err != nil {
		log.Printf("h.signup: %v", err)
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
		ID:      newUser.ID, // Store the ObjectId
		IsAdmin: newUser.IsAdmin,
		Email:   newUser.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		log.Printf("Error signing jwt with claims: %v", err)
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
		Message: "user has been signed up and signed in",
		Data:    newUser,
	})
}

// Signin godoc
// @Summary User signin
// @Description Authenticate a user and start a session
// @Tags auth
// @Accept json
// @Produce json
// @Param creds body Credentials true "Signin Credentials"
// @Success 200 {object} common.APIResponse "user has been signed in"
// @Failure 401 {object} common.APIResponse "invalid Credentials"
// @Failure 404 {object} common.APIResponse "user Not Found"
// @Failure 500 {object} common.APIResponse "internal Server Error"
// @Router /auth/signin [post]
func (ah *AuthHandler) Signin(c echo.Context) error {
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

	existingUser, err := ah.authService.Signin(email, password)
	if err != nil {
		log.Printf("h.signin: %v", err)
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
		ID:      existingUser.ID, // Store the ObjectId
		IsAdmin: existingUser.IsAdmin,
		Email:   existingUser.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		log.Printf("error signing jwt with claims: %v", err)
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

// Signout godoc
// @Summary User signout
// @Description End a user's session
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} common.APIResponse "user has been signed out"
// @Router /auth/signout [post]
func (ah *AuthHandler) Signout(c echo.Context) error {
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

// CheckAuth godoc
// @Summary Check authentication status
// @Description Check if the user is authenticated and if they are an admin
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} common.APIResponse "success"
// @Failure 400 {object} common.APIResponse "invalid claims data"
// @Failure 400 {object} common.APIResponse "admin status not found"
// @Router /auth/check [get]
func (ah *AuthHandler) CheckAuth(c echo.Context) error {
	user, ok := c.Get("user").(*jwt.Token) //echo handles missing/malformed token response
	if !ok {
		log.Printf("auth.check_auth: error asserting token")
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("auth.check_auth: error asserting claims: %v", user.Claims)
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "invalid claims data",
			Data:    nil,
		})
	}

	isAdmin, ok := claims["isAdmin"].(bool)
	if !ok {
		log.Printf("auth.check_auth: error asserting isAdmin: %v", claims["isAdmin"])
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "admin status not found",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "success",
		Data: echo.Map{
			"Authenticated": true,
			"IsAdmin":       isAdmin,
		},
	})
}
