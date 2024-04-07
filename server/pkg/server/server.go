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

type HTTPServer struct {
	logger *zap.Logger
	server *echo.Echo
	scope  string
}

type Params struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
}

func Module(scope string) fx.Option {
	var s *HTTPServer

	return fx.Module(
		scope,
		fx.Provide(func(p Params) *HTTPServer {
			logger := p.Logger.Named(scope)
			server := echo.New()

			s := &HTTPServer{
				server: server,
				logger: logger,
				scope:  scope,
			}

			s.initDefaultConfigs()

			return s
		}),
		fx.Populate(&s),
		fx.Invoke(func(p Params) {
			p.Lifecycle.Append(
				fx.Hook{
					OnStart: s.onStart,
					OnStop:  s.onStop,
				},
			)
		}),
	)
}

func (s *HTTPServer) getConfigPath(key string) string {
	return fmt.Sprintf("%s.%s", s.scope, key)
}

func (s *HTTPServer) initDefaultConfigs() {
	viper.SetDefault(s.getConfigPath("host"), DefaultHost)
	viper.SetDefault(s.getConfigPath("port"), DefaultPort)
	viper.SetDefault(s.getConfigPath("log_level"), DefaultLogLevel)
}

func (s *HTTPServer) onStart(ctx context.Context) error {
	//! Configurations ---------------------
	port := viper.GetInt(s.getConfigPath("port"))
	host := viper.GetString(s.getConfigPath("host"))
	addr := fmt.Sprintf("%s:%d", host, port)

	logLevel := viper.GetString(s.getConfigPath("log_level"))

	allowOrigins := viper.GetString(s.getConfigPath("allow_origins"))
	allowMethods := viper.GetString(s.getConfigPath("allow_methods"))
	allowHeaders := viper.GetString(s.getConfigPath("allow_headers"))

	s.logger.Info("Starting HTTPServer",
		zap.String("address", addr),
	)

	//! Setup CORS ---------------------
	corsConfig := middleware.CORSConfig{
		AllowOrigins: strings.Split(allowOrigins, ","),
		AllowMethods: strings.Split(allowMethods, ","),
		AllowHeaders: strings.Split(allowHeaders, ","),
	}
	if allowOrigins == "" {
		corsConfig.AllowOrigins = []string{"*"}
	}
	if allowMethods == "" {
		corsConfig.AllowMethods = []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete}
	}
	if allowHeaders == "" {
		corsConfig.AllowHeaders = []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept}
	}

	s.server.Use(middleware.CORSWithConfig(corsConfig))

	//! Setup Echo RequestLogger with zap based on log level ---------------------

	requestLoggerConfig := s.setUpLogger(logLevel)

	s.server.Use(middleware.RequestLoggerWithConfig(requestLoggerConfig))

	//! Hello World ---------------------
	s.server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})

	//! Start the server ---------------------
	go func() {
		err := s.server.Start(addr)
		if err != nil && err != http.ErrServerClosed {
			s.logger.Fatal(err.Error())
		}
	}()

	return nil
}

func (s *HTTPServer) setUpLogger(logLevel string) middleware.RequestLoggerConfig {
	if logLevel == "" {
		logLevel = DefaultLogLevel
	}

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
			switch logLevel {
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
					// todo: add more debug logs
				)
			default:
				s.logger.Error("invalid log level", zap.String("log_level", logLevel))
			}
			return nil
		},
	}

	return requestLoggerConfig
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

func (s *HTTPServer) GetRouter() *echo.Echo {
	return s.server
}
