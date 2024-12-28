package config

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSetupConfig(t *testing.T) {
	baseConfigName := "config"
	baseConfigFileType := "yaml"

	// write base config file with a test value
	os.WriteFile(baseConfigName+"."+baseConfigFileType, []byte("test_key: new_test_value"), 0644)
	defer os.RemoveAll(baseConfigName + "." + baseConfigFileType)

	t.Run("TestWithEnvPrefix", func(t *testing.T) {
		defer viper.Reset()
		os.Setenv("SERVER_ENV_KEY", "SERVER_ENV_VALUE")
		defer os.Unsetenv("SERVER_ENV_KEY")

		SetUpConfig("SERVER", baseConfigFileType, "./")

		assert.Equal(t, "SERVER_ENV_VALUE", viper.GetString("env_key"))
	})

	t.Run("TestWithoutEnvPrefix", func(t *testing.T) {
		defer viper.Reset()
		os.Setenv("ENV_KEY", "ENV_VALUE")
		defer os.Unsetenv("ENV_KEY")

		SetUpConfig("", baseConfigFileType, "./")

		assert.Equal(t, "ENV_VALUE", viper.GetString("env_key"))
	})

	t.Run("TestCustomFileType", func(t *testing.T) {
		defer viper.Reset()
		customConfigName := "config"
		customConfigFileType := "json"

		// write custom config file with a test value
		os.WriteFile(customConfigName+"."+customConfigFileType, []byte(`{"custom_key": "custom_value"}`), 0644)
		defer os.RemoveAll(customConfigName + "." + customConfigFileType)

		viper.SetConfigType(customConfigFileType)
		viper.AddConfigPath("./")

		SetUpConfig("", customConfigFileType, "./")

		assert.Equal(t, "custom_value", viper.GetString("custom_key"))

		// assert that the config file used is the custom config file
		assert.Contains(t, viper.ConfigFileUsed(), customConfigFileType)
	})

	t.Run("TestCustomConfigPath", func(t *testing.T) {
		defer viper.Reset()
		tempDir, _ := os.MkdirTemp("./test", "config_test")
		defer os.RemoveAll(tempDir)

		os.WriteFile(tempDir+"/config.yaml", []byte("custom_key: custom_value"), 0644)

		SetUpConfig("", baseConfigFileType, tempDir)

		assert.Equal(t, tempDir, viper.ConfigFileUsed())
	})
}

func TestReadConfigWithOverride(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	baseConfigName := "config"
	overrideConfigName := "config.override"

	getConfigNameFromPath := func(configFilePath string) string {
		configFilePathArray := strings.Split(configFilePath, "/")
		configFileName := configFilePathArray[len(configFilePathArray)-1]

		return configFileName
	}

	t.Run("TestReadConfigWithNoConfig", func(t *testing.T) {
		viper.Reset()
		defer viper.Reset()

		//set a test key
		viper.SetDefault("test_key", "test_value")

		viper.SetConfigType("yaml")

		// Read config with no config files
		readConfigWithOptionalOverride(baseConfigName, overrideConfigName, "yaml", "./")

		// Assert that no config file was used
		assert.Equal(t, "", viper.ConfigFileUsed())

		// Assert that the test key is still set
		assert.Equal(t, "test_value", viper.GetString("test_key"))
	})

	t.Run("TestReadConfigWithBaseConfigOnly", func(t *testing.T) {
		viper.Reset()
		defer os.Remove(baseConfigName + ".yaml")
		defer viper.Reset()

		//set a default test key
		viper.SetDefault("test_key", "test_value")

		// write base config file with new test value
		os.WriteFile(baseConfigName+".yaml", []byte("test_key: new_test_value"), 0644)

		viper.SetConfigType("yaml")
		viper.AddConfigPath("./")

		readConfigWithOptionalOverride(baseConfigName, overrideConfigName, "yaml", "./")

		usedConfigFileName := getConfigNameFromPath(viper.ConfigFileUsed())
		assert.Equal(t, baseConfigName+".yaml", usedConfigFileName, "Base config file should be used")

		// check if test key is overwritten by the new value in the base config file
		assert.Equal(t, "new_test_value", viper.GetString("test_key"))
	})

	t.Run("TestReadConfigWithOverrideConfig", func(t *testing.T) {
		viper.Reset()
		defer viper.Reset()
		defer os.Remove("config.yaml")
		defer os.Remove("config.override.yaml")

		//set a default test key
		viper.SetDefault("test_key", "test_value")

		// write base and override config files with new test values
		os.WriteFile(baseConfigName+".yaml", []byte("test_key: new_test_value"), 0644)
		os.WriteFile(overrideConfigName+".yaml", []byte("test_key: override_test_value"), 0644)

		viper.SetConfigType("yaml")
		viper.AddConfigPath("./")

		readConfigWithOptionalOverride(baseConfigName, overrideConfigName, "yaml", "./")

		usedConfigFileName := getConfigNameFromPath(viper.ConfigFileUsed())

		assert.Equal(t, overrideConfigName+".yaml", usedConfigFileName, "Override config file should be used")

		// check if test key is overwritten by the new value in the override config file
		assert.Equal(t, "override_test_value", viper.GetString("test_key"))
	})

	t.Run("TestReadConfigWithPartialConfigs", func(t *testing.T) {
		viper.Reset()
		defer viper.Reset()
		defer os.Remove("config.yaml")
		defer os.Remove("config.override.yaml")

		//set a default test key
		viper.SetDefault("test_key1", "test_value1")
		viper.SetDefault("test_key2", "test_value2")
		viper.SetDefault("fallback_key", "fallback_value")

		// write base and override config files with new test values
		os.WriteFile(baseConfigName+".yaml", []byte("test_key1: new_test_value1\ntest_key2: new_test_value2"), 0644)
		os.WriteFile(overrideConfigName+".yaml", []byte("test_key1: override_test_value1"), 0644)

		viper.SetConfigType("yaml")
		viper.AddConfigPath("./")

		readConfigWithOptionalOverride(baseConfigName, overrideConfigName, "yaml", "./")

		// check if default keys are overwritten if they exist in the config file
		assert.NotEqual(t, "test_value1", viper.GetString("test_key1"))
		assert.NotEqual(t, "test_value2", viper.GetString("test_key2"))
		// check if default keys are not overwritten if they do not exist in the config file
		assert.Equal(t, "fallback_value", viper.GetString("fallback_key"))
		// check if value from base config file is overwritten if it exists in the override config file
		assert.Equal(t, "override_test_value1", viper.GetString("test_key1"))
		// check if value from base config file is not overwritten if it does not exist in the override config file
		assert.Equal(t, "new_test_value2", viper.GetString("test_key2"))
	})

}

