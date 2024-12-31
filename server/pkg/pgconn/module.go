package pgconn

import (
	"context"
	"fmt"
	"strings"

	"github.com/alsey89/people-matter/pkg/util"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
)

type Module struct {
	config *Config
	db     *gorm.DB
	logger *zap.Logger
	scope  string
}

type Params struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
}

type Config struct {
	DBName   string
	Host     string
	LogLevel string
	Password string
	Port     int
	SSLMode  string
	User     string
}

const (
	DefaultHost     = "0.0.0.0"
	DefaultPort     = 5432
	DefaultDbName   = "postgres"
	DefaultUser     = "postgres"
	DefaultPassword = "postgres"
	DefaultSSLMode  = "allow"
	DefaultLogLevel = "info"
)

//! Module ---------------------------------------------------------------

// Provides the module to the fx framework
func InjectModule(scope string) fx.Option {
	return fx.Module(
		scope,
		fx.Provide(func(p Params) *Module {

			m := &Module{scope: scope}
			m.config = m.setupConfig(scope)
			m.logger = m.setupLogger(scope, p)
			m.db = m.setUpDB()

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

// Instantiates new Module without using the fx framework.
func NewPGConn(scope string, logger *zap.Logger) *Module {
	m := &Module{scope: scope}
	m.logger = logger.Named("[" + scope + "]")
	m.config = m.setupConfig(scope)
	m.db = m.setUpDB()

	m.onStart(context.Background())

	return m
}

//! INTERNAL ---------------------------------------------------------------

func (m *Module) setupConfig(scope string) *Config {
	//set default values
	viper.SetDefault(util.GetConfigPath(scope, "host"), DefaultHost)
	viper.SetDefault(util.GetConfigPath(scope, "port"), DefaultPort)
	viper.SetDefault(util.GetConfigPath(scope, "dbname"), DefaultDbName)
	viper.SetDefault(util.GetConfigPath(scope, "user"), DefaultUser)
	viper.SetDefault(util.GetConfigPath(scope, "password"), DefaultPassword)
	viper.SetDefault(util.GetConfigPath(scope, "sslmode"), DefaultSSLMode)
	viper.SetDefault(util.GetConfigPath("global", "log_level"), DefaultLogLevel)

	return &Config{
		Host:     viper.GetString(util.GetConfigPath(scope, "host")),
		Port:     viper.GetInt(util.GetConfigPath(scope, "port")),
		DBName:   viper.GetString(util.GetConfigPath(scope, "dbname")),
		User:     viper.GetString(util.GetConfigPath(scope, "user")),
		Password: viper.GetString(util.GetConfigPath(scope, "password")),
		SSLMode:  viper.GetString(util.GetConfigPath(scope, "sslmode")),
		LogLevel: viper.GetString(util.GetConfigPath("global", "log_level")),
	}
}

func (m *Module) setupLogger(scope string, p Params) *zap.Logger {
	logger := p.Logger.Named("[" + scope + "]")
	return logger
}

func (m *Module) setUpDB() *gorm.DB {
	dsn := m.getConnectionStringFromConfig()
	loglevel := m.getLogLevelFromConfig()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gorm_logger.Default.LogMode(loglevel),
	})
	if err != nil {
		m.logger.Fatal("Error connecting to database", zap.Error(err))
	}

	return db
}

func (m *Module) getConnectionStringFromConfig() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		m.config.Host, m.config.Port, m.config.User, m.config.Password, m.config.DBName, m.config.SSLMode)
}

func (m *Module) getLogLevelFromConfig() gorm_logger.LogLevel {
	lowerCaseLogLevel := strings.ToLower(m.config.LogLevel)

	switch lowerCaseLogLevel {
	case "debug":
		return gorm_logger.Info // Log all queries
	case "dev":
		return gorm_logger.Warn // Log warnings and errors
	case "prod":
		return gorm_logger.Silent // No logs for production
	default:
		return gorm_logger.Info // Default to warnings and errors
	}
}

func (m *Module) onStart(context.Context) error {
	m.logger.Info("Starting database connection.")

	if viper.GetString("global.log_level") == "DEBUG" || viper.GetString("global.log_level") == "debug" {
		m.logConfigurations()
	}

	return nil
}

func (m *Module) onStop(context.Context) error {
	m.logger.Info("Stopping database connection.")

	db, err := m.db.DB()
	if err != nil {
		m.logger.Error("Error getting DB from GORM", zap.Error(err))
		return err
	}

	err = db.Close()
	if err != nil {
		m.logger.Error("Error closing DB", zap.Error(err))
	}

	return nil
}

func (m *Module) logConfigurations() {
	m.logger.Debug("----- Database Configuration -----")
	m.logger.Debug("Host", zap.String("host", m.config.Host))
	m.logger.Debug("Port", zap.Int("port", m.config.Port))
	m.logger.Debug("DBName", zap.String("dbname", m.config.DBName))
	m.logger.Debug("User", zap.String("user", m.config.User))
	m.logger.Debug("SSLMode", zap.String("sslmode", m.config.SSLMode))
	m.logger.Debug("LogLevel", zap.String("log_level", m.config.LogLevel))
}

//! EXTERNAL ---------------------------------------------------------------

// Applies the passed in schema.
// If autoMigrate is true, it will automatically migrate the schema at startup.
// If autoMigrate is false, it will skip the migration.
func (m *Module) ApplySchema(autoMigrate bool, schema ...interface{}) {
	if !autoMigrate {
		m.logger.Info("Skipping auto migration.")
		return
	}
	m.logger.Info("Migration started.")
	err := m.db.AutoMigrate(schema...)
	if err != nil {
		m.logger.Error("Error with migration.", zap.Error(err))
	}

	m.logger.Info("Migration completed.")
}

// Returns the GORM DB instance
func (m *Module) GetDB() *gorm.DB {
	return m.db
}
