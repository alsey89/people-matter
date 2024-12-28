package fsp

import (
	"context"
	"fmt"

	"github.com/alsey89/people-matter/internal/common/util"
	"github.com/alsey89/people-matter/internal/identity"
	"github.com/alsey89/people-matter/internal/transmail"
	"github.com/alsey89/people-matter/pkg/pgconn"
	"github.com/alsey89/people-matter/pkg/server"
	"github.com/alsey89/people-matter/pkg/token"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Domain struct {
	scope  string
	logger *zap.Logger
	config *Config
	params Params
}

type Params struct {
	fx.In
	Lifecycle    fx.Lifecycle
	Logger       *zap.Logger
	Server       *server.Module
	DB           *pgconn.Module
	TokenManager *token.Module
	Identity     *identity.Domain
	TransMail    *transmail.Domain
}

type Config struct {
	ClientDomain          string
	JWTPasswordResetScope string
}

const (
	defaultClientDomain = "http://localhost:3000"
)

// ! Domain ---------------------------------------------------------------

func InjectDomain(scope string) fx.Option {
	return fx.Module(
		scope,
		fx.Provide(func(p Params) *Domain {
			m := &Domain{scope: scope}
			m.params = p
			m.logger = m.setupLogger(scope, p)
			m.config = m.setupConfig(scope)

			return m
		}),
		fx.Invoke(func(m *Domain, p Params) {
			p.Lifecycle.Append(
				fx.Hook{
					OnStart: m.onStart,
					OnStop:  m.onStop,
				},
			)
		}),
	)
}

// ! Internal ---------------------------------------------------------------
func (d *Domain) setupLogger(scope string, p Params) *zap.Logger {
	logger := p.Logger.Named("[" + scope + "]")
	return logger
}

func (d *Domain) setupConfig(scope string) *Config {
	viper.SetDefault(util.GetConfigPath("global", "client_domain"), defaultClientDomain)
	// viper.SetDefault(util.GetConfigPath(scope, "jwt_auth_scope"), defaultJWTAuthScope)
	// viper.SetDefault(util.GetConfigPath(scope, "jwt_email_scope"), defaultJWTEmailConfirmationScope)
	// viper.SetDefault(util.GetConfigPath(scope, "jwt_pw_reset_scope"), defaultJWTPasswordResetScope)

	return &Config{
		ClientDomain: viper.GetString(util.GetConfigPath("global", "client_domain")),
		// JWTAuthScope:              viper.GetString(util.GetConfigPath(scope, "jwt_auth_scope")),
		// JWTEmailConfirmationScope: viper.GetString(util.GetConfigPath(scope, "jwt_email_scope")),
		// JWTPasswordResetScope:     viper.GetString(util.GetConfigPath(scope, "jwt_pw_reset_scope")),
	}
}

func (d *Domain) onStart(ctx context.Context) error {
	d.logger.Info(fmt.Sprintf("Starting %s domain.", d.scope))

	if viper.GetString("global.log_level") == "DEBUG" || viper.GetString("global.log_level") == "debug" {
		d.logConfigurations()
	}

	d.registerRoutes()

	return nil
}

func (m *Domain) onStop(ctx context.Context) error {
	m.logger.Info(fmt.Sprintf("Stopping %s domain.", m.scope))
	return nil
}

func (d *Domain) logConfigurations() {
	d.logger.Debug("----- Tenant Configuration -----")
	d.logger.Debug("Client Base URL: " + d.config.ClientDomain)
	d.logger.Debug("-------------------------------")
}
