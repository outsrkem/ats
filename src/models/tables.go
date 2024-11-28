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
	Message    *string `gorm:"column:message"`
	Extras     *string `gorm:"column:extras"`
	ETime      *int64  `gorm:"column:etime"`
	CreateTime *int64  `gorm:"column:create_time"`
}

func (OrmAuditLog) TableName() string {
	return "auditlog"
}

type OrmExtras struct {
	Id      *int64  `gorm:"column:id"`
	Exid    *string `gorm:"column:exid"`
	Reqdata *string `gorm:"column:reqdata"`
	Uagent  *string `gorm:"column:uagent"`
}

func (OrmExtras) TableName() string {
	return "extras"
}

// QueryCon 日志查询条件 query 参数
type QueryCon struct {
	PageSize int
	Page     int
	From     int64 // 起始时间
	To       int64 // 结束时间
}
