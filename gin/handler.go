package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/micrease/micrease-core/context"
	"github.com/micrease/micrease-core/errs"
)

type GinHandler struct {
	errs.Error
}

func (h GinHandler) ResponseData(ctx *context.Context, data interface{}) {
	ctx.GinCtx.JSON(200, gin.H{"status": errs.StatusSuccess, "message": "操作成功", "data": data})
}

func (h GinHandler) Success(ctx *context.Context) {
	ctx.GinCtx.JSON(200, gin.H{"status": errs.StatusSuccess, "message": "操作成功", "data": ""})
}

func (h GinHandler) Response(ctx *context.Context, data interface{}) {
	h.ResponseData(ctx, data)
}
