package mysql

import (
	"ats/src/config"
	"fmt"
	"log"

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
	log.Println("Connect database:", cfg.Host+":"+cfg.Port)

	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	// 配置连接池最大连接数
	DB.SetMaxOpenConns(300)
	// 配置连接池最大空闲连接数
	DB.SetMaxIdleConns(20)
	return nil
}
