package rpc

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/micrease/micrease-core/errs"
	"github.com/micro/go-micro/v2/errors"
	"runtime"
)

/**
 * Service中的一些公共方法
 */
type ServiceHandler struct {
	Orm *gorm.DB
	Ctx context.Context
}

/**
 * 异常状态码
 */
func (this ServiceHandler) Error(code int32, message string, err error) error {
	traceInfo := this.TraceInfo()
	return errors.InternalServerError(fmt.Sprint(code), "%s%s%s%s%s", message, errs.ERR_DS, err.Error(), errs.ERR_DS, traceInfo)
}

func (this ServiceHandler) TraceInfo() string {
	pc, file, line, _ := runtime.Caller(2)
	pcName := runtime.FuncForPC(pc).Name() //获取函数名
	traceInfo := fmt.Sprintf("in file:%s,at line:%d,%s", file, line, pcName)
	return traceInfo
}
