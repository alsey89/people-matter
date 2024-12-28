package seeder

import (
	"context"

	"github.com/alsey89/people-matter/internal/fsp"
	"github.com/alsey89/people-matter/internal/identity"
	"github.com/alsey89/people-matter/internal/memorial"
	"github.com/alsey89/people-matter/pkg/pgconn"

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
	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
	DB        *pgconn.Module
	Identity  *identity.Domain
	Memorial  *memorial.Domain
	Tenant    *fsp.Domain
}

type Config struct {
}

const ()

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
	// viper.SetDefault(util.GetConfigPath("global", "client_domain"), defaultClientDomain)
	// viper.SetDefault(util.GetConfigPath(scope, "jwt_auth_scope"), defaultJWTAuthScope)
	// viper.SetDefault(util.GetConfigPath(scope, "jwt_email_scope"), defaultJWTEmailConfirmationScope)
	// viper.SetDefault(util.GetConfigPath(scope, "jwt_pw_reset_scope"), defaultJWTPasswordResetScope)

	return &Config{
		// ClientDomain:             viper.GetString(util.GetConfigPath("global", "client_domain")),
		// JWTAuthScope:              viper.GetString(util.GetConfigPath(scope, "jwt_auth_scope")),
		// JWTEmailConfirmationScope: viper.GetString(util.GetConfigPath(scope, "jwt_email_scope")),
		// JWTPasswordResetScope:     viper.GetString(util.GetConfigPath(scope, "jwt_pw_reset_scope")),
	}
}

func (d *Domain) onStart(ctx context.Context) error {
	d.logger.Info("Starting seeder domain.")

	if viper.GetString("global.log_level") == "DEBUG" || viper.GetString("global.log_level") == "debug" {
		d.logConfigurations()
	}

	d.logger.Info("Intializing essential data.")
	d.initializeFSPRoles()
	d.initializeMemorialRoles()

	d.logger.Info("Seeding test data.")
	d.seedTenants()
	d.seedUsersAndFSPRoles()
	d.seedMemorials()
	d.seedUserMemorialRoles()
	d.seedContributorApplications()
	d.seedContributorInvitations()

	return nil
}

func (m *Domain) onStop(ctx context.Context) error {
	m.logger.Info("Stopping auth domain.")
	return nil
}

func (d *Domain) logConfigurations() {
	d.logger.Debug("----- Seeder Configuration -----")

	d.logger.Debug("-------------------------------")
}
