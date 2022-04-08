package validate

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/micrease/micrease-core/context"
	"github.com/micrease/micrease-core/errs"
	"reflect"
)

//如果绑定参数失败，则中断请求并抛出错误返回给客户端
func BindWithPanic(ctx *context.Context, obj interface{}) error {
	return bind(ctx, obj, true)
}

//如果绑定失败，则返回错误
func Bind(ctx *context.Context, obj interface{}) error {
	return bind(ctx, obj, false)
}

func bind(ctx *context.Context, obj interface{}, isPanic bool) error {
	//绑定参数
	err := ctx.GinCtx.Bind(obj)
	if err != nil {
		//解析错误信息
		value := reflect.TypeOf(obj)
		if validErrs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validErrs {
				if f, exist := value.Elem().FieldByName(e.Field()); exist {
					var paramError error
					if value, ok := f.Tag.Lookup("tips"); ok {
						paramError = fmt.Errorf("%s", value)
					} else {
						paramError = fmt.Errorf("%s", e.Value())
					}

					if isPanic {
						er := errs.Err(errs.StatusParamError, paramError.Error(), paramError)
						panic(er)
					}
					return paramError
				}
			}
		}
	}
	return nil
}
