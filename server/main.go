package main

import (
	"go.uber.org/fx"

	config "github.com/alsey89/gogetter/config/viper"
	postgres "github.com/alsey89/gogetter/database/postgres"
	jwt "github.com/alsey89/gogetter/jwt/echo"
	logger "github.com/alsey89/gogetter/logging/zap"
	mailer "github.com/alsey89/gogetter/mail/gomail"
	server "github.com/alsey89/gogetter/server/echo"

	"github.com/alsey89/people-matter/internal/auth"
	"github.com/alsey89/people-matter/internal/company"
	"github.com/alsey89/people-matter/schema"
)

var configuration *config.Module

func init() {
	config.SetSystemLogLevel("debug")
	configuration = config.SetUpConfig("SERVER", "yaml")
	//! CONFIG PRECEDENCE: ENV > CONFIG FILE > FALLBACK
	configuration.SetFallbackConfigs(map[string]interface{}{
		"server.host":      "0.0.0.0",
		"server.port":      3001,
		"server.log_level": "DEV",

		// "server.allow_headers":   "*",
		// "server.allow_methods":   "*",
		"server.allow_origins":   "http://localhost:3000, http://localhost:3001",
		"server.csrf_protection": true,
		"server.csrf_secure":     false,
		"server.csrf_domain":     "localhost",

		"database.host":         "postgres",
		"database.port":         5432,
		"database.dbname":       "postgres",
		"database.user":         "postgres",
		"database.password":     "password",
		"database.sslmode":      "prefer",
		"databse.loglevel":      "error",
		"database.auto_migrate": true,

		"mailer.host":         "smtp.gmail.com",
		"mailer.port":         587,
		"mailer.username":     "phyokyawsoe89@gmail.com",
		"mailer.app_password": "lyzo mila fxha dupi",
		"mailer.tls":          true,

		"jwt_mw.signing_key":    "authsecret",
		"jwt_mw.token_lookup":   "cookie:jwt",
		"jwt_mw.signing_method": "HS256",
		"jwt_mw.exp_in_hours":   72,

		"auth.signing_key":    "confirmationsecret",
		"auth.signing_method": "HS256",
		"auth.exp_in_hours":   1,
	})
}
func main() {
	app := fx.New(
		fx.Supply(configuration),
		logger.InitiateModule(),
		server.InitiateModule("server"),
		postgres.InitiateModuleAndSchema(
			"database",
			// ...schema,
			schema.Company{},
			schema.Department{},
			schema.Location{},
			schema.User{},
			schema.ContactInfo{},
			schema.EmergencyContact{},
			schema.Position{},
			schema.UserPosition{},
			schema.Leave{},
			schema.Attendance{},
			schema.Salary{},
			schema.Payment{},
			schema.Adjustments{},
			schema.Document{},
		),
		jwt.InitiateModule("jwt_mw"),
		mailer.InitiateModule("mailer"),

		//-- Internal Domains Start --
		auth.InitiateDomain("auth"),
		company.InitiateDomain("company"),
		//-- Internal Domains End --
		fx.NopLogger,
	)
	app.Run()
}
