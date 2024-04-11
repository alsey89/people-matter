package mailer

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

const (
	DefaultHost        = "0.0.0.0"
	DefaultPort        = 25
	DefaultUsername    = ""
	DefaultAppPassword = ""
	DefaultTLS         = false

	DefaultSubject = "From HRMS"
	DefaultBody    = "This is an email from HRMS."
	DefaultFrom    = "hrms@hrms.com"
	DefaultTo      = "hrms@hrms.com"
)

type Config struct {
	Host        string
	Port        int
	Username    string
	AppPassword string //! Important: Use app password, not the account password.
	TLS         bool
}

type Mailer struct {
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

func InitiateModule(scope string) fx.Option {

	var m *Mailer

	return fx.Module(
		scope,
		fx.Provide(func(p Params) *Mailer {
			logger := p.Logger.Named("[" + scope + "]")
			config := loadConfig(scope)
			dialer := gomail.NewDialer(config.Host, config.Port, config.Username, config.AppPassword)

			m := &Mailer{
				logger: logger,
				config: config,
				scope:  scope,
				dialer: dialer,
			}

			return m
		}),
		fx.Populate(&m),
		fx.Invoke(func(p Params) *Mailer {

			p.Lifecycle.Append(
				fx.Hook{
					OnStart: m.onStart,
					OnStop:  m.onStop,
				},
			)

			return m
		}),
	)

}

func loadConfig(scope string) *Config {
	getConfigPath := func(key string) string {
		return fmt.Sprintf("%s.%s", scope, key)
	}

	//set defaults
	viper.SetDefault(getConfigPath("%s.host"), DefaultHost)
	viper.SetDefault(getConfigPath("%s.port"), DefaultPort)
	viper.SetDefault(getConfigPath("%s.username"), DefaultUsername)
	viper.SetDefault(getConfigPath("%s.password"), DefaultAppPassword)
	viper.SetDefault(getConfigPath("%s.tls"), DefaultTLS)
	//populate config
	return &Config{
		Host:        viper.GetString(getConfigPath("host")),
		Port:        viper.GetInt(getConfigPath("port")),
		Username:    viper.GetString(getConfigPath("username")),
		AppPassword: viper.GetString(getConfigPath("app_password")),
		TLS:         viper.GetBool(getConfigPath("tls")),
	}
}

func (m *Mailer) onStart(ctx context.Context) error {
	m.logger.Info("Mailer initiated")

	err := m.TestSMTPConnection()
	if err != nil {
		m.logger.Error("Failed to connect to the SMTP server", zap.Error(err))
		return err
	}

	//* Debug logs
	m.logger.Debug("----- Mailer Configuration -----")
	m.logger.Debug("Host", zap.String("Host", m.config.Host))
	m.logger.Debug("Port", zap.Int("Port", m.config.Port))
	m.logger.Debug("Username", zap.String("Username", m.config.Username))
	m.logger.Debug("AppPassword", zap.String("AppPassword", m.config.AppPassword))
	m.logger.Debug("TLS", zap.Bool("TLS", m.config.TLS))

	return nil
}

func (m *Mailer) onStop(ctx context.Context) error {

	m.logger.Info("Gomail stopped")

	return nil
}

func (m *Mailer) NewMessage() *gomail.Message {
	return gomail.NewMessage()
}

func (m *Mailer) Send(msg *gomail.Message) error {
	return m.dialer.DialAndSend(msg)
}

// ----------------------------------------------------------

func (m *Mailer) TestSMTPConnection() error {
	// m.logger.Info("Testing SMTP connection...")
	s, err := m.dialer.Dial()
	if err != nil {
		m.logger.Error("Failed to connect to the SMTP server", zap.Error(err))
		return err
	}
	defer s.Close()

	// m.logger.Info("Successfully connected to the SMTP server.")
	return nil
}

func (m *Mailer) SendTestMail(to, subject, body string) error {
	msg := m.NewMessage()
	msg.SetHeader("From", DefaultFrom)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	err := m.Send(msg)
	if err != nil {
		m.logger.Error("Failed to send test email", zap.Error(err))
		return err
	}

	m.logger.Info("Test email sent successfully.")
	return nil
}
