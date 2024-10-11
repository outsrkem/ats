package models

// DbAuditLog table: auditlog
type DbAuditLog struct {
	Id         *int64  `db:"id"`
	Eid        *string `db:"eid"`
	UserId     *string `db:"user_id"`
	Account    *string `db:"account"`
	SourceIp   *string `db:"source_ip"`
	Service    *string `db:"service"`
	Name       *string `db:"name"`
	Rating     *string `db:"rating"`
	ETime      *int64  `db:"etime"`
	Message    *string `db:"message"`
	CreateTime *int64  `db:"create_time"`
}
