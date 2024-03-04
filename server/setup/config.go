package setup

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func InitConfig() {
	// Set the base name of the config file, without the file extension.
	viper.SetConfigName("config") // default config file name
	viper.SetConfigType("yaml")   // config file type
	viper.AddConfigPath(".")      // look for config in the working directory

	// Use environment variables where available. Viper can read a set of predefined env variables.
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // replace dots with underscores in env vars

	// Attempt to load the default configuration file.
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Error reading config file, %s", err)
	}

	// Override with values from the local configuration file if it exists.
	localConfig := "config.override"
	_, err = os.Stat(localConfig + ".yaml")
	if err == nil {
		viper.SetConfigName(localConfig) // specific local config file name
		err = viper.MergeInConfig()      // merges config values, overriding with local config
		if err != nil {
			log.Printf("Error merging local config file, %s", err)
		}
	}

	log.Printf("Config loaded: %v", viper.ConfigFileUsed())
}
