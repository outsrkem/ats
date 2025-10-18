package models

import (
	"ats/src/database/mysql"
	"ats/src/slog"
	"errors"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
)

// SelectAuditLog 查询日志列表
func SelectAuditLog(domainId string, q QueryCon, count *int64) ([]OrmAuditLog, error) {
	var alog []OrmAuditLog
	query := mysql.DB.Model(&OrmAuditLog{}).Where("domain_id = ?", domainId)
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
	err := query.Order("id DESC").Count(count).Limit(q.Limit).Offset(q.Offset).Find(&alog).Error
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
	var totalRowsAffected int64

	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 查询符合条件的seid
		var seids []string
		if err := tx.Model(&OrmSupEve{}).Select("seid").Where("etime < ?", t).Find(&seids).Error; err != nil {
			return fmt.Errorf("query seids failed: %w", err)
		}

		// 2. 若seids为空，直接返回不执行删除
		if len(seids) == 0 {
			return nil
		}

		// 3. 删除SupEve表记录
		supEveResult := tx.Where("seid IN (?)", seids).Delete(&OrmSupEve{})
		if supEveResult.Error != nil {
			return fmt.Errorf("delete OrmSupEve failed: %w", supEveResult.Error)
		}

		totalRowsAffected = supEveResult.RowsAffected

		// 4. 删除AuditLog表记录
		if err := tx.Where("seid IN (?)", seids).Delete(&OrmAuditLog{}).Error; err != nil {
			return fmt.Errorf("delete OrmAuditLog failed: %w", err)
		}

		// 5. 删除Extras表记录
		if err := tx.Where("seid IN (?)", seids).Delete(&OrmExtras{}).Error; err != nil {
			return fmt.Errorf("delete OrmExtras failed: %w", err)
		}

		return nil
	})

	return totalRowsAffected, err
}
