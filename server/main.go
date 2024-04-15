package main

import (
	"os"

	"github.com/alsey89/hrms/internal/auth"
	"github.com/alsey89/hrms/pkg/configs"
	echo_jwt "github.com/alsey89/hrms/pkg/jwt"
	"github.com/alsey89/hrms/pkg/logger"
	"github.com/alsey89/hrms/pkg/mailer"
	"github.com/alsey89/hrms/pkg/postgres"
	"github.com/alsey89/hrms/pkg/server"

	"go.uber.org/fx"
)

var config *configs.Config

func init() {
	//system log level
	os.Setenv("LOG_LEVEL", "debug")

	// load config
	config = configs.NewConfig("SERVER", true, true)
	config.SetConfigs(map[string]interface{}{
		"server.host": "0.0.0.0",
		"server.port": 3001,

		"database.host":     "0.0.0.0",
		"database.port":     5432,
		"database.dbname":   "postgres",
		"database.user":     "postgres",
		"database.password": "password",
		"database.sslmode":  "prefer",
		"databse.loglevel":  "error",

		"mailer.host":         "smtp.gmail.com",
		"mailer.port":         587,
		"mailer.username":     "phyokyawsoe89@gmail.com",
		"mailer.app_password": "lyzo mila fxha dupi",
		"mailer.tls":          true,

		"auth_jwt.signing_key":  "thisisasecret",
		"auth_jwt.token_lookup": "cookie:jwt",
	})
}

func main() {

	app := fx.New(
		fx.Supply(config),

		logger.InitiateModule(),
		server.InitiateModule("server"),
		echo_jwt.InitiateModule("echo_jwt"),
		postgres.InitiateModule("database"),
		mailer.InitiateModule("mailer"),

		auth.InitiateDomain("auth"),

		fx.NopLogger,
	)

	app.Run()
}
