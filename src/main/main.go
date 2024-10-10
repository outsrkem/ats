package main

import (
	"ats/src/config"
	"ats/src/database/mysql"
	"ats/src/route"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/requestid"
)

func main() {
	cfg := config.InitConfig()
	app := cfg.Ats.App
	// 连接数据库MySql
	mysql.InitDB(&cfg.Ats.Database)
	h := server.Default(server.WithHostPorts(app.Bind), server.WithExitWaitTime(0*time.Second))
	h.Use(requestid.New())
	route.RouteAts(h)
	h.Spin()
}
