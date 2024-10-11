package mysql

import (
	"ats/src/config"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/jmoiron/sqlx"

	// mysql 连接驱动
	_ "github.com/go-sql-driver/mysql"
)

// DB mysql
// var DB *sql.DB
var DB *sqlx.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.Database) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		cfg.User, cfg.Passwd, cfg.Host, cfg.Port, cfg.Name)
	hlog.Info("Connect database: ", cfg.Host+":"+cfg.Port)

	retries := 0
	backoff := time.Second

	for {
		DB, err = sqlx.Connect("mysql", dsn)
		if err == nil {
			break
		}

		retries++
		if retries >= 100 {
			panic(err)
		}

		hlog.Error("Failed to connect to database. Retrying in %v...", backoff)
		time.Sleep(backoff)

		backoff += time.Second
	}

	// 配置连接池最大连接数
	DB.SetMaxOpenConns(300)
	// 配置连接池最大空闲连接数
	DB.SetMaxIdleConns(20)
	return nil
}
