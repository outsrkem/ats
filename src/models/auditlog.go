package models

import (
	"ats/src/database/mysql"
	"ats/src/slog"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
)

// SelectAuditLog 查询日志列表
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
	if q.Service != "" {
		query.Where("service = ?", q.Service)
	}
	if q.ResourceId != "" {
		query.Where("resource_id = ?", q.ResourceId)
	}
	if q.EventName != "" {
		query.Where("name = ?", q.EventName)
	}
	err := query.Count(count).Limit(q.PageSize).Offset((q.Page - 1) * q.PageSize).Find(&alog).Error
	return alog, err
}

// InstAuditLog 保存日志
func InstAuditLog(c *app.RequestContext, supeve []*OrmSupEve, extras []*OrmExtras, alog []*OrmAuditLog) error {
	klog := slog.FromContext(c)
	if supeve == nil || extras == nil {
		return errors.New("supeve and extras cannot be nil")
	}

	// 通用的错误处理函数
	createRecord := func(tx *gorm.DB, record interface{}, name string) error {
		if err := tx.Create(record).Error; err != nil {
			klog.Errorf("Failed to create %s: %v", name, err)
			return err
		}
		return nil
	}

	return mysql.DB.Transaction(func(tx *gorm.DB) error { // 使用自动事务
		// 创建 supeve 记录
		if err := createRecord(tx, supeve, "supeve"); err != nil {
			return err
		}

		// 创建 extras 记录
		if err := createRecord(tx, extras, "extras"); err != nil {
			return err
		}

		// 创建 alog 记录
		if err := createRecord(tx, alog, "alog"); err != nil {
			return err
		}

		return nil
	})
}

// FindAlogExtras 查询日志扩展数据
func FindAlogExtras(exid string) (*OrmExtras, error) {
	var extras OrmExtras
	query := mysql.DB.Where("exid = ?", exid).First(&extras)
	return &extras, query.Error
}

func DeleteAuditLog(t int64) (int64, error) {
	var rowsAffected int64
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		var seids []string // 存储筛选出的 seid
		err := tx.Model(&OrmSupEve{}).Select("seid").Where("etime < ?", t).Find(&seids).Error
		if err != nil {
			return err
		}

		// 删除 OrmSupEve 表中符合条件的记录
		result := tx.Where("seid IN (?)", seids).Delete(&OrmSupEve{})
		if result.Error != nil {
			return result.Error
		}
		rowsAffected = result.RowsAffected
		if len(seids) > 0 {
			// 删除 OrmAuditLog 表中符合条件的记录
			if err := tx.Where("seid IN (?)", seids).Delete(&OrmAuditLog{}).Error; err != nil {
				return err
			}
			// 删除 OrmExtras 表中符合条件的记录
			if err := tx.Where("seid IN (?)", seids).Delete(&OrmExtras{}).Error; err != nil {
				return err
			}
		}
		return nil
	})

	return rowsAffected, err
}
