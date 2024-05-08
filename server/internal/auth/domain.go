package auth

import (
	"context"
	// "net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	postgres "github.com/alsey89/gogetter/database/postgres"
	jwt "github.com/alsey89/gogetter/jwt/echo"
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
	JWT       *jwt.JWT
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

	// d.AddDefaultData(ctx)

	// Router
	server := d.params.Server.GetServer()
	authGroup := server.Group("api/v1/auth")

	// authGroup.POST("/signup", d.SignupHandler)
	authGroup.POST("/signin", d.SigninHandler)
	authGroup.POST("/signout", d.SignoutHandler)

	authGroup.GET("/check", d.CheckAuth)
	authGroup.GET("/csrf", d.GetCSRFToken)

	authGroup.GET("/confirmation", d.ConfirmationHandler)

	return nil
}

func (d *Domain) onStop(ctx context.Context) error {
	d.logger.Info("Stopped APIs")

	return nil
}

// func (d *Domain) AddDefaultData(ctx context.Context) {
// 	db := d.params.Database.GetDB()

// 	u := User{}
// 	id := uuid.MustParse(viper.GetString(d.getConfigPath("super_admin_id")))
// 	result := db.Where("id = ?", id).First(&u)
// 	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		d.logger.Fatal(result.Error.Error())
// 		return
// 	}

// 	if result.RowsAffected > 0 {
// 		return
// 	}

// 	hashedPwd, err := HashPassword(viper.GetString(d.getConfigPath("super_admin_password")))
// 	if err != nil {
// 		d.logger.Error(err.Error())
// 		return
// 	}

// 	data := User{
// 		BaseModel: BaseModel{
// 			ID: id,
// 		},
// 		Email:    viper.GetString(d.getConfigPath("super_admin_email")),
// 		Password: hashedPwd,
// 	}
// 	if err := db.Create(&data).Error; err != nil {
// 		d.logger.Error(err.Error())
// 	}
// }
