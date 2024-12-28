package mailer

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSetupConfig(t *testing.T) {
	scope := "mailer"
	m := Module{scope: scope}

	t.Run("TestSetupWithNoConfig", func(t *testing.T) {
		viper.Reset()
		defer viper.Reset()

		m.config = m.setupConfig(m.scope)

		assert.NotNil(t, m.config)
		assert.Equal(t, DefaultHost, m.config.Host)
		assert.Equal(t, DefaultPort, m.config.Port)
		assert.Equal(t, DefaultUsername, m.config.Username)
		assert.Equal(t, DefaultPassword, m.config.Password)
		assert.Equal(t, DefaultTLS, m.config.TLS)
	})

	t.Run("TestSetupWithConfig", func(t *testing.T) {
		viper.Reset()
		defer viper.Reset()

		viper.Set("mailer.host", "localhost_test")
		viper.Set("mailer.port", 25)
		viper.Set("mailer.username", "mail_test")
		viper.Set("mailer.password", "password_test")
		viper.Set("mailer.tls", true)

		m.config = m.setupConfig(scope)

		assert.Equal(t, "localhost_test", m.config.Host)
		assert.Equal(t, 25, m.config.Port)
		assert.Equal(t, "mail_test", m.config.Username)
		assert.Equal(t, "password_test", m.config.Password)
		assert.Equal(t, true, m.config.TLS)
	})

	t.Run("TestSetupWithPartialConfig", func(t *testing.T) {
		viper.Reset()
		defer viper.Reset()

		viper.Set("mailer.host", "localhost_test")
		viper.Set("mailer.port", 25)

		m.config = m.setupConfig(scope)

		assert.Equal(t, "localhost_test", m.config.Host)
		assert.Equal(t, 25, m.config.Port)
		assert.Equal(t, DefaultUsername, m.config.Username)
		assert.Equal(t, DefaultPassword, m.config.Password)
		assert.Equal(t, DefaultTLS, m.config.TLS)
	})
}

func TestSetupLogger(t *testing.T) {
	scope := "mailer"
	m := Module{
		scope: scope,
	}
	p := Params{
		Logger: zap.NewExample(),
	}
	m.logger = m.setupLogger(scope, p)

	assert.NotNil(t, m.logger)
	assert.IsType(t, &zap.Logger{}, m.logger)
}

func TestSetupMailer(t *testing.T) {
	scope := "mailer"
	m := Module{
		scope: scope,
		config: &Config{
			Host:     "localhost",
			Port:     25,
			Username: "testuser",
			Password: "testpassword",
		},
	}

	dialer := m.setupMailer()

	assert.NotNil(t, dialer)
	assert.Equal(t, "localhost", dialer.Host)
	assert.Equal(t, 25, dialer.Port)
	assert.Equal(t, "testuser", dialer.Username)
	assert.Equal(t, "testpassword", dialer.Password)
}
