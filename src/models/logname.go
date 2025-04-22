package models

import (
	"ats/src/database/mysql"
)

type OrmLogName struct {
	ID         uint32 `gorm:"column:id;primaryKey;"`
	Name       string `gorm:"column:name;"`
	Zhcn       string `gorm:"column:zhcn;"`
	Enus       string `gorm:"column:enus;"`
	CreateTime int64  `gorm:"column:create_time;"`
	UpdateTime int64  `gorm:"column:update_time;"`
}

func (*OrmLogName) TableName() string {
	return "logname"
}

func FineAllLogName() ([]*OrmLogName, error) {
	var logName []*OrmLogName
	if err := mysql.DB.Find(&logName).Error; err != nil {
		return nil, err
	}
	return logName, nil
}
