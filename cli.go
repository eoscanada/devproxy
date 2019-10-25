package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func errorCheck(logger *zap.Logger, prefix string, err error, fields ...zapcore.Field) {
	if err != nil {
		logger.Error(prefix, append([]zapcore.Field{zap.Error(err)}, fields...)...)
		logger.Sync()
		zlog.Sync()
		os.Exit(1)
	}
}
