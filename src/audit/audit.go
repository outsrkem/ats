package audit

import (
	"ats/src/models"
	"ats/src/pkg/answer"
	"ats/src/pkg/common"
	"ats/src/pkg/uuid4"
	"ats/src/slog"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
)

// 检查事件时间的有效性
// 在过去的1小时整之内, 即当前时间减去1小时(包含该时刻)至当前时间之间的事件为有效事件
// e.g. 当前时间18:30:00, 则在17:30:00~18:30:00之间的事件为有效
func checkEtime(etime int64) bool {
	now := time.Now().UnixMilli()
	// 在过去的1小时之内
	return etime > now-3600000 && etime <= now
}

// processEvents 处理事件数据
func processEvents(data ReqCreateAudLogRaw) ([]*models.OrmSupEve, []*models.OrmAuditLog, []*models.OrmExtras) {
	supeve := make([]*models.OrmSupEve, 0)
	eventAlog := make([]*models.OrmAuditLog, 0)
	extras := make([]*models.OrmExtras, 0)
	t := time.Now().UnixMilli()
	for _, event := range data.Events {
		seid := uuid4.Uuid4Str()
		supeve = append(supeve, &models.OrmSupEve{
			Seid:       seid,
			Etime:      event.Etime,
			CreateTime: t,
		})
		exid := uuid4.Uuid4Str()
		extras = append(extras, &models.OrmExtras{
			Seid:     seid,
			Exid:     exid,
			Reqdata:  event.Reqdata,
			Uagent:   event.Uagent,
			SourceIp: event.SourceIP,
			Method:   event.Method,
			ReqUrl:   event.ReqUrl,
		})
		baseEvent := models.OrmAuditLog{
			Seid:    seid,
			UserId:  event.UserID,
			Account: event.Account,
			Service: data.Service,
			Name:    event.Name,
			Rating:  event.Rating,
			Message: event.Message,
			Extras:  exid,
			ETime:   event.Etime,
		}
		if len(event.ResourceId) > 0 {
			for _, v := range event.ResourceId {
				_event := baseEvent
				_event.Eid = uuid4.Uuid4Str()
				_event.ResourceId = v
				eventAlog = append(eventAlog, &_event)
			}
		} else {
			_event := baseEvent
			_event.Eid = uuid4.Uuid4Str()
			eventAlog = append(eventAlog, &_event)
		}
	}
	return supeve, eventAlog, extras
}

// SaveAuditLog 保存审计日志
func SaveAuditLog() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)
		raw, _ := c.Body()
		var data ReqCreateAudLogRaw
		if err := json.Unmarshal(raw, &data); err != nil {
			klog.Warn("json Unmarshal err", err)
			c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Invalid request data.", ""))
			return
		}
		if len(data.Events) > 10 {
			klog.Warn("More than 10 events uploaded.")
			c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "The number of events exceeded the upper limit by 10.", ""))
			return
		}
		klog.Info("Create event service: ", data.Service)
		klog.Infof("The number of created events is %d", len(data.Events))

		for _, event := range data.Events {
			// 事件时间不符合要求(在过去的1小时之内)
			if ok := checkEtime(event.Etime); !ok {
				msg := "The event time does not meet the requirements"
				klog.Warnf("%s [%d]", msg, event.Etime)
				c.JSON(http.StatusUnprocessableEntity, answer.ResBody(common.EcodeError, msg, ""))
				return
			}
		}

		supeve, eventAlog, extras := processEvents(data)
		klog.Debugf("supeve: %+v", supeve)
		klog.Debugf("extras: %+v", extras)
		klog.Debugf("eventAlog: %+v", eventAlog)

		if err := models.InstAuditLog(c, supeve, extras, eventAlog); err != nil {
			klog.Error("Event creation failure, ", err)
			c.JSON(http.StatusInternalServerError, answer.ResBody(common.EcodeError, "Insert data to db failed", ""))
			return
		}

		klog.Info("Event creation success.")
		c.JSON(http.StatusOK, answer.ResBody(common.EcodeOK, "", ""))
	}
}

