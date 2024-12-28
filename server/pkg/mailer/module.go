package mailer

import (
	"context"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"

	"github.com/alsey89/people-matter/pkg/util"
)

type Module struct {
	logger *zap.Logger
	config *Config

	scope  string
	dialer *gomail.Dialer
}

type Params struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
}

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	TLS      bool
}

const (
	DefaultHost     = "0.0.0.0"
	DefaultPort     = 25
	DefaultUsername = ""
	DefaultPassword = ""
	DefaultTLS      = false

	DefaultSubject = "From Gogetter Mail Module"
	DefaultBody    = "This is an email from Gogetter Mail Module."
	DefaultFrom    = "mail@gogetter.com"
	DefaultTo      = "mail@gogetter.com"
)

//! MODULE ----------------------------------------------------------

// provides Mailer
func InjectModule(scope string, connectOrFatal bool) fx.Option {
	return fx.Module(
		scope,
		fx.Provide(func(p Params) *Module {
			m := &Module{scope: scope}

			m.logger = m.setupLogger(scope, p)
			m.config = m.setupConfig(scope)
			m.dialer = m.setupMailer()

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

// Instantiate the mailer without using the fx framework
func NewMailer(scope string, logger *zap.Logger) *Module {
	m := &Module{scope: scope}
	m.logger = logger.Named("[" + scope + "]")
	m.config = m.setupConfig(scope)
	m.dialer = m.setupMailer()

	m.onStart(context.Background())

	return m
}

//! INTERNAL ----------------------------------------------------------

func (m *Module) setupLogger(scope string, p Params) *zap.Logger {
	logger := p.Logger.Named("[" + scope + "]")

	return logger
}

func (m *Module) setupConfig(scope string) *Config {
	//set defaults
	viper.SetDefault(util.GetConfigPath(scope, "host"), DefaultHost)
	viper.SetDefault(util.GetConfigPath(scope, "port"), DefaultPort)
	viper.SetDefault(util.GetConfigPath(scope, "username"), DefaultUsername)
	viper.SetDefault(util.GetConfigPath(scope, "password"), DefaultPassword)
	viper.SetDefault(util.GetConfigPath(scope, "tls"), DefaultTLS)
	//populate config
	return &Config{
		Host:     viper.GetString(util.GetConfigPath(scope, "host")),
		Port:     viper.GetInt(util.GetConfigPath(scope, "port")),
		Username: viper.GetString(util.GetConfigPath(scope, "username")),
		Password: viper.GetString(util.GetConfigPath(scope, "password")),
		TLS:      viper.GetBool(util.GetConfigPath(scope, "tls")),
	}
}

func (m *Module) setupMailer() *gomail.Dialer {
	dialer := gomail.NewDialer(
		m.config.Host,
		m.config.Port,
		m.config.Username,
		m.config.Password,
	)

	return dialer
}

func (m *Module) onStart(ctx context.Context) error {
	m.logger.Info("Starting mailer module.")

	if viper.GetString("global.log_level") == "DEBUG" || viper.GetString("global.log_level") == "debug" {
		m.logConfigurations()
	}

	return nil
}

func (m *Module) onStop(ctx context.Context) error {
	m.logger.Info("Stopping mailer module.")

	return nil
}

func (m *Module) logConfigurations() {
	m.logger.Debug("----- Mailer Configuration -----")
	m.logger.Debug("Host", zap.String("Host", m.config.Host))
	m.logger.Debug("Port", zap.Int("Port", m.config.Port))
	m.logger.Debug("Username", zap.String("Username", m.config.Username))
	m.logger.Debug("Password", zap.String("Password", m.config.Password))
	m.logger.Debug("TLS", zap.Bool("TLS", m.config.TLS))
}

func (m *Module) testSMTPConnection() error {
	m.logger.Info("Testing SMTP connection...")
	s, err := m.dialer.Dial()
	if err != nil {
		m.logger.Error("Failed to connect to the SMTP server", zap.Error(err))
		return err
	}
	defer s.Close()

	m.logger.Info("Successfully connected to the SMTP server.")
	return nil
}

//! EXTERNAL ----------------------------------------------------------

// Creates a new email message
// Set details on the message and send using SendMail method
func (m *Module) NewMessage() *gomail.Message {
	return gomail.NewMessage()
}

// Sends the email message
// Create the message using NewMessage method
func (m *Module) SendMail(msg *gomail.Message) error {
	return m.dialer.DialAndSend(msg)
}

// Single method to create and send email
// Can be used instead of NewMessage and SendMail methods
func (m *Module) SendTransactionalMail(from string, to string, subject string, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	err := m.dialer.DialAndSend(msg)
	if err != nil {
		m.logger.Error("Failed to send email", zap.Error(err))
		return err
	}

	m.logger.Info("Email sent successfully.")
	return nil
}
