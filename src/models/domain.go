package models

import "ats/src/database/mysql"

// FindDomainById 查询域
func FindDomainById(domainId string) (OrmDomain, error) {
	var domain OrmDomain
	r := mysql.DB.Model(&OrmDomain{}).
		Where("domain_id = ?", domainId).First(&domain)
	return domain, r.Error
}

// InsertDomain 创建域
func InsertDomain(domain *OrmDomain) error {
	r := mysql.DB.Model(&OrmDomain{}).Create(domain)
	return r.Error
}
