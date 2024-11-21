package route

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func apc(action string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		hlog.Debug(action)
		c.Next(ctx)
	}
}
