package logger

import (
	"ExprCalc/pkg/config"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(config *config.LoggerConfig) *zap.Logger {
	return zap.New(configurateLogger(config))
}

func configurateLogger(config *config.LoggerConfig) zapcore.Core {
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(config.Path),
		MaxSize:    config.MaxSize,
		MaxAge:     config.MaxAge,
		MaxBackups: config.MaxBackups,
	})

	var levelEnabler zap.LevelEnablerFunc
	switch config.Level {
	case "debug":
		levelEnabler = zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level >= zapcore.DebugLevel
		})
	case "warn":
		levelEnabler = zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level >= zapcore.WarnLevel
		})
	default:
		levelEnabler = zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level >= zapcore.InfoLevel
		})
	}

	consoleWriter := zapcore.Lock(os.Stdout)
	jsonEndcoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	consoleEndcoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	return zapcore.NewTee(
		zapcore.NewCore(jsonEndcoder, fileWriter, levelEnabler),
		zapcore.NewCore(consoleEndcoder, consoleWriter, levelEnabler),
	)
}
