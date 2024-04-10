package server

import (
	"context"
	"fmt"
	"log"
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

func InitiateModule(scope string) fx.Option {
	return fx.Module(
		scope,
		fx.Provide(func(p Params) *HTTPServer {
			logger := p.Logger.Named(scope)
			server := echo.New()
			config := loadConfig(scope)

			s := &HTTPServer{
				logger: logger,
				scope:  scope,

				config: config,
				server: server,
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
	//set defaults
	viper.SetDefault(fmt.Sprintf("%s.host", scope), DefaultHost)
	viper.SetDefault(fmt.Sprintf("%s.port", scope), DefaultPort)
	viper.SetDefault(fmt.Sprintf("%s.log_level", scope), DefaultLogLevel)

	viper.SetDefault(fmt.Sprintf("%s.allow_origins", scope), "*")
	viper.SetDefault(fmt.Sprintf("%s.allow_methods", scope), "GET,PUT,POST,DELETE")
	viper.SetDefault(fmt.Sprintf("%s.allow_headers", scope), "Origin,Content-Type,Accept")

	getConfigPath := func(key string) string {
		return fmt.Sprintf("%s.%s", scope, key)
	}

	return &Config{
		Host:         viper.GetString(getConfigPath("host")),
		Port:         viper.GetInt(getConfigPath("port")),
		LogLevel:     viper.GetString(getConfigPath("log_level")),
		AllowOrigins: viper.GetString(getConfigPath("allow_origins")),
		AllowMethods: viper.GetString(getConfigPath("allow_methods")),
		AllowHeaders: viper.GetString(getConfigPath("allow_headers")),
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
				log.Printf("|--------------------------------------------\n")
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
				log.Printf("--------------------------------------------|\n")
			case "PROD":
				log.Printf("|--------------------------------------------\n")
				s.logger.Info("request",
					zap.String("URI", v.URI),
					zap.Int("status", v.Status),
					zap.Any("error", v.Error),
					zap.String("request_id", v.RequestID),
					zap.Duration("latency", v.Latency),
				)
				log.Printf("--------------------------------------------|\n")
			case "DEBUG":
				log.Printf("|--------------------------------------------\n")
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
				log.Printf("--------------------------------------------|\n")
			default:
				s.logger.Error("invalid log level", zap.String("log_level", s.config.LogLevel))
			}
			return nil
		},
	}
	// add request logger middleware
	s.server.Use(middleware.RequestLoggerWithConfig(requestLoggerConfig))
}

func (s *HTTPServer) GetServer() *echo.Echo {
	return s.server
}
