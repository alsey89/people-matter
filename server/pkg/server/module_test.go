package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSetupConfig(t *testing.T) {
	m := Module{
		scope: "server",
	}

	t.Run("TestSetupWithNoConfig", func(t *testing.T) {
		viper.Reset()
		defer viper.Reset()

		m.config = m.setupConfig(m.scope)

		assert.NotNil(t, m.config)
		assert.Equal(t, DefaultAllowHeaders, m.config.AllowHeaders)
		assert.Equal(t, DefaultAllowMethods, m.config.AllowMethods)
		assert.Equal(t, DefaultAllowOrigins, m.config.AllowOrigins)
		assert.Equal(t, DefaultCSRFProtection, m.config.CSRFProtection)
		assert.Equal(t, DefaultCSRFSecure, m.config.CSRFSecure)
		assert.Equal(t, DefaultCSRFDomain, m.config.CSRFDomain)
		assert.Equal(t, DefaultHost, m.config.Host)
		assert.Equal(t, DefaultPort, m.config.Port)
		assert.Equal(t, DefaultServerLogLevel, m.config.ServerLogLevel)
	})

	t.Run("TestSetupWithConfig", func(t *testing.T) {
		viper.Reset()
		defer viper.Reset()

		viper.Set("server.allow_headers", "Content-Type,Authorization, X-CSRF-Token, Set-Cookie, Cookie, jwt")
		viper.Set("server.allow_methods", "GET, PUT, POST, DELETE, OPTIONS")
		viper.Set("server.allow_origins", "http://localhost:3000")
		viper.Set("server.csrf_protection", true)
		viper.Set("server.csrf_secure", true)
		viper.Set("server.csrf_domain", "localhost")
		viper.Set("server.host", "localhost")
		viper.Set("server.port", 3001)
		viper.Set("server.server_log_level", "DEBUG")
		viper.Set("server.system_log_level", "DEBUG")

		m.config = m.setupConfig(m.scope)

		assert.NotNil(t, m.config)
		assert.Equal(t, "Content-Type,Authorization, X-CSRF-Token, Set-Cookie, Cookie, jwt", m.config.AllowHeaders)
		assert.Equal(t, "GET, PUT, POST, DELETE, OPTIONS", m.config.AllowMethods)
		assert.Equal(t, "http://localhost:3000", m.config.AllowOrigins)
		assert.Equal(t, true, m.config.CSRFProtection)
		assert.Equal(t, true, m.config.CSRFSecure)
		assert.Equal(t, "localhost", m.config.CSRFDomain)
		assert.Equal(t, "localhost", m.config.Host)
		assert.Equal(t, 3001, m.config.Port)
		assert.Equal(t, "DEBUG", m.config.ServerLogLevel)
	})

	t.Run("TestSetupWithPartialConfig", func(t *testing.T) {
		viper.Reset()
		defer viper.Reset()

		viper.Set("server.host", "localhost_test")
		viper.Set("server.port", 3001)

		m.config = m.setupConfig(m.scope)

		assert.Equal(t, "localhost_test", m.config.Host)
		assert.Equal(t, 3001, m.config.Port)
		assert.Equal(t, DefaultAllowHeaders, m.config.AllowHeaders)
		assert.Equal(t, DefaultAllowMethods, m.config.AllowMethods)
		assert.Equal(t, DefaultAllowOrigins, m.config.AllowOrigins)
		assert.Equal(t, DefaultCSRFProtection, m.config.CSRFProtection)
		assert.Equal(t, DefaultCSRFSecure, m.config.CSRFSecure)
		assert.Equal(t, DefaultCSRFDomain, m.config.CSRFDomain)
		assert.Equal(t, DefaultServerLogLevel, m.config.ServerLogLevel)
	})
}

func TestSetupLogger(t *testing.T) {
	scope := "server"
	m := Module{
		scope: "server",
	}
	p := Params{
		Logger: zap.NewExample(),
	}
	m.logger = m.setupLogger(scope, p)

	assert.NotNil(t, m.logger)
	assert.IsType(t, &zap.Logger{}, m.logger)
}

func TestSetupServer(t *testing.T) {
	m := Module{
		scope: "server",
	}
	m.server = m.setupServer()

	assert.NotNil(t, m.server)
	assert.IsType(t, &echo.Echo{}, m.server)
}

func TestSetUpCorsMiddleware(t *testing.T) {
	m := Module{
		scope: "server",
		config: &Config{
			AllowOrigins: "*",
			AllowMethods: "GET, POST",
			AllowHeaders: "Content-Type, Authorization",
		},
		server: echo.New(),
	}
	m.setUpCorsMiddleware()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type, Authorization")

	m.server.ServeHTTP(rec, req)

	assert.True(t, rec.Code == http.StatusOK || rec.Code == http.StatusNoContent)
	assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST", rec.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Type, Authorization", rec.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCSRFMiddleware(t *testing.T) {

	// creates module + middleware + route
	newModuleWithCSRF := func() Module {
		m := Module{
			scope: "server",
			config: &Config{
				CSRFProtection: true,
				CSRFSecure:     false,
				CSRFDomain:     "localhost",
			},
			server: echo.New(),
		}
		m.setUpCSRFMiddleware()
		m.server.POST("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})
		return m
	}

	t.Run("TestRequestWithoutCSRFToken", func(t *testing.T) {
		m := newModuleWithCSRF()

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		m.server.ServeHTTP(rec, req)

		// When token is not found, expect:
		// 400: Bad Request, message: "missing value in cookies"
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("TestRequestWithCSRFToken", func(t *testing.T) {
		m := newModuleWithCSRF()

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.AddCookie(&http.Cookie{Name: "_csrf", Value: "test"})

		m.server.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "test", rec.Body.String())
	})
}

// ?? test request logger middleware

func TestStartServer(t *testing.T) {
	m := Module{
		scope: "server",
		config: &Config{
			Host: "localhost",
			Port: 3001,
		},
		server: echo.New(),
		logger: zap.NewNop(),
	}

	// Define a simple route for testing
	m.server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Start the server in a goroutine
	go m.startServer(true, false)

	// Sleep for 5 seconds to allow the server to start
	time.Sleep(5 * time.Second)

	var err error
	var resp *http.Response

	resp, err = http.Get("http://localhost:3001")

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	if err := m.server.Close(); err != nil {
		t.Errorf("Failed to close server: %v", err)
	}
}

// External --------------------------------------------------------------------

func TestGetServer(t *testing.T) {
	m := Module{
		scope:  "server",
		server: echo.New(),
	}

	s := m.GetServer()

	assert.NotNil(t, s)
	assert.IsType(t, &echo.Echo{}, s)
}
