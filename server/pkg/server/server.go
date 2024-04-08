package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	DefaultHost     = "0.0.0.0"
	DefaultPort     = 3001
	DefaultLogLevel = "DEV"
)

type Config struct {
	Host         string
	Port         int
	LogLevel     string
	AllowOrigins string
	AllowMethods string
	AllowHeaders string
}

type HTTPServer struct {
	logger *zap.Logger
	config *Config

	scope  string
	server *echo.Echo
}

type Params struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
}

func Module(scope string) fx.Option {
	return fx.Module(
		scope,
		fx.Provide(func(p Params) *HTTPServer {
			logger := p.Logger.Named(scope)
			server := echo.New()
			config := loadConfig(scope)

			s := &HTTPServer{
				logger: logger,
				server: server,
				config: config,
				scope:  scope,
			}

			return s
		}),
		fx.Invoke(func(s *HTTPServer, p Params) {
			p.Lifecycle.Append(
				fx.Hook{
					OnStart: s.onStart,
					OnStop:  s.onStop,
				},
			)
		}),
	)
}

func loadConfig(scope string) *Config {
	getConfigWithDefault := func(key string, defaultVal interface{}) interface{} {
		scopedKey := fmt.Sprintf("%s.%s", scope, key)
		if viper.IsSet(scopedKey) {
			return viper.Get(scopedKey)
		}
		return defaultVal
	}

	return &Config{
		Host:         getConfigWithDefault("host", DefaultHost).(string),
		Port:         getConfigWithDefault("port", DefaultPort).(int),
		LogLevel:     getConfigWithDefault("log_level", DefaultLogLevel).(string),
		AllowOrigins: getConfigWithDefault("allow_origins", "*").(string),
		AllowMethods: getConfigWithDefault("allow_methods", "GET,PUT,POST,DELETE").(string),
		AllowHeaders: getConfigWithDefault("allow_headers", "Origin,Content-Type,Accept").(string),
	}
}

func (s *HTTPServer) onStart(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	s.logger.Info("Starting HTTPServer", zap.String("address", addr), zap.String("log_level", s.config.LogLevel))

	s.configureCORS()
	s.setUpRequestLogger()
	s.server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})

	go func() {
		if err := s.server.Start(addr); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal(err.Error())
		}
	}()

	return nil
}

func (s *HTTPServer) onStop(context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Error("server shutdown error", zap.Error(err))
	}

	s.logger.Info("server stopped")
	return nil
}

func (s *HTTPServer) configureCORS() {
	// configure CORS middleware
	corsConfig := middleware.CORSConfig{
		AllowOrigins: strings.Split(s.config.AllowOrigins, ","),
		AllowMethods: strings.Split(s.config.AllowMethods, ","),
		AllowHeaders: strings.Split(s.config.AllowHeaders, ","),
	}
	if s.config.AllowOrigins == "" {
		corsConfig.AllowOrigins = []string{"*"}
	}
	if s.config.AllowMethods == "" {
		corsConfig.AllowMethods = []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete}
	}
	if s.config.AllowHeaders == "" {
		corsConfig.AllowHeaders = []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept}
	}
	// add CORS middleware
	s.server.Use(middleware.CORSWithConfig(corsConfig))
}

// set up Echo RequestLogger with zap based on log level
func (s *HTTPServer) setUpRequestLogger() {

	// configure request logger according to log level
	requestLoggerConfig := middleware.RequestLoggerConfig{
		LogProtocol:  true,
		LogMethod:    true,
		LogURI:       true,
		LogStatus:    true,
		LogRequestID: true,
		LogRemoteIP:  true,
		LogLatency:   true,
		LogError:     true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			switch s.config.LogLevel {
			case "DEV":
				s.logger.Info("log level is set to DEV")
				s.logger.Info("request",
					zap.String("URI", v.URI),
					zap.String("method", v.Method),
					zap.Int("status", v.Status),
					zap.Any("error", v.Error),
					zap.String("remote_ip", v.RemoteIP),
					zap.String("request_id", v.RequestID),
					zap.Duration("latency", v.Latency),
					zap.String("protocol", v.Protocol),
				)
			case "PROD":
				s.logger.Info("log level is set to PROD")
				s.logger.Info("request",
					zap.String("URI", v.URI),
					zap.Int("status", v.Status),
					zap.Any("error", v.Error),
					zap.String("request_id", v.RequestID),
					zap.Duration("latency", v.Latency),
				)
			case "DEBUG":
				s.logger.Info("log level is set to DEBUG")
				s.logger.Debug("request",
					zap.String("URI", v.URI),
					zap.String("method", v.Method),
					zap.Int("status", v.Status),
					zap.String("remote_ip", v.RemoteIP),
					zap.String("request_id", v.RequestID),
					zap.Duration("latency", v.Latency),
					zap.String("protocol", v.Protocol),
					zap.Any("error", v.Error),
					zap.Any("request_body", c.Request().Body),
					// todo: add more debug logs if needed
				)
			default:
				s.logger.Error("invalid log level", zap.String("log_level", s.config.LogLevel))
			}
			return nil
		},
	}
	// add request logger middleware
	s.server.Use(middleware.RequestLoggerWithConfig(requestLoggerConfig))
}

func (s *HTTPServer) GetRouter() *echo.Echo {
	return s.server
}
