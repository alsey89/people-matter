package logger

import (
	"fmt"
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

const (
	defaultLogLevel = zap.InfoLevel
)

type Params struct {
	fx.In
}

func InitiateModule() fx.Option {
	return fx.Options(
		fx.Provide(SetupLogger),
	)
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

func SetupLogger() *zap.Logger {
	debugLevel := setupLevel()
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(NewCustomEncoderConfig()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		debugLevel,
	)

	if os.Getenv("DEBUG_MODE") == "debug" {
		logger.Info(fmt.Sprintf("Debug mode is set to \"%s\"\n", debugLevel.String()))
		logger = zap.New(core, zap.AddCaller(), zap.Development())
	} else {
		logger = zap.New(core)
	}

	zap.ReplaceGlobals(logger)

	logger.Named("[logger]").Info(fmt.Sprintf("Debug level is set to \"%s\"\n", debugLevel.String()))

	return logger
}

func GetLogger() *zap.Logger {
	return logger
}

func setupLevel() zap.AtomicLevel {

	logLevel := defaultLogLevel

	switch os.Getenv("LOG_LEVEL") {
	case zap.DebugLevel.String():
		logLevel = zap.DebugLevel
	case zap.InfoLevel.String():
		logLevel = zap.InfoLevel
	case zap.WarnLevel.String():
		logLevel = zap.WarnLevel
	case zap.ErrorLevel.String():
		logLevel = zap.ErrorLevel
	case zap.DPanicLevel.String():
		logLevel = zap.DPanicLevel
	case zap.PanicLevel.String():
		logLevel = zap.PanicLevel
	case zap.FatalLevel.String():
		logLevel = zap.FatalLevel
	default:
		logLevel = defaultLogLevel
	}

	return zap.NewAtomicLevelAt(logLevel)
}
