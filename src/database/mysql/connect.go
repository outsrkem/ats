package mysql

import (
	"ats/src/config"
	"database/sql"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

var DB *gorm.DB

// InitDB initializes the database connection.
func InitDB(cfg *config.Database) {
	// Construct the DSN string.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", cfg.User, cfg.Passwd, cfg.Host, cfg.Port, cfg.Name)
	hlog.Debug("database passwd: ", cfg.Passwd)
	hlog.Info("Connect database: ", cfg.Host+":"+cfg.Port)
	retries := 0
	backoff := time.Second
	var (
		err error
		db  *gorm.DB
		_db *sql.DB
	)
	for {
		// Attempt to connect to the database.
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   "",   // 表名前缀，`User`表为`t_users`
				SingularTable: true, // 使用单数表名，启用该选项后，`User` 表将是`user`
			},
		})
		if err == nil {
			hlog.Info("database connection is successful.")
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
	_db, _ = db.DB()
	_db.SetMaxOpenConns(50)               // 设置打开数据库链接的最大数量
	_db.SetMaxIdleConns(5)                // 设置空闲连接池中链接的最大数量
	_db.SetConnMaxLifetime(time.Hour * 2) // 最大连接有效时间，防止异常连接一直占用
	_db.SetConnMaxIdleTime(time.Hour)     // 空闲连接如果持续xx时间没有被使用，就会被关闭

	DB = db
}
