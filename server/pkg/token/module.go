package token

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/alsey89/people-matter/pkg/util"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Module struct {
	scope   string
	logger  *zap.Logger
	configs map[string]*Config
}

type Params struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
}

type Config struct {
	TokenLookup   string
	SigningKey    string
	SigningMethod string
	ExpInHours    int
	ClientDomain  string
}

const (
	defaultTokenScope    = "default"
	defaultSigningKey    = "secret"
	defaultTokenLookup   = "cookie:jwt"
	defaultSigningMethod = "HS256"
	defaultExpInHours    = 72
	defaultClientDomain  = "http://localhost:3000"
)

// ! Module ---------------------------------------------------------------

// Provides the Module struct to the fx framework, and registers Lifecycle hooks.
// Takes tokenScopes as variadic parameter.
// Each tokenScope will generate a Config struct, holding separate configurations.
func InjectModule(moduleScope string, tokenScopes ...string) fx.Option {
	return fx.Module(
		moduleScope,
		fx.Provide(func(p Params) *Module {

			m := &Module{scope: moduleScope}
			m.logger = m.setupLogger(moduleScope, p)
			m.configs = m.setupConfig(tokenScopes...)

			return m
		}),
		fx.Invoke(func(m *Module, p Params) {
			p.Lifecycle.Append(
				fx.Hook{
					OnStart: m.onStart,
					OnStop:  m.onStop,
				},
			)
		}),
	)
}

// Instantiate the Module without using the fx framework
func NewTokenManager(moduleScope string, logger *zap.Logger, tokenScopes ...string) *Module {
	m := &Module{scope: moduleScope}
	m.logger = logger.Named("[" + moduleScope + "]")
	m.configs = m.setupConfig(tokenScopes...)

	m.onStart(context.Background())

	return m
}

//! INTERNAL ---------------------------------------------------------------

func (m *Module) setupLogger(moduleScope string, p Params) *zap.Logger {
	logger := p.Logger.Named("[" + moduleScope + "]")
	return logger
}

func (m *Module) setupConfig(tokenScopes ...string) map[string]*Config {
	configs := make(map[string]*Config)

	// defaultTokenscope will be used if no tokenScopes provided
	if len(tokenScopes) == 0 {
		m.logger.Warn("No token scopes provided. Using default token scope.")
		tokenScopes = append(tokenScopes, defaultTokenScope)
	}

	for _, scope := range tokenScopes {

		viper.SetDefault(util.GetConfigPath(scope, "token_lookup"), defaultTokenLookup)
		viper.SetDefault(util.GetConfigPath(scope, "signing_key"), defaultSigningKey)
		viper.SetDefault(util.GetConfigPath(scope, "signing_method"), defaultSigningMethod)
		viper.SetDefault(util.GetConfigPath(scope, "exp_in_hours"), defaultExpInHours)
		viper.SetDefault(util.GetConfigPath(scope, "client_domain"), defaultClientDomain)

		configs[scope] = &Config{
			TokenLookup:   viper.GetString(util.GetConfigPath(scope, "token_lookup")),
			SigningKey:    viper.GetString(util.GetConfigPath(scope, "signing_key")),
			SigningMethod: viper.GetString(util.GetConfigPath(scope, "signing_method")),
			ExpInHours:    viper.GetInt(util.GetConfigPath(scope, "exp_in_hours")),
			ClientDomain:  viper.GetString(util.GetConfigPath("global", "client_domain")),
		}
	}

	return configs
}

func (m *Module) onStart(ctx context.Context) error {
	m.logger.Info("Starting token manager.")

	if viper.GetString("global.log_level") == "DEBUG" || viper.GetString("global.log_level") == "debug" {
		m.logConfigurations()
	}

	return nil
}

func (m *Module) onStop(ctx context.Context) error {
	m.logger.Info("Stopping token manager.")
	return nil
}

func (m *Module) logConfigurations() {
	for scope, config := range m.configs {
		m.logger.Debug("----- Token Manager Configuration -----")
		m.logger.Debug("TokenScope", zap.String("TokenScope", scope))
		m.logger.Debug("TokenLookup", zap.String("TokenLookup", config.TokenLookup))
		m.logger.Debug("SigningKey", zap.String("SigningKey", config.SigningKey))
		m.logger.Debug("SigningMethod", zap.String("SigningMethod", config.SigningMethod))
		m.logger.Debug("ExpInHours", zap.Int("ExpInHours", config.ExpInHours))
		m.logger.Debug("ClientDomain", zap.String("ClientDomain", config.ClientDomain))
	}
}

//! EXTERNAL ---------------------------------------------------------------

/*
Generates a JWT token with the provided additional claims for a specific scope.
Use jwt.MapClaims from "github.com/golang-jwt/jwt/v5"
*/
func (m *Module) GenerateToken(tokenScope string, additionalClaims jwt.MapClaims) (*string, error) {
	scopeConfig, err := m.getConfigHelper(tokenScope)
	if err != nil {
		m.logger.Error("Config not found", zap.String("Scope:", tokenScope))
		return nil, err
	}

	claims := jwt.MapClaims{
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(scopeConfig.ExpInHours))),
	}

	for key, value := range additionalClaims {
		claims[key] = value
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(scopeConfig.SigningMethod), claims)
	t, err := token.SignedString([]byte(scopeConfig.SigningKey))
	if err != nil {
		m.logger.Error("Failed to generate token", zap.Error(err))
		return nil, err
	}
	return &t, nil
}

/*
Generates a JWT token with the provided additional claims for a specific scope.
Use jwt.MapClaims from "github.com/golang-jwt/jwt/v5"
*/
func (m *Module) GenerateTokenAndHTTPonlyCookie(tokenScope string, additionalClaims jwt.MapClaims) (*http.Cookie, error) {
	scopeConfig, err := m.getConfigHelper(tokenScope)
	if err != nil {
		m.logger.Error("Config not found", zap.String("Scope:", tokenScope))
		return nil, err
	}

	claims := jwt.MapClaims{
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(scopeConfig.ExpInHours))),
	}

	for key, value := range additionalClaims {
		claims[key] = value
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(scopeConfig.SigningMethod), claims)
	t, err := token.SignedString([]byte(scopeConfig.SigningKey))
	if err != nil {
		m.logger.Error("Failed to generate token", zap.Error(err))
		return nil, err
	}

	cookie := &http.Cookie{
		Name:    "jwt",
		Value:   t,
		Expires: time.Now().Add(time.Hour * time.Duration(scopeConfig.ExpInHours)),
		Path:    "/",
		// Domain:   strings.Split(scopeConfig.ClientDomain, ":")[0],
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	m.logger.Debug("Generated cookie", zap.Any("Cookie", cookie))

	return cookie, nil
}

/*
Returns an echo middleware that validates JWT tokens for a specific scope.
Middleware validates the JWT token, parses claims, and stores them in context under the key "user".
*/
func (m *Module) GetJWTMiddleware(tokenScope string) echo.MiddlewareFunc {
	scopeConfig, err := m.getConfigHelper(tokenScope)
	if err != nil {
		m.logger.Error("Config not found", zap.String("Scope:", tokenScope))
		return nil
	}

	return echojwt.WithConfig(echojwt.Config{
		SigningKey:    []byte(scopeConfig.SigningKey),
		SigningMethod: scopeConfig.SigningMethod,
		TokenLookup:   scopeConfig.TokenLookup,
	})
}

func (m *Module) getConfigHelper(scope string) (*Config, error) {
	config, exists := m.configs[scope]
	if !exists {
		return nil, fmt.Errorf("config for scope %s not found", scope)
	}
	return config, nil
}
