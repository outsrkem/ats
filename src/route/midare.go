package route

import (
	"ats/src/config"
	"ats/src/pkg/answer"
	"ats/src/pkg/uuid4"
	"ats/src/slog"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

const xRequestIdKey = "X-Request-Id"
const xAuthTokenKey = "X-Auth-Token"

func RequestId() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)
		xRequestId := string(c.GetHeader(xRequestIdKey))
		if xRequestId == "" {
			xRequestId = uuid4.Uuid4Str()
			c.Response.Header.Set(xRequestIdKey, xRequestId)
			klog.Warnf("request id is empty, Set a new request id: %s", xRequestId)
		}
		c.Set("xRequestId", xRequestId)
		c.Next(ctx)
	}
}

func RequestRecorder() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)
		start := time.Now()
		c.Next(ctx)
		stop := time.Now()
		latency := stop.Sub(start)
		klog.Infof("|%14s | %d |%7s %s",
			latency, c.Response.StatusCode(), string(c.Request.Method()), c.Request.URI().String())
	}
}

func apc(action string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)
		xRequestId := c.Request.Header.Get(xRequestIdKey)
		xAuthToken := c.Request.Header.Get(xAuthTokenKey)
		klog.Debug("token: ", xAuthToken)
		klog.Infof("start check action, [%s] [%s]", xRequestId, action)
		if xAuthToken == "" {
			klog.Error(xAuthTokenKey + " is empty.")
			c.JSON(http.StatusForbidden, answer.ResBody(answer.EcodeInvalidTokenError, xAuthTokenKey+" is empty.", ""))
			c.Abort()
			return
		}

		type actionRaw struct {
			Uias struct {
				Action string `json:"action"`
			} `json:"uias"`
		}
		var raw actionRaw
		raw.Uias.Action = action
		rawJson, err := json.Marshal(raw)
		if err != nil {
			klog.Errorf("Error marshaling audit log: %v", err)
			c.Abort()
			return
		}

		url := config.Cfg.Ats.Uias.Endpoint + "/v1/uias/action/check"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(rawJson))
		if err != nil {
			klog.Errorf("Error creating request: %v", err)
			c.JSON(http.StatusForbidden, answer.ResBody(answer.EcodeInvalidTokenError, "Internal service error.", ""))
			c.Abort()
			return
		}
		req.Header.Set(xAuthTokenKey, xAuthToken)
		req.Header.Set(xRequestIdKey, xRequestId)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			klog.Errorf("Error sending req log: %v", err)
			c.JSON(http.StatusForbidden, answer.ResBody(answer.EcodeInvalidTokenError, "Internal service error.", ""))
			c.Abort()
			return
		}
		defer func() {
			if resp != nil {
				if err := resp.Body.Close(); err != nil {
					klog.Errorf("Close request failed: %v", err)
				}
			}
		}()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			klog.Error("io.ReadAll", err)
			c.JSON(http.StatusForbidden, answer.ResBody(answer.EcodeInvalidTokenError, "Internal service error.", ""))
			c.Abort()
			return
		}

		var result AuthResultData
		if resp.StatusCode != http.StatusOK {
			klog.Errorf("Request failed with status code %d", resp.StatusCode)
			c.JSON(resp.StatusCode, result)
			c.Abort()
			return
		}

		if err := json.Unmarshal(body, &result); err != nil {
			klog.Warn("json Unmarshal err: ", err)
			klog.Error(result)
			c.JSON(http.StatusForbidden, answer.ResBody(answer.EcodeInvalidTokenError, "Internal service error.", ""))
			c.Abort()
			return
		}

		klog.Debugf("result: %+v", result)
		authentication := result.Payload.Authentication
		if authentication != 1 {
			// 没有权限，返回403和上游返回体，便于查看问题
			klog.Warnf("Permission denial. result: %+v", result)
			c.JSON(403, result)
			c.Abort()
			return
		}

		klog.Info("Permission is granted, and the operation is authorized.")
		c.Set("userId", result.Payload.User.ID)
		c.Set("account", result.Payload.User.Name.Account)
		klog.Debug("end check action")
		c.Next(ctx)
	}
}
