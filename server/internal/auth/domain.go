package auth

import (
	"context"
	// "net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/alsey89/hrms/pkg/postgres_connector"
	"github.com/alsey89/hrms/pkg/server"

	"github.com/alsey89/hrms/internal/schema"
)

type Domain struct {
	params Params
	scope  string

	logger *zap.Logger
	User   *User
}

type Params struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
	Server    *server.HTTPServer
	Database  *postgres_connector.Database
}

func Module(scope string) fx.Option {

	var a *Domain

	return fx.Options(
		fx.Provide(func(p Params) *Domain {

			a := &Domain{
				params: p,
				scope:  scope,
				logger: p.Logger.Named(scope),
				User:   &User{},
			}

			return a
		}),
		fx.Populate(&a),
		fx.Invoke(func(p Params) {

			p.Lifecycle.Append(
				fx.Hook{
					OnStart: a.onStart,
					OnStop:  a.onStop,
				},
			)
		}),
	)

}

func (a *Domain) onStart(ctx context.Context) error {

	a.logger.Info("Starting APIs")

	// Auto migration
	db := a.params.Database.GetDB()
	db.AutoMigrate(&schema.User{})

	// a.AddDefaultData(ctx)

	// Router
	server := a.params.Server.GetServer()

	authGroup := server.Group("/auth")

	authGroup.POST("/signin", Signin)
	authGroup.POST("/signup", Signup)
	authGroup.POST("/signout", Signout)

	return nil
}

func (a *Domain) onStop(ctx context.Context) error {
	a.logger.Info("Stopped APIs")

	return nil
}

// func (a *Domain) AddDefaultData(ctx context.Context) {
// 	db := a.params.Database.GetDB()

// 	u := User{}
// 	id := uuid.MustParse(viper.GetString(a.getConfigPath("super_admin_id")))
// 	result := db.Where("id = ?", id).First(&u)
// 	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		a.logger.Fatal(result.Error.Error())
// 		return
// 	}

// 	if result.RowsAffected > 0 {
// 		return
// 	}

// 	hashedPwd, err := HashPassword(viper.GetString(a.getConfigPath("super_admin_password")))
// 	if err != nil {
// 		a.logger.Error(err.Error())
// 		return
// 	}

// 	data := User{
// 		BaseModel: BaseModel{
// 			ID: id,
// 		},
// 		Email:    viper.GetString(a.getConfigPath("super_admin_email")),
// 		Password: hashedPwd,
// 	}
// 	if err := db.Create(&data).Error; err != nil {
// 		a.logger.Error(err.Error())
// 	}
// }
