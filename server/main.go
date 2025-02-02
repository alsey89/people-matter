package main

import (
	"github.com/alsey89/people-matter/internal/schema"
	"github.com/alsey89/people-matter/internal/transmail"
	"github.com/alsey89/people-matter/pkg/config"
	"github.com/alsey89/people-matter/pkg/logger"
	"github.com/alsey89/people-matter/pkg/pgconn"
	"github.com/alsey89/people-matter/pkg/server"

	"go.uber.org/fx"
)

func init() {
	//! CONFIG PRECEDENCE: OVERRIDE > ENV > CONFIG FILE > FALLBACK

	// *OPTIONAL* OVERRIDE GLOBAL LOG LEVEL, INTENDED FOR DEVELOPMENT
	// config.OverrideGlobalLogLevel("debug")

	config.SetUpConfig("SERVER", "yaml", "./")
}

func main() {
	app := fx.New(
		//* Modules ---------------------------------------------------------------
		logger.InjectModule("logger"),
		pgconn.InjectModule("database"),
		server.InjectModule("server"),
		//* Domains ---------------------------------------------------------------
		transmail.InjectDomain("transmail"),
		//* Migration -------------------------------------------------------------
		fx.Invoke(func(m *pgconn.Module) {
			m.ApplySchema(
				true,
				schema.Adjustment{},
				schema.Bonus{},
				schema.Company{},
				schema.Compensation{},
				schema.Document{},
				schema.Expense{},
				schema.Location{},
				schema.Payment{},
				schema.Permission{},
				schema.Position{},
				schema.PositionPermission{},
				schema.User{},
				schema.UserPosition{},
			)
		}),
		//* fx logs ---------------------------------------------------------------
		fx.NopLogger,
	)
	app.Run()
}
