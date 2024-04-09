package main

import (
	"github.com/alsey89/hrms/internal/auth"
	"github.com/alsey89/hrms/pkg/configs"
	"github.com/alsey89/hrms/pkg/logger"
	"github.com/alsey89/hrms/pkg/postgres_connector"
	"github.com/alsey89/hrms/pkg/server"

	"go.uber.org/fx"
)

var config *configs.Config

func init() {
	config = configs.NewConfig("SERVER")
}

func main() {
	config.SetConfigs(map[string]interface{}{
		"server.host": "0.0.0.0",
		"server.port": 3001,

		"database.host":     "0.0.0.0",
		"database.port":     5432,
		"database.dbname":   "postgres",
		"database.user":     "postgres",
		"database.password": "password",
		"database.sslmode":  "prefer",
	})

	app := fx.New(
		fx.Supply(config),

		logger.Module(),
		server.Module("server"),
		postgres_connector.Module("database"),

		auth.Module("auth"),
	)

	app.Run()
}
