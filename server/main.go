package main

import (
	"verve-hrms/pkg/configs"
	"verve-hrms/pkg/logger"
	"verve-hrms/pkg/postgres_connector"
	"verve-hrms/pkg/server"

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
	)

	app.Run()
}
