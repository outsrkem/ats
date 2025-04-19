package autotask

import (
	"ats/src/config"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/robfig/cron/v3"
)

func StartClean() {
	hlog.Info("Scheduled task start")
	crontab := config.Cfg.Ats.Cron.Cleanlog.Time // * * * * *
	if crontab == "" {
		crontab = "10 3 */3 * *"
		hlog.Warnf("Cron time is empty, use default %s", crontab)
	}
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
