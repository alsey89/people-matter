package main

import (
	"github.com/alsey89/people-matter/internal/schema"
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

		//* Migration -------------------------------------------------------------
		fx.Invoke(func(m *pgconn.Module) {
			m.ApplySchema(
				true,
				schema.Company{},
				schema.User{},
				schema.Location{},
				schema.Position{},
				schema.Permission{},
				schema.UserPosition{},
				schema.PositionPermission{},
				schema.Compensation{},
				schema.Payment{},
			)
		}),
		//* fx logs ---------------------------------------------------------------
		fx.NopLogger,
	)
	app.Run()
}
