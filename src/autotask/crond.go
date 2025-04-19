package autotask

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/robfig/cron/v3"
	"time"
)

const crontab = "0 * * * *"
const daysToDelete = 2 // 删除N天前的记录

func StartClean() {
	hlog.Info("Scheduled task start")
	c := cron.New()
	_, err := c.AddFunc(crontab, func() {
		hlog.Info("auto clean task. time now date: ", time.Now().Format("2006-01-02 15:04:05"))
		job1()
	})
	if err != nil {
		hlog.Error("auto task error: ", err)
	}
	c.Start()
}
