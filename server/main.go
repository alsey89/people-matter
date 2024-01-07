package main

import (
	"log"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/spf13/viper"

	_ "verve-hrms/docs"
	"verve-hrms/internal/auth"
	"verve-hrms/internal/user"
	"verve-hrms/setup"
)

// @title Verve HRMS API
// @version 1.0
// @description This server provides APIs for the Verve HRMS application

// @contact.name alsey89
// @contact.email phyokyawsoe89@gmail.com

// @license.name GPL 3.0
// @license.url https://www.gnu.org/licenses/gpl-3.0.en.html

// @host localhost:3001
// @BasePath /api/v1

func main() {
	e := echo.New()

	viper.SetConfigFile("dev.env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	client := setup.GetClient()

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

	//swagger routes
	e.Static("/swagger", "docs")
	e.GET("/swagger/*", echoSwagger.WrapHandler)

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
