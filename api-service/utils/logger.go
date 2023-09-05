package utils

import (
	"sync"

	"go.uber.org/zap"
)

var logger *zap.Logger
var loggerLock sync.Mutex

func ZapLogger() *zap.Logger {
	if logger == nil {
		loggerLock.Lock()
		defer loggerLock.Unlock()
		if logger == nil {
			logger, _ = zap.NewProduction()
		}
	}
	return logger
}

func SugarLogger() *zap.SugaredLogger {
	sugarLogger := ZapLogger().Sugar()
	return sugarLogger
}
