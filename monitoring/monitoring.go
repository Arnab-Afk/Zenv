package monitoring

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func SetupLogging() {
	logger, _ := zap.NewProduction()
	Logger = logger
	defer logger.Sync()
}
