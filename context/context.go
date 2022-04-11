package context

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micrease/micrease-core/errs"
	"gorm.io/gorm"
)

type Context struct {
	errs.Error //
	Orm        *gorm.DB
	Ctx        context.Context
	GinCtx     *gin.Context
}
