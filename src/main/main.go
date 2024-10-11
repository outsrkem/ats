package main

import (
	"ats/src/config"
	"ats/src/database/mysql"
	"ats/src/route"
	"ats/src/slog"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/requestid"
)

func init() {
	slog.InitLog()
}

func main() {
	cfg := config.InitConfig()
	app := &cfg.Ats.App
	mysql.InitDB(&cfg.Ats.Database) // 连接数据库MySql

	hlog.Info("start server")
	h := server.Default(server.WithHostPorts(app.Bind), server.WithExitWaitTime(0*time.Second))
	h.Use(requestid.New())
	route.AtsRoute(h)
	h.Spin()
}
