package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

func InitLogger() {
	// 创建日志目录
	logDir := "log"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("create log directory failed: " + err.Error())
	}

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{filepath.Join(logDir, "rocketmq2nats.log")}
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	config.DisableCaller = true

	logger, err := config.Build()
	if err != nil {
		panic("init logger failed: " + err.Error())
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	zap.L().Info("init logger success")
}
