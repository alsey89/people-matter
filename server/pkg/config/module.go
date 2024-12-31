package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultConfigFileType = "yaml"
	defaultEnvPrefix      = ""
	defaultPrintSettings  = false
	defaultConfigFilePath = "./"
)

/*
Sets up the configuration manager.
prefix: Prefix for environment variables. Defaults to "".
configFileType: Type of config file. Defaults to "yaml".
configFilePath: Path to config file. Defaults to "./".
configFilePath is relative to where the function is called, usually main.go.
*/
func SetUpConfig(prefix string, configFileType string, configFilePath string) {

	// configure how viper reads environment variables
	if prefix != "" {
		viper.SetEnvPrefix(prefix)
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// set config file type
	validatedConfigFileType := validateConfigFileType(configFileType)
	viper.SetConfigType(validatedConfigFileType)

	// set config file paths
	if configFilePath == "" {
		viper.AddConfigPath(defaultConfigFilePath)
	} else {
		viper.AddConfigPath(configFilePath)
	}

	readConfigWithOptionalOverride("config", "config.override", validatedConfigFileType, configFilePath)

	if viper.GetString("global.log_level") == "DEBUG" || viper.GetString("global.log_level") == "debug" {
		logConfigurations(prefix, validatedConfigFileType)
	}

}

//! INTERNAL ---------------------------------------------------------

func readConfigWithOptionalOverride(baseConfig string, overrideConfig string, configFileType string, configFilePath string) {
	var err error

	//check for presence of base config file
	_, err = os.Stat(configFilePath + "/" + baseConfig + "." + configFileType)
	if err != nil {
		log.Printf("Base config -- %s.%s -- not found.", baseConfig, configFileType)
		return
	}

	// read the base config file
	viper.SetConfigName(baseConfig)
	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("Error reading base config file, %s.", err)
		return
	}
	log.Printf("Base config applied from -- %s\n --", viper.ConfigFileUsed())

	// check for presence of override config file
	_, err = os.Stat(configFilePath + "/" + overrideConfig + "." + configFileType)
	if err != nil {
		log.Printf("Skipping override. Override config -- %s.%s -- not found.", overrideConfig, configFileType)
		return
	}

	// merge the override config file with the base config file
	viper.SetConfigName(overrideConfig)
	err = viper.MergeInConfig()
	if err != nil {
		log.Printf("Error reading override config file, %s.", err)
		return
	}
	log.Printf("Override config applied from -- %s\n --", viper.ConfigFileUsed())
}

func validateConfigFileType(configFileType string) string {
	if configFileType == "" {
		return defaultConfigFileType
	}

	validConfigFileTypes := map[string]bool{
		"json":    true,
		"toml":    true,
		"yaml":    true,
		"hcl":     true,
		"envfile": true,
	}

	lowerCaseConfigFileType := strings.ToLower(configFileType)

	if !validConfigFileTypes[lowerCaseConfigFileType] {
		log.Printf("Invalid config file type: %s. Defaulting to %s", lowerCaseConfigFileType, defaultConfigFileType)
		return defaultConfigFileType
	}

	return lowerCaseConfigFileType
}

func logConfigurations(prefix string, configFileType string) {
	log.Println("----- Config Setup -----")
	log.Printf("|| Prefix   %s", prefix)
	log.Printf("|| replacer %s", ". -> _")
	log.Printf("|| autoEnv  %s", "true")
	log.Printf("|| name     %s", "config AND config.override[OPTIONAL]")
	log.Printf("|| paths    %s", defaultConfigFilePath)
	log.Printf("|| fileType %s", configFileType)
	log.Println("------------------------")
}

//! External ---------------------------------------------------------

// Sets values for the config keys that are not provided.
// ! IMPORTANT: Set is absolute. Run this function last to avoid overriding.
func SetFallbackConfigs(configs map[string]interface{}) {
	for k, v := range configs {
		if !viper.IsSet(k) {
			viper.Set(k, v)
		}
	}
}

// Sets value for global.log_level, defaults to "INFO"
func OverrideGlobalLogLevel(logLevel string) {
	if logLevel == "" {
		viper.Set("global.log_level", "dev")
		return
	}

	viper.Set("global.log_level", logLevel)
}