func TestValidateConfigFileType(t *testing.T) {
	t.Run("TestValidConfigFileType", func(t *testing.T) {
		configFileType := "yaml"
		expectedResult := "yaml"
		result := validateConfigFileType(configFileType)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("TestInvalidConfigFileType", func(t *testing.T) {
		configFileType := "invalid"
		expectedResult := "yaml"
		result := validateConfigFileType(configFileType)
		assert.Equal(t, expectedResult, result)
	})
}

func TestSetSystemLogLevel(t *testing.T) {
	SetUpConfig("", "yaml", "./")

	t.Run("TestSetSystemLogLevel", func(t *testing.T) {
		OverrideGlobalLogLevel("debug")

		assert.Equal(t, "debug", viper.GetString("global.log_level"))
	})

	t.Run("TestSetSystemLogLevelWithEmptyLogLevel", func(t *testing.T) {
		OverrideGlobalLogLevel("")

		assert.Equal(t, "debug", viper.GetString("global.log_level"))
	})
}

func TestSetFallbackConfigs(t *testing.T) {
	SetUpConfig("", "yaml", "./")

	t.Run("TestSetFallbackConfigs", func(t *testing.T) {
		viper.Reset()
		defer viper.Reset()

		fallbackConfigs := map[string]interface{}{
			"global.log_level": "DEBUG",
			"database.host":    "localhost",
			"database.port":    5432,
		}

		SetFallbackConfigs(fallbackConfigs)

		// Assert that the fallback configs are set
		assert.Equal(t, "DEBUG", viper.GetString("global.log_level"))
		assert.Equal(t, "localhost", viper.GetString("database.host"))
		assert.Equal(t, 5432, viper.GetInt("database.port"))
	})

	t.Run("TestSetFallbackConfigsWithExistingValues", func(t *testing.T) {
		viper.Reset()
		defer viper.Reset()

		// SetDefault has a lower precedence than Set
		// Using Set after SetDefault would normally override the default value
		viper.SetDefault("global.log_level", "INFO")
		viper.SetDefault("database.host", "example.com")

		fallbackConfigs := map[string]interface{}{
			"global.log_level": "DEBUG",
			"database.port":    5432,
		}

		SetFallbackConfigs(fallbackConfigs)

		// SetFallbackConfigs should not override the SetDefault values
		assert.Equal(t, "INFO", viper.GetString("global.log_level"))
		assert.Equal(t, "example.com", viper.GetString("database.host"))
		assert.Equal(t, 5432, viper.GetInt("database.port"))
	})
}
