package main

import (
	"ats/src/audit"
	"ats/src/autotask"
	"ats/src/config"
	"ats/src/database/mysql"
	"ats/src/database/sql"
	"ats/src/pkg/migrator"
	"ats/src/route"
	"ats/src/slog"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/sirupsen/logrus"
)

func init() {
	// slog.InitLogger("info")
}

// AutoMigrator performs automatic database migrations using embedded SQL scripts
func AutoMigrator(klog *logrus.Entry) {
	// Initialize migrator with embedded scripts
	m := migrator.NewDefault(mysql.DB,
		&migrator.Config{
			Logger:         migrator.NewKlogLogger(klog),
			TableName:      "ats_migrations",
			ScriptProvider: migrator.NewEmbeddedScriptProvider("script", sql.SQLScript),
		})

	// Execute migration workflow
	if err := migrator.RunWithMigrator(m); err != nil {
		panic("Database migration failed: " + err.Error())
	}
}

func main() {
	cfg := config.InitConfig()
	slog.InitLogger(&cfg.Ats.Log)
	klog := slog.FromContext(nil)
	mysql.InitDB(&cfg.Ats) // 连接数据库MySql

	AutoMigrator(klog)

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
