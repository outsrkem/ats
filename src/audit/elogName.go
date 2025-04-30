package audit

import (
	"ats/src/models"
	"ats/src/slog"
	"sync"
)

type logName struct {
	zh string
	en string
}

var (
	logCache    = make(map[string]*logName, 1000)
	cacheMutex  sync.RWMutex
	initialized bool
)

// GetElogName 获取日志名称
func GetElogName(logName string, lang string) string {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	if names, exists := logCache[logName]; exists {
		switch lang {
		case "zhcn":
			return names.zh
		case "enus":
			return names.en
		}
	}
	return logName
}

// InitLogCache 初始化加载日志事件名称
func InitLogCache() error {
	klog := slog.FromContext(nil)
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	if initialized {
		return nil
	}

	logTypes, err := models.FineAllLogName()
	if err != nil {
		klog.Error(err.Error())
	}
	for _, lt := range logTypes {
		logCache[lt.Name] = &logName{
			zh: lt.Zhcn,
			en: lt.Enus,
		}
	}
	initialized = true
	return nil
}
