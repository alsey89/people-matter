package configs

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
}

func NewConfig(prefix string, printDebugLogs bool) *Config {

	//SETUP -----------------------------
	if prefix == "" {
		prefix = "APP"
	}
	viper.SetEnvPrefix(prefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	//CONFIG FILES -------------------------------------
	viper.AddConfigPath("./")
	viper.AddConfigPath("./configs")

	// read the primary config file
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			log.Printf("Error reading config file, %s", err)
		} else {
			log.Println("No base config")
		}
	} else {
		log.Printf("Base config applied: %s\n", viper.ConfigFileUsed())
	}

	// check and read the override config file if it exists
	viper.SetConfigName("config.override")
	err = viper.MergeInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			log.Printf("Error reading override config file, %s", err)
		} else {
			log.Println("No override config")
		}
	} else {
		log.Printf("Override config applied: %s\n", viper.ConfigFileUsed())
	}

	config := &Config{}

	if printDebugLogs {
		config.PrintDebugLogs(prefix)
	}

	return config
}

func (config *Config) SetConfigs(configs map[string]interface{}) {

	for k, v := range configs {
		if !viper.IsSet(k) {
			viper.Set(k, v)
		}
	}
}

func (config *Config) PrintDebugLogs(prefix string) {
	log.Println("----- Config Setup -----")
	log.Printf("|| Prefix:   %s", prefix)
	log.Printf("|| replacer: %s", ". -> _")
	log.Printf("|| autoEnv:  %s", "true")
	log.Printf("|| name:     %s", "config")
	log.Printf("|| path:     %s", "./ OR ./configs")
}
