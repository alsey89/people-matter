package echo_jwt

import (
	"context"
	"fmt"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	// defaultSecret      = "secret"
	defaultTokenLookup   = "cookie:jwt"
	defaultSigningMethod = "HS256"
)

type JWT struct {
	logger *zap.Logger
	config *echojwt.Config

	scope      string
	middleware echo.MiddlewareFunc
}

type Params struct {
	fx.In

	Logger    *zap.Logger
	Lifecycle fx.Lifecycle
}

func InitiateModule(scope string) fx.Option {
	return fx.Module(
		scope,
		fx.Provide(func(p Params) (*JWT, error) {
			logger := p.Logger.Named("[" + scope + "]")
			config := loadConfig(scope)

			middleware := echojwt.WithConfig(config)

			j := &JWT{
				logger:     logger,
				config:     &config,
				scope:      scope,
				middleware: middleware,
			}

			return j, nil
		}),
		fx.Invoke(func(j *JWT, p Params) {
			p.Lifecycle.Append(
				fx.Hook{
					OnStart: j.onStart,
					OnStop:  j.onStop,
				},
			)
		}),
	)
}

func loadConfig(scope string) echojwt.Config {
	getConfigPath := func(key string) string {
		return fmt.Sprintf("%s.%s", scope, key)
	}

	//set defaults
	viper.SetDefault(getConfigPath("token_lookup"), defaultTokenLookup)
	viper.SetDefault(getConfigPath("signing_key"), scope) //* default key is the scope
	viper.SetDefault(getConfigPath("signing_method"), defaultSigningMethod)

	return echojwt.Config{
		TokenLookup:   viper.GetString(getConfigPath("token_lookup")),
		SigningKey:    viper.GetString(getConfigPath("signing_key")),
		SigningMethod: viper.GetString(getConfigPath("signing_method")),
	}
}

func (j *JWT) onStart(ctx context.Context) error {
	j.logger.Info("JWT Middleware initiated")

	j.PrintDebugLogs()

	return nil
}

func (j *JWT) onStop(ctx context.Context) error {
	j.logger.Info("Stopping JWT")
	return nil
}

func (j *JWT) Middleware() echo.MiddlewareFunc {
	return j.middleware
}

func (j *JWT) PrintDebugLogs() {
	j.logger.Debug("----- JWT Middleware Configuration -----")
	j.logger.Debug("TokenLookup", zap.String("TokenLookup", j.config.TokenLookup))
	j.logger.Debug("SigningKey", zap.Any("SigningKey", j.config.SigningKey))
	j.logger.Debug("SigningMethod", zap.String("SigningMethod", j.config.SigningMethod))
}
