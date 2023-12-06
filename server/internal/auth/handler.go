package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"extesy-fullstack/internal/shared"
	"extesy-fullstack/internal/user"
)

func SignupHandler(c echo.Context, userRepository *user.UserRepository) error {
	creds := new(Credentials)
	err := c.Bind(creds)
	if err != nil {
		return err
	}

	username := creds.Username
	if username == "" {
		username = "New User" // default username
	}
	email := creds.Email
	password := creds.Password

	emailIsAvailable, err := userRepository.CheckEmailAvailability(email) //* using availability over existence because of return type (bool)
	if !emailIsAvailable {
		return c.JSON(http.StatusConflict, shared.APIResponse{
			Message: "user already exists",
			Data:    nil,
		})
	}
	if err != nil {
		log.Printf("error checking email availability: %v", err)
		return c.JSON(http.StatusInternalServerError, shared.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "error hashing password",
		})
	}

	newUUID, err := uuid.NewUUID()
	if err != nil {
		log.Printf("Error generating UUID: %v", err)
		return c.JSON(http.StatusInternalServerError, shared.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	newUser := user.User{
		Username:  username,
		UserID:    newUUID,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newUser, err = userRepository.Create(newUser) //* this adds the ID to newUser
	if err != nil {
		log.Printf("Error inserting new user: %v", err)
		return c.JSON(http.StatusInternalServerError, shared.APIResponse{
			Message: err.Error(),
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
		return c.JSON(http.StatusInternalServerError, shared.APIResponse{
			Message: err.Error(),
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

	return c.JSON(http.StatusOK, shared.APIResponse{
		Message: "user has been signed up and signed in",
		Data:    newUser,
	})
}

func SigninHandler(c echo.Context, userRepository *user.UserRepository) error {
	creds := new(Credentials)
	err := c.Bind(creds)
	if err != nil {
		log.Printf("error binding credentials: %v", err)
		return c.JSON(http.StatusInternalServerError, shared.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	email := creds.Email
	password := creds.Password

	// Check if user exists
	existingUser, err := userRepository.ReadByEmail(email)
	if err == mongo.ErrNoDocuments {
		return c.JSON(http.StatusBadRequest, shared.APIResponse{
			Message: "user does not exist",
			Data:    nil,
		})
	} else if err != nil {
		log.Printf("error reading user by email: %v", err)
		return c.JSON(http.StatusBadRequest, shared.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password))
	if err != nil {
		log.Printf("error comparing password: %v", err)
		return c.JSON(http.StatusBadRequest, shared.APIResponse{
			Message: err.Error(),
			Data:    nil,
		})
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
		return c.JSON(http.StatusInternalServerError, shared.APIResponse{
			Message: err.Error(),
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

	return c.JSON(http.StatusOK, shared.APIResponse{
		Message: "user has been signed in",
		Data:    existingUser,
	})
}

func SignoutHandler(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = ""
	cookie.HttpOnly = true
	cookie.Secure = viper.GetBool("IS_PRODUCTION")
	cookie.Path = "/"
	cookie.Expires = time.Unix(0, 0) //* set the cookie to expire immediately

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, shared.APIResponse{
		Message: "user has been signed out",
		Data:    nil,
	})
}

func CheckAuthHandler(c echo.Context) error {
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

	isAdmin, ok := claims["isAdmin"].(bool)
	if !ok {
		log.Printf("error asserting isAdmin: %v", claims["isAdmin"])
		return c.JSON(http.StatusBadRequest, shared.APIResponse{
			Message: "admin status not found",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, shared.APIResponse{
		Message: "success",
		Data: echo.Map{
			"Authenticated": true,
			"IsAdmin":       isAdmin,
		},
	})
}
