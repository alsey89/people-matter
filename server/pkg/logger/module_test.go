package logger

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSetupLogger(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	viper.Set("global.log_level", "DEBUG")
	logger := setupLogger()

	assert.NotNil(t, logger)
	assert.IsType(t, &zap.Logger{}, logger)
}

func TestSetupLevel(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	t.Run("TestSetupLevelWithDebugLevel", func(t *testing.T) {
		viper.Set("global.log_level", "debug")
		level := setupLevel()
		assert.Equal(t, zap.DebugLevel, level.Level())
	})

	t.Run("TestSetupLevelWithInfoLevel", func(t *testing.T) {
		viper.Set("global.log_level", "info")
		level := setupLevel()
		assert.Equal(t, zap.InfoLevel, level.Level())
	})

	t.Run("TestSetupLevelWithWarnLevel", func(t *testing.T) {
		viper.Set("global.log_level", "warn")
		level := setupLevel()
		assert.Equal(t, zap.WarnLevel, level.Level())
	})

	t.Run("TestSetupLevelWithErrorLevel", func(t *testing.T) {
		viper.Set("global.log_level", "error")
		level := setupLevel()
		assert.Equal(t, zap.ErrorLevel, level.Level())
	})

	t.Run("TestSetupLevelWithDPanicLevel", func(t *testing.T) {
		viper.Set("global.log_level", "dpanic")
		level := setupLevel()
		assert.Equal(t, zap.DPanicLevel, level.Level())
	})

	t.Run("TestSetupLevelWithPanicLevel", func(t *testing.T) {
		viper.Set("global.log_level", "panic")
		level := setupLevel()
		assert.Equal(t, zap.PanicLevel, level.Level())
	})

	t.Run("TestSetupLevelWithFatalLevel", func(t *testing.T) {
		viper.Set("global.log_level", "fatal")
		level := setupLevel()
		assert.Equal(t, zap.FatalLevel, level.Level())
	})

	t.Run("TestSetupLevelWithInvalidLevel", func(t *testing.T) {
		viper.Set("global.log_level", "invalid")
		level := setupLevel()
		assert.Equal(t, DefaultSystemLogLevel, level.Level())
	})
}
