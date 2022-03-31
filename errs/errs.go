package errs

import (
	"fmt"
	"github.com/micro/go-micro/v2/errors"
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
	StatusCodeSuccess = 200
	//参数错误4000+
	StatusCodeParamError = 4000
	//服务器错误5000+
	StatusCodeServerError = 5000
)

func NewError() *Error {
	return &Error{}
}

//RPC返回的error, error.Error信息中是一个json字符串结构，因此和常规error不一样的
func PanicIfRpcError(err error) {
	if err == nil {
		return
	}
	rpcError := errors.Parse(err.Error())
	msg := fmt.Sprintf("%s%s%s", rpcError.Id, ERR_DS, rpcError.Detail)
	panic(msg)
}

func PanicIfParamError(err error) {
	PanicIfError(err, StatusCodeParamError, err.Error())
}

func PanicIfServerError(err error) {
	PanicIfError(err, StatusCodeServerError, err.Error())
}

//如果错误抛出信息
func PanicIfError(err error, code int, message string) {
	if err == nil {
		return
	}
	msg := fmt.Sprintf("%d%s%s", code, ERR_DS, message)
	panic(msg)
}

//如果条件不成立抛出信息
func PanicIfFalse(ok bool, code int, message string) {
	if ok {
		return
	}
	msg := fmt.Sprintf("%d%s%s", code, ERR_DS, message)
	panic(msg)
}

//抛出信息
func PanicMessage(code int, message string) {
	msg := fmt.Sprintf("%d%s%s", code, ERR_DS, message)
	panic(msg)
}
