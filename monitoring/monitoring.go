package monitoring

import (
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Logger *zap.Logger

// SetupLogging initializes the global logger
func SetupLogging() error {
	// Ensure logs directory exists
	logPath := "logs"
	if err := os.MkdirAll(logPath, 0755); err != nil {
		return err
	}

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{
		"stdout",
		filepath.Join(logPath, "app.log"),
	}

	var err error
	Logger, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

// RequestLogger middleware for logging HTTP requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		Logger.Info("request completed",
			zap.String("path", path),
			zap.String("method", c.Request.Method),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}

// LogError logs error details with context
func LogError(err error, context map[string]interface{}) {
	fields := make([]zap.Field, 0)
	for k, v := range context {
		fields = append(fields, zap.Any(k, v))
	}
	Logger.Error("error occurred", append(fields, zap.Error(err))...)
}

// AuditLog records security-sensitive operations
func AuditLog(userID string, action string, resource string, success bool) {
	Logger.Info("audit",
		zap.String("user_id", userID),
		zap.String("action", action),
		zap.String("resource", resource),
		zap.Bool("success", success),
		zap.Time("timestamp", time.Now()),
	)
}

// MetricsMiddleware tracks basic performance metrics
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start)
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)

		Logger.Info("metrics",
			zap.Duration("response_time", duration),
			zap.String("endpoint", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Uint64("memory_alloc", mem.Alloc),
			zap.Int("goroutines", runtime.NumGoroutine()),
			zap.String("method", c.Request.Method),
			zap.Int64("request_size", c.Request.ContentLength),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}

// LogFatal logs fatal errors and exits
func LogFatal(msg string, err error) {
	Logger.Fatal(msg,
		zap.Error(err),
		zap.Time("timestamp", time.Now()),
	)
}

// LogPanic logs panic information
func LogPanic(err interface{}) {
	Logger.Error("panic occurred",
		zap.Any("error", err),
		zap.String("stack", string(debug.Stack())),
	)
}
