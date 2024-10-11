package models

import (
	"ats/src/database/mysql"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func SelectAuditLog(from, to int64, page, pageSize int, count *int) (*[]DbAuditLog, error) {
	sqlStr := `SELECT * FROM auditlog WHERE etime >= ? AND etime <= ? ORDER BY etime LIMIT ? OFFSET ?;`
	var alog []DbAuditLog
	offset := (page - 1) * pageSize
	if err := mysql.DB.Select(&alog, sqlStr, from, to, pageSize, offset); err != nil {
		hlog.Errorf("get data failed, err:%v", err)
		return nil, err
	}
	countSql := `SELECT COUNT(*) FROM auditlog WHERE etime >= ? AND etime <= ?;`
	err := mysql.DB.Get(count, countSql, from, to)
	if err != nil {
		hlog.Debugf("Count Records sql: %v", countSql)
		hlog.Errorf("Get count err: %v", err)
		return nil, err
	}
	return &alog, nil
}

func InstAuditLog(d *DbAuditLog) error {
	_sql := `INSERT INTO auditlog(eid,user_id,account,source_ip,service,name,rating,etime,message,create_time) VALUES (?, ?,?, ?, ?, ?, ?, ?, ?, ?);`
	_args := []interface{}{d.Eid, d.UserId, d.Account, d.SourceIp, d.Service, d.Name, d.Rating, d.ETime, d.Message, d.CreateTime}
	result, err := mysql.DB.Exec(_sql, _args...)
	if err != nil {
		hlog.Error("install data error: ", err)
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		hlog.Error("install data error: ", err)
	}
	hlog.Info("install data successfully, id: ", id)
	return nil
}
