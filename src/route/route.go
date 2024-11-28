package route

import (
	"ats/src/audit"
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func AtsRoute(h *server.Hertz) {
	h.GET("/", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, utils.H{"message": "Hello World"})
	})

	// 上传审计日志
	h.POST("/v1/ats/upload", apc("ats:traces:create"), audit.SaveAuditLog())

	// 日志查询
	h.GET("/v1/ats/traces", apc("ats:traces:listAlog"), audit.TracesAuditLog())      // 查询事件列表
	h.GET("/v1/ats/extras/:exid", apc("ats:traces:getExtras"), audit.TracesExtras()) // 查询日志扩展数据
}
