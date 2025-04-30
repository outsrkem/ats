package main

import (
	"ats/src/audit"
	"ats/src/autotask"
	"ats/src/config"
	"ats/src/database/mysql"
	"ats/src/route"
	"ats/src/slog"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func init() {
	// slog.InitLogger("info")
}
func main() {
	cfg := config.InitConfig()
	slog.InitLogger(&cfg.Ats.Log)
	klog := slog.FromContext(nil)
	mysql.InitDB(&cfg.Ats) // 连接数据库MySql

	if err := audit.InitLogCache(); err != nil {
		klog.Errorf("Failed to init log type cache: %v", err)
	}

	autotask.StartClean()

	klog.Info("start server")
	svc := server.Default(server.WithHostPorts(cfg.Ats.App.Bind), server.WithExitWaitTime(0*time.Second))
	route.Middleware(svc)
	route.AtsRoute(svc)
	svc.Spin()
}
