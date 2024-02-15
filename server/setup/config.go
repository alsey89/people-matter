package setup

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func InitConfig() {
	// Attempt to load the local configuration file
	viper.SetConfigName("config.local")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// if local config file does not exist, use the default config file
	_, err := os.Stat("./config.local.yaml")
	if os.IsNotExist(err) {
		viper.SetConfigName("config")
	}

	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	log.Printf("Config loaded: %v", viper.ConfigFileUsed())
}
