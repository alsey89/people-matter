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

// func (d *Domain) SignupHandler(c echo.Context) error {

// 	creds := new(RootUserSignupCredentials)
// 	err := c.Bind(creds)
// 	if err != nil {
// 		d.logger.Error("[signupHandler] error binding credentials", zap.Error(err))
// 		return c.JSON(http.StatusInternalServerError, common.APIResponse{
// 			Message: "credentials error",
// 			Data:    nil,
// 		})
// 	}

// 	// validate email
// 	email := creds.Email
// 	if !common.EmailValidator(email) {
// 		d.logger.Error("[signupHandler] email validation failed")
// 		return c.JSON(http.StatusBadRequest, common.APIResponse{
// 			Message: "invalid email",
// 			Data:    nil,
// 		})
// 	}

// 	// validate password
// 	password := creds.Password
// 	confirmPassword := creds.ConfirmPassword
// 	if password != confirmPassword {
// 		d.logger.Error("[signupHandler] password confirmation check failed")
// 		return c.JSON(http.StatusBadRequest, common.APIResponse{
// 			Message: "passwords do not match",
// 			Data:    nil,
// 		})
// 	}

// 	// validate company
// 	companyName := creds.CompanyName
// 	if companyName == "" {
// 		d.logger.Error("[signupHandler] company name is empty")
// 		return c.JSON(http.StatusBadRequest, common.APIResponse{
// 			Message: "company name is required",
// 			Data:    nil,
// 		})
// 	}

// 	// todo: create company
// 	companyID := func() uint {
// 		return 1
// 	}()

// 	// create user
// 	newUser, err := d.SignupService(email, password, companyID)
// 	if err != nil {
// 		d.logger.Error("[signupHandler] error signing up user", zap.Error(err))
// 		if errors.Is(err, gorm.ErrDuplicatedKey) {
// 			return c.JSON(http.StatusConflict, common.APIResponse{
// 				Message: "email not available",
// 				Data:    nil,
// 			})
// 		}
// 		return c.JSON(http.StatusInternalServerError, common.APIResponse{
// 			Message: "something went wrong",
// 			Data:    nil,
// 		})
// 	}

// 	claims := Claims{
// 		ID:        newUser.ID, // Store the ObjectId
// 		CompanyID: newUser.CompanyID,
// 		Role:      newUser.Role,
// 		Email:     newUser.Email,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	t, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))
// 	if err != nil {
// 		d.logger.Error("[signupHandler] error signing jwt with claims", zap.Error(err))
// 		return c.JSON(http.StatusInternalServerError, common.APIResponse{
// 			Message: "token error",
// 			Data:    nil,
// 		})
// 	}

// 	cookie := new(http.Cookie)
// 	cookie.Name = "jwt"
// 	cookie.Value = t
// 	cookie.HttpOnly = true
// 	cookie.Secure = viper.GetBool("IS_PRODUCTION")
// 	cookie.Path = "/"
// 	cookie.Expires = time.Now().Add(time.Hour * 72)

// 	c.SetCookie(cookie)

// 	return c.JSON(http.StatusOK, common.APIResponse{
// 		Message: "user has been signed up and signed in",
// 		Data:    newUser,
// 	})
// }

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
	if err != nil {
		d.logger.Error("[SigninHandler]", zap.Error(err))
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
