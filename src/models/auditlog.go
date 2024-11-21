package models

import (
	"ats/src/database/mysql"
)

func SelectAuditLog(from, to int64, page, pageSize int, count *int64) ([]OrmAuditLog, error) {
	var alog []OrmAuditLog
	err := mysql.DB.Model(&OrmAuditLog{}).
		Where("etime>=? AND etime<=?", from, to).
		Count(count).Limit(pageSize).Offset((page - 1) * pageSize).Order("id").Find(&alog).Error
	return alog, err
}

func InstAuditLog(d []OrmAuditLog) error {
	return mysql.DB.Create(&d).Error
}
