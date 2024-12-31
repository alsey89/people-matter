package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alsey89/people-matter/pkg/util"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Module struct {
	config *Config
	logger *zap.Logger
	scope  string
	server *echo.Echo
}

type Params struct {
	fx.In
	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
}

type Config struct {
	AllowHeaders string
	AllowMethods string
	AllowOrigins string

	CSRFProtection bool
	CSRFSecure     bool
	CSRFDomain     string

	Host           string
	Port           int
	ServerLogLevel string
}

const (
	DefaultAllowHeaders = "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, X-Requested-With, Origin, Cache-Control, Pragma, Expires, Set-Cookie, Cookie, jwt"
	DefaultAllowMethods = "GET, PUT, POST, DELETE, OPTIONS, PATCH"
	DefaultAllowOrigins = "*"

	DefaultCSRFProtection = false
	DefaultCSRFSecure     = false
	DefaultCSRFDomain     = "localhost"

	DefaultHost           = "localhost"
	DefaultPort           = 3001
	DefaultServerLogLevel = "PROD"
)

// Custom validator for Echo using go-playground/validator.
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

//! MODULE ---------------------------------------------------------------

// Provides the Module struct to the fx framework, and registers lifecycle hooks.
func InjectModule(scope string) fx.Option {
	return fx.Module(
		scope,
		fx.Provide(func(p Params) *Module {

			m := &Module{
				scope: scope,
			}
			m.config = m.setupConfig(scope)
			m.logger = m.setupLogger(scope, p)
			m.server = m.setupServer()

			return m
		}),
		fx.Invoke(func(m *Module, p Params) {
			p.Lifecycle.Append(fx.Hook{
				OnStart: m.onStart,
				OnStop:  m.onStop,
			})
		}),
	)
}

// Instantiates new Module without using the fx framework.
func NewServer(scope string, logger *zap.Logger) *Module {
	m := &Module{
		scope: scope,
	}
	m.logger = logger.Named("[" + scope + "]")
	m.config = m.setupConfig(scope)
	m.server = m.setupServer()

	m.onStart(context.Background())

	return m
}

// ! INTERNAL ---------------------------------------------------------------

func (m *Module) setupConfig(scope string) *Config {
	// searches for pattern: "scope.key"
	viper.SetDefault(util.GetConfigPath(scope, "allow_headers"), DefaultAllowHeaders)
	viper.SetDefault(util.GetConfigPath(scope, "allow_methods"), DefaultAllowMethods)
	viper.SetDefault(util.GetConfigPath(scope, "allow_origins"), DefaultAllowOrigins)

	viper.SetDefault(util.GetConfigPath(scope, "csrf_protection"), DefaultCSRFProtection)
	viper.SetDefault(util.GetConfigPath(scope, "csrf_secure"), DefaultCSRFSecure)
	viper.SetDefault(util.GetConfigPath(scope, "csrf_domain"), DefaultCSRFDomain)

	viper.SetDefault(util.GetConfigPath(scope, "host"), DefaultHost)
	viper.SetDefault(util.GetConfigPath(scope, "port"), DefaultPort)

	return &Config{
		AllowHeaders: viper.GetString(util.GetConfigPath(scope, "allow_headers")),
		AllowMethods: viper.GetString(util.GetConfigPath(scope, "allow_methods")),
		AllowOrigins: viper.GetString(util.GetConfigPath(scope, "allow_origins")),

		CSRFProtection: viper.GetBool(util.GetConfigPath(scope, "csrf_protection")),
		CSRFSecure:     viper.GetBool(util.GetConfigPath(scope, "csrf_secure")),
		CSRFDomain:     viper.GetString(util.GetConfigPath(scope, "csrf_domain")),

		Host:           viper.GetString(util.GetConfigPath(scope, "host")),
		Port:           viper.GetInt(util.GetConfigPath(scope, "port")),
		ServerLogLevel: viper.GetString(util.GetConfigPath("global", "log_level")),
	}
}

