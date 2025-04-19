package autotask

import (
	"ats/src/models"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"time"
)

func job1() {
	hlog.Info("Start the cleaning task.")
	now := time.Now()
	todayZero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	deleteBefore := todayZero.AddDate(0, 0, -daysToDelete).UnixMilli()
	fmt.Println(deleteBefore)
	hlog.Infof("Cleaning time point: %d", deleteBefore)
	rowsAffected, err := models.DeleteAuditLog(deleteBefore)
	if err != nil {
		hlog.Errorf("Error when deleting auditlog: %v", err)
	}
	hlog.Infof("Delete the record of condition %d this time ", rowsAffected)
}
