package company

import (
	"context"
	// "net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	postgres "github.com/alsey89/gogetter/database/postgres"
	mailer "github.com/alsey89/gogetter/mail/gomail"
	server "github.com/alsey89/gogetter/server/echo"
	"github.com/alsey89/people-matter/internal/auth"
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
	Mailer    *mailer.Module

	// internal domains
	Auth *auth.Domain
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
	companyGroup := server.Group("api/v1/company")

	adminGroup := companyGroup.Group("/admin")

	// Routes
	companyGroup.GET("", c.GetCompanyHandler)
	companyGroup.POST("", c.CreateCompanyHandler)
	adminGroup.PUT("", c.UpdateCompanyHandler)
	adminGroup.DELETE("", c.DeleteCompanyHandler)

	adminGroup.POST("/department", c.CreateDepartmentHandler)
	adminGroup.PUT("/department/:department_id", c.UpdateDepartmentHandler)
	adminGroup.DELETE("/department/:department_id", c.DeleteDepartmentHandler)

	adminGroup.POST("/location", c.CreateLocationHandler)
	adminGroup.PUT("/location/:location_id", c.UpdateLocationHandler)
	adminGroup.DELETE("/location/:location_id", c.DeleteLocationHandler)

	adminGroup.POST("/position", c.CreatePositionHandler)
	adminGroup.PUT("/position/:position_id", c.UpdatePositionHandler)
	adminGroup.DELETE("/position/:position_id", c.DeletePositionHandler)

	return nil
}

func (c *Domain) onStop(ctx context.Context) error {
	c.logger.Info("Stopped APIs")

	return nil
}