func (m *Module) setupLogger(scope string, p Params) *zap.Logger {
	logger := p.Logger.Named("[" + scope + "]")
	return logger
}

func (m *Module) setupServer() *echo.Echo {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	return e
}

func (m *Module) onStart(context.Context) error {
	m.logger.Info("Starting server")

	m.setUpCorsMiddleware()
	m.setUpCSRFMiddleware()

	if strings.EqualFold(m.config.ServerLogLevel, "debug") ||
		strings.EqualFold(m.config.ServerLogLevel, "dev") {
		m.setUpRequestLoggerMiddleware()
	}

	// server must be started in a goroutine to prevent blocking the hooks
	// context will timeout otherwise
	go m.startServer(true, false)

	if strings.EqualFold(m.config.ServerLogLevel, "debug") {
		m.logConfigurations()
	}

	return nil
}

func (m *Module) onStop(context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.server.Shutdown(ctx)
	if err != nil {
		m.logger.Error("server shutdown error", zap.Error(err))
	}

	m.logger.Info("Server module stopped")
	return nil
}

func (m *Module) setUpCorsMiddleware() {

	// Helper function to split and trim spaces
	splitAndTrim := func(s string) []string {
		parts := strings.Split(s, ",")
		for i, part := range parts {
			parts[i] = strings.TrimSpace(part)
		}
		return parts
	}

	corsConfig := middleware.CORSConfig{
		AllowOrigins:     splitAndTrim(m.config.AllowOrigins),
		AllowMethods:     splitAndTrim(m.config.AllowMethods),
		AllowHeaders:     splitAndTrim(m.config.AllowHeaders),
		AllowCredentials: true,
	}

	// Default values if unspecified
	if m.config.AllowOrigins == "" || m.config.AllowOrigins == "*" {
		corsConfig.AllowOrigins = []string{"*"}
	}
	if m.config.AllowMethods == "" || m.config.AllowMethods == "*" {
		corsConfig.AllowMethods = []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete}
	}
	if m.config.AllowHeaders == "" || m.config.AllowHeaders == "*" {
		corsConfig.AllowHeaders = []string{
			"accept",
			"content-type",
			"authorization",
			"x-csrf-token",
			"x-requested-with",
			"origin",
			"cache-control",
			"pragma",
			"expires",
			"set-cookie",
			"cookie",
			"X-Tenant", // tenant identifier
			"jwt",      // jwt token for authentication
			"token",    // confirmation tokens
		}
	}

	m.server.Use(middleware.CORSWithConfig(corsConfig))
}

func (m *Module) setUpCSRFMiddleware() {
	// defaults to not using CSRF protection if unspecified
	if !m.config.CSRFProtection {
		return
	}

	CSRFConfig := middleware.CSRFConfig{
		TokenLookup: "cookie:_csrf",
		CookiePath:  "/",
		// CookieDomain:   m.config.CSRFDomain,
		CookieSecure:   m.config.CSRFSecure,
		CookieSameSite: http.SameSiteNoneMode,
		CookieHTTPOnly: true,
		Skipper: func(c echo.Context) bool {
			// skip CSRF for OPTIONS requests
			if c.Request().Method == http.MethodOptions {
				return true
			}
			// skip CSRF for signout
			if c.Request().URL.Path == "/api/v1/auth/signout" {
				return true
			}
			// skip CSRF for "/api/v1/memorial/:memorialID/app/**"
			if strings.Contains(c.Request().URL.Path, "/api/v1/memorial/") && strings.Contains(c.Request().URL.Path, "/app/") {
				return true
			}

			return false
		},
	}
	m.server.Use(middleware.CSRFWithConfig(CSRFConfig))
}

