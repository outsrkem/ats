package autotask

import (
	"ats/src/config"
	"ats/src/models"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func job1() {
	hlog.Info("Start the cleaning task.")
	daysToDelete := config.Cfg.Ats.Cron.Cleanlog.Days
	if daysToDelete <= 0 {
		daysToDelete = 7
		hlog.Warnf("days is empty, use default %d.", daysToDelete)
	}
	now := time.Now()
	todayZero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	deleteBefore := todayZero.AddDate(0, 0, -daysToDelete).UnixMilli()
	hlog.Infof("Cleaning time point: %d", deleteBefore)
	rowsAffected, err := models.DeleteAuditLog(deleteBefore)
	if err != nil {
		hlog.Errorf("Error when deleting auditlog: %v", err)
		return
	}
	hlog.Infof("Delete the record of condition %d this time ", rowsAffected)
}
