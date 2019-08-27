package main

import (
	"github.com/eoscanada/derr"
	"github.com/eoscanada/logging"
	"go.uber.org/zap"
)

var zlog = zap.NewNop()

func setupLogger() {
	zlog = logging.MustCreateLogger()

	// Setting up all package loggers (dependencies)
	derr.SetLogger(zlog)
}
