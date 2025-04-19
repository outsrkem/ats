package audit

import (
	"ats/src/models"
	"ats/src/pkg/answer"
	"ats/src/pkg/common"
	"ats/src/pkg/uuid4"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

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
		raw, _ := c.Body()
		var data ReqCreateAudLogRaw
		if err := json.Unmarshal(raw, &data); err != nil {
			hlog.Warn("json Unmarshal err", err)
			c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Invalid request data.", ""))
			return
		}
		if len(data.Events) > 10 {
			hlog.Warn("More than 10 events uploaded.")
			c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "The number of events exceeded the upper limit by 10.", ""))
			return
		}
		hlog.Info("Create event service: ", data.Service)
		hlog.Infof("The number of created events is %d", len(data.Events))

		supeve, eventAlog, extras := processEvents(data)
		hlog.Debugf("supeve: %+v", supeve)
		hlog.Debugf("extras: %+v", extras)
		hlog.Debugf("eventAlog: %+v", eventAlog)

		if err := models.InstAuditLog(supeve, extras, eventAlog); err != nil {
			hlog.Error("Event creation failure, ", err)
			c.JSON(http.StatusInternalServerError, answer.ResBody(common.EcodeError, "Insert data to db failed", ""))
			return
		}

		hlog.Info("Event creation success.")
		c.JSON(http.StatusOK, answer.ResBody(common.EcodeOK, "", ""))
	}
}

// TracesAuditLog 查询审计日志列表
// Parameter from int64: 事件起始时间,包含该时间
// Parameter to int64: 事件截至事件,包含该时间，有to时必须有from，否则to无效
// Parameter page int:
// Parameter page_size int:
func TracesAuditLog() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		var (
			err error
			q   models.QueryCon
		)
		if c.DefaultQuery("from", "") != "" {
			// from 是使用 to 前提条件，否则to参数无效
			q.From, err = strToInt64(c.DefaultQuery("from", ""))
			if err != nil {
				c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Query parameter error", ""))
				return
			}
			if c.DefaultQuery("to", "") != "" {
				q.To, err = strToInt64(c.DefaultQuery("to", ""))
				if err != nil {
					c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Query parameter error", ""))
					return
				}
			}
			hlog.Infof("Query Audit Log, from:%v, to:%v", q.From, q.To)
		}

		if q.From == 0 {
			// 获取1月之前的时间戳，即默查询2月之内的事件
			q.From = time.Now().AddDate(0, -2, 0).UnixNano() / 1e6
		}
		q.Page, err = strToInt(c.DefaultQuery("page", "1"))
		if err != nil {
			c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Query parameter error", ""))
			return
		}
		q.PageSize, err = strToInt(c.DefaultQuery("page_size", "10"))
		if err != nil {
			c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Query parameter error", ""))
			return
		}

		var count int64
		row, err := models.SelectAuditLog(q, &count) // 查询日志
		if err != nil {
			hlog.Error("Database query failure, err: ", err)
			c.JSON(http.StatusInternalServerError, answer.ResBody(common.EcodeError, "Internal service error", ""))
			return
		}

		alogs := make([]ResTracesAuditLog, 0)
		for _, item := range row {
			a := ResTracesAuditLog{
				Eid:        item.UserId,
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
			alogs = append(alogs, a)
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
		// 查询待修改策略
		exid := c.Param("exid")
		if ok := common.CheckUuId(exid); !ok {
			hlog.Error("Invalid policy id format ", exid)
			c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Invalid extras id format.", ""))
			return
		}

		result, err := models.FindAlogExtras(exid)
		if err != nil {
			hlog.Error("Database query failure, err: ", err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, answer.ResBody(common.EcodeError, err.Error(), ""))
			} else {
				c.JSON(http.StatusInternalServerError, answer.ResBody(common.EcodeError, err.Error(), ""))
			}
			return
		}

		hlog.Debug(result)
		var _reqdata interface{}
		if err := json.Unmarshal([]byte(result.Reqdata), &_reqdata); err != nil {
			hlog.Error("json.Unmarshal ", err)
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
		payload := map[string]interface{}{
			"extras": extras,
		}
		c.JSON(http.StatusOK, answer.ResBody(common.EcodeOK, "", payload))
	}
}
