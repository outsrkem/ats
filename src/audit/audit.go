package audit

import (
	"ats/src/models"
	"ats/src/pkg/answer"
	"ats/src/pkg/common"
	"ats/src/pkg/uuid4"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// SaveAuditLog 保存审计日志
func SaveAuditLog() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		raw, _ := c.Body()
		var data ReqCreateAudLogRaw
		if err := json.Unmarshal(raw, &data); err != nil {
			hlog.Warn("json Unmarshal err", err)
			c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Insert data failed", ""))
			return
		}
		if len(data.Events) > 100 {
			hlog.Warn("More than 20 events uploaded.")
			c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "The number of events exceeded the upper limit by 100.", ""))
			return
		}
		hlog.Info("Create event service: ", data.Service)
		hlog.Infof("The number of created events is %d", len(data.Events))
		t := time.Now().UnixNano() / 1e6
		d := make([]models.OrmAuditLog, 0)
		for _, item := range data.Events {
			uuid := uuid4.Uuid4Str()
			d = append(d, models.OrmAuditLog{
				Eid:        &uuid,
				UserId:     &item.UserID,
				Account:    &item.Account,
				SourceIp:   &item.SourceIP,
				Service:    &data.Service,
				ResourceId: &item.ResourceId,
				Name:       &item.Name,
				Rating:     &item.Rating,
				ETime:      &item.Etime,
				Message:    &item.Message,
				CreateTime: &t,
			})
		}
		if err := models.InstAuditLog(d); err != nil {
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
// Parameter to int64: 事件截至事件,包含该时间
// Parameter page int:
// Parameter page_size int:
func TracesAuditLog() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		var (
			err  error
			to   int64
			from int64
		)
		if c.DefaultQuery("to", "") != "" && c.DefaultQuery("from", "") != "" {
			// from 和 to 要一起使用，否则无效
			to, err = strToInt64(c.DefaultQuery("to", ""))
			if err != nil {
				c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Query parameter error", ""))
				return
			}
			from, err = strToInt64(c.DefaultQuery("from", ""))
			if err != nil {
				c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Query parameter error", ""))
				return
			}
		}
		if from == 0 || to == 0 {
			now := time.Now()
			to = now.UnixNano() / 1e6
			from = now.Add(-time.Hour).UnixNano() / int64(time.Millisecond)
			hlog.Warn("If the event interval is abnormal, the default interval is used")
		}
		page, err := strToInt(c.DefaultQuery("page", "1"))
		if err != nil {
			c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Query parameter error", ""))
			return
		}
		pageSize, err := strToInt(c.DefaultQuery("page_size", "10"))
		if err != nil {
			c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Query parameter error", ""))
			return
		}

		hlog.Infof("Query Audit Log, from:%v, to:%v", from, to)
		var count int64
		row, err := models.SelectAuditLog(from, to, page, pageSize, &count)
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
				SourceIP:   item.SourceIp,
				Service:    item.Service,
				ResourceId: item.ResourceId,
				Name:       item.Name,
				Rating:     item.Rating,
				ETime:      item.ETime,
				Message:    item.Message,
				CreateTime: item.CreateTime,
			}
			alogs = append(alogs, a)
		}
		pageInfo := answer.SetPageInfo(pageSize, page, count)
		payload := map[string]interface{}{
			"items":     alogs,
			"page_info": pageInfo,
		}
		c.JSON(http.StatusOK, answer.ResBody(common.EcodeOK, "", payload))
	}
}
