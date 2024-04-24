package postgres

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
)

const (
	DefaultHost        = "0.0.0.0"
	DefaultPort        = 5432
	DefaultDbName      = "postgres"
	DefaultUser        = "postgres"
	DefaultPassword    = "password"
	DefaultSSLMode     = "allow"
	DefaultLogLevel    = gorm_logger.Error
	DefaultAutoMigrate = false
)

type Config struct {
	Host        string
	Port        int
	DBName      string
	User        string
	Password    string
	SSLMode     string
	LogLevel    gorm_logger.LogLevel
	AutoMigrate bool
}

type Database struct {
	logger *zap.Logger
	config *Config

	scope string
	db    *gorm.DB
}

type Params struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
}

func InitiateModule(scope string) fx.Option {
	return fx.Module(
		scope,
		fx.Provide(func(p Params) (*Database, error) {
			logger := p.Logger.Named("[" + scope + "]")
			config := loadConfig(scope)

			db := setupDatabaseOrFatal(config, logger)

			database := &Database{
				logger: logger,
				db:     db,
				config: config,
				scope:  scope,
			}

			return database, nil
		}),
		fx.Invoke(func(d *Database, p Params) {
			p.Lifecycle.Append(
				fx.Hook{
					OnStart: d.onStart,
					OnStop:  d.onStop,
				},
			)
		}),
	)
}

func loadConfig(scope string) *Config {
	getConfigPath := func(key string) string {
		return fmt.Sprintf("%s.%s", scope, key)
	}

	//set default values
	viper.SetDefault(getConfigPath("%s.host"), DefaultHost)
	viper.SetDefault(getConfigPath("%s.port"), DefaultPort)
	viper.SetDefault(getConfigPath("%s.dbname"), DefaultDbName)
	viper.SetDefault(getConfigPath("%s.user"), DefaultUser)
	viper.SetDefault(getConfigPath("%s.password"), DefaultPassword)
	viper.SetDefault(getConfigPath("%s.sslmode"), DefaultSSLMode)
	viper.SetDefault(getConfigPath("%s.log_level"), DefaultLogLevel)
	viper.SetDefault(getConfigPath("%s.autoMigrate"), DefaultAutoMigrate)

	return &Config{
		Host:        viper.GetString(getConfigPath("host")),
		Port:        viper.GetInt(getConfigPath("port")),
		DBName:      viper.GetString(getConfigPath("dbname")),
		User:        viper.GetString(getConfigPath("user")),
		Password:    viper.GetString(getConfigPath("password")),
		SSLMode:     viper.GetString(getConfigPath("sslmode")),
		LogLevel:    gorm_logger.Info,
		AutoMigrate: viper.GetBool(getConfigPath("autoMigrate")),
	}
}

func setupDatabaseOrFatal(config *Config, logger *zap.Logger) *gorm.DB {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
	gormConfig := &gorm.Config{
		Logger: gorm_logger.Default.LogMode(config.LogLevel),
	}

	db, err := gorm.Open(postgres.Open(connectionString), gormConfig)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	return db
}

func (d *Database) onStart(context.Context) error {
	d.logger.Info("Database initiated")

	d.printDebugLogs()

	return nil
}

func (d *Database) onStop(context.Context) error {
	dbSql, err := d.db.DB()
	if err != nil {
		d.logger.Error("Error getting DB from GORM", zap.Error(err))
		return err
	}

	err = dbSql.Close()
	if err != nil {
		d.logger.Error("Error closing DB", zap.Error(err))
	}

	d.logger.Info("Database connection stopped")
	return nil
}

func (d *Database) GetDB() *gorm.DB {
	return d.db
}

func (d *Database) printDebugLogs() {
	//* Debug Logs
	d.logger.Debug("----- Database Configuration -----")
	d.logger.Debug("Host", zap.String("host", d.config.Host))
	d.logger.Debug("Port", zap.Int("port", d.config.Port))
	d.logger.Debug("DBName", zap.String("dbname", d.config.DBName))
	d.logger.Debug("User", zap.String("user", d.config.User))
	d.logger.Debug("SSLMode", zap.String("sslmode", d.config.SSLMode))
}

func (d *Database) AutoMigrate(models []interface{}) {
	if d.config.AutoMigrate {
		d.logger.Info("Auto migrating database")
		err := d.db.AutoMigrate(models)
		if err != nil {
			d.logger.Error("Error auto migrating database", zap.Error(err))
		}
	}
}
