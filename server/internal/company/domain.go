package company

import (
	"context"
	// "net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	postgres "github.com/alsey89/gogetter/database/postgres"
	server "github.com/alsey89/gogetter/server/echo"
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
	Database  *postgres.Database
}

func InitiateDomain(scope string) fx.Option {

	var c *Domain

	return fx.Options(
		fx.Provide(func(p Params) *Domain {

			c := &Domain{
				params: p,
				scope:  scope,
				logger: p.Logger.Named("[" + scope + "]"),
			}

			return c
		}),
		fx.Populate(&c),
		fx.Invoke(func(p Params) {

			p.Lifecycle.Append(
				fx.Hook{
					OnStart: c.onStart,
					OnStop:  c.onStop,
				},
			)
		}),
	)

}

func (c *Domain) onStart(ctx context.Context) error {

	c.logger.Info("Starting APIs")

	// Router
	server := c.params.Server.GetServer()
	companyGroup := server.Group("/company")

	// Routes
	companyGroup.GET("/", c.GetCompanyHandler)

	companyGroup.POST("/", c.CreateCompanyHandler)
	companyGroup.PUT("/", c.UpdateCompanyHandler)
	companyGroup.DELETE("/", c.DeleteCompanyHandler)

	companyGroup.POST("/department", c.CreateDepartmentHandler)
	companyGroup.PUT("/department/:department_id", c.UpdateDepartmentHandler)
	companyGroup.DELETE("/department/:department_id", c.DeleteDepartmentHandler)

	companyGroup.POST("/location", c.CreateLocationHandler)
	companyGroup.PUT("/location/:location_id", c.UpdateLocationHandler)
	companyGroup.DELETE("/location/:location_id", c.DeleteLocationHandler)

	companyGroup.POST("/position", c.CreatePositionHandler)
	companyGroup.PUT("/position/:position_id", c.UpdatePositionHandler)
	companyGroup.DELETE("/position/:position_id", c.DeletePositionHandler)

	return nil
}

func (c *Domain) onStop(ctx context.Context) error {
	c.logger.Info("Stopped APIs")

	return nil
}
