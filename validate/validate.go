package validate

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	errs "micrease-core/errs"
	"reflect"
)

//如果绑定参数失败，则中断请求并抛出错误返回给客户端
func BindWithPanic(ctx *gin.Context, obj interface{}) error {
	return bind(ctx, obj, true)
}

//如果绑定失败，则返回错误
func Bind(ctx *gin.Context, obj interface{}) error {
	return bind(ctx, obj, false)
}

func bind(ctx *gin.Context, obj interface{}, isPanic bool) error {
	//绑定参数
	err := ctx.Bind(obj)
	if err != nil {
		//解析错误信息
		value := reflect.TypeOf(obj)
		if validErrs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validErrs {
				if f, exist := value.Elem().FieldByName(e.Field()); exist {
					if value, ok := f.Tag.Lookup("tips"); ok {
						paramError := fmt.Errorf("%s", value)
						if isPanic {
							errs.PanicIfParamError(paramError)
						}
						return paramError
					} else {
						paramError := fmt.Errorf("%s", e.Value())
						if isPanic {
							errs.PanicIfParamError(paramError)
						}
						return paramError
					}
				}
			}
		}
	}
	return nil
}
