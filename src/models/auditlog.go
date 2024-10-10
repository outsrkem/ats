package models

import (
	"ats/src/database/mysql"
	"fmt"
)

type AuditLog struct {
	Id        int64  `db:"id"         json:"id"`
	Uuid      int64  `db:"uuid"       json:"uuid"`
	UserId    string `db:"user_id"    json:"user_id"`
	EventTime string `db:"event_time" json:"event_time"`
	SourceIp  string `db:"source_ip"  json:"source_ip"`
}

func SelectAuditLog() (*[]AuditLog, error) {
	sqlStr := "SELECT * FROM auditlog"
	var alog []AuditLog
	if err := mysql.DB.Select(&alog, sqlStr); err != nil {
		fmt.Printf("get data failed, err:%v\n", err)
		return nil, err
	}
	return &alog, nil
}
