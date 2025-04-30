package mysql

import (
	"ats/src/cfgtypts"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

func NewGormLogger(cfg *cfgtypts.Ats) logger.Interface {
	return logger.New(
		logrus.StandardLogger(),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // 慢查询阈值
			LogLevel:                  logger.Info,            // 日志级别
			IgnoreRecordNotFoundError: true,                   // 忽略记录未找到错误
			Colorful:                  false,                  // 禁用颜色
		},
	)
}
