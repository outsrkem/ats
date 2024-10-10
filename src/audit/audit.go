package audit

import (
	"ats/src/models"
	"ats/src/pkg/common"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
)

func CheckAction(action string) bool {
	return os.Getenv("PERMISSION_CHECK_ON") == "ON"
}

// 保存审计日志
func SaveAuditLog(action string) func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		// 进行权限校验
		fmt.Println(CheckAction(action))
		c.JSON(http.StatusOK, common.ResBody(common.EcodeOK, "", ""))
	}
}

// 查询事件列表
func TracesAuditLog(action string) func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		fmt.Println(CheckAction(action))
		from := c.DefaultQuery("from", "")
		to := c.DefaultQuery("to", "")
		fmt.Println(from, to)
		row, _ := models.SelectAuditLog()
		payload := map[string]*[]models.AuditLog{
			"items": row,
		}
		c.JSON(http.StatusOK, common.ResBody(common.EcodeOK, "", payload))
	}
}
