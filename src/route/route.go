package route

import (
	"context"
	"net/http"

	"ats/src/audit"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func RouteAts(h *server.Hertz) {
	h.GET("/", func(ctx context.Context, c *app.RequestContext) {
		c.String(http.StatusOK, "get")
	})

	// 上传审计日志
	h.POST("/v1/ats/upload", audit.SaveAuditLog("ats:traces:create"))

	// 查询事件列表
	h.GET("/v1/ats/traces", audit.TracesAuditLog("ats:traces:list"))

}
