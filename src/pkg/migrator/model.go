package migrator

import (
	"gorm.io/gorm"
)

// ISO8601TZ defines the time format for ISO 8601 with timezone offset
// Format: "2006-01-02T15:04:05.000-0700"
const ISO8601TZ = "2006-01-02T15:04:05.000-0700"

// MigrationRecord represents a migration history record
type MigrationRecord struct {
	ID           uint   `gorm:"column:id;primaryKey"`
	Version      string `gorm:"column:version;size:50;not null"`
	ScriptName   string `gorm:"column:script_name;size:255;not null"`
	Direction    string `gorm:"column:direction;size:10;not null"`     // up/down
	ExecutedTime string `gorm:"column:executed_time;size:30;not null"` // ISO 8601
	Status       string `gorm:"column:status;size:20;not null"`        // success/fail
	ErrorMsg     string `gorm:"column:error_msg;type:text"`
	DeletedTime  string `gorm:"column:deleted_time;size:30;index"`
}

// InitTable initializes the migration history table
func InitTable(db *gorm.DB, tableName string) error {
	return db.Table(tableName).AutoMigrate(&MigrationRecord{})
}
