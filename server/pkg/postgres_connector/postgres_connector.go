package postgres_connector

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
	DefaultHost     = "0.0.0.0"
	DefaultPort     = 5432
	DefaultDbName   = "postgres"
	DefaultUser     = "postgres"
	DefaultPassword = "password"
	DefaultSSLMode  = "allow"
	DefaultLogLevel = gorm_logger.Error
)

type DBConfig struct {
	Host     string
	Port     int
	DBName   string
	User     string
	Password string
	SSLMode  string
}

type Database struct {
	logger *zap.Logger
	db     *gorm.DB
	scope  string
}

type DBParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
}

func DBModule(scope string) fx.Option {
	return fx.Module(
		scope,
		fx.Provide(func(p DBParams) *Database {
			logger := p.Logger.Named(scope)
			config := loadDBConfig(scope)

			dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
				config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				logger.Fatal("Failed to connect to database", zap.Error(err))
			}

			database := &Database{
				logger: logger,
				db:     db,
				scope:  scope,
			}

			return database
		}),
		fx.Invoke(func(d *Database, p DBParams) {
			p.Lifecycle.Append(
				fx.Hook{
					OnStart: d.onStart,
					OnStop:  d.onStop,
				},
			)
		}),
	)
}

func loadDBConfig(scope string) *DBConfig {
	getConfigWithDefault := func(key string, defaultVal interface{}) interface{} {
		scopedKey := fmt.Sprintf("%s.%s", scope, key)
		if viper.IsSet(scopedKey) {
			return viper.Get(scopedKey)
		}
		return defaultVal
	}

	return &DBConfig{
		Host:     getConfigWithDefault("host", DefaultHost).(string),
		Port:     getConfigWithDefault("port", DefaultPort).(int),
		DBName:   getConfigWithDefault("dbname", DefaultDbName).(string),
		User:     getConfigWithDefault("user", DefaultUser).(string),
		Password: getConfigWithDefault("password", DefaultPassword).(string),
		SSLMode:  getConfigWithDefault("sslmode", DefaultSSLMode).(string),
	}
}

func (d *Database) onStart(context.Context) error {
	d.logger.Info("Starting database connection", zap.String("scope", d.scope))
	// todo: add startup logic if applicable

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
