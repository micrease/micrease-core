package errs

import (
	"errors"
	"fmt"
	micro_errors "github.com/micro/go-micro/v2/errors"
	"runtime"
)

//对异常处理进行封装
//grpc的error信息解析
const (
	//error信息中的分隔符
	ERR_DS = "##"
)

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	TraceId string      `json:"traceId"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}

const (
	//成功
	StatusSuccess = 200
	//参数错误4000+
	StatusParamError = 4000
	//服务器错误5000+
	StatusServerError = 5000
)

func Err(code int32, message string, err error) error {
	traceInfo := traceInfo()
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	detail := fmt.Sprintf("%s%s%s%s%s", message, ERR_DS, errMsg, ERR_DS, traceInfo)
	//e := micro_errors.New("", detail, code)
	e := micro_errors.InternalServerError(fmt.Sprint(code), detail)
	return e
}

//参数只接受bool和error类型
//如果real为true或error!=nil时panic
func PanicIf(real interface{}, code int32, message string) {
	switch real.(type) {
	case error:
		e := real.(error)
		if e != nil {
			panic(Err(code, message, e))
		}
	case bool:
		b := real.(bool)
		if b {
			panic(Err(code, message, nil))
		}
	}
}

//可以使用如下方式抛出异常
//panic("查询商品列表失败")
//panic(err)
//PanicIf(err, 5001, "查询商品列表失败")
func Recover(err *error) {
	if e := recover(); e != nil {
		//断言比反射要快很多
		me, ok := e.(*micro_errors.Error)
		if ok {
			*err = fmt.Errorf("%v", me)
			return
		}

		er, ok := e.(error)
		if ok {
			*err = Err(StatusServerError, er.Error(), er)
			return
		}

		str, ok := e.(string)
		if ok {
			*err = Err(StatusServerError, str, errors.New(str))
			return
		}
		*err = Err(StatusServerError, "服务器错误", errors.New("unkonw error"))
	}
}

func traceInfo() string {
	pc, file, line, _ := runtime.Caller(3)
	pcName := runtime.FuncForPC(pc).Name() //获取函数名
	traceInfo := fmt.Sprintf("in file:%s,at line:%d,%s", file, line, pcName)
	return traceInfo
}
