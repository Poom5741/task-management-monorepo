package logger

import (
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func Init(level string) {
	var config zap.Config
	if level == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	logger, _ := config.Build()
	log = logger.Sugar()
}

func Info(msg string, args ...interface{}) {
	log.Infof(msg, args...)
}

func Error(msg string, args ...interface{}) {
	log.Errorf(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	log.Debugf(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	log.Warnf(msg, args...)
}
