package main

import (
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/spf13/viper"

	_ "verve-hrms/docs"
	"verve-hrms/internal/auth"
	"verve-hrms/internal/company"
	"verve-hrms/internal/job"
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

	//! Config
	setup.InitConfig()

	//! Load DB
	client := setup.GetClient()
	//! Middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{viper.GetString("CLIENT_URL")},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
	}))
	if viper.GetBool("IS_PRODUCTION") {
		e.Use(middleware.Secure())
		e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
			TokenLookup:    "cookie:_csrf",
			CookiePath:     "/",
			CookieDomain:   viper.GetString("PRODUCTION_DOMAIN"),
			CookieSecure:   true,
			CookieHTTPOnly: true,
			CookieSameSite: http.SameSiteStrictMode,
		}))
	} else {
		// e.Use(middleware.Secure()) //! this is not needed for development
		e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
			TokenLookup:    "cookie:_csrf",
			CookiePath:     "/",
			CookieDomain:   "localhost",
			CookieSecure:   false,
			CookieHTTPOnly: true,
			CookieSameSite: http.SameSiteLaxMode,
		}))
	}
	e.Use(echojwt.WithConfig(echojwt.Config{
		Skipper: func(c echo.Context) bool {
			if c.Request().URL.Path == "/api/v1/auth/signin" ||
				c.Request().URL.Path == "/api/v1/auth/signup" ||
				c.Request().URL.Path == "/api/v1/auth/signout" ||
				c.Request().URL.Path == "/api/v1/auth/csrf" ||
				c.Request().URL.Path == "/api/v1/auth/password/reset" ||
				c.Request().URL.Path == "/swagger" {
				return true
			}
			return false
		},
		SigningKey:    []byte(viper.GetString("JWT_SECRET")),
		SigningMethod: "HS256",
		TokenLookup:   "cookie:jwt",
	}))
	e.Use(middleware.Logger())
	//! Swagger
	e.Static("/swagger", "docs")
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	//! Domains
	//* Instantiate User Domain
	userRepository := user.NewUserRepository(client)
	userService := user.NewUserService(userRepository)
	userHandler := user.NewUserHandler(userService)
	//* Instantiate Auth Domain
	authService := auth.NewAuthService(userService)
	authHandler := auth.NewAuthHandler(authService)
	//*Instantiate Job Domain
	jobRepository := job.NewJobRepository(client)
	// jobService := job.NewJobService(jobRepository)
	// jobHandler := job.NewJobHandler(jobService)
	//*Instantiate Company Domain
	companyRepository := company.NewCompanyRepository(client)
	companyService := company.NewCompanyService(companyRepository, jobRepository)
	companyHandler := company.NewCompanyHandler(companyService)

	//! Routes
	//* Auth Routes
	authRoutes := e.Group("api/v1/auth")
	auth.SetupAuthRoutes(authRoutes, authHandler)
	//* User Routes
	userRoutes := e.Group("api/v1/user")
	user.SetupUserRoutes(userRoutes, userHandler)
	//* Company Routes
	companyRoutes := e.Group("api/v1/company")
	company.SetupCompanyRoutes(companyRoutes, companyHandler)

	//! START THE SERVER
	e.Logger.Fatal(e.Start(":" + viper.GetString("SERVER_PORT")))
}
