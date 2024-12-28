package memorial

import (
	"context"
	"fmt"
	"text/template"
	"time"

	"github.com/alsey89/people-matter/internal/common/util"
	"github.com/alsey89/people-matter/internal/fsp"
	"github.com/alsey89/people-matter/internal/identity"
	"github.com/alsey89/people-matter/internal/transmail"
	"github.com/alsey89/people-matter/pkg/pgconn"
	"github.com/alsey89/people-matter/pkg/s3conn"
	"github.com/alsey89/people-matter/pkg/server"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Domain struct {
	scope            string
	logger           *zap.Logger
	config           *Config
	params           Params
	s3ContribKeyTmpl *template.Template
	s3DeployKeyTmpl  *template.Template
}

type Params struct {
	fx.In
	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
	Server    *server.Module
	DB        *pgconn.Module
	S3        *s3conn.Module
	Tenant    *fsp.Domain
	Identity  *identity.Domain
	TransMail *transmail.Domain
}

type Config struct {
	ClientDomain         string
	S3PresignedURLExpiry time.Duration
	S3KeyFormat          string
	InvitationExpiry     time.Duration

	S3ContribKeyFormat string
	S3DeployKeyFormat  string

	GitLabTriggerToken string
}

const (
	defaultClientDomain     = "http://localhost:3000"
	defaultInvitationExpiry = time.Hour * 24 * 7 // 1 week

	defaultS3PresignedURLExpiry = 15 * time.Minute
	defaultS3ContribKeyFormat   = "{{.TenantID}}/{{.MemorialID}}/{{.ContributorID}}/{{.Date}}_{{.UUID}}_og"
	defaultS3DeployKeyFormat    = "{{.TenantID}}/{{.MemorialID}}/{{.Date}}_{{.UUID}}"

	defaultGitLabTriggerToken = "token-from-gitlab-ci"
)

type S3KeyData struct {
	TenantID      int
	MemorialID    int
	ContributorID int
	Date          string
	UUID          string
}

func InjectDomain(scope string) fx.Option {
	return fx.Module(
		scope,
		fx.Provide(func(p Params) *Domain {
			m := &Domain{scope: scope}
			m.params = p
			m.logger = m.setupLogger(scope, p)
			m.config = m.setupConfig(scope)
			m.s3ContribKeyTmpl = m.parseS3KeyFormat("contribKeyTmpl", m.config.S3ContribKeyFormat)
			m.s3DeployKeyTmpl = m.parseS3KeyFormat("deployKeyTmpl", m.config.S3DeployKeyFormat)

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

func (d *Domain) setupLogger(scope string, p Params) *zap.Logger {
	logger := p.Logger.Named("[" + scope + "]")
	return logger
}

func (d *Domain) setupConfig(scope string) *Config {
	viper.SetDefault(util.GetConfigPath("global", "client_domain"), defaultClientDomain)
	viper.SetDefault(util.GetConfigPath(scope, "s3_presigned_url_expiry"), defaultS3PresignedURLExpiry)
	viper.SetDefault(util.GetConfigPath(scope, "s3_contrib_key_format"), defaultS3ContribKeyFormat)
	viper.SetDefault(util.GetConfigPath(scope, "s3_deploy_key_format"), defaultS3DeployKeyFormat)
	viper.SetDefault(util.GetConfigPath(scope, "invitation_expiry"), defaultInvitationExpiry)
	viper.SetDefault(util.GetConfigPath(scope, "gitlab_trigger_token"), defaultGitLabTriggerToken)

	return &Config{
		ClientDomain:         viper.GetString(util.GetConfigPath("global", "client_domain")),
		S3PresignedURLExpiry: viper.GetDuration(util.GetConfigPath(scope, "s3_presigned_url_expiry")),
		S3ContribKeyFormat:   viper.GetString(util.GetConfigPath(scope, "s3_contrib_key_format")),
		S3DeployKeyFormat:    viper.GetString(util.GetConfigPath(scope, "s3_deploy_key_format")),
		InvitationExpiry:     viper.GetDuration(util.GetConfigPath(scope, "invitation_expiry")),
		GitLabTriggerToken:   viper.GetString(util.GetConfigPath(scope, "gitlab_trigger_token")),
	}
}

func (d *Domain) parseS3KeyFormat(templateName string, s3KeyFormat string) *template.Template {
	tmpl, err := template.New(templateName).Parse(string(s3KeyFormat))
	if err != nil {
		d.logger.Fatal("Failed to parse S3 key format", zap.Error(err))
	}

	return tmpl
}

func (d *Domain) onStart(ctx context.Context) error {
	d.logger.Info(fmt.Sprintf("Starting %s domain.", d.scope))

	if viper.GetString("global.log_level") == "DEBUG" || viper.GetString("global.log_level") == "debug" {
		d.logConfigurations()
	}

	d.registerRoutes()

	return nil
}

func (d *Domain) onStop(ctx context.Context) error {
	d.logger.Info(fmt.Sprintf("Stopping %s domain.", d.scope))
	return nil
}

func (d *Domain) logConfigurations() {
	d.logger.Debug("----- Memorial Configuration -----")
	d.logger.Debug("Client Base URL: " + d.config.ClientDomain)
	d.logger.Debug("S3 Presigned URL Expiry: " + d.config.S3PresignedURLExpiry.String())
	d.logger.Debug("S3 Contrib Key Format: " + d.config.S3ContribKeyFormat)
	d.logger.Debug("S3 Deploy Key Format: " + d.config.S3DeployKeyFormat)
	d.logger.Debug("-------------------------------")
}
