package logger

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// to be provided to the fx framework
var logger *zap.Logger

// default values
const (
	DefaultSystemLogLevel = zap.InfoLevel
)

//! MODULE ---------------------------------------------------------------

// Provides the logger to the fx framework
func InjectModule(scope string) fx.Option {
	return fx.Options(
		fx.Provide(setupLogger),
	)
}

// Instantiate the logger without using the fx framework
func NewLogger() *zap.Logger {
	return setupLogger()
}

// ! INTERNAL ---------------------------------------------------------------

func setupLogger() *zap.Logger {
	logLevel := setupLevel()
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(NewCustomEncoderConfig()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		logLevel,
	)

	if viper.GetString("global.log_level") == zap.DebugLevel.String() || viper.GetString("global.log_level") == zap.DebugLevel.CapitalString() {
		logger = zap.New(core, zap.AddCaller(), zap.Development())
	} else {
		logger = zap.New(core)
	}

	zap.ReplaceGlobals(logger)

	logger.Named("[logger]").Info(fmt.Sprintf("System log level is set to \"%s\"\n", logLevel.Level().CapitalString()))

	return logger
}

func setupLevel() zap.AtomicLevel {
	var logLevel zapcore.Level

	switch viper.GetString("global.log_level") {
	case zap.DebugLevel.String(), zap.DebugLevel.CapitalString():
		logLevel = zap.DebugLevel
	case zap.InfoLevel.String(), zap.InfoLevel.CapitalString():
		logLevel = zap.InfoLevel
	case zap.WarnLevel.String(), zap.WarnLevel.CapitalString():
		logLevel = zap.WarnLevel
	case zap.ErrorLevel.String(), zap.ErrorLevel.CapitalString():
		logLevel = zap.ErrorLevel
	case zap.DPanicLevel.String(), zap.DPanicLevel.CapitalString():
		logLevel = zap.DPanicLevel
	case zap.PanicLevel.String(), zap.PanicLevel.CapitalString():
		logLevel = zap.PanicLevel
	case zap.FatalLevel.String(), zap.FatalLevel.CapitalString():
		logLevel = zap.FatalLevel
	default:
		logLevel = DefaultSystemLogLevel
	}

	return zap.NewAtomicLevelAt(logLevel)
}

func NewCustomEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
