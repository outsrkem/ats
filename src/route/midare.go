package route

import (
	"ats/src/config"
	"ats/src/pkg/answer"
	"ats/src/slog"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
	"github.com/hiuias/uias-sdk-go"
)

const xRequestIdKey = "X-Request-Id"
const xAuthTokenKey = "X-Auth-Token"

func RequestId() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		klog := slog.FromContext(c)
		xRequestId := string(c.GetHeader("X-Request-Id"))
		if xRequestId == "" {
			xRequestId = strings.ReplaceAll(uuid.New().String(), "-", "")
			c.Request.Header.Set("X-Request-Id", xRequestId)
			klog.Warnf("request id is empty, Set a new request id: %s", xRequestId)
		}
		c.Set("xRequestId", xRequestId)
		c.Next(ctx)
		// 如果响应头中没有 X-Request-Id，则添加它
		if c.Response.Header.Get("X-Request-Id") == "" {
			c.Response.Header.Set("X-Request-Id", xRequestId)
			klog.Debugf("Set X-Request-Id in response: %s", xRequestId)
		}
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

		client := uias.NewClientBuilder().
			WithEndpoint(config.Cfg.Ats.Uias.Endpoint).
			WithTimeout(10 * time.Second).
			WithSkipTlsVerify(false).
			Build()

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
		rawBody, err := json.Marshal(raw)
		if err != nil {
			klog.Errorf("Error marshaling audit log: %v", err)
			c.JSON(http.StatusForbidden, answer.ResBody(answer.EcodeInvalidTokenError, xAuthTokenKey+" is empty.", ""))
			c.Abort()
			return
		}

		resp, err := client.VerifyAction(ctx, c.Request.Header.Get(xRequestIdKey), c.Request.Header.Get(xAuthTokenKey), rawBody)
		if err != nil {
			klog.Warnf("Permission verification error, upstream exception: %v", err)
			c.JSON(500, resp)
			c.Abort()
			return
		}

		// 状态码非200
		if resp.StatusCode != http.StatusOK {
			klog.Warnf("http status code is not 200: %+v", resp)
			c.JSON(resp.StatusCode, resp.Data)
			c.Abort()
			return
		}

		// 没有权限，返回403和上游返回体
		if resp.Data.Payload.Authentication != 1 {
			klog.Warnf("Permission denial. result: %+v", resp)
			c.JSON(403, resp)
			c.Abort()
			return
		}

		// 获取响应结果
		data := resp.Data
		klog.Info("Permission is granted, and the operation is authorized.")
		c.Set("domainId", data.Payload.User.Domain.Id)
		c.Set("domainName", data.Payload.User.Domain.Name)
		c.Set("userId", data.Payload.User.Id)
		c.Set("account", data.Payload.User.Name.Account)
		klog.Debug("end check action")
		c.Next(ctx)
	}
}
