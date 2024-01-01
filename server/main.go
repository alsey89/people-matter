package main

import (
	"context"
	"log"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/spf13/viper"

	"verve-hrms/internal/auth"
	"verve-hrms/internal/user"
	"verve-hrms/setup"
)

func main() {
	e := echo.New()

	viper.SetConfigFile("dev.env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	client := setup.GetMongoClient()
	defer client.Disconnect(context.Background())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{viper.GetString("CLIENT_URL")},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
	}))
	e.Use(middleware.Logger())
	// e.Use(middleware.Secure()) //todo: enable in production
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "cookie:_csrf",
		CookiePath:     "/",
		CookieDomain:   "localhost",
		CookieSecure:   viper.GetBool("IS_PRODUCTION"),
		CookieHTTPOnly: true,
		// CookieSameSite: http.SameSiteStrictMode, //todo: enable in production
	}))

	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey:    []byte(viper.GetString("JWT_SECRET")),
		SigningMethod: "HS256",
		TokenLookup:   "cookie:jwt",
	})

	//* Instantiate User Domain
	userRepository := user.NewUserRepository(client)
	userService := user.NewUserService(userRepository)
	userHandler := user.NewUserHandler(userService)

	//* Instantiate Auth Domain
	authService := auth.NewAuthService(userService)
	authHandler := auth.NewAuthHandler(authService)

	authRoutes := e.Group("api/v1/auth")
	authRoutes.POST("/signin", authHandler.Signin)
	authRoutes.POST("/signup", authHandler.Signup)
	authRoutes.POST("/signout", authHandler.Signout)
	authRoutes.GET("/check", authHandler.CheckAuth, jwtMiddleware)

	userRoutes := e.Group("api/v1/user")
	userRoutes.GET("/data", userHandler.GetUser, jwtMiddleware)
	userRoutes.PUT("/data", userHandler.EditUser, jwtMiddleware)

	// Start the server
	e.Logger.Fatal(e.Start(":" + viper.GetString("SERVER_PORT")))
}
