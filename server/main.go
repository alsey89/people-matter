package main

import (
	"os"

	"github.com/alsey89/hrms/internal/auth"
	"github.com/alsey89/hrms/pkg/configs"
	"github.com/alsey89/hrms/pkg/logger"
	"github.com/alsey89/hrms/pkg/mailer"
	"github.com/alsey89/hrms/pkg/postgres"
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

		"mailer.host":         "smtp.gmail.com",
		"mailer.port":         587,
		"mailer.username":     "phyokyawsoe89@gmail.com",
		"mailer.app_password": "lyzo mila fxha dupi",
		"mailer.tls":          true,
	})

	os.Setenv("LOG_LEVEL", "info")

	app := fx.New(
		fx.Supply(config),

		logger.InitiateModule(),
		server.InitiateModule("server"),
		postgres.InitiateModule("database"),
		mailer.InitiateModule("mailer"),

		auth.InitiateModule("auth"),

		fx.NopLogger,
	)

	app.Run()
}