// TracesAuditLog 查询审计日志列表
// Parameter from int64: 事件起始时间,包含该时间
// Parameter to int64: 事件截至事件,包含该时间，有to时必须有from，否则to无效
// Parameter page int:
// Parameter page_size int:
// Parameter svc string: 按服务查询
// Parameter resid string: 按资源ID查询
func TracesAuditLog() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)
		var (
			err error
			q   models.QueryCon
		)
		if c.DefaultQuery("from", "") != "" {
			// from 是使用 to 前提条件，否则to参数无效
			q.From, err = strToInt64(c.DefaultQuery("from", ""))
			if err != nil {
				klog.Error("strToInt64 error", err)
				c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Query parameter error", ""))
				return
			}
			if c.DefaultQuery("to", "") != "" {
				q.To, err = strToInt64(c.DefaultQuery("to", ""))
				if err != nil {
					klog.Error("strToInt64 error", err)
					c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Query parameter error", ""))
					return
				}
			}
			klog.Infof("Query Audit Log, from:%v, to:%v", q.From, q.To)
		}

		if q.From == 0 {
			// 默查询N天之内的事件
			day := 15
			q.From = time.Now().AddDate(0, 0, -day).UnixNano() / 1e6
			klog.Infof("No specific time is specified for querying event logs within %d days.", day)
		}
		q.Page, err = strToInt(c.DefaultQuery("page", "1"))
		if err != nil {
			klog.Error("strToInt error", err)
			c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Query parameter error", ""))
			return
		}
		q.PageSize, err = strToInt(c.DefaultQuery("page_size", "10"))
		if err != nil {
			klog.Error("strToInt error", err)
			c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Query parameter error", ""))
			return
		}

		// 按服务查询
		q.Service = c.DefaultQuery("svc", "")
		if q.Service != "" && !isAlphaASCIILoop(q.Service) {
			klog.Errorf("Query parameter error: svc [%s]", q.Service)
			c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Query parameter error: svc", ""))
			return
		}

		// 按资源id查询
		q.ResourceId = c.DefaultQuery("resid", "")
		if q.ResourceId != "" && !common.CheckUuId(q.ResourceId) {
			klog.Errorf("Query parameter error: resid [%s]", q.ResourceId)
			c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Query parameter error: resid", ""))
			return
		}

		// 按事件名称查询
		q.EventName = c.DefaultQuery("name", "")
		if q.EventName != "" && !isAlphaASCIILoop(q.EventName) {
			klog.Errorf("Query parameter error: name [%s]", q.EventName)
			c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Query parameter error: name", ""))
			return
		}

		var count int64
		result, err := models.SelectAuditLog(q, &count) // 查询日志
		if err != nil {
			klog.Error("Database query failure, err: ", err)
			c.JSON(http.StatusInternalServerError, answer.ResBody(common.EcodeError, "Internal service error", ""))
			return
		}

		alogs := make([]ResTracesAuditLog, len(result))
		for i, item := range result {
			alogs[i] = ResTracesAuditLog{
				Eid:        item.Eid,
				UserID:     item.UserId,
				Account:    item.Account,
				Service:    item.Service,
				ResourceId: item.ResourceId,
				Name:       GetElogName(item.Name, "zhcn"),
				Rating:     item.Rating,
				ETime:      item.ETime,
				Message:    item.Message,
				Extras:     item.Extras,
			}
		}

		pageInfo := answer.SetPageInfo(q.PageSize, q.Page, count)
		payload := map[string]interface{}{
			"items":     alogs,
			"page_info": pageInfo,
		}
		c.JSON(http.StatusOK, answer.ResBody(common.EcodeOK, "", payload))
	}
}

// TracesExtras 查询日志扩展数据
func TracesExtras() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)
		// 查询待修改策略
		exid := c.Param("exid")
		if ok := common.CheckUuId(exid); !ok {
			klog.Error("Invalid policy id format ", exid)
			c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Invalid extras id format.", ""))
			return
		}

		result, err := models.FindAlogExtras(exid)
		if err != nil {
			klog.Error("Database query failure, err: ", err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, answer.ResBody(common.EcodeError, err.Error(), ""))
			} else {
				c.JSON(http.StatusInternalServerError, answer.ResBody(common.EcodeError, err.Error(), ""))
			}
			return
		}

		klog.Debug(result)
		var _reqdata interface{}
		if err := json.Unmarshal([]byte(result.Reqdata), &_reqdata); err != nil {
			klog.Error("json.Unmarshal ", err)
			_reqdata = result.Reqdata
		}

		// 日志扩展信息
		type ReqData struct {
			Reqdata  interface{} `json:"reqdata"`
			Uagent   string      `json:"uagent"`
			SourceIp string      `json:"source_ip"`
			Method   string      `json:"method"`
			ReqUrl   string      `json:"requrl"`
		}
		extras := ReqData{
			Reqdata:  &_reqdata,
			Uagent:   result.Uagent,
			SourceIp: result.SourceIp,
			Method:   result.Method,
			ReqUrl:   result.ReqUrl,
		}
		payload := map[string]ReqData{
			"extras": extras,
		}
		c.JSON(http.StatusOK, answer.ResBody(common.EcodeOK, "", payload))
	}
}

// isAlphaASCIILoop 检查服务名称是否都是字母, 即 a-z 和 A-Z
// 字符串全是a-z和A-Z 返回true
func isAlphaASCIILoop(s string) bool {
	for _, c := range s {
		if !(('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')) {
			return false
		}
	}
	return true
}
