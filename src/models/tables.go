package models

// OrmSupEve 事件表
type OrmSupEve struct {
	ID         uint32 `gorm:"column:id"`
	DomainId   string `gorm:"column:domain_id"`
	Seid       string `gorm:"column:seid"` // 一批事件id
	Etime      int64  `gorm:"column:etime"`
	CreateTime int64  `gorm:"column:create_time"`
}

func (*OrmSupEve) TableName() string {
	return "ats_supeve"
}

// Domain 账号域表
type OrmDomain struct {
	Kid        int64  `gorm:"column:kid"`
	DomainId   string `gorm:"column:domain_id"`
	CreateTime int64  `gorm:"column:create_time"`
}

func (*OrmDomain) TableName() string {
	return "ats_domain"
}

// OrmAuditLog 主日志表
type OrmAuditLog struct {
	Id         int64  `gorm:"column:id"`
	DomainId   string `gorm:"column:domain_id"`
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
	return "ats_auditlog"
}

// OrmExtras 日志扩展数据
type OrmExtras struct {
	Id       int64  `gorm:"column:id"`
	DomainId string `gorm:"column:domain_id"`
	Seid     string `gorm:"column:seid"`
	Exid     string `gorm:"column:exid"`
	Reqdata  string `gorm:"column:reqdata"`
	Uagent   string `gorm:"column:uagent"`
	SourceIp string `gorm:"column:source_ip"`
	Method   string `gorm:"column:method"`
	ReqUrl   string `gorm:"column:requrl"`
}

func (OrmExtras) TableName() string {
	return "ats_extras"
}

// QueryCon 日志查询条件 query 参数
type QueryCon struct {
	Limit      int
	Offset     int
	From       int64  // 起始时间
	To         int64  // 结束时间
	EventName  string // 事件名称
	Service    string // 服务名 svc
	ResourceId string // 资源ID resid
}
