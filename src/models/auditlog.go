package models

import (
	"ats/src/database/mysql"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
)

func SelectAuditLog(q QueryCon, count *int64) ([]OrmAuditLog, error) {
	var alog []OrmAuditLog
	query := mysql.DB.Model(&OrmAuditLog{}).Order("id DESC")
	if q.From != 0 {
		query.Where("etime>=?", q.From)
		if q.To != 0 {
			// 必须有起始时间才能搭配结束时间查询
			query.Where("etime<=?", q.To)
		}
	}
	err := query.Count(count).Limit(q.PageSize).Offset((q.Page - 1) * q.PageSize).Find(&alog).Error
	hlog.Infof("query condition: %v", q)
	return alog, err
}

func InstAuditLog(extras *OrmExtras, alog []OrmAuditLog) error {
	// 使用自动事务
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(extras).Error; err != nil {
			hlog.Error("Transaction rollback. err: ", err)
			return err
		}
		b := tx.Create(&alog)
		if b.Error != nil {
			hlog.Error("Transaction rollback. err: ", b.Error)
			return b.Error
		}
		return nil
	})
	return err
}
