package models

// OrmAuditLog table: event
type OrmAuditLog struct {
	Id         int64  `gorm:"column:id"`
	Seid       string `gorm:"column:seid"`
	Eid        string `gorm:"column:eid"`
	UserId     string `gorm:"column:user_id"`
	Account    string `gorm:"column:account"`
	Service    string `gorm:"column:service"`
	ResourceId string `gorm:"column:resource_id"`
	Name       string `gorm:"column:name"`
	Rating     string `gorm:"column:rating"`
	Message    string `gorm:"column:message"`
	Extras     string `gorm:"column:extras"`
	ETime      int64  `gorm:"column:etime"`
}

func (OrmAuditLog) TableName() string {
	return "auditlog"
}

type OrmExtras struct {
	Id       int64  `gorm:"column:id"`
	Seid     string `gorm:"column:seid"`
	Exid     string `gorm:"column:exid"`
	Reqdata  string `gorm:"column:reqdata"`
	Uagent   string `gorm:"column:uagent"`
	SourceIp string `gorm:"column:source_ip"`
	Method   string `gorm:"column:method"`
	ReqUrl   string `gorm:"column:requrl"`
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

type OrmSupEve struct {
	ID         uint32 `gorm:"column:id"`
	Seid       string `gorm:"column:seid"` // 一批事件id，同一个请求创建的事件为一批事件
	Etime      int64  `gorm:"column:etime"`
	CreateTime int64  `gorm:"column:create_time"`
}

func (*OrmSupEve) TableName() string {
	return "supeve"
}
