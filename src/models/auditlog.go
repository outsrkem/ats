package models

import (
	"ats/src/database/mysql"

	"github.com/cloudwego/hertz/pkg/common/hlog"
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

func InstAuditLog(d []OrmAuditLog) error {
	return mysql.DB.Create(&d).Error
}
