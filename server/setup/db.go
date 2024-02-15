package setup

import (
	"fmt"
	"verve-hrms/internal/schema"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var client *gorm.DB

func GetClient() *gorm.DB {
	if client != nil {
		return client
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		viper.GetString("db.host"),
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetString("db.dbname"),
		viper.GetString("db.port"),
		viper.GetString("db.sslmode"),
	)

	var err error
	client, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true, //* this is needed to translate postgres errors to gorm errors
	})
	if err != nil {
		panic(err)
	}

	err = client.AutoMigrate(
		&schema.User{},
		&schema.ContactInfo{},
		&schema.EmergencyContact{},
		&schema.Company{},
		&schema.Department{},
		&schema.Location{},
		&schema.Job{},
		&schema.AssignedJob{},
		&schema.Salary{},
		&schema.Payment{},
		&schema.Leave{},
		&schema.Attendance{},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to PostgreSQL!")

	return client
}
