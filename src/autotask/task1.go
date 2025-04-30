package autotask

import (
	"ats/src/config"
	"ats/src/models"
	"ats/src/slog"
	"time"
)

func job1() {
	klog := slog.FromContext(nil)
	klog.Info("Start the cleaning task.")
	daysToDelete := config.Cfg.Ats.Cron.Cleanlog.Days
	if daysToDelete <= 0 {
		daysToDelete = 15
		klog.Warnf("days is empty, use default %d.", daysToDelete)
	}
	now := time.Now()
	todayZero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	deleteBefore := todayZero.AddDate(0, 0, -daysToDelete).UnixMilli()
	klog.Infof("Cleaning time point: %d", deleteBefore)
	rowsAffected, err := models.DeleteAuditLog(deleteBefore)
	if err != nil {
		klog.Errorf("Error when deleting auditlog: %v", err)
		return
	}
	klog.Infof("Delete the record of condition %d this time ", rowsAffected)
}
