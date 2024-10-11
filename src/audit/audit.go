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
func SaveAuditLog(action string) func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		hlog.Info("action: ", action)
		// 进行权限校验
		//token := c.Request.Header.Get("X-Auth-Token")
		//flag, result := CheckAction(action, token)
		//if !flag {
		//	hlog.Warn("No operation permission")
		//	c.JSON(http.StatusForbidden, result)
		//	return
		//}

		raw, _ := c.Body()
		var data ReqCreateAudLogRaw
		if err := json.Unmarshal(raw, &data); err != nil {
			hlog.Warn("json Unmarshal err", err)
			c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Insert data failed", ""))
			return
		}
		uuid := uuid4.Uuid4Str()
		t := time.Now().UnixNano() / 1e6
		d := models.DbAuditLog{
			Eid:        &uuid,
			UserId:     &data.UserID,
			Account:    &data.Account,
			SourceIp:   &data.SourceIP,
			Service:    &data.Service,
			Name:       &data.Name,
			Rating:     &data.Rating,
			ETime:      &data.Etime,
			Message:    &data.Message,
			CreateTime: &t,
		}
		if err := models.InstAuditLog(&d); err != nil {
			c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Insert data failed", ""))
			return
		}
		c.JSON(http.StatusOK, answer.ResBody(common.EcodeOK, "", ""))
	}
}

// TracesAuditLog 查询审计列表
func TracesAuditLog(action string) func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		hlog.Info("action: ", action)
		// 进行权限校验
		// token := c.Request.Header.Get("X-Auth-Token")
		// flag, result := CheckAction(action, token)
		// if !flag {
		// 	c.JSON(http.StatusForbidden, result)
		// 	return
		// }

		to, err := strToInt64(c.DefaultQuery("to", ""))
		if err != nil {
			c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Query parameter error", ""))
			return
		}
		from, err := strToInt64(c.DefaultQuery("from", ""))
		if err != nil {
			c.JSON(http.StatusBadRequest, answer.ResBody(common.EcodeError, "Query parameter error", ""))
			return
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
		var count int
		row, err := models.SelectAuditLog(from, to, page, pageSize, &count)
		if err != nil {
			hlog.Error("Database query failure, err: ", err)
			c.JSON(http.StatusInternalServerError, answer.ResBody(common.EcodeError, "Internal service error", ""))
			return
		}

		alogs := make([]ResTracesAuditLog, 0)
		for _, item := range *row {
			a := ResTracesAuditLog{
				Eid:        item.UserId,
				UserID:     item.UserId,
				Account:    item.Account,
				SourceIP:   item.SourceIp,
				Service:    item.Service,
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
