package main

import (
	"github.com/alsey89/people-matter/internal/fsp"
	"github.com/alsey89/people-matter/internal/identity"
	"github.com/alsey89/people-matter/internal/memorial"
	"github.com/alsey89/people-matter/internal/seeder"
	"github.com/alsey89/people-matter/pkg/config"
	"github.com/alsey89/people-matter/pkg/logger"
	"github.com/alsey89/people-matter/pkg/pgconn"
	"github.com/alsey89/people-matter/pkg/server"
	"github.com/alsey89/people-matter/pkg/token"

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
		// s3conn.InjectModule("s3"),
		token.InjectModule("jwt", "jwt_auth", "jwt_email", "jwt_pw_reset"),
		// mailer.InjectModule("mailer", false),
		server.InjectModule("server"),
		//* Domains ---------------------------------------------------------------
		seeder.InjectDomain("seeder"),
		// transmail.InjectDomain("transmail"),
		identity.InjectDomain("identity"),
		fsp.InjectDomain("fsp"),
		memorial.InjectDomain("memorial"),
		//* Migration -------------------------------------------------------------
		fx.Invoke(func(m *pgconn.Module) {
			m.ApplySchema(
				true,
				// schema.FSP{},
				// schema.User{},
				// schema.Country{},
				// schema.StateProvince{},
				// schema.Memorial{},
				// schema.FSPRole{},
				// schema.UserFSPRole{},
				// schema.MemorialRole{},
				// schema.UserMemorialRole{},
				// schema.Invitation{},
				// schema.Application{},
				// schema.ContributionCondolenceElement{},
				// schema.ContributionGalleryElement{},
				// schema.ContributionStoryElement{},
				// schema.ContributionTimelineElement{},
				// schema.Export{},
			)
		}),
		//* fx logs ---------------------------------------------------------------
		fx.NopLogger,
	)
	app.Run()
}
