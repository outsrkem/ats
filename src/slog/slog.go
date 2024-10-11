package slog

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzslog "github.com/hertz-contrib/logger/slog"
)

func InitLog() {
	logger := hertzslog.NewLogger()
	logger.SetLevel(hlog.LevelDebug)
	hlog.SetLogger(logger)
}
