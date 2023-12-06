package main

import (
	"context"
	"log"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/spf13/viper"

	"extesy-fullstack/internal/auth"
	"extesy-fullstack/internal/user"
	"extesy-fullstack/setup"
)

func init() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}

func main() {
	e := echo.New()

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

	//* Instantiate Repositories
	userRepository := user.NewUserRepository(client)

	//todo: inject repository to service, service to handler

	authRoutes := e.Group("api/v1/auth")
	authRoutes.POST("/signin", func(c echo.Context) error {
		return auth.SigninHandler(c, userRepository)
	})
	authRoutes.POST("/signup", func(c echo.Context) error {
		return auth.SignupHandler(c, userRepository)
	})
	authRoutes.POST("/signout", func(c echo.Context) error {
		return auth.SignoutHandler(c)
	})
	authRoutes.GET("/check", func(c echo.Context) error {
		return auth.CheckAuthHandler(c)
	}, jwtMiddleware)

	userRoutes := e.Group("api/v1/user")
	userRoutes.GET("/data", func(c echo.Context) error {
		return user.GetCurrentUserHandler(c, userRepository)
	}, jwtMiddleware)
	userRoutes.PUT("/data", func(c echo.Context) error {
		return user.EditCurrentUserHandler(c, userRepository)
	}, jwtMiddleware)

	// Start the server
	e.Logger.Fatal(e.Start(":" + viper.GetString("SERVER_PORT")))
}