func (m *Module) setUpRequestLoggerMiddleware() {
	// Defaults to PROD log level if unspecified
	// Valid log levels: DEV, PROD, DEBUG
	requestLoggerConfig := middleware.RequestLoggerConfig{
		LogProtocol:   true,
		LogMethod:     true,
		LogURI:        true,
		LogStatus:     true,
		LogRequestID:  true,
		LogRemoteIP:   true,
		LogLatency:    true,
		LogError:      true,
		LogValuesFunc: m.logRequest,
	}

	m.server.Use(middleware.RequestLoggerWithConfig(requestLoggerConfig))
}

// helper function for setUpRequestLoggerMiddleware
func (m *Module) logRequest(c echo.Context, v middleware.RequestLoggerValues) error {

	lowerCaseLogLevel := strings.ToLower(m.config.ServerLogLevel)

	switch lowerCaseLogLevel {
	case "dev":
		m.logger.Info("request",
			zap.String("method", v.Method),
			zap.String("URI", v.URI),
			zap.String("request_id", v.RequestID),
			zap.Int("status", v.Status),
			zap.Any("error", v.Error),
		)
	case "prod":
		m.logger.Info("request",
			zap.String("URI", v.URI),
			zap.Int("status", v.Status),
			zap.Any("error", v.Error),
			zap.String("request_id", v.RequestID),
			zap.Duration("latency", v.Latency),
		)
	case "debug":
		m.logger.Debug("request",
			zap.String("method", v.Method),
			zap.String("URI", v.URI),
			zap.Int("status", v.Status),
			zap.String("remote_ip", v.RemoteIP),
			zap.String("request_id", v.RequestID),
			zap.Duration("latency", v.Latency),
			zap.String("protocol", v.Protocol),
			zap.Any("error", v.Error),
		)
	default:
		m.logger.Error("invalid log level", zap.String("log_level", m.config.ServerLogLevel))
		return fmt.Errorf("invalid log level: %s", m.config.ServerLogLevel)
	}

	return nil
}

func (m *Module) startServer(HideBanner bool, HidePort bool) {
	m.server.HideBanner = HideBanner || false
	m.server.HidePort = HidePort || false

	addr := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)

	log.Printf("Server started at %s", addr)

	err := m.server.Start(addr)
	if err != nil && err != http.ErrServerClosed {
		m.logger.Fatal(err.Error())
	}
}

func (m *Module) logConfigurations() {
	m.logger.Debug("----- Server Configuration -----")
	m.logger.Debug("Host", zap.String("Host", m.config.Host))
	m.logger.Debug("Port", zap.Int("Port", m.config.Port))

	m.logger.Debug("----- Cors Configuration -----")
	m.logger.Debug("AllowOrigins", zap.String("AllowOrigins", m.config.AllowOrigins))
	m.logger.Debug("AllowMethods", zap.String("AllowMethods", m.config.AllowMethods))
	m.logger.Debug("AllowHeaders", zap.String("AllowHeaders", m.config.AllowHeaders))

	m.logger.Debug("----- CSRF Configuration -----")
	m.logger.Debug("CSRFProtection", zap.Bool("CSRFProtection", m.config.CSRFProtection))
	if m.config.CSRFProtection {
		m.logger.Debug("CSRFTokenLookup", zap.String("CSRFTokenLookup", "cookie:_csrf"))
		m.logger.Debug("CSRFCookiePath", zap.String("CSRFCookiePath", "/"))
		m.logger.Debug("CSRFCookieDomain", zap.String("CSRFCookieDomain", m.config.CSRFDomain))
		m.logger.Debug("CSRFCookieSecure", zap.Bool("CSRFCookieSecure", m.config.CSRFSecure))
		m.logger.Debug("CSRFCookieSameSite", zap.String("CSRFCookieSameSite", "Default"))
		m.logger.Debug("CSRFCookieHTTPOnly", zap.Bool("CSRFCookieHTTPOnly", true))
	}
}

//! EXTERNAL ---------------------------------------------------------------

// Returns the echo server instance
func (m *Module) GetServer() *echo.Echo {
	return m.server
}
