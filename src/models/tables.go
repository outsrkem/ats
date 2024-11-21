package models

// OrmAuditLog table: auditlog
type OrmAuditLog struct {
	Id         *int64  `gorm:"column:id"`
	Eid        *string `gorm:"column:eid"`
	UserId     *string `gorm:"column:user_id"`
	Account    *string `gorm:"column:account"`
	SourceIp   *string `gorm:"column:source_ip"`
	Service    *string `gorm:"column:service"`
	ResourceId *string `gorm:"column:resource_id"`
	Name       *string `gorm:"column:name"`
	Rating     *string `gorm:"column:rating"`
	ETime      *int64  `gorm:"column:etime"`
	Message    *string `gorm:"column:message"`
	CreateTime *int64  `gorm:"column:create_time"`
}

func (OrmAuditLog) TableName() string {
	return "auditlog"
}

// QueryCon 日志查询条件 query 参数
type QueryCon struct {
	To       int64
	From     int64
	Page     int
	PageSize int
}
