package user

import (
	"context"
	// "net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	postgres "github.com/alsey89/gogetter/database/postgres"
	server "github.com/alsey89/gogetter/server/echo"
	"github.com/alsey89/people-matter/internal/common"
)

type Domain struct {
	params Params
	scope  string
	logger *zap.Logger
}

type Params struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
	Server    *server.HTTPServer
	Database  *postgres.Module
}

func InitiateDomain(scope string) fx.Option {

	var d *Domain

	return fx.Options(
		fx.Provide(func(p Params) *Domain {

			d := &Domain{
				params: p,
				scope:  scope,
				logger: p.Logger.Named("[" + scope + "]"),
			}

			return d
		}),
		fx.Populate(&d),
		fx.Invoke(func(p Params) {

			p.Lifecycle.Append(
				fx.Hook{
					OnStart: d.onStart,
					OnStop:  d.onStop,
				},
			)
		}),
	)

}

func (d *Domain) onStart(ctx context.Context) error {

	d.logger.Info("Starting APIs")

	// Router
	server := d.params.Server.GetServer()
	userGroup := server.Group("/user")

	// Role Precedence: Admin > Manager > User

	//Admin group, requires admin role
	adminGroup := userGroup.Group("/admin")
	adminGroup.Use(common.MustBeAdmin)
	adminGroup.GET("/", d.GetAllUsersHandler)
	// adminGroup.GET("/:id", d.GetUserHandler)
	// adminGroup.POST("/", d.CreateUserHandler)
	// adminGroup.PUT("/:id", d.UpdateUserHandler)
	// adminGroup.DELETE("/:id", d.DeleteUserHandler)

	//Manager group, requires manager role or above
	managerGroup := userGroup.Group("/manager")
	managerGroup.Use(common.MustBeManager)
	//todo: add middleware to check if user is manager
	managerGroup.GET("/", d.GetAllLocationUsersHandler)
	// managerGroup.GET("/:user_id", d.GetLocationUserHandler)
	// managerGroup.POST("/:user_id", d.CreateLocationUserHandler)
	// managerGroup.PUT("/:user_id", d.UpdateLocationUserHandler)
	// managerGroup.DELETE("/:user_id", d.DeleteLocationUserHandler)

	//User
	userGroup.GET("/me", d.GetCurrentUserHandler)
	// userGroup.PUT("/me", d.UpdateCurrentUserHandler)

	return nil
}

func (d *Domain) onStop(ctx context.Context) error {
	d.logger.Info("Stopped APIs")

	return nil
}
