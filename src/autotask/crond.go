package autotask

import (
	"ats/src/config"
	"ats/src/slog"
	"time"

	"github.com/robfig/cron/v3"
)

func StartClean() {
	klog := slog.FromContext(nil)
	klog.Info("Scheduled task starting...")
	crontab := config.Cfg.Ats.Cron.Cleanlog.Time
	if crontab == "" {
		crontab = "10 3 */3 * *"
		klog.Warn("Cron time is empty, use default.")
	}

	klog.Infof("Cron time is [%s]", crontab)

	c := cron.New()
	_, err := c.AddFunc(crontab, func() {
		klog.Info("auto clean task. time now date: ", time.Now().Format(time.DateTime))
		job1()
	})

	if err != nil {
		klog.Error("auto task error: ", err)
		return
	}

	klog.Info("auto clean task start success.")
	c.Start()
}
