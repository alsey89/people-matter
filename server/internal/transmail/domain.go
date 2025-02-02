package transmail

import (
	"context"
	"fmt"

	"github.com/alsey89/people-matter/internal/common/errmgr"
	"github.com/alsey89/people-matter/internal/common/util"
	"github.com/alsey89/people-matter/internal/schema"
	"github.com/alsey89/people-matter/pkg/pgconn"

	"github.com/mailjet/mailjet-apiv3-go"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Domain struct {
	scope    string
	logger   *zap.Logger
	config   *Config
	params   Params
	MJClient *mailjet.Client
}

type Params struct {
	fx.In
	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
	DB        *pgconn.Module
}

type Config struct {
	senderEmail  string
	publicAPIKey string
	secretAPIKey string
	clientDomain string
}

const (
	defaultSenderEmail  = "team@peoplematter.app"
	defaultPublicAPIKey = "pub-api-key"
	defaultSecretAPIKey = "secret-api-key"
	defaultClientDomain = "localhost:3000"
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
			m.MJClient = m.setupMailjetClient()

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
	viper.SetDefault(util.GetConfigPath(scope, "sender_email"), defaultSenderEmail)
	viper.SetDefault(util.GetConfigPath(scope, "public_api_key"), defaultPublicAPIKey)
	viper.SetDefault(util.GetConfigPath(scope, "secret_api_key"), defaultSecretAPIKey)
	viper.SetDefault(util.GetConfigPath("global", "client_domain"), defaultClientDomain)

	return &Config{
		senderEmail:  viper.GetString(util.GetConfigPath(scope, "sender_email")),
		publicAPIKey: viper.GetString(util.GetConfigPath(scope, "public_api_key")),
		secretAPIKey: viper.GetString(util.GetConfigPath(scope, "secret_api_key")),
		clientDomain: viper.GetString(util.GetConfigPath("global", "client_domain")),
	}
}

func (d *Domain) setupMailjetClient() *mailjet.Client {
	mjClient := mailjet.NewMailjetClient(d.config.publicAPIKey, d.config.secretAPIKey)
	return mjClient
}

func (d *Domain) onStart(ctx context.Context) error {
	d.logger.Info("Starting transactional email domain.")

	if viper.GetString("global.log_level") == "DEBUG" || viper.GetString("global.log_level") == "debug" {
		d.logConfigurations()
	}

	return nil
}

func (m *Domain) onStop(ctx context.Context) error {
	m.logger.Info("Stopping auth domain.")
	return nil
}

func (d *Domain) logConfigurations() {
	d.logger.Debug("----- Seeder Configuration -----")
	d.logger.Debug("Sender Email: ", zap.String("sender_email", d.config.senderEmail))
	d.logger.Debug("Client Base URL: ", zap.String("client_base_url", d.config.clientDomain))
	d.logger.Debug("Public API Key: ", zap.String("public_api_key", d.config.publicAPIKey))
	d.logger.Debug("Secret API Key: ", zap.String("secret_api_key", d.config.secretAPIKey))
	d.logger.Debug("-------------------------------")
}

// ! Public ---------------------------------------------------------------

func (d *Domain) GetCompanyByID(ComnpanyID uint) (*schema.Company, error) {
	db := d.params.DB.GetDB()

	var company schema.Company
	err := db.Model(&schema.Company{}).Where("id = ?", ComnpanyID).First(&company).Error
	if err != nil {
		return nil, fmt.Errorf("getcompanyByID: %w", err)
	}

	return &company, nil
}

// Sends an email to the recipient with the given templateID.
// urlPath and variables are optional.
// urlPath will be converted to a full URL using the company's tenant identifier and the client domain.
func (d *Domain) SendMail(ComnpanyID uint, recipientEmail string, templateID int, urlPath *string, variables map[string]interface{}) error {
	if ComnpanyID == 0 || recipientEmail == "" || templateID == 0 {
		d.logger.Error("SendMail: Invalid parameters", zap.Any("ComnpanyID", ComnpanyID), zap.Any("recipientEmail", recipientEmail), zap.Any("templateID", templateID))
		return nil
	}
	if variables == nil {
		variables = make(map[string]interface{})
	}

	company, err := d.GetCompanyByID(ComnpanyID)
	if err != nil {
		d.logger.Error("SendMail: Failed to get company", zap.Error(err))
		return nil
	}
	if company == nil {
		d.logger.Error("SendMail: %w", zap.Error(errmgr.ErrNilCheckFailed))
		return nil
	}

	if variables["company"] == nil {
		variables["company"] = company.Name
	}

	if urlPath != nil {
		variables["url"], err = util.PathToFullURL(
			*urlPath,              // path string
			company.TenantID,      // subdomain string
			d.config.clientDomain, // domain string
		)
		if err != nil {
			d.logger.Error("SendMail:", zap.Error(err))
			return nil
		}
	}

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: fmt.Sprintf("%s_noreply@peoplematter.app", company.TenantID),
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: recipientEmail,
				},
			},
			TemplateID:       templateID,
			TemplateLanguage: true,
			Variables:        variables,
		},
	}

	d.logger.Debug("SendMail: Sending email", zap.Any("messagesInfo", messagesInfo))

	messages := mailjet.MessagesV31{Info: messagesInfo}

	// Send email via Mailjet
	_, err = d.MJClient.SendMailV31(&messages)
	if err != nil {
		d.logger.Error("SendMail: %w", zap.Error(err))
		return nil
	}

	return nil
}
