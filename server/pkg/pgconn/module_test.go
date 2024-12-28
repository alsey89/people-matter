package pgconn

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	gorm_logger "gorm.io/gorm/logger"
)

func TestSetupConfig(t *testing.T) {
	scope := "database"
	m := Module{scope: scope}

	t.Run("TestSetupWithNoConfig", func(t *testing.T) {
		viper.Reset()
		defer viper.Reset()

		m.config = m.setupConfig(m.scope)

		assert.NotNil(t, m.config)
		assert.Equal(t, DefaultHost, m.config.Host)
		assert.Equal(t, DefaultPort, m.config.Port)
		assert.Equal(t, DefaultDbName, m.config.DBName)
		assert.Equal(t, DefaultUser, m.config.User)
		assert.Equal(t, DefaultPassword, m.config.Password)
		assert.Equal(t, DefaultSSLMode, m.config.SSLMode)
		assert.Equal(t, DefaultLogLevel, m.config.LogLevel)
	})

	t.Run("TestSetupWithConfig", func(t *testing.T) {
		viper.Reset()
		defer viper.Reset()

		viper.Set("database.host", "localhost_test")
		viper.Set("database.port", 5432)
		viper.Set("database.dbname", "postgres_test")
		viper.Set("database.user", "postgres_test")
		viper.Set("database.password", "password_test")
		viper.Set("database.sslmode", "false")
		viper.Set("database.log_level", "error")
		viper.Set("database.auto_migrate", true)

		m.config = m.setupConfig(scope)

		assert.Equal(t, "localhost_test", m.config.Host)
		assert.Equal(t, 5432, m.config.Port)
		assert.Equal(t, "postgres_test", m.config.DBName)
		assert.Equal(t, "postgres_test", m.config.User)
		assert.Equal(t, "password_test", m.config.Password)
		assert.Equal(t, "false", m.config.SSLMode)
		assert.Equal(t, "error", m.config.LogLevel)
	})

	t.Run("TestSetupWithPartialConfig", func(t *testing.T) {
		viper.Reset()
		defer viper.Reset()

		viper.Set("database.host", "localhost_test")
		viper.Set("database.port", 5432)

		m.config = m.setupConfig(scope)

		assert.Equal(t, "localhost_test", m.config.Host)
		assert.Equal(t, 5432, m.config.Port)
		assert.Equal(t, DefaultDbName, m.config.DBName)
		assert.Equal(t, DefaultUser, m.config.User)
		assert.Equal(t, DefaultPassword, m.config.Password)
		assert.Equal(t, DefaultSSLMode, m.config.SSLMode)
		assert.Equal(t, DefaultLogLevel, m.config.LogLevel)
	})
}

func TestSetupLogger(t *testing.T) {
	scope := "database"
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

//todo: ?? test setupDB?

func TestGetConnectionStringFromConfig(t *testing.T) {
	m := &Module{
		config: &Config{
			Host:     "testhost",
			Port:     1234,
			User:     "testuser",
			Password: "testpassword",
			DBName:   "testdb",
			SSLMode:  "testssl",
		},
	}

	expected := "host=testhost port=1234 user=testuser password=testpassword dbname=testdb sslmode=testssl"
	actual := m.getConnectionStringFromConfig()

	assert.Equal(t, expected, actual)
}

func TestGetLogLevelFromConfig(t *testing.T) {
	d := &Module{
		config: &Config{
			LogLevel: "debug",
		},
	}
	t.Run("Silent", func(t *testing.T) {
		d.config.LogLevel = "silent"
		expected := gorm_logger.Silent
		actual := d.getLogLevelFromConfig()
		assert.Equal(t, expected, actual)
	})

	t.Run("Error", func(t *testing.T) {
		d.config.LogLevel = "error"
		expected := gorm_logger.Error
		actual := d.getLogLevelFromConfig()
		assert.Equal(t, expected, actual)
	})

	t.Run("Warn", func(t *testing.T) {
		d.config.LogLevel = "warn"
		expected := gorm_logger.Warn
		actual := d.getLogLevelFromConfig()
		assert.Equal(t, expected, actual)
	})

	t.Run("Info", func(t *testing.T) {
		d.config.LogLevel = "info"
		expected := gorm_logger.Info
		actual := d.getLogLevelFromConfig()
		assert.Equal(t, expected, actual)
	})

	t.Run("Default", func(t *testing.T) {
		d.config.LogLevel = "invalid"
		expected := gorm_logger.Info
		actual := d.getLogLevelFromConfig()
		assert.Equal(t, expected, actual)
	})
}

//todo: ?? test migrate?
