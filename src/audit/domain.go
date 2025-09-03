package audit

import (
	"ats/src/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

// IfNotExists 检查域不存在
func IfNotExists(domainId string) bool {
	_, err := models.FindDomainById(domainId)
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func CreateDomain(domainId string) error {
	if IfNotExists(domainId) {
		// 不存在则创建域
		domain := &models.OrmDomain{
			DomainId:   domainId,
			CreateTime: time.Now().UnixMilli(),
		}
		err := models.InsertDomain(domain)
		if err != nil {
			return err
		}
	}
	return nil
}
