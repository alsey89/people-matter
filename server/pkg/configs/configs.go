package configs

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
}

func NewConfig(prefix string) *Config {
	// From the environment
	viper.SetEnvPrefix(prefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// From config file
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./configs")

	log.Printf("\n || ***LOADING CONFIG FILES*** \n || prefix: %s\n || replacer: %s\n || autoEnv: %s\n || name: %s\n || path: %s\n", prefix, ". -> _", "true", "config", "./ OR ./configs")

	log.Printf("Reading base config file...\n")
	// read the primary config file
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

	log.Printf("Reading override config file...\n")
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
		log.Printf("Override applied: %s\n", viper.ConfigFileUsed())
	}

	config := &Config{}

	return config
}

func (config *Config) SetConfigs(configs map[string]interface{}) {

	for k, v := range configs {
		if !viper.IsSet(k) {
			viper.Set(k, v)
		}
	}
}

// func (config *Config) GetAllSettings() map[string]interface{} {
// 	return viper.AllSettings()
// }

// func (config *Config) PrintAllSettings() {

// 	log.Println("[List of current configs]")
// 	config.PrintSettings("", config.GetAllSettings())
// 	log.Println("")
// }

// func (config *Config) PrintSettings(scope string, settings map[string]interface{}) {

// 	for k, v := range settings {

// 		domain := k
// 		if len(scope) > 0 {
// 			domain = scope + "." + k
// 		}

// 		switch reflect.TypeOf(v).Kind() {
// 		case reflect.Map:
// 			config.PrintSettings(domain, v.(map[string]interface{}))
// 		default:
// 			log.Printf("%s=%v\n", domain, v)
// 		}
// 	}
// }
