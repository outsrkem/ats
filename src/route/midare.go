package route

import (
	"ats/src/config"
	"ats/src/pkg/answer"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func apc(action string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		hlog.Debugf("start check action, [%s]", action)
		token := c.Request.Header.Get("X-Auth-Token")
		if token == "" {
			hlog.Error("X-Auth-Token is empty.")
			c.JSON(http.StatusForbidden, answer.ResBody(answer.EcodeInvalidTokenError, "X-Auth-Token is empty.", ""))
			c.Abort()
			return
		}
		hlog.Debug("token: ", token)
		type actionRaw struct {
			Uias struct {
				Action string `json:"action"`
			} `json:"uias"`
		}
		var raw actionRaw
		raw.Uias.Action = action
		rawJson, err := json.Marshal(raw)
		if err != nil {
			hlog.Errorf("Error marshaling audit log: %v", err)
			c.Abort()
			return
		}

		url := config.Cfg.Ats.Uias.Endpoint + "/v1/uias/action/check"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(rawJson))
		if err != nil {
			hlog.Errorf("Error creating request: %v", err)
			c.JSON(http.StatusForbidden, answer.ResBody(answer.EcodeInvalidTokenError, "Internal service error.", ""))
			c.Abort()
			return
		}

		req.Header.Set("X-Auth-Token", token)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			hlog.Errorf("Error sending req log: %v", err)
			c.JSON(http.StatusForbidden, answer.ResBody(answer.EcodeInvalidTokenError, "Internal service error.", ""))
			c.Abort()
			return
		}
		defer func() {
			if resp != nil {
				if err := resp.Body.Close(); err != nil {
					hlog.Error("Close request failed: %v", err)
				}
			}
		}()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			hlog.Error("io.ReadAll", err)
			c.JSON(http.StatusForbidden, answer.ResBody(answer.EcodeInvalidTokenError, "Internal service error.", ""))
			c.Abort()
			return
		}

		var result AuthResultData
		if resp.StatusCode != http.StatusOK {
			hlog.Errorf("Request failed with status code %d", resp.StatusCode)
			c.JSON(resp.StatusCode, result)
			c.Abort()
			return
		}

		if err := json.Unmarshal(body, &result); err != nil {
			hlog.Warn("json Unmarshal err: ", err)
			hlog.Error(result)
			c.JSON(http.StatusForbidden, answer.ResBody(answer.EcodeInvalidTokenError, "Internal service error.", ""))
			c.Abort()
			return
		}

		hlog.Debugf("result: %+v", result)
		authentication := result.Payload.Authentication
		if authentication != 1 {
			// 没有权限，返回403和上游返回体，便于查看问题
			hlog.Warnf("Permission denial. result: %+v", result)
			c.JSON(403, result)
			c.Abort()
			return
		}

		hlog.Info("Permission is granted, and the operation is authorized.")
		c.Set("userId", result.Payload.User.ID)
		c.Set("account", result.Payload.User.Name.Account)
		hlog.Debug("end check action")
		c.Next(ctx)
	}
}
