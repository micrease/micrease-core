package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micrease/micrease-core/errs"
	"github.com/micrease/micrease-core/trace"
	micro_errors "github.com/micro/go-micro/v2/errors"
	"net/http"
	"strconv"
	"strings"
)

func Recover(isDebug bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if c.IsAborted() {
					c.Status(errs.StatusSuccess)
				}
				switch errStr := err.(type) {
				case *micro_errors.Error:
					str := fmt.Sprintf("%s%s%s", errStr.Id, errs.ERR_DS, errStr.Detail)
					me := GetError(str, isDebug)
					me.TraceId = c.Request.Header.Get(trace.TrafficKey)
					c.JSON(http.StatusOK, me)
				case string:
					//格式如:5001#message
					me := GetError(errStr, isDebug)
					me.TraceId = c.Request.Header.Get(trace.TrafficKey)
					c.JSON(http.StatusOK, me)
				default:
					panic(err)
				}
			}
		}()
		c.Next()
	}
}

func GetError(str string, isDebug bool) errs.Error {
	me := errs.Error{}
	me.Code = errs.StatusServerError
	p := strings.Split(str, errs.ERR_DS)
	l := len(p)
	if l >= 2 {
		//501#message
		statusCode, e := strconv.Atoi(p[0])
		if e != nil {
			me.Error = e.Error()
			me.Message = e.Error()
			return me
		}
		me.Code = statusCode
		me.Message = p[1]
		//debug模式显示错误信息
		if isDebug {
			//Error
			if l >= 3 {
				me.Error = p[2]
			}

			//traceInfo
			if l >= 4 {
				me.Error = me.Error + "," + p[3]
			}
		}
	} else {
		me.Message = str
	}
	return me
}
