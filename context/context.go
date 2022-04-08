package context

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/micrease/micrease-core/errs"
)

type Context struct {
	errs.Error //
	Orm        *gorm.DB
	Ctx        context.Context
	GinCtx     *gin.Context
}
